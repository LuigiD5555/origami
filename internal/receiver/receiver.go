package receiver

import (
	"fmt"
	"sort"
)

const ContractID = "origami.hybrid-receiver.r0"

// BootSpec describes the minimum self-boot information a carrier must expose.
// The external model must not need prior knowledge of the carrier's private
// symbol assignments in order to begin decoding it.
type BootSpec struct {
	Version          string   `json:"version"`
	EntrySymbol      string   `json:"entry_symbol"`
	StartState       string   `json:"start_state"`
	RosettaRef       string   `json:"rosetta_ref"`
	ProgramRef       string   `json:"program_ref"`
	IndexRef         string   `json:"index_ref"`
	MemoryRef        string   `json:"memory_ref"`
	VerificationRef  string   `json:"verification_ref"`
	RequiredSections []string `json:"required_sections"`
	WorkingWindow    int      `json:"working_window_token_eq"`
}

// SymbolBinding is carrier-local. Physical symbols are deliberately not given
// global meaning; ROSETTA binds them to semantic tokens for this carrier.
type SymbolBinding struct {
	Physical string `json:"physical"`
	Semantic string `json:"semantic"`
}

// MicroRule is a deliberately small deterministic transition. Complex receiver
// behavior should emerge from many local rules rather than from hidden model
// reasoning. Action is preserved from the distilled swarm behavior even when
// the reference executor only needs it as declarative semantics/provenance.
type MicroRule struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Token     string `json:"token"`
	Action    string `json:"action,omitempty"`
	Emit      string `json:"emit,omitempty"`
	NextState string `json:"next_state"`
	Halt      bool   `json:"halt,omitempty"`
}

// Spec is the carrier-facing executable receiver contract. It is transport
// neutral: the same semantics can be projected visually, computationally, or
// through a hybrid carrier.
type Spec struct {
	ContractID string          `json:"contract_id"`
	Boot       BootSpec        `json:"boot"`
	Rosetta    []SymbolBinding `json:"rosetta"`
	Program    []MicroRule     `json:"program"`
}

type Step struct {
	Index     int    `json:"index"`
	Physical  string `json:"physical"`
	Semantic  string `json:"semantic"`
	RuleID    string `json:"rule_id"`
	Action    string `json:"action,omitempty"`
	FromState string `json:"from_state"`
	ToState   string `json:"to_state"`
	Emit      string `json:"emit,omitempty"`
}

type Trace struct {
	FinalState string   `json:"final_state"`
	Outputs    []string `json:"outputs"`
	Steps      []Step   `json:"steps"`
	Halted     bool     `json:"halted"`
	Unknown    bool     `json:"unknown"`
	Reason     string   `json:"reason,omitempty"`
}

func Validate(s Spec) error {
	if s.ContractID != ContractID {
		return fmt.Errorf("receiver contract must be %q", ContractID)
	}
	if s.Boot.Version == "" || s.Boot.EntrySymbol == "" || s.Boot.StartState == "" {
		return fmt.Errorf("boot version, entry symbol and start state are required")
	}
	if s.Boot.RosettaRef == "" || s.Boot.ProgramRef == "" || s.Boot.IndexRef == "" || s.Boot.MemoryRef == "" || s.Boot.VerificationRef == "" {
		return fmt.Errorf("boot must reference rosetta, program, index, memory and verification")
	}
	if s.Boot.WorkingWindow <= 0 {
		return fmt.Errorf("working window must be positive")
	}
	required := map[string]bool{"BOOT": false, "ROSETTA": false, "PROGRAM": false, "INDEX": false, "MEMORY": false, "VERIFICATION": false}
	for _, section := range s.Boot.RequiredSections {
		if _, ok := required[section]; ok {
			required[section] = true
		}
	}
	var missing []string
	for section, seen := range required {
		if !seen {
			missing = append(missing, section)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return fmt.Errorf("boot missing required sections: %v", missing)
	}
	physical := map[string]bool{}
	semantic := map[string]bool{}
	for _, b := range s.Rosetta {
		if b.Physical == "" || b.Semantic == "" {
			return fmt.Errorf("rosetta bindings require physical and semantic values")
		}
		if physical[b.Physical] {
			return fmt.Errorf("duplicate physical symbol %q", b.Physical)
		}
		if semantic[b.Semantic] {
			return fmt.Errorf("duplicate semantic token %q", b.Semantic)
		}
		physical[b.Physical] = true
		semantic[b.Semantic] = true
	}
	if len(s.Rosetta) == 0 || len(s.Program) == 0 {
		return fmt.Errorf("rosetta and micro-program are required")
	}
	ruleIDs := map[string]bool{}
	for _, r := range s.Program {
		if r.ID == "" || r.State == "" || r.Token == "" || r.NextState == "" {
			return fmt.Errorf("micro-rules require id, state, token and next_state")
		}
		if ruleIDs[r.ID] {
			return fmt.Errorf("duplicate micro-rule id %q", r.ID)
		}
		ruleIDs[r.ID] = true
		if !semantic[r.Token] {
			return fmt.Errorf("micro-rule %q references undeclared semantic token %q", r.ID, r.Token)
		}
	}
	return nil
}

// Execute interprets carrier-local physical symbols through the embedded
// ROSETTA and then runs the deterministic micro-program. Unknown physical
// symbols or missing transitions fail closed as UNKNOWN, not as guessed output.
func Execute(s Spec, physicalInput []string) (Trace, error) {
	if err := Validate(s); err != nil {
		return Trace{}, err
	}
	rosetta := make(map[string]string, len(s.Rosetta))
	for _, b := range s.Rosetta {
		rosetta[b.Physical] = b.Semantic
	}
	rules := make(map[string]MicroRule, len(s.Program))
	for _, r := range s.Program {
		key := r.State + "\x00" + r.Token
		if _, exists := rules[key]; exists {
			return Trace{}, fmt.Errorf("ambiguous transition for state=%q token=%q", r.State, r.Token)
		}
		rules[key] = r
	}

	trace := Trace{FinalState: s.Boot.StartState}
	state := s.Boot.StartState
	for i, physical := range physicalInput {
		semantic, ok := rosetta[physical]
		if !ok {
			trace.Unknown = true
			trace.Reason = fmt.Sprintf("unmapped physical symbol at position %d", i)
			trace.FinalState = state
			return trace, nil
		}
		rule, ok := rules[state+"\x00"+semantic]
		if !ok {
			trace.Unknown = true
			trace.Reason = fmt.Sprintf("no transition for state=%q token=%q", state, semantic)
			trace.FinalState = state
			return trace, nil
		}
		step := Step{Index: i, Physical: physical, Semantic: semantic, RuleID: rule.ID, Action: rule.Action, FromState: state, ToState: rule.NextState, Emit: rule.Emit}
		trace.Steps = append(trace.Steps, step)
		if rule.Emit != "" {
			trace.Outputs = append(trace.Outputs, rule.Emit)
		}
		state = rule.NextState
		trace.FinalState = state
		if rule.Halt {
			trace.Halted = true
			return trace, nil
		}
	}
	return trace, nil
}
