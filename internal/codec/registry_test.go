package codec

import "testing"

func TestReferenceRegistryValid(t *testing.T) {
	r := ReferenceRegistry()
	if err := r.Validate(); err != nil {
		t.Fatalf("reference registry invalid: %v", err)
	}
}

func TestSemanticCodecDoesNotRequireExactCapability(t *testing.T) {
	r := ReferenceRegistry()
	for _, e := range r.Entries {
		if e.Family != SemanticDecoder && e.Family != SemanticEncoder {
			continue
		}
		if e.Exact {
			t.Fatalf("semantic codec %s marked exact", e.ID)
		}
		for _, cap := range e.RequiredCapabilities {
			if cap == "EXACT_READ" || cap == "EXACT_WRITE" || cap == "VERIFIER" {
				t.Fatalf("semantic codec %s depends on exact capability %s", e.ID, cap)
			}
		}
	}
}

func TestSelectS2WithSemanticReadOnly(t *testing.T) {
	r := ReferenceRegistry()
	e, err := r.Select("READ_SUPERINDEX", SemanticDecoder, map[string]bool{"SEMANTIC_READ": true})
	if err != nil {
		t.Fatal(err)
	}
	if e.ID != "S2" {
		t.Fatalf("got %s want S2", e.ID)
	}
}

func TestS2E2RoundTrip(t *testing.T) {
	want := SemanticIndex{Entries: []IndexEntry{
		{ID: "part-1", Label: "Foundations", Address: "ohf://book/index/part-1"},
		{ID: "part-2", Label: "Deep Learning", Address: "ohf://book/index/part-2"},
		{ID: "part-3", Label: "Attention", Address: "ohf://book/index/part-3"},
	}}
	ir, err := EncodeSuperIndex(want, "origami.fixed-carrier.r2.profile-3")
	if err != nil {
		t.Fatal(err)
	}
	if ir.CodecID != "E2" {
		t.Fatalf("got codec %s want E2", ir.CodecID)
	}
	got, err := DecodeSuperIndex(ir)
	if err != nil {
		t.Fatal(err)
	}
	if !EqualSemanticIndex(want, got) {
		t.Fatalf("roundtrip mismatch: %#v != %#v", want, got)
	}
}

func TestEncodeSuperIndexRejectsDuplicateIDs(t *testing.T) {
	_, err := EncodeSuperIndex(SemanticIndex{Entries: []IndexEntry{
		{ID: "x", Label: "A"},
		{ID: "x", Label: "B"},
	}}, "profile")
	if err == nil {
		t.Fatal("expected duplicate ID error")
	}
}
