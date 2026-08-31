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
	out := flag.String("out", "evidence/r6-execution-compliance", "output directory")
	flag.Parse()
	must(os.MkdirAll(*out, 0o755))

	// Keep the R5 unseen causal program exactly unchanged. R6 is allowed to alter
	// only the visible execution-compliance surface.
	program := temporal.Program{Schema:temporal.Schema,ID:"film-r5-unseen",Automaton:automaton.Automaton{Schema:automaton.Schema,ID:"reverse-chain",Cells:[]automaton.Cell{{ID:"A",InitialState:"IDLE",Neighbors:[]string{"B","C"}},{ID:"B",InitialState:"IDLE",Neighbors:[]string{"A","C"}},{ID:"C",InitialState:"ACTIVE",Neighbors:[]string{"A","B"}}},Rules:[]automaton.Rule{
		{ID:"r1",TargetCell:"A",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"C",State:"ACTIVE"}}},
		{ID:"r2",TargetCell:"B",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"A",State:"ACTIVE"}}},
		{ID:"r3",TargetCell:"C",FromState:"ACTIVE",ToState:"DONE",Requires:[]automaton.Predicate{{CellID:"B",State:"ACTIVE"}}},
		{ID:"r4",TargetCell:"A",FromState:"ACTIVE",ToState:"DONE",Requires:[]automaton.Predicate{{CellID:"C",State:"DONE"}}},
	}},MaxSteps:7,CheckpointEvery:2}

	base, baseMeta, err := fixedcarrier.RenderTemporalCarrier(program); must(err)
	grammarMutation := fixedcarrier.CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"T2_SEMANTIC_TEMPORAL_SUPERGRAPH",Value:"VISIBLE_RULE_MICROGRAMMAR_R1",Experimental:true}
	execMutation := fixedcarrier.CandidateMutation{Kind:"PROMPT",Target:"EXECUTION_POLICY",Value:"EXECUTE_VISIBLE_RULES_TO_STABLE_R1",Experimental:true}
	identityMutation := fixedcarrier.CandidateMutation{Kind:"REDUNDANCY",Target:"CELL_IDENTITY_ENCODING",Value:fixedcarrier.VisibleCellIdentityRedundancyR1,Experimental:true}
	fromMutation := fixedcarrier.CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"FROM_STATE_PRECONDITION_VISIBILITY",Value:fixedcarrier.VisibleFromStatePreconditionR1,Experimental:true}
	roleMutation := fixedcarrier.CandidateMutation{Kind:"TEMPORAL_STRUCTURE",Target:"RULE_ROLE_BINDING",Value:fixedcarrier.VisibleRuleRoleBindingR1,Experimental:true}

	grammar, _, err := fixedcarrier.BuildTemporalCandidate(base, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r6-temporal-grammar",Mutations:[]fixedcarrier.CandidateMutation{grammarMutation}}); must(err)
	execCarrier, _, err := fixedcarrier.BuildTemporalCandidate(grammar, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r6-execute-stable",Mutations:[]fixedcarrier.CandidateMutation{execMutation}}); must(err)
	identity, _, err := fixedcarrier.BuildTemporalCandidate(execCarrier, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r6-cell-identity",Mutations:[]fixedcarrier.CandidateMutation{identityMutation}}); must(err)
	fromCarrier, _, err := fixedcarrier.BuildTemporalCandidate(identity, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r6-from-state",Mutations:[]fixedcarrier.CandidateMutation{fromMutation}}); must(err)
	r5, _, err := fixedcarrier.BuildRuleRoleBindingCandidate(fromCarrier, "r6-parent-rule-role-binding"); must(err)
	r6, report, err := fixedcarrier.BuildExecutionComplianceCandidate(r5, "execution-policy-compliance-cross-model-r1"); must(err)
	decoded, err := fixedcarrier.DecodeTemporalCarrierPNG(r6); must(err)
	if decoded.ProgramSHA256 != baseMeta.ProgramSHA256 || !report.ExactProgramPreserved { panic("R6 changed exact TemporalProgram") }

	inherited := []fixedcarrier.CandidateMutation{grammarMutation,execMutation,identityMutation,fromMutation,roleMutation}
	build := fixedcarrier.TemporalInteropBuildManifestWithInherited(report,decoded,inherited)
	parentDecoded, err := fixedcarrier.DecodeTemporalCarrierPNG(r5); must(err)
	parentSemantics := fixedcarrier.TemporalSemanticManifest(parentDecoded,inherited)
	trace, err := temporal.Run(program); must(err)
	if trace.Final["A"]!="DONE" || trace.Final["B"]!="ACTIVE" || trace.Final["C"]!="DONE" { panic(fmt.Sprintf("unexpected final state: %#v",trace.Final)) }

	panel := []map[string]any{
		{"model_id":"deepseek-unspecified","mode":"PRESERVE_PASS","baseline_pass":true,"required_candidate_pass":true},
		{"model_id":"qwen-unspecified","mode":"IMPROVE_TO_PASS","baseline_pass":false,"required_candidate_pass":true},
	}
	candidate := map[string]any{
		"schema":"tlaloc.candidate-manifest.r1",
		"id":"execution-policy-compliance-cross-model-r1",
		"parent_id":"rule-role-binding-unseen-rules-r1",
		"program_sha256":baseMeta.ProgramSHA256,
		"mutations":[]map[string]string{{"kind":"PROMPT","target":"EXECUTION_POLICY_COMPLIANCE","value":fixedcarrier.ExecuteDontSummarizeToStableR1}},
		"changed_modules":[]string{"EXECUTION_POLICY_COMPLIANCE"},
		"preserved_modules":[]string{"CELL_IDENTITY_ENCODING","FROM_STATE_PRECONDITION_VISIBILITY","RULE_ROLE_BINDING","TEMPORAL_GRAMMAR","EXECUTION_POLICY","PROGRAM_SEMANTICS","PAYLOAD","INITIAL_STATES"},
		"forbidden_changes":[]string{"CELL_IDENTITY_MUTATION","RULE_ROLE_BINDING_MUTATION","FROM_STATE_PRECONDITION_MUTATION","RULE_MUTATION","STATE_MUTATION","CHECKPOINT_MUTATION","PAYLOAD_MUTATION"},
		"expected_semantic_changes":[]map[string]string{{"key":"EXECUTION_POLICY_COMPLIANCE","value":fixedcarrier.ExecuteDontSummarizeToStableR1}},
		"expected_effect":"make execution mandatory rather than descriptive while preserving DeepSeek R5 PASS and requiring Qwen to improve to PASS",
		"compatibility_panel":panel,
	}
	generation := map[string]any{
		"schema":"origami.r6-execution-compliance-generation.r1",
		"candidate_id":"execution-policy-compliance-cross-model-r1",
		"parent_id":"rule-role-binding-unseen-rules-r1",
		"program_sha256":baseMeta.ProgramSHA256,
		"artifact_sha256":build.ArtifactSHA256,
		"artifact_bytes":build.ArtifactBytes,
		"exact_program_preserved":report.ExactProgramPreserved,
		"compatibility_panel":panel,
		"expected_stable_state":map[string]string{"A":"DONE","B":"ACTIVE","C":"DONE"},
	}

	must(os.WriteFile(filepath.Join(*out,"execution-policy-compliance-cross-model-r1.png"),r6,0o644))
	writeJSON(filepath.Join(*out,"candidate-manifest.json"),candidate)
	writeJSON(filepath.Join(*out,"build-manifest.json"),build)
	writeJSON(filepath.Join(*out,"parent-semantic-manifest.json"),parentSemantics)
	writeJSON(filepath.Join(*out,"semantic-manifest.json"),build.VisibleSemantics)
	writeJSON(filepath.Join(*out,"visible-text-manifest.json"),build.VisibleText)
	writeJSON(filepath.Join(*out,"trace.json"),trace)
	writeJSON(filepath.Join(*out,"generation-report.json"),generation)
	fmt.Printf("R6_PNG=%s\nPROGRAM_SHA256=%s\nARTIFACT_SHA256=%s\nBYTES=%d\n",filepath.Join(*out,"execution-policy-compliance-cross-model-r1.png"),baseMeta.ProgramSHA256,build.ArtifactSHA256,build.ArtifactBytes)
}

func writeJSON(path string,v any){b,err:=json.MarshalIndent(v,"","  ");must(err);b=append(b,'\n');must(os.WriteFile(path,b,0o644))}
func must(err error){if err!=nil{panic(err)}}
