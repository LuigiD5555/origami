# Changelog

## Unreleased

### Reversible installation lifecycle

- added `uninstall.sh` with tracked restore semantics;
- `install.sh` now records SHA-256 installation state and preserves the binaries that existed before the first tracked Origami install;
- reinstalls preserve the original rollback point instead of backing up a previous Origami binary as if it were user state;
- uninstall refuses to overwrite/remove binaries changed after installation unless `--force` is explicit;
- added opt-in legacy uninstall support for the previous `.bak`-based installer;
- normal smoke tests are now ephemeral and leave no persistent `runs/install-smoke` directory.


This file summarizes major repository-level milestones. Detailed experimental state belongs in `state/` and `changes/`.


## 2026-08-27 — Dimensional Visual Register / PSS formalization

- added hierarchical MICRO/MESO/MACRO visual-dimension registry;
- added executable DVR validation and scope-aware capacity reporting;
- separated nominal PSS from evidence-gated `SAFE_PSS`;
- added Perceptual Orthogonality Law and pair-evidence model;
- integrated Macro-Gestalt as the MACRO scale of the same dimensional system;
- added machine-readable state and `CHG-DVR-0001`;
- added local inspection script with zero external requests.

## 2026-08-27 — R3.10-LAB repository modernization

- moved the public repository from R2 documentation-only status to an executable Go laboratory;
- added reproducible F1 laboratory infrastructure;
- added Native isolation, observations, scoring and failure-regression plumbing;
- added Glyph Calculus;
- added SAFE_MICRO_ISA and Context SIMD experiments;
- added conservative external-execution policy;
- added installer, Makefile and CI;
- changed the Go module path to `github.com/LuigiD5555/origami`;
- replaced the R2-centric README with current R3.10-LAB documentation;
- retained R2 prompt as explicit legacy material;
- added machine-readable state/change records and current-state documentation.

## 2026-08-26 — R2 public documentation baseline

- added scientific foundations document;
- added R2 Master Prompt;
- added practical R2 README.
# Origami changelog

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
- runtime support for these operations is not yet complete and unsupported behavior must be explicit.

### Preserved

- Tlaloc/Origami project boundary from `6.0.0-alpha.1`;
- coherent-state profile semantics;
- historical OHF R2 material as non-authoritative research history.

## 6.0.0-alpha.1 — State Semantics R0

### Added / formalized

- `determinate` state;
- `superposed` state;
- `coupled` joint state;
- `observed` state;
- complex branch amplitudes for the quantum-inspired profile;
- `TRANSFORM` without implicit observation;
- `INTERFERE` and cancellation;
- `CONSTRAIN`;
- explicit `OBSERVE`;
- `FOLD` / `UNFOLD` relationship to state resolution/expansion.

### Removed from Origami ownership

The following concepts belong to Tlaloc and must not be versioned as Origami features:

- Behavior Compiler;
- PromptIR;
- Tlaloque coordination/training;
- target-model orchestration;
- behavior-prompt promotion arbiter.
