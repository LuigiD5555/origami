package reference

import (
	"reflect"
	"testing"
)

// ADV-001 remains explicitly unsupported until higher-order participating-set
// semantics exist. ADV-002..005 are now executable semantic gates.
type adversarialCapability struct {
	ID        string
	Supported bool
	Reason    string
}

func requireUnsupported(t *testing.T, c adversarialCapability) {
	t.Helper()
	if c.Supported { t.Fatalf("%s unexpectedly marked supported before its semantic gate exists", c.ID) }
	if c.Reason == "" { t.Fatalf("%s must record why the current engine is unsupported", c.ID) }
	t.Logf("%s: UNSUPPORTED: %s", c.ID, c.Reason)
}

func TestADV001HigherOrderInteractionUnsupported(t *testing.T) {
	requireUnsupported(t, adversarialCapability{ID:"ADV-001",Reason:"formal core does not yet model participating-set/higher-order conditions"})
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
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}}
	ctx:=[]Context{{ID:"C"}}
	one:=[]Rule{{ID:"R-one",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"one"}}}}
	two:=[]Rule{{ID:"R-two",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"two"}}}}
	a,err:=ExecuteCore(initial,ctx,one); if err!=nil{t.Fatal(err)}
	b,err:=ExecuteCore(initial,ctx,two); if err!=nil{t.Fatal(err)}
	if reflect.DeepEqual(a.Trace,b.Trace){t.Fatal("ADV-003: changing only declared rules did not change trace")}
	if !reflect.DeepEqual(a.Transitions[0].FiredRules,[]string{"R-one"}) || !reflect.DeepEqual(b.Transitions[0].FiredRules,[]string{"R-two"}){t.Fatal("ADV-003: rule provenance missing")}
}

func TestADV004SemanticStatusesAreNotCollapsed(t *testing.T) {
	state:=CoreState{
		"absent":{Status:StatusAbsent},
		"unknown":{Status:StatusUnknown},
		"inhibited":{Status:StatusInhibited,Reason:"blocked transition"},
		"cancelled":{Status:StatusCancelled,Reason:"opposing contributions"},
	}
	statuses:=map[SemanticStatus]bool{}
	for _,key:=range []string{"absent","unknown","inhibited","cancelled"}{ v:=state[key]; if statuses[v.Status]{t.Fatalf("ADV-004: collapsed semantic status %s",v.Status)}; statuses[v.Status]=true }
	if state["unknown"].Status==state["absent"].Status{t.Fatal("ADV-004: UNKNOWN collapsed to ABSENT")}
	if state["inhibited"].Reason=="" || state["cancelled"].Reason==""{t.Fatal("ADV-004: inhibition/cancellation evidence missing")}
}

func TestADV005ObservationDoesNotTransition(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}}
	rules:=[]Rule{{ID:"R",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"resolved"}}}}
	run,err:=ExecuteCore(initial,[]Context{{ID:"C"}},rules); if err!=nil{t.Fatal(err)}
	before,err:=TraceHash(run.Trace); if err!=nil{t.Fatal(err)}
	obsA,err:=Observe(run,ObservationContract{ID:"OBS-A",Authority:"authority-a",Entity:"A"}); if err!=nil{t.Fatal(err)}
	obsB,err:=Observe(run,ObservationContract{ID:"OBS-B",Authority:"authority-b",Entity:"B"}); if err!=nil{t.Fatal(err)}
	after,err:=TraceHash(run.Trace); if err!=nil{t.Fatal(err)}
	if before!=after || obsA.TraceHash!=before || obsB.TraceHash!=before{t.Fatal("ADV-005: observation mutated transition evidence")}
	if obsA.Authority==obsB.Authority{t.Fatal("ADV-005: observation authority collapsed")}
	if obsA.Value==obsB.Value{t.Fatal("ADV-005: observers did not expose distinct aspects of same immutable trace")}
}
