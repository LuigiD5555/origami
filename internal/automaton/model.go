package automaton

import (
	"fmt"
	"sort"
)

const Schema = "origami.automaton.r0"

type Cell struct {
	ID           string   `json:"id"`
	Kind         string   `json:"kind,omitempty"`
	InitialState string   `json:"initial_state"`
	Neighbors    []string `json:"neighbors,omitempty"`
}

type Predicate struct {
	CellID string `json:"cell_id"`
	State  string `json:"state"`
}

type Rule struct {
	ID         string      `json:"id"`
	TargetCell string      `json:"target_cell"`
	FromState  string      `json:"from_state,omitempty"`
	ToState    string      `json:"to_state"`
	Requires   []Predicate `json:"requires,omitempty"`
	Priority   int         `json:"priority,omitempty"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind,omitempty"`
}

type Automaton struct {
	Schema string `json:"schema"`
	ID     string `json:"id"`
	Cells  []Cell `json:"cells"`
	Rules  []Rule `json:"rules"`
	Edges  []Edge `json:"edges,omitempty"`
}

type State map[string]string

type Change struct {
	CellID string `json:"cell_id"`
	From   string `json:"from"`
	To     string `json:"to"`
	RuleID string `json:"rule_id"`
}

type StepResult struct {
	Step    int      `json:"step"`
	Changes []Change `json:"changes"`
	Stable  bool     `json:"stable"`
}

func (a Automaton) Validate() error {
	if a.Schema != "" && a.Schema != Schema {
		return fmt.Errorf("unsupported automaton schema %q", a.Schema)
	}
	if a.ID == "" {
		return fmt.Errorf("automaton id is required")
	}
	cells := map[string]Cell{}
	for _, c := range a.Cells {
		if c.ID == "" || c.InitialState == "" {
			return fmt.Errorf("cell id and initial_state are required")
		}
		if _, ok := cells[c.ID]; ok {
			return fmt.Errorf("duplicate cell %q", c.ID)
		}
		cells[c.ID] = c
	}
	for _, c := range a.Cells {
		for _, n := range c.Neighbors {
			if _, ok := cells[n]; !ok {
				return fmt.Errorf("cell %q references unknown neighbor %q", c.ID, n)
			}
		}
	}
	seenRules := map[string]bool{}
	for _, r := range a.Rules {
		if r.ID == "" || r.TargetCell == "" || r.ToState == "" {
			return fmt.Errorf("rule id, target_cell and to_state are required")
		}
		if seenRules[r.ID] {
			return fmt.Errorf("duplicate rule %q", r.ID)
		}
		seenRules[r.ID] = true
		if _, ok := cells[r.TargetCell]; !ok {
			return fmt.Errorf("rule %q targets unknown cell %q", r.ID, r.TargetCell)
		}
		for _, p := range r.Requires {
			if _, ok := cells[p.CellID]; !ok {
				return fmt.Errorf("rule %q requires unknown cell %q", r.ID, p.CellID)
			}
		}
	}
	for _, e := range a.Edges {
		if _, ok := cells[e.From]; !ok {
			return fmt.Errorf("edge references unknown from cell %q", e.From)
		}
		if _, ok := cells[e.To]; !ok {
			return fmt.Errorf("edge references unknown to cell %q", e.To)
		}
	}
	return nil
}

func (a Automaton) Initial() State {
	out := State{}
	for _, c := range a.Cells {
		out[c.ID] = c.InitialState
	}
	return out
}

func CloneState(in State) State {
	out := make(State, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func SortedCellIDs(s State) []string {
	ids := make([]string, 0, len(s))
	for id := range s {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
