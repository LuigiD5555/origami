package hybridruntime

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/hybridcarrier"
	"github.com/LuigiD5555/origami/internal/lab/memoryvisual"
	"github.com/LuigiD5555/origami/internal/memory"
)

func TestRuntimeVirtualQueryExpandAndSignatureFromCompositePNG(t *testing.T){
	nodes:=[]hybridcarrier.MemoryEntry{
		{Address:"ohf://demo/c1/root",CarrierID:"demo",ClusterID:"c1",Kind:"concept",SemanticClass:"transformer",Label:"Transformer attention",Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelityAbstract:{Text:"Attention architecture."},memory.FidelitySummary:{Text:"Transformer uses attention to relate sequence positions."},memory.FidelityDetail:{Text:"Detailed transformer material with dependencies and evidence references."}},Links:map[string][]string{"supports":{"ohf://demo/c1/evidence"}}},
		{Address:"ohf://demo/c1/evidence",CarrierID:"demo",ClusterID:"c1",Kind:"evidence",SemanticClass:"paper",Label:"Attention evidence",Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelitySummary:{Text:"Evidence summary."},memory.FidelityEvidence:{Text:"Verified source evidence.",Verified:true,SourceRef:"paper:attention"}},Parent:"ohf://demo/c1/root"},
	}
	envelope,err:=hybridcarrier.BuildVirtualMemoryEnvelope(testReceiver(),"demo",nil,nodes);if err!=nil{t.Fatal(err)};rendered,err:=hybridcarrier.RenderMemory(envelope,0,0,0,hybridcarrier.DefaultMaxPNGBytes,memoryvisual.DefaultBandHeight);if err!=nil{t.Fatal(err)}
	runtime,err:=OpenPNG(rendered.PNG,Transport{MemoryBandHeight:memoryvisual.DefaultBandHeight});if err!=nil{t.Fatal(err)}
	sig:=runtime.Signature();if sig.CarrierID!="demo"||sig.NodeCount!=2{t.Fatalf("unexpected signature: %+v",sig)}
	packet:=runtime.MemoryQuery("What evidence supports transformer attention?",memory.ContextBudget{MaxTokens:100});if packet.Unknown||len(packet.Items)==0{t.Fatalf("unexpected context packet: %+v",packet)};if packet.Budget.UsedTokens>100{t.Fatalf("query exceeded budget: %+v",packet.Budget)}
	expanded:=runtime.Expand("ohf://demo/c1/root",memory.FidelityDetail,memory.ContextBudget{MaxTokens:100});if expanded.Unknown||len(expanded.Items)!=1{t.Fatalf("unexpected expand result: %+v",expanded)}
	verify:=runtime.Verify();if verify.Verified==nil||!*verify.Verified{t.Fatalf("unexpected verification: %+v",verify)}
}
