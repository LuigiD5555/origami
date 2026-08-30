package core

import "fmt"

type RepresentationAccounting struct {
	SourceBytes            int `json:"source_bytes"`
	SemanticBytes          int `json:"semantic_bytes"`
	LiteralBytes           int `json:"literal_bytes"`
	ResidualBytes          int `json:"residual_bytes"`
	MachineDerivedClaims   int `json:"machine_derived_claims"`
	VerifiedSemanticClaims int `json:"verified_semantic_claims"`
}

type ClaimClass string

const (
	ClaimSemantic ClaimClass = "SEMANTIC"
	ClaimExact    ClaimClass = "EXACT"
)

type AccountingReport struct {
	Class             ClaimClass `json:"class"`
	Allowed           bool       `json:"allowed"`
	Reason            string     `json:"reason,omitempty"`
	LiteralRatio      float64    `json:"literal_ratio"`
	ResidualRatio     float64    `json:"residual_ratio"`
	SemanticAuthority bool       `json:"semantic_authority"`
}

func EvaluateAccounting(class ClaimClass, accounting RepresentationAccounting) AccountingReport {
	report := AccountingReport{Class: class}
	denominator := accounting.SourceBytes
	if denominator <= 0 {
		denominator = accounting.SemanticBytes + accounting.LiteralBytes + accounting.ResidualBytes
	}
	if denominator > 0 {
		report.LiteralRatio = float64(accounting.LiteralBytes) / float64(denominator)
		report.ResidualRatio = float64(accounting.ResidualBytes) / float64(denominator)
	}
	report.SemanticAuthority = accounting.SemanticBytes > 0 && accounting.MachineDerivedClaims > 0

	switch class {
	case ClaimSemantic:
		if !report.SemanticAuthority {
			report.Reason = "semantic claims require represented semantic state and machine-derived evidence; literal/residual transport is not semantic authority"
			return report
		}
		if accounting.VerifiedSemanticClaims <= 0 {
			report.Reason = "semantic claims require at least one verified semantic claim"
			return report
		}
		report.Allowed = true
		return report
	case ClaimExact:
		if accounting.SourceBytes <= 0 {
			report.Reason = "exactness accounting requires declared source bytes"
			return report
		}
		if accounting.LiteralBytes+accounting.ResidualBytes < 0 {
			report.Reason = "literal/residual accounting cannot be negative"
			return report
		}
		report.Allowed = true
		return report
	default:
		report.Reason = fmt.Sprintf("unknown claim class %q", class)
		return report
	}
}
