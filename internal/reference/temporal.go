package reference

import "fmt"

type TemporalRelation string

const (
	TemporalSequence TemporalRelation = "SEQUENCE"
	TemporalChange TemporalRelation = "CHANGE"
)

type TemporalPredicate struct {
	Entity string `json:"entity"`
	Values []string `json:"values"`
	Relation TemporalRelation `json:"relation"`
}

type TemporalObservationContract struct {
	ID string `json:"id"`
	Authority string `json:"authority"`
	Predicate TemporalPredicate `json:"predicate"`
}

type TemporalObservationResult struct {
	ContractID string `json:"contract_id"`
	Authority string `json:"authority"`
	Matched bool `json:"matched"`
	Entity string `json:"entity"`
	Relation TemporalRelation `json:"relation"`
	ObservedValues []string `json:"observed_values"`
	TraceHash string `json:"trace_sha256"`
}

func trajectoryValues(trace []CoreState,entity string)[]string{out:=make([]string,0,len(trace));for _,state:=range trace{v,ok:=state[entity];if !ok{out=append(out,"<ABSENT>");continue};if v.Value!=""{out=append(out,v.Value)}else{out=append(out,string(v.Status))}};return out}

func containsContiguousSequence(have,want []string)bool{if len(want)==0||len(want)>len(have){return false};for i:=0;i<=len(have)-len(want);i++{ok:=true;for j:=range want{if have[i+j]!=want[j]{ok=false;break}};if ok{return true}};return false}

// ObserveTemporal derives a property from ordered transition evidence. The
// property is not written back into CoreState and therefore does not become a
// hidden state transition or a stored answer.
func ObserveTemporal(run CoreRun,contract TemporalObservationContract)(TemporalObservationResult,error){
	if contract.ID==""||contract.Authority==""{return TemporalObservationResult{},fmt.Errorf("temporal observation requires id and authority")};if len(run.Trace)==0{return TemporalObservationResult{},fmt.Errorf("run has no trace")}
	values:=trajectoryValues(run.Trace,contract.Predicate.Entity);matched:=false
	switch contract.Predicate.Relation{
	case TemporalSequence: matched=containsContiguousSequence(values,contract.Predicate.Values)
	case TemporalChange: for i:=1;i<len(values);i++{if values[i]!=values[i-1]{matched=true;break}}
	default:return TemporalObservationResult{},fmt.Errorf("unsupported temporal relation %q",contract.Predicate.Relation)
	}
	h,err:=TraceHash(run.Trace);if err!=nil{return TemporalObservationResult{},err};return TemporalObservationResult{ContractID:contract.ID,Authority:contract.Authority,Matched:matched,Entity:contract.Predicate.Entity,Relation:contract.Predicate.Relation,ObservedValues:values,TraceHash:h},nil
}
