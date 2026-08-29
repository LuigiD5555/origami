# Origami 6.0.0-alpha.3

Origami is an experimental **visual/computational representation and state-machine language** for expressing complex states, relations, dynamics, selectively accessible memory and perceptual emergence.

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
├── Fold / Unfold / addressability
├── Machine / dynamics / visual projection
├── Experimental self-boot receiver
│   └── BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
└── Research profiles and carriers
    └── OHF — Origami HyperFold
        ├── R3.10 / R3.10-LAB workbench
        ├── Dependency Folding / Accessible Memory Lab
        ├── Hybrid Receiver experiments
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

### Hybrid Receiver R0 — experimental, not promoted

The feature branch introduces [`spec/HYBRID_RECEIVER_R0.json`](spec/HYBRID_RECEIVER_R0.json) and a deterministic reference implementation under `internal/receiver/`.

A carrier can describe how to start and interpret itself through:

```text
BOOT
  ↓
ROSETTA
  ↓
PROGRAM
  ↓
INDEX
  ↓
MEMORY
  ↓
VERIFICATION
```

Physical marks are carrier-local; a triangle, circle or color does not acquire universal semantic meaning merely because another carrier used it that way. Symbol-permutation tests require equivalent carriers with different physical mappings to produce the same semantics through their own Rosetta tables.

The preferred end-to-end mode is **Hybrid**:

```text
model perception
  -> BOOT / carrier-local Rosetta / high-level navigation
Origami deterministic runtime
  -> micro-agents / addressability / Fold-Unfold / computation / verification
model
  -> compact retained semantic state / next bounded access / answer
```

Native and Computational are retained as diagnostic baselines.

The experimental external bootstrap is stored at `generated/MASTER_PROMPT.md`, with provenance/status under `receiver/registry/`. It is currently `REFERENCE_CANDIDATE`, not PROMOTED.

## Accessible Memory direction

The current research question is not how much knowledge can be squeezed directly into a model's context window. The working target is to make a much larger Origami memory usable through a bounded interface:

```text
large Origami memory
  -> SuperIndex / selective unfolding / internal computation / Dependency Folding
  -> <= ~4000 token-equivalent active model-facing context
  -> compact retained semantic state
  -> release obsolete evidence
  -> next bounded access
```

The ~4000 token-equivalent value is an **active interface budget**, not an Origami storage cap and not a cumulative task-token limit. Internal work, touched memory and model-facing exposure are measured separately.

## Tlaloc relationship

Tlaloc may use its Behavior Compiler, Tlaloque and richer swarm behavior to discover how a model can bootstrap/navigate Origami more reliably. Successful behavior can be distilled into candidate receiver prompts, BOOT/Rosetta strategies and simple deterministic micro-agent rules.

Origami remains authoritative for its semantics. A Tlaloc candidate must pass Origami validation before Origami stores it as a promoted receiver artifact.

## Executable OHF research track

The repository also contains the experimental Go workbench inherited from OHF R3.10-LAB. It currently includes deterministic experiment identity, truth-first fixtures, artifact hashing, Native isolation, Hybrid public-bundle isolation, observation/scoring, failure-to-regression plumbing, Dependency Folding, Accessible Memory/SuperIndex work, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible installation for the OHF laboratory CLIs.

Glyph Calculus and SAFE_MICRO_ISA remain low-level experimental mechanisms. Their physical transport encodings do not define universal carrier semantics; under the Hybrid Receiver they can act as deterministic execution substrates selected after BOOT/ROSETTA interpretation.

These implementations are **experimental profile/tooling**, not proof that every Origami 6.x semantic operation has a renderer, detector or reference executor.

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
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/PERCEPTUAL_CHANNELS_R0.json
spec/HYBRID_RECEIVER_R0.json
state/ORIGAMI_STATE.json
receiver/registry/
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

For receiver work, the immediate deterministic gates cover local symbol permutation, fail-closed UNKNOWN behavior, receiver registry identity and Hybrid public-bundle isolation. Model-facing Native/Hybrid promotion still requires external evidence; deterministic tests do not substitute for VLM readability.

For OHF dimensional/laboratory work, run the affected scripts and regression gates documented under `docs/`, `docs/ohf/` and `scripts/`. External model campaigns remain progressive, cached and conservative.

`install.sh` installs the **OHF laboratory command-line tools** (`ohf-lab`, `ohf-glyphcalc`) under the user prefix. It is not a global installer for every future Origami runtime/profile.

## Version

`6.0.0-alpha.3`

The promoted release version remains alpha.3. Hybrid Receiver, Dependency Folding and Accessible Memory work in the current feature PR is experimental and must not be read as a released alpha.4 until explicit promotion/versioning occurs.
