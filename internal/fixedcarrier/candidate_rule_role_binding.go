package fixedcarrier

import (
	"fmt"
	"image"
	"strings"
)

const VisibleRuleRoleBindingR1 = "VISIBLE_RULE_ROLE_BINDING_R1"

// drawVisibleRuleRoleBindingR1 separates rule operands by semantic role so a
// receiver does not need to infer which cell is the source condition, which is
// the target, or which target state is required before applying the SET.
func drawVisibleRuleRoleBindingR1(img *image.Gray, decoded TemporalCarrierDecoded) {
	labels := temporalVisibleCellLabels(decoded, true)
	p := decoded.Program
	if len(p.Automaton.Cells) > 4 {
		drawText(img, 24, 300, 1, "RULE ROLE BINDING: DENSE T2 FALLBACK", 0)
		emphasizeTemporalCheckpoints(img, decoded)
		return
	}

	fill(img, 18, 202, 604, 112, 0xff)
	box(img, 18, 202, 604, 112, 0)
	drawText(img, 24, 208, 1, "T2 RULE MICROGRAMMAR | EXPLICIT OPERAND ROLES", 0)
	drawText(img, 24, 220, 1, "EACH STEP: TEST ALL CONDITIONS ON SAME PRE-STEP SNAPSHOT", 0)

	y := 234
	visible := 0
	for _, r := range p.Automaton.Rules {
		if visible >= 4 { break }
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
		drawText(img, 28, y, 1, shortLabel(fmt.Sprintf("%s WHEN %s", strings.ToUpper(r.ID), req), 94), 0)
		y += 10
		drawText(img, 36, y, 1, shortLabel(fmt.Sprintf("TARGET %s | REQUIRE %s | SET %s", shortLabel(labels[r.TargetCell], 16), shortLabel(from, 10), shortLabel(r.ToState, 10)), 94), 0)
		y += 12
		visible++
	}
	if len(p.Automaton.Rules) > visible {
		drawText(img, 28, y, 1, fmt.Sprintf("+%d MORE RULES IN EXACT PROGRAM", len(p.Automaton.Rules)-visible), 0)
	}
}
