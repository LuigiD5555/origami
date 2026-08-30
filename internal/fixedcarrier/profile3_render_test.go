package fixedcarrier

import (
	"bytes"
	"testing"
)

func TestProfile3RoundTripFrozenSize(t *testing.T) {
	meta := Metadata{
		CarrierID: "p3-demo", StoreRoot: "root", SourceSHA256: "source",
		PageCount: 100, BlockCount: 900, DocumentCount: 1,
		GraphSignature: []byte(`{"i":["PART I FOUNDATIONS","PART II ATTENTION","PART III SYSTEMS"]}`),
	}
	pngBytes, want, err := RenderProfile3(meta)
	if err != nil { t.Fatal(err) }
	if len(pngBytes) != FixedPNGBytes { t.Fatalf("profile-3 size=%d want=%d", len(pngBytes), FixedPNGBytes) }
	got, err := DecodeProfile3PNG(pngBytes)
	if err != nil { t.Fatal(err) }
	if got.Profile != Profile3ID { t.Fatalf("profile=%s", got.Profile) }
	if got.CarrierDigest != want.CarrierDigest { t.Fatalf("digest mismatch") }
	if got.PageCount != meta.PageCount || got.BlockCount != meta.BlockCount { t.Fatalf("metadata mismatch: %+v", got.Metadata) }
	if !bytes.Contains([]byte(got.BootText[2]), []byte("SMALLEST S*")) { t.Fatalf("semantic decoder route absent: %v", got.BootText) }
	if !bytes.Contains([]byte(got.BootText[3]), []byte("SMALLEST E*")) { t.Fatalf("semantic encoder route absent: %v", got.BootText) }
}

func TestDecodeAnyAcceptsProfile2AndProfile3(t *testing.T) {
	meta := Metadata{CarrierID:"x",StoreRoot:"r",SourceSHA256:"s",PageCount:1,BlockCount:2,GraphSignature:[]byte(`{"i":["ONE","TWO"]}`)}
	p2, _, err := Render(meta); if err != nil { t.Fatal(err) }
	p3, _, err := RenderProfile3(meta); if err != nil { t.Fatal(err) }
	got2, err := DecodeAnyPNG(p2); if err != nil { t.Fatal(err) }
	got3, err := DecodeAnyPNG(p3); if err != nil { t.Fatal(err) }
	if got2.Profile != ProfileID { t.Fatalf("profile2 regression: %s", got2.Profile) }
	if got3.Profile != Profile3ID { t.Fatalf("profile3 decode: %s", got3.Profile) }
}

func TestProfile3DeclaresRequiredZones(t *testing.T) {
	p := Profile3()
	want := map[string]bool{"T0":false,"T1":false,"T2":false,"PROGRAM":false,"T3":false,"EXACT":false,"VERIFY":false}
	for _, z := range p.Zones { if _, ok := want[z.ID]; ok { want[z.ID] = true } }
	for id, ok := range want { if !ok { t.Fatalf("missing profile-3 zone %s", id) } }
	if p.FixedPNGBytes != 8192 || p.Width != 640 || p.Height != 640 { t.Fatalf("physical invariant drift: %+v", p) }
}
