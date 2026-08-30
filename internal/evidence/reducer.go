package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/memory"
)

var hardInvariants = []string{
	"VERIFIED_EVIDENCE_ONLY",
	"CID_MISMATCH_REJECTED",
	"SOURCE_HASH_MISMATCH_REJECTED",
	"DUPLICATE_EVIDENCE_COUNTS_ONCE",
	"AGENT_CONFIDENCE_IS_NOT_AUTHORITY",
	"UNKNOWN_OVER_UNSUPPORTED_CLAIM",
	"EXACT_CLAIM_REQUIRES_BYTE_EQUAL_EXACT_EVIDENCE",
	"FALSE_EXACT=0",
}

type MemoryResolver struct { store *memory.Store }
func NewMemoryResolver(store *memory.Store) *MemoryResolver { return &MemoryResolver{store: store} }

func (r *MemoryResolver) Resolve(ref EvidenceRef) (ResolvedEvidence, string, bool) {
	if r == nil || r.store == nil { return ResolvedEvidence{}, "memory store unavailable", false }
	address := strings.TrimSpace(ref.Address)
	if address == "" { return ResolvedEvidence{}, "evidence address is empty", false }
	node, ok := r.store.Lookup(address)
	if !ok { return ResolvedEvidence{}, "evidence address not found", false }
	cid := node.DerivedCID()
	if ref.CID != "" && ref.CID != cid { return ResolvedEvidence{}, "evidence CID mismatch", false }
	if ref.SourceHash != "" && ref.SourceHash != node.Provenance.SourceHash { return ResolvedEvidence{}, "evidence source hash mismatch", false }
	fidelity := ref.Fidelity
	if fidelity == "" { fidelity = memory.FidelityEvidence }
	level, ok := node.Level(fidelity)
	if !ok { return ResolvedEvidence{}, fmt.Sprintf("requested fidelity %q unavailable", fidelity), false }
	if !level.Verified { return ResolvedEvidence{}, "evidence level is not verified", false }
	return ResolvedEvidence{
		Address: node.Address, Fidelity: fidelity, CID: cid,
		SourceID: node.Provenance.SourceID,
		SourceRef: firstNonEmpty(level.SourceRef, node.Provenance.SourceRef),
		SourceHash: node.Provenance.SourceHash,
		Verified: true, Exact: fidelity == memory.FidelityExact, Content: level.Text,
	}, "", true
}

type Reducer struct { resolver Resolver; policy Policy }
func NewReducer(resolver Resolver, policy Policy) *Reducer { return &Reducer{resolver: resolver, policy: normalizePolicy(policy)} }

func (r *Reducer) Reduce(proposals []Proposal) (Reduction, error) {
	if r == nil || r.resolver == nil { return Reduction{}, fmt.Errorf("evidence resolver is required") }
	policy := normalizePolicy(r.policy)
	prepared := append([]Proposal(nil), proposals...)
	for i := range prepared {
		if err := validateProposal(prepared[i]); err != nil { return Reduction{}, err }
		prepared[i].Proposition = strings.TrimSpace(prepared[i].Proposition)
		prepared[i].AgentID = strings.TrimSpace(prepared[i].AgentID)
		prepared[i].RunID = strings.TrimSpace(prepared[i].RunID)
	}
	sort.Slice(prepared, func(i, j int) bool {
		if prepared[i].ClaimKey == prepared[j].ClaimKey { return prepared[i].ID < prepared[j].ID }
		return prepared[i].ClaimKey < prepared[j].ClaimKey
	})
	result := Reduction{Schema: SchemaR0 + ".reduction", Policy: policy, Invariants: append([]string(nil), hardInvariants...)}
	result.Metrics.Proposals = len(prepared)
	groups := groupByClaim(prepared)
	keys := make([]string, 0, len(groups)); for key := range groups { keys = append(keys, key) }; sort.Strings(keys)
	for _, key := range keys {
		claim, metrics := r.reduceClaim(key, groups[key], policy)
		result.Claims = append(result.Claims, claim)
		result.Metrics.EvidenceRequested += metrics.EvidenceRequested
		result.Metrics.EvidenceAccepted += metrics.EvidenceAccepted
		result.Metrics.EvidenceRejected += metrics.EvidenceRejected
		result.Metrics.EvidenceDeduplicated += metrics.EvidenceDeduplicated
	}
	result.Metrics.Claims = len(result.Claims)
	result.Digest = reductionDigest(result)
	return result, nil
}

func (r *Reducer) reduceClaim(key string, proposals []Proposal, policy Policy) (ClaimResult, Metrics) {
	claim := ClaimResult{ClaimKey: key, Exactness: ExactnessNotClaimed}
	metrics := Metrics{}
	variants := map[string]struct{}{}
	for _, proposal := range proposals {
		claim.ProposalIDs = append(claim.ProposalIDs, proposal.ID)
		if proposal.Proposition != "" { variants[proposal.Proposition] = struct{}{} }
	}
	claim.PropositionVariants = sortedSet(variants)
	if len(claim.PropositionVariants) > 0 { claim.Proposition = claim.PropositionVariants[0] }

	supportAnchors := map[string]EvidenceAudit{}
	opposeAnchors := map[string]EvidenceAudit{}
	supportSources := map[string]struct{}{}
	opposeSources := map[string]struct{}{}
	exactMatchedPropositions := map[string]struct{}{}

	for _, proposal := range proposals {
		refs := append([]EvidenceRef(nil), proposal.Evidence...)
		sort.Slice(refs, func(i, j int) bool { return evidenceRefKey(refs[i]) < evidenceRefKey(refs[j]) })
		for _, ref := range refs {
			metrics.EvidenceRequested++
			resolved, reason, ok := r.resolver.Resolve(ref)
			exactMatch := ok && proposal.Position == PositionSupport && resolved.Exact && proposal.Proposition != "" && proposal.Proposition == resolved.Content
			audit := EvidenceAudit{ProposalID: proposal.ID, Position: proposal.Position, Requested: ref, Resolved: resolved, Accepted: ok, Reason: reason, ExactMatch: exactMatch}
			if !ok {
				metrics.EvidenceRejected++; claim.RejectedEvidence = append(claim.RejectedEvidence, audit); continue
			}
			audit.Anchor = evidenceAnchor(resolved); audit.SourceKey = evidenceSourceKey(resolved)
			switch proposal.Position {
			case PositionSupport:
				if _, exists := supportAnchors[audit.Anchor]; exists { metrics.EvidenceDeduplicated++; continue }
				supportAnchors[audit.Anchor] = audit; supportSources[audit.SourceKey] = struct{}{}; metrics.EvidenceAccepted++
				if exactMatch { exactMatchedPropositions[proposal.Proposition] = struct{}{} }
			case PositionOppose:
				if _, exists := opposeAnchors[audit.Anchor]; exists { metrics.EvidenceDeduplicated++; continue }
				opposeAnchors[audit.Anchor] = audit; opposeSources[audit.SourceKey] = struct{}{}; metrics.EvidenceAccepted++
			case PositionUnknown:
				audit.Accepted = false; audit.Reason = "UNKNOWN proposals cannot establish or refute a claim"
				metrics.EvidenceRejected++; claim.RejectedEvidence = append(claim.RejectedEvidence, audit)
			}
		}
	}

	claim.SupportEvidence = sortedAudits(supportAnchors); claim.OpposeEvidence = sortedAudits(opposeAnchors)
	sort.Slice(claim.RejectedEvidence, func(i, j int) bool { return auditKey(claim.RejectedEvidence[i]) < auditKey(claim.RejectedEvidence[j]) })
	claim.VerifiedSupportAnchors = len(claim.SupportEvidence); claim.VerifiedOpposeAnchors = len(claim.OpposeEvidence)
	claim.IndependentSupport = len(nonEmptySet(supportSources)); claim.IndependentOppose = len(nonEmptySet(opposeSources))
	supportOK := claim.VerifiedSupportAnchors >= policy.MinVerifiedEvidence && claim.IndependentSupport >= policy.MinIndependentSources
	opposeOK := claim.VerifiedOpposeAnchors >= policy.MinVerifiedEvidence && claim.IndependentOppose >= policy.MinIndependentSources
	switch {
	case supportOK && opposeOK:
		claim.Status = StatusConflict; claim.Exactness = ExactnessNotClaimed; claim.Reason = "verified evidence exists on both sides; deterministic reduction preserves the conflict"
	case supportOK:
		claim.Status = StatusVerified; claim.Reason = "verified support evidence satisfies the reduction policy"
		if _, ok := exactMatchedPropositions[claim.Proposition]; ok { claim.Exactness = ExactnessVerified }
	case opposeOK:
		claim.Status = StatusRejected; claim.Exactness = ExactnessNotClaimed; claim.Reason = "verified opposing evidence satisfies the reduction policy"
	default:
		claim.Status = StatusUnknown; claim.Exactness = ExactnessNotClaimed; claim.Reason = "insufficient verified evidence; UNKNOWN is preferred to invention"
	}
	claim.CanonicalID = claimDigest(claim)
	return claim, metrics
}

func validateProposal(p Proposal) error {
	if strings.TrimSpace(p.ID) == "" { return fmt.Errorf("proposal id is required") }
	if strings.TrimSpace(p.ClaimKey) == "" { return fmt.Errorf("proposal %q claim_key is required", p.ID) }
	switch p.Position { case PositionSupport, PositionOppose, PositionUnknown: default: return fmt.Errorf("proposal %q has invalid position %q", p.ID, p.Position) }
	if p.Confidence < 0 || p.Confidence > 1 { return fmt.Errorf("proposal %q confidence must be between 0 and 1", p.ID) }
	return nil
}
func normalizePolicy(policy Policy) Policy { if policy.MinVerifiedEvidence <= 0 { policy.MinVerifiedEvidence = 1 }; if policy.MinIndependentSources <= 0 { policy.MinIndependentSources = 1 }; return policy }
func groupByClaim(proposals []Proposal) map[string][]Proposal { groups:=map[string][]Proposal{}; for _,proposal:=range proposals{key:=strings.TrimSpace(proposal.ClaimKey);proposal.ClaimKey=key;groups[key]=append(groups[key],proposal)};return groups }
func evidenceRefKey(ref EvidenceRef) string { return strings.Join([]string{strings.TrimSpace(ref.Address),string(ref.Fidelity),ref.CID,ref.SourceHash},"|") }
func evidenceAnchor(resolved ResolvedEvidence) string { return strings.Join([]string{resolved.CID,string(resolved.Fidelity),resolved.SourceHash},"|") }
func evidenceSourceKey(resolved ResolvedEvidence) string { if resolved.SourceHash!=""{return "hash:"+resolved.SourceHash};if resolved.SourceID!=""{return "source:"+resolved.SourceID};return "cid:"+resolved.CID }
func sortedAudits(in map[string]EvidenceAudit) []EvidenceAudit { keys:=make([]string,0,len(in));for key:=range in{keys=append(keys,key)};sort.Strings(keys);out:=make([]EvidenceAudit,0,len(keys));for _,key:=range keys{out=append(out,in[key])};return out }
func sortedSet(in map[string]struct{}) []string { out:=make([]string,0,len(in));for value:=range in{out=append(out,value)};sort.Strings(out);return out }
func nonEmptySet(in map[string]struct{}) map[string]struct{} { out:=map[string]struct{}{};for key:=range in{if key!=""{out[key]=struct{}{}}};return out }
func auditKey(a EvidenceAudit) string { return strings.Join([]string{a.ProposalID,string(a.Position),evidenceRefKey(a.Requested),a.Reason},"|") }
func claimDigest(claim ClaimResult) string { var b strings.Builder;fmt.Fprintf(&b,"%s\x00%s\x00%s\x00",claim.ClaimKey,claim.Status,claim.Exactness);for _,id:=range claim.ProposalIDs{fmt.Fprintf(&b,"p:%s\x00",id)};for _,item:=range claim.SupportEvidence{fmt.Fprintf(&b,"s:%s:%t\x00",item.Anchor,item.ExactMatch)};for _,item:=range claim.OpposeEvidence{fmt.Fprintf(&b,"o:%s\x00",item.Anchor)};sum:=sha256.Sum256([]byte(b.String()));return "sha256:"+hex.EncodeToString(sum[:]) }
func reductionDigest(reduction Reduction) string { var b strings.Builder;fmt.Fprintf(&b,"%s\x00%d\x00%d\x00",reduction.Schema,reduction.Policy.MinVerifiedEvidence,reduction.Policy.MinIndependentSources);for _,claim:=range reduction.Claims{fmt.Fprintf(&b,"%s\x00%s\x00%s\x00",claim.ClaimKey,claim.Status,claim.CanonicalID)};sum:=sha256.Sum256([]byte(b.String()));return "sha256:"+hex.EncodeToString(sum[:]) }
func firstNonEmpty(values ...string) string { for _,value:=range values{if value!=""{return value}};return "" }
