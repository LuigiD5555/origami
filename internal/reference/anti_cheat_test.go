package reference

import "testing"

func strictRoutePolicy()RouteIntegrityPolicy{return RouteIntegrityPolicy{MaxLiteralFraction:0.10,MaxResidualFraction:0.20,RequireEvidence:true,ForbidExternalCore:true}}

func TestADV010RejectsStoredFinalAnswer(t *testing.T){
	claim:=SemanticClaim{ID:"answer",Entity:"X",Value:SemanticValue{Status:StatusPresent,Value:"expected-answer"},Derivation:DerivedLiteral,EvidenceIDs:[]string{"literal:answer"},LiteralBytes:20,SourceBytes:100}
	if CheckRouteIntegrity(claim,strictRoutePolicy()).Passed{t.Fatal("ADV-010: stored final answer passed as semantic computation")}
}

func TestADV010RejectsSourceSizedLiteralOrResidual(t *testing.T){
	literal:=SemanticClaim{ID:"literal-copy",Entity:"X",Derivation:DerivedTransition,EvidenceIDs:[]string{"R1"},LiteralBytes:950,SourceBytes:1000}
	residual:=SemanticClaim{ID:"residual-copy",Entity:"X",Derivation:DerivedTransition,EvidenceIDs:[]string{"R1"},ResidualBytes:900,SourceBytes:1000}
	if CheckRouteIntegrity(literal,strictRoutePolicy()).Passed{t.Fatal("ADV-010: near-source-sized literal bypass passed")}
	if CheckRouteIntegrity(residual,strictRoutePolicy()).Passed{t.Fatal("ADV-010: near-source-sized residual bypass passed")}
}

func TestADV010RejectsExternalModelAsCoreSolver(t *testing.T){
	claim:=SemanticClaim{ID:"model-answer",Entity:"X",Derivation:DerivedExternal,EvidenceIDs:[]string{"vlm-output"},SourceBytes:1000}
	if CheckRouteIntegrity(claim,strictRoutePolicy()).Passed{t.Fatal("ADV-010: external model solved core semantics but gate passed")}
}

func TestADV010RejectsAnswerWithoutEvidence(t *testing.T){
	claim:=SemanticClaim{ID:"unsupported-answer",Entity:"X",Derivation:DerivedTransition,SourceBytes:1000}
	if CheckRouteIntegrity(claim,strictRoutePolicy()).Passed{t.Fatal("ADV-010: evidence-free semantic answer passed")}
}

func TestADV010AcceptsMachineDerivedBoundedClaim(t *testing.T){
	claim:=SemanticClaim{ID:"derived",Entity:"X",Value:SemanticValue{Status:StatusCancelled},Derivation:DerivedInteraction,EvidenceIDs:[]string{"R-plus","R-minus","contribution:C-plus","contribution:C-minus"},LiteralBytes:20,ResidualBytes:50,SourceBytes:1000}
	r:=CheckRouteIntegrity(claim,strictRoutePolicy());if !r.Passed{t.Fatalf("ADV-010: legitimate bounded machine-derived claim rejected: %v",r.Violations)}
}

func TestADV010ExactRecoveryBudgetDoesNotProveSemanticDerivation(t *testing.T){
	claim:=SemanticClaim{ID:"exact-but-literal",Entity:"document.answer",Value:SemanticValue{Status:StatusPresent,Value:"42"},Derivation:DerivedLiteral,EvidenceIDs:[]string{"sha256:exact"},LiteralBytes:1,SourceBytes:1000}
	if CheckRouteIntegrity(claim,strictRoutePolicy()).Passed{t.Fatal("ADV-010: exact recovery evidence was incorrectly accepted as semantic derivation")}
}
