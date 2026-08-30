package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/lab/memoryscale"
)

func main(){
	configPath:=flag.String("config","experiments/memory-scale-r0/config.json","memory scale laboratory config")
	outPath:=flag.String("out","runs/memory-scale-r0/report.json","summary report JSON")
	tracePath:=flag.String("trace","runs/memory-scale-r0/trace.jsonl","per-query trace JSONL")
	flag.Parse()

	cfgBytes,err:=os.ReadFile(*configPath);must(err)
	var cfg memoryscale.Config
	must(json.Unmarshal(cfgBytes,&cfg));must(cfg.Validate())
	must(os.MkdirAll(filepath.Dir(*outPath),0o755));must(os.MkdirAll(filepath.Dir(*tracePath),0o755))
	trace,err:=os.Create(*tracePath);must(err)
	report,runErr:=memoryscale.Run(cfg,trace)
	closeErr:=trace.Close();must(runErr);must(closeErr)
	reportBytes,err:=json.MarshalIndent(report,"","  ");must(err);reportBytes=append(reportBytes,'\n');must(os.WriteFile(*outPath,reportBytes,0o644))

	fmt.Printf("SCHEMA=%s\n",report.Schema)
	fmt.Printf("REPORT=%s\nTRACE=%s\n",*outPath,*tracePath)
	fmt.Printf("ACTIVE_BUDGET=%d\n",cfg.BudgetTokens-cfg.ReserveTokens)
	for _,s:=range report.Scales{
		fmt.Printf("SCALE carriers=%d nodes=%d top1=%.4f target=%.4f evidence=%.4f unknown=%.4f uce=%.4f ksd=%.4f signatures=%.2f metadata_nodes=%.2f exposed_tokens=%.2f false_exact=%d budget_violations=%d\n",
			s.CarrierCount,s.TotalNodes,s.CarrierTop1Accuracy,s.TargetHitAccuracy,s.EvidenceHitAccuracy,s.NegativeUnknownAccuracy,s.UsefulContextEfficiency,s.KnowledgeScaleDegradation,s.AvgSignaturesScanned,s.AvgMetadataNodesLoaded,s.AvgTokensExposed,s.FalseExact,s.BudgetViolations)
	}
}

func must(err error){if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
