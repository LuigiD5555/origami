package machine

import (
	"bytes"
	"fmt"

	"github.com/LuigiD5555/origami/internal/origami/ir"
)

// Budget bounds materialization so generative representations cannot expand
// without an explicit finite limit.
type Budget struct {
	MaxBytes uint64
	MaxNodes uint64
	MaxDepth uint64
}

// Stats records the work performed by an unfold operation.
type Stats struct {
	BytesMaterialized uint64
	NodesVisited      uint64
	MaxDepth          uint64
}

// UnfoldExact deterministically materializes the representation root.
func UnfoldExact(r ir.Representation, budget Budget) ([]byte, Stats, error) {
	if err := r.Validate(); err != nil {
		return nil, Stats{}, err
	}
	if budget.MaxBytes == 0 || budget.MaxNodes == 0 || budget.MaxDepth == 0 {
		return nil, Stats{}, fmt.Errorf("origami machine: finite non-zero budget required")
	}

	var out bytes.Buffer
	var stats Stats
	var emit func(string, uint64) error
	emit = func(id string, depth uint64) error {
		if depth > budget.MaxDepth {
			return fmt.Errorf("origami machine: max depth exceeded")
		}
		stats.NodesVisited++
		if stats.NodesVisited > budget.MaxNodes {
			return fmt.Errorf("origami machine: max nodes exceeded")
		}
		if depth > stats.MaxDepth {
			stats.MaxDepth = depth
		}

		n := r.Nodes[id]
		switch n.Kind {
		case ir.KindLiteral:
			if uint64(out.Len())+uint64(len(n.Literal)) > budget.MaxBytes {
				return fmt.Errorf("origami machine: max bytes exceeded")
			}
			_, _ = out.Write(n.Literal)
		case ir.KindRef:
			return emit(n.Target, depth+1)
		case ir.KindConcat:
			for _, child := range n.Children {
				if err := emit(child, depth+1); err != nil {
					return err
				}
			}
		case ir.KindRepeat:
			for i := uint64(0); i < n.Count; i++ {
				if err := emit(n.Target, depth+1); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("origami machine: unsupported node kind %q", n.Kind)
		}
		return nil
	}

	if err := emit(r.Root, 1); err != nil {
		return nil, stats, err
	}
	stats.BytesMaterialized = uint64(out.Len())
	return out.Bytes(), stats, nil
}
