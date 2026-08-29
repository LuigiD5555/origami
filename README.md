# Origami 6.0.0-alpha.3

Origami is an experimental **relational, generative, visual/computational representation and state-machine language** for folding complex information, relations, rules and dynamics into addressable executable representations that can be selectively unfolded and verified.

A screenshot, rasterized PDF page or visual carrier is **not itself the definition of Origami**. Visual projection is one possible surface of an underlying Origami representation. The recovery model is defined in [`docs/ORIGAMI_MODEL.md`](docs/ORIGAMI_MODEL.md).

Origami is independent of **Tlaloc**. Tlaloc is the work/orchestration system; Origami defines representation and transformation semantics and may be used with Tlaloc or independently.

## Core model

Origami is organized around five core verbs:

```text
FOLD -> RELATE -> EVOLVE / OBSERVE -> UNFOLD
```

- `FOLD` creates an explicit generative representation from source structure.
- `RELATE` preserves addressable dependencies, hierarchy, references and constraints.
- `EVOLVE` applies declared rules and operations to machine state.
- `OBSERVE` evaluates direct or latent properties through finite observation contracts.
- `UNFOLD` materializes only the required dependency closure when possible, with exact full recovery when a profile explicitly guarantees it.

The recovery invariant `RASTER_IS_NOT_FOLD` prevents a rendered page or screenshot from being accepted as evidence that generative Fold occurred.

## Project hierarchy

```text
ORIGAMI
├── Generative Fold / Unfold
│   ├── literals / references / repetition
│   ├── transforms / rules / graphs
│   ├── residual / verification
│   └── addressing / selective unfolding
├── State semantics
│   └── coherent-state profile
├── Perceptual channels
│   ├── spatial
│   ├── interference / moiré
│   ├── depth / stereo / parallax
│   ├── temporal
│   └── emergent
├── Machine / dynamics / projection
└── Research profiles and carriers
    └── OHF — Origami HyperFold
        ├── R3.10 / R3.10-LAB workbench
        ├── DVR / PSS
        ├── Glyph Calculus
        ├── SAFE_MICRO_ISA
        └── Context SIMD
```

**OHF is an Origami research track, not the name or complete architecture of Origami.** Its R3.10-LAB Go code remains useful and executable, but OHF-specific carrier goals and Native experiments do not redefine the whole language.

Detailed pre-reconciliation OHF documentation is preserved under [`docs/ohf/`](docs/ohf/README.md), rather than being discarded.

## Current semantic contracts

### Formal machine core

[`docs/FORMAL_CORE_R0.md`](docs/FORMAL_CORE_R0.md) defines deterministic state transition and finite observation discipline. The recovery work extends that formal direction with explicit generative Fold/Unfold, addressability and selective materialization without changing the existing state/observation laws.

### Coherent-state profile

Defined in [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md): `determinate`, `superposed`, `coupled`, `observed`, explicit observation authority, coherent interference/cancellation and Fold/Unfold semantics.

The profile is quantum-inspired as a computational analogy, not a claim of physical quantum computation.

### Perceptual Channels R0

Defined in [`docs/PERCEPTUAL_CHANNELS_R0.md`](docs/PERCEPTUAL_CHANNELS_R0.md) and [`spec/PERCEPTUAL_CHANNELS_R0.json`](spec/PERCEPTUAL_CHANNELS_R0.json).

It formalizes spatial, interference, depth, temporal and emergent channels, including **Temporal Latent Images (TLI)** and operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL` and `TEMPORAL_INTEGRATE`.

A valid Origami percept may exist only through a declared relation or trajectory; failure to reveal it in one static frame is not proof of absence.

## Executable OHF research track

The repository also contains the experimental Go workbench inherited from OHF R3.10-LAB. It currently includes deterministic experiment identity, truth-first fixtures, artifact hashing, Native isolation, observation/scoring, failure-to-regression plumbing, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible installation for the OHF laboratory CLIs.

These implementations are **experimental profile/tooling**, not proof that every Origami 6.x semantic operation has a renderer, detector or reference executor. Likewise, successful raster/VLM perception is not evidence of generative Fold unless the underlying representation gates pass.

The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` remains missing; that blocks byte-identical historical R3.10 baseline claims, but it does not block independent Origami 6.x semantic development or OHF-LAB experiments.

## Source of truth and change control

Current project identity and state must be recoverable from the repository:

```text
VERSION
README.md
PROJECT_BOUNDARY.md
GATEKEEPER.md
gatekeeper.json
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/ORIGAMI_MODEL.md
docs/FORMAL_CORE_R0.md
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/PERCEPTUAL_CHANNELS_R0.json
state/ORIGAMI_STATE.json
state/components/
changes/
CHANGE_CONTROL_*.md
```

Historical documents and OHF R2/R3.10 artifacts remain evidence, but they do not override the current project boundary.

## Project Gatekeeper

Origami follows the project-wide Gatekeeper R0 owned by Tonal. The local `gatekeeper.json` is a CI mirror, not a second policy authority.

- owner PR (`LuigiD5555` from this canonical repository): Origami CI and evidence still run; explicit owner promotion override is allowed;
- external PR: CI/evidence still run and an `APPROVED` review from `LuigiD5555` is mandatory; external override/auto-promotion is denied.

See `GATEKEEPER.md`. The canonical operational `gatekeeper` skill lives in Tonal alongside `repo-flow`.

## Development

Cheap local gates first:

```bash
go test ./...
go vet ./...
```

The generative-machine recovery starts with small deterministic text/graph fixtures. PDF ingestion and visual carriers come only after Fold/Unfold, addressing and verification are executable; PDF is an ingestion format, not the definition of Origami.

For OHF dimensional/laboratory work, run the affected scripts and regression gates documented under `docs/`, `docs/ohf/` and `scripts/`. External model campaigns remain progressive, cached and conservative.

`install.sh` installs the **OHF laboratory command-line tools** (`ohf-lab`, `ohf-glyphcalc`) under the user prefix. It is not a global installer for every future Origami runtime/profile.

## Version

`6.0.0-alpha.3`

This remains the current released version while the generative-machine recovery is developed on a feature branch. No recovery capability is promoted merely by documenting its intended contract.
