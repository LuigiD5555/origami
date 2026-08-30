package writer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/core"
)

var canonicalRosetta = []RosettaEntry{
	{VisualToken: "RINGED_SQUARE", SemanticRole: "ROOT_OR_VERIFY_ANCHOR", Dimension: "D01_GEOMETRY"},
	{VisualToken: "BLACK_CENTER_SQUARE", SemanticRole: "PAGE_OR_MAJOR_CONTAINER", Dimension: "D01_GEOMETRY"},
	{VisualToken: "WHITE_SQUARE", SemanticRole: "CONCEPT", Dimension: "D01_GEOMETRY"},
	{VisualToken: "WHITE_DIAMOND", SemanticRole: "SOURCE_OR_EVIDENCE_ORIGIN", Dimension: "D01_GEOMETRY"},
	{VisualToken: "SMALL_BLACK_SQUARE", SemanticRole: "GRAPH_OR_RELATION_CLUSTER", Dimension: "D01_GEOMETRY"},
	{VisualToken: "LINE_OR_ARROW", SemanticRole: "DECLARED_RELATION", Dimension: "D04_TOPOLOGY"},
	{VisualToken: "ENCLOSURE_BOX", SemanticRole: "SPACE_SCOPE_OR_ADDRESS_REGION", Dimension: "D05_ENCLOSURE"},
	{VisualToken: "BLACK_CELL", SemanticRole: "VISUAL_BIT_1", Dimension: "D02_FILL_CONTRAST"},
	{VisualToken: "WHITE_CELL", SemanticRole: "VISUAL_BIT_0", Dimension: "D02_FILL_CONTRAST"},
	{VisualToken: "MACRO_MESO_MICRO_SCALE", SemanticRole: "REPRESENTATION_LEVEL", Dimension: "D06_SCALE"},
	{VisualToken: "DUPLICATED_PROBE", SemanticRole: "REDUNDANCY_AND_AGREEMENT", Dimension: "D07_REPETITION"},
	{VisualToken: "DENSE_REGION", SemanticRole: "DATA_PLANE_OR_HIGH_INFORMATION_DENSITY", Dimension: "D08_DENSITY"},
	{VisualToken: "SMALL_DECLARATIVE_TEXT", SemanticRole: "BOOT_ROSETTA_LABEL", Dimension: "D09_TEXT"},
}

var hardInvariants = []string{
	"WRITE_USES_CANONICAL_VISUAL_PROFILE",
	"ROSETTA_IS_REQUIRED_IN_READ_AND_WRITE",
	"MODEL_DOES_NOT_DEFINE_A_NEW_AESTHETIC_PER_DOCUMENT",
	"SOURCE_SCREENSHOT_NE_ORIGAMI",
	"CONSTRUCTION_SPEC_PRECEDES_COMPILATION",
	"COMPILER_NOT_MODEL_PIXELS_IS_CARRIER_AUTHORITY",
	"ROUNDTRIP_REQUIRED",
	"FALSE_EXACT=0",
}

func Plan(request Request) (ConstructionPlan, error) {
	if request.Schema != "" && request.Schema != RequestSchemaR0 {
		return ConstructionPlan{}, fmt.Errorf("request schema must be %q", RequestSchemaR0)
	}
	if strings.TrimSpace(request.SourceIdentity) == "" {
		return ConstructionPlan{}, fmt.Errorf("source_identity is required")
	}
	folded, _, err := core.FoldGraph(request.SemanticGraph)
	if err != nil {
		return ConstructionPlan{}, fmt.Errorf("semantic graph: %w", err)
	}
	canonical := canonicalGraphCopy(request.SemanticGraph)
	intents, err := deriveVisualIntents(canonical, request.IntentHints)
	if err != nil {
		return ConstructionPlan{}, err
	}
	active := activeDimensions(intents)
	rosetta := Rosetta{
		ProfileID:        CanonicalProfileR0,
		ProfileVersion:   "r0",
		ActiveDimensions: active,
		Entries:          append([]RosettaEntry(nil), canonicalRosetta...),
	}
	sort.Slice(rosetta.Entries, func(i, j int) bool {
		if rosetta.Entries[i].Dimension == rosetta.Entries[j].Dimension {
			return rosetta.Entries[i].VisualToken < rosetta.Entries[j].VisualToken
		}
		return rosetta.Entries[i].Dimension < rosetta.Entries[j].Dimension
	})
	return ConstructionPlan{
		Schema:             PlanSchemaR0,
		ProfileID:          CanonicalProfileR0,
		SourceIdentity:     strings.TrimSpace(request.SourceIdentity),
		SemanticCommitment: folded.Commitment,
		VisualIntents:      intents,
		Rosetta:            rosetta,
		Sections: []Section{
			{ID: "BOOT", Required: true, Purpose: "entry point, profile/version discovery and section addresses"},
			{ID: "ROSETTA", Required: true, Purpose: "declare canonical profile and active visual-to-semantic bindings"},
			{ID: "PROGRAM", Required: true, Purpose: "deterministic micro-machine/rule execution"},
			{ID: "INDEX", Required: true, Purpose: "bounded addressable navigation"},
			{ID: "MEMORY", Required: true, Purpose: "semantic/addressable payload references"},
			{ID: "VERIFICATION", Required: true, Purpose: "proof, StoreRoot/CID/hash and roundtrip checks"},
		},
		Verification: VerificationPlan{
			SemanticRoundtrip: true,
			ExactFields:       []string{"profile_id", "source_identity", "semantic_commitment", "addresses", "hashes", "verification_refs"},
			RequiredChecks: []string{
				"DECODED_PROFILE_MATCHES_PLAN",
				"ROSETTA_PRESENT_AND_PROFILE_BOUND",
				"SEMANTIC_ROUNDTRIP_EQUIVALENT",
				"EXACT_FIELDS_MATCH",
				"FALSE_EXACT_EQUALS_ZERO",
			},
			FailureState: "INVALID_CARRIER",
		},
		CanonicalGraph: canonical,
		Invariants:     append([]string(nil), hardInvariants...),
	}, nil
}

func deriveVisualIntents(graph core.SemanticGraph, hints []IntentHint) ([]VisualIntent, error) {
	byAddress := map[string][]Intent{}
	addresses := make([]string, 0, len(graph.Nodes))
	for address, node := range graph.Nodes {
		addresses = append(addresses, address)
		byAddress[address] = append(byAddress[address], IntentIdentity, IntentAddressability, IntentState)
		if len(node.Dependencies) > 0 {
			byAddress[address] = append(byAddress[address], IntentDependency)
		}
		if len(node.Alternatives) > 0 || node.Value.Status == core.StatusUnknown {
			byAddress[address] = append(byAddress[address], IntentUncertainty)
		}
		if len(node.Relations) > 0 {
			byAddress[address] = append(byAddress[address], IntentRelation)
		}
		if len(node.EvidenceRefs) > 0 {
			byAddress[address] = append(byAddress[address], IntentEvidence)
		}
		if node.Verified {
			byAddress[address] = append(byAddress[address], IntentVerification, IntentRedundancy)
		}
	}
	for _, hint := range hints {
		if _, ok := graph.Nodes[hint.Address]; !ok {
			return nil, fmt.Errorf("intent hint address %q is not in semantic graph", hint.Address)
		}
		if !validIntent(hint.Intent) {
			return nil, fmt.Errorf("unsupported intent %q", hint.Intent)
		}
		byAddress[hint.Address] = append(byAddress[hint.Address], hint.Intent)
	}
	sort.Strings(addresses)
	var out []VisualIntent
	for _, address := range addresses {
		intents := dedupeIntents(byAddress[address])
		for _, intent := range intents {
			dimension, primitive, role := bindingFor(intent, graph.Nodes[address])
			out = append(out, VisualIntent{
				Address:      address,
				Intent:       intent,
				Dimension:    dimension,
				Primitive:    primitive,
				SemanticRole: role,
				Authority:    "CANONICAL_PROFILE_R0",
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Address == out[j].Address {
			return out[i].Intent < out[j].Intent
		}
		return out[i].Address < out[j].Address
	})
	return out, nil
}

func bindingFor(intent Intent, node core.SemanticNode) (string, string, string) {
	switch intent {
	case IntentIdentity:
		primitive := "WHITE_SQUARE"
		role := "CONCEPT"
		switch strings.ToLower(node.Kind) {
		case "source", "evidence", "source-block":
			primitive, role = "WHITE_DIAMOND", "SOURCE_OR_EVIDENCE_ORIGIN"
		case "page", "container", "space":
			primitive, role = "BLACK_CENTER_SQUARE", "PAGE_OR_MAJOR_CONTAINER"
		case "graph", "relation", "cluster":
			primitive, role = "SMALL_BLACK_SQUARE", "GRAPH_OR_RELATION_CLUSTER"
		case "root", "verify", "verification":
			primitive, role = "RINGED_SQUARE", "ROOT_OR_VERIFY_ANCHOR"
		}
		return "D01_GEOMETRY", primitive, role
	case IntentHierarchy:
		return "D03_POSITION", "VERTICAL_OR_TREE_POSITION", "HIERARCHY"
	case IntentRelation, IntentDependency:
		return "D04_TOPOLOGY", "LINE_OR_ARROW", string(intent)
	case IntentScope:
		return "D05_ENCLOSURE", "ENCLOSURE_BOX", "SCOPE"
	case IntentImportance:
		return "D06_SCALE", "RELATIVE_SCALE", "IMPORTANCE"
	case IntentState:
		return "D02_FILL_CONTRAST", "FILL_OR_CONTRAST_STATE", "SEMANTIC_STATE"
	case IntentUncertainty:
		return "D02_FILL_CONTRAST", "DECLARED_UNCERTAINTY_MARK", "UNCERTAINTY"
	case IntentEvidence:
		return "D01_GEOMETRY", "WHITE_DIAMOND", "EVIDENCE_ORIGIN"
	case IntentVerification:
		return "D07_REPETITION", "DUPLICATED_OR_RINGED_VERIFY_MARK", "VERIFICATION"
	case IntentAddressability:
		return "D05_ENCLOSURE", "ADDRESS_REGION", "ADDRESSABILITY"
	case IntentTemporalOrder:
		// R0 does not make temporal graphics primary authority. The writer uses a
		// canonical positional sequence until a temporal channel is promoted.
		return "D03_POSITION", "ORDERED_SEQUENCE", "TEMPORAL_ORDER_AS_STATIC_SEQUENCE"
	case IntentRedundancy:
		return "D07_REPETITION", "DUPLICATED_MARK", "REDUNDANCY"
	case IntentRouting:
		return "D03_POSITION", "ROUTE_POSITION", "ROUTING"
	default:
		return "D01_GEOMETRY", "WHITE_SQUARE", "UNSPECIFIED"
	}
}

func validIntent(intent Intent) bool {
	switch intent {
	case IntentIdentity, IntentHierarchy, IntentRelation, IntentDependency, IntentScope,
		IntentImportance, IntentState, IntentUncertainty, IntentEvidence, IntentVerification,
		IntentAddressability, IntentTemporalOrder, IntentRedundancy, IntentRouting:
		return true
	default:
		return false
	}
}

func dedupeIntents(values []Intent) []Intent {
	seen := map[Intent]struct{}{}
	for _, value := range values {
		seen[value] = struct{}{}
	}
	out := make([]Intent, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func activeDimensions(intents []VisualIntent) []string {
	seen := map[string]struct{}{"D09_TEXT": {}}
	for _, intent := range intents {
		seen[intent.Dimension] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for dimension := range seen {
		out = append(out, dimension)
	}
	sort.Strings(out)
	return out
}

func canonicalGraphCopy(graph core.SemanticGraph) core.SemanticGraph {
	out := core.SemanticGraph{Schema: core.SchemaR1 + ".semantic-graph", Nodes: map[string]core.SemanticNode{}}
	for address, node := range graph.Nodes {
		copyNode := node
		copyNode.Dependencies = append([]string(nil), node.Dependencies...)
		copyNode.Alternatives = append([]string(nil), node.Alternatives...)
		copyNode.EvidenceRefs = append([]string(nil), node.EvidenceRefs...)
		sort.Strings(copyNode.Dependencies)
		sort.Strings(copyNode.Alternatives)
		sort.Strings(copyNode.EvidenceRefs)
		if node.Relations != nil {
			copyNode.Relations = map[string][]string{}
			for relation, targets := range node.Relations {
				copyTargets := append([]string(nil), targets...)
				sort.Strings(copyTargets)
				copyNode.Relations[relation] = copyTargets
			}
		}
		out.Nodes[address] = copyNode
	}
	return out
}
