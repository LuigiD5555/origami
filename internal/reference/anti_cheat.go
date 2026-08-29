package reference

import (
	"fmt"
	"strings"
)

type DerivationKind string

const (
	DerivedTransition DerivationKind = "TRANSITION"
	DerivedInteraction DerivationKind = "INTERACTION"
	DerivedTemporal DerivationKind = "TEMPORAL_OBSERVATION"
	DerivedSelectiveUnfold DerivationKind = "SELECTIVE_UNFOLD"
	DerivedLiteral DerivationKind = "LITERAL"
	DerivedResidual DerivationKind = "RESIDUAL"
	DerivedExternal DerivationKind = "EXTERNAL_MODEL"
)

type SemanticClaim struct {
	ID string `json:"id"`
	Entity string `json:"entity"`
	Value SemanticValue `json:"value"`
	Derivation DerivationKind `json:"derivation"`
	EvidenceIDs []string `json:"evidence_ids,omitempty"`
	LiteralBytes int `json:"literal_bytes,omitempty"`
	ResidualBytes int `json:"residual_bytes,omitempty"`
	SourceBytes int `json:"source_bytes,omitempty"`
}

type RouteIntegrityPolicy struct {
	MaxLiteralFraction float64 `json:"max_literal_fraction"`
	MaxResidualFraction float64 `json:"max_residual_fraction"`
	RequireEvidence bool `json:"require_evidence"`
	ForbidExternalCore bool `json:"forbid_external_core"`
}

type RouteIntegrityResult struct { Passed bool `json:"passed"`; Violations []string `json:"violations,omitempty"` }

// CheckRouteIntegrity does not prove semantic correctness by itself. It rejects
// known easy routes that would make a test meaningless: stored answers,
// source-sized literals/residuals, or delegating core derivation to a model.
func CheckRouteIntegrity(claim SemanticClaim,policy RouteIntegrityPolicy) RouteIntegrityResult {
	var violations []string
	if claim.ID==""||claim.Entity==""{violations=append(violations,"claim identity is incomplete")}
	if policy.RequireEvidence&&len(claim.EvidenceIDs)==0{violations=append(violations,"semantic claim has no derivation evidence")}
	if claim.Derivation==DerivedLiteral{violations=append(violations,"final semantic claim is stored as literal")}
	if claim.Derivation==DerivedExternal&&policy.ForbidExternalCore{violations=append(violations,"core semantic derivation delegated to external model")}
	if claim.SourceBytes>0{
		if float64(claim.LiteralBytes)/float64(claim.SourceBytes)>policy.MaxLiteralFraction{violations=append(violations,"literal payload exceeds semantic-route budget")}
		if float64(claim.ResidualBytes)/float64(claim.SourceBytes)>policy.MaxResidualFraction{violations=append(violations,"residual payload exceeds semantic-route budget")}
	}
	return RouteIntegrityResult{Passed:len(violations)==0,Violations:violations}
}

func RequireRouteIntegrity(claim SemanticClaim,policy RouteIntegrityPolicy) error {r:=CheckRouteIntegrity(claim,policy);if r.Passed{return nil};return fmt.Errorf("route integrity failed: %s",strings.Join(r.Violations,"; "))}
