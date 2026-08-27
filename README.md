# Origami HyperFold (OHF)

**Origami HyperFold is an experimental protocol and Go workbench for representing structured information as a compact, addressable, selectively unfoldable and verifiable carrier.**

OHF is not intended to be a giant QR code or a PNG wrapper around conventional compression. The research direction is to move as much work as possible away from open-ended LLM reasoning and into deterministic structure: indexes, bounded generative instructions, local visual states, small dependency windows and cryptographic verification.

> **Research status:** R3.10 is the current protocol reference. This repository contains the **R3.10-LAB experimental workbench**, not a stable production release. The exact historical R3.10 SDK artifact is not currently present in the repository, so the laboratory reports that baseline as blocked rather than reconstructing it from documentation.

## Core invariants

- `FALSE_EXACT = 0` — exactness is never inferred from plausible visual interpretation.
- `UNKNOWN` is preferable to invented certainty.
- Go is the canonical deterministic runtime.
- Native input is isolated to `Master Prompt + image + query`.
- Perception, resolution, execution and verification are separate stages.
- The target visual carrier remains `<= 500 KB` unless the protocol is explicitly changed.
- Experimental visual dimensions are not promoted until Perception Lab evidence supports them.
- External VLM tests are progressive, cached and rate-limited; the project avoids unnecessary requests.

## Current implementation

The repository now contains a reproducible experimental foundation rather than documentation only.

### R3.10-LAB foundation

Implemented:

- canonical `ExperimentSpec` and content-addressed experiment identity;
- deterministic namespaced seeds;
- deterministic synthetic source generation and truth-first ground truth;
- artifact manifests and SHA-256 verification;
- deterministic LAB-only geometry renderer;
- PNG save/reopen verification with separate file and pixel hashes;
- Native blind bundles with an exact three-input boundary;
- structured `ObservationEnvelope` values: `KNOWN`, `AMBIGUOUS`, `UNKNOWN`, `INVALID`;
- external scoring against ground truth;
- failure registration and reproducible regressions;
- conservative external-execution policy.

### Glyph Calculus

`ohf-glyphcalc` explores a lower-cost execution model in which the VLM performs only small perceptual operations over glyphs and Go performs deterministic reconstruction.

The experimental bounded instruction set currently includes operations such as:

```text
HALT
LIT
RUN
REF
SEQ
```

The compiler can use structured instructions when they reduce work and fall back to raw transport when the data does not offer useful structure.

### SAFE_MICRO_ISA and Context SIMD

Instead of assuming which visual operations are easy for an LLM, the laboratory profiles them empirically.

Candidate operations include:

```text
OBSERVE4   observe four physical glyph states
TUPLE4     emit four base-4 state digits
DECODE8    map the four states to one byte
```

The profiler measures correctness first and can then compare retries, output tokens, latency and cost when those metrics are available. `Context SIMD` batches independent lanes in one prompt without claiming hardware-level SIMD inside the model.

The widening policy is adaptive: test one lane first, select the safest operation, then try `4 -> 8 -> 16` lanes and stop at the first unsafe boundary.

## Architecture

```text
SOURCE
  ↓
Canonical Source Model
  ↓
Structure / Graph / Grammar / Transform discovery
  ↓
Representation Tournament
  ↓
Generative IR + Residual
  ↓
SuperIndex + Dependency Manifests
  ↓
Attention Router + Window Sliders
  ↓
Visual Compiler / Perceptual ABI
  ↓
OHF Carrier
```

Receiver:

```text
QUERY
  ↓
BOOT / ROSETTA / SUPERINDEX
  ↓
minimum dependency closure
  ↓
PERCEPTION  -> evidence only
  ↓
RESOLUTION  -> deterministic candidate collapse
  ↓
EXECUTION   -> bounded generative operations
  ↓
RESIDUAL
  ↓
VERIFICATION
  ↓
ANSWER / EXACT DATA
```

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for the current design and [`docs/SCIENTIFIC_FOUNDATIONS.md`](docs/SCIENTIFIC_FOUNDATIONS.md) for the longer foundational research document.

## Quick start

Prerequisite: **Go 1.23+**.

```bash
git clone https://github.com/LuigiD5555/origami.git
cd origami
./install.sh --check-only
./install.sh
```

The installer:

- does not use `sudo`;
- does not download packages;
- makes zero external VLM/API requests;
- installs `ohf-lab` and `ohf-glyphcalc` under `~/.local/bin` by default;
- runs a local deterministic smoke test.

Useful alternatives:

```bash
./install.sh --prepare-local
./install.sh --full-check
./install.sh --prefix "$HOME/bin"
```

Or work directly with Go:

```bash
go test ./...
go vet ./...
go run ./cmd/ohf-lab doctor
go run ./cmd/ohf-glyphcalc --help
```

## Reproducible local experiments

Foundation checks:

```bash
./scripts/reproduce_f1.sh
./scripts/reproduce_f1_visual.sh
./scripts/reproduce_f1_evaluator.sh
```

Glyph Calculus local matrix:

```bash
./scripts/glyphcalc_local_matrix.sh
```

Prepare the SAFE_MICRO_ISA campaign locally:

```bash
./scripts/microisa_prepare_local.sh
```

This generates local blind bundles but makes **zero external requests**.

Then ask the state machine which single probe is allowed next:

```bash
go run ./cmd/ohf-glyphcalc microbench-next \
  --root runs/microbench \
  --min-exact 0.98
```

Do not manually skip ahead in the campaign. The adaptive procedure exists to reduce cost, avoid redundant requests and stop early when a candidate is unsafe.

## Repository layout

```text
origami/
├── cmd/
│   ├── ohf-lab/              # reproducible OHF laboratory CLI
│   └── ohf-glyphcalc/        # glyph calculus / micro-ISA CLI
├── internal/lab/             # deterministic laboratory implementation
├── experiments/              # declarative experiment specifications
├── profiles/                 # execution policies
├── state/                    # current machine-readable protocol/component state
├── changes/                  # immutable-ish change records
├── regressions/              # reproducible failures (local evidence ignored by default)
├── scripts/                  # repeatable workflows
├── testdata/                 # minimal fixtures
├── docs/                     # architecture, state, experiment protocol and research docs
├── prompts/                  # prompt history; R2 is legacy
├── install.sh
├── Makefile
└── go.mod
```

## State and change control

The chat history is not the source of truth for the implementation.

Machine-readable state lives in:

```text
state/ORIGAMI_STATE.json
state/components/*.json
changes/*.json
```

Every material change should record:

1. component changed;
2. before/after state;
3. evidence/run identifiers;
4. affected downstream components;
5. regressions created or resolved;
6. promotion state.

See [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md).

## R2 legacy baseline

The repository originally documented an R2 binary-grid ZIP carrier. That prompt remains under `prompts/OHF_R2_MASTER_PROMPT.txt` for historical reproducibility, but it is **not the current OHF research direction** and must not be confused with R3.10-LAB Native experiments.

See [`docs/LEGACY_R2.md`](docs/LEGACY_R2.md).

## What is not yet proven

OHF does **not** currently claim that:

- arbitrary multi-gigabyte data can be losslessly represented in a <=500 KB self-contained carrier;
- Macro-Gestalt visual compounds are production-safe;
- one SAFE_MICRO_ISA profile transfers across all VLMs;
- Native-only reconstruction is solved for large carriers;
- a visual description that looks plausible is exact.

Those are experiment targets or open questions, not completed claims.

## Documentation

Start with [`docs/README.md`](docs/README.md).

Key documents:

- [`docs/CURRENT_STATE.md`](docs/CURRENT_STATE.md) — what exists now and what is blocked.
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — current sender/receiver architecture.
- [`docs/EXPERIMENT_PROTOCOL.md`](docs/EXPERIMENT_PROTOCOL.md) — strict reproducibility and external-test procedure.
- [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md) — contribution and change-control workflow.
- [`docs/ROADMAP.md`](docs/ROADMAP.md) — current research sequence.
- [`docs/SCIENTIFIC_FOUNDATIONS.md`](docs/SCIENTIFIC_FOUNDATIONS.md) — long-form scientific motivation and foundations.

## Core law

> **UNKNOWN is preferable to false exactness.**

A model may perceive probabilistically. Exact recovery is committed only by deterministic verification.
