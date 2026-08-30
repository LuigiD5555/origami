package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/memory"
)

type CarrierCounts struct {
	Pages     uint32
	Blocks    uint32
	Documents uint32
	Objects   uint32
}

// FixedCarrierMetadata binds the frozen visual control plane to the current
// canonical memory state. The pixels remain a boot/navigation/control plane;
// StoreRoot commits to the external semantic/exact memory universe.
func FixedCarrierMetadata(system *System, sourceHash string, counts CarrierCounts) (fixedcarrier.Metadata, error) {
	nodes := system.Nodes()
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Address < nodes[j].Address })
	body, err := json.Marshal(nodes)
	if err != nil {
		return fixedcarrier.Metadata{}, err
	}
	root := sha256.Sum256(body)
	signature := system.Signature()
	sigBody, err := compactSignature(signature, nodes)
	if err != nil {
		return fixedcarrier.Metadata{}, err
	}
	objects := counts.Objects
	if objects == 0 {
		objects = uint32(len(nodes))
	}
	return fixedcarrier.Metadata{
		CarrierID:      system.CarrierID,
		StoreRoot:      hex.EncodeToString(root[:]),
		SourceSHA256:   sourceHash,
		PageCount:      counts.Pages,
		BlockCount:     counts.Blocks,
		DocumentCount:  counts.Documents,
		ObjectCount:    objects,
		GraphSignature: sigBody,
		Flags:          0x0003, // external memory binding + semantic T2 hint present
	}, nil
}

func CanonicalStoreRoot(nodes []memory.Node) string {
	copyNodes := append([]memory.Node(nil), nodes...)
	sort.Slice(copyNodes, func(i, j int) bool { return copyNodes[i].Address < copyNodes[j].Address })
	body, _ := json.Marshal(copyNodes)
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// compactSignature is deliberately a valid bounded JSON object. Earlier code
// truncated JSON at 256 bytes, which could leave an unreadable routing hint.
// T2 now gives semantic navigation priority to a tiny set of top-level labels.
func compactSignature(signature memory.GraphSignature, nodes []memory.Node) ([]byte, error) {
	type payload struct {
		Schema   string   `json:"s"`
		Carrier string   `json:"c,omitempty"`
		Index   []string `json:"i,omitempty"`
		Roots   []string `json:"r,omitempty"`
		Classes []string `json:"k,omitempty"`
		Groups  []string `json:"g,omitempty"`
		Nodes   int      `json:"n"`
	}
	p := payload{
		Schema:   signature.Schema,
		Carrier: trimASCII(signature.CarrierID, 20),
		Index:   semanticIndexLabels(nodes, 4),
		Roots:   trimStrings(signature.RootTerms, 3, 18),
		Classes: trimStrings(signature.SemanticClasses, 2, 16),
		Groups:  trimStrings(signature.Clusters, 2, 16),
		Nodes:   signature.NodeCount,
	}
	marshal := func() ([]byte, error) { return json.Marshal(p) }
	body, err := marshal()
	if err != nil {
		return nil, err
	}
	if len(body) <= 256 {
		return body, nil
	}
	// Remove lower-value routing hints before touching the semantic index.
	p.Groups = nil
	body, err = marshal()
	if err != nil { return nil, err }
	if len(body) <= 256 { return body, nil }
	p.Classes = nil
	body, err = marshal()
	if err != nil { return nil, err }
	if len(body) <= 256 { return body, nil }
	p.Roots = nil
	body, err = marshal()
	if err != nil { return nil, err }
	if len(body) <= 256 { return body, nil }
	p.Carrier = ""
	p.Index = trimStrings(p.Index, 4, 16)
	body, err = marshal()
	if err != nil { return nil, err }
	if len(body) <= 256 { return body, nil }
	p.Index = trimStrings(p.Index, 3, 12)
	return json.Marshal(p)
}

func semanticIndexLabels(nodes []memory.Node, limit int) []string {
	type ranked struct {
		label string
		addr  string
		score float64
	}
	seen := map[string]bool{}
	var candidates []ranked
	for _, n := range nodes {
		label := strings.TrimSpace(n.Label)
		if label == "" || n.Parent != "" {
			continue
		}
		key := strings.ToUpper(label)
		if seen[key] { continue }
		seen[key] = true
		candidates = append(candidates, ranked{label: label, addr: n.Address, score: n.Importance.Global})
	}
	if len(candidates) < limit {
		for _, n := range nodes {
			label := strings.TrimSpace(n.Label)
			if label == "" { continue }
			key := strings.ToUpper(label)
			if seen[key] { continue }
			seen[key] = true
			candidates = append(candidates, ranked{label: label, addr: n.Address, score: n.Importance.Global - 0.01})
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score { return candidates[i].score > candidates[j].score }
		return candidates[i].addr < candidates[j].addr
	})
	out := make([]string, 0, limit)
	for _, c := range candidates {
		if len(out) >= limit { break }
		out = append(out, trimASCII(c.label, 24))
	}
	return out
}

func trimStrings(values []string, n, width int) []string {
	if n <= 0 || len(values) == 0 { return nil }
	if len(values) > n { values = values[:n] }
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = trimASCII(value, width)
		if value != "" { out = append(out, value) }
	}
	return out
}

func trimASCII(value string, width int) string {
	value = strings.Join(strings.Fields(value), " ")
	if width > 0 && len(value) > width { value = value[:width] }
	return value
}
