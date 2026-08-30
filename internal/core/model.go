package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const SchemaR1 = "origami.core.r1"

type Status string
const(StatusPresent Status="PRESENT";StatusAbsent Status="ABSENT";StatusUnknown Status="UNKNOWN";StatusInhibited Status="INHIBITED";StatusCancelled Status="CANCELLED")
type Value struct{Status Status `json:"status"`;Data string `json:"data,omitempty"`}
type State struct{Values map[string]Value `json:"values"`}
type Context struct{ID string `json:"id"`;Values map[string]Value `json:"values,omitempty"`}
type Predicate struct{Scope string `json:"scope"`;Key string `json:"key"`;Op string `json:"op"`;Data string `json:"data,omitempty"`;Status Status `json:"status,omitempty"`}
type Effect struct{Key string `json:"key"`;Mode string `json:"mode"`;Data string `json:"data,omitempty"`;Status Status `json:"status,omitempty"`;Delta int `json:"delta,omitempty"`}
type Contradiction struct{Reason string `json:"reason"`}
type Rule struct{ID string `json:"id"`;When []Predicate `json:"when,omitempty"`;Then []Effect `json:"then,omitempty"`;Contradiction *Contradiction `json:"contradiction,omitempty"`}
type Machine struct{Schema string `json:"schema"`;Rules []Rule `json:"rules"`;MaxSteps int `json:"max_steps"`}
type RuleFire struct{RuleID string `json:"rule_id"`;Participants []string `json:"participants,omitempty"`;Effects []Effect `json:"effects,omitempty"`}
type ContradictionHit struct{RuleID string `json:"rule_id"`;Reason string `json:"reason"`;Participants []string `json:"participants,omitempty"`}
type Resolution struct{Key string `json:"key"`;Kind string `json:"kind"`;RuleIDs []string `json:"rule_ids,omitempty"`;Contributions []int `json:"contributions,omitempty"`;Result Value `json:"result"`}
type TransitionEvidence struct{Step int `json:"step"`;ContextID string `json:"context_id,omitempty"`;ContextHash string `json:"context_hash"`;StateBeforeHash string `json:"state_before_hash"`;StateAfterHash string `json:"state_after_hash"`;FiredRules []RuleFire `json:"fired_rules,omitempty"`;Contradictions []ContradictionHit `json:"contradictions,omitempty"`;Resolutions []Resolution `json:"resolutions,omitempty"`}
type Trajectory struct{Schema string `json:"schema"`;Initial State `json:"initial"`;States []State `json:"states"`;Evidence []TransitionEvidence `json:"evidence"`;TraceHash string `json:"trace_hash"`;Terminated string `json:"terminated"`}

func NewMachine(rules []Rule,maxSteps int)Machine{if maxSteps<=0{maxSteps=32};return Machine{Schema:SchemaR1+".machine",Rules:append([]Rule(nil),rules...),MaxSteps:maxSteps}}
func CloneState(s State)State{out:=State{Values:map[string]Value{}};for k,v:=range s.Values{out.Values[k]=v};return out}
func CloneContext(c Context)Context{out:=Context{ID:c.ID,Values:map[string]Value{}};for k,v:=range c.Values{out.Values[k]=v};return out}
func ValidateStatus(status Status)bool{switch status{case StatusPresent,StatusAbsent,StatusUnknown,StatusInhibited,StatusCancelled:return true;default:return false}}
func ValidateState(s State)error{if s.Values==nil{return fmt.Errorf("state values cannot be nil")};for key,value:=range s.Values{if strings.TrimSpace(key)==""{return fmt.Errorf("state key cannot be empty")};if !ValidateStatus(value.Status){return fmt.Errorf("state %q has invalid status %q",key,value.Status)}};return nil}
func ValidateMachine(m Machine)error{
	if m.Schema!=""&&m.Schema!=SchemaR1+".machine"{return fmt.Errorf("unexpected machine schema %q",m.Schema)};if m.MaxSteps<=0{return fmt.Errorf("max_steps must be positive")};seen:=map[string]struct{}{}
	for _,rule:=range m.Rules{if strings.TrimSpace(rule.ID)==""{return fmt.Errorf("rule id cannot be empty")};if _,ok:=seen[rule.ID];ok{return fmt.Errorf("duplicate rule id %q",rule.ID)};seen[rule.ID]=struct{}{};if len(rule.Then)==0&&rule.Contradiction==nil{return fmt.Errorf("rule %q has neither effects nor contradiction",rule.ID)};if len(rule.Then)>0&&rule.Contradiction!=nil{return fmt.Errorf("rule %q cannot mutate state and declare contradiction simultaneously",rule.ID)};if rule.Contradiction!=nil&&strings.TrimSpace(rule.Contradiction.Reason)==""{return fmt.Errorf("rule %q contradiction requires reason",rule.ID)};for _,predicate:=range rule.When{if err:=validatePredicate(predicate);err!=nil{return fmt.Errorf("rule %q: %w",rule.ID,err)}};for _,effect:=range rule.Then{if err:=validateEffect(effect);err!=nil{return fmt.Errorf("rule %q: %w",rule.ID,err)}}};return nil
}
func validatePredicate(p Predicate)error{if p.Scope!="state"&&p.Scope!="context"{return fmt.Errorf("predicate scope must be state or context")};if strings.TrimSpace(p.Key)==""{return fmt.Errorf("predicate key cannot be empty")};switch p.Op{case "data_eq","data_ne","present","not_present":case "status_eq":if !ValidateStatus(p.Status){return fmt.Errorf("invalid predicate status %q",p.Status)};default:return fmt.Errorf("unsupported predicate op %q",p.Op)};return nil}
func validateEffect(e Effect)error{if strings.TrimSpace(e.Key)==""{return fmt.Errorf("effect key cannot be empty")};switch e.Mode{case "set":if !ValidateStatus(e.Status){return fmt.Errorf("set effect requires valid status")};case "contribute":if e.Delta==0{return fmt.Errorf("contribute effect delta cannot be zero")};case "inhibit":default:return fmt.Errorf("unsupported effect mode %q",e.Mode)};return nil}
func StateHash(s State)string{return hashJSON(s)}
func ContextHash(c Context)string{return hashJSON(c)}
func hashJSON(v any)string{body,_:=json.Marshal(v);sum:=sha256.Sum256(body);return "sha256:"+hex.EncodeToString(sum[:])}
func SortedKeys(values map[string]Value)[]string{keys:=make([]string,0,len(values));for key:=range values{keys=append(keys,key)};sort.Strings(keys);return keys}
