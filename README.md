<p align="center">
  <img src="assets/origami-ant.webp" alt="Origami ant mark" width="360">
</p>

# Origami 6.0.0-alpha.3

Origami is an experimental **visual/computational representation and state-machine language** for expressing complex states, relations, dynamics and perceptual emergence.

Origami is independent of **Tlaloc**. Tlaloc is the work/orchestration system; Origami defines representation and transformation semantics and may be used with Tlaloc or independently.

## Project hierarchy

```text
ORIGAMI
├── State semantics
│   └── coherent-state profile
├── Perceptual channels
│   ├── spatial
│   ├── interference / moiré
│   ├── depth / stereo / parallax
│   ├── temporal
│   └── emergent
├── Machine / dynamics / visual projection
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

### Formal Core R0

Defined in [`docs/FORMAL_CORE_R0.md`](docs/FORMAL_CORE_R0.md), [`spec/FORMAL_CORE_R0.json`](spec/FORMAL_CORE_R0.json) and [`spec/OBSERVATION_CONTRACT_R0.json`](spec/OBSERVATION_CONTRACT_R0.json). Non-direct claims require a finite predeclared observation procedure capable of returning an explicit result rather than becoming unfalsifiable.

### Coherent-state profile

Defined in [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md): `determinate`, `superposed`, `coupled`, `observed`, explicit observation authority, coherent interference/cancellation and Fold/Unfold semantics.

The profile is quantum-inspired as a computational analogy, not a claim of physical quantum computation.

### Perceptual Channels R0

Defined in [`docs/PERCEPTUAL_CHANNELS_R0.md`](docs/PERCEPTUAL_CHANNELS_R0.md) and [`spec/PERCEPTUAL_CHANNELS_R0.json`](spec/PERCEPTUAL_CHANNELS_R0.json).

It formalizes spatial, interference, depth, temporal and emergent channels, including **Temporal Latent Images (TLI)** and operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL` and `TEMPORAL_INTEGRATE`.

A valid Origami percept may require a declared relation or trajectory. Failure to reveal it in an insufficient/static observation is not proof of absence; exhaustion of its valid finite Observation Contract without satisfying the expected predicate is `FAIL`.

## Executable OHF research track

The repository also contains the experimental Go workbench inherited from OHF R3.10-LAB. It currently includes deterministic experiment identity, truth-first fixtures, artifact hashing, Native isolation, observation/scoring, failure-to-regression plumbing, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible installation for the OHF laboratory CLIs.

These implementations are **experimental profile/tooling**, not proof that every Origami 6.x semantic operation has a renderer, detector or reference executor.

The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` remains missing; that blocks byte-identical historical R3.10 baseline claims, but it does not block independent Origami 6.x semantic development or OHF-LAB experiments.

## Source of truth and change control

Current project identity and state must be recoverable from the repository:

```text
VERSION
README.md
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/FORMAL_CORE_R0.md
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/FORMAL_CORE_R0.json
spec/OBSERVATION_CONTRACT_R0.json
spec/PERCEPTUAL_CHANNELS_R0.json
state/ORIGAMI_STATE.json
state/components/
changes/
CHANGE_CONTROL_*.md
```

Historical documents and OHF R2/R3.10 artifacts remain evidence, but they do not override the current project boundary.

## Development

Cheap local gates first:

```bash
bash tests/test-formal-contracts.sh
go test ./...
go vet ./...
```

For OHF dimensional/laboratory work, run the affected scripts and regression gates documented under `docs/`, `docs/ohf/` and `scripts/`. External model campaigns remain progressive, cached and conservative.

`install.sh` installs the **OHF laboratory command-line tools** (`ohf-lab`, `ohf-glyphcalc`) under the user prefix. It is not a global installer for every future Origami runtime/profile.

## Brand asset

The ant mark shown above is stored at `assets/origami-ant.webp`. It is a project identity asset and has no semantic or runtime meaning.

## Version

`6.0.0-alpha.3`

The repository now includes Formal Core R0 and finite Observation Contract R0 as experimental contracts. Runtime support remains explicitly separate from semantic specification.
