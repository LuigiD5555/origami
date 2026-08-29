package reference

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// SemanticStatus makes non-equivalent semantic conditions explicit instead of
// encoding them as arbitrary strings.
type SemanticStatus string

const (
	StatusPresent   SemanticStatus = "PRESENT"
	StatusAbsent    SemanticStatus = "ABSENT"
	StatusUnknown   SemanticStatus = "UNKNOWN"
	StatusInhibited SemanticStatus = "INHIBITED"
	StatusCancelled SemanticStatus = "CANCELLED"
)

type SemanticValue struct {
	Status SemanticStatus `json:"status"`
	Value  string         `json:"value,omitempty"`
	Reason string         `json:"reason,omitempty"`
}

type CoreState map[string]SemanticValue

type Context struct {
	ID     string            `json:"id"`
	Values map[string]string `json:"values,omitempty"`
}

type Predicate struct {
	Source string `json:"source"` // state or context
	Key    string `json:"key"`
	Equals string `json:"equals"`
}

type Effect struct {
	Entity string        `json:"entity"`
	Status SemanticStatus `json:"status"`
	Value  string        `json:"value,omitempty"`
	Reason string        `json:"reason,omitempty"`
}

type Rule struct {
	ID   string      `json:"id"`
	When []Predicate `json:"when"`
	Then []Effect    `json:"then"`
}

type TransitionEvidence struct {
	Step       int      `json:"step"`
	ContextID  string   `json:"context_id"`
	FiredRules []string `json:"fired_rules"`
}

type CoreRun struct {
	Initial    CoreState            `json:"initial"`
	Trace      []CoreState          `json:"trace"`
	Transitions []TransitionEvidence `json:"transitions"`
}

type ObservationContract struct {
	ID        string `json:"id"`
	Authority string `json:"authority"`
	Entity    string `json:"entity"`
}

type ObservationResult struct {
	ContractID string        `json:"contract_id"`
	Authority  string        `json:"authority"`
	Entity     string        `json:"entity"`
	Value      SemanticValue `json:"value"`
	TraceHash  string        `json:"trace_sha256"`
}

func cloneCoreState(s CoreState) CoreState {
	r := CoreState{}
	for k, v := range s { r[k] = v }
	return r
}

func predicateMatches(p Predicate, state CoreState, context Context) bool {
	switch p.Source {
	case "context":
		return context.Values[p.Key] == p.Equals
	case "state":
		v, ok := state[p.Key]
		if !ok { return false }
		return string(v.Status) == p.Equals || v.Value == p.Equals
	default:
		return false
	}
}

func ruleMatches(r Rule, state CoreState, context Context) bool {
	for _, p := range r.When {
		if !predicateMatches(p, state, context) { return false }
	}
	return true
}

// ApplyRules performs one deterministic transition from declared finite rules.
// Rules are evaluated against the same pre-step state. Effects are then applied
// in declaration order; conflicting writes are rejected instead of silently
// becoming priority semantics.
func ApplyRules(state CoreState, context Context, rules []Rule) (CoreState, []string, error) {
	next := cloneCoreState(state)
	writes := map[string]SemanticValue{}
	var fired []string
	for _, r := range rules {
		if r.ID == "" { return nil, nil, fmt.Errorf("rule id must not be empty") }
		if !ruleMatches(r, state, context) { continue }
		fired = append(fired, r.ID)
		for _, e := range r.Then {
			v := SemanticValue{Status:e.Status, Value:e.Value, Reason:e.Reason}
			if prior, ok := writes[e.Entity]; ok && prior != v {
				return nil, nil, fmt.Errorf("conflicting effects for %s", e.Entity)
			}
			writes[e.Entity] = v
		}
	}
	for entity, value := range writes { next[entity] = value }
	return next, fired, nil
}

func ExecuteCore(initial CoreState, contexts []Context, rules []Rule) (CoreRun, error) {
	cur := cloneCoreState(initial)
	run := CoreRun{Initial:cloneCoreState(initial), Trace:[]CoreState{cloneCoreState(initial)}}
	for i, context := range contexts {
		next, fired, err := ApplyRules(cur, context, rules)
		if err != nil { return CoreRun{}, err }
		run.Transitions = append(run.Transitions, TransitionEvidence{Step:i+1, ContextID:context.ID, FiredRules:fired})
		run.Trace = append(run.Trace, cloneCoreState(next))
		cur = next
	}
	return run, nil
}

func TraceHash(trace []CoreState) (string, error) {
	b, err := json.Marshal(trace)
	if err != nil { return "", err }
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// Observe reads immutable transition evidence. It does not execute a transition
// and cannot rewrite the run trace.
func Observe(run CoreRun, contract ObservationContract) (ObservationResult, error) {
	if contract.ID == "" || contract.Authority == "" { return ObservationResult{}, fmt.Errorf("observation contract requires id and authority") }
	if len(run.Trace) == 0 { return ObservationResult{}, fmt.Errorf("run has no trace") }
	final := run.Trace[len(run.Trace)-1]
	v, ok := final[contract.Entity]
	if !ok { v = SemanticValue{Status:StatusAbsent} }
	h, err := TraceHash(run.Trace)
	if err != nil { return ObservationResult{}, err }
	return ObservationResult{ContractID:contract.ID,Authority:contract.Authority,Entity:contract.Entity,Value:v,TraceHash:h}, nil
}
