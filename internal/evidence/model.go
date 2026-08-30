package evidence

import "github.com/LuigiD5555/origami/internal/memory"

const SchemaR0 = "origami.evidence-reduction.r0"

type Position string

const (
	PositionSupport Position = "SUPPORT"
	PositionOppose  Position = "OPPOSE"
	PositionUnknown Position = "UNKNOWN"
)

type ClaimStatus string

const (
	StatusVerified ClaimStatus = "VERIFIED"
	StatusRejected ClaimStatus = "REJECTED"
	StatusConflict ClaimStatus = "CONFLICT"
	StatusUnknown  ClaimStatus = "UNKNOWN"
)

type Exactness string

const (
	ExactnessVerified   Exactness = "VERIFIED_EXACT"
	ExactnessNotClaimed Exactness = "NOT_CLAIMED"
)

type EvidenceRef struct {
	Address    string          `json:"address"`
	Fidelity   memory.Fidelity `json:"fidelity,omitempty"`
	CID        string          `json:"cid,omitempty"`
	SourceHash string          `json:"source_hash,omitempty"`
}

type Proposal struct {
	ID          string        `json:"id"`
	ClaimKey    string        `json:"claim_key"`
	Proposition string        `json:"proposition,omitempty"`
	AgentID     string        `json:"agent_id,omitempty"`
	RunID       string        `json:"run_id,omitempty"`
	Position    Position      `json:"position"`
	Evidence    []EvidenceRef `json:"evidence,omitempty"`
	Confidence  float64       `json:"confidence,omitempty"`
}

type Policy struct {
	MinVerifiedEvidence   int `json:"min_verified_evidence"`
	MinIndependentSources int `json:"min_independent_sources"`
}

func DefaultPolicy() Policy {
	return Policy{MinVerifiedEvidence: 1, MinIndependentSources: 1}
}

type ResolvedEvidence struct {
	Address    string          `json:"address"`
	Fidelity   memory.Fidelity `json:"fidelity"`
	CID        string          `json:"cid"`
	SourceID   string          `json:"source_id,omitempty"`
	SourceRef  string          `json:"source_ref,omitempty"`
	SourceHash string          `json:"source_hash,omitempty"`
	Verified   bool            `json:"verified"`
	Exact      bool            `json:"exact"`
	// Content is intentionally private to the reducer. It is required to decide
	// whether a claim is byte-for-byte exact, but it is never emitted as hidden
	// output in an audit record merely because the source was opened.
	Content string `json:"-"`
}

type EvidenceAudit struct {
	ProposalID string           `json:"proposal_id"`
	Position   Position         `json:"position"`
	Requested  EvidenceRef      `json:"requested"`
	Resolved   ResolvedEvidence `json:"resolved,omitempty"`
	Accepted   bool             `json:"accepted"`
	Reason     string           `json:"reason,omitempty"`
	Anchor     string           `json:"anchor,omitempty"`
	SourceKey  string           `json:"source_key,omitempty"`
	ExactMatch bool             `json:"exact_match,omitempty"`
}

type ClaimResult struct {
	ClaimKey               string          `json:"claim_key"`
	Proposition            string          `json:"proposition,omitempty"`
	PropositionVariants    []string        `json:"proposition_variants,omitempty"`
	Status                 ClaimStatus     `json:"status"`
	Exactness              Exactness       `json:"exactness"`
	CanonicalID            string          `json:"canonical_id"`
	ProposalIDs            []string        `json:"proposal_ids"`
	SupportEvidence        []EvidenceAudit `json:"support_evidence,omitempty"`
	OpposeEvidence         []EvidenceAudit `json:"oppose_evidence,omitempty"`
	RejectedEvidence       []EvidenceAudit `json:"rejected_evidence,omitempty"`
	VerifiedSupportAnchors int             `json:"verified_support_anchors"`
	VerifiedOpposeAnchors  int             `json:"verified_oppose_anchors"`
	IndependentSupport     int             `json:"independent_support_sources"`
	IndependentOppose      int             `json:"independent_oppose_sources"`
	Reason                 string          `json:"reason,omitempty"`
}

type Metrics struct {
	Proposals            int `json:"proposals"`
	Claims               int `json:"claims"`
	EvidenceRequested    int `json:"evidence_requested"`
	EvidenceAccepted     int `json:"evidence_accepted"`
	EvidenceRejected     int `json:"evidence_rejected"`
	EvidenceDeduplicated int `json:"evidence_deduplicated"`
}

type Reduction struct {
	Schema     string        `json:"schema"`
	Policy     Policy        `json:"policy"`
	Claims     []ClaimResult `json:"claims"`
	Metrics    Metrics       `json:"metrics"`
	Digest     string        `json:"digest"`
	Invariants []string      `json:"invariants"`
}

type Resolver interface {
	Resolve(EvidenceRef) (ResolvedEvidence, string, bool)
}
