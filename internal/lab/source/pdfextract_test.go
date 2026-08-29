package source

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/spec"
)

func TestGeneratePDFBuildsReconstructablePageIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fixture.pdf")
	if err := os.WriteFile(path, minimalPDF("Thinking swarms coordinate cognition."), 0o600); err != nil {
		t.Fatal(err)
	}
	experiment := spec.ExperimentSpec{Source: spec.SourceSpec{Kind: "file", Generator: GeneratorPDFExtractV1, Path: path}}
	material, index, err := GeneratePDF(experiment)
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Pages) != 1 || len(material.PageIndex) != 1 {
		t.Fatalf("unexpected page index: %+v", index)
	}
	page := index.Pages[0]
	if got := material.Bytes[page.ByteOffset : page.ByteOffset+page.ByteLength]; !bytes.Equal(got, material.Bytes) {
		t.Fatal("indexed ranges do not reconstruct canonical material")
	}
}

func minimalPDF(text string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
	}
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
