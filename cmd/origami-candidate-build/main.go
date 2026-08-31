package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "capabilities" {
		writeJSON(fixedcarrier.TemporalCandidateCapabilities())
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "semantic-manifest" {
		semanticManifest(os.Args[2:])
		return
	}
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "build" { args = args[1:] }
	fs := flag.NewFlagSet("origami-candidate-build", flag.ExitOnError)
	parent := fs.String("parent", env("TLALOC_PARENT_PNG"), "parent temporal PNG")
	out := fs.String("out", env("TLALOC_OUTPUT_PNG"), "candidate PNG output")
	id := fs.String("id", env("TLALOC_CANDIDATE_ID"), "candidate id")
	mutationsRaw := fs.String("mutations-json", env("TLALOC_MUTATIONS_JSON"), "JSON array of experimental mutations")
	specPath := fs.String("spec", "", "optional full CandidateSpec JSON")
	reportPath := fs.String("report", "", "optional legacy build report JSON output")
	interopPath := fs.String("interop-report", "", "optional origami.candidate-build-manifest.r1 output")
	fs.Parse(args)
	if strings.TrimSpace(*parent) == "" || strings.TrimSpace(*out) == "" { die(fmt.Errorf("parent and out are required (flags or TLALOC_PARENT_PNG/TLALOC_OUTPUT_PNG)")) }
	var spec fixedcarrier.CandidateSpec
	if *specPath != "" {
		body, err := os.ReadFile(*specPath); die(err)
		dec := json.NewDecoder(bytes.NewReader(body)); dec.DisallowUnknownFields(); die(dec.Decode(&spec))
	} else {
		if strings.TrimSpace(*id) == "" { die(fmt.Errorf("candidate id is required")) }
		mutations, err := fixedcarrier.ParseCandidateMutationsJSON(*mutationsRaw); die(err)
		spec = fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:*id,Mutations:mutations}
	}
	parentBytes, err := os.ReadFile(*parent); die(err)
	candidate, report, err := fixedcarrier.BuildTemporalCandidate(parentBytes, spec); die(err)
	if err := os.WriteFile(*out, candidate, 0o644); err != nil { die(err) }
	encoded, err := json.MarshalIndent(report, "", "  "); die(err); encoded=append(encoded,'\n')
	if *reportPath != "" { die(os.WriteFile(*reportPath,encoded,0o644)) }
	if *interopPath != "" {
		decoded,err:=fixedcarrier.DecodeTemporalCarrierPNG(candidate);die(err)
		interop:=fixedcarrier.TemporalInteropBuildManifest(report,decoded)
		b,err:=json.MarshalIndent(interop,"","  ");die(err);die(os.WriteFile(*interopPath,append(b,'\n'),0o644))
	}
	_, err = os.Stdout.Write(encoded); die(err)
}

func semanticManifest(args []string){
	fs:=flag.NewFlagSet("semantic-manifest",flag.ExitOnError)
	in:=fs.String("in","","temporal carrier PNG")
	mutationsRaw:=fs.String("mutations-json","[]","optional mutations represented visibly by this carrier")
	out:=fs.String("out","","optional semantic manifest output")
	fs.Parse(args)
	if strings.TrimSpace(*in)==""{die(fmt.Errorf("-in is required"))}
	body,err:=os.ReadFile(*in);die(err);decoded,err:=fixedcarrier.DecodeTemporalCarrierPNG(body);die(err)
	mutations,err:=fixedcarrier.ParseCandidateMutationsJSON(*mutationsRaw);die(err)
	manifest:=fixedcarrier.TemporalSemanticManifest(decoded,mutations);b,err:=json.MarshalIndent(manifest,"","  ");die(err);b=append(b,'\n')
	if *out!=""{die(os.WriteFile(*out,b,0o644))}
	_,err=os.Stdout.Write(b);die(err)
}

func env(name string) string { return strings.TrimSpace(os.Getenv(name)) }
func writeJSON(v any){b,err:=json.MarshalIndent(v,"","  ");die(err);b=append(b,'\n');_,err=os.Stdout.Write(b);die(err)}
func die(err error){if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
