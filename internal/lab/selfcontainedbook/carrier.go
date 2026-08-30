package selfcontainedbook

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

const (
	Schema             = "origami.self-contained-book.r0-lab"
	ProfileID          = "origami.self-contained-book.r0-lab.profile-2"
	HeaderBytes        = 144
	DefaultDataWidth   = 2304
	TargetPNGBytes     = 512000
	ControlPlaneHeight = 640
)

var magic = [8]byte{'O', 'S', 'C', 'B', 'R', '0', 0, 2}

type Header struct {
	SourceBytes        uint64 `json:"source_bytes"`
	CompressedBytes    uint64 `json:"compressed_bytes"`
	MasterPromptBytes  uint32 `json:"master_prompt_bytes"`
	PageCount          uint32 `json:"page_count"`
	BlockCount         uint32 `json:"block_count"`
	SourceSHA256       string `json:"source_sha256"`
	CompressedSHA256   string `json:"compressed_sha256"`
	MasterPromptSHA256 string `json:"master_prompt_sha256"`
	Compression        string `json:"compression"`
}

type Rendered struct {
	Schema                 string `json:"schema"`
	Profile                string `json:"profile"`
	PNG                    []byte `json:"-"`
	PNGBytes               int    `json:"png_bytes"`
	TargetPNGBytes         int    `json:"target_png_bytes"`
	OverTargetBytes        int    `json:"over_target_bytes"`
	PayloadBytes           int    `json:"payload_bytes"`
	MasterPromptBytes      int    `json:"master_prompt_bytes"`
	DataWidthBits          int    `json:"data_width_bits"`
	DataHeightBits         int    `json:"data_height_bits"`
	ControlX               int    `json:"control_x"`
	ControlY               int    `json:"control_y"`
	VisibleMasterBootstrap bool   `json:"visible_master_bootstrap"`
}

// Pack creates one exact in-image payload. The Master Prompt is carried as
// literal UTF-8 bytes inside the image, not merely referenced by hash.
// The current book body remains an explicit literal residual baseline.
func Pack(source, masterPrompt, compressed []byte, pageCount, blockCount uint32, compression string) ([]byte, Header, error) {
	if len(source) == 0 {
		return nil, Header{}, fmt.Errorf("source cannot be empty")
	}
	if len(masterPrompt) == 0 {
		return nil, Header{}, fmt.Errorf("Master Prompt cannot be empty")
	}
	if len(masterPrompt) > int(^uint32(0)) {
		return nil, Header{}, fmt.Errorf("Master Prompt is too large")
	}
	if len(compressed) == 0 {
		return nil, Header{}, fmt.Errorf("compressed residual cannot be empty")
	}
	if len(compression) == 0 || len(compression) > 8 {
		return nil, Header{}, fmt.Errorf("compression label must be 1..8 ASCII bytes")
	}

	sourceHash := sha256.Sum256(source)
	compressedHash := sha256.Sum256(compressed)
	promptHash := sha256.Sum256(masterPrompt)
	header := make([]byte, HeaderBytes)
	copy(header[0:8], magic[:])
	binary.BigEndian.PutUint64(header[8:16], uint64(len(source)))
	binary.BigEndian.PutUint64(header[16:24], uint64(len(compressed)))
	binary.BigEndian.PutUint32(header[24:28], uint32(len(masterPrompt)))
	binary.BigEndian.PutUint32(header[28:32], pageCount)
	binary.BigEndian.PutUint32(header[32:36], blockCount)
	copy(header[36:68], sourceHash[:])
	copy(header[68:100], compressedHash[:])
	copy(header[100:132], promptHash[:])
	copy(header[132:140], []byte(compression))
	binary.BigEndian.PutUint32(header[140:144], crc32.ChecksumIEEE(header[:140]))

	payload := make([]byte, 0, HeaderBytes+len(masterPrompt)+len(compressed))
	payload = append(payload, header...)
	payload = append(payload, masterPrompt...)
	payload = append(payload, compressed...)
	return payload, Header{
		SourceBytes:        uint64(len(source)),
		CompressedBytes:    uint64(len(compressed)),
		MasterPromptBytes:  uint32(len(masterPrompt)),
		PageCount:          pageCount,
		BlockCount:         blockCount,
		SourceSHA256:       hex.EncodeToString(sourceHash[:]),
		CompressedSHA256:   hex.EncodeToString(compressedHash[:]),
		MasterPromptSHA256: hex.EncodeToString(promptHash[:]),
		Compression:        compression,
	}, nil
}

func Unpack(payload []byte) (Header, []byte, []byte, error) {
	if len(payload) < HeaderBytes {
		return Header{}, nil, nil, fmt.Errorf("payload too short")
	}
	header := payload[:HeaderBytes]
	if !bytes.Equal(header[:8], magic[:]) {
		return Header{}, nil, nil, fmt.Errorf("self-contained book magic mismatch")
	}
	if binary.BigEndian.Uint32(header[140:144]) != crc32.ChecksumIEEE(header[:140]) {
		return Header{}, nil, nil, fmt.Errorf("self-contained book header crc mismatch")
	}

	promptBytes := binary.BigEndian.Uint32(header[24:28])
	compressedBytes := binary.BigEndian.Uint64(header[16:24])
	total := uint64(HeaderBytes) + uint64(promptBytes) + compressedBytes
	if total > uint64(len(payload)) {
		return Header{}, nil, nil, fmt.Errorf("declared prompt/residual exceeds payload")
	}
	promptStart := HeaderBytes
	promptEnd := promptStart + int(promptBytes)
	compressedEnd := promptEnd + int(compressedBytes)
	masterPrompt := append([]byte(nil), payload[promptStart:promptEnd]...)
	compressed := append([]byte(nil), payload[promptEnd:compressedEnd]...)

	gotPrompt := sha256.Sum256(masterPrompt)
	if !bytes.Equal(gotPrompt[:], header[100:132]) {
		return Header{}, nil, nil, fmt.Errorf("Master Prompt sha256 mismatch")
	}
	gotCompressed := sha256.Sum256(compressed)
	if !bytes.Equal(gotCompressed[:], header[68:100]) {
		return Header{}, nil, nil, fmt.Errorf("compressed residual sha256 mismatch")
	}

	return Header{
		SourceBytes:        binary.BigEndian.Uint64(header[8:16]),
		CompressedBytes:    compressedBytes,
		MasterPromptBytes:  promptBytes,
		PageCount:          binary.BigEndian.Uint32(header[28:32]),
		BlockCount:         binary.BigEndian.Uint32(header[32:36]),
		SourceSHA256:       hex.EncodeToString(header[36:68]),
		CompressedSHA256:   hex.EncodeToString(header[68:100]),
		MasterPromptSHA256: hex.EncodeToString(header[100:132]),
		Compression:        string(bytes.TrimRight(header[132:140], "\x00")),
	}, masterPrompt, compressed, nil
}

func Render(controlPNG, payload []byte, dataWidth int) (Rendered, error) {
	if dataWidth <= 0 {
		dataWidth = DefaultDataWidth
	}
	if dataWidth < 640 {
		return Rendered{}, fmt.Errorf("data width must be >= 640")
	}
	control, err := png.Decode(bytes.NewReader(controlPNG))
	if err != nil {
		return Rendered{}, fmt.Errorf("decode control plane: %w", err)
	}
	if control.Bounds().Dx() != 640 || control.Bounds().Dy() != 640 {
		return Rendered{}, fmt.Errorf("control plane must be 640x640")
	}
	header, masterPrompt, _, err := Unpack(payload)
	if err != nil {
		return Rendered{}, fmt.Errorf("inspect self-contained payload: %w", err)
	}

	bits := len(payload) * 8
	dataHeight := (bits + dataWidth - 1) / dataWidth
	canvas := image.NewGray(image.Rect(0, 0, dataWidth, ControlPlaneHeight+dataHeight))
	for i := range canvas.Pix {
		canvas.Pix[i] = 0xff
	}
	controlX := (dataWidth - 640) / 2
	draw.Draw(canvas, image.Rect(controlX, 0, controlX+640, 640), control, control.Bounds().Min, draw.Src)
	drawSelfBootstrap(canvas, controlX, header, masterPrompt)

	bit := 0
	for _, value := range payload {
		for shift := 7; shift >= 0; shift-- {
			x := bit % dataWidth
			y := ControlPlaneHeight + bit/dataWidth
			if (value>>uint(shift))&1 == 1 {
				canvas.SetGray(x, y, color.Gray{Y: 0})
			}
			bit++
		}
	}
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, canvas); err != nil {
		return Rendered{}, err
	}
	over := buf.Len() - TargetPNGBytes
	if over < 0 {
		over = 0
	}
	return Rendered{
		Schema:                 Schema,
		Profile:                ProfileID,
		PNG:                    buf.Bytes(),
		PNGBytes:               buf.Len(),
		TargetPNGBytes:         TargetPNGBytes,
		OverTargetBytes:        over,
		PayloadBytes:           len(payload),
		MasterPromptBytes:      len(masterPrompt),
		DataWidthBits:          dataWidth,
		DataHeightBits:         dataHeight,
		ControlX:               controlX,
		ControlY:               0,
		VisibleMasterBootstrap: true,
	}, nil
}

func DecodePNG(data []byte) ([]byte, image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	b := img.Bounds()
	if b.Dx() < 640 || b.Dy() <= ControlPlaneHeight {
		return nil, nil, fmt.Errorf("invalid self-contained carrier dimensions")
	}
	controlX := (b.Dx() - 640) / 2
	control := image.NewGray(image.Rect(0, 0, 640, 640))
	draw.Draw(control, control.Bounds(), img, image.Pt(b.Min.X+controlX, b.Min.Y), draw.Src)

	readBytes := func(count int) []byte {
		out := make([]byte, count)
		for bit := 0; bit < count*8; bit++ {
			x := bit % b.Dx()
			y := ControlPlaneHeight + bit/b.Dx()
			r, g, bb, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if (r+g+bb)/3 < 0x8000 {
				out[bit/8] |= 1 << uint(7-(bit%8))
			}
		}
		return out
	}

	prefix := readBytes(HeaderBytes)
	if !bytes.Equal(prefix[:8], magic[:]) {
		return nil, nil, fmt.Errorf("self-contained book magic mismatch")
	}
	promptBytes := binary.BigEndian.Uint32(prefix[24:28])
	compressedBytes := binary.BigEndian.Uint64(prefix[16:24])
	total := HeaderBytes + int(promptBytes) + int(compressedBytes)
	capacityBytes := b.Dx() * (b.Dy()-ControlPlaneHeight) / 8
	if total > capacityBytes {
		return nil, nil, fmt.Errorf("declared payload %d exceeds pixel capacity %d", total, capacityBytes)
	}
	return readBytes(total), control, nil
}

func drawSelfBootstrap(img *image.Gray, controlX int, header Header, masterPrompt []byte) {
	leftWidth := controlX - 48
	if leftWidth >= 300 {
		box(img, 20, 20, leftWidth, 600, 0)
		drawText(img, 36, 38, 2, "S0 SELF BOOT", 0)
		drawText(img, 36, 66, 1, "THIS IMAGE CONTAINS ITS MASTER PROMPT", 0)
		drawText(img, 36, 82, 1, "MP0 IS EXACT UTF-8 INSIDE THE PIXEL PAYLOAD", 0)
		drawText(img, 36, 98, 1, "LOAD MP0 THEN ENTER T0 BOOT AND T1 ROSETTA", 0)
		drawText(img, 36, 122, 1, fmt.Sprintf("MP0 BYTES: %d", header.MasterPromptBytes), 0)
		drawText(img, 36, 138, 1, "MP0 SHA256:", 0)
		h := header.MasterPromptSHA256
		for i := 0; i < 4 && i*16 < len(h); i++ {
			end := (i + 1) * 16
			if end > len(h) {
				end = len(h)
			}
			drawText(img, 52, 154+i*14, 1, h[i*16:end], 0)
		}
		drawText(img, 36, 224, 1, "READ ORDER:", 0)
		drawText(img, 52, 242, 1, "S0 > MP0 > T0 > T1 > T2 > T3 > VERIFY", 0)
		drawText(img, 36, 278, 1, "REQUIRED EXTERNAL INPUTS:", 0)
		drawText(img, 52, 296, 1, "THIS IMAGE + USER QUESTION", 0)
		drawText(img, 36, 332, 1, "MP0 PROMPT START:", 0)
		drawText(img, 52, 350, 1, "YOU ARE INTERACTING WITH ORIGAMI", 0)
		drawText(img, 52, 366, 1, "READ BOOT ROSETTA PROGRAM INDEX MEMORY VERIFY", 0)
		drawText(img, 36, 410, 1, "PROMPT-ONLY NATIVE READABILITY IS AN", 0)
		drawText(img, 36, 426, 1, "EMPIRICAL CLAIM AND REMAINS UNPROMOTED", 0)
	}

	rightX := controlX + 640 + 28
	rightWidth := img.Bounds().Dx() - rightX - 20
	if rightWidth >= 300 {
		box(img, rightX, 20, rightWidth, 600, 0)
		drawText(img, rightX+16, 38, 2, "MP0 MASTER PROMPT", 0)
		drawText(img, rightX+16, 72, 1, "EMBEDDED: YES", 0)
		drawText(img, rightX+16, 88, 1, "EXACT HASH BOUND: YES", 0)
		drawText(img, rightX+16, 104, 1, "SIDE-CAR FILE REQUIRED: NO", 0)
		drawText(img, rightX+16, 132, 1, "THE FULL MARKDOWN PROMPT IS NOT PRINTED", 0)
		drawText(img, rightX+16, 148, 1, "AT HUMAN FONT SIZE. IT IS CARRIED EXACTLY", 0)
		drawText(img, rightX+16, 164, 1, "IN THE SAME PIXEL PAYLOAD AS THIS BOOK.", 0)
		drawText(img, rightX+16, 202, 1, "DECODE MUST VERIFY MP0 SHA BEFORE USE.", 0)
		drawText(img, rightX+16, 240, 1, "MASTER PROMPT PURPOSE:", 0)
		drawText(img, rightX+32, 258, 1, "READ OR WRITE ORIGAMI", 0)
		drawText(img, rightX+32, 274, 1, "READ ROSETTA BEFORE SEMANTICS", 0)
		drawText(img, rightX+32, 290, 1, "USE SELECTIVE UNFOLD", 0)
		drawText(img, rightX+32, 306, 1, "UNKNOWN OVER INVENTED EXACTNESS", 0)
		drawText(img, rightX+32, 322, 1, "FALSE EXACT = 0", 0)
		drawText(img, rightX+16, 372, 1, fmt.Sprintf("PROMPT PAYLOAD BYTES: %d", len(masterPrompt)), 0)
		drawText(img, rightX+16, 410, 1, "THIS PANEL IS BOOTSTRAP METADATA.", 0)
		drawText(img, rightX+16, 426, 1, "MP0 BYTES ARE THE AUTHORITATIVE PROMPT.", 0)
	}
}
