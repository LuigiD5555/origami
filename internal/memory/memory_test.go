package memory

import "testing"

func richNodes() []Node {
	return []Node{
		{Address:"n:blt",CarrierID:"book-a",ClusterID:"models",Kind:"concept",SemanticClass:"architecture",Label:"Byte Latent Transformer BLT",Levels:map[Fidelity]ContentLevel{
			FidelityAbstract:{Text:"BLT groups bytes into dynamic patches."},
			FidelitySummary:{Text:"BLT allocates expensive global computation over dynamically sized byte patches instead of a fixed token vocabulary."},
			FidelityDetail:{Text:"BLT uses a lightweight byte encoder and decoder around a large latent transformer. Patch boundaries can follow next-byte entropy so more computation is spent on difficult regions."},
			FidelityEvidence:{Text:"Reported experiments compare byte-patch models with tokenized baselines under controlled inference cost.",SourceRef:"paper:blt#results",Verified:true},
			FidelityExact:{Text:"EXACT-BLT-SOURCE-BLOCK-001",SourceRef:"paper:blt#exact",Verified:true},
		},Importance:ImportanceProfile{Global:.8,Local:.9,Confidence:.95},Links:map[string][]string{"supports":{"n:evidence"}},Provenance:Provenance{SourceID:"paper:blt",SourceRef:"paper:blt"}},
		{Address:"n:evidence",CarrierID:"book-a",ClusterID:"models",Kind:"evidence",SemanticClass:"experiment",Label:"BLT inference evidence",Levels:map[Fidelity]ContentLevel{FidelitySummary:{Text:"Experiments report improved inference efficiency."},FidelityEvidence:{Text:"The evaluation reports lower inference FLOPs at comparable performance.",SourceRef:"paper:blt#table",Verified:true}},Parent:"n:blt",Importance:ImportanceProfile{Confidence:.9}},
		{Address:"n:other",CarrierID:"book-a",ClusterID:"other",Kind:"concept",SemanticClass:"unrelated",Label:"Cake recipe",Levels:map[Fidelity]ContentLevel{FidelitySummary:{Text:"Unrelated cooking material."}}},
	}
}

func TestManagerUsesBoundedSelectiveUnfolding(t *testing.T) {
	store,err:=BuildStore("book-a",richNodes());if err!=nil{t.Fatal(err)}
	packet:=NewManager(store).Query("What evidence supports BLT inference efficiency?",ContextBudget{MaxTokens:120,ReserveTokens:20})
	if packet.Unknown{t.Fatalf("unexpected unknown: %+v",packet)}
	if packet.Budget.UsedTokens>100{t.Fatalf("budget exceeded: %+v",packet.Budget)}
	if len(packet.Items)==0{t.Fatal("expected at least one selected item")}
	for _,item:=range packet.Items{if item.Address=="n:other"{t.Fatal("irrelevant node should not be unfolded")}}
	if packet.Metrics.Touched>=store.Len()+10{t.Fatalf("unexpected touch count: %+v",packet.Metrics)}
}

func TestExpandDowngradesFidelityBeforeExceedingBudget(t *testing.T) {
	store,err:=BuildStore("book-a",richNodes());if err!=nil{t.Fatal(err)}
	packet:=NewManager(store).Expand("n:blt",FidelityExact,ContextBudget{MaxTokens:12})
	if packet.Unknown{t.Fatalf("expected lower fidelity to fit: %+v",packet)}
	if len(packet.Items)!=1{t.Fatalf("unexpected items: %+v",packet.Items)}
	if packet.Items[0].Fidelity==FidelityExact && packet.Items[0].TokenCost>12{t.Fatal("exact content must not exceed budget")}
	if packet.Budget.UsedTokens>12{t.Fatalf("budget exceeded: %+v",packet.Budget)}
}

func TestFederationRanksCarrierSignaturesBeforeLocalSearch(t *testing.T) {
	a,err:=BuildStore("a",[]Node{{Address:"root-a",ClusterID:"models",Kind:"concept",SemanticClass:"transformer",Label:"Transformer attention",Levels:map[Fidelity]ContentLevel{FidelitySummary:{Text:"Attention model material."}}}});if err!=nil{t.Fatal(err)}
	b,err:=BuildStore("b",[]Node{{Address:"root-b",ClusterID:"finance",Kind:"concept",SemanticClass:"accounting",Label:"Accounting ledger",Levels:map[Fidelity]ContentLevel{FidelitySummary:{Text:"Finance material."}}}});if err!=nil{t.Fatal(err)}
	fed:=NewFederation();if err:=fed.Mount("a",a);err!=nil{t.Fatal(err)};if err:=fed.Mount("b",b);err!=nil{t.Fatal(err)}
	ranked:=fed.Rank("Explain transformer attention",2);if len(ranked)!=1||ranked[0].ID!="a"{t.Fatalf("unexpected carrier ranking: %+v",ranked)}
	packet:=fed.Query("Explain transformer attention",ContextBudget{MaxTokens:100},2);if packet.Unknown||len(packet.Items)==0{t.Fatalf("unexpected federation packet: %+v",packet)}
	if packet.Items[0].CarrierID!="a"{t.Fatalf("wrong carrier exposed: %+v",packet.Items[0])}
	if !IsCanonicalAddress(packet.Items[0].Address){t.Fatalf("federated legacy address must be qualified: %s",packet.Items[0].Address)}
}

func TestWorkingSetEvictsLeastRecentlyUsedByTokenBudget(t *testing.T) {
	w:=NewWorkingSet(5)
	w.Put(ContextItem{Address:"a",TokenCost:2});w.Put(ContextItem{Address:"b",TokenCost:2});if _,ok:=w.Get("a");!ok{t.Fatal("missing cached a")};w.Put(ContextItem{Address:"c",TokenCost:2})
	if _,ok:=w.Get("b");ok{t.Fatal("b should be evicted as least recently used")}
	if _,ok:=w.Get("a");!ok{t.Fatal("a should remain")};if _,ok:=w.Get("c");!ok{t.Fatal("c should remain")}
	if w.UsedTokens()>w.MaxTokens(){t.Fatalf("working set exceeded budget: %d",w.UsedTokens())}
}

func TestLegacyValueRemainsQueryableWithoutRichMetadata(t *testing.T) {
	store,err:=BuildStore("legacy",[]Node{{Address:"mem:1",Value:"AMBER-10593"}});if err!=nil{t.Fatal(err)}
	packet:=NewManager(store).Query("AMBER-10593",ContextBudget{MaxTokens:20});if packet.Unknown||len(packet.Items)!=1{t.Fatalf("legacy memory should remain usable: %+v",packet)}
	if packet.Items[0].Fidelity!=FidelitySummary{t.Fatalf("legacy query must default to summary semantics, got %s",packet.Items[0].Fidelity)}
}
