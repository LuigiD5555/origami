# Origami Capacity Lab R0

Primary question: how much knowledge can Origami make usable while the answering working set remains bounded?

Default target: `W = 4000 token-equivalent`.

## Required counters

Every run records REQUESTED, TOUCHED/VISITED, UNFOLDED, EXPOSED, active bytes, token-equivalent, physical depth and cognitive depth. TOUCHED, UNFOLDED and EXPOSED are distinct counters and must not be collapsed.

The answering model may depend only on EXPOSED material. Access to the complete source outside that set contaminates the trial.

## Gates

A scale is promoted only when route integrity, answer correctness, evidence integrity, budget compliance, no-global-scan and determinism pass.

Missing information returns UNKNOWN. A miss must never trigger a hidden global scan.

## Campaign

- ACC: W = 256, 512, 1024, 2048, 4000, 8000.
- DEP: dependency depth, comparing canonical traversal with Dependency Folding.
- BRN: branching factors and selective branch resolution.
- HOR: higher-order interactions where the answer depends on multiple relations.
- TMP: temporal distance/span.
- NAV: multi-step LOOKUP/FOLLOW/EXPAND/UNFOLD/TRACE/VERIFY/BACKTRACK/STOP behavior.
- NEG: absent facts and unresolved guards must produce UNKNOWN.

Synthetic memory scale proceeds progressively: 1K, 10K, 100K, 1M, then 10M only after the previous gate passes.

## Dependency Folding accounting

A folded result does not erase physical depth. Reports carry both physical and cognitive depth. Conditions, exceptions and verification material exposed to the answerer count against W.

If folding cannot preserve semantics across a boundary, refine only that region. Never flatten an irreducible interaction merely to satisfy the window budget.
