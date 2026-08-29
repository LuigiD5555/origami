package reference

import (
	"reflect"
	"testing"
)

func TestADV007SelectiveUnfoldReturnsOnlyDependencyClosure(t *testing.T){
	src:=FoldSource{State:CoreState{
		"query":{Status:StatusPresent,Value:"q"},
		"concept":{Status:StatusPresent,Value:"semantic-core"},
		"evidence":{Status:StatusPresent,Value:"verified"},
		"unrelated":{Status:StatusPresent,Value:"must-stay-folded"},
		"other":{Status:StatusUnknown},
	}}
	folded,err:=Fold(src);if err!=nil{t.Fatal(err)}
	graph:=AddressGraph{Dependencies:map[string][]string{"query":{"concept"},"concept":{"evidence"},"evidence":{},"unrelated":{"other"}}}
	view,err:=SelectiveUnfold(folded,graph,[]string{"query"});if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(view.Closure,[]string{"concept","evidence","query"}){t.Fatalf("ADV-007: wrong closure %v",view.Closure)}
	if _,ok:=view.State["unrelated"];ok{t.Fatal("ADV-007: unrelated semantic state was unfolded")}
	if _,ok:=view.State["other"];ok{t.Fatal("ADV-007: unrelated dependency region was unfolded")}
	if !reflect.DeepEqual(view.Folded,[]string{"other","unrelated"}){t.Fatalf("ADV-007: still-folded set incorrect: %v",view.Folded)}
	if view.Commitment!=folded.Commitment{t.Fatal("ADV-007: selective view detached from verified Fold commitment")}
}

func TestADV007SelectiveEqualsSelectionAfterFullUnfold(t *testing.T){
	src:=FoldSource{State:CoreState{"A":{Status:StatusPresent,Value:"a"},"B":{Status:StatusPresent,Value:"b"},"C":{Status:StatusPresent,Value:"c"},"Z":{Status:StatusPresent,Value:"z"}}}
	folded,err:=Fold(src);if err!=nil{t.Fatal(err)};graph:=AddressGraph{Dependencies:map[string][]string{"A":{"B"},"B":{"C"}}}
	selective,err:=SelectiveUnfold(folded,graph,[]string{"A"});if err!=nil{t.Fatal(err)}
	full,err:=Unfold(folded);if err!=nil{t.Fatal(err)}
	selectedFromFull:=CoreState{};for _,id:=range selective.Closure{selectedFromFull[id]=full.State[id]}
	if !reflect.DeepEqual(selective.State,selectedFromFull){t.Fatalf("ADV-007: selective semantics differ from full-unfold selection\nselective=%v\nfull-selection=%v",selective.State,selectedFromFull)}
}

func TestADV007CrossDependencyIsIncluded(t *testing.T){
	src:=FoldSource{State:CoreState{"regionA":{Status:StatusPresent},"bridge":{Status:StatusPresent},"regionB":{Status:StatusPresent},"noise":{Status:StatusPresent}}}
	folded,err:=Fold(src);if err!=nil{t.Fatal(err)};graph:=AddressGraph{Dependencies:map[string][]string{"regionA":{"bridge"},"bridge":{"regionB"}}}
	view,err:=SelectiveUnfold(folded,graph,[]string{"regionA"});if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(view.Closure,[]string{"bridge","regionA","regionB"}){t.Fatalf("ADV-007: cross dependency omitted: %v",view.Closure)}
	if _,ok:=view.State["noise"];ok{t.Fatal("ADV-007: independent region leaked into selective unfold")}
}

func TestADV007RejectsDependencyCycles(t *testing.T){
	graph:=AddressGraph{Dependencies:map[string][]string{"A":{"B"},"B":{"A"}}}
	if _,err:=DependencyClosure(graph,[]string{"A"});err==nil{t.Fatal("ADV-007: cyclic dependency graph accepted without explicit semantics")}
}
