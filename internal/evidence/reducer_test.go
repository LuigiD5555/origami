package evidence

import (
	"reflect"
	"testing"

	"github.com/LuigiD5555/origami/internal/memory"
)

func testStore(t *testing.T) *memory.Store {
	t.Helper()
	nodes := []memory.Node{
		{
			Address: "ohf://paper/models/blt-claim",
			Kind: "concept", SemanticClass: "architecture", Label: "BLT dynamic patching",
			Levels: map[memory.Fidelity]memory.ContentLevel{
				memory.FidelitySummary: {Text: "BLT groups bytes into dynamic patches."},
			},
		},
		{
			Address: "ohf://paper/evidence/blt-support",
			Kind: "evidence", SemanticClass: "source-block", Label: "BLT support evidence",
			Levels: map[memory.Fidelity]memory.ContentLevel{
				memory.FidelityEvidence: {Text: "Patch boundaries can follow next-byte entropy.", SourceRef: "paper:blt#p2", Verified: true},
				memory.FidelityExact: {Text: "EXACT-BLT-BLOCK", SourceRef: "paper:blt#p2-exact", Verified: true},
			},
			Provenance: memory.Provenance{SourceID: "paper:blt", SourceRef: "paper:blt", SourceHash: "sha256:source-blt"},
		},
		{
			Address: "ohf://paper/evidence/blt-oppose",
			Kind: "evidence", SemanticClass: "source-block", Label: "Opposing evidence",
			Levels: map[memory.Fidelity]memory.ContentLevel{
				memory.FidelityEvidence: {Text: "A contradictory test fixture.", SourceRef: "paper:other#p1", Verified: true},
			},
			Provenance: memory.Provenance{SourceID: "paper:other", SourceHash: "sha256:source-other"},
		},
		{
			Address: "ohf://paper/evidence/unverified",
			Kind: "evidence", SemanticClass: "source-block", Label: "Unverified note",
			Levels: map[memory.Fidelity]memory.ContentLevel{
				memory.FidelityEvidence: {Text: "Unverified model interpretation.", SourceRef: "note:1", Verified: false},
			},
			Provenance: memory.Provenance{SourceID: "note:1"},
		},
	}
	store, err := memory.BuildStore("paper", nodes)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestReducerVerifiesClaimFromAddressedEvidence(t *testing.T) {
	store := testStore(t)
	ref := EvidenceRef{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityEvidence, SourceHash: "sha256:source-blt"}
	result, err := NewReducer(NewMemoryResolver(store), DefaultPolicy()).Reduce([]Proposal{{
		ID: "agent-a-1", ClaimKey: "blt.dynamic-patching", Proposition: "BLT can use entropy-driven dynamic patches.",
		AgentID: "reader-a", Position: PositionSupport, Evidence: []EvidenceRef{ref}, Confidence: .91,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Claims) != 1 || result.Claims[0].Status != StatusVerified {
		t.Fatalf("expected VERIFIED, got %+v", result)
	}
	if result.Claims[0].Exactness != ExactnessNotClaimed {
		t.Fatalf("evidence fidelity must not claim exactness: %+v", result.Claims[0])
	}
}

func TestReducerIsOrderDeterministic(t *testing.T) {
	store := testStore(t)
	ref := EvidenceRef{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityEvidence}
	proposalsA := []Proposal{
		{ID: "p2", ClaimKey: "blt.dynamic-patching", Proposition: "Dynamic patches", Position: PositionSupport, Evidence: []EvidenceRef{ref}},
		{ID: "p1", ClaimKey: "blt.dynamic-patching", Proposition: "BLT dynamic patching", Position: PositionSupport, Evidence: []EvidenceRef{ref}},
	}
	proposalsB := []Proposal{proposalsA[1], proposalsA[0]}
	r := NewReducer(NewMemoryResolver(store), DefaultPolicy())
	a, err := r.Reduce(proposalsA)
	if err != nil { t.Fatal(err) }
	b, err := r.Reduce(proposalsB)
	if err != nil { t.Fatal(err) }
	if a.Digest != b.Digest || !reflect.DeepEqual(a.Claims, b.Claims) {
		t.Fatalf("reduction must be deterministic across proposal order\nA=%+v\nB=%+v", a, b)
	}
	if a.Metrics.EvidenceDeduplicated != 1 {
		t.Fatalf("duplicate swarm evidence must count once: %+v", a.Metrics)
	}
}

func TestReducerPreservesVerifiedConflict(t *testing.T) {
	store := testStore(t)
	result, err := NewReducer(NewMemoryResolver(store), DefaultPolicy()).Reduce([]Proposal{
		{ID: "support", ClaimKey: "claim-x", Position: PositionSupport, Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityEvidence}}},
		{ID: "oppose", ClaimKey: "claim-x", Position: PositionOppose, Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-oppose", Fidelity: memory.FidelityEvidence}}},
	})
	if err != nil { t.Fatal(err) }
	if got := result.Claims[0].Status; got != StatusConflict {
		t.Fatalf("expected CONFLICT, got %s", got)
	}
}

func TestReducerReturnsUnknownForUnverifiedEvidence(t *testing.T) {
	store := testStore(t)
	result, err := NewReducer(NewMemoryResolver(store), DefaultPolicy()).Reduce([]Proposal{{
		ID: "guess", ClaimKey: "claim-y", Position: PositionSupport,
		Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/unverified", Fidelity: memory.FidelityEvidence}},
	}})
	if err != nil { t.Fatal(err) }
	claim := result.Claims[0]
	if claim.Status != StatusUnknown {
		t.Fatalf("unverified evidence must resolve UNKNOWN, got %+v", claim)
	}
	if len(claim.RejectedEvidence) != 1 || claim.RejectedEvidence[0].Reason == "" {
		t.Fatalf("expected explicit rejection reason: %+v", claim)
	}
}

func TestReducerRejectsMismatchedCIDAndSourceHash(t *testing.T) {
	store := testStore(t)
	result, err := NewReducer(NewMemoryResolver(store), DefaultPolicy()).Reduce([]Proposal{{
		ID: "bad-anchor", ClaimKey: "claim-z", Position: PositionSupport,
		Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityEvidence, CID: "sha256:not-the-node", SourceHash: "sha256:not-the-source"}},
	}})
	if err != nil { t.Fatal(err) }
	claim := result.Claims[0]
	if claim.Status != StatusUnknown || len(claim.RejectedEvidence) != 1 {
		t.Fatalf("mismatched identity must not establish a claim: %+v", claim)
	}
}

func TestReducerClaimsExactnessOnlyFromVerifiedExactLevel(t *testing.T) {
	store := testStore(t)
	result, err := NewReducer(NewMemoryResolver(store), DefaultPolicy()).Reduce([]Proposal{{
		ID: "exact", ClaimKey: "claim-exact", Position: PositionSupport,
		Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityExact}},
	}})
	if err != nil { t.Fatal(err) }
	claim := result.Claims[0]
	if claim.Status != StatusVerified || claim.Exactness != ExactnessVerified {
		t.Fatalf("verified exact evidence should earn exactness: %+v", claim)
	}
}

func TestReducerRequiresIndependentSourcesWhenPolicyDemandsIt(t *testing.T) {
	store := testStore(t)
	policy := Policy{MinVerifiedEvidence: 1, MinIndependentSources: 2}
	result, err := NewReducer(NewMemoryResolver(store), policy).Reduce([]Proposal{
		{ID: "a", ClaimKey: "claim-sources", Position: PositionSupport, Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityEvidence}}},
		{ID: "b", ClaimKey: "claim-sources", Position: PositionSupport, Evidence: []EvidenceRef{{Address: "ohf://paper/evidence/blt-support", Fidelity: memory.FidelityExact}}},
	})
	if err != nil { t.Fatal(err) }
	if got := result.Claims[0].Status; got != StatusUnknown {
		t.Fatalf("two fidelities of one source are not two independent sources: %s", got)
	}
}
