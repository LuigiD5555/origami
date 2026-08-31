package fixedcarrier

import (
	"fmt"
	"image"
)

const VisibleCellIdentityRedundancyR1 = "VISIBLE_CELL_ID_REDUNDANCY_R1"

// drawVisibleCellIdentityRedundancyR1 changes only the visible identity of cells.
// Canonical IDs, states, rules, execution policy, checkpoints and exact payload
// remain untouched. Labels are derived deterministically from structured cell
// order: A -> A[01], B -> B[02], ... .
func drawVisibleCellIdentityRedundancyR1(img *image.Gray, decoded TemporalCarrierDecoded) {
	labels := temporalVisibleCellLabels(decoded, true)
	cells := decoded.Program.Automaton.Cells
	if len(cells) > 8 {
		cells = cells[:8]
	}
	for i, c := range cells {
		x := 34 + (i%4)*145
		y := 148 + (i/4)*82
		// Replace only the first textual row inside the existing cell box.
		fill(img, x+5, y+5, 100, 16, 0xff)
		drawText(img, x+8, y+8, 1, "CELL "+shortLabel(labels[c.ID], 16), 0)
	}

	// R1 visible rules are regenerated from the exact TemporalProgram using the
	// same grammar renderer and only the redundant cell labels. This avoids any
	// free-form semantic rewrite of preconditions or transitions.
	drawTemporalRuleMicrogrammarWithCellLabels(img, decoded, func(id string) string {
		if v, ok := labels[id]; ok {
			return v
		}
		return id
	})
}

func temporalVisibleCellLabels(decoded TemporalCarrierDecoded, redundant bool) map[string]string {
	out := make(map[string]string, len(decoded.Program.Automaton.Cells))
	for i, c := range decoded.Program.Automaton.Cells {
		if redundant {
			out[c.ID] = fmt.Sprintf("%s[%02d]", c.ID, i+1)
		} else {
			out[c.ID] = c.ID
		}
	}
	return out
}
