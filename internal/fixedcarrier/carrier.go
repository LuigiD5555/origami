package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strings"
)

const (
	Schema          = "origami.fixed-carrier.r2"
	ToolProtocol    = "tlaloc.origami-tools.r2"
	AddressABI      = "ohf-address.r2"
	ProfileID       = "origami.fixed-carrier.r2.profile-2"
	LegacyProfileID = "origami.fixed-carrier.r2.profile-1"
	Width           = 640
	Height          = 640
	GridBits        = 64
	RecordBytes     = GridBits * GridBits / 8
	Cell            = 3
	GridX           = 224
	GridY           = 316
	MaxPNGBytes     = 512000
	FixedPNGBytes   = 8192
)

var magic = [8]byte{'O', 'F', 'C', 'R', '2', 0, 0, 1}

var bootText = []string{
	"ORIGAMI FIXED CARRIER R2",
	"THIS IMAGE IS A COMPUTATIONAL CARRIER",
	"BOOT: TEXT > PROBE > ROSETTA > T2",
	"SEMANTIC: READ T2. DO NOT DECODE BINARY",
	"EXACT: QUERY EXPAND VERIFY IF TOOLS EXIST",
	"UNKNOWN > INVENTION. FALSE_EXACT=0",
}

var legacyBootText = []string{
	"ORIGAMI FIXED CARRIER R2",
	"THIS IMAGE IS A COMPUTATIONAL CARRIER",
	"BOOT: TEXT > PROBE > ROSETTA > INDEX > TOOLS",
	"TOOLS: BOOT QUERY EXPAND VERIFY",
	"ADDRESS: OHF://<CARRIER>/<SPACE>/<ID>",
	"OCR OPTIONAL. NEVER INVENT.",
}

type Metadata struct {
	CarrierID      string `json:"carrier_id"`
	StoreRoot      string `json:"store_root_sha256"`
	SourceSHA256   string `json:"source_sha256"`
	PageCount      uint32 `json:"page_count"`
	BlockCount     uint32 `json:"block_count"`
	DocumentCount  uint32 `json:"document_count,omitempty"`
	ObjectCount    uint32 `json:"object_count,omitempty"`
	GraphSignature []byte `json:"graph_signature,omitempty"`
	GraphSketch    []byte `json:"graph_sketch,omitempty"`
	Flags          uint16 `json:"flags,omitempty"`
}

type Decoded struct {
	Schema       string `json:"schema"`
	Profile      string `json:"profile"`
	ToolProtocol string `json:"tool_protocol"`
	AddressABI   string `json:"address_abi"`
	Metadata
	CarrierDigest string   `json:"carrier_digest_sha256"`
	VisualProbe   string   `json:"visual_probe"`
	BootText      []string `json:"boot_text"`
	Zones         []Zone   `json:"zones"`
}

func Render(meta Metadata) ([]byte, Decoded, error) {
	rec, decoded, err := encodeRecord(meta)
	if err != nil {
		return nil, Decoded{}, err
	}
	img := image.NewGray(image.Rect(0, 0, Width, Height))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	drawFrame(img)
	drawT0(img)
	drawT1(img, decoded.VisualProbe)
	drawT2(img, decoded.GraphSignature)
	drawT3(img, rec)
	drawVerify(img, decoded.VisualProbe)
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		return nil, Decoded{}, err
	}
	if buf.Len() > FixedPNGBytes-12 {
		return nil, Decoded{}, fmt.Errorf("R2 raw PNG exceeds frozen envelope: %d > %d", buf.Len(), FixedPNGBytes-12)
	}
	padded, err := padPNG(buf.Bytes(), FixedPNGBytes)
	if err != nil {
		return nil, Decoded{}, err
	}
	if len(padded) > MaxPNGBytes {
		return nil, Decoded{}, fmt.Errorf("fixed carrier exceeds hard limit: %d > %d", len(padded), MaxPNGBytes)
	}
	return padded, decoded, nil
}

func DecodePNG(data []byte) (Decoded, error) {
	if len(data) > MaxPNGBytes {
		return Decoded{}, fmt.Errorf("carrier exceeds hard limit: %d", len(data))
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return Decoded{}, err
	}
	b := img.Bounds()
	if b.Dx() != Width || b.Dy() != Height {
		return Decoded{}, fmt.Errorf("unexpected carrier dimensions %dx%d", b.Dx(), b.Dy())
	}
	rec := make([]byte, RecordBytes)
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			gx := GridX + x*Cell + Cell/2
			gy := GridY + y*Cell + Cell/2
			r, g, bb, _ := img.At(gx, gy).RGBA()
			if (r+g+bb)/3 < 0x8000 {
				rec[bit/8] |= 1 << uint(7-(bit%8))
			}
			bit++
		}
	}
	return decodeRecord(rec)
}

func VisualProbeFromPNG(data []byte) (top, bottom string, err error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return "", "", err
	}
	if img.Bounds().Dx() != Width || img.Bounds().Dy() != Height {
		return "", "", fmt.Errorf("unexpected carrier dimensions")
	}
	return readProbeRow(img, 154), readProbeRow(img, 566), nil
}

func WritePNG(path string, meta Metadata) (Decoded, int, error) {
	data, dec, err := Render(meta)
	if err != nil {
		return Decoded{}, 0, err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return Decoded{}, 0, err
	}
	return dec, len(data), nil
}

func encodeRecord(meta Metadata) ([]byte, Decoded, error) {
	cid, err := parseHashish(meta.CarrierID)
	if err != nil {
		return nil, Decoded{}, fmt.Errorf("carrier_id: %w", err)
	}
	root, err := parseHashish(meta.StoreRoot)
	if err != nil {
		return nil, Decoded{}, fmt.Errorf("store_root: %w", err)
	}
	src, err := parseHashish(meta.SourceSHA256)
	if err != nil {
		return nil, Decoded{}, fmt.Errorf("source_sha256: %w", err)
	}
	rec := make([]byte, RecordBytes)
	copy(rec[0:8], magic[:])
	binary.BigEndian.PutUint16(rec[8:10], meta.Flags)
	binary.BigEndian.PutUint32(rec[10:14], meta.PageCount)
	binary.BigEndian.PutUint32(rec[14:18], meta.BlockCount)
	binary.BigEndian.PutUint32(rec[18:22], max1(meta.DocumentCount))
	objects := meta.ObjectCount
	if objects == 0 {
		objects = meta.PageCount + meta.BlockCount + max1(meta.DocumentCount)
	}
	binary.BigEndian.PutUint32(rec[22:26], objects)
	copy(rec[26:58], cid)
	copy(rec[58:90], root)
	copy(rec[90:122], src)
	tool := sha256.Sum256([]byte(ToolProtocol))
	copy(rec[122:138], tool[:16])
	profile := sha256.Sum256([]byte(ProfileID))
	copy(rec[138:154], profile[:16])
	boot := sha256.Sum256([]byte(joinBootText(bootText)))
	copy(rec[154:170], boot[:16])
	sig := meta.GraphSignature
	if len(sig) == 0 {
		sig = meta.GraphSketch
	}
	if len(sig) > 256 {
		sig = sig[:256]
	}
	copy(rec[170:426], sig)
	fillRedundancy(rec[426:508], rec[:426])
	crc := crc32.ChecksumIEEE(rec[:508])
	binary.BigEndian.PutUint32(rec[508:512], crc)
	digest := sha256.Sum256(rec)
	return rec, Decoded{Schema: Schema, Profile: ProfileID, ToolProtocol: ToolProtocol, AddressABI: AddressABI, Metadata: Metadata{
		CarrierID: hex.EncodeToString(cid), StoreRoot: hex.EncodeToString(root), SourceSHA256: hex.EncodeToString(src), PageCount: meta.PageCount, BlockCount: meta.BlockCount, DocumentCount: max1(meta.DocumentCount), ObjectCount: objects, GraphSignature: append([]byte(nil), sig...), Flags: meta.Flags,
	}, CarrierDigest: hex.EncodeToString(digest[:]), VisualProbe: probeBits(digest[0]), BootText: append([]string(nil), bootText...), Zones: append([]Zone(nil), profileR2.Zones...)}, nil
}

func decodeRecord(rec []byte) (Decoded, error) {
	if len(rec) != RecordBytes {
		return Decoded{}, fmt.Errorf("invalid record length")
	}
	if !bytes.Equal(rec[:8], magic[:]) {
		return Decoded{}, fmt.Errorf("fixed carrier R2 magic mismatch")
	}
	if binary.BigEndian.Uint32(rec[508:512]) != crc32.ChecksumIEEE(rec[:508]) {
		return Decoded{}, fmt.Errorf("carrier crc mismatch")
	}
	tool := sha256.Sum256([]byte(ToolProtocol))
	if !bytes.Equal(rec[122:138], tool[:16]) {
		return Decoded{}, fmt.Errorf("tool protocol mismatch")
	}
	profileID, profileBoot, err := detectProfile(rec[138:154])
	if err != nil {
		return Decoded{}, err
	}
	boot := sha256.Sum256([]byte(joinBootText(profileBoot)))
	if !bytes.Equal(rec[154:170], boot[:16]) {
		return Decoded{}, fmt.Errorf("T0 boot-text binding mismatch")
	}
	wantRed := make([]byte, 82)
	fillRedundancy(wantRed, rec[:426])
	if !bytes.Equal(wantRed, rec[426:508]) {
		return Decoded{}, fmt.Errorf("carrier redundancy mismatch")
	}
	digest := sha256.Sum256(rec)
	sig := make([]byte, 256)
	copy(sig, rec[170:426])
	return Decoded{Schema: Schema, Profile: profileID, ToolProtocol: ToolProtocol, AddressABI: AddressABI, Metadata: Metadata{
		CarrierID: hex.EncodeToString(rec[26:58]), StoreRoot: hex.EncodeToString(rec[58:90]), SourceSHA256: hex.EncodeToString(rec[90:122]), PageCount: binary.BigEndian.Uint32(rec[10:14]), BlockCount: binary.BigEndian.Uint32(rec[14:18]), DocumentCount: binary.BigEndian.Uint32(rec[18:22]), ObjectCount: binary.BigEndian.Uint32(rec[22:26]), GraphSignature: sig, Flags: binary.BigEndian.Uint16(rec[8:10]),
	}, CarrierDigest: hex.EncodeToString(digest[:]), VisualProbe: probeBits(digest[0]), BootText: append([]string(nil), profileBoot...), Zones: append([]Zone(nil), profileR2.Zones...)}, nil
}

func detectProfile(digest []byte) (string, []string, error) {
	current := sha256.Sum256([]byte(ProfileID))
	if bytes.Equal(digest, current[:16]) {
		return ProfileID, bootText, nil
	}
	legacy := sha256.Sum256([]byte(LegacyProfileID))
	if bytes.Equal(digest, legacy[:16]) {
		return LegacyProfileID, legacyBootText, nil
	}
	return "", nil, fmt.Errorf("profile mismatch")
}

func drawT0(img *image.Gray) {
	box(img, 14, 12, 612, 112, 0)
	y := 20
	for i, s := range bootText {
		scale := 2
		if i == 0 {
			scale = 2
		}
		x := 24
		if w := textWidth(s, scale); w > Width-48 {
			scale = 1
		}
		drawText(img, x, y, scale, s, 0)
		y += 16
	}
}

func drawT1(img *image.Gray, probe string) {
	drawText(img, 24, 132, 1, "T1 ROSETTA: NODE=ENTRY LINK=RELATION BOX=SCOPE TEXT=SEMANTIC LABEL", 0)
	drawProbeRow(img, 154, probe)
	drawText(img, 24, 190, 1, "READ BOTH PROBE ROWS. THEY MUST AGREE BEFORE SEMANTIC NAVIGATION.", 0)
}

func drawT2(img *image.Gray, signature []byte) {
	box(img, 14, 214, 612, 82, 0)
	drawText(img, 24, 220, 1, "T2 SEMANTIC SUPERINDEX", 0)
	entries := semanticIndexFromSignature(signature)
	if len(entries) == 0 {
		drawText(img, 42, 238, 1, "INDEX UNKNOWN: NO SEMANTIC HINT", 0)
		drawText(img, 42, 254, 1, "DO NOT DECODE T3 TO INVENT A SEMANTIC ANSWER", 0)
		drawText(img, 24, 280, 1, "T3 IS CONTROL/EXACT; SEMANTIC NAVIGATION STAYS SEPARATE", 0)
		return
	}
	for i, label := range entries {
		if i >= 4 {
			break
		}
		y := 236 + i*12
		fill(img, 26, y+1, 7, 7, 0)
		drawText(img, 42, y, 1, fmt.Sprintf("%d %s", i+1, label), 0)
		if i > 0 {
			line(img, 29, y-3, 29, y+1, 0)
		}
	}
	drawText(img, 310, 280, 1, "SEMANTIC FIRST | EXACT OPTIONAL", 0)
}

func semanticIndexFromSignature(signature []byte) []string {
	type hint struct {
		Index   []string `json:"i"`
		Roots   []string `json:"r"`
		Groups  []string `json:"g"`
		Classes []string `json:"k"`
	}
	trimmed := bytes.TrimRight(signature, "\x00")
	if len(trimmed) == 0 {
		return nil
	}
	var h hint
	if err := json.Unmarshal(trimmed, &h); err != nil {
		return nil
	}
	for _, values := range [][]string{h.Index, h.Roots, h.Groups, h.Classes} {
		if len(values) == 0 {
			continue
		}
		out := make([]string, 0, 4)
		for _, value := range values {
			value = sanitizeIndexLabel(value, 30)
			if value == "" {
				continue
			}
			out = append(out, value)
			if len(out) == 4 {
				break
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return nil
}

func sanitizeIndexLabel(value string, width int) string {
	value = strings.ToUpper(strings.Join(strings.Fields(value), " "))
	var b strings.Builder
	for _, r := range value {
		ok := (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '_' || r == '/' || r == ':' || r == '.'
		if ok {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
		if width > 0 && b.Len() >= width {
			break
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func drawT3(img *image.Gray, rec []byte) {
	drawText(img, 24, 302, 1, "T3 CONTROL / EXACT RECORD", 0)
	drawText(img, 28, 320, 1, "MACRO", 0)
	drawNode(img, 92, 344, 0)
	drawNode(img, 58, 382, 1)
	drawNode(img, 126, 382, 2)
	line(img, 92, 353, 58, 373, 0)
	line(img, 92, 353, 126, 373, 0)
	drawNode(img, 42, 422, 3)
	drawNode(img, 74, 422, 3)
	drawNode(img, 110, 422, 3)
	drawNode(img, 142, 422, 3)
	line(img, 58, 391, 42, 413, 0)
	line(img, 58, 391, 74, 413, 0)
	line(img, 126, 391, 110, 413, 0)
	line(img, 126, 391, 142, 413, 0)
	drawText(img, 30, 446, 1, "MESO GRAPH", 0)
	drawText(img, 474, 320, 1, "MICRO", 0)
	for i := 0; i < 4; i++ {
		box(img, 462+i*8, 346+i*12, 118-i*16, 44, 0)
	}
	drawText(img, 474, 404, 1, "ADDRESS", 0)
	drawText(img, 458, 418, 1, "OHF://.../ID", 0)
	drawText(img, 462, 450, 1, "CID", 0)
	drawText(img, 462, 464, 1, "MERKLE", 0)
	drawText(img, 462, 478, 1, "VERIFY", 0)
	box(img, GridX-5, GridY-5, GridBits*Cell+10, GridBits*Cell+10, 0)
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			b := (rec[bit/8] >> uint(7-(bit%8))) & 1
			if b == 1 {
				fill(img, GridX+x*Cell, GridY+y*Cell, Cell, Cell, 0)
			}
			bit++
		}
	}
}

func drawVerify(img *image.Gray, probe string) {
	drawProbeRow(img, 566, probe)
	drawText(img, 24, 610, 1, "VERIFY: ROOT BOUND | FALSE EXACT 0 | UNKNOWN > INVENT", 0)
}

func drawProbeRow(img *image.Gray, y int, bits string) {
	start := 144
	for i := 0; i < 8; i++ {
		x := start + i*44
		box(img, x, y, 36, 28, 0)
		fill(img, x+4, y+4, 28, 20, 255)
		if i < len(bits) && bits[i] == '1' {
			fill(img, x+8, y+7, 20, 14, 0)
		}
	}
}

func readProbeRow(img image.Image, y int) string {
	out := make([]byte, 8)
	start := 144
	for i := 0; i < 8; i++ {
		x := start + i*44 + 18
		yy := y + 14
		r, g, b, _ := img.At(x, yy).RGBA()
		if (r+g+b)/3 < 0x8000 {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return string(out)
}

func drawNode(img *image.Gray, cx, cy, kind int) {
	switch kind % 4 {
	case 0:
		fill(img, cx-8, cy-8, 16, 16, 0)
		fill(img, cx-4, cy-4, 8, 8, 255)
	case 1:
		box(img, cx-9, cy-7, 18, 14, 0)
	case 2:
		diamond(img, cx, cy, 9, 0)
	case 3:
		box(img, cx-7, cy-7, 14, 14, 0)
		line(img, cx-7, cy, cx+7, cy, 0)
	}
}

func drawFrame(img *image.Gray) {
	box(img, 2, 2, Width-4, Height-4, 0)
	for _, p := range [][2]int{{8, 8}, {Width - 46, 8}, {8, Height - 46}, {Width - 46, Height - 46}} {
		fill(img, p[0], p[1], 38, 38, 0)
		fill(img, p[0]+7, p[1]+7, 24, 24, 255)
		fill(img, p[0]+13, p[1]+13, 12, 12, 0)
	}
}

func fillRedundancy(dst, seed []byte) {
	for i := range dst {
		h := sha256.Sum256(append(append([]byte(nil), seed...), byte(i), byte(i>>8)))
		dst[i] = h[0]
	}
}

func probeBits(v byte) string {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		if v&(1<<uint(7-i)) != 0 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}

func parseHashish(s string) ([]byte, error) {
	s = bytesToString(bytes.TrimSpace([]byte(s)))
	if len(s) == 64 {
		if b, err := hex.DecodeString(s); err == nil {
			return b, nil
		}
	}
	h := sha256.Sum256([]byte(s))
	return h[:], nil
}

func bytesToString(b []byte) string { return string(b) }

func max1(v uint32) uint32 {
	if v == 0 {
		return 1
	}
	return v
}

func joinBootText(lines []string) string {
	var b bytes.Buffer
	for _, s := range lines {
		b.WriteString(s)
		b.WriteByte('\n')
	}
	return b.String()
}

func padPNG(data []byte, target int) ([]byte, error) {
	if len(data) > target {
		return nil, fmt.Errorf("PNG larger than frozen target: %d > %d", len(data), target)
	}
	iend := bytes.LastIndex(data, []byte{'I', 'E', 'N', 'D'})
	if iend < 4 {
		return nil, fmt.Errorf("IEND chunk missing")
	}
	chunkStart := iend - 4
	gap := target - len(data)
	if gap == 0 {
		return data, nil
	}
	if gap < 12 {
		return nil, fmt.Errorf("padding gap too small")
	}
	payload := make([]byte, gap-12)
	chunk := pngChunk("oPAD", payload)
	out := make([]byte, 0, target)
	out = append(out, data[:chunkStart]...)
	out = append(out, chunk...)
	out = append(out, data[chunkStart:]...)
	if len(out) != target {
		return nil, fmt.Errorf("padding failed: %d", len(out))
	}
	return out, nil
}

func pngChunk(kind string, payload []byte) []byte {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.BigEndian, uint32(len(payload)))
	b.WriteString(kind)
	b.Write(payload)
	crc := crc32.NewIEEE()
	_, _ = io.WriteString(crc, kind)
	_, _ = crc.Write(payload)
	_ = binary.Write(&b, binary.BigEndian, crc.Sum32())
	return b.Bytes()
}

func box(img *image.Gray, x, y, w, h int, v uint8) {
	line(img, x, y, x+w, y, v)
	line(img, x+w, y, x+w, y+h, v)
	line(img, x+w, y+h, x, y+h, v)
	line(img, x, y+h, x, y, v)
}

func fill(img *image.Gray, x, y, w, h int, v uint8) {
	c := color.Gray{Y: v}
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			if image.Pt(xx, yy).In(img.Bounds()) {
				img.SetGray(xx, yy, c)
			}
		}
	}
}

func line(img *image.Gray, x0, y0, x1, y1 int, v uint8) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		if image.Pt(x0, y0).In(img.Bounds()) {
			img.SetGray(x0, y0, color.Gray{Y: v})
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func diamond(img *image.Gray, cx, cy, r int, v uint8) {
	line(img, cx, cy-r, cx+r, cy, v)
	line(img, cx+r, cy, cx, cy+r, v)
	line(img, cx, cy+r, cx-r, cy, v)
	line(img, cx-r, cy, cx, cy-r, v)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
