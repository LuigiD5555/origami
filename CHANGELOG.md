# Origami changelog

This file records the maintained/current release line. Older per-alpha implementation notes remain available in Git history instead of occupying the active repository surface.

## 6.0.0-alpha.15 — Candidate compilation, causal trace foundations, behavior-in-image lab

### Active additions

- `origami.experimental-candidate.r0` and `origami-candidate-build` for bounded representation mutations proposed by external experiment systems;
- mechanical parent/candidate exact-program preservation inside the fixed 640×640 / exactly-8192-byte temporal carrier envelope;
- explicit mutation-capability negotiation and unsupported-mutation failure rather than silent fallback;
- causal temporal trace foundations preserving parallel frontiers and rejecting cyclic or unknown-parent traces;
- decoder training/inference adapters for the SAFE_MICRO_ISA research track;
- experimental `origami.behavior-recipe.r0` plus `origami-behavior-image`, rendering bounded forward-only visible microinstructions with fail-closed `UNKNOWN` routing.

### Evidence boundary

- candidate compilation does not imply model improvement or canonical promotion;
- deterministic behavior-recipe rendering does not prove that LFM2-VL or any other VLM can execute those instructions;
- Tlaloc may run experiments and collect evidence, while Origami remains authority for representation/protocol promotion.

## 6.0.0-alpha.14 — Temporal Graph Automata R0

### Added

- deterministic semantic graph automata and temporal programs using rules, deltas and checkpoints rather than literal video frames;
- temporal graph navigation and declared temporal codec families;
- deterministic run/seek/replay verification;
- `origami-temporal-carrier`, a self-contained fixed-envelope PNG carrying visible semantic temporal structure plus an exact verified TemporalProgram.

### Evidence boundary

- semantic film is not literal video;
- unsupported temporal capability returns `UNKNOWN` / `NOT_VERIFIED`;
- deterministic runtime success is not held-out Native VLM evidence.

## 6.0.0-alpha.13 — Self-describing read/write protocol R0

### Added

- `origami.protocol.r0`: Origami is a bidirectional self-describing communication protocol, not only a visual memory representation;
- `origami.codec-registry.r0` with semantic decoder/encoder families `S0..S6` / `E0..E6` and exact/control families `X0..X5` / `Q0..Q5`;
- `origami.capability-negotiation.r0`, making unsupported work fail explicitly with `UNKNOWN` / `NOT_VERIFIED`;
- ROSETTA-declared decoding and encoding semantics;
- first semantic roundtrip gate `S2(E2(INDEX)) ~= INDEX`;
- protocol-aware Master Prompt R4 reference candidate;
- experimental profile-3 renderer with visible protocol/codec entry points while profile-2 remains the default.

### Corrected interpretation

- self-declared semantic codecs are valid Origami behavior;
- undeclared external decoder/file/binary dependencies remain forbidden for Native semantic navigation;
- semantic queries should not escalate to exact/binary decode unless exactness is actually required;
- construction specification is not a compiled/verified carrier.

### Evidence boundary

Held-out Native S2 recovery, held-out E2 writing, and real cross-model A→B→C interoperability remain evidence-pending.

## Historical releases

Alpha.1–alpha.12 and the older OHF/R2/R3.10 laboratory milestones remain permanently recoverable from Git history. Their durable contracts that are still active are represented by current `docs/`, `spec/`, `state/`, tests, and machine-readable change/claim ledgers; obsolete narrative snapshots are intentionally not duplicated in the active root.
