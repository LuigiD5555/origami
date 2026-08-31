package fixedcarrier

import "image"

// drawExecuteToStableDirective adds only execution-policy guidance. It does not
// author or rewrite automaton rules; those remain sourced from structured IR and
// the existing temporal rule microgrammar.
func drawExecuteToStableDirective(img *image.Gray) {
	// Keep the rule microgrammar band (y=202..314) untouched. Place a compact
	// directive immediately above the timeline.
	fill(img, 360, 286, 252, 34, 0xff)
	box(img, 360, 286, 252, 34, 0)
	drawText(img, 366, 292, 1, "EXECUTE: INIT > APPLY ALL SAME PRE-STEP > NEXT", 0)
	drawText(img, 366, 304, 1, "REPEAT UNTIL UNCHANGED > REPORT STABLE STATE", 0)
}
