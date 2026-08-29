package reference

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type SemanticStatus string

const (
	StatusPresent SemanticStatus = "PRESENT"
	StatusAbsent SemanticStatus = "ABSENT"
	StatusUnknown SemanticStatus = "UNKNOWN"
	StatusInhibited SemanticStatus = "INHIBITED"
	StatusCancelled SemanticStatus = "CANCELLED"
)

type SemanticValue struct { Status SemanticStatus `json:"status"`; Value string `json:"value,omitempty"`; Reason string `json:"reason,omitempty"` }
type CoreState map[string]SemanticValue
type Context struct { ID string `json:"id"`; Values map[string]string `json:"values,omitempty"` }
type Predicate struct { Source string `json:"source"`; Key string `json:"key"`; Equals string `json:"equals"` }
type Effect struct { Entity string `json:"entity"`; Status SemanticStatus `json:"status"`; Value string `json:"value,omitempty"`; Reason string `json:"reason,omitempty"` }

// Participant requires an entity to participate with an exact semantic status
// and, optionally, an exact value. A rule with 3+ participants is genuinely
// higher-order: the complete participating set is a condition of the rule.
type Participant struct { Entity string `json:"entity"`; Status SemanticStatus `json:"status"`; Value string `json:"value,omitempty"` }

type Rule struct {
	ID string `json:"id"`
	When []Predicate `json:"when,omitempty"`
	Participants []Participant `json:"participants,omitempty"`
	Then []Effect `json:"then"`
}

type InteractionEvidence struct { RuleID string `json:"rule_id"`; Participants []string `json:"participants"` }
type TransitionEvidence struct { Step int `json:"step"`; ContextID string `json:"context_id"`; FiredRules []string `json:"fired_rules"`; Interactions []InteractionEvidence `json:"interactions,omitempty"` }
type CoreRun struct { Initial CoreState `json:"initial"`; Trace []CoreState `json:"trace"`; Transitions []TransitionEvidence `json:"transitions"` }
type ObservationContract struct { ID string `json:"id"`; Authority string `json:"authority"`; Entity string `json:"entity"` }
type ObservationResult struct { ContractID string `json:"contract_id"`; Authority string `json:"authority"`; Entity string `json:"entity"`; Value SemanticValue `json:"value"`; TraceHash string `json:"trace_sha256"` }

func cloneCoreState(s CoreState) CoreState { r:=CoreState{}; for k,v:=range s{r[k]=v}; return r }
func predicateMatches(p Predicate,state CoreState,context Context) bool { switch p.Source { case "context": return context.Values[p.Key]==p.Equals; case "state": v,ok:=state[p.Key]; return ok&&(string(v.Status)==p.Equals||v.Value==p.Equals); default:return false } }
func participantMatches(p Participant,state CoreState) bool { v,ok:=state[p.Entity]; if !ok||v.Status!=p.Status{return false}; return p.Value==""||v.Value==p.Value }
func ruleMatches(r Rule,state CoreState,context Context) bool { for _,p:=range r.When{if !predicateMatches(p,state,context){return false}}; for _,p:=range r.Participants{if !participantMatches(p,state){return false}}; return true }

func participantIDs(r Rule) []string { ids:=make([]string,0,len(r.Participants)); for _,p:=range r.Participants{ids=append(ids,p.Entity)}; sort.Strings(ids); return ids }
func validateParticipants(r Rule) error { seen:=map[string]bool{}; for _,p:=range r.Participants{if p.Entity==""{return fmt.Errorf("rule %s has empty participant",r.ID)}; if seen[p.Entity]{return fmt.Errorf("rule %s repeats participant %s",r.ID,p.Entity)}; seen[p.Entity]=true}; return nil }

func ApplyRulesWithEvidence(state CoreState,context Context,rules []Rule)(CoreState,[]string,[]InteractionEvidence,error){
	next:=cloneCoreState(state); writes:=map[string]SemanticValue{}; var fired []string; var interactions []InteractionEvidence
	for _,r:=range rules{
		if r.ID==""{return nil,nil,nil,fmt.Errorf("rule id must not be empty")}; if err:=validateParticipants(r);err!=nil{return nil,nil,nil,err}; if !ruleMatches(r,state,context){continue}
		fired=append(fired,r.ID); if len(r.Participants)>0{interactions=append(interactions,InteractionEvidence{RuleID:r.ID,Participants:participantIDs(r)})}
		for _,e:=range r.Then{v:=SemanticValue{Status:e.Status,Value:e.Value,Reason:e.Reason}; if prior,ok:=writes[e.Entity];ok&&prior!=v{return nil,nil,nil,fmt.Errorf("conflicting effects for %s",e.Entity)}; writes[e.Entity]=v}
	}
	for entity,value:=range writes{next[entity]=value}; return next,fired,interactions,nil
}
func ApplyRules(state CoreState,context Context,rules []Rule)(CoreState,[]string,error){next,fired,_,err:=ApplyRulesWithEvidence(state,context,rules); return next,fired,err}
func ExecuteCore(initial CoreState,contexts []Context,rules []Rule)(CoreRun,error){cur:=cloneCoreState(initial); run:=CoreRun{Initial:cloneCoreState(initial),Trace:[]CoreState{cloneCoreState(initial)}}; for i,context:=range contexts{next,fired,interactions,err:=ApplyRulesWithEvidence(cur,context,rules);if err!=nil{return CoreRun{},err};run.Transitions=append(run.Transitions,TransitionEvidence{Step:i+1,ContextID:context.ID,FiredRules:fired,Interactions:interactions});run.Trace=append(run.Trace,cloneCoreState(next));cur=next};return run,nil}
func TraceHash(trace []CoreState)(string,error){b,err:=json.Marshal(trace);if err!=nil{return "",err};h:=sha256.Sum256(b);return hex.EncodeToString(h[:]),nil}
func Observe(run CoreRun,contract ObservationContract)(ObservationResult,error){if contract.ID==""||contract.Authority==""{return ObservationResult{},fmt.Errorf("observation contract requires id and authority")};if len(run.Trace)==0{return ObservationResult{},fmt.Errorf("run has no trace")};final:=run.Trace[len(run.Trace)-1];v,ok:=final[contract.Entity];if !ok{v=SemanticValue{Status:StatusAbsent}};h,err:=TraceHash(run.Trace);if err!=nil{return ObservationResult{},err};return ObservationResult{ContractID:contract.ID,Authority:contract.Authority,Entity:contract.Entity,Value:v,TraceHash:h},nil}

// ParticipatingSetKey is canonical and order-independent. It is useful for
// evidence, addressing and later Fold/Unfold dependency closures.
func ParticipatingSetKey(ids []string) string { cp:=append([]string(nil),ids...);sort.Strings(cp);return strings.Join(cp,"+") }
