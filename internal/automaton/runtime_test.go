package automaton

import "testing"

func propagationAutomaton() Automaton {
	return Automaton{
		Schema: Schema,
		ID:     "signal-chain",
		Cells: []Cell{
			{ID: "A", Kind: "SOURCE", InitialState: "ACTIVE", Neighbors: []string{"B"}},
			{ID: "B", Kind: "RELAY", InitialState: "IDLE", Neighbors: []string{"A", "C"}},
			{ID: "C", Kind: "SINK", InitialState: "IDLE", Neighbors: []string{"B"}},
		},
		Edges: []Edge{{From: "A", To: "B", Kind: "EMISSION"}, {From: "B", To: "C", Kind: "EMISSION"}},
		Rules: []Rule{
			{ID: "r-b-on-a", TargetCell: "B", FromState: "IDLE", ToState: "ACTIVE", Requires: []Predicate{{CellID: "A", State: "ACTIVE"}}},
			{ID: "r-a-done", TargetCell: "A", FromState: "ACTIVE", ToState: "DONE", Requires: []Predicate{{CellID: "B", State: "ACTIVE"}}},
			{ID: "r-c-on-b", TargetCell: "C", FromState: "IDLE", ToState: "ACTIVE", Requires: []Predicate{{CellID: "B", State: "ACTIVE"}}},
			{ID: "r-b-done", TargetCell: "B", FromState: "ACTIVE", ToState: "DONE", Requires: []Predicate{{CellID: "C", State: "ACTIVE"}}},
		},
	}
}

func TestPropagation(t *testing.T) {
	a := propagationAutomaton()
	s0 := a.Initial()
	s1, r1, err := Step(a, s0, 1)
	if err != nil { t.Fatal(err) }
	if s1["A"] != "ACTIVE" || s1["B"] != "ACTIVE" || s1["C"] != "IDLE" { t.Fatalf("unexpected s1: %#v", s1) }
	if r1.Stable { t.Fatal("step 1 unexpectedly stable") }

	s2, _, err := Step(a, s1, 2)
	if err != nil { t.Fatal(err) }
	if s2["A"] != "DONE" || s2["B"] != "ACTIVE" || s2["C"] != "ACTIVE" { t.Fatalf("unexpected s2: %#v", s2) }

	s3, _, err := Step(a, s2, 3)
	if err != nil { t.Fatal(err) }
	if s3["B"] != "DONE" { t.Fatalf("expected B DONE, got %#v", s3) }
}

func TestPriorityThenRuleID(t *testing.T) {
	a := Automaton{Schema: Schema, ID: "conflict", Cells: []Cell{{ID: "A", InitialState: "IDLE"}}, Rules: []Rule{
		{ID: "z", TargetCell: "A", FromState: "IDLE", ToState: "LOW", Priority: 1},
		{ID: "b", TargetCell: "A", FromState: "IDLE", ToState: "WIN", Priority: 2},
		{ID: "a", TargetCell: "A", FromState: "IDLE", ToState: "LEX", Priority: 2},
	}}
	next, result, err := Step(a, a.Initial(), 1)
	if err != nil { t.Fatal(err) }
	if next["A"] != "LEX" { t.Fatalf("expected lexicographic tie winner, got %q", next["A"]) }
	if len(result.Changes) != 1 || result.Changes[0].RuleID != "a" { t.Fatalf("unexpected result: %#v", result) }
}
