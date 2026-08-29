# Origami changelog

## Unreleased

### Added

- added the explicit `toggles` relation to the deterministic reference engine so bounded oscillations can be classified as `CYCLE` without changing existing relation semantics;
- added EXP-002 hostile-cycle evidence and a gate requiring at least one observed cycle.

## 6.0.0-alpha.3 — Repository hierarchy and state reconciliation

### Corrected

- reconciled the executable OHF R3.10-LAB workbench with the Origami 6.x project boundary;
- made the repository hierarchy explicit: Origami is the representation/state-machine language; OHF is a research carrier/protocol track inside Origami;
- replaced the stale OHF-only `state/ORIGAMI_STATE.json` root model with a project-level state that tracks semantic contracts and the OHF substate separately;
- repaired README and documentation that described R3.10/OHF as the current state or architecture of all Origami;
- clarified that `install.sh` / `uninstall.sh` manage OHF laboratory CLIs, not the complete future Origami runtime;
- preserved the detailed pre-reconciliation OHF R3.10-LAB architecture/state/development/installation/roadmap documents under `docs/ohf/`;
- preserved R3.10-LAB, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible-installation work as active experimental assets rather than discarding them.

### Preserved

- `origami.perceptual-channels.r0` semantics from alpha.2;
- coherent-state semantics from alpha.1;
- Tlaloc/Origami ownership boundary;
- historical R2/R3.10 evidence and change records;
- existing Go experimental behavior.

## 6.0.0-alpha.2 — Perceptual Channels R0

### Added

- perceptual-channel contract `origami.perceptual-channels.r0`;
- spatial, interference, depth, temporal and emergent channel families;
- explicit-vs-emergent channel distinction;
- Temporal Latent Image (TLI) semantic object;
- perceptual availability conditions: `latent`, `masked`, `motion-bound`, `phase-bound`, `stereo-bound`, `emergent`, `revealed`, `decayed`;
- operations `MOIRE`, `PHASE_SHIFT`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE`, `TEMPORAL_DECAY`;
- machine-readable `spec/PERCEPTUAL_CHANNELS_R0.json`;
- invariants preventing static-frame assumptions and false absence claims.

### Clarified

- coherent-state kinds and perceptual conditions are separate semantic layers;
- a valid Origami percept may exist only across a declared temporal trajectory;
- an emergent percept belongs to the interaction that produces it, not automatically to one source layer;
- runtime support for these operations is incomplete and unsupported behavior must be explicit.

## 6.0.0-alpha.1 — State Semantics R0

### Added / formalized

- `determinate`, `superposed`, `coupled` and `observed` states;
- complex branch amplitudes for the quantum-inspired profile;
- `TRANSFORM` without implicit observation;
- `INTERFERE` and cancellation;
- `CONSTRAIN`;
- explicit `OBSERVE`;
- `FOLD` / `UNFOLD` relationship to state resolution/expansion.

### Removed from Origami ownership

Behavior compilation, PromptIR, Tlaloque coordination/training, target-model orchestration and behavioral-artifact promotion belong to Tlaloc.

## Historical OHF / R3.10 laboratory milestones

### 2026-08-27 — Reversible installation lifecycle

- added tracked SHA-256 restore semantics for OHF laboratory binaries;
- preserved the original pre-install rollback point across reinstalls;
- added conflict detection and opt-in legacy `.bak` removal support;
- kept normal smoke runs ephemeral.

### 2026-08-27 — Dimensional Visual Register / PSS formalization

- added hierarchical MICRO/MESO/MACRO visual-dimension registry;
- added executable DVR validation and scope-aware capacity reporting;
- separated nominal PSS from evidence-gated `SAFE_PSS`;
- added Perceptual Orthogonality Law and pair-evidence model;
- integrated Macro-Gestalt as a MACRO-scale experimental candidate;
- added `CHG-DVR-0001` and local inspection gates.

### 2026-08-27 — R3.10-LAB repository modernization

- moved the public OHF research track from R2 documentation-only status to an executable Go laboratory;
- added reproducible experiment infrastructure, Native isolation, observations/scoring and failure-regression plumbing;
- added Glyph Calculus, SAFE_MICRO_ISA and Context SIMD experiments;
- added conservative external-execution policy, installer, Makefile and CI;
- changed the Go module path to `github.com/LuigiD5555/origami`;
- retained R2 material as explicit legacy evidence;
- added machine-readable OHF/LAB state and change records.

### 2026-08-26 — R2 public documentation baseline

- added scientific foundations document;
- added R2 Master Prompt;
- added practical R2 README.
