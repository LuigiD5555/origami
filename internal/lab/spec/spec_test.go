package spec

import "testing"

func TestValidatePDFSource(t *testing.T) {
	experiment := ExperimentSpec{Schema: SchemaV1, ID: "pdf", Baseline: "R3.10", Campaign: "PDF", Seed: SeedSpec{Algorithm: "splitmix64-v1"}, Source: SourceSpec{Kind: "file", Generator: "pdf-extract-v1", Path: "document.pdf"}, Renderer: RendererSpec{ID: "pdf-carrier", Version: "v1", Width: 640, Height: 480}, Variants: 1, Evaluators: []string{"computational"}}
	if err := experiment.Validate(); err != nil {
		t.Fatal(err)
	}
	experiment.Source.Length = 1
	if err := experiment.Validate(); err == nil {
		t.Fatal("file source with declared length should fail")
	}
}
