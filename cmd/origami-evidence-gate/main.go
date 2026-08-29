package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/reference"
)

func main() {
	experimentPath := flag.String("experiment", "experiments/EXP-001-relational-state/experiment.json", "experiment specification")
	evidencePath := flag.String("evidence", "evidence/EXP-001-reference.json", "reference evidence")
	flag.Parse()

	e, err := reference.Load(*experimentPath)
	if err != nil { fatal(err) }
	b, err := os.ReadFile(*evidencePath)
	if err != nil { fatal(err) }
	var ev reference.Evidence
	if err := json.Unmarshal(b, &ev); err != nil { fatal(err) }

	want := 1
	for range e.Entities { want *= len(e.Domain) }
	if ev.Schema != "origami.reference-evidence.r0" { fatal(fmt.Errorf("unexpected evidence schema %q", ev.Schema)) }
	if ev.ExperimentID != e.ID { fatal(fmt.Errorf("experiment mismatch: evidence=%q spec=%q", ev.ExperimentID, e.ID)) }
	if !ev.Deterministic { fatal(fmt.Errorf("evidence is not marked deterministic")) }
	if len(ev.Runs) != want { fatal(fmt.Errorf("incomplete state space: runs=%d want=%d", len(ev.Runs), want)) }

	allowed := map[string]bool{"FIXED_POINT":true,"CYCLE":true,"CONTRADICTION":true,"BUDGET_EXHAUSTED":true}
	counts := map[string]int{}
	for i,r := range ev.Runs {
		if !allowed[r.Classification] { fatal(fmt.Errorf("run %d has invalid classification %q", i, r.Classification)) }
		if len(r.Trace) == 0 { fatal(fmt.Errorf("run %d has empty trace", i)) }
		if len(r.Trace) > e.Protocol.MaxSteps+1 { fatal(fmt.Errorf("run %d exceeds trace budget: %d > %d", i, len(r.Trace), e.Protocol.MaxSteps+1)) }
		counts[r.Classification]++
	}
	for k,v := range counts { if ev.Counts[k] != v { fatal(fmt.Errorf("count mismatch for %s: evidence=%d computed=%d", k, ev.Counts[k], v)) } }
	for k,v := range ev.Counts { if counts[k] != v { fatal(fmt.Errorf("unexpected count for %s: evidence=%d computed=%d", k, v, counts[k])) } }

	fmt.Printf("evidence gate PASS: experiment=%s runs=%d", e.ID, len(ev.Runs))
	for _,k := range []string{"FIXED_POINT","CYCLE","CONTRADICTION","BUDGET_EXHAUSTED"} { fmt.Printf(" %s=%d", k, counts[k]) }
	fmt.Println()
}

func fatal(err error) { fmt.Fprintln(os.Stderr, "origami-evidence-gate:", err); os.Exit(1) }
