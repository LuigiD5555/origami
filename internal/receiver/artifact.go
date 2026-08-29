package receiver

import (
	"fmt"
	"strings"
)

const ArtifactSchemaR0 = "origami.receiver-artifact.r0"

type ArtifactStatus string

const (
	ArtifactReferenceCandidate ArtifactStatus = "REFERENCE_CANDIDATE"
	ArtifactPromoted           ArtifactStatus = "PROMOTED"
	ArtifactRejected           ArtifactStatus = "REJECTED"
)

// ArtifactManifest records the provenance of a receiver artifact stored by
// Origami. Tlaloc may produce the candidate, but Origami is the registry and
// semantic-validation authority for promoted receiver artifacts.
type ArtifactManifest struct {
	Schema             string         `json:"schema"`
	Status             ArtifactStatus `json:"status"`
	ReceiverContract   string         `json:"receiver_contract"`
	DistillerContract  string         `json:"distiller_contract,omitempty"`
	CandidateID        string         `json:"candidate_id"`
	PromptPath         string         `json:"prompt_path"`
	PromptSHA256       string         `json:"prompt_sha256"`
	SourceTraceSHA256  string         `json:"source_trace_sha256,omitempty"`
	ValidationSuite    []string       `json:"validation_suite"`
	FalseExact         int            `json:"false_exact"`
	Contaminated       bool           `json:"contaminated"`
	PeakActiveTokenEq  int            `json:"peak_active_token_eq"`
	WorkingWindowToken int            `json:"working_window_token_eq"`
	Notes              string         `json:"notes,omitempty"`
}

func ValidateArtifactManifest(m ArtifactManifest) error {
	if m.Schema != ArtifactSchemaR0 {
		return fmt.Errorf("artifact schema must be %q", ArtifactSchemaR0)
	}
	if m.ReceiverContract != ContractID {
		return fmt.Errorf("receiver contract must be %q", ContractID)
	}
	if m.CandidateID == "" || m.PromptPath == "" || m.PromptSHA256 == "" {
		return fmt.Errorf("candidate_id, prompt_path and prompt_sha256 are required")
	}
	if len(m.PromptSHA256) != 64 {
		return fmt.Errorf("prompt_sha256 must be a 64-character SHA-256 hex string")
	}
	if m.WorkingWindowToken <= 0 {
		return fmt.Errorf("working window must be positive")
	}
	if m.PeakActiveTokenEq < 0 {
		return fmt.Errorf("peak active token-equivalent cannot be negative")
	}
	if m.Status == ArtifactPromoted {
		if m.FalseExact != 0 {
			return fmt.Errorf("promoted receiver cannot have false exactness")
		}
		if m.Contaminated {
			return fmt.Errorf("promoted receiver cannot be contaminated")
		}
		if m.PeakActiveTokenEq > m.WorkingWindowToken {
			return fmt.Errorf("promoted receiver exceeds active interface budget")
		}
		if len(m.ValidationSuite) == 0 {
			return fmt.Errorf("promoted receiver requires validation evidence")
		}
	}
	if strings.Contains(m.PromptPath, "..") || strings.HasPrefix(m.PromptPath, "/") {
		return fmt.Errorf("prompt_path must be repository-relative and non-traversing")
	}
	switch m.Status {
	case ArtifactReferenceCandidate, ArtifactPromoted, ArtifactRejected:
		return nil
	default:
		return fmt.Errorf("unknown artifact status %q", m.Status)
	}
}
