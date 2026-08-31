package fixedcarrier

import (
	"bytes"
	"image/png"
	"testing"
)

func TestSynchronousExecutionFidelityPreservesProgramAndFrozenSurfaces(t *testing.T) {
	base, meta, err := RenderTemporalCarrier(temporalCarrierFixture()); if err != nil { t.Fatal(err) }
	r6, _, err := BuildExecutionComplianceCandidate(base, "r6"); if err != nil { t.Fatal(err) }
	r7, report, err := BuildSynchronousExecutionFidelityCandidate(r6, "r7"); if err != nil { t.Fatal(err) }
	if !report.ExactProgramPreserved || report.CandidateProgramSHA256 != meta.ProgramSHA256 { t.Fatal("R7 changed exact TemporalProgram") }
	if bytes.Equal(r6, r7) { t.Fatal("R7 must alter visible synchronous execution guidance") }

	r6img, err := png.Decode(bytes.NewReader(r6)); if err != nil { t.Fatal(err) }
	r7img, err := png.Decode(bytes.NewReader(r7)); if err != nil { t.Fatal(err) }
	for y:=88; y<102; y++ { for x:=20; x<616; x++ { if r6img.At(x,y)!=r7img.At(x,y) { t.Fatalf("R6 mode row changed at %d,%d",x,y) } } }
	for y:=415; y<620; y++ { for x:=219; x<430; x++ { if r6img.At(x,y)!=r7img.At(x,y) { t.Fatalf("payload surface changed at %d,%d",x,y) } } }

	decoded, err := DecodeTemporalCarrierPNG(r7); if err != nil { t.Fatal(err) }
	manifest := TemporalInteropBuildManifestWithInherited(report, decoded, []CandidateMutation{{Kind:"PROMPT",Target:"EXECUTION_POLICY_COMPLIANCE",Value:ExecuteDontSummarizeToStableR1,Experimental:true}})
	if factValue(manifest.VisibleSemantics.Facts,"SYNCHRONOUS_EXECUTION_FIDELITY") != FreezeSelectApplyTogetherR1 { t.Fatal("R7 semantic fact missing") }
	if factValue(manifest.VisibleText.Facts,"SYNCHRONOUS_EXECUTION_FIDELITY.NO_ORDER_TEXT") != SyncFidelityNoOrderTextR1 { t.Fatal("NO ORDER / NO CASCADE text missing") }
	if factValue(manifest.VisibleText.Facts,"EXECUTION_POLICY_COMPLIANCE.MODE_TEXT") != ExecutionComplianceModeTextR1 { t.Fatal("R6 execution mode lost") }
}
