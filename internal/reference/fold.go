package reference

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// AlternativeSet preserves unresolved alternatives as alternatives. Fold is
// forbidden from selecting a member merely to reduce representation size.
type AlternativeSet struct {
	ID      string          `json:"id"`
	Members []SemanticValue `json:"members"`
}

// FoldSource is carrier-neutral semantic material. It deliberately contains no
// image, glyph, byte-stream or OHF-specific fields.
type FoldSource struct {
	State        CoreState        `json:"state"`
	Alternatives []AlternativeSet `json:"alternatives,omitempty"`
	Rules        []Rule           `json:"rules,omitempty"`
}

type FoldedCore struct {
	Schema       string           `json:"schema"`
	State        CoreState        `json:"state"`
	Alternatives []AlternativeSet `json:"alternatives,omitempty"`
	Rules        []Rule           `json:"rules,omitempty"`
	Commitment   string           `json:"commitment_sha256"`
}

func canonicalFoldSource(src FoldSource) FoldSource {
	out:=FoldSource{State:cloneCoreState(src.State),Alternatives:append([]AlternativeSet(nil),src.Alternatives...),Rules:append([]Rule(nil),src.Rules...)}
	for i:=range out.Alternatives{out.Alternatives[i].Members=append([]SemanticValue(nil),out.Alternatives[i].Members...)}
	sort.Slice(out.Alternatives,func(i,j int)bool{return out.Alternatives[i].ID<out.Alternatives[j].ID})
	sort.Slice(out.Rules,func(i,j int)bool{return out.Rules[i].ID<out.Rules[j].ID})
	return out
}

func FoldCommitment(src FoldSource)(string,error){b,err:=json.Marshal(canonicalFoldSource(src));if err!=nil{return "",err};h:=sha256.Sum256(b);return hex.EncodeToString(h[:]),nil}

func Fold(src FoldSource)(FoldedCore,error){
	for _,a:=range src.Alternatives{if a.ID==""{return FoldedCore{},fmt.Errorf("alternative set id must not be empty")};if len(a.Members)<2{return FoldedCore{},fmt.Errorf("alternative set %s must preserve at least two unresolved members",a.ID)}}
	canonical:=canonicalFoldSource(src);commitment,err:=FoldCommitment(canonical);if err!=nil{return FoldedCore{},err}
	return FoldedCore{Schema:"origami.folded-core.r0",State:canonical.State,Alternatives:canonical.Alternatives,Rules:canonical.Rules,Commitment:commitment},nil
}

func Unfold(f FoldedCore)(FoldSource,error){
	if f.Schema!="origami.folded-core.r0"{return FoldSource{},fmt.Errorf("unsupported folded schema %q",f.Schema)}
	src:=canonicalFoldSource(FoldSource{State:f.State,Alternatives:f.Alternatives,Rules:f.Rules});commitment,err:=FoldCommitment(src);if err!=nil{return FoldSource{},err};if commitment!=f.Commitment{return FoldSource{},fmt.Errorf("fold commitment mismatch")};return src,nil
}
