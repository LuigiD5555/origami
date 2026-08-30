package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/document"
	"github.com/LuigiD5555/origami/internal/evidence"
	"github.com/LuigiD5555/origami/internal/integration"
	"github.com/LuigiD5555/origami/internal/memory"
)

type proposalInput struct {
	Policy    evidence.Policy     `json:"policy,omitempty"`
	Proposals []evidence.Proposal `json:"proposals"`
}

type verifyResult struct {
	Schema      string                `json:"schema"`
	SourceID    string                `json:"source_id"`
	SourceHash  string                `json:"source_hash"`
	StoreRoot   string                `json:"store_root_sha256"`
	Signature   memory.GraphSignature `json:"graph_signature"`
	NodeCount   int                   `json:"node_count"`
}

func main(){
	bundlePath:=flag.String("bundle","","origami.document.r0 bundle JSON")
	op:=flag.String("op","QUERY","QUERY|EXPAND|LOOKUP|SIGNATURE|VERIFY|REDUCE|PROMOTE")
	query:=flag.String("query","","natural-language query or lookup address")
	address:=flag.String("address","","address for EXPAND/LOOKUP")
	budget:=flag.Int("budget",4000,"active model-facing token-equivalent budget")
	reserve:=flag.Int("reserve",0,"tokens reserved for answer/control state")
	fidelity:=flag.String("fidelity","detail","label|abstract|summary|detail|evidence|exact")
	proposalsPath:=flag.String("proposals","","proposal JSON for REDUCE/PROMOTE")
	writeBundle:=flag.String("write-bundle","","write updated bundle after PROMOTE")
	flag.Parse()
	if *bundlePath==""{die(fmt.Errorf("-bundle is required"))}
	body,err:=os.ReadFile(*bundlePath);die(err);var bundle document.Bundle;die(json.Unmarshal(body,&bundle));if bundle.Schema!=document.SchemaR0+".bundle"{die(fmt.Errorf("unexpected bundle schema %q",bundle.Schema))};die(document.Validate(bundle.Document))
	system,err:=integration.NewSystem(bundle.Document.SourceID,bundle.Nodes);die(err)
	var result any
	switch *op{
	case "QUERY":if *query==""{die(fmt.Errorf("QUERY requires -query"))};result=system.Query(*query,memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve})
	case "EXPAND":if *address==""{die(fmt.Errorf("EXPAND requires -address"))};result=system.Expand(*address,memory.Fidelity(*fidelity),memory.ContextBudget{MaxTokens:*budget,ReserveTokens:*reserve})
	case "LOOKUP":target:=*address;if target==""{target=*query};if target==""{die(fmt.Errorf("LOOKUP requires -address or -query"))};node,ok:=system.Lookup(target);if !ok{die(fmt.Errorf("address not found: %s",target))};result=node
	case "SIGNATURE":result=system.Signature()
	case "VERIFY":result=verifyResult{Schema:integration.SchemaR0+".verify",SourceID:bundle.Document.SourceID,SourceHash:bundle.Document.SourceHash,StoreRoot:integration.CanonicalStoreRoot(system.Nodes()),Signature:system.Signature(),NodeCount:len(system.Nodes())}
	case "REDUCE","PROMOTE":
		if *proposalsPath==""{die(fmt.Errorf("%s requires -proposals",*op))};pb,err:=os.ReadFile(*proposalsPath);die(err);var input proposalInput;die(json.Unmarshal(pb,&input));reduction,err:=system.Reduce(input.Proposals,input.Policy);die(err)
		if *op=="REDUCE"{result=reduction;break}
		promoted,err:=system.Promote(reduction);die(err);bundle.Nodes=system.Nodes();if *writeBundle!=""{updated,err:=json.MarshalIndent(bundle,"","  ");die(err);updated=append(updated,'\n');die(os.WriteFile(*writeBundle,updated,0644))};result=struct{Reduction evidence.Reduction `json:"reduction"`;Promoted []memory.Node `json:"promoted"`;StoreRoot string `json:"store_root_sha256"`}{reduction,promoted,integration.CanonicalStoreRoot(system.Nodes())}
	default:die(fmt.Errorf("unsupported operation %q",*op))
	}
	out,err:=json.MarshalIndent(result,"","  ");die(err);fmt.Println(string(out))
}

func die(err error){if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
