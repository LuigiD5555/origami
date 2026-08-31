package fixedcarrier

import (
	"bytes"
	"image/png"
	"testing"
)

func TestExecutionComplianceR1PreservesR5AndExactProgram(t *testing.T) {
	base,meta,err:=RenderTemporalCarrier(temporalCarrierFixture());if err!=nil{t.Fatal(err)}
	grammarMutation:=CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true}
	execMutation:=CandidateMutation{Kind:"PROMPT",Target:"EXECUTION_POLICY",Value:"EXECUTE_VISIBLE_RULES_TO_STABLE_R1",Experimental:true}
	identityMutation:=CandidateMutation{Kind:"REDUNDANCY",Target:"CELL_IDENTITY_ENCODING",Value:VisibleCellIdentityRedundancyR1,Experimental:true}
	fromMutation:=CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"FROM_STATE_PRECONDITION_VISIBILITY",Value:VisibleFromStatePreconditionR1,Experimental:true}
	roleMutation:=CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"RULE_ROLE_BINDING",Value:VisibleRuleRoleBindingR1,Experimental:true}
	grammar,_,err:=BuildTemporalCandidate(base,CandidateSpec{Schema:CandidateSpecSchema,ID:"g",Mutations:[]CandidateMutation{grammarMutation}});if err!=nil{t.Fatal(err)}
	exec,_,err:=BuildTemporalCandidate(grammar,CandidateSpec{Schema:CandidateSpecSchema,ID:"e",Mutations:[]CandidateMutation{execMutation}});if err!=nil{t.Fatal(err)}
	identity,_,err:=BuildTemporalCandidate(exec,CandidateSpec{Schema:CandidateSpecSchema,ID:"i",Mutations:[]CandidateMutation{identityMutation}});if err!=nil{t.Fatal(err)}
	from,_,err:=BuildTemporalCandidate(identity,CandidateSpec{Schema:CandidateSpecSchema,ID:"f",Mutations:[]CandidateMutation{fromMutation}});if err!=nil{t.Fatal(err)}
	r5,_,err:=BuildRuleRoleBindingCandidate(from,"r5");if err!=nil{t.Fatal(err)}
	r6,report,err:=BuildExecutionComplianceCandidate(r5,"r6");if err!=nil{t.Fatal(err)}
	if report.CandidateProgramSHA256!=meta.ProgramSHA256||!report.ExactProgramPreserved{t.Fatal("R6 changed exact program")}
	if bytes.Equal(r5,r6){t.Fatal("R6 must alter execution-policy surface")}

	r5img,err:=png.Decode(bytes.NewReader(r5));if err!=nil{t.Fatal(err)}
	r6img,err:=png.Decode(bytes.NewReader(r6));if err!=nil{t.Fatal(err)}
	// R5 rule-role grammar is frozen.
	for y:=202;y<314;y++{for x:=18;x<622;x++{if r5img.At(x,y)!=r6img.At(x,y){t.Fatalf("R5 rule surface changed at %d,%d",x,y)}}}
	// Exact program visual payload is frozen.
	for y:=415;y<617;y++{for x:=219;x<422;x++{if r5img.At(x,y)!=r6img.At(x,y){t.Fatalf("exact payload pixel changed at %d,%d",x,y)}}}

	decoded,err:=DecodeTemporalCarrierPNG(r6);if err!=nil{t.Fatal(err)}
	manifest:=TemporalInteropBuildManifestWithInherited(report,decoded,[]CandidateMutation{grammarMutation,execMutation,identityMutation,fromMutation,roleMutation})
	if factValue(manifest.VisibleSemantics.Facts,"RULE_ROLE_BINDING")!=VisibleRuleRoleBindingR1{t.Fatal("R6 lost R5 role binding")}
	if factValue(manifest.VisibleSemantics.Facts,"EXECUTION_POLICY_COMPLIANCE")!=ExecuteDontSummarizeToStableR1{t.Fatal("R6 compliance semantic fact missing")}
	if factValue(manifest.VisibleText.Facts,"EXECUTION_POLICY_COMPLIANCE.MODE_TEXT")!=ExecutionComplianceModeTextR1{t.Fatal("R6 mode text missing")}
	if factValue(manifest.VisibleText.Facts,"EXECUTION_POLICY_COMPLIANCE.STOP_TEXT")!=ExecutionComplianceStopTextR1{t.Fatal("R6 stop text missing")}
}
