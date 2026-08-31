package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func main() {
	out := flag.String("out", "evidence/r3-cell-identity", "output directory")
	flag.Parse()
	must(os.MkdirAll(*out, 0o755))

	program := temporal.Program{
		Schema: temporal.Schema, ID: "film",
		Automaton: automaton.Automaton{
			Schema: automaton.Schema, ID: "chain",
			Cells: []automaton.Cell{
				{ID:"A", InitialState:"ACTIVE", Neighbors:[]string{"B"}},
				{ID:"B", InitialState:"IDLE", Neighbors:[]string{"A","C"}},
				{ID:"C", InitialState:"IDLE", Neighbors:[]string{"B"}},
			},
			Rules: []automaton.Rule{
				{ID:"r1",TargetCell:"B",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"A",State:"ACTIVE"}}},
				{ID:"r2",TargetCell:"A",FromState:"ACTIVE",ToState:"DONE",Requires:[]automaton.Predicate{{CellID:"B",State:"ACTIVE"}}},
				{ID:"r3",TargetCell:"C",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"B",State:"ACTIVE"}}},
				{ID:"r4",TargetCell:"B",FromState:"ACTIVE",ToState:"DONE",Requires:[]automaton.Predicate{{CellID:"C",State:"ACTIVE"}}},
			},
		},
		MaxSteps: 6, CheckpointEvery: 2,
	}

	base, baseMeta, err := fixedcarrier.RenderTemporalCarrier(program); must(err)
	grammarMutation := fixedcarrier.CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true}
	execMutation := fixedcarrier.CandidateMutation{Kind:"PROMPT",Target:"EXECUTION_POLICY",Value:"EXECUTE_VISIBLE_RULES_TO_STABLE_R1",Experimental:true}
	identityMutation := fixedcarrier.CandidateMutation{Kind:"REDUNDANCY",Target:"CELL_IDENTITY_ENCODING",Value:fixedcarrier.VisibleCellIdentityRedundancyR1,Experimental:true}

	grammar, _, err := fixedcarrier.BuildTemporalCandidate(base, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"temporal-grammar-visible-r1",Mutations:[]fixedcarrier.CandidateMutation{grammarMutation}}); must(err)
	execCarrier, _, err := fixedcarrier.BuildTemporalCandidate(grammar, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"execute-to-stable-text-r1",Mutations:[]fixedcarrier.CandidateMutation{execMutation}}); must(err)
	r3, report, err := fixedcarrier.BuildTemporalCandidate(execCarrier, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"cell-identity-redundancy-r1",Mutations:[]fixedcarrier.CandidateMutation{identityMutation}}); must(err)
	decoded, err := fixedcarrier.DecodeTemporalCarrierPNG(r3); must(err)
	if decoded.ProgramSHA256 != baseMeta.ProgramSHA256 || !report.ExactProgramPreserved { panic("R3 changed exact TemporalProgram") }

	inherited := []fixedcarrier.CandidateMutation{grammarMutation, execMutation}
	build := fixedcarrier.TemporalInteropBuildManifestWithInherited(report, decoded, inherited)
	parentDecoded, err := fixedcarrier.DecodeTemporalCarrierPNG(execCarrier); must(err)
	parentSemantics := fixedcarrier.TemporalSemanticManifest(parentDecoded, inherited)
	trace, err := temporal.Run(program); must(err)
	if trace.Final["A"]!="DONE" || trace.Final["B"]!="DONE" || trace.Final["C"]!="ACTIVE" { panic(fmt.Sprintf("unexpected final state: %#v",trace.Final)) }

	candidate := map[string]any{
		"schema":"tlaloc.candidate-manifest.r1",
		"id":"cell-identity-redundancy-r1",
		"parent_id":"execute-to-stable-text-r1",
		"program_sha256":baseMeta.ProgramSHA256,
		"mutations":[]map[string]string{{"kind":"REDUNDANCY","target":"CELL_IDENTITY_ENCODING","value":fixedcarrier.VisibleCellIdentityRedundancyR1}},
		"changed_modules":[]string{"CELL_IDENTITY_ENCODING"},
		"preserved_modules":[]string{"TEMPORAL_GRAMMAR","EXECUTION_POLICY","PROGRAM_SEMANTICS","PAYLOAD","INITIAL_STATES"},
		"forbidden_changes":[]string{"RULE_MUTATION","STATE_MUTATION","EXECUTION_POLICY_MUTATION","CHECKPOINT_MUTATION","PAYLOAD_MUTATION"},
		"expected_semantic_changes":[]map[string]string{{"key":"VISIBLE_CELL_ID_A","value":"A[01]"},{"key":"VISIBLE_CELL_ID_B","value":"B[02]"},{"key":"VISIBLE_CELL_ID_C","value":"C[03]"}},
		"expected_effect":"reduce A/B/C confusion while preserving rule recovery and execution",
	}
	generation := map[string]any{
		"schema":"origami.r3-cell-identity-generation.r1",
		"candidate_id":"cell-identity-redundancy-r1",
		"parent_id":"execute-to-stable-text-r1",
		"program_sha256":baseMeta.ProgramSHA256,
		"artifact_sha256":build.ArtifactSHA256,
		"artifact_bytes":build.ArtifactBytes,
		"exact_program_preserved":report.ExactProgramPreserved,
		"four_rule_program":true,
		"expected_stable_state":map[string]string{"A":"DONE","B":"DONE","C":"ACTIVE"},
	}

	must(os.WriteFile(filepath.Join(*out,"cell-identity-redundancy-r1.png"), r3, 0o644))
	writeJSON(filepath.Join(*out,"candidate-manifest.json"), candidate)
	writeJSON(filepath.Join(*out,"build-manifest.json"), build)
	writeJSON(filepath.Join(*out,"parent-semantic-manifest.json"), parentSemantics)
	writeJSON(filepath.Join(*out,"semantic-manifest.json"), build.VisibleSemantics)
	writeJSON(filepath.Join(*out,"visible-text-manifest.json"), build.VisibleText)
	writeJSON(filepath.Join(*out,"trace.json"), trace)
	writeJSON(filepath.Join(*out,"generation-report.json"), generation)
	fmt.Printf("R3_PNG=%s\nPROGRAM_SHA256=%s\nARTIFACT_SHA256=%s\nBYTES=%d\n", filepath.Join(*out,"cell-identity-redundancy-r1.png"), baseMeta.ProgramSHA256, build.ArtifactSHA256, build.ArtifactBytes)
}

func writeJSON(path string, v any) {
	b, err := json.MarshalIndent(v,"","  "); must(err)
	b=append(b,'\n'); must(os.WriteFile(path,b,0o644))
}

func must(err error) { if err != nil { panic(err) } }
