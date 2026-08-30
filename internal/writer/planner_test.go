package writer

import (
	"reflect"
	"testing"

	"github.com/LuigiD5555/origami/internal/core"
)

func sampleGraph() core.SemanticGraph {
	return core.SemanticGraph{Nodes: map[string]core.SemanticNode{
		"root": {
			Address: "root", Kind: "root", Label: "Document root",
			Value: core.Value{Status: core.StatusPresent, Data: "paper"},
			Dependencies: []string{"concept"}, Verified: true,
		},
		"concept": {
			Address: "concept", Kind: "concept", Label: "Dynamic programming",
			Value: core.Value{Status: core.StatusUnknown},
			Alternatives: []string{"source"},
			EvidenceRefs: []string{"ohf://paper/page-1/block-2"},
		},
		"source": {
			Address: "source", Kind: "source", Label: "Source block",
			Value: core.Value{Status: core.StatusPresent, Data: "evidence"},
			Relations: map[string][]string{"supports": {"concept"}},
			Verified: true,
		},
	}}
}

func TestPlanUsesOneCanonicalProfile(t *testing.T) {
	plan, err := Plan(Request{Schema: RequestSchemaR0, SourceIdentity: "paper.pdf", SemanticGraph: sampleGraph()})
	if err != nil { t.Fatal(err) }
	if plan.ProfileID != CanonicalProfileR0 || plan.Rosetta.ProfileID != CanonicalProfileR0 {
		t.Fatalf("writer drifted from canonical profile: %+v", plan)
	}
	if len(plan.Rosetta.Entries) == 0 || len(plan.Rosetta.ActiveDimensions) == 0 {
		t.Fatalf("Rosetta must be generated: %+v", plan.Rosetta)
	}
	if !plan.Verification.SemanticRoundtrip || plan.Verification.FailureState != "INVALID_CARRIER" {
		t.Fatalf("roundtrip gate missing: %+v", plan.Verification)
	}
}

func TestPlanDeterministicAcrossMapOrder(t *testing.T) {
	graphA := sampleGraph()
	graphB := core.SemanticGraph{Nodes: map[string]core.SemanticNode{}}
	graphB.Nodes["source"] = graphA.Nodes["source"]
	graphB.Nodes["root"] = graphA.Nodes["root"]
	graphB.Nodes["concept"] = graphA.Nodes["concept"]
	a, err := Plan(Request{SourceIdentity: "paper.pdf", SemanticGraph: graphA}); if err != nil { t.Fatal(err) }
	b, err := Plan(Request{SourceIdentity: "paper.pdf", SemanticGraph: graphB}); if err != nil { t.Fatal(err) }
	if a.SemanticCommitment != b.SemanticCommitment || !reflect.DeepEqual(a.VisualIntents, b.VisualIntents) || !reflect.DeepEqual(a.Rosetta, b.Rosetta) {
		t.Fatalf("construction plan is not deterministic\nA=%+v\nB=%+v", a, b)
	}
}

func TestSourceNodeUsesCanonicalEvidenceGeometry(t *testing.T) {
	plan, err := Plan(Request{SourceIdentity: "paper.pdf", SemanticGraph: sampleGraph()}); if err != nil { t.Fatal(err) }
	found := false
	for _, intent := range plan.VisualIntents {
		if intent.Address == "source" && intent.Intent == IntentIdentity {
			found = true
			if intent.Dimension != "D01_GEOMETRY" || intent.Primitive != "WHITE_DIAMOND" {
				t.Fatalf("source node escaped canonical visual grammar: %+v", intent)
			}
		}
	}
	if !found { t.Fatal("missing source identity intent") }
}

func TestTemporalIntentUsesStaticCanonicalFallbackUntilTemporalChannelPromoted(t *testing.T) {
	plan, err := Plan(Request{
		SourceIdentity: "sequence",
		SemanticGraph: sampleGraph(),
		IntentHints: []IntentHint{{Address: "concept", Intent: IntentTemporalOrder}},
	}); if err != nil { t.Fatal(err) }
	for _, intent := range plan.VisualIntents {
		if intent.Address == "concept" && intent.Intent == IntentTemporalOrder {
			if intent.Dimension != "D03_POSITION" || intent.SemanticRole != "TEMPORAL_ORDER_AS_STATIC_SEQUENCE" {
				t.Fatalf("unpromoted temporal dimension became primary authority: %+v", intent)
			}
			return
		}
	}
	t.Fatal("temporal intent missing")
}

func TestUnknownHintAddressRejected(t *testing.T) {
	_, err := Plan(Request{SourceIdentity: "paper.pdf", SemanticGraph: sampleGraph(), IntentHints: []IntentHint{{Address: "missing", Intent: IntentRouting}}})
	if err == nil { t.Fatal("unknown address hint must fail") }
}

func TestArbitraryIntentRejected(t *testing.T) {
	_, err := Plan(Request{SourceIdentity: "paper.pdf", SemanticGraph: sampleGraph(), IntentHints: []IntentHint{{Address: "root", Intent: Intent("MAKE_IT_PRETTY")}}})
	if err == nil { t.Fatal("arbitrary aesthetic intent must fail") }
}
