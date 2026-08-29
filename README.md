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

### Coherent-state profile

Defined in [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md): `determinate`, `superposed`, `coupled`, `observed`, explicit observation authority, coherent interference/cancellation and Fold/Unfold semantics.

The profile is quantum-inspired as a computational analogy, not a claim of physical quantum computation.

### Perceptual Channels R0

Defined in [`docs/PERCEPTUAL_CHANNELS_R0.md`](docs/PERCEPTUAL_CHANNELS_R0.md) and [`spec/PERCEPTUAL_CHANNELS_R0.json`](spec/PERCEPTUAL_CHANNELS_R0.json).

It formalizes spatial, interference, depth, temporal and emergent channels, including **Temporal Latent Images (TLI)** and operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL` and `TEMPORAL_INTEGRATE`.

A valid Origami percept may exist only through a declared relation or trajectory; failure to reveal it in one static frame is not proof of absence.

## Executable OHF research track

The repository also contains the experimental Go workbench inherited from OHF R3.10-LAB. It currently includes deterministic experiment identity, truth-first fixtures, artifact hashing, Native isolation, observation/scoring, failure-to-regression plumbing, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible installation for the OHF laboratory CLIs.

These implementations are **experimental profile/tooling**, not proof that every Origami 6.x semantic operation has a renderer, detector or reference executor.

The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` remains missing; that blocks byte-identical historical R3.10 baseline claims, but it does not block independent Origami 6.x semantic development or OHF-LAB experiments.

Accordingly, `ohf-lab doctor` reports readiness only for the historical OHF R3.10 baseline. `OHF R3.10 READY = NO` with `MISSING_R310_SDK` is an expected scoped block, not a failure of Origami core. Use `make check` to verify the current Go implementation independently.

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

For OHF dimensional/laboratory work, run the affected scripts and regression gates documented under `docs/`, `docs/ohf/` and `scripts/`. External model campaigns remain progressive, cached and conservative.

`install.sh` installs the **OHF laboratory command-line tools** (`ohf-lab`, `ohf-glyphcalc`) under the user prefix. It is not a global installer for every future Origami runtime/profile.

## Version

`6.0.0-alpha.3`

This release reconciles repository hierarchy and machine-readable state after integrating the OHF R3.10-LAB workbench. It does not change the semantic laws introduced in alpha.1 or alpha.2.
