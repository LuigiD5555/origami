package codec

import (
	"testing"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func temporalFixture() temporal.Program {
	a := automaton.Automaton{
		Schema: automaton.Schema, ID: "a",
		Cells: []automaton.Cell{{ID: "A", InitialState: "ACTIVE"}, {ID: "B", InitialState: "IDLE"}},
		Rules: []automaton.Rule{{ID: "r1", TargetCell: "B", FromState: "IDLE", ToState: "ACTIVE", Requires: []automaton.Predicate{{CellID: "A", State: "ACTIVE"}}}},
	}
	return temporal.Program{Schema: temporal.Schema, ID: "p", Automaton: a, MaxSteps: 4, CheckpointEvery: 2}
}

func TestTemporalRegistryValid(t *testing.T) {
	r := TemporalRegistry()
	if err := r.Validate(); err != nil { t.Fatal(err) }
	entry, err := r.Select("READ_TIMELINE", SemanticDecoder, map[string]bool{"SEMANTIC_READ": true, "TEMPORAL_READ": true})
	if err != nil { t.Fatal(err) }
	if entry.ID != "ST2" { t.Fatalf("expected ST2, got %s", entry.ID) }
	if _, err := r.Select("READ_TIMELINE", SemanticDecoder, map[string]bool{"SEMANTIC_READ": true}); err == nil {
		t.Fatal("expected missing temporal capability")
	}
}

func TestAutomatonRoundtrip(t *testing.T) {
	p := temporalFixture()
	ir, err := EncodeAutomaton(p.Automaton, "profile")
	if err != nil { t.Fatal(err) }
	out, err := DecodeAutomaton(ir)
	if err != nil { t.Fatal(err) }
	if !EqualAutomaton(p.Automaton, out) { t.Fatal("automaton roundtrip drift") }
}

func TestTimelineRoundtripReplay(t *testing.T) {
	p := temporalFixture()
	ir, err := EncodeTimeline(p, "profile")
	if err != nil { t.Fatal(err) }
	out, err := DecodeTimeline(ir)
	if err != nil { t.Fatal(err) }
	trace, err := temporal.Run(out)
	if err != nil { t.Fatal(err) }
	if err := temporal.VerifyReplay(p, trace); err != nil { t.Fatal(err) }
}
