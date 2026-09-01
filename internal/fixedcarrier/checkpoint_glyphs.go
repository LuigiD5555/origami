package fixedcarrier

import (
	"image"
	"strings"

	"github.com/LuigiD5555/origami/internal/automaton"
	"github.com/LuigiD5555/origami/internal/temporal"
)

// stateToGlyph maps a cell's textual state to a compact one-character symbol.
// Symbols chosen to fit within the font5x7 bitmap font used by drawText.
func stateToGlyph(state string) string {
	switch strings.ToUpper(state) {
	case "IDLE":
		return "O"
	case "ACTIVE":
		return "*"
	case "DONE":
		return "X"
	default:
		// Unmapped states fall back to '?' (already in font).
		return "?"
	}
}

// drawCheckpointStateRow renders one row of per-cell state glyphs at a given
// checkpoint step. Each cell's current state is displayed as a compact glyph
// so a reader can see the frozen computation instead of inferring it.
// Format: "T{step}" + glyph-per-cell separated by spaces.
// Positioned below the timeline but above the exact-plane grid (y >= 375, y < 415).
func drawCheckpointStateRow(img *image.Gray, baseX, baseY int, step int, cellOrder []string, state automaton.State) {
	// Convert step to decimal string inline (avoid fmt import).
	stepStr := "T"
	remaining := step
	if remaining >= 100 {
		stepStr += string(rune('0' + remaining/100))
		remaining = remaining % 100
	}
	if remaining >= 10 {
		stepStr += string(rune('0' + remaining/10))
		remaining = remaining % 10
	}
	stepStr += string(rune('0' + remaining))

	// Render label and glyphs on a single line.
	// scale=1 for minimal byte footprint; each character is 6 pixels wide + 1px spacing.
	drawText(img, baseX, baseY, 1, stepStr, 0)
	xOffset := baseX + textWidth(stepStr, 1) + 2 // 2px gap after label

	for _, cellID := range cellOrder {
		cellState := state[cellID]
		glyph := stateToGlyph(cellState)
		drawText(img, xOffset, baseY, 1, glyph, 0)
		xOffset += textWidth(glyph, 1) + 1 // 1px spacing between glyphs
	}
}

// drawCheckpointStates renders a frozen trace: for each checkpoint in the
// trace, draw its state snapshot as a compact row of glyphs below the
// timeline. This allows readers to verify the program's evolution without
// simulating it themselves.
func drawCheckpointStates(img *image.Gray, trace *temporal.Trace) {
	if trace == nil || len(trace.Checkpoints) == 0 {
		return
	}

	// Position rows below the timeline (which ends ~y=372) and well above the
	// exact-plane grid (which starts at y0=420). Use y=375 to y=413 (39px span),
	// allowing ~5 rows at scale=1 (each row is ~8px high), with 5px margin to grid.
	const baseY = 375
	const rowHeight = 8

	// Get stable cell order (alphabetically sorted).
	cellOrder := automaton.SortedCellIDs(trace.Initial)

	// Draw one row per checkpoint.
	for i, checkpoint := range trace.Checkpoints {
		y := baseY + i*rowHeight
		if y+8 >= 415 { // Don't draw past the margin to the grid.
			break
		}
		drawCheckpointStateRow(img, 38, y, checkpoint.Step, cellOrder, checkpoint.State)
	}
}
