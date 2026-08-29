package fold

import (
	"bytes"
	"fmt"

	"github.com/LuigiD5555/origami/internal/origami/ir"
)

// ExactRepeat recognizes the deliberately small first Fold profile: a source
// made of an exact repetition of one non-empty unit. It is deterministic and
// refuses sources that cannot be represented by this profile.
func ExactRepeat(source, unit []byte) (ir.Representation, error) {
	if len(source) == 0 || len(unit) == 0 {
		return ir.Representation{}, fmt.Errorf("origami fold: source and unit must be non-empty")
	}
	if len(source)%len(unit) != 0 {
		return ir.Representation{}, fmt.Errorf("origami fold: source is not an exact repetition")
	}
	count := len(source) / len(unit)
	if count < 2 || !bytes.Equal(source, bytes.Repeat(unit, count)) {
		return ir.Representation{}, fmt.Errorf("origami fold: source is not an exact repetition")
	}

	return ir.Representation{
		Root: "root",
		Nodes: map[string]ir.Node{
			"pattern": {ID: "pattern", Kind: ir.KindLiteral, Literal: append([]byte(nil), unit...)},
			"root":    {ID: "root", Kind: ir.KindRepeat, Target: "pattern", Count: uint64(count)},
		},
	}, nil
}

// ExactRepeatAuto finds the shortest exact repeating unit. This is deliberately
// a simple deterministic reference Fold, not a general-purpose compressor.
func ExactRepeatAuto(source []byte) (ir.Representation, error) {
	if len(source) < 2 {
		return ir.Representation{}, fmt.Errorf("origami fold: source has no repeated structure")
	}
	for size := 1; size <= len(source)/2; size++ {
		if len(source)%size != 0 {
			continue
		}
		unit := source[:size]
		count := len(source) / size
		if count >= 2 && bytes.Equal(source, bytes.Repeat(unit, count)) {
			return ExactRepeat(source, unit)
		}
	}
	return ir.Representation{}, fmt.Errorf("origami fold: source has no exact repeating unit")
}
