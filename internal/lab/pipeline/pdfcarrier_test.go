package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/query"
	"github.com/LuigiD5555/origami/internal/lab/source"
	"github.com/LuigiD5555/origami/internal/lab/spec"
)

func TestPDFCarrierEndToEnd(t *testing.T) {
	root := t.TempDir()
	for _, relativePath := range []string{"internal/lab/renderer/renderer.go", "internal/lab/renderer/pdfcarrier/pdfcarrier.go", "internal/lab/canonicaldocument/document.go", "internal/lab/superindex/superindex.go", "internal/lab/capacity/capacity.go"} {
		path := filepath.Join(root, filepath.FromSlash(relativePath))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture renderer source"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	pdfPath := filepath.Join(root, "fixture.pdf")
	if err := os.WriteFile(pdfPath, pipelineFixturePDF("Thinking swarms coordinate cognition."), 0o600); err != nil {
		t.Fatal(err)
	}
	experiment := spec.ExperimentSpec{Schema: spec.SchemaV1, ID: "pdf-e2e", Baseline: "R3.10", Campaign: "PDF_CARRIER_R0", Seed: spec.SeedSpec{Algorithm: "splitmix64-v1", Value: 7}, Source: spec.SourceSpec{Kind: "file", Generator: source.GeneratorPDFExtractV1, Path: pdfPath}, Renderer: spec.RendererSpec{ID: "pdf-carrier", Version: "v1", Width: 640, Height: 480}, Variants: 1, Evaluators: []string{"computational"}}
	result, err := Run(root, experiment)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "BLOCKED_BASELINE" || result.PNGBytes <= 0 || result.PNGBytes > MaxCarrierBytes {
		t.Fatalf("unexpected pipeline result: %+v", result)
	}
	indexBytes, err := os.ReadFile(filepath.Join(result.RunDir, "page_index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var index source.PageIndex
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		t.Fatal(err)
	}
	material, err := os.ReadFile(filepath.Join(result.RunDir, "source", "source.bin"))
	if err != nil {
		t.Fatal(err)
	}
	if err := source.VerifyPageIndex(material, index); err != nil {
		t.Fatal(err)
	}
	plan := query.Plan(index.Pages, "What coordinates a thinking swarm?")
	if len(plan.Addresses) != 1 || plan.ClosureBytes != plan.TotalBytes {
		t.Fatalf("unexpected one-page query plan: %+v", plan)
	}
}

func pipelineFixturePDF(text string) []byte {
	objects := []string{"<< /Type /Catalog /Pages 2 0 R >>", "<< /Type /Pages /Kids [3 0 R] /Count 1 >>", "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>", "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>"}
	stream := fmt.Sprintf("BT /F1 12 Tf 72 720 Td (%s) Tj ET", text)
	objects = append(objects, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream))
	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n")
	offsets := []int{0}
	for objectIndex, object := range objects {
		offsets = append(offsets, document.Len())
		fmt.Fprintf(&document, "%d 0 obj\n%s\nendobj\n", objectIndex+1, object)
	}
	xrefOffset := document.Len()
	fmt.Fprintf(&document, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&document, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&document, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)
	return document.Bytes()
}
