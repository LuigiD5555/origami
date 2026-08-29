# Origami architecture

Origami is a visual/computational representation and state-machine language. Its architecture separates semantic state, perceptual availability, dynamics and profile/carrier implementations.

```text
STATE / RELATIONS
      ↓
ORIGAMI SEMANTICS
  ├── coherent-state profile
  └── perceptual-channel contracts
      ↓
MACHINE / DYNAMICS
      ↓
VISUAL / COMPUTATIONAL PROJECTION
      ↓
PROFILE OR CARRIER
  └── OHF research track
```

## Semantic layers

Coherent-state kinds (`determinate`, `superposed`, `coupled`, `observed`) are distinct from perceptual availability conditions (`latent`, `motion-bound`, `phase-bound`, `stereo-bound`, `revealed`, `decayed`).

A representation may therefore have a valid semantic state while its intended percept is available only through interference, depth, motion or temporal integration.

## Perceptual channels

The current R0 contract defines spatial, interference, depth, temporal and emergent channels. Emergent meaning may belong to an interaction rather than a single source layer. A Temporal Latent Image may require a declared trajectory to become perceptually available.

The contract defines meaning and validity conditions; it does not imply every operation already has a production renderer/detector.

## OHF research profile

OHF remains the main executable carrier/laboratory track currently in this repository. Its architecture explores compact addressable carriers, selective unfolding, deterministic verification and bounded model perception.

Representative OHF pipeline:

```text
source
  ↓
structure / grammar / graph / transforms
  ↓
representation candidate
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

For representation changes, record before/after state, evidence, regressions, downstream impact and promotion decision. Mathematical capacity must remain separate from evidence-gated perceptual capacity. Unsupported perceptual behavior must be explicit rather than guessed.
