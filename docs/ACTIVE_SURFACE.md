# Origami active surface

This is the short maintained map of the repository. Git history is the archive; old per-alpha change-control notes are not part of the active root.

## Canonical project boundary

Origami owns representation, addressable state, visual/computational carriers, self-description, codecs, temporal semantics and exact verification. Tlaloc may experiment against Origami, but does not own Origami pixel or protocol authority.

## Active protocol/runtime surface

- `internal/core`, `internal/document`, `internal/memory` — semantic/addressable state and virtual memory.
- `internal/codec` — declared read/write semantic and exact codecs.
- `internal/fixedcarrier` — fixed self-contained carrier profiles and deterministic exact decode.
- `internal/temporal`, `internal/automaton` — temporal program and graph-automaton runtime.
- `internal/receiver` — embedded BOOT/ROSETTA/micro-program receiver contract.
- `internal/writer` — canonical construction planning.
- `internal/promotion` and `internal/evidence` — evidence gates and authority boundary.

## Active laboratory surface

- `internal/lab/microisa` — bounded visual micro-operations and safe-width evidence.
- `internal/lab/behaviorimage` — **experimental visible behavior recipes** for testing whether tiny VLMs can follow precompiled microinstructions from the carrier.
- `experiments/behavior-in-image-r0/` — behavior-recipe fixtures; no model success is claimed by their presence.

## Compatibility retained deliberately

Profile-1 decoding and other compatibility paths remain only where tests/current contracts still exercise them. Compatibility code is not automatically deleted merely because it is called legacy.

## Removed from the active surface

Per-alpha root change-control files are deleted once their durable information is folded into `CHANGELOG.md`, current state, tests and Git history.
