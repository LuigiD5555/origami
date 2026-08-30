package fixedcarrier

import (
	"encoding/json"
	"testing"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func temporalCarrierFixture() temporal.Program {
	return temporal.Program{
		Schema: temporal.Schema, ID: "film",
		Automaton: automaton.Automaton{
			Schema: automaton.Schema, ID: "chain",
			Cells: []automaton.Cell{
				{ID:"A", InitialState:"ACTIVE", Neighbors:[]string{"B"}},
				{ID:"B", InitialState:"IDLE", Neighbors:[]string{"A","C"}},
				{ID:"C", InitialState:"IDLE", Neighbors:[]string{"B"}},
			},
			Rules: []automaton.Rule{
				{ID:"r1",TargetCell:"B",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"A",State:"ACTIVE"}}},
				{ID:"r2",TargetCell:"A",FromState:"ACTIVE",ToState:"DONE",Requires:[]automaton.Predicate{{CellID:"B",State:"ACTIVE"}}},
				{ID:"r3",TargetCell:"C",FromState:"IDLE",ToState:"ACTIVE",Requires:[]automaton.Predicate{{CellID:"B",State:"ACTIVE"}}},
			},
		},
		MaxSteps: 6, CheckpointEvery: 2,
	}
}

func TestTemporalCarrierSelfContainedRoundtrip(t *testing.T) {
	p := temporalCarrierFixture()
	pngData, built, err := RenderTemporalCarrier(p)
	if err != nil { t.Fatal(err) }
	if len(pngData) != FixedPNGBytes { t.Fatalf("expected %d bytes, got %d", FixedPNGBytes, len(pngData)) }
	if !built.ExactProgramInPNG || built.VisibleCellCount != 3 { t.Fatalf("bad build metadata: %#v", built) }
	decoded, err := DecodeTemporalCarrierPNG(pngData)
	if err != nil { t.Fatal(err) }
	a, _ := json.Marshal(p)
	b, _ := json.Marshal(decoded.Program)
	if string(a) != string(b) { t.Fatal("program roundtrip drift") }
	trace, err := temporal.Run(decoded.Program); if err != nil { t.Fatal(err) }
	if err := temporal.VerifyReplay(p, trace); err != nil { t.Fatal(err) }
}

func TestTemporalCarrierRejectsOversizeProgram(t *testing.T) {
	p := temporalCarrierFixture()
	for i:=0;i<200;i++ {
		p.Automaton.Cells = append(p.Automaton.Cells, automaton.Cell{ID: "cell-long-unique-id-"+string(rune(0x1000+i)), InitialState:"STATE-WITH-UNIQUE-CONTENT"})
	}
	if _, _, err := RenderTemporalCarrier(p); err == nil { t.Fatal("expected oversize program failure") }
}
