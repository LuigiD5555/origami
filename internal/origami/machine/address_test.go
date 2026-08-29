package machine_test

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/origami/ir"
	"github.com/LuigiD5555/origami/internal/origami/machine"
)

func TestUnfoldAddressMaterializesOnlySelectedClosure(t *testing.T) {
	r := ir.Representation{
		Root: "document",
		Nodes: map[string]ir.Node{
			"intro":    {ID: "intro", Kind: ir.KindLiteral, Literal: []byte("INTRO")},
			"pattern":  {ID: "pattern", Kind: ir.KindLiteral, Literal: []byte("ABC")},
			"repeated": {ID: "repeated", Kind: ir.KindRepeat, Target: "pattern", Count: 1000},
			"document": {ID: "document", Kind: ir.KindConcat, Children: []string{"intro", "repeated"}},
		},
	}

	got, stats, err := machine.UnfoldAddress(r, "intro", machine.Budget{MaxBytes: 5, MaxNodes: 1, MaxDepth: 1})
	if err != nil { t.Fatal(err) }
	if string(got) != "INTRO" { t.Fatalf("got %q, want INTRO", got) }
	if stats.NodesVisited != 1 { t.Fatalf("visited %d nodes, want 1", stats.NodesVisited) }
	if stats.BytesMaterialized != 5 { t.Fatalf("materialized %d bytes, want 5", stats.BytesMaterialized) }
}

func TestUnfoldAddressFollowsReferenceClosure(t *testing.T) {
	r := ir.Representation{
		Root: "root",
		Nodes: map[string]ir.Node{
			"literal": {ID: "literal", Kind: ir.KindLiteral, Literal: []byte("VALUE")},
			"alias":   {ID: "alias", Kind: ir.KindRef, Target: "literal"},
			"root":    {ID: "root", Kind: ir.KindConcat, Children: []string{"alias", "literal"}},
		},
	}
	got, stats, err := machine.UnfoldAddress(r, "alias", machine.Budget{MaxBytes: 5, MaxNodes: 2, MaxDepth: 2})
	if err != nil { t.Fatal(err) }
	if string(got) != "VALUE" { t.Fatalf("got %q, want VALUE", got) }
	if stats.NodesVisited != 2 { t.Fatalf("visited %d nodes, want 2", stats.NodesVisited) }
}

func TestUnfoldAddressRejectsUnknownAddress(t *testing.T) {
	r := ir.Representation{Root: "root", Nodes: map[string]ir.Node{"root": {ID: "root", Kind: ir.KindLiteral, Literal: []byte("X")}}}
	if _, _, err := machine.UnfoldAddress(r, "missing", machine.Budget{MaxBytes: 1, MaxNodes: 1, MaxDepth: 1}); err == nil {
		t.Fatal("expected unknown address to be rejected")
	}
}
