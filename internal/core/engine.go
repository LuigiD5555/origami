package core

import (
	"fmt"
	"sort"
	"strings"
)

type appliedEffect struct {
	RuleID string
	Effect Effect
}

func Step(state State, context Context, machine Machine, stepIndex int) (State, TransitionEvidence, error) {
	if err := ValidateState(state); err != nil {
		return State{}, TransitionEvidence{}, err
	}
	if err := validateContext(context); err != nil {
		return State{}, TransitionEvidence{}, err
	}
	if err := ValidateMachine(machine); err != nil {
		return State{}, TransitionEvidence{}, err
	}

	rules := append([]Rule(nil), machine.Rules...)
	sort.Slice(rules, func(i, j int) bool { return rules[i].ID < rules[j].ID })

	fires := make([]RuleFire, 0)
	byKey := map[string][]appliedEffect{}
	for _, rule := range rules {
		if !matches(rule.When, state, context) {
			continue
		}
		participants := predicateParticipants(rule.When)
		effects := append([]Effect(nil), rule.Then...)
		sort.Slice(effects, func(i, j int) bool { return effectKey(effects[i]) < effectKey(effects[j]) })
		fires = append(fires, RuleFire{RuleID: rule.ID, Participants: participants, Effects: effects})
		for _, effect := range effects {
			byKey[effect.Key] = append(byKey[effect.Key], appliedEffect{RuleID: rule.ID, Effect: effect})
		}
	}

	next := CloneState(state)
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	resolutions := make([]Resolution, 0, len(keys))
	for _, key := range keys {
		value, resolution := resolveEffects(key, byKey[key])
		next.Values[key] = value
		resolutions = append(resolutions, resolution)
	}

	evidence := TransitionEvidence{
		Step:            stepIndex,
		ContextID:       context.ID,
		ContextHash:     ContextHash(context),
		StateBeforeHash: StateHash(state),
		StateAfterHash:  StateHash(next),
		FiredRules:      fires,
		Resolutions:     resolutions,
	}
	return next, evidence, nil
}

func Execute(initial State, machine Machine, contexts []Context) (Trajectory, error) {
	if machine.Schema == "" {
		machine.Schema = SchemaR1 + ".machine"
	}
	if err := ValidateMachine(machine); err != nil {
		return Trajectory{}, err
	}
	if err := ValidateState(initial); err != nil {
		return Trajectory{}, err
	}

	trajectory := Trajectory{
		Schema:  SchemaR1 + ".trajectory",
		Initial: CloneState(initial),
		States:  []State{CloneState(initial)},
	}
	current := CloneState(initial)

	if len(contexts) > 0 {
		limit := len(contexts)
		if limit > machine.MaxSteps {
			limit = machine.MaxSteps
		}
		for i := 0; i < limit; i++ {
			next, evidence, err := Step(current, CloneContext(contexts[i]), machine, i+1)
			if err != nil {
				return Trajectory{}, err
			}
			trajectory.States = append(trajectory.States, CloneState(next))
			trajectory.Evidence = append(trajectory.Evidence, evidence)
			current = next
		}
		if len(contexts) > machine.MaxSteps {
			trajectory.Terminated = "BUDGET_EXHAUSTED"
		} else {
			trajectory.Terminated = "CONTEXTS_EXHAUSTED"
		}
		trajectory.TraceHash = trajectoryHash(trajectory)
		return trajectory, nil
	}

	empty := Context{ID: "default", Values: map[string]Value{}}
	seen := map[string]struct{}{StateHash(current): {}}
	for i := 0; i < machine.MaxSteps; i++ {
		next, evidence, err := Step(current, empty, machine, i+1)
		if err != nil {
			return Trajectory{}, err
		}
		trajectory.States = append(trajectory.States, CloneState(next))
		trajectory.Evidence = append(trajectory.Evidence, evidence)
		if StateHash(next) == StateHash(current) {
			trajectory.Terminated = "FIXED_POINT"
			trajectory.TraceHash = trajectoryHash(trajectory)
			return trajectory, nil
		}
		hash := StateHash(next)
		if _, ok := seen[hash]; ok {
			trajectory.Terminated = "CYCLE"
			trajectory.TraceHash = trajectoryHash(trajectory)
			return trajectory, nil
		}
		seen[hash] = struct{}{}
		current = next
	}
	trajectory.Terminated = "BUDGET_EXHAUSTED"
	trajectory.TraceHash = trajectoryHash(trajectory)
	return trajectory, nil
}

func validateContext(context Context) error {
	for key, value := range context.Values {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("context key cannot be empty")
		}
		if !ValidateStatus(value.Status) {
			return fmt.Errorf("context %q has invalid status %q", key, value.Status)
		}
	}
	return nil
}

func matches(predicates []Predicate, state State, context Context) bool {
	for _, predicate := range predicates {
		var value Value
		var ok bool
		if predicate.Scope == "context" {
			value, ok = context.Values[predicate.Key]
		} else {
			value, ok = state.Values[predicate.Key]
		}
		if !ok {
			value = Value{Status: StatusUnknown}
		}
		switch predicate.Op {
		case "data_eq":
			if value.Data != predicate.Data {
				return false
			}
		case "data_ne":
			if value.Data == predicate.Data {
				return false
			}
		case "status_eq":
			if value.Status != predicate.Status {
				return false
			}
		case "present":
			if value.Status != StatusPresent {
				return false
			}
		case "not_present":
			if value.Status == StatusPresent {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func predicateParticipants(predicates []Predicate) []string {
	seen := map[string]struct{}{}
	for _, predicate := range predicates {
		prefix := predicate.Scope + ":"
		seen[prefix+predicate.Key] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func effectKey(effect Effect) string {
	return fmt.Sprintf("%s|%s|%s|%s|%d", effect.Key, effect.Mode, effect.Status, effect.Data, effect.Delta)
}

func resolveEffects(key string, effects []appliedEffect) (Value, Resolution) {
	sort.Slice(effects, func(i, j int) bool {
		if effects[i].RuleID == effects[j].RuleID {
			return effectKey(effects[i].Effect) < effectKey(effects[j].Effect)
		}
		return effects[i].RuleID < effects[j].RuleID
	})
	ruleIDs := make([]string, 0, len(effects))
	seenRule := map[string]struct{}{}
	for _, item := range effects {
		if _, ok := seenRule[item.RuleID]; !ok {
			seenRule[item.RuleID] = struct{}{}
			ruleIDs = append(ruleIDs, item.RuleID)
		}
	}

	for _, item := range effects {
		if item.Effect.Mode == "inhibit" {
			value := Value{Status: StatusInhibited, Data: item.Effect.Data}
			return value, Resolution{Key: key, Kind: "inhibition", RuleIDs: ruleIDs, Result: value}
		}
	}

	sets := make([]Value, 0)
	contributions := make([]int, 0)
	contributionData := ""
	mixedContributionData := false
	for _, item := range effects {
		switch item.Effect.Mode {
		case "set":
			sets = append(sets, Value{Status: item.Effect.Status, Data: item.Effect.Data})
		case "contribute":
			contributions = append(contributions, item.Effect.Delta)
			if contributionData == "" {
				contributionData = item.Effect.Data
			} else if contributionData != item.Effect.Data {
				mixedContributionData = true
			}
		}
	}

	if len(sets) > 0 && len(contributions) > 0 {
		value := Value{Status: StatusUnknown}
		return value, Resolution{Key: key, Kind: "conflict", RuleIDs: ruleIDs, Contributions: contributions, Result: value}
	}
	if len(sets) > 0 {
		first := sets[0]
		for _, candidate := range sets[1:] {
			if candidate != first {
				value := Value{Status: StatusUnknown}
				return value, Resolution{Key: key, Kind: "conflict", RuleIDs: ruleIDs, Result: value}
			}
		}
		return first, Resolution{Key: key, Kind: "set", RuleIDs: ruleIDs, Result: first}
	}
	if len(contributions) > 0 {
		if mixedContributionData {
			value := Value{Status: StatusUnknown}
			return value, Resolution{Key: key, Kind: "conflict", RuleIDs: ruleIDs, Contributions: contributions, Result: value}
		}
		total := 0
		for _, delta := range contributions {
			total += delta
		}
		var value Value
		kind := "contribution"
		switch {
		case total == 0:
			value = Value{Status: StatusCancelled, Data: contributionData}
			kind = "cancellation"
		case total > 0:
			value = Value{Status: StatusPresent, Data: contributionData}
		case total < 0:
			value = Value{Status: StatusAbsent, Data: contributionData}
		}
		return value, Resolution{Key: key, Kind: kind, RuleIDs: ruleIDs, Contributions: contributions, Result: value}
	}

	value := Value{Status: StatusUnknown}
	return value, Resolution{Key: key, Kind: "conflict", RuleIDs: ruleIDs, Result: value}
}

func trajectoryHash(trajectory Trajectory) string {
	clone := trajectory
	clone.TraceHash = ""
	return hashJSON(clone)
}
