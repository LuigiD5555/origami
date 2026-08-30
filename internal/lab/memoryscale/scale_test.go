package memoryscale

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunPreservesAccuracyWhileKnowledgeScales(t *testing.T) {
	cfg:=Config{CarrierCounts:[]int{1,10,100},NodesPerCarrier:16,PositiveQueriesPerScale:8,NegativeQueriesPerScale:2,BudgetTokens:4000,CarrierLimit:4,Seed:7}
	var trace bytes.Buffer
	report,err:=Run(cfg,&trace);if err!=nil{t.Fatal(err)}
	if report.Schema!=SchemaR0{t.Fatalf("unexpected schema %q",report.Schema)}
	if len(report.Scales)!=3{t.Fatalf("unexpected scale count %d",len(report.Scales))}
	for _,scale:=range report.Scales{
		if scale.CarrierTop1Accuracy!=1{t.Fatalf("top1 degraded at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.CarrierRecallAtK!=1{t.Fatalf("recall degraded at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.TargetHitAccuracy!=1{t.Fatalf("target retrieval degraded at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.EvidenceHitAccuracy!=1{t.Fatalf("evidence retrieval degraded at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.NegativeUnknownAccuracy!=1{t.Fatalf("negative fail-closed degraded at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.BudgetViolations!=0{t.Fatalf("budget violation at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.FalseExact!=0{t.Fatalf("false exact at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.KnowledgeScaleDegradation!=0{t.Fatalf("unexpected KSD at %d carriers: %+v",scale.CarrierCount,scale)}
		if scale.AvgSignaturesScanned!=float64(scale.CarrierCount){t.Fatalf("signature work must be explicit: %+v",scale)}
		if scale.AvgTokensExposed>=float64(scale.ActiveBudgetTokens){t.Fatalf("benchmark should remain below active budget: %+v",scale)}
		if scale.UsefulContextEfficiency<=0||scale.UsefulContextEfficiency>1{t.Fatalf("invalid UCE: %+v",scale)}
	}
	if report.Scales[2].NodesPerActiveToken<=report.Scales[0].NodesPerActiveToken{t.Fatalf("knowledge/context ratio did not grow: %+v",report.Scales)}
	lines:=strings.Count(strings.TrimSpace(trace.String()),"\n")+1
	want:=(cfg.PositiveQueriesPerScale+cfg.NegativeQueriesPerScale)*len(cfg.CarrierCounts)
	if lines!=want{t.Fatalf("trace lines=%d want=%d",lines,want)}
}

func TestMetadataLoadingIsBoundedBySelectedCarriers(t *testing.T){
	cfg:=Config{CarrierCounts:[]int{100},NodesPerCarrier:32,PositiveQueriesPerScale:4,NegativeQueriesPerScale:0,BudgetTokens:4000,CarrierLimit:3,Seed:11}
	report,err:=Run(cfg,nil);if err!=nil{t.Fatal(err)}
	scale:=report.Scales[0]
	if scale.AvgSelectedCarriers>3{t.Fatalf("selected too many carriers: %+v",scale)}
	if scale.AvgMetadataNodesLoaded>float64(3*cfg.NodesPerCarrier){t.Fatalf("loaded metadata outside selected carrier bound: %+v",scale)}
	if scale.ExposureFraction>=0.1{t.Fatalf("too much of total memory unfolded: %+v",scale)}
}

func TestThousandCarrierCanonicalDensitySmoke(t *testing.T){
	cfg:=Config{CarrierCounts:[]int{1000},NodesPerCarrier:64,PositiveQueriesPerScale:4,NegativeQueriesPerScale:1,BudgetTokens:4000,CarrierLimit:4,Seed:20260829}
	report,err:=Run(cfg,nil);if err!=nil{t.Fatal(err)}
	s:=report.Scales[0]
	if s.TotalNodes!=64000{t.Fatalf("unexpected total nodes: %+v",s)}
	if s.CarrierTop1Accuracy!=1||s.TargetHitAccuracy!=1||s.EvidenceHitAccuracy!=1{t.Fatalf("1000-carrier retrieval failed: %+v",s)}
	if s.NegativeUnknownAccuracy!=1{t.Fatalf("1000-carrier UNKNOWN failure: %+v",s)}
	if s.BudgetViolations!=0||s.FalseExact!=0{t.Fatalf("1000-carrier invariant failure: %+v",s)}
	if s.AvgSignaturesScanned!=1000{t.Fatalf("global routing work must remain visible: %+v",s)}
	if s.AvgMetadataNodesLoaded>float64(cfg.CarrierLimit*cfg.NodesPerCarrier){t.Fatalf("local metadata loading escaped selected carriers: %+v",s)}
	if s.ExposureFraction>=0.01{t.Fatalf("too much of 64k-node universe was unfolded: %+v",s)}
}

func TestConfigRejectsUnsafeOrAmbiguousRanges(t *testing.T){
	cases:=[]Config{
		{CarrierCounts:nil,NodesPerCarrier:3,PositiveQueriesPerScale:1,BudgetTokens:10,CarrierLimit:1},
		{CarrierCounts:[]int{10,10},NodesPerCarrier:3,PositiveQueriesPerScale:1,BudgetTokens:10,CarrierLimit:1},
		{CarrierCounts:[]int{1},NodesPerCarrier:2,PositiveQueriesPerScale:1,BudgetTokens:10,CarrierLimit:1},
		{CarrierCounts:[]int{1},NodesPerCarrier:3,PositiveQueriesPerScale:1,BudgetTokens:10,ReserveTokens:10,CarrierLimit:1},
	}
	for i,cfg:=range cases{if err:=cfg.Validate();err==nil{t.Fatalf("case %d should fail validation",i)}}
}
