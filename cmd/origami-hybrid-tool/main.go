package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/lab/hybridruntime"
	"github.com/LuigiD5555/origami/internal/memory"
	"github.com/LuigiD5555/origami/internal/receiver"
)

type ModelPacket struct { Schema string `json:"schema"`; ReceiverContract string `json:"receiver_contract"`; CarrierSHA256 string `json:"carrier_sha256"`; Transport hybridruntime.Transport `json:"transport"`; AllowedTools []string `json:"allowed_tools"` }

func main(){
	carrierPath:=flag.String("carrier","runs/hybrid-carrier-synthetic-r0/public/carrier.png","public carrier PNG");packetPath:=flag.String("packet","runs/hybrid-carrier-synthetic-r0/public/model_packet.json","public model packet");op:=flag.String("op","BOOT","BOOT|LOOKUP|FOLLOW|TRACE|QUERY|EXPAND|SIGNATURE|VERIFY|STOP");query:=flag.String("query","","key/address or natural-language query");address:=flag.String("address","","address for EXPAND");relation:=flag.String("relation","depends","relation for FOLLOW/TRACE");depth:=flag.Int("depth",1,"maximum relation depth for FOLLOW/TRACE");budget:=flag.Int("budget",4000,"active context token-equivalent budget");reserve:=flag.Int("reserve",0,"tokens reserved outside retrieved context");fidelity:=flag.String("fidelity","detail","fidelity for EXPAND");flag.Parse()
	carrier,err:=os.ReadFile(*carrierPath);must(err);packetBytes,err:=os.ReadFile(*packetPath);must(err);var packet ModelPacket;must(json.Unmarshal(packetBytes,&packet));if packet.Schema!="origami.model-packet.r0"{must(fmt.Errorf("unexpected model packet schema %q",packet.Schema))};if packet.ReceiverContract!=receiver.ContractID{must(fmt.Errorf("unexpected receiver contract %q",packet.ReceiverContract))};if got:=hash(carrier);packet.CarrierSHA256==""||got!=packet.CarrierSHA256{must(fmt.Errorf("carrier identity mismatch: packet=%s actual=%s",packet.CarrierSHA256,got))};if !allowed(packet.AllowedTools,*op){must(fmt.Errorf("operation %s is not declared in model_packet.json",*op))}
	runtime,err:=hybridruntime.OpenPNG(carrier,packet.Transport);must(err);var result any
	switch *op{case "BOOT":result=runtime.Boot();case "LOOKUP":if *query==""{must(fmt.Errorf("LOOKUP requires -query"))};result=runtime.Lookup(*query);case "FOLLOW","TRACE":if *query==""{must(fmt.Errorf("%s requires -query",*op))};r:=runtime.Follow(*query,*relation,*depth);r.Operation=*op;result=r;case "QUERY":if *query==""{must(fmt.Errorf("QUERY requires -query"))};result=runtime.MemoryQuery(*query,memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve});case "EXPAND":target:=*address;if target==""{target=*query};if target==""{must(fmt.Errorf("EXPAND requires -address or -query"))};result=runtime.Expand(target,memory.Fidelity(*fidelity),memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve});case "SIGNATURE":result=runtime.Signature();case "VERIFY":result=runtime.Verify();case "STOP":result=hybridruntime.Result{Operation:"STOP"};default:must(fmt.Errorf("unsupported operation %q",*op))}
	out:=struct{Open hybridruntime.OpenMetrics `json:"open_metrics"`;Result any `json:"result"`}{Open:runtime.OpenMetrics(),Result:result};b,err:=json.MarshalIndent(out,"","  ");must(err);fmt.Println(string(b))
}
func allowed(ops []string,want string)bool{for _,op:=range ops{if op==want{return true}};return false}
func hash(b []byte)string{s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
func must(err error){if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
