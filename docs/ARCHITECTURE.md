# Origami architecture

Origami is a relational, generative, visual/computational representation and state-machine language. Its architecture separates source ingestion, generative representation, semantic state, dynamics, observation, selective unfolding and optional projection/carrier implementations.

```text
SOURCE
  ↓
CANONICAL STRUCTURE
  ↓
FOLD
  ↓
ORIGAMI REPRESENTATION
├── state / relations
├── literal / reference / repetition
├── rules / graph / grammar / transforms
├── residual / verification
└── addresses / index
  ↓
RELATE / MACHINE DYNAMICS
  ↓
OBSERVATION + QUERY/ADDRESS RESOLUTION
  ↓
SELECTIVE UNFOLD
  ↓
RESULT + PROVENANCE + EVIDENCE
  ↓ optional
VISUAL / COMPUTATIONAL PROJECTION
  ↓
PROFILE OR CARRIER
  └── OHF research track
```

## Fold/Unfold layer

`FOLD` is not rasterization. It produces an explicit representation that can reuse structure through references, repetition, transforms, rules, graphs or other promoted mechanisms while retaining residual information where exactness requires it.

`UNFOLD` reconstructs a requested region of logical structure. The default engineering goal is selective unfolding: resolve an address/query, compute the bounded dependency closure and materialize only what is required. Full exact unfolding is a separate capability and must be verified when claimed.

`RASTER_IS_NOT_FOLD` is a recovery invariant. Rendering a PDF/page/document to pixels may participate in ingestion, projection or perception experiments, but rasterization alone cannot satisfy a gate for generative Fold.

## Semantic layers

Coherent-state kinds (`determinate`, `superposed`, `coupled`, `observed`) are distinct from perceptual availability conditions (`latent`, `motion-bound`, `phase-bound`, `stereo-bound`, `revealed`, `decayed`).

A representation may therefore have a valid semantic state while its intended percept is available only through interference, depth, motion or temporal integration.

State transition remains separate from observation. The deterministic R0 model is defined by `docs/FORMAL_CORE_R0.md`; recovery work must extend executable representation without silently changing that law.

## Addressability

Origami representations require stable addresses for meaningful nodes/regions. An index or equivalent resolver maps declared identifiers or queries to representation addresses and their dependency closure. This is what permits selective unfolding and avoids treating the entire carrier/source as one opaque image.

Addressability claims must be measured. Useful evidence includes total logical expansion, node count, nodes visited, nodes unfolded, bytes materialized, residual bytes and verification outcome.

## Perceptual channels

The current R0 contract defines spatial, interference, depth, temporal and emergent channels. Emergent meaning may belong to an interaction rather than a single source layer. A Temporal Latent Image may require a declared trajectory to become perceptually available.

The contract defines meaning and validity conditions; it does not imply every operation already has a production renderer/detector.

## Projection boundary

Projection is downstream of the semantic representation. One Origami state may have multiple projections: structured data, static visual output, temporal output, an OHF carrier or another computational surface. A projection does not become semantic source of truth merely because a VLM can inspect it.

## OHF research profile

OHF remains the main executable carrier/laboratory track currently in this repository. Its architecture explores compact addressable carriers, selective unfolding, deterministic verification and bounded model perception.

Representative recovered relationship:

```text
source
  ↓
canonical structure
  ↓
Origami Fold / representation tournament
  ↓
graph / grammar / rules / references / residual
  ↓
index + selective-unfold contract
  ↓
visual compiler / perceptual ABI
  ↓
OHF carrier
  ↓
perception evidence
  ↓
resolution
  ↓
bounded Go execution
  ↓
verification
```

DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA and Context SIMD are experimental OHF mechanisms. They may inform future general Origami machinery, but they are not automatically universal language semantics.

## Tlaloc boundary

Tlaloc owns work orchestration, behavior compilation, Tlaloque coordination, model-facing execution/evaluation and promotion control. Origami owns representation semantics. Tlaloc may consume Origami contracts without becoming part of the Origami runtime.

## Verification discipline

For representation changes, record before/after state, evidence, regressions, downstream impact and promotion decision. Mathematical/generative capacity must remain separate from evidence-gated perceptual capacity. Unsupported behavior must be explicit rather than guessed.
