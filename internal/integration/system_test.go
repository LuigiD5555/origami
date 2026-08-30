package integration

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/core"
	"github.com/LuigiD5555/origami/internal/document"
	"github.com/LuigiD5555/origami/internal/evidence"
	"github.com/LuigiD5555/origami/internal/memory"
)

func TestEndToEndExactSourceToCanonicalSemanticMemory(t *testing.T) {
	doc,err:=document.FromText("paper","application/pdf",[]byte("Byte models can use entropy-driven dynamic patches.\n\nThe reported experiment lowers inference cost."));if err!=nil{t.Fatal(err)}
	bundle,err:=document.BuildBundle(doc);if err!=nil{t.Fatal(err)}
	system,err:=NewSystem("paper",bundle.Nodes);if err!=nil{t.Fatal(err)}

	packet:=system.Query("entropy dynamic patches",memory.ContextBudget{MaxTokens:120,ReserveTokens:20})
	if packet.Unknown||len(packet.Items)==0{t.Fatalf("bounded source query failed: %+v",packet)}
	if packet.Budget.UsedTokens>100{t.Fatalf("active context budget exceeded: %+v",packet.Budget)}
	var evidenceAddress string
	for _,item:=range packet.Items{
		if item.Kind=="source-block"{
			if item.Fidelity==memory.FidelitySummary{t.Fatal("raw source was falsely promoted to semantic summary")}
			evidenceAddress=item.Address;break
		}
	}
	if evidenceAddress==""{t.Fatalf("query did not return addressable source block: %+v",packet.Items)}

	reduction,err:=system.Reduce([]evidence.Proposal{{ID:"reader-1",ClaimKey:"model.dynamic-patches",Proposition:"Byte models can use entropy-driven dynamic patches.",Position:evidence.PositionSupport,Evidence:[]evidence.EvidenceRef{{Address:evidenceAddress,Fidelity:memory.FidelityEvidence}}}},evidence.DefaultPolicy());if err!=nil{t.Fatal(err)}
	if reduction.Claims[0].Status!=evidence.StatusVerified{t.Fatalf("evidence-backed claim not verified: %+v",reduction.Claims[0])}
	if reduction.Claims[0].Exactness!=evidence.ExactnessNotClaimed{t.Fatal("semantic interpretation must not silently become exact")}
	promoted,err:=system.Promote(reduction);if err!=nil{t.Fatal(err)}
	if len(promoted)!=1||promoted[0].Kind!="canonical-claim"{t.Fatalf("canonical semantic promotion failed: %+v",promoted)}

	semanticPacket:=system.Query("dynamic patches",memory.ContextBudget{MaxTokens:80})
	foundClaim:=false
	for _,item:=range semanticPacket.Items{if item.Kind=="canonical-claim"&&strings.Contains(item.Content,"dynamic patches"){foundClaim=true}}
	if !foundClaim{t.Fatalf("later query did not reuse canonical semantic memory: %+v",semanticPacket.Items)}
}

func TestCoreSemanticGraphAndTrajectoryAreDirectlyQueryableWithoutLLMReconstruction(t *testing.T) {
	graph:=core.SemanticGraph{Nodes:map[string]core.SemanticNode{
		"cause":{Address:"cause",Kind:"concept",Label:"Cause",Value:core.Value{Status:core.StatusPresent,Data:"high entropy"},Verified:true},
		"effect":{Address:"effect",Kind:"claim",Label:"Selective compute",Value:core.Value{Status:core.StatusPresent,Data:"more computation on difficult regions"},Dependencies:[]string{"cause"},Verified:true},
	}}
	graphNodes,err:=SemanticGraphNodes("semantic-demo",graph);if err!=nil{t.Fatal(err)}
	machine:=core.NewMachine([]core.Rule{{ID:"derive",When:[]core.Predicate{{Scope:"state",Key:"entropy",Op:"present"}},Then:[]core.Effect{{Key:"compute",Mode:"set",Status:core.StatusPresent,Data:"selective"}}}},2)
	trajectory,err:=core.Execute(core.State{Values:map[string]core.Value{"entropy":{Status:core.StatusPresent,Data:"high"}}},machine,[]core.Context{{ID:"document-context",Values:map[string]core.Value{}}});if err!=nil{t.Fatal(err)}
	trajectoryNodes,err:=TrajectoryNodes("semantic-demo","run-1",trajectory);if err!=nil{t.Fatal(err)}
	all:=append(graphNodes,trajectoryNodes...)
	system,err:=NewSystem("semantic-demo",all);if err!=nil{t.Fatal(err)}
	packet:=system.Query("difficult regions selective compute",memory.ContextBudget{MaxTokens:100})
	if packet.Unknown||len(packet.Items)==0{t.Fatalf("core-derived semantics not queryable: %+v",packet)}
	for _,item:=range packet.Items{if item.SourceRef==""&&item.Kind=="transition-evidence"{t.Fatalf("machine-derived transition lost provenance: %+v",item)}}
}

func TestExactClaimPromotionRequiresExactByteEquality(t *testing.T) {
	doc,err:=document.FromText("literal","text/plain",[]byte("EXACT SENTENCE"));if err!=nil{t.Fatal(err)}
	bundle,err:=document.BuildBundle(doc);if err!=nil{t.Fatal(err)}
	system,err:=NewSystem("literal",bundle.Nodes);if err!=nil{t.Fatal(err)}
	block:="";for _,node:=range bundle.Nodes{if node.Kind=="source-block"{block=node.Address}}
	paraphrase,err:=system.Reduce([]evidence.Proposal{{ID:"p",ClaimKey:"literal",Proposition:"Exact sentence",Position:evidence.PositionSupport,Evidence:[]evidence.EvidenceRef{{Address:block,Fidelity:memory.FidelityExact}}}},evidence.DefaultPolicy());if err!=nil{t.Fatal(err)}
	if paraphrase.Claims[0].Exactness!=evidence.ExactnessNotClaimed{t.Fatalf("case/paraphrase drift falsely exact: %+v",paraphrase.Claims[0])}
	exact,err:=system.Reduce([]evidence.Proposal{{ID:"e",ClaimKey:"literal",Proposition:"EXACT SENTENCE",Position:evidence.PositionSupport,Evidence:[]evidence.EvidenceRef{{Address:block,Fidelity:memory.FidelityExact}}}},evidence.DefaultPolicy());if err!=nil{t.Fatal(err)}
	if exact.Claims[0].Exactness!=evidence.ExactnessVerified{t.Fatalf("byte-equal exact claim not recognized: %+v",exact.Claims[0])}
}
