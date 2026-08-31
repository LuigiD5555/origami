package fixedcarrier

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestTemporalCandidateBuildPreservesExactProgram(t *testing.T) {
	parent, parentMeta, err := RenderTemporalCarrier(temporalCarrierFixture())
	if err != nil { t.Fatal(err) }
	sum := sha256.Sum256(parent)
	spec := CandidateSpec{
		Schema: CandidateSpecSchema,
		ID: "t2-route-candidate",
		ParentSHA256: fmt.Sprintf("%x", sum[:]),
		Mutations: []CandidateMutation{
			{Kind:"LAYOUT",Target:"T1_TO_T2_ENTRY_ROUTE",Value:"EXPLICIT_DIRECTIONAL_ANCHOR",Experimental:true},
			{Kind:"PROMPT",Target:"ROSETTA.S2.READ_SUPERINDEX",Value:"DECLARE_T2_LOCATION_BEFORE_DECODE",Experimental:true},
			{Kind:"REDUNDANCY",Target:"T2_ENTRY_MARKER",Value:"REPEAT_AT_BOOT_AND_ROSETTA",Experimental:true},
		},
	}
	candidate, report, err := BuildTemporalCandidate(parent, spec)
	if err != nil { t.Fatal(err) }
	if len(candidate) != FixedPNGBytes { t.Fatalf("expected %d bytes, got %d", FixedPNGBytes, len(candidate)) }
	if bytes.Equal(parent, candidate) { t.Fatal("candidate should differ visually from parent") }
	if !report.ExactProgramPreserved { t.Fatal("exact program must be preserved") }
	if report.ParentProgramSHA256 != parentMeta.ProgramSHA256 || report.CandidateProgramSHA256 != parentMeta.ProgramSHA256 {
		t.Fatalf("program SHA drift: %#v", report)
	}
	decoded, err := DecodeTemporalCarrierPNG(candidate)
	if err != nil { t.Fatal(err) }
	if decoded.ProgramSHA256 != parentMeta.ProgramSHA256 { t.Fatal("candidate exact program drift") }
}

func TestTemporalCandidateBuildSupportsTemporalAndNumericSignals(t *testing.T) {
	parent, parentMeta, err := RenderTemporalCarrier(temporalCarrierFixture())
	if err != nil { t.Fatal(err) }
	spec := CandidateSpec{Schema:CandidateSpecSchema,ID:"timeline-candidate",Mutations:[]CandidateMutation{
		{Kind:"TEMPORAL_STRUCTURE",Target:"TEMPORAL_ROUTING",Value:"EXPLICIT_PHASE_EVENT_CHECKPOINT_STRUCTURE",Experimental:true},
		{Kind:"NUMERIC_STRUCTURE",Target:"TEMPORAL_ROUTING",Value:"EXPERIMENTAL_ADDRESS_ORDER_SIGNAL",Experimental:true},
	}}
	candidate, report, err := BuildTemporalCandidate(parent, spec)
	if err != nil { t.Fatal(err) }
	if len(candidate) != FixedPNGBytes { t.Fatalf("bad size %d",len(candidate)) }
	if report.CandidateProgramSHA256 != parentMeta.ProgramSHA256 { t.Fatal("temporal mutation touched exact program") }
}

func TestVisibleRuleMicrogrammarIsDistinctAndPreservesExactProgram(t *testing.T) {
	parent, parentMeta, err := RenderTemporalCarrier(temporalCarrierFixture())
	if err != nil { t.Fatal(err) }

	microSpec := CandidateSpec{Schema:CandidateSpecSchema,ID:"t2-temporal-grammar-visible-r1",Mutations:[]CandidateMutation{
		{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true},
	}}
	micro, report, err := BuildTemporalCandidate(parent, microSpec)
	if err != nil { t.Fatal(err) }
	if len(micro) != FixedPNGBytes { t.Fatalf("expected %d-byte candidate, got %d",FixedPNGBytes,len(micro)) }
	if bytes.Equal(parent,micro) { t.Fatal("visible rule microgrammar must change visual pixels") }
	if !report.ExactProgramPreserved { t.Fatal("visible rule microgrammar must preserve exact program") }
	if report.ParentProgramSHA256 != parentMeta.ProgramSHA256 || report.CandidateProgramSHA256 != parentMeta.ProgramSHA256 {
		t.Fatalf("visible rule microgrammar changed TemporalProgram SHA: %#v",report)
	}
	decoded, err := DecodeTemporalCarrierPNG(micro)
	if err != nil { t.Fatal(err) }
	if decoded.ProgramSHA256 != parentMeta.ProgramSHA256 { t.Fatal("decoded candidate exact program drift") }

	checkpointSpec := CandidateSpec{Schema:CandidateSpecSchema,ID:"checkpoint-only",Mutations:[]CandidateMutation{
		{Kind:"TEMPORAL_STRUCTURE",Target:"TEMPORAL_ROUTING",Value:"EXPLICIT_PHASE_EVENT_CHECKPOINT_STRUCTURE",Experimental:true},
	}}
	checkpointOnly, _, err := BuildTemporalCandidate(parent,checkpointSpec)
	if err != nil { t.Fatal(err) }
	if bytes.Equal(micro,checkpointOnly) { t.Fatal("visible rule microgrammar must not collapse to the legacy checkpoint-only temporal mutation") }
}

func TestTemporalCandidateRejectsUnsupportedMutation(t *testing.T) {
	parent, _, err := RenderTemporalCarrier(temporalCarrierFixture())
	if err != nil { t.Fatal(err) }
	spec := CandidateSpec{Schema:CandidateSpecSchema,ID:"bad",Mutations:[]CandidateMutation{{Kind:"DEPTH_STRUCTURE",Target:"T2",Value:"EXPERIMENTAL_DEPTH_SIGNAL",Experimental:true}}}
	if _, _, err := BuildTemporalCandidate(parent, spec); err == nil { t.Fatal("expected unsupported mutation failure") }
}

func TestTemporalCandidateRejectsParentMismatch(t *testing.T) {
	parent, _, err := RenderTemporalCarrier(temporalCarrierFixture())
	if err != nil { t.Fatal(err) }
	spec := CandidateSpec{Schema:CandidateSpecSchema,ID:"bad-parent",ParentSHA256:"00",Mutations:[]CandidateMutation{{Kind:"LAYOUT",Target:"T1_TO_T2_ENTRY_ROUTE",Value:"EXPLICIT_DIRECTIONAL_ANCHOR",Experimental:true}}}
	if _, _, err := BuildTemporalCandidate(parent, spec); err == nil { t.Fatal("expected parent SHA mismatch") }
}
