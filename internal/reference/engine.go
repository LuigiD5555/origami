package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Relation struct { Type, From, To string }
type Protocol struct {
	MaxSteps int `json:"max_steps_per_state"`
	RequiredClassifications map[string]int `json:"required_classifications,omitempty"`
}
type Experiment struct {
	ID string `json:"experiment_id"`
	Entities []string `json:"entities"`
	Domain []string `json:"domain"`
	Relations []Relation `json:"relations"`
	Protocol Protocol `json:"protocol"`
}
type State map[string]string

type Run struct {
	Initial State `json:"initial"`
	Trace []State `json:"trace"`
	Classification string `json:"classification"`
	Reason string `json:"reason,omitempty"`
}
type Evidence struct {
	Schema string `json:"schema"`
	ExperimentID string `json:"experiment_id"`
	Deterministic bool `json:"deterministic"`
	Runs []Run `json:"runs"`
	Counts map[string]int `json:"counts"`
}

func Load(path string) (Experiment, error) {
	b, err := os.ReadFile(path); if err != nil { return Experiment{}, err }
	var e Experiment
	if err := json.Unmarshal(b, &e); err != nil { return Experiment{}, err }
	if e.Protocol.MaxSteps <= 0 { return Experiment{}, fmt.Errorf("max_steps_per_state must be positive") }
	return e, nil
}

func clone(s State) State { r:=State{}; for k,v:=range s { r[k]=v }; return r }
func key(s State, entities []string) string { var b strings.Builder; for _,e:=range entities { b.WriteString(e); b.WriteByte('='); b.WriteString(s[e]); b.WriteByte(';') }; return b.String() }
func active(v string) bool { return v == "active" }

func contradiction(s State, rels []Relation) (bool,string) {
	for _,r:=range rels {
		switch r.Type {
		case "requires": if active(s[r.From]) && s[r.To] == "inactive" { return true, r.From+" requires "+r.To+" but target is inactive" }
		case "excludes": if active(s[r.From]) && active(s[r.To]) { return true, r.From+" excludes active "+r.To }
		}
	}
	return false,""
}

func step(s State, rels []Relation) State {
	n:=clone(s)
	for _,r:=range rels { if r.Type=="requires" && active(s[r.From]) && s[r.To]!="inactive" { n[r.To]="active" } }
	s=clone(n)
	for _,r:=range rels { if r.Type=="excludes" && active(s[r.From]) && s[r.To]!="active" { n[r.To]="inactive" } }
	s=clone(n)
	for _,r:=range rels { if r.Type=="couples" { if active(s[r.From]) && s[r.To]!="inactive" { n[r.To]="active" }; if s[r.From]=="inactive" && s[r.To]!="active" { n[r.To]="inactive" } } }
	s=clone(n)
	for _,r:=range rels { if r.Type=="toggles" && active(s[r.From]) { if active(s[r.To]) { n[r.To]="inactive" } else { n[r.To]="active" } } }
	return n
}

func Execute(initial State, e Experiment) Run {
	cur:=clone(initial); trace:=[]State{clone(cur)}; seen:=map[string]int{key(cur,e.Entities):0}
	for i:=0;i<e.Protocol.MaxSteps;i++ {
		if bad,why:=contradiction(cur,e.Relations); bad { return Run{Initial:clone(initial),Trace:trace,Classification:"CONTRADICTION",Reason:why} }
		next:=step(cur,e.Relations)
		if key(next,e.Entities)==key(cur,e.Entities) { trace=append(trace,clone(next)); return Run{Initial:clone(initial),Trace:trace,Classification:"FIXED_POINT"} }
		k:=key(next,e.Entities); trace=append(trace,clone(next))
		if _,ok:=seen[k]; ok { return Run{Initial:clone(initial),Trace:trace,Classification:"CYCLE"} }
		seen[k]=len(trace)-1; cur=next
	}
	return Run{Initial:clone(initial),Trace:trace,Classification:"BUDGET_EXHAUSTED"}
}

func Enumerate(e Experiment) []State {
	entities:=append([]string(nil),e.Entities...); sort.Strings(entities)
	var out []State; var rec func(int,State)
	rec=func(i int,s State){ if i==len(entities){ out=append(out,clone(s)); return }; for _,v:=range e.Domain { s[entities[i]]=v; rec(i+1,s) } }
	rec(0,State{}); return out
}

func RunExperiment(e Experiment) Evidence {
	ev:=Evidence{Schema:"origami.reference-evidence.r0",ExperimentID:e.ID,Deterministic:true,Counts:map[string]int{}}
	for _,s:=range Enumerate(e){ r:=Execute(s,e); ev.Runs=append(ev.Runs,r); ev.Counts[r.Classification]++ }
	return ev
}
