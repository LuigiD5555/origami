package fixedcarrier

import (
	"bytes"
	"image/png"
	"testing"
)

func TestExecutionPolicyCandidatePreservesGrammarAndProgram(t *testing.T){
	base,meta,err:=RenderTemporalCarrier(temporalCarrierFixture());if err!=nil{t.Fatal(err)}
	grammar,_,err:=BuildTemporalCandidate(base,CandidateSpec{Schema:CandidateSpecSchema,ID:"grammar",Mutations:[]CandidateMutation{{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true}}});if err!=nil{t.Fatal(err)}
	exec,report,err:=BuildTemporalCandidate(grammar,CandidateSpec{Schema:CandidateSpecSchema,ID:"execute",Mutations:[]CandidateMutation{{Kind:"PROMPT",Target:"EXECUTION_POLICY",Value:"EXECUTE_VISIBLE_RULES_TO_STABLE_R1",Experimental:true}}});if err!=nil{t.Fatal(err)}
	if report.CandidateProgramSHA256!=meta.ProgramSHA256{t.Fatal("execution directive changed exact program")}
	if bytes.Equal(grammar,exec){t.Fatal("execution directive must alter visible carrier")}

	gi,err:=png.Decode(bytes.NewReader(grammar));if err!=nil{t.Fatal(err)}
	ei,err:=png.Decode(bytes.NewReader(exec));if err!=nil{t.Fatal(err)}
	// The R1 microgrammar lives in y=202..313 and must remain pixel-identical.
	for y:=202;y<314;y++{for x:=18;x<622;x++{if gi.At(x,y)!=ei.At(x,y){t.Fatalf("grammar pixel changed at %d,%d",x,y)}}}

	decoded,err:=DecodeTemporalCarrierPNG(exec);if err!=nil{t.Fatal(err)}
	manifest:=TemporalInteropBuildManifest(report,decoded)
	if manifest.Schema!=InteropBuildManifestSchemaR1{t.Fatalf("schema=%q",manifest.Schema)}
	if manifest.VisibleSemantics.ProgramSHA256!=meta.ProgramSHA256{t.Fatal("manifest program sha drift")}
	if factValue(manifest.VisibleSemantics.Facts,"EXECUTION_POLICY")!="EXECUTE_VISIBLE_RULES_TO_STABLE_R1"{t.Fatalf("execution policy missing: %#v",manifest.VisibleSemantics.Facts)}
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r1.REQUIRES")!="A=ACTIVE"{t.Fatalf("R1 semantic drift: %#v",manifest.VisibleSemantics.Facts)}
}

func TestCellIdentityRedundancyR1ChangesOnlyVisibleIdentity(t *testing.T){
	base,meta,err:=RenderTemporalCarrier(temporalCarrierFixture());if err!=nil{t.Fatal(err)}
	grammarMutation:=CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true}
	execMutation:=CandidateMutation{Kind:"PROMPT",Target:"EXECUTION_POLICY",Value:"EXECUTE_VISIBLE_RULES_TO_STABLE_R1",Experimental:true}
	grammar,_,err:=BuildTemporalCandidate(base,CandidateSpec{Schema:CandidateSpecSchema,ID:"grammar",Mutations:[]CandidateMutation{grammarMutation}});if err!=nil{t.Fatal(err)}
	exec,_,err:=BuildTemporalCandidate(grammar,CandidateSpec{Schema:CandidateSpecSchema,ID:"execute-to-stable-text-r1",Mutations:[]CandidateMutation{execMutation}});if err!=nil{t.Fatal(err)}
	identityMutation:=CandidateMutation{Kind:"REDUNDANCY",Target:"CELL_IDENTITY_ENCODING",Value:VisibleCellIdentityRedundancyR1,Experimental:true}
	r3,report,err:=BuildTemporalCandidate(exec,CandidateSpec{Schema:CandidateSpecSchema,ID:"cell-identity-redundancy-r1",Mutations:[]CandidateMutation{identityMutation}});if err!=nil{t.Fatal(err)}
	if report.CandidateProgramSHA256!=meta.ProgramSHA256||!report.ExactProgramPreserved{t.Fatal("R3 changed exact TemporalProgram")}
	if bytes.Equal(exec,r3){t.Fatal("R3 must change the visible carrier")}

	parentImg,err:=png.Decode(bytes.NewReader(exec));if err!=nil{t.Fatal(err)}
	r3Img,err:=png.Decode(bytes.NewReader(r3));if err!=nil{t.Fatal(err)}
	// R2 execution-policy row is frozen and must remain pixel-identical.
	for y:=88;y<102;y++{for x:=20;x<616;x++{if parentImg.At(x,y)!=r3Img.At(x,y){t.Fatalf("execution policy pixel changed at %d,%d",x,y)}}}

	decoded,err:=DecodeTemporalCarrierPNG(r3);if err!=nil{t.Fatal(err)}
	manifest:=TemporalInteropBuildManifestWithInherited(report,decoded,[]CandidateMutation{grammarMutation,execMutation})
	if factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_A")!="A[01]"{t.Fatalf("A visible identity=%q",factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_A"))}
	if factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_B")!="B[02]"{t.Fatalf("B visible identity=%q",factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_B"))}
	if factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_C")!="C[03]"{t.Fatalf("C visible identity=%q",factValue(manifest.VisibleSemantics.Facts,"VISIBLE_CELL_ID_C"))}
	// Canonical semantic rule identities remain raw A/B/C.
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r1.REQUIRES")!="A=ACTIVE"{t.Fatal("R3 rewrote canonical rule precondition")}
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r1.TARGET")!="B"{t.Fatal("R3 rewrote canonical rule target")}
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r4.REQUIRES")!="C=ACTIVE"{t.Fatal("R3 rewrote canonical r4 precondition")}
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r4.TARGET")!="B"{t.Fatal("R3 rewrote canonical r4 target")}
	if factValue(manifest.VisibleText.Facts,"CELL.A.LABEL")!="CELL A[01]"{t.Fatalf("visible A label missing: %#v",manifest.VisibleText.Facts)}
	if factValue(manifest.VisibleText.Facts,"RULE.r1.TEXT")!="IF A[01]=ACTIVE => B[02]:IDLE>ACTIVE"{t.Fatalf("visible r1 text=%q",factValue(manifest.VisibleText.Facts,"RULE.r1.TEXT"))}
	if factValue(manifest.VisibleText.Facts,"RULE.r4.TEXT")!="IF C[03]=ACTIVE => B[02]:ACTIVE>DONE"{t.Fatalf("visible r4 text=%q",factValue(manifest.VisibleText.Facts,"RULE.r4.TEXT"))}
	if factValue(manifest.VisibleText.Facts,"EXECUTION_POLICY.TEXT")!=ExecuteVisibleRulesToStableTextR1{t.Fatal("R3 lost frozen execution directive")}
}

func TestSemanticManifestBaselineHasNoExecutionPolicy(t *testing.T){
	base,_,err:=RenderTemporalCarrier(temporalCarrierFixture());if err!=nil{t.Fatal(err)}
	decoded,err:=DecodeTemporalCarrierPNG(base);if err!=nil{t.Fatal(err)}
	m:=TemporalSemanticManifest(decoded,nil)
	if factValue(m.Facts,"EXECUTION_POLICY")!="NONE"{t.Fatalf("unexpected execution policy: %#v",m.Facts)}
	if factValue(m.Facts,"VISIBLE_CELL_ID_A")!="A"{t.Fatalf("baseline visible identity=%q",factValue(m.Facts,"VISIBLE_CELL_ID_A"))}
}

func factValue(f []SemanticFact,key string)string{for _,x:=range f{if x.Key==key{return x.Value}};return ""}
