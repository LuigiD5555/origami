package fold_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/LuigiD5555/origami/internal/origami/fold"
	"github.com/LuigiD5555/origami/internal/origami/ir"
	"github.com/LuigiD5555/origami/internal/origami/machine"
)

func TestExactRepeatFoldUnfoldHash(t *testing.T) {
	unit := []byte("ABC")
	source := bytes.Repeat(unit, 100000)

	rep, err := fold.ExactRepeat(source, unit)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Nodes[rep.Root].Kind != ir.KindRepeat {
		t.Fatalf("root kind = %q, want REPEAT", rep.Nodes[rep.Root].Kind)
	}
	if len(rep.Nodes) != 2 {
		t.Fatalf("nodes = %d, want 2", len(rep.Nodes))
	}

	got, stats, err := machine.UnfoldExact(rep, machine.Budget{
		MaxBytes: uint64(len(source)),
		MaxNodes: 100001,
		MaxDepth: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(got) != sha256.Sum256(source) {
		t.Fatal("unfolded SHA-256 differs from source")
	}
	if stats.BytesMaterialized != uint64(len(source)) {
		t.Fatalf("materialized bytes = %d, want %d", stats.BytesMaterialized, len(source))
	}
}

func TestExactRepeatRejectsNonRepeatingSource(t *testing.T) {
	if _, err := fold.ExactRepeat([]byte("ABCABD"), []byte("ABC")); err == nil {
		t.Fatal("expected non-repeating source to be rejected")
	}
}

func TestUnfoldRequiresFiniteBudget(t *testing.T) {
	rep, err := fold.ExactRepeat([]byte("ABCABC"), []byte("ABC"))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := machine.UnfoldExact(rep, machine.Budget{}); err == nil {
		t.Fatal("expected zero budget to be rejected")
	}
}
