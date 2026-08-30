package integration

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/LuigiD5555/origami/internal/core"
	"github.com/LuigiD5555/origami/internal/memory"
)

func SemanticGraphNodes(carrierID string, graph core.SemanticGraph) ([]memory.Node, error) {
	folded, store, err := core.FoldGraph(graph)
	if err != nil { return nil, err }
	if !core.VerifyFold(folded, store) { return nil, fmt.Errorf("semantic graph fold verification failed") }
	addresses := make([]string, 0, len(graph.Nodes)); for address := range graph.Nodes { addresses = append(addresses, address) }; sort.Strings(addresses)
	out := make([]memory.Node, 0, len(addresses))
	for _, address := range addresses {
		node := graph.Nodes[address]
		vmAddress := memory.CanonicalAddress(carrierID, "semantic", cleanID(address))
		edges := make([]memory.Edge, 0)
		for _, dependency := range node.Dependencies { edges = append(edges, memory.Edge{Relation:"depends",Target:memory.CanonicalAddress(carrierID,"semantic",cleanID(dependency)),Weight:1,Confidence:1,ProofRef:folded.Commitment}) }
		for _, alternative := range node.Alternatives { edges = append(edges, memory.Edge{Relation:"alternative",Target:memory.CanonicalAddress(carrierID,"semantic",cleanID(alternative)),Weight:1,Confidence:1,ProofRef:folded.Commitment}) }
		for relation, targets := range node.Relations { for _, target := range targets { edges = append(edges,memory.Edge{Relation:relation,Target:memory.CanonicalAddress(carrierID,"semantic",cleanID(target)),Weight:1,Confidence:1,ProofRef:folded.Commitment}) } }
		sort.Slice(edges,func(i,j int)bool{if edges[i].Relation==edges[j].Relation{return edges[i].Target<edges[j].Target};return edges[i].Relation<edges[j].Relation})
		evidenceBody,_:=json.Marshal(struct{Address string `json:"address"`;Status core.Status `json:"status"`;Data string `json:"data,omitempty"`;FoldCommitment string `json:"fold_commitment"`;EvidenceRefs []string `json:"evidence_refs,omitempty"`}{node.Address,node.Value.Status,node.Value.Data,folded.Commitment,node.EvidenceRefs})
		levels:=map[memory.Fidelity]memory.ContentLevel{
			memory.FidelityAbstract:{Text:fmt.Sprintf("semantic %s %s",node.Kind,node.Value.Status)},
			memory.FidelitySummary:{Text:semanticSummary(node)},
			memory.FidelityEvidence:{Text:string(evidenceBody),SourceRef:folded.Commitment,Verified:node.Verified},
		}
		out=append(out,memory.Node{Address:vmAddress,CarrierID:carrierID,ClusterID:"semantic",Kind:firstNonEmptyString(node.Kind,"semantic-node"),SemanticClass:string(node.Value.Status),Label:firstNonEmptyString(node.Label,node.Address),Levels:levels,Edges:edges,Provenance:memory.Provenance{SourceID:"origami:core",SourceRef:node.Address,SourceHash:folded.Commitment},VerificationRef:core.SemanticNodeCID(node),Importance:memory.ImportanceProfile{Local:1,Confidence:boolConfidence(node.Verified)}})
	}
	return out,nil
}

func TrajectoryNodes(carrierID, trajectoryID string, trajectory core.Trajectory) ([]memory.Node,error) {
	if trajectory.TraceHash=="" { return nil,fmt.Errorf("trajectory hash required") }
	out:=make([]memory.Node,0,len(trajectory.States)+len(trajectory.Evidence))
	stateAddresses:=make([]string,len(trajectory.States))
	for i,state:=range trajectory.States { stateAddresses[i]=memory.CanonicalAddress(carrierID,"trajectory",fmt.Sprintf("%s-state-%06d",cleanID(trajectoryID),i)) }
	for i,state:=range trajectory.States {
		body,_:=json.Marshal(state);edges:=[]memory.Edge{}
		if i>0 { edges=append(edges,memory.Edge{Relation:"previous-state",Target:stateAddresses[i-1],Weight:1,Confidence:1,ProofRef:trajectory.TraceHash}) }
		if i+1<len(stateAddresses) { edges=append(edges,memory.Edge{Relation:"next-state",Target:stateAddresses[i+1],Weight:1,Confidence:1,ProofRef:trajectory.TraceHash}) }
		out=append(out,memory.Node{Address:stateAddresses[i],CarrierID:carrierID,ClusterID:"trajectory",Kind:"semantic-state",SemanticClass:"logical-time",Label:fmt.Sprintf("%s state %d",trajectoryID,i),Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelityAbstract:{Text:fmt.Sprintf("logical state %d",i)},memory.FidelitySummary:{Text:stateSummary(state)},memory.FidelityEvidence:{Text:string(body),SourceRef:trajectory.TraceHash,Verified:true}},Edges:edges,Provenance:memory.Provenance{SourceID:"origami:core-trajectory",SourceRef:trajectoryID,SourceHash:trajectory.TraceHash},VerificationRef:core.StateHash(state)})
	}
	for _,ev:=range trajectory.Evidence {
		body,_:=json.Marshal(ev);address:=memory.CanonicalAddress(carrierID,"trajectory",fmt.Sprintf("%s-transition-%06d",cleanID(trajectoryID),ev.Step));edges:=[]memory.Edge{}
		if ev.Step-1>=0&&ev.Step-1<len(stateAddresses){edges=append(edges,memory.Edge{Relation:"from-state",Target:stateAddresses[ev.Step-1],Weight:1,Confidence:1,ProofRef:trajectory.TraceHash})}
		if ev.Step>=0&&ev.Step<len(stateAddresses){edges=append(edges,memory.Edge{Relation:"to-state",Target:stateAddresses[ev.Step],Weight:1,Confidence:1,ProofRef:trajectory.TraceHash})}
		out=append(out,memory.Node{Address:address,CarrierID:carrierID,ClusterID:"trajectory",Kind:"transition-evidence",SemanticClass:"rule-provenance",Label:fmt.Sprintf("%s transition %d",trajectoryID,ev.Step),Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelityAbstract:{Text:fmt.Sprintf("transition %d context %s",ev.Step,ev.ContextID)},memory.FidelityEvidence:{Text:string(body),SourceRef:trajectory.TraceHash,Verified:true}},Edges:edges,Provenance:memory.Provenance{SourceID:"origami:core-trajectory",SourceRef:trajectoryID,SourceHash:trajectory.TraceHash},VerificationRef:ev.StateAfterHash})
	}
	sort.Slice(out,func(i,j int)bool{return out[i].Address<out[j].Address});return out,nil
}

func semanticSummary(node core.SemanticNode) string {
	if node.Value.Status==core.StatusPresent&&node.Value.Data!="" { return node.Value.Data }
	if node.Value.Data!="" { return fmt.Sprintf("%s: %s",node.Value.Status,node.Value.Data) }
	return string(node.Value.Status)
}
func stateSummary(state core.State) string { keys:=core.SortedKeys(state.Values);out:="";for _,key:=range keys{value:=state.Values[key];if out!=""{out+="; "};out+=fmt.Sprintf("%s=%s:%s",key,value.Status,value.Data)};return out }
func firstNonEmptyString(values ...string) string { for _,value:=range values{if value!=""{return value}};return "" }
func boolConfidence(value bool) float64 { if value{return 1};return 0 }
