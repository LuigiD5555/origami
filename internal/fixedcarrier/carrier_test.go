package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"image"
	"image/draw"
	"image/png"
	"testing"
)

func TestR2RoundTripFrozenSizeAndBootBinding(t *testing.T) {
	meta := Metadata{CarrierID: "demo", StoreRoot: "root", SourceSHA256: "source", PageCount: 800, BlockCount: 5500, DocumentCount: 1, ObjectCount: 6301, GraphSignature: bytes.Repeat([]byte{0x5a}, 256)}
	pngBytes, want, err := Render(meta)
	if err != nil { t.Fatal(err) }
	if len(pngBytes) != FixedPNGBytes { t.Fatalf("carrier size not frozen: %d", len(pngBytes)) }
	if len(pngBytes) > MaxPNGBytes { t.Fatalf("hard max exceeded: %d", len(pngBytes)) }
	got, err := DecodePNG(pngBytes)
	if err != nil { t.Fatal(err) }
	if got.Schema != Schema || got.Profile != ProfileID || got.ToolProtocol != ToolProtocol || got.AddressABI != AddressABI { t.Fatalf("contract mismatch: %+v", got) }
	if got.StoreRoot != want.StoreRoot || got.SourceSHA256 != want.SourceSHA256 || got.PageCount != 800 || got.BlockCount != 5500 || got.DocumentCount != 1 || got.ObjectCount != 6301 { t.Fatalf("roundtrip mismatch: got=%+v want=%+v", got, want) }
	if len(got.BootText) < 5 { t.Fatalf("T0 boot text missing: %+v", got.BootText) }
	top, bottom, err := VisualProbeFromPNG(pngBytes)
	if err != nil { t.Fatal(err) }
	if top != want.VisualProbe || bottom != want.VisualProbe { t.Fatalf("probe mismatch top=%s bottom=%s want=%s", top, bottom, want.VisualProbe) }
}

func TestCarrierSizeDoesNotGrowWithCorpusOrEntropy(t *testing.T) {
	for i := 0; i < 128; i++ {
		seed := sha256.Sum256([]byte{byte(i), byte(i >> 8)})
		sig := make([]byte, 256)
		for j := 0; j < 8; j++ { h := sha256.Sum256(append(seed[:], byte(j))); copy(sig[j*32:], h[:]) }
		meta := Metadata{CarrierID: string(seed[:]), StoreRoot: string(sig[:32]), SourceSHA256: string(sig[32:64]), PageCount: uint32(1 + i*7000), BlockCount: uint32(1 + i*23000), DocumentCount: uint32(1 + i), ObjectCount: uint32(3 + i*30000), GraphSignature: sig}
		b, _, err := Render(meta)
		if err != nil { t.Fatalf("case %d: %v", i, err) }
		if len(b) != FixedPNGBytes { t.Fatalf("case %d size changed: %d", i, len(b)) }
	}
}

func TestProbeSurvivesNearestResizeRoundTrip(t *testing.T) {
	b, want, err := Render(Metadata{CarrierID: "resize", StoreRoot: "r", SourceSHA256: "s", PageCount: 800, BlockCount: 5000, GraphSignature: bytes.Repeat([]byte{0xa5}, 256)})
	if err != nil { t.Fatal(err) }
	img, err := png.Decode(bytes.NewReader(b)); if err != nil { t.Fatal(err) }
	small := nearest(img, 320, 320); restored := nearest(small, Width, Height)
	var buf bytes.Buffer; if err := png.Encode(&buf, restored); err != nil { t.Fatal(err) }
	top, bottom, err := VisualProbeFromPNG(buf.Bytes()); if err != nil { t.Fatal(err) }
	if top != want.VisualProbe || bottom != want.VisualProbe { t.Fatalf("resized probe mismatch top=%s bottom=%s want=%s", top, bottom, want.VisualProbe) }
}

func TestR2ProfileHasFourBootstrapLevels(t *testing.T) {
	p := R2Profile()
	if p.Width != Width || p.Height != Height || p.FixedPNGBytes != FixedPNGBytes { t.Fatalf("profile physical mismatch: %+v", p) }
	want := map[string]bool{"T0": false, "T1": false, "T2": false, "T3": false, "VERIFY": false}
	for _, z := range p.Zones { if _, ok := want[z.ID]; ok { want[z.ID] = true } }
	for id, ok := range want { if !ok { t.Fatalf("missing zone %s", id) } }
	if len(p.FamilyBindings) < 10 { t.Fatalf("family profile too small: %d", len(p.FamilyBindings)) }
}

func nearest(src image.Image, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h)); sb := src.Bounds()
	for y := 0; y < h; y++ { sy := sb.Min.Y + y*sb.Dy()/h; for x := 0; x < w; x++ { sx := sb.Min.X + x*sb.Dx()/w; dst.Set(x, y, src.At(sx, sy)) } }
	return dst
}
var _ = draw.Src
