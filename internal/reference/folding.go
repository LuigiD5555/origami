package reference

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// FoldGuard is a semantic predicate that must remain visible on a compiled edge.
type FoldGuard struct {
	Kind  string `json:"kind"` // precondition, condition, exception
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FoldStep is one canonical dependency in a derivation path.
type FoldStep struct {
	From   string      `json:"from"`
	To     string      `json:"to"`
	Rule   string      `json:"rule,omitempty"`
	Guards []FoldGuard `json:"guards,omitempty"`
}

// CompiledEdge is a derived view. It never replaces the canonical steps.
// ProofRef commits to the exact ordered derivation so VERIFY can reopen it.
type CompiledEdge struct {
	From       string      `json:"from"`
	To         string      `json:"to"`
	Depth      int         `json:"depth"`
	Guards     []FoldGuard `json:"guards,omitempty"`
	ProofRef   string      `json:"proof_ref"`
	Canonical  []FoldStep  `json:"-"`
}

func canonicalCommitment(steps []FoldStep) string {
	h := sha256.New()
	for _, s := range steps {
		fmt.Fprintf(h, "%s\x00%s\x00%s\x00", s.From, s.To, s.Rule)
		guards := append([]FoldGuard(nil), s.Guards...)
		sort.Slice(guards, func(i, j int) bool {
			if guards[i].Kind != guards[j].Kind { return guards[i].Kind < guards[j].Kind }
			if guards[i].Key != guards[j].Key { return guards[i].Key < guards[j].Key }
			return guards[i].Value < guards[j].Value
		})
		for _, g := range guards { fmt.Fprintf(h, "%s=%s=%s\x00", g.Kind, g.Key, g.Value) }
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

// FoldDependency compiles a contiguous canonical path into one derived edge.
// Guards are preserved instead of being summarized away.
func FoldDependency(steps []FoldStep) (CompiledEdge, error) {
	if len(steps) == 0 { return CompiledEdge{}, fmt.Errorf("cannot fold empty dependency") }
	for i := 1; i < len(steps); i++ {
		if steps[i-1].To != steps[i].From {
			return CompiledEdge{}, fmt.Errorf("non-contiguous dependency at step %d", i)
		}
	}
	var guards []FoldGuard
	for _, s := range steps { guards = append(guards, s.Guards...) }
	return CompiledEdge{
		From: steps[0].From, To: steps[len(steps)-1].To, Depth: len(steps),
		Guards: guards, ProofRef: canonicalCommitment(steps), Canonical: append([]FoldStep(nil), steps...),
	}, nil
}

// VerifyFold proves that a compiled edge still commits to its canonical path.
func VerifyFold(edge CompiledEdge) bool {
	if len(edge.Canonical) != edge.Depth || edge.Depth == 0 { return false }
	if edge.Canonical[0].From != edge.From || edge.Canonical[len(edge.Canonical)-1].To != edge.To { return false }
	return strings.EqualFold(edge.ProofRef, canonicalCommitment(edge.Canonical))
}

// FoldLevel groups a long path into bounded blocks. Repeating this operation
// over the returned summaries yields logarithmic cognitive depth while the
// canonical derivation remains independently addressable.
func FoldLevel(steps []FoldStep, blockSize int) ([]CompiledEdge, error) {
	if blockSize < 2 { return nil, fmt.Errorf("block size must be >= 2") }
	out := make([]CompiledEdge, 0, (len(steps)+blockSize-1)/blockSize)
	for start := 0; start < len(steps); start += blockSize {
		end := start + blockSize
		if end > len(steps) { end = len(steps) }
		edge, err := FoldDependency(steps[start:end]); if err != nil { return nil, err }
		out = append(out, edge)
	}
	return out, nil
}
