package reference

import (
	"reflect"
	"testing"
)

func TestPhaseBContextChangesEvolution(t *testing.T) {
	initial := CoreState{"A":{Status:StatusPresent,Value:"ready"}, "B":{Status:StatusUnknown}}
	rules := []Rule{{ID:"R-CONTEXT",When:[]Predicate{{Source:"context",Key:"mode",Equals:"alpha"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"active"}}}}
	alpha, err := ExecuteCore(initial, []Context{{ID:"C-alpha",Values:map[string]string{"mode":"alpha"}}}, rules); if err!=nil{t.Fatal(err)}
	beta, err := ExecuteCore(initial, []Context{{ID:"C-beta",Values:map[string]string{"mode":"beta"}}}, rules); if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(alpha.Initial,beta.Initial){t.Fatal("initial state changed across contexts")}
	if reflect.DeepEqual(alpha.Trace,beta.Trace){t.Fatal("different contexts did not produce different traces")}
	if got:=alpha.Transitions[0].FiredRules; len(got)!=1 || got[0]!="R-CONTEXT"{t.Fatalf("missing rule provenance: %v",got)}
	if len(beta.Transitions[0].FiredRules)!=0{t.Fatalf("unexpected beta rules: %v",beta.Transitions[0].FiredRules)}
}

func TestPhaseBRulesAreData(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}}
	ctx:=[]Context{{ID:"C"}}
	r1:=[]Rule{{ID:"R1",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"one"}}}}
	r2:=[]Rule{{ID:"R2",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"two"}}}}
	a,err:=ExecuteCore(initial,ctx,r1); if err!=nil{t.Fatal(err)}
	b,err:=ExecuteCore(initial,ctx,r2); if err!=nil{t.Fatal(err)}
	if reflect.DeepEqual(a.Trace,b.Trace){t.Fatal("changing declared rules did not change trace")}
	if a.Transitions[0].FiredRules[0]!="R1" || b.Transitions[0].FiredRules[0]!="R2"{t.Fatal("rule provenance missing")}
}

func TestPhaseBSemanticStatusesRemainDistinct(t *testing.T) {
	values:=[]SemanticValue{{Status:StatusAbsent},{Status:StatusUnknown},{Status:StatusInhibited,Reason:"blocked"},{Status:StatusCancelled,Reason:"opposing contributions"}}
	seen:=map[SemanticStatus]bool{}
	for _,v:=range values{ if seen[v.Status]{t.Fatalf("duplicate status %s",v.Status)}; seen[v.Status]=true }
	if len(seen)!=4{t.Fatalf("statuses collapsed: %v",seen)}
}

func TestPhaseBObservationDoesNotMutateTransition(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"},"B":{Status:StatusUnknown}}
	rules:=[]Rule{{ID:"R",When:[]Predicate{{Source:"state",Key:"A",Equals:"active"}},Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"resolved"}}}}
	run,err:=ExecuteCore(initial,[]Context{{ID:"C"}},rules); if err!=nil{t.Fatal(err)}
	before,err:=TraceHash(run.Trace); if err!=nil{t.Fatal(err)}
	a,err:=Observe(run,ObservationContract{ID:"OBS-A",Authority:"test-a",Entity:"A"}); if err!=nil{t.Fatal(err)}
	b,err:=Observe(run,ObservationContract{ID:"OBS-B",Authority:"test-b",Entity:"B"}); if err!=nil{t.Fatal(err)}
	after,err:=TraceHash(run.Trace); if err!=nil{t.Fatal(err)}
	if before!=after || a.TraceHash!=before || b.TraceHash!=before{t.Fatal("observation mutated or detached from transition trace")}
	if a.Authority==b.Authority{t.Fatal("observation authority not preserved")}
	if a.Value.Value==b.Value.Value{t.Fatal("distinct observations did not expose distinct values")}
}

func TestPhaseBConflictingWritesAreRejected(t *testing.T) {
	initial:=CoreState{"A":{Status:StatusPresent,Value:"active"}}
	rules:=[]Rule{
		{ID:"R1",Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"x"}}},
		{ID:"R2",Then:[]Effect{{Entity:"B",Status:StatusPresent,Value:"y"}}},
	}
	if _,_,err:=ApplyRules(initial,Context{ID:"C"},rules); err==nil{t.Fatal("conflicting writes silently became priority semantics")}
}
