package reference

import (
	"reflect"
	"testing"
)

func TestADV006FoldPreservesUnresolvedAlternatives(t *testing.T){
	src:=FoldSource{
		State:CoreState{"signal":{Status:StatusPresent,Value:"ready"},"decision":{Status:StatusUnknown}},
		Alternatives:[]AlternativeSet{{ID:"decision-path",Members:[]SemanticValue{{Status:StatusPresent,Value:"left"},{Status:StatusPresent,Value:"right"},{Status:StatusCancelled,Reason:"balanced paths"}}}},
		Rules:[]Rule{{ID:"R-relate",Participants:[]Participant{{Entity:"signal",Status:StatusPresent,Value:"ready"}},Then:[]Effect{{Entity:"decision",Status:StatusUnknown,Reason:"await observation"}}}},
	}
	folded,err:=Fold(src);if err!=nil{t.Fatal(err)}
	if folded.Commitment==""{t.Fatal("ADV-006: Fold emitted no verifiable commitment")}
	if len(folded.Alternatives)!=1||len(folded.Alternatives[0].Members)!=3{t.Fatalf("ADV-006: unresolved alternatives were reduced: %+v",folded.Alternatives)}
	if folded.State["decision"].Status!=StatusUnknown{t.Fatal("ADV-006: Fold resolved UNKNOWN without observation authority")}
	unfolded,err:=Unfold(folded);if err!=nil{t.Fatal(err)}
	if !reflect.DeepEqual(canonicalFoldSource(src),unfolded){t.Fatalf("ADV-006: Unfold did not restore semantic source\nwant=%+v\ngot=%+v",canonicalFoldSource(src),unfolded)}
}

func TestADV006FoldIsDeterministicAcrossDeclarationOrder(t *testing.T){
	a:=FoldSource{State:CoreState{"A":{Status:StatusPresent},"B":{Status:StatusUnknown}},Alternatives:[]AlternativeSet{{ID:"z",Members:[]SemanticValue{{Status:StatusPresent,Value:"1"},{Status:StatusPresent,Value:"2"}}},{ID:"a",Members:[]SemanticValue{{Status:StatusPresent,Value:"x"},{Status:StatusPresent,Value:"y"}}}},Rules:[]Rule{{ID:"R2"},{ID:"R1"}}}
	b:=FoldSource{State:cloneCoreState(a.State),Alternatives:[]AlternativeSet{a.Alternatives[1],a.Alternatives[0]},Rules:[]Rule{a.Rules[1],a.Rules[0]}}
	fa,err:=Fold(a);if err!=nil{t.Fatal(err)};fb,err:=Fold(b);if err!=nil{t.Fatal(err)}
	if fa.Commitment!=fb.Commitment{t.Fatalf("ADV-006: declaration order changed commitment: %s != %s",fa.Commitment,fb.Commitment)}
}

func TestADV006TamperedFoldCannotUnfold(t *testing.T){
	src:=FoldSource{State:CoreState{"A":{Status:StatusUnknown}},Alternatives:[]AlternativeSet{{ID:"choice",Members:[]SemanticValue{{Status:StatusPresent,Value:"x"},{Status:StatusPresent,Value:"y"}}}}}
	folded,err:=Fold(src);if err!=nil{t.Fatal(err)}
	folded.Alternatives[0].Members[0]=SemanticValue{Status:StatusPresent,Value:"tampered"}
	if _,err:=Unfold(folded);err==nil{t.Fatal("ADV-006: tampered folded semantics passed commitment verification")}
}

func TestADV006RejectsFakeSingleAlternativeFold(t *testing.T){
	src:=FoldSource{State:CoreState{"A":{Status:StatusUnknown}},Alternatives:[]AlternativeSet{{ID:"choice",Members:[]SemanticValue{{Status:StatusPresent,Value:"preselected"}}}}}
	if _,err:=Fold(src);err==nil{t.Fatal("ADV-006: Fold accepted a preselected singleton disguised as unresolved alternatives")}
}
