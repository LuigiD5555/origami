package codec

import (
	"fmt"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/temporal"
)

// TemporalRegistry extends the protocol registry without making temporal
// support mandatory for Core semantic compatibility.
func TemporalRegistry() Registry {
	r := ReferenceRegistry()
	r.Entries = append(r.Entries,
		temporalEntry("ST0", SemanticDecoder, "READ_AUTOMATON", "ET0", "UNKNOWN", "TEMPORAL_READ", "READ_DECLARED_CELLS_RULES_GRAPH"),
		temporalEntry("ET0", SemanticEncoder, "ENCODE_AUTOMATON", "ST0", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_AUTOMATON_IR"),
		temporalEntry("ST1", SemanticDecoder, "READ_CELL", "ET1", "UNKNOWN", "TEMPORAL_READ", "LOCATE_CELL", "READ_LOCAL_STATE_AND_NEIGHBORS"),
		temporalEntry("ET1", SemanticEncoder, "ENCODE_CELL", "ST1", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_CELL_IR"),
		temporalEntry("ST2", SemanticDecoder, "READ_TIMELINE", "ET2", "UNKNOWN", "TEMPORAL_READ", "LOCATE_TIMELINE", "READ_CHECKPOINTS_AND_DELTAS"),
		temporalEntry("ET2", SemanticEncoder, "ENCODE_TIMELINE", "ST2", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_TEMPORAL_PROGRAM_IR"),
		temporalEntry("ST3", SemanticDecoder, "LOCATE_EVENT", "ET3", "UNKNOWN", "TEMPORAL_READ", "ROUTE_EVENT_THROUGH_T2"),
		temporalEntry("ET3", SemanticEncoder, "ENCODE_EVENT", "ST3", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_EVENT_ADDRESS"),
		temporalEntry("ST4", SemanticDecoder, "READ_TRANSITION", "ET4", "UNKNOWN", "TEMPORAL_READ", "READ_RULE_AND_DELTA"),
		temporalEntry("ET4", SemanticEncoder, "ENCODE_TRANSITION", "ST4", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_DECLARED_RULE_AND_DELTA"),
		temporalEntry("ST5", SemanticDecoder, "UNFOLD_TEMPORAL_REGION", "ET5", "UNKNOWN", "TEMPORAL_READ", "SEEK_NEAREST_CHECKPOINT", "APPLY_BOUNDED_DELTAS"),
		temporalEntry("ET5", SemanticEncoder, "FOLD_TEMPORAL_REGION", "ST5", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "FACTOR_RULES_DELTAS_CHECKPOINTS"),
		temporalEntry("ST6", SemanticDecoder, "SIMULATE_DECLARED_STEP", "ET6", "UNKNOWN", "TEMPORAL_READ", "APPLY_DECLARED_SYNCHRONOUS_STEP"),
		temporalEntry("ET6", SemanticEncoder, "ENCODE_CHECKPOINT", "ST6", "CONSTRUCTION_SPEC_ONLY", "TEMPORAL_WRITE", "EMIT_CHECKPOINT"),
	)
	return r
}

func temporalEntry(id string, family Family, operation, pair, failure, temporalCap string, steps ...string) Entry {
	baseCap := "SEMANTIC_READ"
	if family == SemanticEncoder {
		baseCap = "SEMANTIC_WRITE"
	}
	return Entry{
		ID: id, Family: family, Operation: operation, Pair: pair, Exact: false,
		RequiredCapabilities: []string{baseCap, temporalCap}, FailureState: failure, Steps: steps,
	}
}

type AutomatonConstructionIR struct {
	Schema    string              `json:"schema"`
	ProfileID string              `json:"profile_id,omitempty"`
	CodecID   string              `json:"codec_id"`
	Automaton automaton.Automaton `json:"automaton"`
}

type TimelineConstructionIR struct {
	Schema    string           `json:"schema"`
	ProfileID string           `json:"profile_id,omitempty"`
	CodecID   string           `json:"codec_id"`
	Program   temporal.Program `json:"program"`
}

func EncodeAutomaton(a automaton.Automaton, profileID string) (AutomatonConstructionIR, error) {
	if err := a.Validate(); err != nil {
		return AutomatonConstructionIR{}, err
	}
	return AutomatonConstructionIR{Schema: "origami.automaton-construction-ir.r0", ProfileID: profileID, CodecID: "ET0", Automaton: a}, nil
}

func DecodeAutomaton(ir AutomatonConstructionIR) (automaton.Automaton, error) {
	if ir.Schema != "origami.automaton-construction-ir.r0" || ir.CodecID != "ET0" {
		return automaton.Automaton{}, fmt.Errorf("unsupported automaton construction IR")
	}
	if err := ir.Automaton.Validate(); err != nil {
		return automaton.Automaton{}, err
	}
	return ir.Automaton, nil
}

func EncodeTimeline(p temporal.Program, profileID string) (TimelineConstructionIR, error) {
	if err := p.Validate(); err != nil {
		return TimelineConstructionIR{}, err
	}
	return TimelineConstructionIR{Schema: "origami.timeline-construction-ir.r0", ProfileID: profileID, CodecID: "ET2", Program: p}, nil
}

func DecodeTimeline(ir TimelineConstructionIR) (temporal.Program, error) {
	if ir.Schema != "origami.timeline-construction-ir.r0" || ir.CodecID != "ET2" {
		return temporal.Program{}, fmt.Errorf("unsupported timeline construction IR")
	}
	if err := ir.Program.Validate(); err != nil {
		return temporal.Program{}, err
	}
	return ir.Program, nil
}

func EqualAutomaton(a, b automaton.Automaton) bool {
	ia, err := EncodeAutomaton(a, "")
	if err != nil { return false }
	ib, err := EncodeAutomaton(b, "")
	if err != nil { return false }
	return fmt.Sprintf("%#v", ia.Automaton) == fmt.Sprintf("%#v", ib.Automaton)
}
