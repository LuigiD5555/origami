package temporal

import "testing"

func TestCausalTraceFrontiersPreserveParallelBranches(t *testing.T) {
	trace := CausalTrace{
		Schema:           CausalTraceSchema,
		ID:               "run-1",
		InitialStateHash: "s0",
		Events: []CausalEvent{
			{ID: "a", Type: EventStart, LogicalClock: 1},
			{ID: "b", Type: EventStart, LogicalClock: 1},
			{ID: "c", Type: EventComplete, ParentIDs: []string{"a"}, LogicalClock: 2},
			{ID: "d", Type: EventCommit, ParentIDs: []string{"b", "c"}, LogicalClock: 3, OutputStateHash: "s1"},
		},
		Checkpoints:    []CausalCheckpoint{{EventID: "d", StateHash: "s1"}},
		FinalStateHash: "s1",
	}
	if err := trace.Validate(); err != nil {
		t.Fatal(err)
	}
	frontiers, err := trace.Frontiers()
	if err != nil {
		t.Fatal(err)
	}
	if len(frontiers) != 3 {
		t.Fatalf("frontiers=%d", len(frontiers))
	}
	if len(frontiers[0]) != 2 || frontiers[0][0].ID != "a" || frontiers[0][1].ID != "b" {
		t.Fatalf("frontier0=%+v", frontiers[0])
	}
	if len(frontiers[1]) != 1 || frontiers[1][0].ID != "c" {
		t.Fatalf("frontier1=%+v", frontiers[1])
	}
	if len(frontiers[2]) != 1 || frontiers[2][0].ID != "d" {
		t.Fatalf("frontier2=%+v", frontiers[2])
	}
}

func TestCausalTraceRejectsCycle(t *testing.T) {
	trace := CausalTrace{
		ID:               "cycle",
		InitialStateHash: "s0",
		Events: []CausalEvent{
			{ID: "a", Type: EventStart, ParentIDs: []string{"b"}},
			{ID: "b", Type: EventComplete, ParentIDs: []string{"a"}},
		},
	}
	if err := trace.Validate(); err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestCausalTraceRejectsUnknownParent(t *testing.T) {
	trace := CausalTrace{ID: "bad-parent", InitialStateHash: "s0", Events: []CausalEvent{{ID: "a", Type: EventStart, ParentIDs: []string{"missing"}}}}
	if err := trace.Validate(); err == nil {
		t.Fatal("expected unknown parent error")
	}
}
