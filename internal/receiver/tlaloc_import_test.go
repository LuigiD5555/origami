package receiver

import "testing"

const testTraceSHA256 = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestImportTlalocCandidatePreservesAction(t *testing.T) {
	got, err := ImportTlalocCandidate(TlalocHybridArtifactSet{
		Schema:             TlalocHybridArtifactSchema,
		CandidateID:        "candidate-1",
		UniversalPrompt:    "bootstrap from carrier",
		BootStrategy:       []string{"find BOOT"},
		RosettaConstraints: []string{"symbols are carrier-local"},
		MicroProgram: []TlalocMicroRule{{
			ID: "m0000", State: "S0", Token: "BOOT_OK", Action: "advance", NextState: "S1", Emit: "ROSETTA",
		}},
		SourceTraceSHA256: testTraceSHA256,
		WorkingWindow:     4000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Program) != 1 || got.Program[0].Action != "advance" {
		t.Fatalf("distilled action was not preserved: %#v", got.Program)
	}
	if got.Program[0].Emit != "ROSETTA" {
		t.Fatalf("distilled emit was not preserved: %#v", got.Program[0])
	}
}

func TestImportTlalocCandidateDoesNotBindPhysicalSymbols(t *testing.T) {
	got, err := ImportTlalocCandidate(TlalocHybridArtifactSet{
		Schema:             TlalocHybridArtifactSchema,
		CandidateID:        "candidate-2",
		UniversalPrompt:    "bootstrap",
		BootStrategy:       []string{"find BOOT"},
		RosettaConstraints: []string{"do not hardcode glyph meaning"},
		MicroProgram:       []TlalocMicroRule{{ID: "m0", State: "S0", Token: "T", NextState: "S1"}},
		SourceTraceSHA256:  testTraceSHA256,
		WorkingWindow:      4000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Program) != 1 {
		t.Fatal("expected imported semantic program")
	}
	// ImportedCandidate intentionally has no Rosetta physical bindings. Origami
	// must create those when constructing a concrete carrier.
}

func TestImportTlalocCandidateRejectsShortTraceHash(t *testing.T) {
	_, err := ImportTlalocCandidate(TlalocHybridArtifactSet{
		Schema:             TlalocHybridArtifactSchema,
		CandidateID:        "candidate-3",
		UniversalPrompt:    "bootstrap",
		BootStrategy:       []string{"find BOOT"},
		RosettaConstraints: []string{"carrier-local"},
		MicroProgram:       []TlalocMicroRule{{ID: "m0", State: "S0", Token: "T", NextState: "S1"}},
		SourceTraceSHA256:  "short",
		WorkingWindow:      4000,
	})
	if err == nil {
		t.Fatal("expected short provenance hash to be rejected")
	}
}
