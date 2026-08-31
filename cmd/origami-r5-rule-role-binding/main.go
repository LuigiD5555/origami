package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func main() {
	out := flag.String("out", "evidence/r5-rule-role-binding", "output directory")
	flag.Parse()
	must(os.MkdirAll(*out, 0o755))

	// R5 deliberately uses an unseen causal program. This prevents success from
	// being explained by memorizing the previous A->B->C rule pattern.
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

	grammar, _, err := fixedcarrier.BuildTemporalCandidate(base, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r5-temporal-grammar",Mutations:[]fixedcarrier.CandidateMutation{grammarMutation}}); must(err)
	execCarrier, _, err := fixedcarrier.BuildTemporalCandidate(grammar, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r5-execute-stable",Mutations:[]fixedcarrier.CandidateMutation{execMutation}}); must(err)
	identity, _, err := fixedcarrier.BuildTemporalCandidate(execCarrier, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r5-cell-identity",Mutations:[]fixedcarrier.CandidateMutation{identityMutation}}); must(err)
	fromCarrier, _, err := fixedcarrier.BuildTemporalCandidate(identity, fixedcarrier.CandidateSpec{Schema:fixedcarrier.CandidateSpecSchema,ID:"r5-from-state",Mutations:[]fixedcarrier.CandidateMutation{fromMutation}}); must(err)
	r5, report, err := fixedcarrier.BuildRuleRoleBindingCandidate(fromCarrier, "rule-role-binding-unseen-rules-r1"); must(err)
	decoded, err := fixedcarrier.DecodeTemporalCarrierPNG(r5); must(err)
	if decoded.ProgramSHA256 != baseMeta.ProgramSHA256 || !report.ExactProgramPreserved { panic("R5 changed exact TemporalProgram") }

	inherited := []fixedcarrier.CandidateMutation{grammarMutation,execMutation,identityMutation,fromMutation}
	build := fixedcarrier.TemporalInteropBuildManifestWithInherited(report,decoded,inherited)
	build.VisibleSemantics.Facts = setFact(build.VisibleSemantics.Facts,"RULE_ROLE_BINDING",fixedcarrier.VisibleRuleRoleBindingR1)
	build.VisibleText.Facts = roleVisibleText(decoded)
	parentDecoded, err := fixedcarrier.DecodeTemporalCarrierPNG(fromCarrier); must(err)
	parentSemantics := fixedcarrier.TemporalSemanticManifest(parentDecoded,inherited)
	trace, err := temporal.Run(program); must(err)
	if trace.Final["A"]!="DONE" || trace.Final["B"]!="ACTIVE" || trace.Final["C"]!="DONE" { panic(fmt.Sprintf("unexpected final state: %#v",trace.Final)) }

	candidate := map[string]any{"schema":"tlaloc.candidate-manifest.r1","id":"rule-role-binding-unseen-rules-r1","parent_id":"from-state-precondition-visible-r1","program_sha256":baseMeta.ProgramSHA256,"mutations":[]map[string]string{{"kind":"TEMPORAL_STRUCTURE","target":"RULE_ROLE_BINDING","value":fixedcarrier.VisibleRuleRoleBindingR1}},"changed_modules":[]string{"RULE_ROLE_BINDING"},"preserved_modules":[]string{"CELL_IDENTITY_ENCODING","FROM_STATE_PRECONDITION_VISIBILITY","TEMPORAL_GRAMMAR","EXECUTION_POLICY","PROGRAM_SEMANTICS","PAYLOAD"},"forbidden_changes":[]string{"CELL_IDENTITY_MUTATION","EXECUTION_POLICY_MUTATION","CHECKPOINT_MUTATION","PAYLOAD_MUTATION"},"expected_semantic_changes":[]map[string]string{{"key":"RULE_ROLE_BINDING","value":fixedcarrier.VisibleRuleRoleBindingR1}},"expected_effect":"bind WHEN source, TARGET identity, REQUIRE target-from-state, and SET target-to-state on unseen rules"}
	generation := map[string]any{"schema":"origami.r5-rule-role-generation.r1","candidate_id":"rule-role-binding-unseen-rules-r1","program_sha256":baseMeta.ProgramSHA256,"artifact_sha256":build.ArtifactSHA256,"artifact_bytes":build.ArtifactBytes,"exact_program_preserved":report.ExactProgramPreserved,"unseen_rule_program":true,"expected_stable_state":map[string]string{"A":"DONE","B":"ACTIVE","C":"DONE"}}

	must(os.WriteFile(filepath.Join(*out,"rule-role-binding-unseen-rules-r1.png"),r5,0o644))
	writeJSON(filepath.Join(*out,"candidate-manifest.json"),candidate)
	writeJSON(filepath.Join(*out,"build-manifest.json"),build)
	writeJSON(filepath.Join(*out,"parent-semantic-manifest.json"),parentSemantics)
	writeJSON(filepath.Join(*out,"semantic-manifest.json"),build.VisibleSemantics)
	writeJSON(filepath.Join(*out,"visible-text-manifest.json"),build.VisibleText)
	writeJSON(filepath.Join(*out,"trace.json"),trace)
	writeJSON(filepath.Join(*out,"generation-report.json"),generation)
	fmt.Printf("R5_PNG=%s\nPROGRAM_SHA256=%s\nARTIFACT_SHA256=%s\nBYTES=%d\n",filepath.Join(*out,"rule-role-binding-unseen-rules-r1.png"),baseMeta.ProgramSHA256,build.ArtifactSHA256,build.ArtifactBytes)
	_ = roleMutation
}

func roleVisibleText(decoded fixedcarrier.TemporalCarrierDecoded) fixedcarrier.VisibleTextManifest {
	labels:=map[string]string{"A":"A[01]","B":"B[02]","C":"C[03]"}
	facts:=[]fixedcarrier.SemanticFact{}
	for _,c:=range decoded.Program.Automaton.Cells{facts=append(facts,fixedcarrier.SemanticFact{Key:"CELL."+c.ID+".LABEL",Value:"CELL "+labels[c.ID]},fixedcarrier.SemanticFact{Key:"CELL."+c.ID+".INITIAL_TEXT",Value:c.InitialState})}
	facts=append(facts,fixedcarrier.SemanticFact{Key:"TEMPORAL_GRAMMAR.SYNC_TEXT",Value:"EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT"})
	for _,r:=range decoded.Program.Automaton.Rules{
		req:="TRUE";if len(r.Requires)>0{parts:=[]string{};for _,p:=range r.Requires{parts=append(parts,labels[p.CellID]+"="+p.State)};req=strings.Join(parts,"&")}
		from:=r.FromState;if from==""{from="*"}
		facts=append(facts,fixedcarrier.SemanticFact{Key:"RULE."+r.ID+".WHEN_TEXT",Value:strings.ToUpper(r.ID)+" WHEN "+req},fixedcarrier.SemanticFact{Key:"RULE."+r.ID+".ROLE_TEXT",Value:"TARGET "+labels[r.TargetCell]+" | REQUIRE "+from+" | SET "+r.ToState})
	}
	facts=append(facts,fixedcarrier.SemanticFact{Key:"EXECUTION_POLICY.TEXT",Value:fixedcarrier.ExecuteVisibleRulesToStableTextR1})
	sort.Slice(facts,func(i,j int)bool{return facts[i].Key<facts[j].Key})
	return fixedcarrier.VisibleTextManifest{Schema:fixedcarrier.VisibleTextManifestSchemaR1,ProgramSHA256:decoded.ProgramSHA256,Facts:facts}
}

func setFact(in []fixedcarrier.SemanticFact,key,value string)[]fixedcarrier.SemanticFact{for i:=range in{if in[i].Key==key{in[i].Value=value;return in}};return append(in,fixedcarrier.SemanticFact{Key:key,Value:value})}
func writeJSON(path string,v any){b,err:=json.MarshalIndent(v,"","  ");must(err);b=append(b,'\n');must(os.WriteFile(path,b,0o644))}
func must(err error){if err!=nil{panic(err)}}
