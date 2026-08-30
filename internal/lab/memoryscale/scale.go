package memoryscale

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/LuigiD5555/origami/internal/memory"
)

const SchemaR0 = "origami.memory-scale-lab.r0"

type Config struct {
	CarrierCounts           []int  `json:"carrier_counts"`
	NodesPerCarrier         int    `json:"nodes_per_carrier"`
	PositiveQueriesPerScale int    `json:"positive_queries_per_scale"`
	NegativeQueriesPerScale int    `json:"negative_queries_per_scale"`
	BudgetTokens            int    `json:"budget_tokens"`
	ReserveTokens           int    `json:"reserve_tokens"`
	CarrierLimit            int    `json:"carrier_limit"`
	Seed                    uint64 `json:"seed"`
}

func DefaultConfig() Config {
	return Config{
		CarrierCounts:           []int{1, 10, 100, 1000},
		NodesPerCarrier:         64,
		PositiveQueriesPerScale: 32,
		NegativeQueriesPerScale: 4,
		BudgetTokens:            4000,
		ReserveTokens:           0,
		CarrierLimit:            4,
		Seed:                    20260829,
	}
}

func (c Config) Validate() error {
	if len(c.CarrierCounts) == 0 { return fmt.Errorf("at least one carrier count is required") }
	last := 0
	for _, n := range c.CarrierCounts {
		if n <= 0 { return fmt.Errorf("carrier counts must be positive") }
		if n <= last { return fmt.Errorf("carrier counts must be strictly increasing") }
		if n > 100000 { return fmt.Errorf("carrier count %d exceeds laboratory safety limit 100000", n) }
		last = n
	}
	if c.NodesPerCarrier < 3 { return fmt.Errorf("nodes_per_carrier must be >= 3") }
	if c.NodesPerCarrier > 10000 { return fmt.Errorf("nodes_per_carrier exceeds laboratory safety limit 10000") }
	if c.PositiveQueriesPerScale <= 0 { return fmt.Errorf("positive_queries_per_scale must be positive") }
	if c.NegativeQueriesPerScale < 0 { return fmt.Errorf("negative_queries_per_scale cannot be negative") }
	if c.BudgetTokens <= 0 { return fmt.Errorf("budget_tokens must be positive") }
	if c.ReserveTokens < 0 || c.ReserveTokens >= c.BudgetTokens { return fmt.Errorf("reserve_tokens must be >=0 and < budget_tokens") }
	if c.CarrierLimit <= 0 { return fmt.Errorf("carrier_limit must be positive") }
	return nil
}

type Report struct {
	Schema      string        `json:"schema"`
	Config      Config        `json:"config"`
	Scales      []ScaleResult `json:"scales"`
	Definitions Definitions   `json:"definitions"`
}

type Definitions struct {
	UsefulContextEfficiency   string `json:"useful_context_efficiency"`
	KnowledgeScaleDegradation string `json:"knowledge_scale_degradation"`
	RoutingSelectivity        string `json:"routing_selectivity"`
	ExposureFraction          string `json:"exposure_fraction"`
}

type ScaleResult struct {
	CarrierCount              int     `json:"carrier_count"`
	NodesPerCarrier           int     `json:"nodes_per_carrier"`
	TotalNodes                int     `json:"total_nodes"`
	ActiveBudgetTokens        int     `json:"active_budget_tokens"`
	PositiveQueries           int     `json:"positive_queries"`
	NegativeQueries           int     `json:"negative_queries"`
	CarrierTop1Accuracy       float64 `json:"carrier_top1_accuracy"`
	CarrierRecallAtK          float64 `json:"carrier_recall_at_k"`
	TargetHitAccuracy         float64 `json:"target_hit_accuracy"`
	EvidenceHitAccuracy       float64 `json:"evidence_hit_accuracy"`
	NegativeUnknownAccuracy  float64 `json:"negative_unknown_accuracy"`
	BudgetViolations          int     `json:"budget_violations"`
	FalseExact                int     `json:"false_exact"`
	AvgSignaturesScanned      float64 `json:"avg_signatures_scanned"`
	AvgCandidateCarriers      float64 `json:"avg_candidate_carriers"`
	AvgSelectedCarriers       float64 `json:"avg_selected_carriers"`
	AvgMetadataNodesLoaded    float64 `json:"avg_metadata_nodes_loaded"`
	AvgNodesTouched           float64 `json:"avg_nodes_touched"`
	AvgGraphHops              float64 `json:"avg_graph_hops"`
	AvgNodesUnfolded          float64 `json:"avg_nodes_unfolded"`
	AvgTokensExposed          float64 `json:"avg_tokens_exposed"`
	AvgIrrelevantTokens       float64 `json:"avg_irrelevant_tokens"`
	UsefulContextEfficiency   float64 `json:"useful_context_efficiency"`
	KnowledgeScaleDegradation float64 `json:"knowledge_scale_degradation"`
	RoutingSelectivity        float64 `json:"routing_selectivity"`
	ExposureFraction          float64 `json:"exposure_fraction"`
	NodesPerActiveToken       float64 `json:"nodes_per_active_token"`
}

type TraceRecord struct {
	Schema           string                    `json:"schema"`
	CarrierCount     int                       `json:"carrier_count"`
	QueryIndex       int                       `json:"query_index"`
	Kind             string                    `json:"kind"`
	Query            string                    `json:"query"`
	TargetCarrier    string                    `json:"target_carrier,omitempty"`
	Ranked           []memory.CarrierCandidate `json:"ranked,omitempty"`
	Federation       memory.FederationMetrics  `json:"federation"`
	PacketMetrics    memory.ContextMetrics     `json:"packet_metrics"`
	Budget           memory.ContextBudget      `json:"budget"`
	Top1Correct      bool                      `json:"top1_correct,omitempty"`
	RecallAtK        bool                      `json:"recall_at_k,omitempty"`
	TargetHit        bool                      `json:"target_hit,omitempty"`
	EvidenceHit      bool                      `json:"evidence_hit,omitempty"`
	UsefulTokens     int                       `json:"useful_tokens"`
	IrrelevantTokens int                       `json:"irrelevant_tokens"`
	FalseExact       int                       `json:"false_exact"`
	Unknown          bool                      `json:"unknown"`
	Reason           string                    `json:"reason,omitempty"`
}

type accum struct {
	top1, recall, targetHit, evidenceHit     int
	negativeUnknown                         int
	budgetViolations, falseExact            int
	signatures, candidates, selected         int
	metadataLoaded, touched, hops, unfolded int
	exposed, useful, irrelevant              int
}

func Run(cfg Config, trace io.Writer) (Report, error) {
	if err := cfg.Validate(); err != nil { return Report{}, err }
	report := Report{
		Schema: SchemaR0,
		Config: cfg,
		Definitions: Definitions{
			UsefulContextEfficiency:   "useful target-carrier tokens / all exposed tokens; higher is better",
			KnowledgeScaleDegradation: "max(0, 1 - target_hit_accuracy_at_scale / baseline_target_hit_accuracy); lower is better",
			RoutingSelectivity:        "average selected carriers / total carriers; lower is better when accuracy is preserved",
			ExposureFraction:          "average unfolded nodes / total available nodes; lower is better when accuracy is preserved",
		},
	}
	var encoder *json.Encoder
	if trace != nil { encoder = json.NewEncoder(trace) }
	baselineAccuracy := 0.0
	for scaleIndex, carrierCount := range cfg.CarrierCounts {
		federation, err := buildFederation(carrierCount, cfg.NodesPerCarrier)
		if err != nil { return Report{}, err }
		a := accum{}
		for q := 0; q < cfg.PositiveQueriesPerScale; q++ {
			target := targetIndex(cfg.Seed, scaleIndex, q, carrierCount)
			targetCarrier := carrierID(target)
			query := fmt.Sprintf("What evidence supports systems %s?", topicID(target))
			result := federation.QueryDetailed(query, memory.ContextBudget{MaxTokens:cfg.BudgetTokens,ReserveTokens:cfg.ReserveTokens}, cfg.CarrierLimit)
			tr := scorePositive(carrierCount, q, query, targetCarrier, result)
			if tr.Top1Correct { a.top1++ }
			if tr.RecallAtK { a.recall++ }
			if tr.TargetHit { a.targetHit++ }
			if tr.EvidenceHit { a.evidenceHit++ }
			if result.Packet.Budget.UsedTokens > cfg.BudgetTokens-cfg.ReserveTokens { a.budgetViolations++ }
			a.falseExact += tr.FalseExact
			addWork(&a, result, tr.UsefulTokens, tr.IrrelevantTokens)
			if encoder != nil { if err := encoder.Encode(tr); err != nil { return Report{}, err } }
		}
		for q := 0; q < cfg.NegativeQueriesPerScale; q++ {
			query := fmt.Sprintf("unseen-needle-%d-%d", scaleIndex, q)
			result := federation.QueryDetailed(query, memory.ContextBudget{MaxTokens:cfg.BudgetTokens,ReserveTokens:cfg.ReserveTokens}, cfg.CarrierLimit)
			tr := scoreNegative(carrierCount, q, query, result)
			if tr.Unknown && result.Packet.Budget.UsedTokens == 0 { a.negativeUnknown++ }
			if result.Packet.Budget.UsedTokens > cfg.BudgetTokens-cfg.ReserveTokens { a.budgetViolations++ }
			a.falseExact += tr.FalseExact
			addWork(&a, result, 0, tr.IrrelevantTokens)
			if encoder != nil { if err := encoder.Encode(tr); err != nil { return Report{}, err } }
		}
		positive := float64(cfg.PositiveQueriesPerScale)
		all := float64(cfg.PositiveQueriesPerScale + cfg.NegativeQueriesPerScale)
		totalNodes := carrierCount * cfg.NodesPerCarrier
		result := ScaleResult{
			CarrierCount: carrierCount, NodesPerCarrier: cfg.NodesPerCarrier, TotalNodes: totalNodes,
			ActiveBudgetTokens: cfg.BudgetTokens-cfg.ReserveTokens,
			PositiveQueries: cfg.PositiveQueriesPerScale, NegativeQueries: cfg.NegativeQueriesPerScale,
			CarrierTop1Accuracy: ratio(a.top1, positive), CarrierRecallAtK: ratio(a.recall, positive),
			TargetHitAccuracy: ratio(a.targetHit, positive), EvidenceHitAccuracy: ratio(a.evidenceHit, positive),
			NegativeUnknownAccuracy: ratio(a.negativeUnknown, float64(cfg.NegativeQueriesPerScale)),
			BudgetViolations:a.budgetViolations, FalseExact:a.falseExact,
			AvgSignaturesScanned:ratio(a.signatures,all), AvgCandidateCarriers:ratio(a.candidates,all), AvgSelectedCarriers:ratio(a.selected,all),
			AvgMetadataNodesLoaded:ratio(a.metadataLoaded,all), AvgNodesTouched:ratio(a.touched,all), AvgGraphHops:ratio(a.hops,all), AvgNodesUnfolded:ratio(a.unfolded,all),
			AvgTokensExposed:ratio(a.exposed,all), AvgIrrelevantTokens:ratio(a.irrelevant,all), UsefulContextEfficiency:ratio(a.useful,float64(a.exposed)),
			RoutingSelectivity:ratioFloat(ratio(a.selected,all),float64(carrierCount)), ExposureFraction:ratioFloat(ratio(a.unfolded,all),float64(totalNodes)),
			NodesPerActiveToken:ratioFloat(float64(totalNodes),float64(cfg.BudgetTokens-cfg.ReserveTokens)),
		}
		if scaleIndex == 0 { baselineAccuracy = result.TargetHitAccuracy }
		if baselineAccuracy > 0 { result.KnowledgeScaleDegradation = max0(1 - result.TargetHitAccuracy/baselineAccuracy) }
		report.Scales = append(report.Scales, result)
	}
	return report,nil
}

func addWork(a *accum, result memory.FederatedResult, useful, irrelevant int) {
	a.signatures += result.Metrics.SignaturesScanned
	a.candidates += result.Metrics.CandidateCarriers
	a.selected += result.Metrics.SelectedCarriers
	a.metadataLoaded += result.Metrics.MetadataNodesLoaded
	a.touched += result.Packet.Metrics.Touched
	a.hops += result.Packet.Metrics.GraphHops
	a.unfolded += result.Packet.Metrics.Unfolded
	a.exposed += result.Packet.Budget.UsedTokens
	a.useful += useful
	a.irrelevant += irrelevant
}

func scorePositive(carrierCount, queryIndex int, query, targetCarrier string, result memory.FederatedResult) TraceRecord {
	tr:=TraceRecord{Schema:SchemaR0+".trace",CarrierCount:carrierCount,QueryIndex:queryIndex,Kind:"positive",Query:query,TargetCarrier:targetCarrier,Ranked:append([]memory.CarrierCandidate(nil),result.Ranked...),Federation:result.Metrics,PacketMetrics:result.Packet.Metrics,Budget:result.Packet.Budget,Unknown:result.Packet.Unknown,Reason:result.Packet.Reason}
	if len(result.Ranked)>0 && result.Ranked[0].ID==targetCarrier { tr.Top1Correct=true }
	for _,candidate:=range result.Ranked{if candidate.ID==targetCarrier{tr.RecallAtK=true;break}}
	for _,item:=range result.Packet.Items{
		if item.CarrierID==targetCarrier { tr.TargetHit=true;tr.UsefulTokens+=item.TokenCost;if item.Kind=="evidence"{tr.EvidenceHit=true} } else { tr.IrrelevantTokens+=item.TokenCost }
		if item.Fidelity==memory.FidelityExact && !item.Verified { tr.FalseExact++ }
	}
	return tr
}

func scoreNegative(carrierCount, queryIndex int, query string, result memory.FederatedResult) TraceRecord {
	tr:=TraceRecord{Schema:SchemaR0+".trace",CarrierCount:carrierCount,QueryIndex:queryIndex,Kind:"negative",Query:query,Ranked:append([]memory.CarrierCandidate(nil),result.Ranked...),Federation:result.Metrics,PacketMetrics:result.Packet.Metrics,Budget:result.Packet.Budget,Unknown:result.Packet.Unknown,Reason:result.Packet.Reason}
	for _,item:=range result.Packet.Items{tr.IrrelevantTokens+=item.TokenCost;if item.Fidelity==memory.FidelityExact&&!item.Verified{tr.FalseExact++}}
	return tr
}

func buildFederation(carrierCount, nodesPerCarrier int) (*memory.Federation,error) {
	fed:=memory.NewFederation()
	for i:=0;i<carrierCount;i++{
		id:=carrierID(i);topic:=topicID(i);cluster:="systems"
		root:=memory.CanonicalAddress(id,cluster,"root")
		evidence:=memory.CanonicalAddress(id,cluster,"evidence")
		nodes:=make([]memory.Node,0,nodesPerCarrier)
		nodes=append(nodes,memory.Node{Address:root,CarrierID:id,ClusterID:cluster,Kind:"concept",SemanticClass:"architecture",Label:"systems "+topic,Levels:map[memory.Fidelity]memory.ContentLevel{
			memory.FidelityAbstract:{Text:"Abstract for "+topic+" in the shared systems domain."},
			memory.FidelitySummary:{Text:"Summary for "+topic+" describes the synthetic system represented by "+id+"."},
			memory.FidelityDetail:{Text:"Detailed material for "+topic+" preserves a reopenable evidence relation while remaining distinct from unrelated carriers."},
		},Links:map[string][]string{"evidence":{evidence}},Importance:memory.ImportanceProfile{Global:.8,Local:.9,Confidence:.95},Provenance:memory.Provenance{SourceID:"synthetic:"+id}})
		nodes=append(nodes,memory.Node{Address:evidence,CarrierID:id,ClusterID:cluster,Kind:"evidence",SemanticClass:"evidence",Label:"systems evidence "+topic,Parent:root,Levels:map[memory.Fidelity]memory.ContentLevel{
			memory.FidelitySummary:{Text:"Evidence summary for "+topic+"."},
			memory.FidelityEvidence:{Text:"Verified synthetic evidence for "+topic+" in "+id+".",SourceRef:"synthetic:"+id+"#evidence",Verified:true},
			memory.FidelityExact:{Text:"EXACT-SYNTHETIC-EVIDENCE-"+topic,SourceRef:"synthetic:"+id+"#exact",Verified:true},
		},Importance:memory.ImportanceProfile{Global:.7,Local:.95,Confidence:1},Provenance:memory.Provenance{SourceID:"synthetic:"+id,SourceRef:"synthetic:"+id+"#evidence"}})
		for n:=2;n<nodesPerCarrier;n++{
			address:=memory.CanonicalAddress(id,cluster,fmt.Sprintf("filler-%04d",n))
			nodes=append(nodes,memory.Node{Address:address,CarrierID:id,ClusterID:cluster,Parent:root,Kind:"context",SemanticClass:"archive",Label:fmt.Sprintf("systems archive %s filler-%04d",id,n),Levels:map[memory.Fidelity]memory.ContentLevel{memory.FidelitySummary:{Text:fmt.Sprintf("Shared systems-domain distractor material from %s node %d.",id,n)}},Importance:memory.ImportanceProfile{Global:.2,Local:.2,Confidence:.8}})
		}
		store,err:=memory.BuildStore(id,nodes);if err!=nil{return nil,err}
		if err:=fed.Mount(id,store);err!=nil{return nil,err}
	}
	return fed,nil
}

func carrierID(i int)string{return fmt.Sprintf("carrier-%06d",i)}
func topicID(i int)string{return fmt.Sprintf("topic-%06d",i)}
func targetIndex(seed uint64,scaleIndex,queryIndex,count int)int{if count<=1{return 0};v:=seed+uint64(scaleIndex+1)*104729+uint64(queryIndex+1)*7919;return int(v%uint64(count))}
func ratio(v int,d float64)float64{if d<=0{return 0};return float64(v)/d}
func ratioFloat(v,d float64)float64{if d<=0{return 0};return v/d}
func max0(v float64)float64{if v<0{return 0};return v}

func SortedCarrierCounts(in []int) []int { out:=append([]int(nil),in...);sort.Ints(out);return out }
