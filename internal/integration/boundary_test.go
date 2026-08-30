package integration

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/document"
	"github.com/LuigiD5555/origami/internal/evidence"
	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/lab/observation"
	"github.com/LuigiD5555/origami/internal/memory"
)

func TestPerceptionObservationCannotDirectlyBecomeVerifiedSemanticEvidence(t *testing.T){
	envelope:=observation.Envelope{Schema:observation.SchemaV1,BundleID:"bundle-1",Model:"vlm",Trial:1,Status:observation.Known,Observations:map[string]any{"region":"looks like a graph"}}
	node,err:=PerceptionNode("carrier",envelope);if err!=nil{t.Fatal(err)}
	level,ok:=node.Level(memory.FidelityEvidence);if !ok{t.Fatal("perception evidence level missing")};if level.Verified{t.Fatal("well-formed perception must not automatically become semantic authority")}
	system,err:=NewSystem("carrier",[]memory.Node{node});if err!=nil{t.Fatal(err)}
	result,err:=system.Reduce([]evidence.Proposal{{ID:"vision",ClaimKey:"graph.exists",Proposition:"a graph exists",Position:evidence.PositionSupport,Evidence:[]evidence.EvidenceRef{{Address:node.Address,Fidelity:memory.FidelityEvidence}}}},evidence.DefaultPolicy());if err!=nil{t.Fatal(err)}
	if result.Claims[0].Status!=evidence.StatusUnknown{t.Fatalf("perception bypassed verification boundary: %+v",result.Claims[0])}
}

func TestFixedCarrierBindsExternalCanonicalMemoryInsteadOfContainingCorpus(t *testing.T){
	doc,err:=document.FromText("paper","application/pdf",[]byte("one block of source truth"));if err!=nil{t.Fatal(err)}
	bundle,err:=document.BuildBundle(doc);if err!=nil{t.Fatal(err)}
	system,err:=NewSystem("paper",bundle.Nodes);if err!=nil{t.Fatal(err)}
	meta,err:=FixedCarrierMetadata(system,doc.SourceHash,CarrierCounts{Pages:1,Blocks:1,Documents:1});if err!=nil{t.Fatal(err)}
	if meta.StoreRoot!=CanonicalStoreRoot(system.Nodes()){t.Fatalf("Fixed Carrier root does not bind canonical memory: %s",meta.StoreRoot)}
	png,decoded,err:=fixedcarrier.Render(meta);if err!=nil{t.Fatal(err)}
	if len(png)!=fixedcarrier.FixedPNGBytes{t.Fatalf("carrier grew with source: %d",len(png))}
	if decoded.StoreRoot!=meta.StoreRoot{t.Fatalf("decoded control plane lost canonical root: %+v",decoded)}
	if string(meta.GraphSignature)=="one block of source truth"{t.Fatal("fixed carrier graph signature became a corpus dump")}
}
