package core

import (
	"fmt"
	"sort"
)

type ObservationStatus string

const (
	ObservationPass            ObservationStatus = "PASS"
	ObservationFail            ObservationStatus = "FAIL"
	ObservationInvalidContract ObservationStatus = "INVALID_CONTRACT"
	ObservationUnsupported     ObservationStatus = "UNSUPPORTED"
)

type TemporalAtom struct {
	Key    string `json:"key"`
	Data   string `json:"data,omitempty"`
	Status Status `json:"status,omitempty"`
}

type ObservationContract struct {
	ID           string         `json:"id"`
	Kind         string         `json:"kind"` // state_predicate | temporal_sequence | rule_fired | participant_set
	Budget       int            `json:"budget"`
	StateStep    int            `json:"state_step,omitempty"`
	Predicate    *Predicate     `json:"predicate,omitempty"`
	Sequence     []TemporalAtom `json:"sequence,omitempty"`
	RuleID       string         `json:"rule_id,omitempty"`
	Participants []string       `json:"participants,omitempty"`
}

type ObservationResult struct {
	ContractID string            `json:"contract_id"`
	Status     ObservationStatus `json:"status"`
	TraceHash  string            `json:"trace_hash"`
	Inspected  int               `json:"inspected"`
	Evidence   []string          `json:"evidence,omitempty"`
	Reason     string            `json:"reason,omitempty"`
}

func Observe(trajectory Trajectory, contract ObservationContract) ObservationResult {
	result := ObservationResult{ContractID: contract.ID, TraceHash: trajectory.TraceHash}
	if contract.ID == "" || contract.Budget <= 0 {
		result.Status = ObservationInvalidContract
		result.Reason = "observation contract requires id and positive finite budget"
		return result
	}
	if trajectory.TraceHash == "" || trajectory.TraceHash != trajectoryHash(trajectory) {
		result.Status = ObservationInvalidContract
		result.Reason = "trajectory hash is missing or invalid"
		return result
	}

	switch contract.Kind {
	case "state_predicate":
		return observeStatePredicate(trajectory, contract, result)
	case "temporal_sequence":
		return observeTemporalSequence(trajectory, contract, result)
	case "rule_fired":
		return observeRuleFired(trajectory, contract, result)
	case "participant_set":
		return observeParticipantSet(trajectory, contract, result)
	default:
		result.Status = ObservationUnsupported
		result.Reason = fmt.Sprintf("observation kind %q is not implemented", contract.Kind)
		return result
	}
}

func observeStatePredicate(trajectory Trajectory, contract ObservationContract, result ObservationResult) ObservationResult {
	if contract.Predicate == nil || contract.Predicate.Scope != "state" {
		result.Status = ObservationInvalidContract
		result.Reason = "state_predicate requires a state-scoped predicate"
		return result
	}
	if err := validatePredicate(*contract.Predicate); err != nil {
		result.Status = ObservationInvalidContract
		result.Reason = err.Error()
		return result
	}
	step := contract.StateStep
	if step < 0 || step >= len(trajectory.States) {
		result.Status = ObservationInvalidContract
		result.Reason = "state_step outside trajectory"
		return result
	}
	if contract.Budget < 1 {
		result.Status = ObservationFail
		return result
	}
	result.Inspected = 1
	state := trajectory.States[step]
	if matches([]Predicate{*contract.Predicate}, state, Context{Values: map[string]Value{}}) {
		result.Status = ObservationPass
		result.Evidence = []string{fmt.Sprintf("state:%d:%s", step, StateHash(state))}
		return result
	}
	result.Status = ObservationFail
	result.Reason = "predicate not satisfied"
	return result
}

func observeTemporalSequence(trajectory Trajectory, contract ObservationContract, result ObservationResult) ObservationResult {
	if len(contract.Sequence) < 2 {
		result.Status = ObservationInvalidContract
		result.Reason = "temporal_sequence requires at least two atoms"
		return result
	}
	limit := len(trajectory.States)
	if limit > contract.Budget {
		limit = contract.Budget
	}
	result.Inspected = limit
	if limit < len(contract.Sequence) {
		result.Status = ObservationFail
		result.Reason = "finite observation budget is shorter than required sequence"
		return result
	}
	for start := 0; start+len(contract.Sequence) <= limit; start++ {
		ok := true
		var evidence []string
		for offset, atom := range contract.Sequence {
			state := trajectory.States[start+offset]
			value, exists := state.Values[atom.Key]
			if !exists {
				value = Value{Status: StatusUnknown}
			}
			if atom.Status != "" && value.Status != atom.Status {
				ok = false
				break
			}
			if atom.Data != "" && value.Data != atom.Data {
				ok = false
				break
			}
			evidence = append(evidence, fmt.Sprintf("state:%d:%s:%s", start+offset, atom.Key, StateHash(state)))
		}
		if ok {
			result.Status = ObservationPass
			result.Evidence = evidence
			return result
		}
	}
	result.Status = ObservationFail
	result.Reason = "ordered temporal predicate not satisfied within budget"
	return result
}

func observeRuleFired(trajectory Trajectory, contract ObservationContract, result ObservationResult) ObservationResult {
	if contract.RuleID == "" {
		result.Status = ObservationInvalidContract
		result.Reason = "rule_fired requires rule_id"
		return result
	}
	limit := len(trajectory.Evidence)
	if limit > contract.Budget {
		limit = contract.Budget
	}
	result.Inspected = limit
	for i := 0; i < limit; i++ {
		for _, fired := range trajectory.Evidence[i].FiredRules {
			if fired.RuleID == contract.RuleID {
				result.Status = ObservationPass
				result.Evidence = []string{fmt.Sprintf("transition:%d:rule:%s", trajectory.Evidence[i].Step, fired.RuleID)}
				return result
			}
		}
	}
	result.Status = ObservationFail
	result.Reason = "rule did not fire within budget"
	return result
}

func observeParticipantSet(trajectory Trajectory, contract ObservationContract, result ObservationResult) ObservationResult {
	if len(contract.Participants) == 0 {
		result.Status = ObservationInvalidContract
		result.Reason = "participant_set requires participants"
		return result
	}
	wanted := append([]string(nil), contract.Participants...)
	sort.Strings(wanted)
	limit := len(trajectory.Evidence)
	if limit > contract.Budget {
		limit = contract.Budget
	}
	result.Inspected = limit
	for i := 0; i < limit; i++ {
		for _, fired := range trajectory.Evidence[i].FiredRules {
			participants := append([]string(nil), fired.Participants...)
			sort.Strings(participants)
			if equalStrings(participants, wanted) {
				result.Status = ObservationPass
				result.Evidence = []string{fmt.Sprintf("transition:%d:participants:%v", trajectory.Evidence[i].Step, participants)}
				return result
			}
		}
	}
	result.Status = ObservationFail
	result.Reason = "exact participant set not observed within budget"
	return result
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
