package ir_test

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/origami/ir"
)

func TestValidateRejectsCycle(t *testing.T) {
	r := ir.Representation{
		Root: "a",
		Nodes: map[string]ir.Node{
			"a": {ID: "a", Kind: ir.KindRef, Target: "b"},
			"b": {ID: "b", Kind: ir.KindRef, Target: "a"},
		},
	}
	if err := r.Validate(); err == nil {
		t.Fatal("expected cycle to be rejected")
	}
}

func TestValidateConcat(t *testing.T) {
	r := ir.Representation{
		Root: "root",
		Nodes: map[string]ir.Node{
			"a":    {ID: "a", Kind: ir.KindLiteral, Literal: []byte("A")},
			"b":    {ID: "b", Kind: ir.KindLiteral, Literal: []byte("B")},
			"root": {ID: "root", Kind: ir.KindConcat, Children: []string{"a", "b"}},
		},
	}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}
