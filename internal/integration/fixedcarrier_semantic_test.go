package integration

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/memory"
)

func TestCompactSignatureCarriesSemanticIndexAsValidJSON(t *testing.T) {
	sig := memory.GraphSignature{
		Schema:          memory.SchemaR0 + ".graph-signature",
		CarrierID:       "book-demo",
		RootTerms:       []string{"dynamic", "programming", "optimization", "algorithms"},
		SemanticClasses: []string{"chapter", "concept"},
		Clusters:        []string{"part-1", "part-2"},
		NodeCount:       5000,
	}
	nodes := []memory.Node{
		{Address: "ohf://book/root/1", Label: "Introduction", Importance: memory.ImportanceProfile{Global: .9}},
		{Address: "ohf://book/root/2", Label: "Dynamic Programming", Importance: memory.ImportanceProfile{Global: 1}},
		{Address: "ohf://book/root/3", Label: "Greedy Algorithms", Importance: memory.ImportanceProfile{Global: .8}},
		{Address: "ohf://book/root/4", Label: "Graph Algorithms", Importance: memory.ImportanceProfile{Global: .7}},
	}
	body, err := compactSignature(sig, nodes)
	if err != nil { t.Fatal(err) }
	if len(body) > 256 { t.Fatalf("semantic T2 hint exceeded control-record budget: %d", len(body)) }
	var decoded struct { Index []string `json:"i"` }
	if err := json.Unmarshal(body, &decoded); err != nil { t.Fatalf("bounded signature is invalid JSON: %v: %q", err, body) }
	if len(decoded.Index) < 3 { t.Fatalf("expected semantic index labels, got %v", decoded.Index) }
	joined := strings.Join(decoded.Index, " ")
	if !strings.Contains(joined, "Dynamic Programming") || !strings.Contains(joined, "Introduction") { t.Fatalf("missing semantic labels: %v", decoded.Index) }
}

func TestCompactSignatureNeverUsesRawJSONTruncation(t *testing.T) {
	long := strings.Repeat("semantic-section-", 20)
	sig := memory.GraphSignature{Schema: memory.SchemaR0 + ".graph-signature", CarrierID: long, RootTerms: []string{long, long + "2"}, SemanticClasses: []string{long}, Clusters: []string{long}, NodeCount: 99999}
	nodes := []memory.Node{{Address: "root", Label: long, Importance: memory.ImportanceProfile{Global: 1}}}
	body, err := compactSignature(sig, nodes)
	if err != nil { t.Fatal(err) }
	if len(body) > 256 { t.Fatalf("bounded signature exceeded 256 bytes: %d", len(body)) }
	var any map[string]any
	if err := json.Unmarshal(body, &any); err != nil { t.Fatalf("signature must remain valid JSON instead of truncating bytes: %v", err) }
}
