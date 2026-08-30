package core

import (
	"fmt"
	"sort"
)

type SemanticNode struct {
	Address      string              `json:"address"`
	Kind         string              `json:"kind,omitempty"`
	Label        string              `json:"label,omitempty"`
	Value        Value               `json:"value"`
	Dependencies []string            `json:"dependencies,omitempty"`
	Alternatives []string            `json:"alternatives,omitempty"`
	Relations    map[string][]string `json:"relations,omitempty"`
	EvidenceRefs []string            `json:"evidence_refs,omitempty"`
	Verified     bool                `json:"verified,omitempty"`
}

type SemanticGraph struct {
	Schema string                  `json:"schema"`
	Nodes  map[string]SemanticNode `json:"nodes"`
}

type FoldIndexEntry struct {
	Address      string   `json:"address"`
	CID          string   `json:"cid"`
	Dependencies []string `json:"dependencies,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
}

type FoldedGraph struct {
	Schema     string           `json:"schema"`
	Commitment string           `json:"commitment"`
	Index      []FoldIndexEntry `json:"index"`
	NodeCount  int              `json:"node_count"`
}

type UnfoldResult struct {
	Schema       string         `json:"schema"`
	Commitment   string         `json:"commitment"`
	Requested    []string       `json:"requested"`
	Closure      []SemanticNode `json:"closure"`
	Addresses    []string       `json:"addresses"`
	Touched      []string       `json:"touched"`
	FullUnfold   bool           `json:"full_unfold"`
	VerifiedFold bool           `json:"verified_fold"`
}

type GraphStore struct {
	nodes   map[string]SemanticNode
	touched map[string]int
}

func FoldGraph(graph SemanticGraph) (FoldedGraph, *GraphStore, error) {
	if graph.Schema != "" && graph.Schema != SchemaR1+".semantic-graph" {
		return FoldedGraph{}, nil, fmt.Errorf("unexpected semantic graph schema %q", graph.Schema)
	}
	if len(graph.Nodes) == 0 {
		return FoldedGraph{}, nil, fmt.Errorf("semantic graph cannot be empty")
	}
	addresses := make([]string, 0, len(graph.Nodes))
	for address, node := range graph.Nodes {
		if address == "" || node.Address == "" || address != node.Address {
			return FoldedGraph{}, nil, fmt.Errorf("semantic node map key/address mismatch for %q", address)
		}
		if !ValidateStatus(node.Value.Status) {
			return FoldedGraph{}, nil, fmt.Errorf("semantic node %q has invalid status %q", address, node.Value.Status)
		}
		addresses = append(addresses, address)
	}
	sort.Strings(addresses)
	for _, address := range addresses {
		node := graph.Nodes[address]
		for _, dependency := range append(append([]string(nil), node.Dependencies...), node.Alternatives...) {
			if _, ok := graph.Nodes[dependency]; !ok {
				return FoldedGraph{}, nil, fmt.Errorf("semantic node %q references missing node %q", address, dependency)
			}
		}
		for relation, targets := range node.Relations {
			if relation == "" {
				return FoldedGraph{}, nil, fmt.Errorf("semantic node %q has empty relation", address)
			}
			for _, target := range targets {
				if _, ok := graph.Nodes[target]; !ok {
					return FoldedGraph{}, nil, fmt.Errorf("semantic node %q relation %q references missing node %q", address, relation, target)
				}
			}
		}
	}

	canonical := canonicalGraph(graph)
	commitment := hashJSON(canonical)
	index := make([]FoldIndexEntry, 0, len(addresses))
	storeNodes := make(map[string]SemanticNode, len(addresses))
	for _, address := range addresses {
		node := canonical.Nodes[address]
		storeNodes[address] = cloneSemanticNode(node)
		index = append(index, FoldIndexEntry{
			Address:      address,
			CID:          SemanticNodeCID(node),
			Dependencies: append([]string(nil), node.Dependencies...),
			Alternatives: append([]string(nil), node.Alternatives...),
		})
	}
	folded := FoldedGraph{Schema: SchemaR1 + ".folded-graph", Commitment: commitment, Index: index, NodeCount: len(index)}
	store := &GraphStore{nodes: storeNodes, touched: map[string]int{}}
	return folded, store, nil
}

func SemanticNodeCID(node SemanticNode) string {
	return hashJSON(cloneSemanticNode(node))
}

func (s *GraphStore) Lookup(address string) (SemanticNode, bool) {
	if s == nil {
		return SemanticNode{}, false
	}
	node, ok := s.nodes[address]
	if !ok {
		return SemanticNode{}, false
	}
	if s.touched == nil {
		s.touched = map[string]int{}
	}
	s.touched[address]++
	return cloneSemanticNode(node), true
}

func (s *GraphStore) ResetTouches() {
	if s != nil {
		s.touched = map[string]int{}
	}
}

func (s *GraphStore) TouchedAddresses() []string {
	if s == nil {
		return nil
	}
	out := make([]string, 0, len(s.touched))
	for address := range s.touched {
		out = append(out, address)
	}
	sort.Strings(out)
	return out
}

func (s *GraphStore) TouchCount() int {
	if s == nil {
		return 0
	}
	count := 0
	for _, n := range s.touched {
		count += n
	}
	return count
}

func VerifyFold(folded FoldedGraph, store *GraphStore) bool {
	if store == nil || folded.NodeCount != len(store.nodes) || folded.NodeCount != len(folded.Index) {
		return false
	}
	graph := SemanticGraph{Schema: SchemaR1 + ".semantic-graph", Nodes: map[string]SemanticNode{}}
	for _, entry := range folded.Index {
		node, ok := store.nodes[entry.Address]
		if !ok || SemanticNodeCID(node) != entry.CID {
			return false
		}
		graph.Nodes[entry.Address] = cloneSemanticNode(node)
	}
	return hashJSON(canonicalGraph(graph)) == folded.Commitment
}

func SelectiveUnfold(folded FoldedGraph, store *GraphStore, requested []string) (UnfoldResult, error) {
	if store == nil {
		return UnfoldResult{}, fmt.Errorf("graph store unavailable")
	}
	if len(requested) == 0 {
		return UnfoldResult{}, fmt.Errorf("selective unfold requires at least one address")
	}
	allowed := map[string]FoldIndexEntry{}
	for _, entry := range folded.Index {
		allowed[entry.Address] = entry
	}
	queue := append([]string(nil), requested...)
	sort.Strings(queue)
	seen := map[string]struct{}{}
	var nodes []SemanticNode
	for len(queue) > 0 {
		address := queue[0]
		queue = queue[1:]
		if _, ok := seen[address]; ok {
			continue
		}
		entry, ok := allowed[address]
		if !ok {
			return UnfoldResult{}, fmt.Errorf("address %q is not part of folded graph", address)
		}
		node, ok := store.Lookup(address)
		if !ok {
			return UnfoldResult{}, fmt.Errorf("address %q missing from canonical store", address)
		}
		if SemanticNodeCID(node) != entry.CID {
			return UnfoldResult{}, fmt.Errorf("CID mismatch for %q", address)
		}
		seen[address] = struct{}{}
		nodes = append(nodes, node)
		next := append([]string(nil), entry.Dependencies...)
		next = append(next, entry.Alternatives...)
		sort.Strings(next)
		for _, dependency := range next {
			if _, ok := seen[dependency]; !ok {
				queue = append(queue, dependency)
			}
		}
		sort.Strings(queue)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Address < nodes[j].Address })
	addresses := make([]string, 0, len(nodes))
	for _, node := range nodes {
		addresses = append(addresses, node.Address)
	}
	requestedCopy := append([]string(nil), requested...)
	sort.Strings(requestedCopy)
	return UnfoldResult{
		Schema:       SchemaR1 + ".unfold-result",
		Commitment:   folded.Commitment,
		Requested:    requestedCopy,
		Closure:      nodes,
		Addresses:    addresses,
		Touched:      store.TouchedAddresses(),
		FullUnfold:   len(addresses) == folded.NodeCount,
		VerifiedFold: VerifyFoldWithoutTouch(folded, store),
	}, nil
}

func FullUnfold(folded FoldedGraph, store *GraphStore) (UnfoldResult, error) {
	addresses := make([]string, 0, len(folded.Index))
	for _, entry := range folded.Index {
		addresses = append(addresses, entry.Address)
	}
	result, err := SelectiveUnfold(folded, store, addresses)
	if err != nil {
		return UnfoldResult{}, err
	}
	result.FullUnfold = true
	return result, nil
}

func VerifyFoldWithoutTouch(folded FoldedGraph, store *GraphStore) bool {
	if store == nil || folded.NodeCount != len(store.nodes) || folded.NodeCount != len(folded.Index) {
		return false
	}
	graph := SemanticGraph{Schema: SchemaR1 + ".semantic-graph", Nodes: map[string]SemanticNode{}}
	for _, entry := range folded.Index {
		node, ok := store.nodes[entry.Address]
		if !ok || SemanticNodeCID(node) != entry.CID {
			return false
		}
		graph.Nodes[entry.Address] = cloneSemanticNode(node)
	}
	return hashJSON(canonicalGraph(graph)) == folded.Commitment
}

func canonicalGraph(graph SemanticGraph) SemanticGraph {
	out := SemanticGraph{Schema: SchemaR1 + ".semantic-graph", Nodes: map[string]SemanticNode{}}
	for address, raw := range graph.Nodes {
		node := cloneSemanticNode(raw)
		sort.Strings(node.Dependencies)
		sort.Strings(node.Alternatives)
		sort.Strings(node.EvidenceRefs)
		if node.Relations != nil {
			for relation, targets := range node.Relations {
				copyTargets := append([]string(nil), targets...)
				sort.Strings(copyTargets)
				node.Relations[relation] = copyTargets
			}
		}
		out.Nodes[address] = node
	}
	return out
}

func cloneSemanticNode(node SemanticNode) SemanticNode {
	out := node
	out.Dependencies = append([]string(nil), node.Dependencies...)
	out.Alternatives = append([]string(nil), node.Alternatives...)
	out.EvidenceRefs = append([]string(nil), node.EvidenceRefs...)
	if node.Relations != nil {
		out.Relations = map[string][]string{}
		for relation, targets := range node.Relations {
			out.Relations[relation] = append([]string(nil), targets...)
		}
	}
	return out
}
