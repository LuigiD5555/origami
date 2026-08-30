package temporal

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/automaton"
)

func chainProgram() Program {
	a := automaton.Automaton{
		Schema: automaton.Schema,
		ID: "signal-chain",
		Cells: []automaton.Cell{
			{ID: "A", InitialState: "ACTIVE", Neighbors: []string{"B"}},
			{ID: "B", InitialState: "IDLE", Neighbors: []string{"A", "C"}},
			{ID: "C", InitialState: "IDLE", Neighbors: []string{"B"}},
		},
		Rules: []automaton.Rule{
			{ID: "r1", TargetCell: "B", FromState: "IDLE", ToState: "ACTIVE", Requires: []automaton.Predicate{{CellID: "A", State: "ACTIVE"}}},
			{ID: "r2", TargetCell: "A", FromState: "ACTIVE", ToState: "DONE", Requires: []automaton.Predicate{{CellID: "B", State: "ACTIVE"}}},
			{ID: "r3", TargetCell: "C", FromState: "IDLE", ToState: "ACTIVE", Requires: []automaton.Predicate{{CellID: "B", State: "ACTIVE"}}},
			{ID: "r4", TargetCell: "B", FromState: "ACTIVE", ToState: "DONE", Requires: []automaton.Predicate{{CellID: "C", State: "ACTIVE"}}},
		},
	}
	return Program{Schema: Schema, ID: "film-1", Automaton: a, MaxSteps: 8, CheckpointEvery: 2}
}

func TestRunSeekReplay(t *testing.T) {
	p := chainProgram()
	trace, err := Run(p)
	if err != nil { t.Fatal(err) }
	if trace.StableAt == 0 { t.Fatal("expected stable trace") }
	state2, err := Seek(trace, 2)
	if err != nil { t.Fatal(err) }
	if state2["A"] != "DONE" || state2["C"] != "ACTIVE" { t.Fatalf("unexpected step 2: %#v", state2) }
	if err := VerifyReplay(p, trace); err != nil { t.Fatal(err) }
}

func TestSeekDetectsCorruptDelta(t *testing.T) {
	p := chainProgram()
	trace, err := Run(p)
	if err != nil { t.Fatal(err) }
	trace.Checkpoints = nil
	if len(trace.Frames) == 0 || len(trace.Frames[0].Changes) == 0 { t.Fatal("missing delta") }
	trace.Frames[0].Changes[0].From = "CORRUPT"
	if _, err := Seek(trace, 1); err == nil { t.Fatal("expected corrupt delta rejection") }
}

func TestMeasure(t *testing.T) {
	trace, err := Run(chainProgram())
	if err != nil { t.Fatal(err) }
	m := Measure(trace)
	if m.Steps == 0 || m.Changes == 0 || m.FoldedTraceByte == 0 { t.Fatalf("bad metrics: %#v", m) }
}
