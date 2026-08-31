package automaton

import (
	"fmt"
	"sort"
)

func applicable(r Rule, state State) bool {
	cur, ok := state[r.TargetCell]
	if !ok {
		return false
	}
	if r.FromState != "" && cur != r.FromState {
		return false
	}
	for _, p := range r.Requires {
		if state[p.CellID] != p.State {
			return false
		}
	}
	return true
}

func betterRule(a, b Rule) bool {
	if a.Priority != b.Priority {
		return a.Priority > b.Priority
	}
	return a.ID < b.ID
}

// Step applies all selected local rules synchronously. Selection is computed
// entirely from the input state. If multiple rules target one cell, the rule
// with the highest priority wins; ties are broken lexicographically by rule ID.
func Step(a Automaton, state State, step int) (State, StepResult, error) {
	if err := a.Validate(); err != nil {
		return nil, StepResult{}, err
	}
	for _, c := range a.Cells {
		if _, ok := state[c.ID]; !ok {
			return nil, StepResult{}, fmt.Errorf("state missing cell %q", c.ID)
		}
	}

	chosen := map[string]Rule{}
	for _, r := range a.Rules {
		if !applicable(r, state) {
			continue
		}
		prev, ok := chosen[r.TargetCell]
		if !ok || betterRule(r, prev) {
			chosen[r.TargetCell] = r
		}
	}

	next := CloneState(state)
	changes := make([]Change, 0, len(chosen))
	for target, r := range chosen {
		from := state[target]
		if from == r.ToState {
			continue
		}
		next[target] = r.ToState
		changes = append(changes, Change{CellID: target, From: from, To: r.ToState, RuleID: r.ID})
	}
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].CellID != changes[j].CellID {
			return changes[i].CellID < changes[j].CellID
		}
		return changes[i].RuleID < changes[j].RuleID
	})
	return next, StepResult{Step: step, Changes: changes, Stable: len(changes) == 0}, nil
}

func EqualState(a, b State) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
