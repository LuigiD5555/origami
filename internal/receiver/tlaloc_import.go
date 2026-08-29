package receiver

import "fmt"

const TlalocHybridArtifactSchema = "tlaloc.origami-hybrid-artifact-set.r0"

type TlalocMicroRule struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Token     string `json:"token"`
	Action    string `json:"action,omitempty"`
	NextState string `json:"next_state"`
	Emit      string `json:"emit,omitempty"`
}

// TlalocHybridArtifactSet is an import DTO only. Accepting this proposal does
// not promote it and does not allow Tlaloc to define carrier-local glyphs or
// Origami semantics.
type TlalocHybridArtifactSet struct {
	Schema             string            `json:"schema"`
	CandidateID        string            `json:"candidate_id"`
	UniversalPrompt    string            `json:"universal_prompt"`
	BootStrategy       []string          `json:"boot_strategy"`
	RosettaConstraints []string          `json:"rosetta_constraints"`
	MicroProgram       []TlalocMicroRule `json:"micro_program"`
	SourceTraceSHA256  string            `json:"source_trace_sha256"`
	WorkingWindow      int               `json:"working_window_token_eq"`
}

type ImportedCandidate struct {
	CandidateID        string      `json:"candidate_id"`
	UniversalPrompt    string      `json:"universal_prompt"`
	BootStrategy       []string    `json:"boot_strategy"`
	RosettaConstraints []string    `json:"rosetta_constraints"`
	Program            []MicroRule `json:"program"`
	SourceTraceSHA256  string      `json:"source_trace_sha256"`
	WorkingWindow      int         `json:"working_window_token_eq"`
}

// ImportTlalocCandidate preserves the behavior Tlaloc distilled while keeping
// the authority boundary explicit: Tlaloc proposes receiver behavior; Origami
// later binds carrier-local physical symbols, constructs BOOT/INDEX/MEMORY and
// validates the final Spec before any artifact may be promoted.
func ImportTlalocCandidate(in TlalocHybridArtifactSet) (ImportedCandidate, error) {
	if in.Schema != TlalocHybridArtifactSchema {
		return ImportedCandidate{}, fmt.Errorf("candidate schema must be %q", TlalocHybridArtifactSchema)
	}
	if in.CandidateID == "" || in.UniversalPrompt == "" || in.SourceTraceSHA256 == "" {
		return ImportedCandidate{}, fmt.Errorf("candidate id, universal prompt and source trace hash are required")
	}
	if in.WorkingWindow <= 0 {
		return ImportedCandidate{}, fmt.Errorf("working window must be positive")
	}
	if len(in.BootStrategy) == 0 || len(in.RosettaConstraints) == 0 || len(in.MicroProgram) == 0 {
		return ImportedCandidate{}, fmt.Errorf("boot strategy, rosetta constraints and micro-program are required")
	}
	program := make([]MicroRule, 0, len(in.MicroProgram))
	seen := map[string]bool{}
	for _, r := range in.MicroProgram {
		if r.ID == "" || r.State == "" || r.Token == "" || r.NextState == "" {
			return ImportedCandidate{}, fmt.Errorf("imported micro-rules require id, state, token and next_state")
		}
		if seen[r.ID] {
			return ImportedCandidate{}, fmt.Errorf("duplicate imported micro-rule id %q", r.ID)
		}
		seen[r.ID] = true
		program = append(program, MicroRule{
			ID:        r.ID,
			State:     r.State,
			Token:     r.Token,
			Action:    r.Action,
			Emit:      r.Emit,
			NextState: r.NextState,
		})
	}
	return ImportedCandidate{
		CandidateID:        in.CandidateID,
		UniversalPrompt:    in.UniversalPrompt,
		BootStrategy:       append([]string(nil), in.BootStrategy...),
		RosettaConstraints: append([]string(nil), in.RosettaConstraints...),
		Program:            program,
		SourceTraceSHA256:  in.SourceTraceSHA256,
		WorkingWindow:      in.WorkingWindow,
	}, nil
}
