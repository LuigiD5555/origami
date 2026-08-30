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
	ProfileID          = "origami.self-contained-book.r0-lab.profile-1"
	HeaderBytes        = 108
	DefaultDataWidth   = 2304
	TargetPNGBytes     = 512000
	ControlPlaneHeight = 640
)

var magic = [8]byte{'O', 'S', 'C', 'B', 'R', '0', 0, 1}

type Header struct {
	SourceBytes      uint64 `json:"source_bytes"`
	CompressedBytes  uint64 `json:"compressed_bytes"`
	PageCount        uint32 `json:"page_count"`
	BlockCount       uint32 `json:"block_count"`
	SourceSHA256     string `json:"source_sha256"`
	CompressedSHA256 string `json:"compressed_sha256"`
	Compression      string `json:"compression"`
}

type Rendered struct {
	Schema          string `json:"schema"`
	Profile         string `json:"profile"`
	PNG             []byte `json:"-"`
	PNGBytes        int    `json:"png_bytes"`
	TargetPNGBytes  int    `json:"target_png_bytes"`
	OverTargetBytes int    `json:"over_target_bytes"`
	PayloadBytes    int    `json:"payload_bytes"`
	DataWidthBits   int    `json:"data_width_bits"`
	DataHeightBits  int    `json:"data_height_bits"`
	ControlX        int    `json:"control_x"`
	ControlY        int    `json:"control_y"`
}

func Pack(source, compressed []byte, pageCount, blockCount uint32, compression string) ([]byte, Header, error) {
	if len(source) == 0 {
		return nil, Header{}, fmt.Errorf("source cannot be empty")
	}
	if len(compressed) == 0 {
		return nil, Header{}, fmt.Errorf("compressed residual cannot be empty")
	}
	if len(compression) == 0 || len(compression) > 8 {
		return nil, Header{}, fmt.Errorf("compression label must be 1..8 ASCII bytes")
	}
	sourceHash := sha256.Sum256(source)
	compressedHash := sha256.Sum256(compressed)
	header := make([]byte, HeaderBytes)
	copy(header[0:8], magic[:])
	binary.BigEndian.PutUint64(header[8:16], uint64(len(source)))
	binary.BigEndian.PutUint64(header[16:24], uint64(len(compressed)))
	binary.BigEndian.PutUint32(header[24:28], pageCount)
	binary.BigEndian.PutUint32(header[28:32], blockCount)
	copy(header[32:64], sourceHash[:])
	copy(header[64:96], compressedHash[:])
	copy(header[96:104], []byte(compression))
	binary.BigEndian.PutUint32(header[104:108], crc32.ChecksumIEEE(header[:104]))
	payload := make([]byte, 0, HeaderBytes+len(compressed))
	payload = append(payload, header...)
	payload = append(payload, compressed...)
	return payload, Header{
		SourceBytes: uint64(len(source)), CompressedBytes: uint64(len(compressed)),
		PageCount: pageCount, BlockCount: blockCount,
		SourceSHA256: hex.EncodeToString(sourceHash[:]), CompressedSHA256: hex.EncodeToString(compressedHash[:]),
		Compression: compression,
	}, nil
}

func Unpack(payload []byte) (Header, []byte, error) {
	if len(payload) < HeaderBytes {
		return Header{}, nil, fmt.Errorf("payload too short")
	}
	header := payload[:HeaderBytes]
	if !bytes.Equal(header[:8], magic[:]) {
		return Header{}, nil, fmt.Errorf("self-contained book magic mismatch")
	}
	if binary.BigEndian.Uint32(header[104:108]) != crc32.ChecksumIEEE(header[:104]) {
		return Header{}, nil, fmt.Errorf("self-contained book header crc mismatch")
	}
	sourceBytes := binary.BigEndian.Uint64(header[8:16])
	compressedBytes := binary.BigEndian.Uint64(header[16:24])
	if compressedBytes > uint64(len(payload)-HeaderBytes) {
		return Header{}, nil, fmt.Errorf("declared compressed residual exceeds payload")
	}
	compression := string(bytes.TrimRight(header[96:104], "\x00"))
	compressed := append([]byte(nil), payload[HeaderBytes:HeaderBytes+int(compressedBytes)]...)
	gotCompressed := sha256.Sum256(compressed)
	if !bytes.Equal(gotCompressed[:], header[64:96]) {
		return Header{}, nil, fmt.Errorf("compressed residual sha256 mismatch")
	}
	return Header{
		SourceBytes:      sourceBytes,
		CompressedBytes:  compressedBytes,
		PageCount:        binary.BigEndian.Uint32(header[24:28]),
		BlockCount:       binary.BigEndian.Uint32(header[28:32]),
		SourceSHA256:     hex.EncodeToString(header[32:64]),
		CompressedSHA256: hex.EncodeToString(header[64:96]),
		Compression:      compression,
	}, compressed, nil
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
	bits := len(payload) * 8
	dataHeight := (bits + dataWidth - 1) / dataWidth
	canvas := image.NewGray(image.Rect(0, 0, dataWidth, ControlPlaneHeight+dataHeight))
	for i := range canvas.Pix {
		canvas.Pix[i] = 0xff
	}
	controlX := (dataWidth - 640) / 2
	draw.Draw(canvas, image.Rect(controlX, 0, controlX+640, 640), control, control.Bounds().Min, draw.Src)
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
		Schema: Schema, Profile: ProfileID, PNG: buf.Bytes(), PNGBytes: buf.Len(),
		TargetPNGBytes: TargetPNGBytes, OverTargetBytes: over, PayloadBytes: len(payload),
		DataWidthBits: dataWidth, DataHeightBits: dataHeight, ControlX: controlX, ControlY: 0,
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
	compressedBytes := binary.BigEndian.Uint64(prefix[16:24])
	total := HeaderBytes + int(compressedBytes)
	capacityBytes := b.Dx() * (b.Dy()-ControlPlaneHeight) / 8
	if total > capacityBytes {
		return nil, nil, fmt.Errorf("declared payload %d exceeds pixel capacity %d", total, capacityBytes)
	}
	return readBytes(total), control, nil
}
