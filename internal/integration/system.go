package integration

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/LuigiD5555/origami/internal/evidence"
	"github.com/LuigiD5555/origami/internal/memory"
)

const SchemaR0 = "origami.semantic-spine.r0"

type System struct {
	CarrierID string
	nodes     []memory.Node
	store     *memory.Store
	manager   *memory.Manager
}

func NewSystem(carrierID string, nodes []memory.Node) (*System, error) {
	if strings.TrimSpace(carrierID) == "" { carrierID = "origami" }
	s := &System{CarrierID: carrierID}
	if err := s.replaceNodes(nodes); err != nil { return nil, err }
	return s, nil
}

func (s *System) Query(query string, budget memory.ContextBudget) memory.ContextPacket {
	if s == nil || s.manager == nil {
		return memory.ContextPacket{Schema: memory.SchemaR0 + ".context-packet", Query: query, Unknown: true, Reason: "semantic spine unavailable"}
	}
	return s.manager.Query(query, budget)
}

func (s *System) Expand(address string, fidelity memory.Fidelity, budget memory.ContextBudget) memory.ContextPacket {
	if s == nil || s.manager == nil {
		return memory.ContextPacket{Schema: memory.SchemaR0 + ".context-packet", Query: "EXPAND " + address, Unknown: true, Reason: "semantic spine unavailable"}
	}
	return s.manager.Expand(address, fidelity, budget)
}

func (s *System) Reduce(proposals []evidence.Proposal, policy evidence.Policy) (evidence.Reduction, error) {
	if s == nil || s.store == nil { return evidence.Reduction{}, fmt.Errorf("semantic spine unavailable") }
	return evidence.NewReducer(evidence.NewMemoryResolver(s.store), policy).Reduce(proposals)
}

func (s *System) Promote(reduction evidence.Reduction) ([]memory.Node, error) {
	if s == nil { return nil, fmt.Errorf("semantic spine unavailable") }
	promoted := ReductionNodes(s.CarrierID, reduction)
	combined := append([]memory.Node(nil), s.nodes...)
	byAddress := map[string]int{}
	for i, node := range combined { byAddress[node.Address] = i }
	for _, node := range promoted {
		if index, ok := byAddress[node.Address]; ok {
			combined[index] = node
		} else {
			byAddress[node.Address] = len(combined)
			combined = append(combined, node)
		}
	}
	if err := s.replaceNodes(combined); err != nil { return nil, err }
	return promoted, nil
}

func (s *System) Signature() memory.GraphSignature {
	if s == nil || s.store == nil { return memory.GraphSignature{} }
	return s.store.Signature()
}

func (s *System) Nodes() []memory.Node {
	return append([]memory.Node(nil), s.nodes...)
}

func (s *System) replaceNodes(nodes []memory.Node) error {
	copyNodes := append([]memory.Node(nil), nodes...)
	sort.Slice(copyNodes, func(i, j int) bool { return copyNodes[i].Address < copyNodes[j].Address })
	store, err := memory.BuildStore(s.CarrierID, copyNodes)
	if err != nil { return err }
	s.nodes = copyNodes
	s.store = store
	s.manager = memory.NewManager(store)
	return nil
}

func ReductionNodes(carrierID string, reduction evidence.Reduction) []memory.Node {
	out := make([]memory.Node, 0, len(reduction.Claims))
	for _, claim := range reduction.Claims {
		address := memory.CanonicalAddress(carrierID, "claims", cleanID(claim.ClaimKey))
		label := claim.Proposition
		if label == "" { label = claim.ClaimKey }
		summary := canonicalClaimSummary(claim)
		evidenceBody, _ := json.Marshal(struct {
			ClaimKey   string                  `json:"claim_key"`
			Status     evidence.ClaimStatus    `json:"status"`
			Exactness  evidence.Exactness      `json:"exactness"`
			Canonical  string                  `json:"canonical_id"`
			Support    []evidence.EvidenceAudit `json:"support_evidence,omitempty"`
			Opposition []evidence.EvidenceAudit `json:"oppose_evidence,omitempty"`
		}{claim.ClaimKey, claim.Status, claim.Exactness, claim.CanonicalID, claim.SupportEvidence, claim.OpposeEvidence})
		edges := make([]memory.Edge, 0, len(claim.SupportEvidence)+len(claim.OpposeEvidence))
		for _, audit := range claim.SupportEvidence {
			if audit.Resolved.Address != "" { edges = append(edges, memory.Edge{Relation: "supported-by", Target: audit.Resolved.Address, Weight: 1, Confidence: 1, ProofRef: claim.CanonicalID}) }
		}
		for _, audit := range claim.OpposeEvidence {
			if audit.Resolved.Address != "" { edges = append(edges, memory.Edge{Relation: "opposed-by", Target: audit.Resolved.Address, Weight: 1, Confidence: 1, ProofRef: claim.CanonicalID}) }
		}
		sort.Slice(edges, func(i, j int) bool { if edges[i].Relation == edges[j].Relation { return edges[i].Target < edges[j].Target }; return edges[i].Relation < edges[j].Relation })
		levels := map[memory.Fidelity]memory.ContentLevel{
			memory.FidelityAbstract: {Text: fmt.Sprintf("canonical claim state %s", claim.Status)},
			memory.FidelitySummary:  {Text: summary},
			memory.FidelityEvidence: {Text: string(evidenceBody), SourceRef: reduction.Digest, Verified: true},
		}
		if claim.Exactness == evidence.ExactnessVerified && claim.Status == evidence.StatusVerified && claim.Proposition != "" {
			levels[memory.FidelityExact] = memory.ContentLevel{Text: claim.Proposition, SourceRef: reduction.Digest, Verified: true}
		}
		out = append(out, memory.Node{
			Address:          address,
			CarrierID:        carrierID,
			ClusterID:        "claims",
			Kind:             "canonical-claim",
			SemanticClass:    strings.ToLower(string(claim.Status)),
			Label:            label,
			Levels:           levels,
			Edges:            edges,
			Importance:       memory.ImportanceProfile{Global: .8, Local: 1, Confidence: 1},
			Provenance:       memory.Provenance{SourceID: "origami:evidence-reduction", SourceRef: reduction.Schema, SourceHash: reduction.Digest},
			VerificationRef:  claim.CanonicalID,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Address < out[j].Address })
	return out
}

func canonicalClaimSummary(claim evidence.ClaimResult) string {
	proposition := claim.Proposition
	if proposition == "" { proposition = claim.ClaimKey }
	switch claim.Status {
	case evidence.StatusVerified:
		return proposition
	case evidence.StatusRejected:
		return "REJECTED: " + proposition
	case evidence.StatusConflict:
		return "CONFLICT: " + proposition
	default:
		return "UNKNOWN: " + proposition
	}
}

func cleanID(value string) string {
	value = strings.TrimSpace(value)
	var b strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' { b.WriteRune(r) } else { b.WriteByte('-') }
	}
	clean := strings.Trim(b.String(), "-")
	if clean == "" { return "claim" }
	return clean
}
