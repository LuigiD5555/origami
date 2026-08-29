package reference

import (
	"reflect"
	"testing"
)

func TestADV009OpposingContributionsCancelWithoutWinner(t *testing.T){
	cs:=[]Contribution{{ID:"C-plus",Entity:"X",Value:"activate",Weight:1,SourceRule:"R-plus"},{ID:"C-minus",Entity:"X",Value:"inhibit",Weight:-1,SourceRule:"R-minus"}}
	result,evidence,err:=ResolveContributions("X",cs);if err!=nil{t.Fatal(err)}
	if result.Status!=StatusCancelled{t.Fatalf("ADV-009: expected CANCELLED, got %+v",result)}
	if len(evidence.Contributions)!=2||evidence.Net!=0{t.Fatalf("ADV-009: opposing evidence was discarded: %+v",evidence)}
	if evidence.Contributions[0].SourceRule==evidence.Contributions[1].SourceRule{t.Fatal("ADV-009: independent provenance collapsed")}
}

func TestADV009DeclarationOrderCannotChooseWinner(t *testing.T){
	a:=[]Contribution{{ID:"A",Entity:"X",Weight:2,SourceRule:"R-a"},{ID:"B",Entity:"X",Weight:-2,SourceRule:"R-b"}}
	b:=[]Contribution{a[1],a[0]}
	ra,ea,err:=ResolveContributions("X",a);if err!=nil{t.Fatal(err)};rb,eb,err:=ResolveContributions("X",b);if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(ra,rb)||!reflect.DeepEqual(ea,eb){t.Fatalf("ADV-009: declaration order changed interference semantics\nA=%+v %+v\nB=%+v %+v",ra,ea,rb,eb)}
}

func TestADV009RemovingEitherContributionChangesOutcome(t *testing.T){
	plus:=Contribution{ID:"plus",Entity:"X",Weight:1,SourceRule:"R-plus"};minus:=Contribution{ID:"minus",Entity:"X",Weight:-1,SourceRule:"R-minus"}
	both,_,err:=ResolveContributions("X",[]Contribution{plus,minus});if err!=nil{t.Fatal(err)};onlyPlus,_,err:=ResolveContributions("X",[]Contribution{plus});if err!=nil{t.Fatal(err)};onlyMinus,_,err:=ResolveContributions("X",[]Contribution{minus});if err!=nil{t.Fatal(err)}
	if both.Status!=StatusCancelled||onlyPlus.Status!=StatusPresent||onlyMinus.Status!=StatusInhibited{t.Fatalf("ADV-009: counterfactual outcomes wrong: both=%s plus=%s minus=%s",both.Status,onlyPlus.Status,onlyMinus.Status)}
}

func TestADV009CancellationIsNotUnknownOrAbsent(t *testing.T){
	result,_,err:=ResolveContributions("X",[]Contribution{{ID:"p",Entity:"X",Weight:3},{ID:"n",Entity:"X",Weight:-3}});if err!=nil{t.Fatal(err)}
	if result.Status==StatusUnknown||result.Status==StatusAbsent{t.Fatalf("ADV-009: cancellation collapsed to %s",result.Status)}
}
