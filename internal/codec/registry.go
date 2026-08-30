package codec

import (
	"fmt"
	"sort"
	"strings"
)

type Family string

const (
	SemanticDecoder Family = "SEMANTIC_DECODER"
	SemanticEncoder Family = "SEMANTIC_ENCODER"
	ExactDecoder    Family = "EXACT_DECODER"
	ExactEncoder    Family = "EXACT_ENCODER"
)

type Entry struct {
	ID                   string   `json:"id"`
	Family               Family   `json:"family"`
	Operation            string   `json:"operation"`
	Pair                 string   `json:"pair,omitempty"`
	Exact                bool     `json:"exact"`
	RequiredCapabilities []string `json:"required_capabilities,omitempty"`
	Fallbacks            []string `json:"fallbacks,omitempty"`
	FailureState         string   `json:"failure_state"`
	Steps                []string `json:"steps,omitempty"`
}

type Registry struct {
	Schema  string  `json:"schema"`
	Entries []Entry `json:"entries"`
}

func ReferenceRegistry() Registry {
	return Registry{
		Schema: "origami.codec-registry.r0",
		Entries: []Entry{
			semantic("S0", SemanticDecoder, "READ_IDENTITY", "E0", "UNKNOWN", "LOCATE_T0", "READ_DECLARED_IDENTITY"),
			semantic("E0", SemanticEncoder, "ENCODE_IDENTITY", "S0", "CONSTRUCTION_SPEC_ONLY", "DERIVE_IDENTITY", "BIND_IDENTITY_TO_PROFILE"),
			semantic("S1", SemanticDecoder, "READ_HIERARCHY", "E1", "UNKNOWN", "LOCATE_DECLARED_HIERARCHY", "READ_PARENT_CHILD_STRUCTURE"),
			semantic("E1", SemanticEncoder, "ENCODE_HIERARCHY", "S1", "CONSTRUCTION_SPEC_ONLY", "DERIVE_HIERARCHY", "BIND_PARENT_CHILD_STRUCTURE"),
			semantic("S2", SemanticDecoder, "READ_SUPERINDEX", "E2", "UNKNOWN", "LOCATE_T2", "READ_BOUNDED_ACTUAL_ENTRIES", "STOP_WITHOUT_EXACT_ESCALATION"),
			semantic("E2", SemanticEncoder, "ENCODE_SUPERINDEX", "S2", "CONSTRUCTION_SPEC_ONLY", "CANONICALIZE_INDEX", "EMIT_T2_CONSTRUCTION_IR"),
			semantic("S3", SemanticDecoder, "LOCATE_TOPIC", "E3", "UNKNOWN", "READ_T2", "MATCH_DECLARED_TOPIC_ENTRY", "RETURN_ADDRESS_OR_REGION"),
			semantic("E3", SemanticEncoder, "ENCODE_CONCEPT", "S3", "CONSTRUCTION_SPEC_ONLY", "DERIVE_CONCEPT", "ASSIGN_DECLARED_SEMANTIC_REGION"),
			semantic("S4", SemanticDecoder, "READ_RELATION", "E4", "UNKNOWN", "READ_DECLARED_RELATION_PRIMITIVE", "RESOLVE_ENDPOINTS"),
			semantic("E4", SemanticEncoder, "ENCODE_RELATION", "S4", "CONSTRUCTION_SPEC_ONLY", "DERIVE_RELATION", "BIND_RELATION_PRIMITIVE"),
			semantic("S5", SemanticDecoder, "EXPAND_CLUSTER", "E5", "UNKNOWN", "LOCATE_CLUSTER", "OPEN_ONE_DECLARED_LEVEL"),
			semantic("E5", SemanticEncoder, "ENCODE_CLUSTER", "S5", "CONSTRUCTION_SPEC_ONLY", "GROUP_SEMANTIC_MEMBERS", "EMIT_CLUSTER_SCOPE"),
			semantic("S6", SemanticDecoder, "SEMANTIC_UNFOLD", "E6", "UNKNOWN", "SELECT_SMALLEST_RELEVANT_REGION", "UNFOLD_WITHIN_CONTEXT_BUDGET"),
			semantic("E6", SemanticEncoder, "SEMANTIC_FOLD", "S6", "CONSTRUCTION_SPEC_ONLY", "FACTOR_SEMANTIC_STRUCTURE", "EMIT_ADDRESSABLE_FOLD"),
			exact("X0", ExactDecoder, "RESOLVE_ADDRESS", "NOT_VERIFIED", "EXACT_READ"),
			exact("X1", ExactDecoder, "READ_CID", "NOT_VERIFIED", "EXACT_READ"),
			exact("X2", ExactDecoder, "VERIFY_HASH", "NOT_VERIFIED", "EXACT_READ", "VERIFIER"),
			exact("X3", ExactDecoder, "VERIFY_MERKLE", "NOT_VERIFIED", "EXACT_READ", "VERIFIER"),
			exact("X4", ExactDecoder, "DECODE_RESIDUAL", "NOT_VERIFIED", "EXACT_READ"),
			exact("X5", ExactDecoder, "EXACT_RECOVERY", "NOT_VERIFIED", "EXACT_READ"),
			exact("Q0", ExactEncoder, "ASSIGN_ADDRESS", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE"),
			exact("Q1", ExactEncoder, "CREATE_CID", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE"),
			exact("Q2", ExactEncoder, "HASH_CONTENT", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE"),
			exact("Q3", ExactEncoder, "BUILD_MERKLE", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE"),
			exact("Q4", ExactEncoder, "ENCODE_RESIDUAL", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE"),
			exact("Q5", ExactEncoder, "BUILD_VERIFICATION", "CONSTRUCTION_SPEC_ONLY", "EXACT_WRITE", "VERIFIER"),
		},
	}
}

func semantic(id string, family Family, operation, pair, failure string, steps ...string) Entry {
	capability := "SEMANTIC_READ"
	if family == SemanticEncoder {
		capability = "SEMANTIC_WRITE"
	}
	return Entry{ID: id, Family: family, Operation: operation, Pair: pair, Exact: false, RequiredCapabilities: []string{capability}, FailureState: failure, Steps: steps}
}

func exact(id string, family Family, operation, failure string, caps ...string) Entry {
	return Entry{ID: id, Family: family, Operation: operation, Exact: true, RequiredCapabilities: append([]string(nil), caps...), FailureState: failure}
}

func (r Registry) Lookup(id string) (Entry, bool) {
	id = strings.ToUpper(strings.TrimSpace(id))
	for _, e := range r.Entries {
		if e.ID == id {
			return e, true
		}
	}
	return Entry{}, false
}

func (r Registry) Validate() error {
	if r.Schema != "origami.codec-registry.r0" {
		return fmt.Errorf("unexpected registry schema %q", r.Schema)
	}
	seen := map[string]bool{}
	for _, e := range r.Entries {
		if e.ID == "" || e.Operation == "" || e.FailureState == "" {
			return fmt.Errorf("incomplete codec entry: %+v", e)
		}
		if seen[e.ID] {
			return fmt.Errorf("duplicate codec id %s", e.ID)
		}
		seen[e.ID] = true
		if (e.Family == SemanticDecoder || e.Family == SemanticEncoder) && e.Exact {
			return fmt.Errorf("semantic codec %s cannot be exact", e.ID)
		}
		if e.Family == SemanticDecoder || e.Family == SemanticEncoder {
			for _, cap := range e.RequiredCapabilities {
				if strings.HasPrefix(cap, "EXACT_") || cap == "VERIFIER" {
					return fmt.Errorf("semantic codec %s hides exact dependency %s", e.ID, cap)
				}
			}
			if e.Pair == "" {
				return fmt.Errorf("semantic codec %s missing pair", e.ID)
			}
		}
	}
	for _, e := range r.Entries {
		if e.Family != SemanticDecoder && e.Family != SemanticEncoder {
			continue
		}
		pair, ok := r.Lookup(e.Pair)
		if !ok {
			return fmt.Errorf("codec %s pair %s missing", e.ID, e.Pair)
		}
		if pair.Pair != e.ID {
			return fmt.Errorf("codec pair not symmetric: %s -> %s -> %s", e.ID, pair.ID, pair.Pair)
		}
		if e.Family == pair.Family {
			return fmt.Errorf("codec %s pair has same family", e.ID)
		}
	}
	return nil
}

func (r Registry) Select(operation string, family Family, capabilities map[string]bool) (Entry, error) {
	operation = strings.ToUpper(strings.TrimSpace(operation))
	candidates := make([]Entry, 0)
	for _, e := range r.Entries {
		if e.Family == family && e.Operation == operation {
			candidates = append(candidates, e)
		}
	}
	if len(candidates) == 0 {
		return Entry{}, fmt.Errorf("codec not declared for %s/%s", family, operation)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })
	for _, e := range candidates {
		if supports(capabilities, e.RequiredCapabilities) {
			return e, nil
		}
	}
	return Entry{}, fmt.Errorf("declared codec unavailable for %s/%s", family, operation)
}

func supports(capabilities map[string]bool, required []string) bool {
	for _, cap := range required {
		if !capabilities[cap] {
			return false
		}
	}
	return true
}

type IndexEntry struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Address string `json:"address,omitempty"`
}

type SemanticIndex struct {
	Entries []IndexEntry `json:"entries"`
}

type ConstructionIR struct {
	Schema    string       `json:"schema"`
	ProfileID string       `json:"profile_id"`
	CodecID   string       `json:"codec_id"`
	Entries   []IndexEntry `json:"entries"`
}

func EncodeSuperIndex(index SemanticIndex, profileID string) (ConstructionIR, error) {
	if len(index.Entries) == 0 {
		return ConstructionIR{}, fmt.Errorf("semantic index is empty")
	}
	entries := make([]IndexEntry, 0, len(index.Entries))
	seen := map[string]bool{}
	for _, in := range index.Entries {
		e := IndexEntry{ID: strings.TrimSpace(in.ID), Label: strings.TrimSpace(in.Label), Address: strings.TrimSpace(in.Address)}
		if e.ID == "" || e.Label == "" {
			return ConstructionIR{}, fmt.Errorf("index entry requires id and label")
		}
		if seen[e.ID] {
			return ConstructionIR{}, fmt.Errorf("duplicate index entry id %s", e.ID)
		}
		seen[e.ID] = true
		entries = append(entries, e)
	}
	return ConstructionIR{Schema: "origami.t2-construction-ir.r0", ProfileID: strings.TrimSpace(profileID), CodecID: "E2", Entries: entries}, nil
}

func DecodeSuperIndex(ir ConstructionIR) (SemanticIndex, error) {
	if ir.Schema != "origami.t2-construction-ir.r0" || ir.CodecID != "E2" {
		return SemanticIndex{}, fmt.Errorf("unsupported T2 construction IR")
	}
	if len(ir.Entries) == 0 {
		return SemanticIndex{}, fmt.Errorf("T2 construction IR has no entries")
	}
	out := make([]IndexEntry, len(ir.Entries))
	copy(out, ir.Entries)
	return SemanticIndex{Entries: out}, nil
}

func EqualSemanticIndex(a, b SemanticIndex) bool {
	if len(a.Entries) != len(b.Entries) {
		return false
	}
	for i := range a.Entries {
		if a.Entries[i] != b.Entries[i] {
			return false
		}
	}
	return true
}
