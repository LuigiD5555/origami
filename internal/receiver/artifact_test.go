package receiver

import "testing"

func validManifest() ArtifactManifest {
	return ArtifactManifest{
		Schema:             ArtifactSchemaR0,
		Status:             ArtifactPromoted,
		ReceiverContract:   ContractID,
		DistillerContract:  "tlaloc.origami-receiver-distillation.r0",
		CandidateID:        "candidate-r0",
		PromptPath:         "generated/MASTER_PROMPT.md",
		PromptSHA256:       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		ValidationSuite:    []string{"symbol-permutation", "unknown-fail-closed"},
		FalseExact:         0,
		Contaminated:       false,
		PeakActiveTokenEq:  900,
		WorkingWindowToken: 4000,
	}
}

func TestPromotedArtifactRequiresCleanEvidence(t *testing.T) {
	m := validManifest()
	if err := ValidateArtifactManifest(m); err != nil { t.Fatal(err) }

	m.FalseExact = 1
	if err := ValidateArtifactManifest(m); err == nil {
		t.Fatal("expected false exactness to reject promoted artifact")
	}

	m = validManifest()
	m.Contaminated = true
	if err := ValidateArtifactManifest(m); err == nil {
		t.Fatal("expected contamination to reject promoted artifact")
	}

	m = validManifest()
	m.PeakActiveTokenEq = 4001
	if err := ValidateArtifactManifest(m); err == nil {
		t.Fatal("expected active-window violation to reject promoted artifact")
	}
}

func TestReferenceCandidateMayExistBeforePromotionEvidence(t *testing.T) {
	m := validManifest()
	m.Status = ArtifactReferenceCandidate
	m.ValidationSuite = nil
	m.PeakActiveTokenEq = 0
	if err := ValidateArtifactManifest(m); err != nil { t.Fatal(err) }
}
