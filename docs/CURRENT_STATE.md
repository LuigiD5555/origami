# Current Origami state

**Status date:** 2026-08-28  
**Project version:** 6.0.0-alpha.3  
**Status:** experimental

## Root project state

Origami is the visual/computational representation and state-machine language. Tlaloc is a separate optional work/orchestration system.

Current semantic layers:

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels | EXPERIMENTAL_CONTRACT — alpha.2 |
| Temporal Latent Image semantics | EXPERIMENTAL_CONTRACT |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| OHF R3.10-LAB research track | EXPERIMENTAL / executable |

The root machine-readable source is `state/ORIGAMI_STATE.json`. OHF-specific component state remains under `state/components/`.

## OHF research track

OHF is a nested carrier/protocol track, not the complete project identity.

**Protocol reference:** R3.10  
**Executable laboratory:** R3.10-LAB  
**Canonical deterministic laboratory language:** Go

The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` is still missing. Therefore byte-identical historical R3.10 baseline claims remain blocked. This blocker is specific to historical OHF baseline recovery; it does not invalidate the Origami 6.x semantic contracts.

Current implemented/experimental OHF assets include deterministic experiment identity, truth-first fixtures, SHA-256 artifact verification, Native blind bundles, observations/scoring, failure-to-regression plumbing, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD, DVR/PSS, perceptual-orthogonality evidence and reversible laboratory-tool installation.

Important OHF invariants remain:

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
carrier target <= 500 KB
nominal PSS != SAFE_PSS
pairwise PASS != higher-order promotion
external tests are progressive/cached/rate-limited
```

## Origami 6.x perceptual state

`origami.perceptual-channels.r0` defines spatial, interference, depth, temporal and emergent channels. It includes Temporal Latent Images and the invariant that a motion-bound or otherwise latent percept is not absent merely because one static observation does not reveal it.

Contract awareness is not runtime implementation. Individual operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE` and `TEMPORAL_DECAY` require explicit implementation/evidence before being reported as executable.

## Source of truth

Use:

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/PERCEPTUAL_CHANNELS_R0.json
changes/
state/components/   # OHF-specific implementation evidence
```

Historical R2/R3.10 materials are evidence and continuity inputs, not authority over the current project boundary.
