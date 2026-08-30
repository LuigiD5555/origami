package hybridcarrier

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/memoryvisual"
	"github.com/LuigiD5555/origami/internal/memory"
)

func TestVirtualMemoryCarrierAddsRedundantVisualMapAndRemainsExact(t *testing.T){
	nodes:=[]MemoryEntry{
		{Address:"ohf://demo/models/blt",CarrierID:"demo",ClusterID:"models",Kind:"concept",SemanticClass:"architecture",Label:"BLT",Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelitySummary:{Text:"Byte patch architecture."},memory.FidelityEvidence:{Text:"Measured evidence.",Verified:true,SourceRef:"src:1"}},Links:map[string][]string{"supports":{"ohf://demo/models/e1"}},Importance:memory.ImportanceProfile{Global:.9}},
		{Address:"ohf://demo/models/e1",CarrierID:"demo",ClusterID:"models",Kind:"evidence",SemanticClass:"experiment",Label:"Evidence",Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelitySummary:{Text:"Evidence summary."}},Parent:"ohf://demo/models/blt"},
	}
	envelope,err:=BuildVirtualMemoryEnvelope(receiverSpec(),"demo",nil,nodes);if err!=nil{t.Fatal(err)}
	if envelope.MemoryProfile==nil||envelope.MemoryProfile.Signature.NodeCount!=2{t.Fatalf("missing virtual-memory profile: %+v",envelope.MemoryProfile)}
	rendered,err:=RenderMemory(envelope,DefaultColumns,DefaultCellPixels,DefaultMargin,DefaultMaxPNGBytes,memoryvisual.DefaultBandHeight);if err!=nil{t.Fatal(err)}
	if rendered.MemoryBandHeight!=memoryvisual.DefaultBandHeight||rendered.VisualPlan==nil{t.Fatalf("missing visual memory map: %+v",rendered)}
	if len(rendered.VisualPlan.Points)==0{t.Fatal("visual memory plan contains no points")}
	decoded,err:=DecodeMemoryPNG(rendered.PNG,DefaultColumns,DefaultCellPixels,DefaultMargin,rendered.MemoryBandHeight);if err!=nil{t.Fatal(err)}
	if decoded.Verification.MemorySHA256!=envelope.Verification.MemorySHA256{t.Fatalf("memory identity changed through composite PNG: got=%s want=%s",decoded.Verification.MemorySHA256,envelope.Verification.MemorySHA256)}
	if decoded.MemoryProfile==nil||decoded.MemoryProfile.AddressABI!=memory.AddressABIR0{t.Fatalf("memory profile did not round-trip: %+v",decoded.MemoryProfile)}
}
