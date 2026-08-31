package fixedcarrier

import (
	"fmt"
	"sort"
	"strings"
)

const (
	SemanticManifestSchemaR1 = "origami.semantic-manifest.r1"
	VisibleTextManifestSchemaR1 = "origami.visible-text-manifest.r1"
	InteropBuildManifestSchemaR1 = "origami.candidate-build-manifest.r1"
)

type SemanticFact struct { Key string `json:"key"`; Value string `json:"value"` }
type SemanticManifest struct { Schema string `json:"schema"`; ProgramSHA256 string `json:"program_sha256"`; PayloadSHA256 string `json:"payload_sha256,omitempty"`; Facts []SemanticFact `json:"facts"` }
type VisibleTextManifest struct { Schema string `json:"schema"`; ProgramSHA256 string `json:"program_sha256"`; Facts []SemanticFact `json:"facts"` }
type InteropMutation struct { Kind string `json:"kind"`; Target string `json:"target"`; Value string `json:"value"` }
type InteropBuildManifest struct { Schema string `json:"schema"`; CandidateID string `json:"candidate_id"`; RendererVersion string `json:"renderer_version"`; ArtifactSHA256 string `json:"artifact_sha256"`; ArtifactBytes int `json:"artifact_bytes"`; ProgramSHA256 string `json:"program_sha256"`; PayloadSHA256 string `json:"payload_sha256,omitempty"`; InheritedMutations []InteropMutation `json:"inherited_mutations,omitempty"`; AppliedMutations []InteropMutation `json:"applied_mutations"`; VisibleSemantics SemanticManifest `json:"visible_semantics"`; VisibleText VisibleTextManifest `json:"visible_text"` }

func TemporalSemanticManifest(decoded TemporalCarrierDecoded, mutations []CandidateMutation) SemanticManifest {
	facts:=[]SemanticFact{{Key:"PROGRAM_ID",Value:decoded.Program.ID},{Key:"AUTOMATON_ID",Value:decoded.Program.Automaton.ID},{Key:"MAX_STEPS",Value:fmt.Sprintf("%d",decoded.Program.MaxSteps)},{Key:"CHECKPOINT_EVERY",Value:fmt.Sprintf("%d",decoded.Program.CheckpointEvery)},{Key:"SYNC_SEMANTICS",Value:"ALL_SELECTED_RULES_USE_SAME_PRE_STEP_SNAPSHOT"},{Key:"EXECUTION_POLICY",Value:"NONE"}}
	for _,c:=range decoded.Program.Automaton.Cells{facts=append(facts,SemanticFact{Key:"CELL."+c.ID+".INITIAL",Value:c.InitialState},SemanticFact{Key:"VISIBLE_CELL_ID_"+c.ID,Value:c.ID})}
	for _,r:=range decoded.Program.Automaton.Rules{
		req:=make([]string,0,len(r.Requires));for _,p:=range r.Requires{req=append(req,p.CellID+"="+p.State)};sort.Strings(req)
		from:=r.FromState;if strings.TrimSpace(from)==""{from="*"}
		facts=append(facts,SemanticFact{Key:"RULE."+r.ID+".TARGET",Value:r.TargetCell},SemanticFact{Key:"RULE."+r.ID+".FROM",Value:from},SemanticFact{Key:"RULE."+r.ID+".TO",Value:r.ToState},SemanticFact{Key:"RULE."+r.ID+".REQUIRES",Value:strings.Join(req,"&")})
	}
	for _,e:=range decoded.Program.Automaton.Edges{facts=append(facts,SemanticFact{Key:"EDGE."+e.From+"->"+e.To,Value:e.Kind})}
	for _,m:=range mutations{
		t:=strings.ToUpper(strings.TrimSpace(m.Target));v:=strings.ToUpper(strings.TrimSpace(m.Value))
		if t=="EXECUTION_POLICY"&&v=="EXECUTE_VISIBLE_RULES_TO_STABLE_R1"{setFact(&facts,"EXECUTION_POLICY","EXECUTE_VISIBLE_RULES_TO_STABLE_R1")}
		if t=="T2_SEMANTIC_TEMPORAL_SUPERGRAPH"&&v=="VISIBLE_RULE_MICROGRAMMAR_R1"{setFact(&facts,"TEMPORAL_GRAMMAR","VISIBLE_RULE_MICROGRAMMAR_R1")}
		if t=="CELL_IDENTITY_ENCODING"&&v==VisibleCellIdentityRedundancyR1{labels:=temporalVisibleCellLabels(decoded,true);for _,c:=range decoded.Program.Automaton.Cells{setFact(&facts,"VISIBLE_CELL_ID_"+c.ID,labels[c.ID])}}
		if t=="FROM_STATE_PRECONDITION_VISIBILITY"&&v==VisibleFromStatePreconditionR1{setFact(&facts,"FROM_STATE_PRECONDITION_VISIBILITY",VisibleFromStatePreconditionR1)}
		if t=="RULE_ROLE_BINDING"&&v==VisibleRuleRoleBindingR1{setFact(&facts,"RULE_ROLE_BINDING",VisibleRuleRoleBindingR1)}
		if t=="EXECUTION_POLICY_COMPLIANCE"&&v==ExecuteDontSummarizeToStableR1{setFact(&facts,"EXECUTION_POLICY_COMPLIANCE",ExecuteDontSummarizeToStableR1)}
		if t=="SYNCHRONOUS_EXECUTION_FIDELITY"&&v==FreezeSelectApplyTogetherR1{setFact(&facts,"SYNCHRONOUS_EXECUTION_FIDELITY",FreezeSelectApplyTogetherR1)}
	}
	sort.Slice(facts,func(i,j int)bool{if facts[i].Key==facts[j].Key{return facts[i].Value<facts[j].Value};return facts[i].Key<facts[j].Key})
	return SemanticManifest{Schema:SemanticManifestSchemaR1,ProgramSHA256:decoded.ProgramSHA256,Facts:facts}
}

func TemporalVisibleTextManifest(decoded TemporalCarrierDecoded, mutations []CandidateMutation) VisibleTextManifest {
	redundant:=hasCandidateMutation(mutations,"CELL_IDENTITY_ENCODING",VisibleCellIdentityRedundancyR1)
	grammar:=hasCandidateMutation(mutations,"T2_SEMANTIC_TEMPORAL_SUPERGRAPH","VISIBLE_RULE_MICROGRAMMAR_R1")
	explicitFrom:=hasCandidateMutation(mutations,"FROM_STATE_PRECONDITION_VISIBILITY",VisibleFromStatePreconditionR1)
	roleBinding:=hasCandidateMutation(mutations,"RULE_ROLE_BINDING",VisibleRuleRoleBindingR1)
	execute:=hasCandidateMutation(mutations,"EXECUTION_POLICY","EXECUTE_VISIBLE_RULES_TO_STABLE_R1")
	executionCompliance:=hasCandidateMutation(mutations,"EXECUTION_POLICY_COMPLIANCE",ExecuteDontSummarizeToStableR1)
	syncFidelity:=hasCandidateMutation(mutations,"SYNCHRONOUS_EXECUTION_FIDELITY",FreezeSelectApplyTogetherR1)
	labels:=temporalVisibleCellLabels(decoded,redundant)
	facts:=[]SemanticFact{}
	for _,c:=range decoded.Program.Automaton.Cells{facts=append(facts,SemanticFact{Key:"CELL."+c.ID+".LABEL",Value:"CELL "+labels[c.ID]},SemanticFact{Key:"CELL."+c.ID+".INITIAL_TEXT",Value:c.InitialState})}

	if roleBinding {
		facts=append(facts,SemanticFact{Key:"TEMPORAL_GRAMMAR.SYNC_TEXT",Value:"EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT"})
		if len(decoded.Program.Automaton.Cells)<=4{
			for i,r:=range decoded.Program.Automaton.Rules{
				if i>=4{break}
				req:="TRUE";if len(r.Requires)>0{parts:=make([]string,0,len(r.Requires));for _,p:=range r.Requires{parts=append(parts,shortLabel(labels[p.CellID],16)+"="+shortLabel(p.State,10))};req=strings.Join(parts,"&")}
				from:=r.FromState;if strings.TrimSpace(from)==""{from="*"}
				facts=append(facts,SemanticFact{Key:"RULE."+r.ID+".WHEN_TEXT",Value:strings.ToUpper(r.ID)+" WHEN "+req},SemanticFact{Key:"RULE."+r.ID+".ROLE_TEXT",Value:"TARGET "+shortLabel(labels[r.TargetCell],16)+" | REQUIRE "+shortLabel(from,10)+" | SET "+shortLabel(r.ToState,10)})
			}
		}
	} else if explicitFrom {
		facts=append(facts,SemanticFact{Key:"TEMPORAL_GRAMMAR.SYNC_TEXT",Value:"EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT"})
		if len(decoded.Program.Automaton.Cells)<=4{for i,r:=range decoded.Program.Automaton.Rules{if i>=6{break};req:="TRUE";if len(r.Requires)>0{parts:=make([]string,0,len(r.Requires));for _,p:=range r.Requires{parts=append(parts,shortLabel(labels[p.CellID],16)+"="+shortLabel(p.State,10))};req=strings.Join(parts,"&")};from:=r.FromState;if strings.TrimSpace(from)==""{from="*"};line:=fmt.Sprintf("IF %s AND %s=%s THEN %s -> %s",req,shortLabel(labels[r.TargetCell],16),shortLabel(from,10),shortLabel(labels[r.TargetCell],16),shortLabel(r.ToState,10));facts=append(facts,SemanticFact{Key:"RULE."+r.ID+".TEXT",Value:shortLabel(line,94)})}}
	} else if grammar {
		facts=append(facts,SemanticFact{Key:"TEMPORAL_GRAMMAR.SYNC_TEXT",Value:"EACH STEP: TEST ALL RULES ON SAME PRE-STEP SNAPSHOT"})
		if len(decoded.Program.Automaton.Cells)<=4{for i,r:=range decoded.Program.Automaton.Rules{if i>=6{break};req:="TRUE";if len(r.Requires)>0{parts:=make([]string,0,len(r.Requires));for _,p:=range r.Requires{parts=append(parts,shortLabel(labels[p.CellID],16)+"="+shortLabel(p.State,10))};req=strings.Join(parts,"&")};from:=r.FromState;if strings.TrimSpace(from)==""{from="*"};line:=fmt.Sprintf("IF %s => %s:%s>%s",req,shortLabel(labels[r.TargetCell],16),shortLabel(from,10),shortLabel(r.ToState,10));facts=append(facts,SemanticFact{Key:"RULE."+r.ID+".TEXT",Value:shortLabel(line,94)})}}
	}
	if execute{facts=append(facts,SemanticFact{Key:"EXECUTION_POLICY.TEXT",Value:ExecuteVisibleRulesToStableTextR1})}
	if executionCompliance{
		facts=append(facts,SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.MODE_TEXT",Value:ExecutionComplianceModeTextR1})
		if !syncFidelity {
			facts=append(facts,
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP1_TEXT",Value:ExecutionComplianceStep1TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP2_TEXT",Value:ExecutionComplianceStep2TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP3_TEXT",Value:ExecutionComplianceStep3TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP4_TEXT",Value:ExecutionComplianceStep4TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP5_TEXT",Value:ExecutionComplianceStep5TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STEP6_TEXT",Value:ExecutionComplianceStep6TextR1},
				SemanticFact{Key:"EXECUTION_POLICY_COMPLIANCE.STOP_TEXT",Value:ExecutionComplianceStopTextR1},
			)
		}
	}
	if syncFidelity {
		facts=append(facts,
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.HEADER_TEXT",Value:SyncFidelityHeaderTextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STEP1_TEXT",Value:SyncFidelityStep1TextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STEP2_TEXT",Value:SyncFidelityStep2TextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STEP2B_TEXT",Value:SyncFidelityStep2BTextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STEP3_TEXT",Value:SyncFidelityStep3TextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STEP4_TEXT",Value:SyncFidelityStep4TextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.NO_ORDER_TEXT",Value:SyncFidelityNoOrderTextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.NO_CASCADE_TEXT",Value:SyncFidelityNoCascadeTextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.NEXT_TEXT",Value:SyncFidelityNextTextR1},
			SemanticFact{Key:"SYNCHRONOUS_EXECUTION_FIDELITY.STOP_TEXT",Value:SyncFidelityStopTextR1},
		)
	}
	sort.Slice(facts,func(i,j int)bool{return facts[i].Key<facts[j].Key})
	return VisibleTextManifest{Schema:VisibleTextManifestSchemaR1,ProgramSHA256:decoded.ProgramSHA256,Facts:facts}
}

func TemporalInteropBuildManifest(report CandidateBuildReport, decoded TemporalCarrierDecoded) InteropBuildManifest { return TemporalInteropBuildManifestWithInherited(report,decoded,nil) }
func TemporalInteropBuildManifestWithInherited(report CandidateBuildReport, decoded TemporalCarrierDecoded, inherited []CandidateMutation) InteropBuildManifest { toInterop:=func(in []CandidateMutation)[]InteropMutation{out:=make([]InteropMutation,0,len(in));for _,m:=range in{out=append(out,InteropMutation{Kind:m.Kind,Target:m.Target,Value:m.Value})};return out};all:=append([]CandidateMutation(nil),inherited...);all=append(all,report.AppliedMutations...);return InteropBuildManifest{Schema:InteropBuildManifestSchemaR1,CandidateID:report.CandidateID,RendererVersion:CandidateBuildReportSchema,ArtifactSHA256:report.CandidateSHA256,ArtifactBytes:report.PNGBytes,ProgramSHA256:report.CandidateProgramSHA256,InheritedMutations:toInterop(inherited),AppliedMutations:toInterop(report.AppliedMutations),VisibleSemantics:TemporalSemanticManifest(decoded,all),VisibleText:TemporalVisibleTextManifest(decoded,all)} }
func hasCandidateMutation(mutations []CandidateMutation,target,value string)bool{for _,m:=range mutations{if strings.EqualFold(strings.TrimSpace(m.Target),target)&&strings.EqualFold(strings.TrimSpace(m.Value),value){return true}};return false}
func setFact(facts *[]SemanticFact,key,value string){for i:=range *facts{if (*facts)[i].Key==key{(*facts)[i].Value=value;return}};*facts=append(*facts,SemanticFact{Key:key,Value:value})}
