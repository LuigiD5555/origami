package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"
	"os"
)

const Profile3ID = "origami.fixed-carrier.r2.profile-3"

var profile3BootText = []string{
	"ORIGAMI FIXED CARRIER R2 PROFILE 3",
	"BOOT: READ T1 ROSETTA CODECS BEFORE T2",
	"READ: CLASSIFY > SMALLEST S* > SEMANTIC STATE",
	"WRITE: SEMANTIC IR > SMALLEST E* > CONSTRUCTION",
	"SEMANTIC != EXACT. X*/Q* ARE OPTIONAL",
	"UNKNOWN / NOT VERIFIED > INVENTION. FALSE_EXACT=0",
}

var profile3 = Profile{
	ID:            Profile3ID,
	Width:         Width,
	Height:        Height,
	FixedPNGBytes: FixedPNGBytes,
	MaxPNGBytes:   MaxPNGBytes,
	BootText:      append([]string(nil), profile3BootText...),
	Zones: []Zone{
		{ID: "T0", Purpose: "protocol/profile bootstrap", X: 18, Y: 16, W: 604, H: 112},
		{ID: "T1", Purpose: "ROSETTA grammar plus decoder/encoder entrypoints and capabilities", X: 18, Y: 128, W: 604, H: 86},
		{ID: "T2", Purpose: "bounded actual semantic superindex", X: 18, Y: 214, W: 604, H: 82},
		{ID: "PROGRAM", Purpose: "self-declared semantic codec procedures", X: 18, Y: 296, W: 190, H: 252},
		{ID: "T3", Purpose: "exact/control record; never required for semantic index", X: 214, Y: 296, W: 212, H: 252},
		{ID: "EXACT", Purpose: "X/Q exact codec declarations and verification boundary", X: 432, Y: 296, W: 190, H: 252},
		{ID: "VERIFY", Purpose: "semantic roundtrip and exactness footer", X: 24, Y: 552, W: 592, H: 72},
	},
	FamilyBindings: append([]FamilyBind(nil), profileR2.FamilyBindings...),
}

func Profile3() Profile { return profile3 }

// RenderProfile3 renders the protocol-aware candidate without changing the
// profile-2 default renderer. The exact/control record is still deterministic,
// while T1/PROGRAM expose the semantic read/write codec route visibly.
func RenderProfile3(meta Metadata) ([]byte, Decoded, error) {
	rec, decoded, err := encodeProfile3Record(meta)
	if err != nil {
		return nil, Decoded{}, err
	}
	img := image.NewGray(image.Rect(0, 0, Width, Height))
	for i := range img.Pix { img.Pix[i] = 0xff }
	drawFrame(img)
	drawProfile3T0(img)
	drawProfile3T1(img, decoded.VisualProbe)
	drawT2(img, decoded.GraphSignature)
	drawProfile3ProgramExact(img, rec)
	drawProfile3Verify(img, decoded.VisualProbe)

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		return nil, Decoded{}, err
	}
	if buf.Len() > FixedPNGBytes-12 {
		return nil, Decoded{}, fmt.Errorf("profile-3 raw PNG exceeds frozen envelope: %d > %d", buf.Len(), FixedPNGBytes-12)
	}
	padded, err := padPNG(buf.Bytes(), FixedPNGBytes)
	if err != nil { return nil, Decoded{}, err }
	if len(padded) > MaxPNGBytes {
		return nil, Decoded{}, fmt.Errorf("profile-3 exceeds hard limit: %d > %d", len(padded), MaxPNGBytes)
	}
	return padded, decoded, nil
}

func WriteProfile3PNG(path string, meta Metadata) (Decoded, int, error) {
	data, dec, err := RenderProfile3(meta)
	if err != nil { return Decoded{}, 0, err }
	if err := os.WriteFile(path, data, 0644); err != nil { return Decoded{}, 0, err }
	return dec, len(data), nil
}

func DecodeProfile3PNG(data []byte) (Decoded, error) {
	if len(data) > MaxPNGBytes { return Decoded{}, fmt.Errorf("carrier exceeds hard limit: %d", len(data)) }
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil { return Decoded{}, err }
	if img.Bounds().Dx() != Width || img.Bounds().Dy() != Height {
		return Decoded{}, fmt.Errorf("unexpected carrier dimensions %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
	rec := make([]byte, RecordBytes)
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			gx := GridX + x*Cell + Cell/2
			gy := GridY + y*Cell + Cell/2
			r, g, b, _ := img.At(gx, gy).RGBA()
			if (r+g+b)/3 < 0x8000 { rec[bit/8] |= 1 << uint(7-(bit%8)) }
			bit++
		}
	}
	return decodeProfile3Record(rec)
}

func DecodeAnyPNG(data []byte) (Decoded, error) {
	if dec, err := DecodePNG(data); err == nil { return dec, nil }
	return DecodeProfile3PNG(data)
}

func encodeProfile3Record(meta Metadata) ([]byte, Decoded, error) {
	rec, decoded, err := encodeRecord(meta)
	if err != nil { return nil, Decoded{}, err }
	profileHash := sha256.Sum256([]byte(Profile3ID))
	copy(rec[138:154], profileHash[:16])
	bootHash := sha256.Sum256([]byte(joinBootText(profile3BootText)))
	copy(rec[154:170], bootHash[:16])
	fillRedundancy(rec[426:508], rec[:426])
	binary.BigEndian.PutUint32(rec[508:512], crc32.ChecksumIEEE(rec[:508]))
	digest := sha256.Sum256(rec)
	decoded.Profile = Profile3ID
	decoded.CarrierDigest = hex.EncodeToString(digest[:])
	decoded.VisualProbe = probeBits(digest[0])
	decoded.BootText = append([]string(nil), profile3BootText...)
	decoded.Zones = append([]Zone(nil), profile3.Zones...)
	return rec, decoded, nil
}

func decodeProfile3Record(rec []byte) (Decoded, error) {
	if len(rec) != RecordBytes { return Decoded{}, fmt.Errorf("invalid record length") }
	if !bytes.Equal(rec[:8], magic[:]) { return Decoded{}, fmt.Errorf("fixed carrier R2 magic mismatch") }
	if binary.BigEndian.Uint32(rec[508:512]) != crc32.ChecksumIEEE(rec[:508]) { return Decoded{}, fmt.Errorf("carrier crc mismatch") }
	tool := sha256.Sum256([]byte(ToolProtocol))
	if !bytes.Equal(rec[122:138], tool[:16]) { return Decoded{}, fmt.Errorf("tool protocol mismatch") }
	profileHash := sha256.Sum256([]byte(Profile3ID))
	if !bytes.Equal(rec[138:154], profileHash[:16]) { return Decoded{}, fmt.Errorf("profile-3 mismatch") }
	bootHash := sha256.Sum256([]byte(joinBootText(profile3BootText)))
	if !bytes.Equal(rec[154:170], bootHash[:16]) { return Decoded{}, fmt.Errorf("profile-3 T0 binding mismatch") }
	wantRed := make([]byte, 82)
	fillRedundancy(wantRed, rec[:426])
	if !bytes.Equal(wantRed, rec[426:508]) { return Decoded{}, fmt.Errorf("carrier redundancy mismatch") }
	digest := sha256.Sum256(rec)
	sig := make([]byte, 256); copy(sig, rec[170:426])
	return Decoded{
		Schema: Schema, Profile: Profile3ID, ToolProtocol: ToolProtocol, AddressABI: AddressABI,
		Metadata: Metadata{
			CarrierID: hex.EncodeToString(rec[26:58]), StoreRoot: hex.EncodeToString(rec[58:90]), SourceSHA256: hex.EncodeToString(rec[90:122]),
			PageCount: binary.BigEndian.Uint32(rec[10:14]), BlockCount: binary.BigEndian.Uint32(rec[14:18]), DocumentCount: binary.BigEndian.Uint32(rec[18:22]), ObjectCount: binary.BigEndian.Uint32(rec[22:26]), GraphSignature: sig, Flags: binary.BigEndian.Uint16(rec[8:10]),
		},
		CarrierDigest: hex.EncodeToString(digest[:]), VisualProbe: probeBits(digest[0]), BootText: append([]string(nil), profile3BootText...), Zones: append([]Zone(nil), profile3.Zones...),
	}, nil
}

func drawProfile3T0(img *image.Gray) {
	box(img, 14, 12, 612, 112, 0)
	y := 20
	for _, s := range profile3BootText {
		scale := 2
		if textWidth(s, scale) > Width-48 { scale = 1 }
		drawText(img, 24, y, scale, s, 0)
		y += 16
	}
}

func drawProfile3T1(img *image.Gray, probe string) {
	drawText(img, 24, 132, 1, "T1 ROSETTA + CODECS", 0)
	drawText(img, 24, 146, 1, "READ S0 ID  S1 HIER  S2 INDEX", 0)
	drawText(img, 318, 146, 1, "WRITE E0 ID  E1 HIER  E2 INDEX", 0)
	drawText(img, 24, 160, 1, "CORE SEMANTIC R/W | EXACT OPTIONAL", 0)
	drawProbeRow(img, 184, probe)
}

func drawProfile3ProgramExact(img *image.Gray, rec []byte) {
	box(img, 14, 296, 194, 252, 0)
	drawText(img, 24, 304, 1, "PROGRAM", 0)
	drawText(img, 24, 324, 1, "S2: T2 -> INDEX", 0)
	drawText(img, 24, 340, 1, "E2: INDEX -> T2", 0)
	drawText(img, 24, 356, 1, "SEMANTIC: S* BEFORE X*", 0)
	drawText(img, 24, 380, 1, "UNSUPPORTED -> UNKNOWN", 0)

	box(img, 432, 296, 194, 252, 0)
	drawText(img, 442, 304, 1, "EXACT X*/Q* OPTIONAL", 0)
	drawText(img, 442, 326, 1, "CID HASH MERKLE", 0)
	drawText(img, 442, 342, 1, "RESIDUAL VERIFY", 0)
	drawText(img, 442, 366, 1, "NOT NEEDED FOR T2", 0)

	box(img, GridX-5, GridY-5, GridBits*Cell+10, GridBits*Cell+10, 0)
	bit := 0
	for y := 0; y < GridBits; y++ {
		for x := 0; x < GridBits; x++ {
			b := (rec[bit/8] >> uint(7-(bit%8))) & 1
			if b == 1 { fill(img, GridX+x*Cell, GridY+y*Cell, Cell, Cell, 0) }
			bit++
		}
	}
}

func drawProfile3Verify(img *image.Gray, probe string) {
	drawProbeRow(img, 566, probe)
	drawText(img, 24, 610, 1, "VERIFY: S2(E2(INDEX)) ~= INDEX | FALSE EXACT 0 | UNKNOWN > INVENT", 0)
}
