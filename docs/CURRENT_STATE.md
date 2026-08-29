# Current Origami state

**Status date:** 2026-08-29  
**Released project version:** 6.0.0-alpha.3  
**Development track:** generative-machine recovery  
**Status:** experimental

## Root project state

Origami is a relational, generative, visual/computational representation and state-machine language. Tlaloc is a separate optional work/orchestration system.

The current recovery contract restores five core responsibilities: `FOLD`, `RELATE`, `EVOLVE`, `OBSERVE`, and `UNFOLD`. The formal state/observation laws remain authoritative; the generative Fold/Unfold runtime described by `docs/ORIGAMI_MODEL.md` is being implemented incrementally and must not be reported as complete until its executable gates pass.

Current layers:

| Area | State |
|---|---|
| formal relational machine | EXPERIMENTAL_CONTRACT |
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels | EXPERIMENTAL_CONTRACT — alpha.2 |
| Temporal Latent Image semantics | EXPERIMENTAL_CONTRACT |
| generative Fold/Unfold recovery model | RECOVERY_CONTRACT — implementation pending |
| stable addressing/selective unfolding core | NOT IMPLEMENTED |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| OHF R3.10-LAB research track | EXPERIMENTAL / executable |

The root machine-readable source is `state/ORIGAMI_STATE.json`. OHF-specific component state remains under `state/components/`.

## Recovery invariant

```text
RASTER_IS_NOT_FOLD
```

A screenshot or rendered PDF page may be an ingestion artifact, projection, perception fixture or control baseline. Rasterization alone is not evidence of Origami generative Fold and must not be promoted as such.

The minimum recovery path is:

```text
source
  ↓
canonical structure
  ↓
non-raster Fold
  ↓
generative / relational representation
  ↓
stable addressing
  ↓
selective Unfold
  ↓
deterministic verification
```

PDF ingestion and visual carriers are downstream experiments. Initial implementation should use small deterministic text/graph fixtures.

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

These OHF invariants do not make a raster carrier equivalent to a successful Origami Fold.

## Origami 6.x perceptual state

`origami.perceptual-channels.r0` defines spatial, interference, depth, temporal and emergent channels. It includes Temporal Latent Images and the invariant that a motion-bound or otherwise latent percept is not absent merely because one static observation does not reveal it.

Contract awareness is not runtime implementation. Individual operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE` and `TEMPORAL_DECAY` require explicit implementation/evidence before being reported as executable.

## Source of truth

Use:

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
docs/ORIGAMI_MODEL.md
docs/FORMAL_CORE_R0.md
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/PERCEPTUAL_CHANNELS_R0.json
changes/
state/components/   # OHF-specific implementation evidence
```

Historical R2/R3.10 materials are evidence and continuity inputs, not authority over the current project boundary.
