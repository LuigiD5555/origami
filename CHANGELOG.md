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
