# Origami 6.0.0-alpha.2

Origami is an experimental **visual/computational representation and state-machine language** for expressing complex states, relations, dynamics and perceptual emergence.

Origami is **not** the complete work system. Orchestration, behavior compilation, Tlaloque coordination and model-facing workflows belong to **Tlaloc**.

## Project boundary

```text
TLALOC                          ORIGAMI
work system                     representation
orchestration                   state language
behavior compilation   ----->   semantics contracts
Tlaloque coordination  ----->   reference semantics / projection contracts
verification                     dynamics / perceptual channels
```

Origami may be used independently of Tlaloc. Tlaloc may also operate without Origami or with another representation provider.

> Tlaloc coordinates work; Tlaloque perform bounded specialist work; Origami defines how a class of states is represented, transformed and made perceptually available.

## Semantic layers

Origami 6.0 currently separates two orthogonal semantic layers.

### Coherent-state profile

Defined in [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md):

- `determinate`
- `superposed`
- `coupled`
- `observed`

This profile is quantum-inspired as a computational analogy. It is not a claim of physical quantum computation.

### Perceptual channels R0

New in `6.0.0-alpha.2`, defined in [`docs/PERCEPTUAL_CHANNELS_R0.md`](docs/PERCEPTUAL_CHANNELS_R0.md):

- spatial channels;
- interference / moiré channels;
- depth / stereoscopic / parallax channels;
- temporal channels;
- emergent channels.

A represented structure does not need to exist as a complete static image. Some structures may become available only through interaction among layers or through a declared observation trajectory over time.

## Temporal Latent Image

A **Temporal Latent Image (TLI)** is a representation whose intended percept is unavailable or incomplete under a valid static observation and becomes available only through a declared temporal transformation.

```text
carrier + layer relation + trajectory + temporal integration
                         -> percept
```

The motion/trajectory is part of the representation contract, not incidental animation.

The corresponding operations include:

- `MOIRE`
- `PHASE_SHIFT`
- `STEREO_BIND`
- `PARALLAX_RESOLVE`
- `KINETIC_REVEAL`
- `TEMPORAL_INTEGRATE`
- `TEMPORAL_DECAY`

Machine-readable form: [`spec/PERCEPTUAL_CHANNELS_R0.json`](spec/PERCEPTUAL_CHANNELS_R0.json).

## Important distinction

Origami keeps state semantics and perceptual availability separate.

```text
state semantics != perceptual availability
```

For example, a `superposed` state may have a `motion-bound` projection; an `observed` state may still be encoded as a Temporal Latent Image.

Perceptual conditions such as `latent`, `motion-bound`, `phase-bound`, `stereo-bound`, `revealed` and `decayed` are therefore annotations on availability, not replacements for coherent-state kinds.

## Fold / Unfold

- `UNFOLD` constructs or expands alternatives under a declared representation rule.
- `FOLD` constrains and/or explicitly resolves according to a declared policy.

These remain Origami operations. Model training, prompt compilation and behavior enforcement remain Tlaloc responsibilities.

## Runtime status

`6.0.0-alpha.2` formalizes semantic contracts. It does **not** yet claim complete production renderers/detectors/decoders for every perceptual operation.

Unsupported operations must be reported explicitly rather than approximated silently.

The exact pre-split Origami 5.1 implementation remains the authoritative implementation base until these 6.0 contracts are integrated and regression-tested against it.

## Historical continuity

`6.0.0-alpha.1` established the Tlaloc/Origami project split and coherent-state semantic overlay.

`6.0.0-alpha.2` adds perceptual-channel semantics while preserving that boundary.

Earlier **Origami HyperFold (OHF) R2** material remains archived as historical research input:

- `docs/SCIENTIFIC_FOUNDATIONS.md`
- `prompts/OHF_R2_MASTER_PROMPT.txt`
- `docs/LEGACY_OHF_R2.md`

Those files do not define the current project boundary.

## Repository layout

```text
origami/
├── README.md
├── VERSION
├── CHANGELOG.md
├── PROJECT_BOUNDARY.md
├── CHANGE_CONTROL_6.0.0-alpha.1.md
├── CHANGE_CONTROL_6.0.0-alpha.2.md
├── docs/
│   ├── STATE_SEMANTICS_R0.md
│   ├── PERCEPTUAL_CHANNELS_R0.md
│   ├── LEGACY_OHF_R2.md
│   └── SCIENTIFIC_FOUNDATIONS.md      # historical OHF material
├── spec/
│   └── PERCEPTUAL_CHANNELS_R0.json
└── prompts/
    └── OHF_R2_MASTER_PROMPT.txt        # historical OHF R2 artifact
```

## Version

`6.0.0-alpha.2`
