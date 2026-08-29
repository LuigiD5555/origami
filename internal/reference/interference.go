package reference

import (
	"fmt"
	"sort"
)

type Contribution struct {
	ID string `json:"id"`
	Entity string `json:"entity"`
	Value string `json:"value"`
	Weight int `json:"weight"`
	SourceRule string `json:"source_rule"`
}

type ContributionEvidence struct {
	Entity string `json:"entity"`
	Contributions []Contribution `json:"contributions"`
	Net int `json:"net"`
	Result SemanticValue `json:"result"`
}

// ResolveContributions combines independent signed contributions. It is
// deliberately order-independent: declaration order cannot become hidden
// priority semantics. Net zero with non-zero opposing contributions is
// CANCELLED, not UNKNOWN and not ABSENT.
func ResolveContributions(entity string, contributions []Contribution)(SemanticValue,ContributionEvidence,error){
	if entity==""{return SemanticValue{},ContributionEvidence{},fmt.Errorf("entity must not be empty")}
	if len(contributions)==0{return SemanticValue{Status:StatusUnknown},ContributionEvidence{Entity:entity,Result:SemanticValue{Status:StatusUnknown}},nil}
	cp:=append([]Contribution(nil),contributions...);sort.Slice(cp,func(i,j int)bool{return cp[i].ID<cp[j].ID})
	net:=0;positive:=false;negative:=false
	for _,c:=range cp{if c.ID==""||c.Entity!=entity||c.Weight==0{return SemanticValue{},ContributionEvidence{},fmt.Errorf("invalid contribution for %s",entity)};net+=c.Weight;if c.Weight>0{positive=true};if c.Weight<0{negative=true}}
	var result SemanticValue
	switch {
	case net==0&&positive&&negative: result=SemanticValue{Status:StatusCancelled,Reason:"opposing contributions cancel"}
	case net>0: result=SemanticValue{Status:StatusPresent,Value:"positive",Reason:"net positive contribution"}
	case net<0: result=SemanticValue{Status:StatusInhibited,Reason:"net negative contribution"}
	default: result=SemanticValue{Status:StatusUnknown,Reason:"no resolvable contribution"}
	}
	return result,ContributionEvidence{Entity:entity,Contributions:cp,Net:net,Result:result},nil
}
