package machine

import (
	"fmt"

	"github.com/LuigiD5555/origami/internal/origami/ir"
)

// UnfoldAddress materializes exactly one address and its dependency closure.
// It is the first selective-unfold primitive: callers need not unfold Root.
func UnfoldAddress(r ir.Representation, address string, budget Budget) ([]byte, Stats, error) {
	if err := r.Validate(); err != nil {
		return nil, Stats{}, err
	}
	if address == "" {
		return nil, Stats{}, fmt.Errorf("origami machine: empty address")
	}
	if _, ok := r.Nodes[address]; !ok {
		return nil, Stats{}, fmt.Errorf("origami machine: address %q not found", address)
	}
	selected := r
	selected.Root = address
	return UnfoldExact(selected, budget)
}
