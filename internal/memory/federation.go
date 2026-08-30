package memory

import (
	"fmt"
	"sort"
	"strings"
)

type MountedCarrier struct {
	ID        string         `json:"id"`
	Signature GraphSignature `json:"signature"`
	Store     *Store         `json:"-"`
}

type Federation struct { carriers map[string]MountedCarrier }

func NewFederation() *Federation { return &Federation{carriers:map[string]MountedCarrier{}} }

func (f *Federation) Mount(id string, store *Store) error {
	if f == nil { return fmt.Errorf("federation is nil") }
	if strings.TrimSpace(id)=="" || store==nil { return fmt.Errorf("carrier id and store are required") }
	if _,exists:=f.carriers[id];exists{return fmt.Errorf("carrier %q already mounted",id)}
	sig:=store.Signature();sig.CarrierID=id
	f.carriers[id]=MountedCarrier{ID:id,Signature:sig,Store:store};return nil
}

func (f *Federation) Signatures() []GraphSignature {
	if f==nil{return nil};ids:=make([]string,0,len(f.carriers));for id:=range f.carriers{ids=append(ids,id)};sort.Strings(ids)
	out:=make([]GraphSignature,0,len(ids));for _,id:=range ids{out=append(out,f.carriers[id].Signature)};return out
}

type CarrierCandidate struct { ID string `json:"id"`; Score float64 `json:"score"`; Signature GraphSignature `json:"signature"` }

func (f *Federation) Rank(raw string, limit int) []CarrierCandidate {
	if f==nil{return nil};if limit<=0{limit=8};terms:=tokenize(raw);var out []CarrierCandidate
	for id,mounted:=range f.carriers{
		hay:=strings.Join(append(append(append([]string{},mounted.Signature.RootTerms...),mounted.Signature.SemanticClasses...),mounted.Signature.Clusters...)," ")
		hayTerms:=map[string]struct{}{};for _,t:=range tokenize(hay){hayTerms[t]=struct{}{}}
		score:=0.0;for _,term:=range terms{if _,ok:=hayTerms[term];ok{score+=1}}
		if score>0{out=append(out,CarrierCandidate{ID:id,Score:score,Signature:mounted.Signature})}
	}
	sort.Slice(out,func(i,j int)bool{if out[i].Score==out[j].Score{return out[i].ID<out[j].ID};return out[i].Score>out[j].Score});if len(out)>limit{out=out[:limit]};return out
}

func (f *Federation) Query(raw string, budget ContextBudget, carrierLimit int) ContextPacket {
	plan:=Analyze(raw);budget=normalizeBudget(budget);ranked:=f.Rank(raw,carrierLimit)
	if len(ranked)==0{return ContextPacket{Schema:SchemaR0+".context-packet",Query:raw,Plan:plan,Budget:budget,Unknown:true,Reason:"no carrier signature candidates"}}
	var nodes []Node
	for _,candidate:=range ranked{
		mounted:=f.carriers[candidate.ID]
		for _,node:=range mounted.Store.Nodes(){nodes=append(nodes,qualifyNode(candidate.ID,node))}
	}
	store,err:=BuildStore("federated",nodes);if err!=nil{return ContextPacket{Schema:SchemaR0+".context-packet",Query:raw,Plan:plan,Budget:budget,Unknown:true,Reason:err.Error()}}
	packet:=NewManager(store).Query(raw,budget)
	return packet
}

func qualifyNode(carrierID string, n Node) Node {
	if n.CarrierID==""{n.CarrierID=carrierID}
	mapping:=func(address string)string{if address==""||IsCanonicalAddress(address){return address};return CanonicalAddress(carrierID,firstNonEmpty(n.ClusterID,"legacy"),address)}
	oldAddress:=n.Address;n.Address=mapping(n.Address)
	// Relation targets must be qualified using the same carrier/cluster scope.
	mapTarget:=func(target string)string{if target==""||IsCanonicalAddress(target){return target};return CanonicalAddress(carrierID,firstNonEmpty(n.ClusterID,"legacy"),target)}
	if n.Parent!=""{n.Parent=mapTarget(n.Parent)}
	for i:=range n.Children{n.Children[i]=mapTarget(n.Children[i])}
	if n.Links!=nil{links:=make(map[string][]string,len(n.Links));for relation,targets:=range n.Links{mapped:=make([]string,len(targets));for i,target:=range targets{mapped[i]=mapTarget(target)};links[relation]=mapped};n.Links=links}
	for i:=range n.Edges{n.Edges[i].Target=mapTarget(n.Edges[i].Target)}
	if n.CID==""{n.CID=n.DerivedCID()}
	_ = oldAddress
	return n
}
