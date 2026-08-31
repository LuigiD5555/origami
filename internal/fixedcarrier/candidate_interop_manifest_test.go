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
	if factValue(manifest.VisibleSemantics.Facts,"RULE.r1-b-on-a.REQUIRES")!="A=ACTIVE"{t.Fatalf("R1 semantic drift: %#v",manifest.VisibleSemantics.Facts)}
}

func TestSemanticManifestBaselineHasNoExecutionPolicy(t *testing.T){
	base,_,err:=RenderTemporalCarrier(temporalCarrierFixture());if err!=nil{t.Fatal(err)}
	decoded,err:=DecodeTemporalCarrierPNG(base);if err!=nil{t.Fatal(err)}
	m:=TemporalSemanticManifest(decoded,nil)
	if factValue(m.Facts,"EXECUTION_POLICY")!="NONE"{t.Fatalf("unexpected execution policy: %#v",m.Facts)}
}

func factValue(f []SemanticFact,key string)string{for _,x:=range f{if x.Key==key{return x.Value}};return ""}
