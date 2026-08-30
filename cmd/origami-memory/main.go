package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/lab/hybridruntime"
	"github.com/LuigiD5555/origami/internal/memory"
)

func main(){
	carrier:=flag.String("carrier","runs/hybrid-carrier-synthetic-r0/public/carrier.png","Origami carrier PNG")
	op:=flag.String("op","QUERY","QUERY|EXPAND|SIGNATURE|VERIFY|LOOKUP")
	query:=flag.String("query","","natural-language query or lookup address/key")
	address:=flag.String("address","","address for EXPAND")
	budget:=flag.Int("budget",4000,"active model-facing token-equivalent budget")
	reserve:=flag.Int("reserve",0,"tokens reserved for model response/control state")
	fidelity:=flag.String("fidelity","detail","label|abstract|summary|detail|evidence|exact")
	columns:=flag.Int("columns",24,"exact transport columns")
	cell:=flag.Int("cell",56,"exact transport cell pixels")
	margin:=flag.Int("margin",8,"exact transport margin pixels")
	band:=flag.Int("memory-band",192,"visual memory-band height; set 0 for legacy carrier")
	flag.Parse()
	pngBytes,err:=os.ReadFile(*carrier);must(err)
	runtime,err:=hybridruntime.OpenPNG(pngBytes,hybridruntime.Transport{Columns:*columns,Cell:*cell,Margin:*margin,MemoryBandHeight:*band});must(err)
	var result any
	switch *op{
	case "QUERY":if *query==""{must(fmt.Errorf("QUERY requires -query"))};result=runtime.MemoryQuery(*query,memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve})
	case "EXPAND":if *address==""{must(fmt.Errorf("EXPAND requires -address"))};result=runtime.Expand(*address,memory.Fidelity(*fidelity),memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve})
	case "SIGNATURE":result=runtime.Signature()
	case "VERIFY":result=runtime.Verify()
	case "LOOKUP":if *query==""{must(fmt.Errorf("LOOKUP requires -query"))};result=runtime.Lookup(*query)
	default:must(fmt.Errorf("unsupported operation %q",*op))
	}
	out:=struct{Open hybridruntime.OpenMetrics `json:"open_metrics"`;Result any `json:"result"`}{Open:runtime.OpenMetrics(),Result:result};b,err:=json.MarshalIndent(out,"","  ");must(err);fmt.Println(string(b))
}
func must(err error){if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
