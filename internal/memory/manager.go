package memory

import (
	"fmt"
	"sort"
	"strings"
)

type ContextBudget struct {
	MaxTokens       int `json:"max_tokens"`
	ReserveTokens   int `json:"reserve_tokens,omitempty"`
	UsedTokens      int `json:"used_tokens"`
	RemainingTokens int `json:"remaining_tokens"`
}

type QueryPlan struct {
	Raw               string   `json:"raw"`
	Terms             []string `json:"terms"`
	Relations         []string `json:"relations,omitempty"`
	NeedEvidence      bool     `json:"need_evidence,omitempty"`
	NeedContradiction bool     `json:"need_contradiction,omitempty"`
	NeedExact         bool     `json:"need_exact,omitempty"`
	PreferredFidelity Fidelity `json:"preferred_fidelity,omitempty"`
	MaxGraphDepth     int      `json:"max_graph_depth"`
}

type ContextItem struct {
	Address    string   `json:"address"`
	CID        string   `json:"cid"`
	CarrierID  string   `json:"carrier_id,omitempty"`
	ClusterID  string   `json:"cluster_id,omitempty"`
	Kind       string   `json:"kind,omitempty"`
	Label      string   `json:"label,omitempty"`
	Fidelity   Fidelity `json:"fidelity"`
	Content    string   `json:"content"`
	Score      float64  `json:"score"`
	TokenCost  int      `json:"token_cost"`
	SourceRef  string   `json:"source_ref,omitempty"`
	Verified   bool     `json:"verified,omitempty"`
}

type ContextRelation struct {
	From       string  `json:"from"`
	Relation   string  `json:"relation"`
	To         string  `json:"to"`
	Weight     float64 `json:"weight,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	ProofRef   string  `json:"proof_ref,omitempty"`
}

type ContextMetrics struct {
	Candidates       int `json:"candidates"`
	Touched          int `json:"touched"`
	Unfolded         int `json:"unfolded"`
	GraphHops        int `json:"graph_hops"`
	Deduplicated     int `json:"deduplicated"`
	IrrelevantSkipped int `json:"irrelevant_skipped"`
}

type ContextPacket struct {
	Schema         string            `json:"schema"`
	Query          string            `json:"query"`
	Plan           QueryPlan         `json:"plan"`
	Items          []ContextItem     `json:"items,omitempty"`
	Relations      []ContextRelation `json:"relations,omitempty"`
	Contradictions []ContextRelation `json:"contradictions,omitempty"`
	Uncertainties  []string          `json:"uncertainties,omitempty"`
	Provenance     []Provenance      `json:"provenance,omitempty"`
	ExpandableRefs []string          `json:"expandable_refs,omitempty"`
	Budget         ContextBudget     `json:"budget"`
	Metrics        ContextMetrics    `json:"metrics"`
	Unknown        bool              `json:"unknown"`
	Reason         string            `json:"reason,omitempty"`
}

type Manager struct { store *Store }

func NewManager(store *Store) *Manager { return &Manager{store: store} }

func Analyze(raw string) QueryPlan {
	lower := strings.ToLower(raw)
	plan := QueryPlan{Raw: raw, Terms: tokenize(raw), MaxGraphDepth: 2}
	containsAny := func(words ...string) bool {
		for _, word := range words { if strings.Contains(lower, word) { return true } }
		return false
	}
	if containsAny("por qué", "porque", "causa", "causal", "why", "cause") { plan.Relations = append(plan.Relations, "causes", "depends") }
	if containsAny("evidencia", "fuente", "prueba", "evidence", "source", "proof") {
		plan.NeedEvidence = true
		plan.Relations = append(plan.Relations, "evidence", "supports", "source")
	}
	if containsAny("contradic", "conflict", "disagree", "contradict") {
		plan.NeedContradiction = true
		plan.Relations = append(plan.Relations, "contradicts")
	}
	if containsAny("exacto", "exacta", "literal", "bytes", "verbatim", "exact") {
		plan.NeedExact = true
		plan.PreferredFidelity = FidelityExact
	}
	plan.Relations = uniqueStrings(plan.Relations)
	return plan
}

func (m *Manager) Query(raw string, budget ContextBudget) ContextPacket {
	plan := Analyze(raw)
	budget = normalizeBudget(budget)
	packet := ContextPacket{Schema: SchemaR0 + ".context-packet", Query: raw, Plan: plan, Budget: budget}
	if m == nil || m.store == nil {
		packet.Unknown = true; packet.Reason = "memory store unavailable"; return packet
	}
	candidates := m.store.Search(plan.Terms, 64)
	packet.Metrics.Candidates = len(candidates)
	if len(candidates) == 0 {
		packet.Unknown = true; packet.Reason = "no metadata candidates"; return packet
	}

	// Expand only declared graph edges from already relevant candidates. There is
	// no fallback full-memory scan when a relation cannot be followed.
	expanded := append([]Candidate(nil), candidates...)
	if len(plan.Relations) > 0 && plan.MaxGraphDepth > 0 {
		expanded = m.expandGraph(expanded, plan.Relations, plan.MaxGraphDepth, &packet.Metrics)
	}
	packet.Metrics.Candidates = len(expanded)

	preferred := defaultFidelity(plan)
	type ranked struct { candidate Candidate; fidelity Fidelity; level ContentLevel; cost int; utility float64 }
	rankedItems := make([]ranked, 0, len(expanded))
	for _, candidate := range expanded {
		packet.Metrics.Touched++
		f, level, cost, ok := selectLevel(candidate.Node, preferred, packet.Budget.RemainingTokens)
		if !ok { packet.Metrics.IrrelevantSkipped++; continue }
		utility := candidate.Score + clamp01(candidate.Node.Importance.Confidence)*0.25
		if cost > 0 { utility = utility / float64(cost) }
		rankedItems = append(rankedItems, ranked{candidate: candidate, fidelity: f, level: level, cost: cost, utility: utility})
	}
	sort.Slice(rankedItems, func(i, j int) bool {
		if rankedItems[i].utility == rankedItems[j].utility { return rankedItems[i].candidate.Node.Address < rankedItems[j].candidate.Node.Address }
		return rankedItems[i].utility > rankedItems[j].utility
	})

	seenCID := map[string]struct{}{}
	seenAddr := map[string]struct{}{}
	selected := map[string]struct{}{}
	for _, item := range rankedItems {
		n := item.candidate.Node
		cid := n.DerivedCID()
		if _, ok := seenCID[cid]; ok { packet.Metrics.Deduplicated++; continue }
		if _, ok := seenAddr[n.Address]; ok { packet.Metrics.Deduplicated++; continue }
		available := packet.Budget.RemainingTokens
		f, level, cost, ok := selectLevel(n, preferred, available)
		if !ok || cost <= 0 || cost > available { continue }
		packet.Items = append(packet.Items, ContextItem{
			Address: n.Address, CID: cid, CarrierID: n.CarrierID, ClusterID: n.ClusterID,
			Kind: n.Kind, Label: n.Label, Fidelity: f, Content: level.Text,
			Score: item.candidate.Score, TokenCost: cost, SourceRef: firstNonEmpty(level.SourceRef, n.Provenance.SourceRef), Verified: level.Verified,
		})
		packet.Budget.UsedTokens += cost
		packet.Budget.RemainingTokens -= cost
		packet.Metrics.Unfolded++
		seenCID[cid] = struct{}{}; seenAddr[n.Address] = struct{}{}; selected[n.Address] = struct{}{}
		if hasDeeper(n, f) { packet.ExpandableRefs = append(packet.ExpandableRefs, n.Address) }
		if n.Provenance.SourceID != "" || n.Provenance.SourceRef != "" || n.Provenance.SourceHash != "" { packet.Provenance = append(packet.Provenance, n.Provenance) }
	}
	packet.ExpandableRefs = uniqueStrings(packet.ExpandableRefs)
	packet.Provenance = uniqueProvenance(packet.Provenance)
	packet.Relations, packet.Contradictions = m.relationsFor(selected)
	if len(packet.Items) == 0 { packet.Unknown = true; packet.Reason = "no candidate fits active context budget" }
	return packet
}

func (m *Manager) Expand(address string, wanted Fidelity, budget ContextBudget) ContextPacket {
	budget = normalizeBudget(budget)
	plan := QueryPlan{Raw: "EXPAND " + address, PreferredFidelity: wanted}
	packet := ContextPacket{Schema: SchemaR0 + ".context-packet", Query: plan.Raw, Plan: plan, Budget: budget}
	if m == nil || m.store == nil { packet.Unknown = true; packet.Reason = "memory store unavailable"; return packet }
	n, ok := m.store.Lookup(address); packet.Metrics.Touched = 1
	if !ok { packet.Unknown = true; packet.Reason = "address not found"; return packet }
	if wanted == "" { wanted = FidelityDetail }
	f, level, cost, ok := selectLevel(n, wanted, packet.Budget.RemainingTokens)
	if !ok || cost <= 0 || cost > packet.Budget.RemainingTokens { packet.Unknown = true; packet.Reason = "no fidelity fits active context budget"; return packet }
	packet.Items = []ContextItem{{Address:n.Address, CID:n.DerivedCID(), CarrierID:n.CarrierID, ClusterID:n.ClusterID, Kind:n.Kind, Label:n.Label, Fidelity:f, Content:level.Text, TokenCost:cost, SourceRef:firstNonEmpty(level.SourceRef,n.Provenance.SourceRef), Verified:level.Verified}}
	packet.Budget.UsedTokens += cost; packet.Budget.RemainingTokens -= cost; packet.Metrics.Unfolded = 1
	if hasDeeper(n, f) { packet.ExpandableRefs = []string{n.Address} }
	if n.Provenance.SourceID != "" || n.Provenance.SourceRef != "" || n.Provenance.SourceHash != "" { packet.Provenance = []Provenance{n.Provenance} }
	return packet
}

func (m *Manager) expandGraph(seed []Candidate, relations []string, depth int, metrics *ContextMetrics) []Candidate {
	wanted := map[string]struct{}{}; for _, r := range relations { wanted[r] = struct{}{} }
	byAddress := map[string]Candidate{}; frontier := make([]Candidate, 0, len(seed))
	for _, c := range seed { byAddress[c.Node.Address] = c; frontier = append(frontier, c) }
	for level := 0; level < depth; level++ {
		var next []Candidate
		for _, c := range frontier {
			targets := map[string]float64{}
			for relation, addresses := range c.Node.Links {
				if _, ok := wanted[relation]; !ok { continue }
				for _, address := range addresses { targets[address] = 1.0 }
			}
			for _, edge := range c.Node.Edges {
				if _, ok := wanted[edge.Relation]; !ok { continue }
				weight := edge.Weight; if weight <= 0 { weight = 1 }
				targets[edge.Target] = weight
			}
			for address, weight := range targets {
				metrics.Touched++; metrics.GraphHops++
				if _, exists := byAddress[address]; exists { continue }
				n, ok := m.store.Lookup(address); if !ok { continue }
				candidate := Candidate{Node:n, Score:c.Score*0.70 + weight*0.30}
				byAddress[address] = candidate; next = append(next, candidate)
			}
		}
		if len(next) == 0 { break }
		frontier = next
	}
	out := make([]Candidate, 0, len(byAddress)); for _, c := range byAddress { out = append(out, c) }
	sort.Slice(out, func(i,j int)bool{if out[i].Score==out[j].Score{return out[i].Node.Address<out[j].Node.Address};return out[i].Score>out[j].Score})
	return out
}

func (m *Manager) relationsFor(selected map[string]struct{}) ([]ContextRelation, []ContextRelation) {
	var relations, contradictions []ContextRelation
	for address := range selected {
		n, ok := m.store.Lookup(address); if !ok { continue }
		for relation, targets := range n.Links {
			for _, target := range targets { if _, ok := selected[target]; ok { r:=ContextRelation{From:address,Relation:relation,To:target}; if relation=="contradicts"{contradictions=append(contradictions,r)}else{relations=append(relations,r)} } }
		}
		for _, edge := range n.Edges { if _, ok := selected[edge.Target]; ok { r:=ContextRelation{From:address,Relation:edge.Relation,To:edge.Target,Weight:edge.Weight,Confidence:edge.Confidence,ProofRef:edge.ProofRef}; if edge.Relation=="contradicts"{contradictions=append(contradictions,r)}else{relations=append(relations,r)} } }
	}
	sort.Slice(relations,func(i,j int)bool{return relationKey(relations[i])<relationKey(relations[j])}); sort.Slice(contradictions,func(i,j int)bool{return relationKey(contradictions[i])<relationKey(contradictions[j])})
	return relations, contradictions
}

func normalizeBudget(b ContextBudget) ContextBudget {
	if b.MaxTokens <= 0 { b.MaxTokens = 4000 }
	if b.ReserveTokens < 0 { b.ReserveTokens = 0 }
	if b.ReserveTokens >= b.MaxTokens { b.ReserveTokens = 0 }
	b.UsedTokens = 0; b.RemainingTokens = b.MaxTokens - b.ReserveTokens
	return b
}

func defaultFidelity(plan QueryPlan) Fidelity {
	if plan.PreferredFidelity != "" { return plan.PreferredFidelity }
	if plan.NeedExact { return FidelityExact }
	if plan.NeedEvidence { return FidelityEvidence }
	return FidelitySummary
}

func selectLevel(n Node, wanted Fidelity, available int) (Fidelity, ContentLevel, int, bool) {
	if available <= 0 { return "", ContentLevel{}, 0, false }
	for _, fidelity := range fallbackOrder(wanted) {
		level, ok := n.Level(fidelity); if !ok { continue }
		cost := n.Estimate(fidelity); if cost <= 0 { cost = EstimateTokens(level.Text) }
		if cost <= available { return fidelity, level, cost, true }
	}
	return "", ContentLevel{}, 0, false
}

func fallbackOrder(wanted Fidelity) []Fidelity {
	switch wanted {
	case FidelityExact: return []Fidelity{FidelityExact,FidelityEvidence,FidelityDetail,FidelitySummary,FidelityAbstract,FidelityLabel}
	case FidelityEvidence: return []Fidelity{FidelityEvidence,FidelityDetail,FidelitySummary,FidelityAbstract,FidelityLabel}
	case FidelityDetail: return []Fidelity{FidelityDetail,FidelitySummary,FidelityAbstract,FidelityLabel}
	case FidelityAbstract: return []Fidelity{FidelityAbstract,FidelityLabel}
	case FidelityLabel: return []Fidelity{FidelityLabel}
	default: return []Fidelity{FidelitySummary,FidelityAbstract,FidelityLabel}
	}
}

func hasDeeper(n Node, current Fidelity) bool {
	order:=[]Fidelity{FidelityLabel,FidelityAbstract,FidelitySummary,FidelityDetail,FidelityEvidence,FidelityExact}; currentRank:=-1
	for i,f:=range order{if f==current{currentRank=i;break}}
	for i:=currentRank+1;i<len(order);i++{if _,ok:=n.Level(order[i]);ok{return true}}
	return false
}

func uniqueStrings(in []string) []string { seen:=map[string]struct{}{}; out:=make([]string,0,len(in)); for _,v:=range in{if v==""{continue};if _,ok:=seen[v];ok{continue};seen[v]=struct{}{};out=append(out,v)}; sort.Strings(out); return out }
func uniqueProvenance(in []Provenance) []Provenance { seen:=map[string]struct{}{}; out:=make([]Provenance,0,len(in)); for _,p:=range in{k:=p.SourceID+"\x00"+p.SourceRef+"\x00"+p.SourceHash;if _,ok:=seen[k];ok{continue};seen[k]=struct{}{};out=append(out,p)};return out }
func firstNonEmpty(values ...string) string { for _,v:=range values{if v!=""{return v}};return "" }
func relationKey(r ContextRelation) string { return fmt.Sprintf("%s\x00%s\x00%s",r.From,r.Relation,r.To) }
