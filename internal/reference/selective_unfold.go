package reference

import (
	"fmt"
	"sort"
)

// AddressGraph indexes semantic dependencies. It indexes meaning-bearing
// entities; it is not a byte offset table, screenshot crop map, or text page map.
type AddressGraph struct {
	Dependencies map[string][]string `json:"dependencies"`
}

type SelectiveView struct {
	Requested []string                    `json:"requested"`
	Closure   []string                    `json:"closure"`
	State     CoreState                   `json:"state"`
	Folded    []string                    `json:"still_folded"`
	Rules     []Rule                      `json:"rules,omitempty"`
	Alternatives []AlternativeSet         `json:"alternatives,omitempty"`
	Commitment string                     `json:"source_commitment_sha256"`
}

func DependencyClosure(graph AddressGraph, requested []string)([]string,error){
	seen:=map[string]bool{}; visiting:=map[string]bool{}
	var visit func(string) error
	visit=func(id string)error{if visiting[id]{return fmt.Errorf("dependency cycle at %s",id)};if seen[id]{return nil};visiting[id]=true;for _,dep:=range graph.Dependencies[id]{if err:=visit(dep);err!=nil{return err}};visiting[id]=false;seen[id]=true;return nil}
	for _,id:=range requested{if id==""{return nil,fmt.Errorf("empty requested address")};if err:=visit(id);err!=nil{return nil,err}}
	out:=make([]string,0,len(seen));for id:=range seen{out=append(out,id)};sort.Strings(out);return out,nil
}

func ruleTouchesClosure(r Rule, closure map[string]bool) bool {
	for _,p:=range r.Participants{if closure[p.Entity]{return true}}
	for _,p:=range r.When{if p.Source=="state"&&closure[p.Key]{return true}}
	for _,e:=range r.Then{if closure[e.Entity]{return true}}
	return false
}

// SelectiveUnfold verifies the complete Fold commitment first, then exposes only
// the requested semantic dependency closure. It never returns unrelated state.
// Verification of a commitment is not equivalent to semantically unfolding all
// nodes: the caller receives only the closure and an explicit still-folded set.
func SelectiveUnfold(f FoldedCore,graph AddressGraph,requested []string)(SelectiveView,error){
	src,err:=Unfold(f);if err!=nil{return SelectiveView{},err}
	closure,err:=DependencyClosure(graph,requested);if err!=nil{return SelectiveView{},err}
	set:=map[string]bool{};for _,id:=range closure{set[id]=true}
	state:=CoreState{};for _,id:=range closure{if v,ok:=src.State[id];ok{state[id]=v}}
	var alternatives []AlternativeSet;for _,a:=range src.Alternatives{if set[a.ID]{alternatives=append(alternatives,a)}}
	var rules []Rule;for _,r:=range src.Rules{if ruleTouchesClosure(r,set){rules=append(rules,r)}}
	var folded []string;for id:=range src.State{if !set[id]{folded=append(folded,id)}};sort.Strings(folded)
	return SelectiveView{Requested:append([]string(nil),requested...),Closure:closure,State:state,Folded:folded,Rules:rules,Alternatives:alternatives,Commitment:f.Commitment},nil
}
