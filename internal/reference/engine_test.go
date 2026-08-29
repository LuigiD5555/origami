package reference

import (
	"reflect"
	"testing"
)

func fixture() Experiment { return Experiment{ID:"test",Entities:[]string{"A","B"},Domain:[]string{"active","inactive"},Relations:[]Relation{{Type:"requires",From:"A",To:"B"}},Protocol:Protocol{MaxSteps:8}} }

func TestDeterministic(t *testing.T){ e:=fixture(); s:=State{"A":"active","B":"active"}; a:=Execute(s,e); b:=Execute(s,e); if !reflect.DeepEqual(a,b){t.Fatal("same input produced different trace")}; if a.Classification!="FIXED_POINT"{t.Fatalf("got %s",a.Classification)} }
func TestContradiction(t *testing.T){ e:=fixture(); r:=Execute(State{"A":"active","B":"inactive"},e); if r.Classification!="CONTRADICTION"{t.Fatalf("got %s",r.Classification)} }
func TestEnumeration(t *testing.T){ e:=fixture(); got:=Enumerate(e); if len(got)!=4{t.Fatalf("got %d states",len(got))} }
func TestEXP001TerminatesInDeclaredClass(t *testing.T){ e,err:=Load("../../experiments/EXP-001-relational-state/experiment.json"); if err!=nil{t.Fatal(err)}; ev:=RunExperiment(e); want:=1; for range e.Entities { want*=len(e.Domain) }; if len(ev.Runs)!=want{t.Fatalf("runs=%d want=%d",len(ev.Runs),want)}; allowed:=map[string]bool{"FIXED_POINT":true,"CYCLE":true,"CONTRADICTION":true,"BUDGET_EXHAUSTED":true}; for _,r:=range ev.Runs{if !allowed[r.Classification]{t.Fatalf("undeclared classification %q",r.Classification)}} }
