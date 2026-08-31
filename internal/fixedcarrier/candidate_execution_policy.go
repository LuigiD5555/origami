package fixedcarrier

import "image"

// drawExecuteToStableDirective adds only execution-policy guidance. It does not
// author or rewrite automaton rules; those remain sourced from structured IR and
// the existing temporal rule microgrammar.
func drawExecuteToStableDirective(img *image.Gray) {
	// Keep the semantic graph and rule microgrammar untouched. The canonical BOOT
	// leaves a free final row at y=92 for experimental protocol guidance.
	fill(img, 20, 88, 596, 14, 0xff)
	drawText(img, 24, 92, 1, "EXEC: INIT > APPLY ALL SAME PRE-STEP > NEXT > REPEAT UNTIL UNCHANGED > REPORT STABLE", 0)
}
