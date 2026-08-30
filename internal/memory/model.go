package memory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	SchemaR0     = "origami.virtual-memory.r0"
	AddressABIR0 = "ohf-address.r0"
)

type Fidelity string

const (
	FidelityLabel    Fidelity = "label"
	FidelityAbstract Fidelity = "abstract"
	FidelitySummary  Fidelity = "summary"
	FidelityDetail   Fidelity = "detail"
	FidelityEvidence Fidelity = "evidence"
	FidelityExact    Fidelity = "exact"
)

type ContentLevel struct {
	Text       string `json:"text"`
	TokenCost  int    `json:"token_cost,omitempty"`
	SourceRef  string `json:"source_ref,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
}

type CostProfile struct {
	Label    int `json:"label,omitempty"`
	Abstract int `json:"abstract,omitempty"`
	Summary  int `json:"summary,omitempty"`
	Detail   int `json:"detail,omitempty"`
	Evidence int `json:"evidence,omitempty"`
	Exact    int `json:"exact,omitempty"`
}

type ImportanceProfile struct {
	Global     float64 `json:"global,omitempty"`
	Local      float64 `json:"local,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type Provenance struct {
	SourceID   string `json:"source_id,omitempty"`
	SourceRef  string `json:"source_ref,omitempty"`
	SourceHash string `json:"source_hash,omitempty"`
}

type Edge struct {
	Relation   string  `json:"relation"`
	Target     string  `json:"target"`
	Weight     float64 `json:"weight,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	ProofRef   string  `json:"proof_ref,omitempty"`
}

// Node is the model-agnostic addressable memory unit. Value and Links remain
// intentionally compatible with Hybrid Carrier R0. Rich nodes add multiple
// fidelity levels, semantic metadata and provenance without invalidating old
// carriers.
type Node struct {
	Address          string                  `json:"address"`
	CID              string                  `json:"cid,omitempty"`
	CarrierID        string                  `json:"carrier_id,omitempty"`
	ClusterID        string                  `json:"cluster_id,omitempty"`
	Kind             string                  `json:"kind,omitempty"`
	SemanticClass    string                  `json:"semantic_class,omitempty"`
	Label            string                  `json:"label,omitempty"`
	Value            string                  `json:"value,omitempty"`
	Levels           map[Fidelity]ContentLevel `json:"levels,omitempty"`
	Parent           string                  `json:"parent,omitempty"`
	Children         []string                `json:"children,omitempty"`
	Links            map[string][]string     `json:"links,omitempty"`
	Edges            []Edge                  `json:"edges,omitempty"`
	Importance       ImportanceProfile       `json:"importance,omitempty"`
	Costs            CostProfile             `json:"costs,omitempty"`
	Provenance       Provenance              `json:"provenance,omitempty"`
	VerificationRef  string                  `json:"verification_ref,omitempty"`
}

func EstimateTokens(text string) int {
	if text == "" { return 0 }
	return (len([]byte(text)) + 3) / 4
}

func (n Node) Level(f Fidelity) (ContentLevel, bool) {
	if f == FidelityLabel && n.Label != "" {
		return ContentLevel{Text: n.Label, TokenCost: EstimateTokens(n.Label)}, true
	}
	if n.Levels != nil {
		if level, ok := n.Levels[f]; ok && level.Text != "" {
			if level.TokenCost <= 0 { level.TokenCost = EstimateTokens(level.Text) }
			return level, true
		}
	}
	// Legacy R0 memory has one Value. Treat it as a summary-level view so old
	// carriers participate in virtual-memory queries without claiming evidence
	// or exactness that was never declared.
	if f == FidelitySummary && n.Value != "" {
		return ContentLevel{Text: n.Value, TokenCost: EstimateTokens(n.Value)}, true
	}
	if f == FidelityExact && n.Value != "" && len(n.Levels) == 0 {
		return ContentLevel{Text: n.Value, TokenCost: EstimateTokens(n.Value)}, true
	}
	return ContentLevel{}, false
}

func (n Node) Estimate(f Fidelity) int {
	var declared int
	switch f {
	case FidelityLabel: declared = n.Costs.Label
	case FidelityAbstract: declared = n.Costs.Abstract
	case FidelitySummary: declared = n.Costs.Summary
	case FidelityDetail: declared = n.Costs.Detail
	case FidelityEvidence: declared = n.Costs.Evidence
	case FidelityExact: declared = n.Costs.Exact
	}
	if declared > 0 { return declared }
	level, ok := n.Level(f)
	if !ok { return 0 }
	return level.TokenCost
}

func (n Node) DerivedCID() string {
	if n.CID != "" { return n.CID }
	h := sha256.New()
	fmt.Fprintf(h, "%s\x00%s\x00%s\x00%s\x00", n.Kind, n.SemanticClass, n.Label, n.Value)
	for _, f := range []Fidelity{FidelityAbstract, FidelitySummary, FidelityDetail, FidelityEvidence, FidelityExact} {
		if level, ok := n.Level(f); ok { fmt.Fprintf(h, "%s\x00%s\x00", f, level.Text) }
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

func ValidateNode(n Node) error {
	if strings.TrimSpace(n.Address) == "" { return fmt.Errorf("memory address cannot be empty") }
	for _, edge := range n.Edges {
		if edge.Relation == "" || edge.Target == "" { return fmt.Errorf("edge requires relation and target") }
	}
	return nil
}

func CanonicalAddress(carrierID, clusterID, nodeID string) string {
	clean := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "/", "-")
		if s == "" { return "_" }
		return s
	}
	return fmt.Sprintf("ohf://%s/%s/%s", clean(carrierID), clean(clusterID), clean(nodeID))
}

func IsCanonicalAddress(address string) bool { return strings.HasPrefix(address, "ohf://") }

func ContentID(text string) string {
	s := sha256.Sum256([]byte(text))
	return "sha256:" + hex.EncodeToString(s[:])
}
