package fixedcarrier

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"
	"io"

	"github.com/LuigiD5555/origami/internal/temporal"
)

const TemporalCarrierProfileID = "origami.temporal-carrier.r0.profile-1"

var temporalMagic = [8]byte{'O', 'T', 'C', 'R', '0', 0, 0, 1}

// TemporalCarrierDecoded reports the exact Temporal Program recovered from the
// image plus the semantic visual coverage rendered in the upper plane.
type TemporalCarrierDecoded struct {
	Schema            string           `json:"schema"`
	Profile           string           `json:"profile"`
	Program           temporal.Program `json:"program"`
	ProgramSHA256     string           `json:"program_sha256"`
	CompressedBytes   int              `json:"compressed_bytes"`
	PayloadBytes      int              `json:"payload_bytes"`
	VisibleCellCount  int              `json:"visible_cell_count"`
	VisibleRuleCount  int              `json:"visible_rule_count"`
	SemanticTopLevel  bool             `json:"semantic_top_level"`
	ExactProgramInPNG bool             `json:"exact_program_in_png"`
}

func RenderTemporalCarrier(program temporal.Program) ([]byte, TemporalCarrierDecoded, error) {
	if err := program.Validate(); err != nil {
		return nil, TemporalCarrierDecoded{}, err
	}
	canonical, err := json.Marshal(program)
	if err != nil {
		return nil, TemporalCarrierDecoded{}, err
	}
	digest := sha256.Sum256(canonical)
	compressed, err := zlibBest(canonical)
	if err != nil {
		return nil, TemporalCarrierDecoded{}, err
	}
	const headerBytes = 46
	const crcBytes = 4
	capacity := RecordBytes - headerBytes - crcBytes
	if len(compressed) > capacity {
		return nil, TemporalCarrierDecoded{}, fmt.Errorf("temporal program does not fit self-contained R0 payload: compressed=%d capacity=%d", len(compressed), capacity)
	}
	rec := make([]byte, RecordBytes)
	copy(rec[:8], temporalMagic[:])
	binary.BigEndian.PutUint32(rec[8:12], uint32(len(canonical)))
	binary.BigEndian.PutUint16(rec[12:14], uint16(len(compressed)))
	copy(rec[14:46], digest[:])
	copy(rec[46:46+len(compressed)], compressed)
	binary.BigEndian.PutUint32(rec[508:512], crc32.ChecksumIEEE(rec[:508]))

	img := image.NewGray(image.Rect(0, 0, Width, Height))
	for i := range img.Pix { img.Pix[i] = 0xff }
	drawFrame(img)
	drawTemporalBoot(img)
	visibleCells, visibleRules := drawTemporalSemanticPlane(img, program)
	drawTemporalExactPlane(img, rec)

	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&raw, img); err != nil {
		return nil, TemporalCarrierDecoded{}, err
	}
	if raw.Len() > FixedPNGBytes-12 {
		return nil, TemporalCarrierDecoded{}, fmt.Errorf("temporal carrier raw PNG exceeds frozen envelope: %d > %d", raw.Len(), FixedPNGBytes-12)
	}
	padded, err := padPNG(raw.Bytes(), FixedPNGBytes)
	if err != nil { return nil, TemporalCarrierDecoded{}, err }
	decoded := TemporalCarrierDecoded{
		Schema: "origami.temporal-carrier.r0", Profile: TemporalCarrierProfileID, Program: program,
		ProgramSHA256: fmt.Sprintf("%x", digest[:]), CompressedBytes: len(compressed), PayloadBytes: headerBytes + len(compressed) + crcBytes,
		VisibleCellCount: visibleCells, VisibleRuleCount: visibleRules, SemanticTopLevel: true, ExactProgramInPNG: true,
	}
	return padded, decoded, nil
}

func DecodeTemporalCarrierPNG(data []byte) (TemporalCarrierDecoded, error) {
	if len(data) > MaxPNGBytes { return TemporalCarrierDecoded{}, fmt.Errorf("carrier exceeds hard limit: %d", len(data)) }
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil { return TemporalCarrierDecoded{}, err }
	if img.Bounds().Dx() != Width || img.Bounds().Dy() != Height {
		return TemporalCarrierDecoded{}, fmt.Errorf("unexpected carrier dimensions %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
	rec := readTemporalRecord(img)
	if !bytes.Equal(rec[:8], temporalMagic[:]) { return TemporalCarrierDecoded{}, fmt.Errorf("temporal carrier magic mismatch") }
	if binary.BigEndian.Uint32(rec[508:512]) != crc32.ChecksumIEEE(rec[:508]) { return TemporalCarrierDecoded{}, fmt.Errorf("temporal carrier crc mismatch") }
	plainLen := int(binary.BigEndian.Uint32(rec[8:12]))
	compressedLen := int(binary.BigEndian.Uint16(rec[12:14]))
	if compressedLen < 0 || 46+compressedLen > 508 { return TemporalCarrierDecoded{}, fmt.Errorf("invalid compressed length") }
	plain, err := zlibDecode(rec[46 : 46+compressedLen])
	if err != nil { return TemporalCarrierDecoded{}, err }
	if len(plain) != plainLen { return TemporalCarrierDecoded{}, fmt.Errorf("program length mismatch: %d != %d", len(plain), plainLen) }
	digest := sha256.Sum256(plain)
	if !bytes.Equal(rec[14:46], digest[:]) { return TemporalCarrierDecoded{}, fmt.Errorf("program sha256 mismatch") }
	var program temporal.Program
	if err := json.Unmarshal(plain, &program); err != nil { return TemporalCarrierDecoded{}, err }
	if err := program.Validate(); err != nil { return TemporalCarrierDecoded{}, err }
	visibleCells := len(program.Automaton.Cells); if visibleCells > 8 { visibleCells = 8 }
	visibleRules := len(program.Automaton.Rules); if visibleRules > 8 { visibleRules = 8 }
	return TemporalCarrierDecoded{
		Schema: "origami.temporal-carrier.r0", Profile: TemporalCarrierProfileID, Program: program,
		ProgramSHA256: fmt.Sprintf("%x", digest[:]), CompressedBytes: compressedLen, PayloadBytes: 46 + compressedLen + 4,
		VisibleCellCount: visibleCells, VisibleRuleCount: visibleRules, SemanticTopLevel: true, ExactProgramInPNG: true,
	}, nil
}

func zlibBest(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
	if err != nil { return nil, err }
	if _, err := w.Write(in); err != nil { return nil, err }
	if err := w.Close(); err != nil { return nil, err }
	return b.Bytes(), nil
}

func zlibDecode(in []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(in)); if err != nil { return nil, err }
	defer r.Close()
	return io.ReadAll(r)
}

func drawTemporalBoot(img *image.Gray) {
	box(img, 14, 12, 612, 98, 0)
	drawText(img, 24, 20, 2, "ORIGAMI TEMPORAL CARRIER R0", 0)
	drawText(img, 24, 48, 1, "T1 ROSETTA: BOX=CELL ARROW=TRANSITION RING=CHECKPOINT X=TIME", 0)
	drawText(img, 24, 62, 1, "READ: T2 GRAPH > ST* | WRITE: ET* > GRAPH+TIMELINE", 0)
	drawText(img, 24, 76, 1, "SEMANTIC FILM != VIDEO | EXACT PROGRAM BELOW | FALSE EXACT 0", 0)
}

func drawTemporalSemanticPlane(img *image.Gray, p temporal.Program) (int, int) {
	drawText(img, 24, 120, 1, "T2 SEMANTIC-TEMPORAL SUPERGRAPH", 0)
	cells := p.Automaton.Cells
	if len(cells) > 8 { cells = cells[:8] }
	positions := map[string][2]int{}
	for i, c := range cells {
		x := 34 + (i%4)*145
		y := 148 + (i/4)*82
		box(img, x, y, 110, 48, 0)
		drawText(img, x+8, y+8, 1, "CELL "+shortLabel(c.ID, 10), 0)
		drawText(img, x+8, y+24, 1, shortLabel(c.InitialState, 12), 0)
		positions[c.ID] = [2]int{x+55, y+24}
	}
	visibleRules := 0
	for _, r := range p.Automaton.Rules {
		if visibleRules >= 8 { break }
		to, ok := positions[r.TargetCell]; if !ok { continue }
		if len(r.Requires) == 0 {
			drawSmallRing(img, to[0], to[1])
			visibleRules++
			continue
		}
		from, ok := positions[r.Requires[0].CellID]; if !ok { continue }
		drawSimpleLine(img, from[0], from[1], to[0], to[1])
		drawArrowTip(img, from[0], from[1], to[0], to[1])
		visibleRules++
	}

	y := 326
	drawText(img, 24, y, 1, "TIMELINE / CHECKPOINT ROUTE", 0)
	x0, x1 := 38, 602
	drawSimpleLine(img, x0, y+28, x1, y+28)
	steps := p.MaxSteps; if steps < 1 { steps = 1 }; if steps > 32 { steps = 32 }
	for s := 0; s <= steps; s++ {
		x := x0 + (x1-x0)*s/steps
		h := 8
		if p.CheckpointEvery > 0 && s%p.CheckpointEvery == 0 { h = 18; drawSmallRing(img, x, y+28) }
		drawSimpleLine(img, x, y+28-h/2, x, y+28+h/2)
	}
	drawText(img, 38, y+46, 1, "t0", 0)
	drawText(img, 560, y+46, 1, "tN", 0)
	return len(cells), visibleRules
}

func drawTemporalExactPlane(img *image.Gray, rec []byte) {
	drawText(img, 24, 398, 1, "EXACT PROGRAM PAYLOAD: ZLIB JSON + SHA256 + CRC", 0)
	// 64x64 bits, one 3x3 cell per bit. The program is self-contained here,
	// while semantic questions are expected to use the graph/timeline above.
	const x0 = 224
	const y0 = 420
	box(img, x0-5, y0-5, GridBits*Cell+10, GridBits*Cell+10, 0)
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			b := (rec[bit/8] >> uint(7-(bit%8))) & 1
			if b == 1 { fill(img, x0+x*Cell, y0+y*Cell, Cell, Cell, 0) }
			bit++
		}
	}
}

func readTemporalRecord(img image.Image) []byte {
	rec := make([]byte, RecordBytes)
	const x0 = 224
	const y0 = 420
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			r, g, b, _ := img.At(x0+x*Cell+Cell/2, y0+y*Cell+Cell/2).RGBA()
			if (r+g+b)/3 < 0x8000 { rec[bit/8] |= 1 << uint(7-(bit%8)) }
			bit++
		}
	}
	return rec
}

func shortLabel(s string, n int) string { if len(s) <= n { return s }; return s[:n] }

func drawSimpleLine(img *image.Gray, x0, y0, x1, y1 int) {
	dx := absInt(x1-x0); sx := -1; if x0 < x1 { sx = 1 }
	dy := -absInt(y1-y0); sy := -1; if y0 < y1 { sy = 1 }
	err := dx + dy
	for {
		if image.Pt(x0,y0).In(img.Bounds()) { fill(img, x0, y0, 1, 1, 0) }
		if x0 == x1 && y0 == y1 { break }
		e2 := 2*err
		if e2 >= dy { err += dy; x0 += sx }
		if e2 <= dx { err += dx; y0 += sy }
	}
}

func drawArrowTip(img *image.Gray, x0, y0, x1, y1 int) {
	dx, dy := x1-x0, y1-y0
	if absInt(dx) >= absInt(dy) {
		if dx >= 0 { drawSimpleLine(img,x1,y1,x1-6,y1-4); drawSimpleLine(img,x1,y1,x1-6,y1+4) } else { drawSimpleLine(img,x1,y1,x1+6,y1-4); drawSimpleLine(img,x1,y1,x1+6,y1+4) }
	} else {
		if dy >= 0 { drawSimpleLine(img,x1,y1,x1-4,y1-6); drawSimpleLine(img,x1,y1,x1+4,y1-6) } else { drawSimpleLine(img,x1,y1,x1-4,y1+6); drawSimpleLine(img,x1,y1,x1+4,y1+6) }
	}
}

func drawSmallRing(img *image.Gray, cx, cy int) {
	box(img, cx-4, cy-4, 9, 9, 0)
	fill(img, cx-2, cy-2, 5, 5, 0xff)
}

func absInt(v int) int { if v < 0 { return -v }; return v }
