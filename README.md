# Origami 6.0.0-alpha.4

Origami is an experimental **visual/computational representation, state-machine language and model-agnostic virtual memory** for expressing complex states, relations, dynamics, selectively accessible knowledge and perceptual emergence.

Origami is independent of **Tlaloc**. Tlaloc is the optional work/orchestration system; Origami defines representation, memory and transformation semantics and can be used by any compatible model/runtime.

## Project hierarchy

```text
ORIGAMI
├── State semantics
├── Perceptual channels
├── Fold / Unfold / addressability
├── Virtual Memory R0
│   ├── address != CID
│   ├── label -> abstract -> summary -> detail -> evidence -> exact
│   ├── local graph + GraphSignature federation
│   ├── bounded ContextBudget / ContextPacket
│   └── model-agnostic QUERY / EXPAND / SIGNATURE
├── Machine / dynamics / visual projection
├── Experimental self-boot receiver
│   └── BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
└── Research profiles and carriers
    └── OHF — Origami HyperFold
        ├── R3.10 / R3.10-LAB workbench
        ├── Hybrid Receiver / image-backed runtime
        ├── visual macro/meso/micro memory map
        ├── DVR / PSS
        ├── Glyph Calculus / SAFE_MICRO_ISA
        └── Context SIMD
```

**OHF is an Origami research track, not the whole language.** Detailed pre-reconciliation OHF material remains under `docs/ohf/`.

## Virtual Memory R0

[`docs/VIRTUAL_MEMORY_R0.md`](docs/VIRTUAL_MEMORY_R0.md) and [`spec/VIRTUAL_MEMORY_R0.json`](spec/VIRTUAL_MEMORY_R0.json) formalize the current memory direction.

The core target is not to fit all knowledge into a model context window:

```text
large / multi-carrier Origami memory
  -> GraphSignature carrier routing
  -> local metadata index / graph traversal
  -> selective fidelity unfolding
  -> <= ~4000 token-equivalent active ContextPacket
  -> model
```

The ~4000-token value is an **active model-facing interface budget**, not a storage cap and not a cumulative task-token limit.

Each memory node can expose multiple fidelities:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Summaries retain a reopenable address. If a requested fidelity does not fit the current budget, the runtime drops to a lower fidelity before violating the budget.

### Address vs content identity

Canonical federated addresses use:

```text
ohf://<carrier>/<cluster>/<node>
```

Address identifies location. CID identifies content. Multiple addresses may share one CID, allowing deterministic deduplication without losing location/provenance.

### No disguised global scan

Semantic search indexes only labels, classes, clusters, abstracts and summaries. `detail`, `evidence` and `exact` are not silently indexed. Missing relations/addresses return `UNKNOWN` instead of triggering an undeclared full-memory scan.

### Any model can use it

The runtime returns a portable `ContextPacket` containing selected items, relations, contradictions, provenance, expandable references, budget and access metrics. A consumer does not need to understand Origami's visual families.

The future specialized Origami LLM is therefore optional: it may become better at requesting/using memory, but it is not required by the format.

## Multi-Origami federation

Every rich carrier exposes a compact `GraphSignature`. Federation first ranks signatures, then opens only selected carriers and performs local graph-aware lookup.

```text
question
  -> global carrier signatures
  -> candidate Origamis
  -> local graph / clusters
  -> exact addresses
  -> selective unfold
```

This lets total available knowledge scale independently from the active context budget.

## Visual Memory Layout R0

[`docs/VISUAL_MEMORY_LAYOUT_R0.md`](docs/VISUAL_MEMORY_LAYOUT_R0.md) and [`spec/VISUAL_MEMORY_PROFILE_R0.json`](spec/VISUAL_MEMORY_PROFILE_R0.json) define an experimental visual projection.

A new memory carrier remains **one PNG**:

```text
macro / meso / micro navigation band
------------------------------------
exact self-decodable Hybrid/Glyph payload
```

The visual band is redundant navigation, not exactness authority. Computational mode strips the declared band and recovers the exact payload. Native/Hybrid perception can use the map only after Perception Lab evidence supports it.

All F01–F42 families now have documented possible memory/navigation roles. They remain generator-tree candidates: **42 registered families do not mean 42 independently reliable channels**. The current renderer uses only a conservative structural subset and does not promote Native perceptual reliability.

Navigation evaluation is specified in [`spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json`](spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json).

## Hybrid Receiver R0

`origami.hybrid-receiver.r0` still provides the carrier-facing self-boot path:

```text
BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
```

Physical marks remain carrier-local. The experimental `generated/MASTER_PROMPT.md` is still a `REFERENCE_CANDIDATE`; it was not manually changed by Virtual Memory R0.

Hybrid remains the preferred end-to-end target:

```text
model perception
  -> bootstrap / choose region
Origami runtime
  -> address / graph / unfold / compute / verify
model
  -> integrate compact ContextPacket / request next access / answer
```

Native and Computational remain diagnostic modes.

## CLI

Build all local tools:

```bash
make build
```

Universal memory query:

```bash
./bin/origami-memory \
  -carrier runs/.../public/carrier.png \
  -op QUERY \
  -query 'What evidence supports X?' \
  -budget 4000
```

Targeted reopen:

```bash
./bin/origami-memory \
  -carrier runs/.../public/carrier.png \
  -op EXPAND \
  -address 'ohf://carrier/cluster/node' \
  -fidelity evidence \
  -budget 800
```

`origami-hybrid-tool` also supports `QUERY`, `EXPAND` and `SIGNATURE` when those operations are declared in the model packet.

The reversible `install.sh` remains intentionally scoped to OHF laboratory binaries. `origami-memory` is currently built through `make build` rather than silently expanding that installer contract.

## Semantic and perceptual contracts

- [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md): coherent-state kinds and Fold/Unfold semantics.
- [`docs/PERCEPTUAL_CHANNELS_R0.md`](docs/PERCEPTUAL_CHANNELS_R0.md): spatial, interference, depth, temporal and emergent channels, including Temporal Latent Images.

Contract availability is not implementation/promotion. Unsupported perceptual behavior must remain explicit.

## Core invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
address != CID
active model interface != total memory
no implicit global exact scan
visual navigation != exactness authority
nominal visual capacity != SAFE perceptual capacity
perception != resolution != execution != verification
carrier target <= 500 KB
model specialization is optional
```

## Development

```bash
go test ./...
go vet ./...
make build
```

Deterministic tests validate memory selection, fidelity fallback, federation, composite-PNG exact recovery and runtime operations. They **do not** substitute for held-out Native VLM navigation evidence.

## Source of truth

```text
VERSION
README.md
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/VIRTUAL_MEMORY_R0.md
docs/VISUAL_MEMORY_LAYOUT_R0.md
spec/VIRTUAL_MEMORY_R0.json
spec/VISUAL_MEMORY_PROFILE_R0.json
spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json
spec/HYBRID_RECEIVER_R0.json
state/ORIGAMI_STATE.json
changes/
receiver/registry/
```

## Version

`6.0.0-alpha.4`

Alpha.4 implements the deterministic Virtual Memory R0 reference path and redundant visual memory projection. Native visual navigation remains experimental and unpromoted pending held-out cross-model evidence.
