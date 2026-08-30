package core

import (
	"reflect"
	"testing"
)

func present(data string) Value { return Value{Status: StatusPresent, Data: data} }

func TestADV001HigherOrderInteractionRequiresExactParticipantSet(t *testing.T) {
	machine := NewMachine([]Rule{{
		ID: "triple",
		When: []Predicate{
			{Scope: "state", Key: "A", Op: "present"},
			{Scope: "state", Key: "B", Op: "present"},
			{Scope: "state", Key: "C", Op: "present"},
		},
		Then: []Effect{{Key: "emergent", Mode: "set", Status: StatusPresent, Data: "ABC"}},
	}}, 4)
	all := State{Values: map[string]Value{"A": present("1"), "B": present("1"), "C": present("1")}}
	trajectory, err := Execute(all, machine, []Context{{ID: "c", Values: map[string]Value{}}})
	if err != nil { t.Fatal(err) }
	if got := trajectory.States[1].Values["emergent"]; got.Status != StatusPresent || got.Data != "ABC" {
		t.Fatalf("triple interaction missing: %+v", got)
	}
	obs := Observe(trajectory, ObservationContract{ID: "participants", Kind: "participant_set", Budget: 1, Participants: []string{"state:A", "state:B", "state:C"}})
	if obs.Status != ObservationPass { t.Fatalf("expected exact triple participant evidence: %+v", obs) }

	pair := State{Values: map[string]Value{"A": present("1"), "B": present("1"), "C": {Status: StatusAbsent}}}
	pairTrajectory, err := Execute(pair, machine, []Context{{ID: "c", Values: map[string]Value{}}})
	if err != nil { t.Fatal(err) }
	if _, ok := pairTrajectory.States[1].Values["emergent"]; ok {
		t.Fatalf("pairwise state must not promote triple behavior: %+v", pairTrajectory.States[1])
	}
}

func TestADV002ContextParticipatesInTransition(t *testing.T) {
	machine := NewMachine([]Rule{
		{ID: "alpha", When: []Predicate{{Scope: "context", Key: "mode", Op: "data_eq", Data: "alpha"}}, Then: []Effect{{Key: "out", Mode: "set", Status: StatusPresent, Data: "A"}}},
		{ID: "beta", When: []Predicate{{Scope: "context", Key: "mode", Op: "data_eq", Data: "beta"}}, Then: []Effect{{Key: "out", Mode: "set", Status: StatusPresent, Data: "B"}}},
	}, 4)
	initial := State{Values: map[string]Value{"seed": present("same")}}
	a, err := Execute(initial, machine, []Context{{ID: "C_alpha", Values: map[string]Value{"mode": present("alpha")}}})
	if err != nil { t.Fatal(err) }
	b, err := Execute(initial, machine, []Context{{ID: "C_beta", Values: map[string]Value{"mode": present("beta")}}})
	if err != nil { t.Fatal(err) }
	if StateHash(a.Initial) != StateHash(b.Initial) { t.Fatal("S0 changed across context experiment") }
	if StateHash(a.States[1]) == StateHash(b.States[1]) { t.Fatal("different contexts must deterministically change evolution") }
	if a.Evidence[0].ContextID != "C_alpha" || b.Evidence[0].ContextID != "C_beta" { t.Fatalf("context provenance missing: %+v %+v", a.Evidence[0], b.Evidence[0]) }
}

func TestADV003RulesAreDataWithStableProvenance(t *testing.T) {
	initial := State{Values: map[string]Value{"x": present("go")}}
	machineA := NewMachine([]Rule{{ID: "rule-A", When: []Predicate{{Scope: "state", Key: "x", Op: "present"}}, Then: []Effect{{Key: "out", Mode: "set", Status: StatusPresent, Data: "A"}}}}, 2)
	machineB := NewMachine([]Rule{{ID: "rule-B", When: []Predicate{{Scope: "state", Key: "x", Op: "present"}}, Then: []Effect{{Key: "out", Mode: "set", Status: StatusPresent, Data: "B"}}}}, 2)
	a, err := Execute(initial, machineA, []Context{{ID: "same", Values: map[string]Value{}}}); if err != nil { t.Fatal(err) }
	b, err := Execute(initial, machineB, []Context{{ID: "same", Values: map[string]Value{}}}); if err != nil { t.Fatal(err) }
	if a.States[1].Values["out"].Data == b.States[1].Values["out"].Data { t.Fatal("changing declared R must change trace without code changes") }
	if a.Evidence[0].FiredRules[0].RuleID != "rule-A" || b.Evidence[0].FiredRules[0].RuleID != "rule-B" { t.Fatal("transition must cite fired rule ids") }
}

func TestADV004SemanticStatusesRemainDistinct(t *testing.T) {
	state := State{Values: map[string]Value{
		"absent":    {Status: StatusAbsent},
		"unknown":   {Status: StatusUnknown},
		"inhibited": {Status: StatusInhibited, Data: "blocked"},
		"cancelled": {Status: StatusCancelled, Data: "signal"},
	}}
	seen := map[Status]bool{}
	for _, value := range state.Values { seen[value.Status] = true }
	if len(seen) != 4 { t.Fatalf("semantic statuses collapsed: %+v", state) }
	for key, value := range state.Values {
		obsTrajectory := Trajectory{Schema: SchemaR1 + ".trajectory", Initial: CloneState(state), States: []State{CloneState(state)}, Terminated: "TEST"}
		obsTrajectory.TraceHash = trajectoryHash(obsTrajectory)
		predicate := Predicate{Scope: "state", Key: key, Op: "status_eq", Status: value.Status}
		obs := Observe(obsTrajectory, ObservationContract{ID: key, Kind: "state_predicate", Budget: 1, StateStep: 0, Predicate: &predicate})
		if obs.Status != ObservationPass { t.Fatalf("status %s not observable for %s: %+v", value.Status, key, obs) }
	}
}

func TestADV005ObservationDoesNotRewriteTransition(t *testing.T) {
	machine := NewMachine([]Rule{{ID: "derive", When: []Predicate{{Scope: "state", Key: "x", Op: "present"}}, Then: []Effect{{Key: "y", Mode: "set", Status: StatusPresent, Data: "derived"}}}}, 2)
	trajectory, err := Execute(State{Values: map[string]Value{"x": present("1")}}, machine, []Context{{ID: "c", Values: map[string]Value{}}})
	if err != nil { t.Fatal(err) }
	before := trajectory.TraceHash
	predicate := Predicate{Scope: "state", Key: "y", Op: "data_eq", Data: "derived"}
	one := Observe(trajectory, ObservationContract{ID: "state", Kind: "state_predicate", Budget: 1, StateStep: 1, Predicate: &predicate})
	two := Observe(trajectory, ObservationContract{ID: "rule", Kind: "rule_fired", Budget: 1, RuleID: "derive"})
	if one.Status != ObservationPass || two.Status != ObservationPass { t.Fatalf("observations failed: %+v %+v", one, two) }
	if trajectory.TraceHash != before || trajectory.TraceHash != trajectoryHash(trajectory) { t.Fatal("observer mutated pre-observation trajectory") }
}

func TestADV006FoldPreservesUnresolvedAlternatives(t *testing.T) {
	graph := SemanticGraph{Nodes: map[string]SemanticNode{
		"root": {Address: "root", Kind: "choice", Label: "unresolved", Value: Value{Status: StatusUnknown}, Alternatives: []string{"a", "b"}},
		"a":    {Address: "a", Kind: "alternative", Label: "A", Value: present("A")},
		"b":    {Address: "b", Kind: "alternative", Label: "B", Value: present("B")},
	}}
	folded, store, err := FoldGraph(graph); if err != nil { t.Fatal(err) }
	if !VerifyFold(folded, store) { t.Fatal("fold commitment failed") }
	store.ResetTouches()
	partial, err := SelectiveUnfold(folded, store, []string{"root"}); if err != nil { t.Fatal(err) }
	if !reflect.DeepEqual(partial.Addresses, []string{"a", "b", "root"}) { t.Fatalf("alternatives lost during selective unfold: %+v", partial.Addresses) }
	if partial.Closure[2].Value.Status != StatusUnknown { t.Fatalf("Fold resolved an alternative without observation: %+v", partial.Closure[2]) }
}

func TestADV007SelectiveUnfoldReadsOnlyDependencyClosure(t *testing.T) {
	graph := SemanticGraph{Nodes: map[string]SemanticNode{
		"region-a": {Address: "region-a", Value: present("A"), Dependencies: []string{"shared"}},
		"shared":   {Address: "shared", Value: present("needed")},
		"region-b": {Address: "region-b", Value: present("B")},
	}}
	folded, store, err := FoldGraph(graph); if err != nil { t.Fatal(err) }
	store.ResetTouches()
	partial, err := SelectiveUnfold(folded, store, []string{"region-a"}); if err != nil { t.Fatal(err) }
	if !reflect.DeepEqual(partial.Addresses, []string{"region-a", "shared"}) { t.Fatalf("wrong dependency closure: %+v", partial.Addresses) }
	if !reflect.DeepEqual(partial.Touched, []string{"region-a", "shared"}) { t.Fatalf("selective unfold touched unrelated memory: %+v", partial.Touched) }
	store.ResetTouches()
	full, err := FullUnfold(folded, store); if err != nil { t.Fatal(err) }
	selected := map[string]SemanticNode{}; for _, node := range full.Closure { if node.Address == "region-a" || node.Address == "shared" { selected[node.Address] = node } }
	if len(selected) != 2 || selected["region-a"].Value.Data != partial.Closure[0].Value.Data { t.Fatal("selective unfold differs from equivalent full-unfold selection") }
}

func TestADV008TemporalLatentPropertyRequiresOrderedTrajectory(t *testing.T) {
	machine := NewMachine([]Rule{
		{ID: "phase-A", When: []Predicate{{Scope: "context", Key: "phase", Op: "data_eq", Data: "A"}}, Then: []Effect{{Key: "phase", Mode: "set", Status: StatusPresent, Data: "A"}}},
		{ID: "phase-B", When: []Predicate{{Scope: "context", Key: "phase", Op: "data_eq", Data: "B"}}, Then: []Effect{{Key: "phase", Mode: "set", Status: StatusPresent, Data: "B"}}},
		{ID: "phase-C", When: []Predicate{{Scope: "context", Key: "phase", Op: "data_eq", Data: "C"}}, Then: []Effect{{Key: "phase", Mode: "set", Status: StatusPresent, Data: "C"}}},
	}, 4)
	initial := State{Values: map[string]Value{"phase": {Status: StatusUnknown}}}
	ordered, err := Execute(initial, machine, []Context{{ID: "1", Values: map[string]Value{"phase": present("A")}}, {ID: "2", Values: map[string]Value{"phase": present("B")}}, {ID: "3", Values: map[string]Value{"phase": present("C")}}}); if err != nil { t.Fatal(err) }
	contract := ObservationContract{ID: "ABC", Kind: "temporal_sequence", Budget: 4, Sequence: []TemporalAtom{{Key: "phase", Data: "A", Status: StatusPresent}, {Key: "phase", Data: "B", Status: StatusPresent}, {Key: "phase", Data: "C", Status: StatusPresent}}}
	if got := Observe(ordered, contract); got.Status != ObservationPass { t.Fatalf("ordered temporal property not detected: %+v", got) }
	for i, state := range ordered.States { if state.Values["phase"].Data == "ABC" { t.Fatalf("temporal answer smuggled into state %d", i) } }
	reordered, err := Execute(initial, machine, []Context{{ID: "1", Values: map[string]Value{"phase": present("A")}}, {ID: "2", Values: map[string]Value{"phase": present("C")}}, {ID: "3", Values: map[string]Value{"phase": present("B")}}}); if err != nil { t.Fatal(err) }
	if got := Observe(reordered, contract); got.Status != ObservationFail { t.Fatalf("reordered trajectory must destroy temporal property: %+v", got) }
}

func TestADV009CancellationIsOrderIndependentAndKeepsContributions(t *testing.T) {
	plus := Rule{ID: "plus", When: []Predicate{{Scope: "state", Key: "trigger", Op: "present"}}, Then: []Effect{{Key: "signal", Mode: "contribute", Data: "wave", Delta: 1}}}
	minus := Rule{ID: "minus", When: []Predicate{{Scope: "state", Key: "trigger", Op: "present"}}, Then: []Effect{{Key: "signal", Mode: "contribute", Data: "wave", Delta: -1}}}
	initial := State{Values: map[string]Value{"trigger": present("on")}}
	a, err := Execute(initial, NewMachine([]Rule{plus, minus}, 2), []Context{{ID: "c", Values: map[string]Value{}}}); if err != nil { t.Fatal(err) }
	b, err := Execute(initial, NewMachine([]Rule{minus, plus}, 2), []Context{{ID: "c", Values: map[string]Value{}}}); if err != nil { t.Fatal(err) }
	if a.States[1].Values["signal"].Status != StatusCancelled || StateHash(a.States[1]) != StateHash(b.States[1]) { t.Fatalf("cancellation depends on declaration order: %+v %+v", a.States[1], b.States[1]) }
	resolution := a.Evidence[0].Resolutions[0]
	if resolution.Kind != "cancellation" || len(resolution.Contributions) != 2 { t.Fatalf("cancellation evidence lost: %+v", resolution) }
	onlyPlus, err := Execute(initial, NewMachine([]Rule{plus}, 2), []Context{{ID: "c", Values: map[string]Value{}}}); if err != nil { t.Fatal(err) }
	if onlyPlus.States[1].Values["signal"].Status != StatusPresent { t.Fatal("removing one contribution must change result") }
}

func TestADV010And011LiteralTransportCannotClaimSemanticComputation(t *testing.T) {
	transportOnly := RepresentationAccounting{SourceBytes: 1000, LiteralBytes: 1000}
	if report := EvaluateAccounting(ClaimSemantic, transportOnly); report.Allowed { t.Fatalf("literal transport falsely earned semantic authority: %+v", report) }
	if report := EvaluateAccounting(ClaimExact, transportOnly); !report.Allowed { t.Fatalf("explicit literal transport may support exact recovery accounting: %+v", report) }
	semantic := RepresentationAccounting{SourceBytes: 1000, SemanticBytes: 120, LiteralBytes: 1000, MachineDerivedClaims: 3, VerifiedSemanticClaims: 2}
	if report := EvaluateAccounting(ClaimSemantic, semantic); !report.Allowed { t.Fatalf("verified semantic representation rejected: %+v", report) }
	withoutExact := semantic; withoutExact.LiteralBytes = 0
	if report := EvaluateAccounting(ClaimSemantic, withoutExact); !report.Allowed { t.Fatalf("semantic state must survive removal of exact transport: %+v", report) }
}

func TestDeterministicReplayProducesSameTraceHash(t *testing.T) {
	machine := NewMachine([]Rule{{ID: "r", When: []Predicate{{Scope: "state", Key: "x", Op: "present"}}, Then: []Effect{{Key: "y", Mode: "set", Status: StatusPresent, Data: "1"}}}}, 3)
	initial := State{Values: map[string]Value{"x": present("1")}}
	contexts := []Context{{ID: "c", Values: map[string]Value{}}}
	a, err := Execute(initial, machine, contexts); if err != nil { t.Fatal(err) }
	b, err := Execute(initial, machine, contexts); if err != nil { t.Fatal(err) }
	if a.TraceHash != b.TraceHash || !reflect.DeepEqual(a, b) { t.Fatalf("deterministic replay drifted\nA=%+v\nB=%+v", a, b) }
}
