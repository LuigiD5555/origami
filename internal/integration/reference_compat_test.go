package integration

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/core"
	"github.com/LuigiD5555/origami/internal/reference"
)

func TestCoreR1MatchesReferenceR0ForAtomicRequiresRelation(t *testing.T){
	experiment:=reference.Experiment{ID:"compat-requires",Entities:[]string{"A","B"},Domain:[]string{"active","inactive"},Relations:[]reference.Relation{{Type:"requires",From:"A",To:"B"}},Protocol:reference.Protocol{MaxSteps:4}}
	machine:=core.NewMachine([]core.Rule{
		{ID:"requires-contradiction",When:[]core.Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"},{Scope:"state",Key:"B",Op:"data_eq",Data:"inactive"}},Contradiction:&core.Contradiction{Reason:"A requires B but B is inactive"}},
		{ID:"requires-transition",When:[]core.Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"},{Scope:"state",Key:"B",Op:"data_ne",Data:"inactive"}},Then:[]core.Effect{{Key:"B",Mode:"set",Status:core.StatusPresent,Data:"active"}}},
	},4)
	for _,initial:=range reference.Enumerate(experiment){
		r0:=reference.Execute(initial,experiment)
		r1initial:=core.State{Values:map[string]core.Value{}}
		for key,value:=range initial{r1initial.Values[key]=core.Value{Status:core.StatusPresent,Data:value}}
		r1,err:=core.Execute(r1initial,machine,nil);if err!=nil{t.Fatal(err)}
		if r1.Terminated!=r0.Classification{t.Fatalf("classification drift initial=%+v R0=%s R1=%s",initial,r0.Classification,r1.Terminated)}
		r0final:=r0.Trace[len(r0.Trace)-1];r1final:=r1.States[len(r1.States)-1]
		for key,value:=range r0final{if r1final.Values[key].Data!=value{t.Fatalf("state drift initial=%+v key=%s R0=%s R1=%s",initial,key,value,r1final.Values[key].Data)}}
	}
}

func TestCoreR1MatchesReferenceR0ForAtomicExcludesRelation(t *testing.T){
	experiment:=reference.Experiment{ID:"compat-excludes",Entities:[]string{"A","B"},Domain:[]string{"active","inactive"},Relations:[]reference.Relation{{Type:"excludes",From:"A",To:"B"}},Protocol:reference.Protocol{MaxSteps:4}}
	machine:=core.NewMachine([]core.Rule{
		{ID:"excludes-contradiction",When:[]core.Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"},{Scope:"state",Key:"B",Op:"data_eq",Data:"active"}},Contradiction:&core.Contradiction{Reason:"A excludes active B"}},
		{ID:"excludes-transition",When:[]core.Predicate{{Scope:"state",Key:"A",Op:"data_eq",Data:"active"},{Scope:"state",Key:"B",Op:"data_ne",Data:"active"}},Then:[]core.Effect{{Key:"B",Mode:"set",Status:core.StatusPresent,Data:"inactive"}}},
	},4)
	for _,initial:=range reference.Enumerate(experiment){
		r0:=reference.Execute(initial,experiment);r1initial:=core.State{Values:map[string]core.Value{}};for key,value:=range initial{r1initial.Values[key]=core.Value{Status:core.StatusPresent,Data:value}};r1,err:=core.Execute(r1initial,machine,nil);if err!=nil{t.Fatal(err)}
		if r1.Terminated!=r0.Classification{t.Fatalf("classification drift initial=%+v R0=%s R1=%s",initial,r0.Classification,r1.Terminated)}
		r0final:=r0.Trace[len(r0.Trace)-1];r1final:=r1.States[len(r1.States)-1];for key,value:=range r0final{if r1final.Values[key].Data!=value{t.Fatalf("state drift initial=%+v key=%s R0=%s R1=%s",initial,key,value,r1final.Values[key].Data)}}
	}
}
