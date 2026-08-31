package fixedcarrier

import (
	"fmt"
	"sort"
	"strings"
)

const (
	SemanticManifestSchemaR1 = "origami.semantic-manifest.r1"
	InteropBuildManifestSchemaR1 = "origami.candidate-build-manifest.r1"
)

type SemanticFact struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type SemanticManifest struct {
	Schema string `json:"schema"`
	ProgramSHA256 string `json:"program_sha256"`
	PayloadSHA256 string `json:"payload_sha256,omitempty"`
	Facts []SemanticFact `json:"facts"`
}

type InteropMutation struct {
	Kind string `json:"kind"`
	Target string `json:"target"`
	Value string `json:"value"`
}

type InteropBuildManifest struct {
	Schema string `json:"schema"`
	CandidateID string `json:"candidate_id"`
	RendererVersion string `json:"renderer_version"`
	ArtifactSHA256 string `json:"artifact_sha256"`
	ArtifactBytes int `json:"artifact_bytes"`
	ProgramSHA256 string `json:"program_sha256"`
	PayloadSHA256 string `json:"payload_sha256,omitempty"`
	InheritedMutations []InteropMutation `json:"inherited_mutations,omitempty"`
	AppliedMutations []InteropMutation `json:"applied_mutations"`
	VisibleSemantics SemanticManifest `json:"visible_semantics"`
}

func TemporalSemanticManifest(decoded TemporalCarrierDecoded, mutations []CandidateMutation) SemanticManifest {
	facts:=[]SemanticFact{
		{Key:"PROGRAM_ID",Value:decoded.Program.ID},
		{Key:"AUTOMATON_ID",Value:decoded.Program.Automaton.ID},
		{Key:"MAX_STEPS",Value:fmt.Sprintf("%d",decoded.Program.MaxSteps)},
		{Key:"CHECKPOINT_EVERY",Value:fmt.Sprintf("%d",decoded.Program.CheckpointEvery)},
		{Key:"SYNC_SEMANTICS",Value:"ALL_SELECTED_RULES_USE_SAME_PRE_STEP_SNAPSHOT"},
		{Key:"EXECUTION_POLICY",Value:"NONE"},
	}
	for _,c:=range decoded.Program.Automaton.Cells{facts=append(facts,SemanticFact{Key:"CELL."+c.ID+".INITIAL",Value:c.InitialState})}
	for _,r:=range decoded.Program.Automaton.Rules{
		req:=make([]string,0,len(r.Requires));for _,p:=range r.Requires{req=append(req,p.CellID+"="+p.State)};sort.Strings(req)
		from:=r.FromState;if strings.TrimSpace(from)==""{from="*"}
		facts=append(facts,SemanticFact{Key:"RULE."+r.ID+".TARGET",Value:r.TargetCell},SemanticFact{Key:"RULE."+r.ID+".FROM",Value:from},SemanticFact{Key:"RULE."+r.ID+".TO",Value:r.ToState},SemanticFact{Key:"RULE."+r.ID+".REQUIRES",Value:strings.Join(req,"&")})
	}
	for _,e:=range decoded.Program.Automaton.Edges{facts=append(facts,SemanticFact{Key:"EDGE."+e.From+"->"+e.To,Value:e.Kind})}
	for _,m:=range mutations{
		t:=strings.ToUpper(strings.TrimSpace(m.Target));v:=strings.ToUpper(strings.TrimSpace(m.Value))
		if t=="EXECUTION_POLICY" && v=="EXECUTE_VISIBLE_RULES_TO_STABLE_R1"{setFact(&facts,"EXECUTION_POLICY","EXECUTE_VISIBLE_RULES_TO_STABLE_R1")}
		if t=="T2_SEMANTIC_TEMPORAL_SUPERGRAPH" && v=="VISIBLE_RULE_MICROGRAMMAR_R1"{setFact(&facts,"TEMPORAL_GRAMMAR","VISIBLE_RULE_MICROGRAMMAR_R1")}
	}
	sort.Slice(facts,func(i,j int)bool{if facts[i].Key==facts[j].Key{return facts[i].Value<facts[j].Value};return facts[i].Key<facts[j].Key})
	return SemanticManifest{Schema:SemanticManifestSchemaR1,ProgramSHA256:decoded.ProgramSHA256,Facts:facts}
}

func TemporalInteropBuildManifest(report CandidateBuildReport, decoded TemporalCarrierDecoded) InteropBuildManifest {
	return TemporalInteropBuildManifestWithInherited(report,decoded,nil)
}

func TemporalInteropBuildManifestWithInherited(report CandidateBuildReport, decoded TemporalCarrierDecoded, inherited []CandidateMutation) InteropBuildManifest {
	toInterop:=func(in []CandidateMutation)[]InteropMutation{out:=make([]InteropMutation,0,len(in));for _,m:=range in{out=append(out,InteropMutation{Kind:m.Kind,Target:m.Target,Value:m.Value})};return out}
	all:=append([]CandidateMutation(nil),inherited...);all=append(all,report.AppliedMutations...)
	return InteropBuildManifest{Schema:InteropBuildManifestSchemaR1,CandidateID:report.CandidateID,RendererVersion:CandidateBuildReportSchema,ArtifactSHA256:report.CandidateSHA256,ArtifactBytes:report.PNGBytes,ProgramSHA256:report.CandidateProgramSHA256,InheritedMutations:toInterop(inherited),AppliedMutations:toInterop(report.AppliedMutations),VisibleSemantics:TemporalSemanticManifest(decoded,all)}
}

func setFact(facts *[]SemanticFact,key,value string){for i:=range *facts{if (*facts)[i].Key==key{(*facts)[i].Value=value;return}};*facts=append(*facts,SemanticFact{Key:key,Value:value})}
