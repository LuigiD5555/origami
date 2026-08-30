package document

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/memory"
)

func TestTextIngestionBuildsExactPlaneWithMetadataOnlyRouting(t *testing.T) {
	input := []byte("# Introduction\n\nTransformers use attention for sequence modeling.\n\nEfficiency depends on selective computation.\f# Results\n\nThe evidence reports lower inference cost.")
	doc, err := FromText("paper-demo", "application/pdf", input)
	if err != nil { t.Fatal(err) }
	if doc.Plane != "EXACT_SOURCE_PLANE" || len(doc.Pages) != 2 { t.Fatalf("unexpected document: %+v", doc) }
	bundle, err := BuildBundle(doc); if err != nil { t.Fatal(err) }
	store, err := memory.BuildStore("paper-demo", bundle.Nodes); if err != nil { t.Fatal(err) }
	packet := memory.NewManager(store).Query("inference cost evidence", memory.ContextBudget{MaxTokens: 120})
	if packet.Unknown || len(packet.Items) == 0 { t.Fatalf("addressable source query failed: %+v", packet) }
	found := false
	for _, item := range packet.Items {
		if item.Kind == "source-block" && item.Fidelity == memory.FidelitySummary {
			t.Fatal("raw source block must not pretend to have a semantic summary")
		}
		if item.Kind == "source-block" && (item.Fidelity == memory.FidelityEvidence || item.Fidelity == memory.FidelityExact || item.Fidelity == memory.FidelityAbstract) {
			found = true
		}
	}
	if !found { t.Fatalf("expected source block or metadata candidate: %+v", packet.Items) }
}

func TestSourceBlockExactnessIsHashBound(t *testing.T) {
	doc, err := FromText("source", "text/plain", []byte("alpha beta gamma")); if err != nil { t.Fatal(err) }
	nodes, err := ToMemoryNodes(doc); if err != nil { t.Fatal(err) }
	var block memory.Node
	for _, node := range nodes { if node.Kind == "source-block" { block = node; break } }
	if block.Address == "" { t.Fatal("missing source block") }
	exact, ok := block.Level(memory.FidelityExact); if !ok || !exact.Verified || exact.Text != "alpha beta gamma" { t.Fatalf("exact source fidelity missing: %+v", block) }
	if block.Provenance.SourceHash == "" || block.VerificationRef != block.Provenance.SourceHash { t.Fatalf("exact block not bound to source hash: %+v", block) }
}

func TestLayoutJSONPreservesImageOrPDFGeometryWithoutGrantingSemantics(t *testing.T) {
	text := "OCR-visible label 2026"
	doc := Document{Schema: SchemaR0, SourceID: "scan", SourceHash: contentHash([]byte("scan-root")), MediaType: "image/png", Plane: "EXACT_SOURCE_PLANE", Pages: []Page{{Number: 1, Blocks: []Block{{ID: "block-000001", Page: 1, Order: 1, Kind: "image_ocr", Text: text, BBox: &BBox{X0: 10, Y0: 20, X1: 200, Y1: 60}, SourceHash: contentHash([]byte(text))}}}}}
	if err := Validate(doc); err != nil { t.Fatal(err) }
	nodes, err := ToMemoryNodes(doc); if err != nil { t.Fatal(err) }
	for _, node := range nodes {
		if node.Kind != "source-block" { continue }
		if node.SemanticClass != "image_ocr" { t.Fatalf("layout kind lost: %+v", node) }
		if _, ok := node.Level(memory.FidelitySummary); ok { t.Fatal("OCR text must not be upgraded into semantic summary automatically") }
		return
	}
	t.Fatal("missing OCR source block")
}
