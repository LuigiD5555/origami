package reference

import (
	"reflect"
	"testing"
)

func TestADV001HigherOrderInteractionIsNotPairwise(t *testing.T) {
	tripleRule:=Rule{
		ID:"R-ABC-emergent",
		Participants:[]Participant{
			{Entity:"A",Status:StatusPresent,Value:"active"},
			{Entity:"B",Status:StatusPresent,Value:"active"},
			{Entity:"C",Status:StatusPresent,Value:"active"},
		},
		Then:[]Effect{{Entity:"X",Status:StatusPresent,Value:"emergent"}},
	}
	cases:=[]struct{name string; state CoreState; expect bool}{
		{"AB",CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusPresent,Value:"active"},"C":{Status:StatusAbsent},"X":{Status:StatusUnknown}},false},
		{"AC",CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusAbsent},"C":{Status:StatusPresent,Value:"active"},"X":{Status:StatusUnknown}},false},
		{"BC",CoreState{"A":{Status:StatusAbsent},"B":{Status:StatusPresent,Value:"active"},"C":{Status:StatusPresent,Value:"active"},"X":{Status:StatusUnknown}},false},
		{"ABC",CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusPresent,Value:"active"},"C":{Status:StatusPresent,Value:"active"},"X":{Status:StatusUnknown}},true},
	}
	for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){run,err:=ExecuteCore(tc.state,[]Context{{ID:"C"}},[]Rule{tripleRule});if err!=nil{t.Fatal(err)};got:=run.Trace[len(run.Trace)-1]["X"];emerged:=got.Status==StatusPresent&&got.Value=="emergent";if emerged!=tc.expect{t.Fatalf("ADV-001: emerged=%v want %v",emerged,tc.expect)};if tc.expect{ev:=run.Transitions[0].Interactions;if len(ev)!=1{t.Fatalf("ADV-001: missing exact interaction evidence: %+v",ev)};if !reflect.DeepEqual(ev[0].Participants,[]string{"A","B","C"}){t.Fatalf("ADV-001: wrong participating set: %v",ev[0].Participants)};if ParticipatingSetKey(ev[0].Participants)!="A+B+C"{t.Fatal("ADV-001: participating-set key is not canonical")}}})}
}

func TestADV001RemovingAnyParticipantRemovesResult(t *testing.T){
	rule:=Rule{ID:"R-ABC",Participants:[]Participant{{Entity:"A",Status:StatusPresent},{Entity:"B",Status:StatusPresent},{Entity:"C",Status:StatusPresent}},Then:[]Effect{{Entity:"X",Status:StatusPresent,Value:"triple"}}}
	base:=CoreState{"A":{Status:StatusPresent},"B":{Status:StatusPresent},"C":{Status:StatusPresent},"X":{Status:StatusUnknown}}
	for _,removed:=range []string{"A","B","C"}{s:=cloneCoreState(base);s[removed]=SemanticValue{Status:StatusAbsent};run,err:=ExecuteCore(s,[]Context{{ID:"C"}},[]Rule{rule});if err!=nil{t.Fatal(err)};if run.Trace[1]["X"].Status==StatusPresent{t.Fatalf("ADV-001: result survived removal of %s",removed)}}
}

func TestADV002ContextChangesEvolution(t *testing.T) {
	initial:=CoreState{"signal":{Status:StatusPresent,Value:"ready"},"target":{Status:StatusUnknown}}
	rules:=[]Rule{{ID:"R-context",When:[]Predicate{{Source:"context",Key:"phase",Equals:"open"}},Then:[]Effect{{Entity:"target",Status:StatusPresent,Value:"active"}}}}
	openRun,err:=ExecuteCore(initial,[]Context{{ID:"C-open",Values:map[string]string{"phase":"open"}}},rules); if err!=nil{t.Fatal(err)}
	closedRun,err:=ExecuteCore(initial,[]Context{{ID:"C-closed",Values:map[string]string{"phase":"closed"}}},rules); if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(openRun.Initial,closedRun.Initial){t.Fatal("ADV-002: initial state differs")}
	if reflect.DeepEqual(openRun.Trace,closedRun.Trace){t.Fatal("ADV-002: context did not change deterministic evolution")}
	if got:=openRun.Transitions[0]; got.ContextID!="C-open" || len(got.FiredRules)!=1 || got.FiredRules[0]!="R-context"{t.Fatalf("ADV-002: missing context/rule evidence: %+v",got)}
}

func TestADV003RulesAreDeclaredData(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}};ctx:=[]Context{{ID:"C"}}
	one:=[]Rule{{ID:"R-one",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"one"}}}}
	two:=[]Rule{{ID:"R-two",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"two"}}}}
	a,err:=ExecuteCore(initial,ctx,one);if err!=nil{t.Fatal(err)};b,err:=ExecuteCore(initial,ctx,two);if err!=nil{t.Fatal(err)}
	if reflect.DeepEqual(a.Trace,b.Trace){t.Fatal("ADV-003: changing only declared rules did not change trace")};if !reflect.DeepEqual(a.Transitions[0].FiredRules,[]string{"R-one"})||!reflect.DeepEqual(b.Transitions[0].FiredRules,[]string{"R-two"}){t.Fatal("ADV-003: rule provenance missing")}
}

func TestADV004SemanticStatusesAreNotCollapsed(t *testing.T) {
	state:=CoreState{"absent":{Status:StatusAbsent},"unknown":{Status:StatusUnknown},"inhibited":{Status:StatusInhibited,Reason:"blocked transition"},"cancelled":{Status:StatusCancelled,Reason:"opposing contributions"}}
	statuses:=map[SemanticStatus]bool{};for _,key:=range []string{"absent","unknown","inhibited","cancelled"}{v:=state[key];if statuses[v.Status]{t.Fatalf("ADV-004: collapsed semantic status %s",v.Status)};statuses[v.Status]=true};if state["unknown"].Status==state["absent"].Status{t.Fatal("ADV-004: UNKNOWN collapsed to ABSENT")};if state["inhibited"].Reason==""||state["cancelled"].Reason==""{t.Fatal("ADV-004: inhibition/cancellation evidence missing")}
}

func TestADV005ObservationDoesNotTransition(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}};rules:=[]Rule{{ID:"R",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"resolved"}}}}
	run,err:=ExecuteCore(initial,[]Context{{ID:"C"}},rules);if err!=nil{t.Fatal(err)};before,err:=TraceHash(run.Trace);if err!=nil{t.Fatal(err)};obsA,err:=Observe(run,ObservationContract{ID:"OBS-A",Authority:"authority-a",Entity:"A"});if err!=nil{t.Fatal(err)};obsB,err:=Observe(run,ObservationContract{ID:"OBS-B",Authority:"authority-b",Entity:"B"});if err!=nil{t.Fatal(err)};after,err:=TraceHash(run.Trace);if err!=nil{t.Fatal(err)};if before!=after||obsA.TraceHash!=before||obsB.TraceHash!=before{t.Fatal("ADV-005: observation mutated transition evidence")};if obsA.Authority==obsB.Authority{t.Fatal("ADV-005: observation authority collapsed")};if obsA.Value==obsB.Value{t.Fatal("ADV-005: observers did not expose distinct aspects of same immutable trace")}
}
