package ir

import "fmt"

// Validate rejects malformed or cyclic R0 representations before execution.
func (r Representation) Validate() error {
	if r.Root == "" {
		return fmt.Errorf("origami ir: empty root")
	}
	if _, ok := r.Nodes[r.Root]; !ok {
		return fmt.Errorf("origami ir: root %q not found", r.Root)
	}

	visiting := make(map[string]bool)
	visited := make(map[string]bool)
	var visit func(string) error
	visit = func(id string) error {
		if visiting[id] {
			return fmt.Errorf("origami ir: cycle at %q", id)
		}
		if visited[id] {
			return nil
		}
		n, ok := r.Nodes[id]
		if !ok {
			return fmt.Errorf("origami ir: node %q not found", id)
		}
		if !n.Kind.Valid() {
			return fmt.Errorf("origami ir: node %q has unsupported kind %q", id, n.Kind)
		}
		if n.ID != id {
			return fmt.Errorf("origami ir: node key %q does not match id %q", id, n.ID)
		}

		visiting[id] = true
		var deps []string
		switch n.Kind {
		case KindLiteral:
			if n.Target != "" || len(n.Children) != 0 || n.Count != 0 {
				return fmt.Errorf("origami ir: literal %q has invalid operands", id)
			}
		case KindRef:
			if n.Target == "" || len(n.Literal) != 0 || len(n.Children) != 0 || n.Count != 0 {
				return fmt.Errorf("origami ir: ref %q requires only a target", id)
			}
			deps = []string{n.Target}
		case KindConcat:
			if len(n.Children) == 0 || len(n.Literal) != 0 || n.Target != "" || n.Count != 0 {
				return fmt.Errorf("origami ir: concat %q requires only children", id)
			}
			deps = n.Children
		case KindRepeat:
			if n.Target == "" || n.Count == 0 || len(n.Literal) != 0 || len(n.Children) != 0 {
				return fmt.Errorf("origami ir: repeat %q requires only target and positive count", id)
			}
			deps = []string{n.Target}
		}
		for _, dep := range deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		return nil
	}
	return visit(r.Root)
}
