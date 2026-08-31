package fixedcarrier

import (
	"fmt"
	"image"
	"strings"
)

const VisibleFromStatePreconditionR1 = "VISIBLE_FROM_STATE_PRECONDITION_R1"

// drawVisibleFromStatePreconditionR1 changes only the visible rule grammar.
// Canonical TemporalProgram semantics, cell identities and execution policy are
// preserved. FROM is rendered as an explicit target-state precondition instead
// of compact transition punctuation.
func drawVisibleFromStatePreconditionR1(img *image.Gray, decoded TemporalCarrierDecoded) {
	labels := temporalVisibleCellLabels(decoded, true)
	p := decoded.Program
	if len(p.Automaton.Cells) > 4 {
		drawText(img, 24, 300, 1, "FROM PRECONDITION: DENSE T2 FALLBACK", 0)
		emphasizeTemporalCheckpoints(img, decoded)
		return
	}

	fill(img, 18, 202, 604, 112, 0xff)
	box(img, 18, 202, 604, 112, 0)
	drawText(img, 24, 208, 1, "T2 RULE MICROGRAMMAR | FROM IS REQUIRED", 0)
	drawText(img, 24, 220, 1, "EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT", 0)

	y := 234
	visible := 0
	for _, r := range p.Automaton.Rules {
		if visible >= 6 { break }
		req := "TRUE"
		if len(r.Requires) > 0 {
			parts := make([]string, 0, len(r.Requires))
			for _, pred := range r.Requires {
				parts = append(parts, shortLabel(labels[pred.CellID], 16)+"="+shortLabel(pred.State, 10))
			}
			req = strings.Join(parts, "&")
		}
		from := r.FromState
		if strings.TrimSpace(from) == "" { from = "*" }
		line := fmt.Sprintf("IF %s AND %s=%s THEN %s -> %s",
			req,
			shortLabel(labels[r.TargetCell], 16), shortLabel(from, 10),
			shortLabel(labels[r.TargetCell], 16), shortLabel(r.ToState, 10))
		drawText(img, 28, y, 1, shortLabel(line, 94), 0)
		y += 12
		visible++
	}
	if len(p.Automaton.Rules) > visible {
		drawText(img, 28, y, 1, fmt.Sprintf("+%d MORE RULES IN EXACT PROGRAM", len(p.Automaton.Rules)-visible), 0)
	}
}
