package memory

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type GraphSignature struct {
	Schema          string   `json:"schema"`
	CarrierID       string   `json:"carrier_id"`
	RootTerms       []string `json:"root_terms,omitempty"`
	SemanticClasses []string `json:"semantic_classes,omitempty"`
	RelationTypes   []string `json:"relation_types,omitempty"`
	Clusters        []string `json:"clusters,omitempty"`
	NodeCount       int      `json:"node_count"`
	MetadataTokens  int      `json:"metadata_tokens"`
}

type Candidate struct {
	Node  Node    `json:"node"`
	Score float64 `json:"score"`
}

type Store struct {
	carrierID string
	byAddress map[string]Node
	byCID     map[string]string
	inverted  map[string]map[string]struct{}
	signature GraphSignature
}

func BuildStore(carrierID string, nodes []Node) (*Store, error) {
	store := &Store{carrierID:carrierID,byAddress:make(map[string]Node,len(nodes)),byCID:make(map[string]string,len(nodes)),inverted:make(map[string]map[string]struct{})}
	for _, raw := range nodes {
		n := raw
		if err := ValidateNode(n); err != nil { return nil, err }
		if _, exists := store.byAddress[n.Address]; exists { return nil, fmt.Errorf("duplicate memory address %q", n.Address) }
		if n.CarrierID == "" { n.CarrierID = carrierID }
		if n.CID == "" { n.CID = n.DerivedCID() }
		store.byAddress[n.Address] = n
		if old, exists := store.byCID[n.CID]; exists && old != n.Address { return nil, fmt.Errorf("duplicate CID %q at %q and %q", n.CID, old, n.Address) }
		store.byCID[n.CID] = n
		store.indexNode(n)
	}
	store.signature = store.buildSignature()
	return store, nil
}

func (s *Store) CarrierID() string { return s.carrierID }
func (s *Store) Signature() GraphSignature { return s.signature }
func (s *Store) Len() int { return len(s.byAddress) }
func (s *Store) Lookup(address string) (Node, bool) { n,ok:=s.byAddress[address];return n,ok }
func (s *Store) LookupCID(cid string) (Node, bool) { address,ok:=s.byCID[cid];if !ok{return Node{},false};return s.Lookup(address) }

func (s *Store) Nodes() []Node {
	addresses:=make([]string,0,len(s.byAddress));for address:=range s.byAddress{addresses=append(addresses,address)};sort.Strings(addresses)
	out:=make([]Node,0,len(addresses));for _,address:=range addresses{out=append(out,s.byAddress[address])};return out
}

func (s *Store) Search(terms []string, limit int) []Candidate {
	if limit<=0{limit=32};scores:=map[string]float64{}
	for _,term:=range normalizeTerms(terms){for address:=range s.inverted[term]{scores[address]+=1.0}}
	out:=make([]Candidate,0,len(scores));for address,score:=range scores{n:=s.byAddress[address];score+=clamp01(n.Importance.Global)*0.35+clamp01(n.Importance.Local)*0.20+clamp01(n.Importance.Confidence)*0.15;out=append(out,Candidate{Node:n,Score:score})}
	sort.Slice(out,func(i,j int)bool{if out[i].Score==out[j].Score{return out[i].Node.Address<out[j].Node.Address};return out[i].Score>out[j].Score});if len(out)>limit{out=out[:limit]};return out
}

func (s *Store) indexNode(n Node) {
	var b strings.Builder;fmt.Fprintf(&b,"%s %s %s %s ",n.Label,n.Kind,n.SemanticClass,n.ClusterID)
	for _,fidelity:=range []Fidelity{FidelityAbstract,FidelitySummary}{if level,ok:=n.Level(fidelity);ok{b.WriteString(level.Text);b.WriteByte(' ')}}
	if n.Label==""&&len(n.Levels)==0&&n.Value!=""{v:=n.Value;if len(v)>512{v=v[:512]};b.WriteString(v)}
	for _,term:=range tokenize(b.String()){bucket:=s.inverted[term];if bucket==nil{bucket=map[string]struct{}{};s.inverted[term]=bucket};bucket[n.Address]=struct{}{}}
}

func (s *Store) buildSignature() GraphSignature {
	roots,classes,relations,clusters:=map[string]struct{}{},map[string]struct{}{},map[string]struct{}{},map[string]struct{}{};metadata:=0
	for _,n:=range s.byAddress{
		if n.Parent==""{for _,t:=range tokenize(n.Label+" "+n.SemanticClass){roots[t]=struct{}{}}}
		if n.SemanticClass!=""{classes[n.SemanticClass]=struct{}{}};if n.ClusterID!=""{clusters[n.ClusterID]=struct{}{}}
		metadata+=EstimateTokens(n.Label+" "+n.Kind+" "+n.SemanticClass+" "+n.ClusterID)
		for relation:=range n.Links{relations[relation]=struct{}{}};for _,edge:=range n.Edges{relations[edge.Relation]=struct{}{}}
	}
	return GraphSignature{Schema:SchemaR0+".graph-signature",CarrierID:s.carrierID,RootTerms:sortedKeys(roots),SemanticClasses:sortedKeys(classes),RelationTypes:sortedKeys(relations),Clusters:sortedKeys(clusters),NodeCount:len(s.byAddress),MetadataTokens:metadata}
}

func tokenize(text string) []string { return normalizeTerms(strings.FieldsFunc(strings.ToLower(text),func(r rune)bool{return !(unicode.IsLetter(r)||unicode.IsDigit(r)||r=='_'||r=='-')})) }
func normalizeTerms(terms []string) []string { seen:=map[string]struct{}{};out:=make([]string,0,len(terms));for _,raw:=range terms{for _,term:=range strings.Fields(strings.ToLower(strings.TrimSpace(raw))){if len(term)<2||isStopWord(term){continue};if _,ok:=seen[term];ok{continue};seen[term]=struct{}{};out=append(out,term)}};sort.Strings(out);return out }
func isStopWord(term string) bool { switch term{case "the","and","for","with","from","that","this","what","which","why","how","una","uno","unos","unas","que","qué","como","cómo","por","para","con","del","las","los","esto","esta","este":return true;default:return false} }
func sortedKeys(m map[string]struct{}) []string { out:=make([]string,0,len(m));for k:=range m{out=append(out,k)};sort.Strings(out);return out }
func clamp01(v float64) float64 { if v<0{return 0};if v>1{return 1};return v }
