package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/temporal"
)

func main() {
	out := flag.String("out", "runs/comparison-baseline-r0", "output directory")
	flag.Parse()
	os.MkdirAll(*out, 0o755)

	program := temporal.Program{
		Schema: temporal.Schema,
		ID:     "film-r5-unseen",
		Automaton: automaton.Automaton{
			Schema: automaton.Schema,
			ID:     "reverse-chain",
			Cells: []automaton.Cell{
				{ID: "A", InitialState: "IDLE", Neighbors: []string{"B", "C"}},
				{ID: "B", InitialState: "IDLE", Neighbors: []string{"A", "C"}},
				{ID: "C", InitialState: "ACTIVE", Neighbors: []string{"A", "B"}},
			},
			Rules: []automaton.Rule{
				{ID: "r1", TargetCell: "A", FromState: "IDLE", ToState: "ACTIVE", Requires: []automaton.Predicate{{CellID: "C", State: "ACTIVE"}}},
				{ID: "r2", TargetCell: "B", FromState: "IDLE", ToState: "ACTIVE", Requires: []automaton.Predicate{{CellID: "A", State: "ACTIVE"}}},
				{ID: "r3", TargetCell: "C", FromState: "ACTIVE", ToState: "DONE", Requires: []automaton.Predicate{{CellID: "B", State: "ACTIVE"}}},
				{ID: "r4", TargetCell: "A", FromState: "ACTIVE", ToState: "DONE", Requires: []automaton.Predicate{{CellID: "C", State: "DONE"}}},
			},
		},
		MaxSteps:        7,
		CheckpointEvery: 2,
	}

	base, baseMeta, err := fixedcarrier.RenderTemporalCarrier(program)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(*out, "baseline-r0.png"), base, 0o644); err != nil {
		panic(err)
	}

	meta := map[string]any{
		"schema":            "origami.carrier-r0.r1",
		"id":                "baseline-r0",
		"program_sha256":    baseMeta.ProgramSHA256,
		"artifact_bytes":    baseMeta.PayloadBytes,
		"visible_cells":     baseMeta.VisibleCellCount,
		"visible_rules":     baseMeta.VisibleRuleCount,
		"exact_program_yes": baseMeta.ExactProgramInPNG,
	}
	b, _ := json.MarshalIndent(meta, "", "  ")
	os.WriteFile(filepath.Join(*out, "metadata.json"), append(b, '\n'), 0o644)

	fmt.Printf("R0_PNG=%s\nPROGRAM_SHA256=%s\nBYTES=%d\n", filepath.Join(*out, "baseline-r0.png"), baseMeta.ProgramSHA256, len(base))
}
