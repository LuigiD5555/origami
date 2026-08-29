package hybridruntime

import (
	"fmt"
	"sort"

	"github.com/LuigiD5555/origami/internal/lab/hybridcarrier"
	"github.com/LuigiD5555/origami/internal/receiver"
)

type Transport struct {
	Columns int `json:"columns"`
	Cell    int `json:"cell_pixels"`
	Margin  int `json:"margin"`
}

type OpenMetrics struct {
	MemoryEntries int `json:"memory_entries"`
	IndexEntries  int `json:"index_entries"`
	BuildTouched  int `json:"build_touched"`
}

type ToolMetrics struct {
	Touched int `json:"touched"`
	Exposed int `json:"exposed"`
}

type Result struct {
	Operation string                      `json:"operation"`
	Entries   []hybridcarrier.MemoryEntry `json:"entries,omitempty"`
	Boot      *receiver.BootSpec          `json:"boot,omitempty"`
	Verified  *bool                       `json:"verified,omitempty"`
	Evidence  string                      `json:"evidence,omitempty"`
	Unknown   bool                        `json:"unknown"`
	Reason    string                      `json:"reason,omitempty"`
	Metrics   ToolMetrics                 `json:"metrics"`
}

type Runtime struct {
	envelope  hybridcarrier.Envelope
	byAddress map[string]hybridcarrier.MemoryEntry
	byKey     map[string]string
	metrics   OpenMetrics
}

// OpenPNG uses the carrier image itself as the source. It performs transport
// decoding and builds deterministic address maps once. BuildTouched is reported
// separately from per-query Touched so preprocessing cannot masquerade as free
// selective access.
func OpenPNG(pngBytes []byte, transport Transport) (*Runtime, error) {
	envelope, err := hybridcarrier.DecodePNG(pngBytes, transport.Columns, transport.Cell, transport.Margin)
	if err != nil { return nil, err }
	byAddress := make(map[string]hybridcarrier.MemoryEntry, len(envelope.Memory))
	for _, entry := range envelope.Memory {
		byAddress[entry.Address] = entry
	}
	byKey := make(map[string]string, len(envelope.Index))
	for _, entry := range envelope.Index {
		if _, exists := byKey[entry.Key]; exists {
			return nil, fmt.Errorf("duplicate index key %q", entry.Key)
		}
		byKey[entry.Key] = entry.Address
	}
	return &Runtime{
		envelope:  envelope,
		byAddress: byAddress,
		byKey:     byKey,
		metrics: OpenMetrics{
			MemoryEntries: len(envelope.Memory),
			IndexEntries:  len(envelope.Index),
			BuildTouched:  len(envelope.Memory) + len(envelope.Index),
		},
	}, nil
}

func (r *Runtime) OpenMetrics() OpenMetrics { return r.metrics }

func (r *Runtime) Boot() Result {
	boot := r.envelope.Receiver.Boot
	return Result{Operation: "BOOT", Boot: &boot, Evidence: "carrier:BOOT", Metrics: ToolMetrics{Touched: 1, Exposed: 1}}
}

func (r *Runtime) Lookup(query string) Result {
	address := query
	touched := 1
	if a, ok := r.byKey[query]; ok {
		address = a
		touched++
	}
	entry, ok := r.byAddress[address]
	if !ok {
		return Result{Operation: "LOOKUP", Unknown: true, Reason: "address/key not found", Metrics: ToolMetrics{Touched: touched, Exposed: 0}}
	}
	return Result{Operation: "LOOKUP", Entries: []hybridcarrier.MemoryEntry{entry}, Evidence: "carrier:" + entry.Address, Metrics: ToolMetrics{Touched: touched, Exposed: 1}}
}

// Follow traverses only declared links reachable from the requested start and
// relation. It never falls back to a scan when an address/link is missing.
func (r *Runtime) Follow(start, relation string, depth int) Result {
	if depth < 0 {
		return Result{Operation: "FOLLOW", Unknown: true, Reason: "depth must be non-negative"}
	}
	address := start
	touched := 1
	if a, ok := r.byKey[start]; ok {
		address = a
		touched++
	}
	current, ok := r.byAddress[address]
	if !ok {
		return Result{Operation: "FOLLOW", Unknown: true, Reason: "start address/key not found", Metrics: ToolMetrics{Touched: touched}}
	}
	entries := []hybridcarrier.MemoryEntry{current}
	frontier := []string{current.Address}
	visited := map[string]bool{current.Address: true}
	for level := 0; level < depth; level++ {
		var next []string
		for _, addr := range frontier {
			node, ok := r.byAddress[addr]
			touched++
			if !ok { continue }
			for _, target := range node.Links[relation] {
				if visited[target] { continue }
				visited[target] = true
				next = append(next, target)
			}
		}
		sort.Strings(next)
		if len(next) == 0 { break }
		for _, addr := range next {
			node, ok := r.byAddress[addr]
			touched++
			if !ok {
				return Result{Operation: "FOLLOW", Entries: entries, Unknown: true, Reason: "declared link target missing", Metrics: ToolMetrics{Touched: touched, Exposed: len(entries)}}
			}
			entries = append(entries, node)
		}
		frontier = next
	}
	return Result{Operation: "FOLLOW", Entries: entries, Evidence: "carrier:path:" + relation, Metrics: ToolMetrics{Touched: touched, Exposed: len(entries)}}
}

func (r *Runtime) Verify() Result {
	ok := r.envelope.Verification.MemorySHA256 != ""
	return Result{Operation: "VERIFY", Verified: &ok, Evidence: "sha256:" + r.envelope.Verification.MemorySHA256, Metrics: ToolMetrics{Touched: 1, Exposed: 1}}
}
