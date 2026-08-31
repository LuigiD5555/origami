package fixedcarrier

import (
	"fmt"
	"image"
	"strings"
)

// drawTemporalRuleMicrogrammar renders an experimental semantic rule plane for
// one-row temporal carriers. It does not alter the exact TemporalProgram plane.
// The goal is to make causal rule preconditions and synchronous step semantics
// visually explicit enough for a VLM to simulate the declared automaton without
// inferring dynamics from graph topology alone.
func drawTemporalRuleMicrogrammar(img *image.Gray, decoded TemporalCarrierDecoded) {
	drawTemporalRuleMicrogrammarWithCellLabels(img, decoded, func(id string) string { return id })
}

// drawTemporalRuleMicrogrammarWithCellLabels keeps rule semantics sourced from
// structured IR while allowing an experimental renderer to change only the
// visible identity label of a cell. The canonical cell IDs inside TemporalProgram
// are never rewritten.
func drawTemporalRuleMicrogrammarWithCellLabels(img *image.Gray, decoded TemporalCarrierDecoded, label func(string) string) {
	p := decoded.Program

	// Preserve dense two-row T2 layouts instead of overwriting cells. The current
	// signal-chain benchmark uses one row; denser programs retain a compact marker
	// and the ordinary checkpoint emphasis until a dedicated dense grammar exists.
	if len(p.Automaton.Cells) > 4 {
		drawText(img, 24, 300, 1, "RULE GRAMMAR: DENSE T2 FALLBACK", 0)
		emphasizeTemporalCheckpoints(img, decoded)
		return
	}

	// The one-row semantic graph ends at y=196 and the timeline begins at y=326.
	// Use the otherwise empty band as a bounded rule microgrammar.
	fill(img, 18, 202, 604, 112, 0xff)
	box(img, 18, 202, 604, 112, 0)
	drawText(img, 24, 208, 1, "T2 RULE MICROGRAMMAR | RULES DRIVE STATE", 0)
	drawText(img, 24, 220, 1, "EACH STEP: TEST ALL RULES ON SAME PRE-STEP SNAPSHOT", 0)

	y := 234
	visible := 0
	for _, r := range p.Automaton.Rules {
		if visible >= 6 {
			break
		}
		req := "TRUE"
		if len(r.Requires) > 0 {
			parts := make([]string, 0, len(r.Requires))
			for _, pred := range r.Requires {
				parts = append(parts, shortLabel(label(pred.CellID), 16)+"="+shortLabel(pred.State, 10))
			}
			req = strings.Join(parts, "&")
		}
		from := r.FromState
		if strings.TrimSpace(from) == "" {
			from = "*"
		}
		line := fmt.Sprintf("IF %s => %s:%s>%s", req, shortLabel(label(r.TargetCell), 16), shortLabel(from, 10), shortLabel(r.ToState, 10))
		drawText(img, 28, y, 1, shortLabel(line, 94), 0)
		y += 12
		visible++
	}

	if len(p.Automaton.Rules) > visible {
		drawText(img, 28, y, 1, fmt.Sprintf("+%d MORE RULES IN EXACT PROGRAM", len(p.Automaton.Rules)-visible), 0)
	}
}
