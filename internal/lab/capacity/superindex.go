package capacity

import (
	"fmt"
	"sort"
)

type Node struct {
	Address string
	Value string
	Links map[string][]string
}

type Metrics struct {
	Requested int `json:"requested"`
	Touched int `json:"touched"`
	Unfolded int `json:"unfolded"`
	Exposed int `json:"exposed"`
	ActiveBytes int `json:"active_bytes"`
	TokenEquivalent int `json:"token_equivalent"`
	PhysicalDepth int `json:"physical_depth"`
	CognitiveDepth int `json:"cognitive_depth"`
}

type Result struct { Nodes []Node; Metrics Metrics }

type SuperIndex struct { byAddress map[string]Node }

func Build(nodes []Node) *SuperIndex {
	i:=&SuperIndex{byAddress:make(map[string]Node,len(nodes))}
	for _,n:=range nodes { i.byAddress[n.Address]=n }
	return i
}

func (s *SuperIndex) Lookup(address string)(Node,bool){ n,ok:=s.byAddress[address]; return n,ok }

// Follow exposes only the requested relation path. Missing addresses/relations
// stop as UNKNOWN rather than triggering a global scan.
func (s *SuperIndex) Follow(start, relation string, depth, budget int) (Result,error) {
	m:=Metrics{Requested:1,PhysicalDepth:depth}; if depth<0{return Result{},fmt.Errorf("depth must be >= 0")}
	addr:=start; out:=make([]Node,0,depth+1)
	for step:=0;step<=depth;step++ {
		n,ok:=s.byAddress[addr]; m.Touched++
		if !ok { return Result{Nodes:out,Metrics:m},nil }
		out=append(out,n); m.Unfolded++; m.Exposed++
		m.ActiveBytes += len(n.Address)+len(n.Value); m.TokenEquivalent=(m.ActiveBytes+3)/4
		m.CognitiveDepth=step
		if budget>0 && m.TokenEquivalent>budget{return Result{Nodes:out,Metrics:m},fmt.Errorf("working window exceeded: %d > %d",m.TokenEquivalent,budget)}
		if step==depth { break }
		next:=append([]string(nil),n.Links[relation]...); sort.Strings(next)
		if len(next)==0 { return Result{Nodes:out,Metrics:m},nil }
		addr=next[0]
	}
	return Result{Nodes:out,Metrics:m},nil
}
