# Origami 6.0.0-alpha.7

Origami is an experimental **visual/computational representation, state-machine language and model-agnostic virtual memory** for complex state, relations, dynamics, selectively accessible knowledge and perceptual emergence.

Origami is independent of **Tlaloc**. Tlaloc owns work/orchestration, document acquisition, swarms and model-facing execution. Origami owns representation semantics, state evolution, Fold/Unfold, memory/evidence authority, carrier semantics and exactness boundaries.

## Current architecture

```text
SOURCE / CANONICAL OBSERVATIONS
           |
           v
EXACT SOURCE PLANE
address + source hash + page/block/layout structure
           |
           v
SEMANTIC SPINE R1
S_t + C_t + declared R
relations / alternatives / higher-order participants
           |
           v
SEMANTIC FOLD
commitment + dependencies + unresolved alternatives
           |
           v
SELECTIVE UNFOLD / VIRTUAL MEMORY
GraphSignature -> local candidates -> declared closure
           |
           v
ContextPacket <= active budget (~4000 token-eq)
           |
           v
Tlaloc / external model or swarm when interpretation is needed
           |
           v
SUPPORT / OPPOSE / UNKNOWN proposals + evidence addresses
           |
           v
ORIGAMI EVIDENCE REDUCTION
           |
           v
VERIFIED / REJECTED / CONFLICT / UNKNOWN
           |
           v
CANONICAL SEMANTIC MEMORY
           |
           v
OPTIONAL FIXED CARRIER / VISUAL CONTROL PROJECTION
```

The carrier is not the corpus. Literal recovery is not semantic Fold. A screenshot/poster is not an Origami semantic representation.

## Project hierarchy

```text
ORIGAMI
├── Semantic Spine R1
│   ├── S_(t+1)=F(S_t,C_t,R)
│   ├── PRESENT / ABSENT / UNKNOWN / INHIBITED / CANCELLED
│   ├── contradictions / signed cancellation / higher-order participants
│   ├── finite Observation Contracts / temporal sequence observations
│   ├── exhaustive finite exploration oracle
│   └── semantic Fold + selective Unfold
├── Exact source / representation boundary
├── Virtual Memory R0
│   ├── address != CID
│   ├── label -> abstract -> summary -> detail -> evidence -> exact
│   ├── local graph + GraphSignature federation
│   ├── bounded ContextBudget / ContextPacket
│   └── model-neutral QUERY / EXPAND / SIGNATURE
├── Evidence Reduction R0
│   ├── SUPPORT / OPPOSE / UNKNOWN candidates
│   ├── address + CID + source-hash + fidelity verification
│   ├── duplicate-evidence collapse
│   ├── VERIFIED / REJECTED / CONFLICT / UNKNOWN
│   └── byte-equal exactness gate
├── Memory Scale Lab R0
│   └── fixed-context 1 -> 10 -> 100 -> 1,000 carrier scaling
├── Fixed Carrier R2
│   └── frozen 640x640 / 8192-byte self-boot control plane
├── Perceptual channels / state semantics
├── Experimental self-boot receiver
│   └── BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
└── OHF research track
    ├── Perception Lab
    ├── R3.10 / R3.10-LAB
    ├── Hybrid Receiver
    ├── visual macro/meso/micro memory map
    ├── DVR / PSS
    ├── Glyph Calculus / SAFE_MICRO_ISA
    └── Context SIMD
```

**OHF is a research/carrier track inside Origami, not the whole language.** Its visual/transport experiments remain preserved and evidence gated.

## Semantic Spine R1

[`docs/SEMANTIC_SPINE_R1.md`](docs/SEMANTIC_SPINE_R1.md) and [`spec/SEMANTIC_SPINE_R1.json`](spec/SEMANTIC_SPINE_R1.json) implement the connected semantic reference path that earlier releases only described in separate contracts.

### State evolution

The executable core uses:

```text
S_(t+1) = F(S_t, C_t, R)
```

`Context` is a real transition input and rules are finite data with stable IDs. Transition evidence records context/state hashes, fired rules, participant sets, contradiction hits and resolutions.

Semantic statuses are distinct:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

`CANCELLED` retains the contributions that cancelled; `INHIBITED` is not absence; `UNKNOWN` is not absence.

### Higher-order / combinatorial behavior

A rule can require any finite participant set. Triple interaction does not follow from pairwise evidence. `internal/core/explore.go` keeps the exhaustive Cartesian product as a reference oracle for small declared domains; future optimizations require equivalence evidence rather than silent pruning.

### Observation and time

Observation is separate from transition. Reference observation kinds are:

```text
state_predicate
temporal_sequence
rule_fired
participant_set
```

A property may exist only over an ordered trajectory. Reordering the same states may destroy it; the runtime does not smuggle the temporal answer into a single frame.

### Semantic Fold / selective Unfold

Fold commits to the canonical semantic graph while preserving dependencies, relations and unresolved alternatives. Selective Unfold opens only the requested address plus its declared dependency/alternative closure and records touched addresses. Loading the full graph and returning a small slice does not count as selective access.

## Source plane is not semantic authority

`origami.document.r0` supplies a neutral source/diagnostic representation. Raw source blocks expose deterministic routing metadata plus evidence/exact payload but do not automatically receive a semantic summary.

`origami-ingest` supports plain text, digitally extractable PDF text through `pdftotext`, and canonical layout JSON from an external adapter. Textless/scanned PDFs fail instead of being rasterized into a fake semantic screenshot.

This path **does not replace Tlaloc Canonical Memory R2**. Tlaloc `6.0.0-alpha.11` remains the existing layout-preserving Canonical Document IR/OCR pipeline with page/region/block CIDs, Merkle exact plane, Tlaloque proposals, CanonicalState, uncertainty/verification queues, ERA H0-H5 and `tlaloc.origami-tools.r2`. Origami defines the representation/evidence semantics that the composition must respect.

## Virtual Memory R0

[`docs/VIRTUAL_MEMORY_R0.md`](docs/VIRTUAL_MEMORY_R0.md) and [`spec/VIRTUAL_MEMORY_R0.json`](spec/VIRTUAL_MEMORY_R0.json) formalize memory that can be much larger than the active model context:

```text
large / multi-carrier memory
  -> GraphSignature routing
  -> local metadata graph
  -> selective fidelity unfold
  -> <= ~4000 token-equivalent ContextPacket
  -> model
```

The ~4000 value is an **active interface budget**, not a storage limit or cumulative task limit.

Fidelity order remains:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Address identifies location; CID identifies content. Canonical federated addresses use:

```text
ohf://<carrier>/<cluster>/<node>
```

Search indexes shallow routing metadata; deep evidence/exact payload is reopened by address. Missing structure returns `UNKNOWN` rather than causing an undeclared global exact scan.

Semantic graphs, deterministic trajectories and canonical evidence-reduction results can now be projected directly into Virtual Memory, so a model does not need to reconstruct deterministic machine state from prose.

## Evidence Reduction R0

[`docs/EVIDENCE_REDUCTION_R0.md`](docs/EVIDENCE_REDUCTION_R0.md) and [`spec/EVIDENCE_REDUCTION_R0.json`](spec/EVIDENCE_REDUCTION_R0.json) define the boundary between probabilistic proposals and canonical state.

```text
ContextPacket
  -> external swarm/model
  -> SUPPORT / OPPOSE / UNKNOWN + evidence addresses
  -> Origami reopens evidence
  -> CID/source-hash/fidelity/Verified checks
  -> duplicate collapse
  -> deterministic reduction
```

Canonical claim states are:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

Verified disagreement remains `CONFLICT`; it is not hidden by majority vote. Agent confidence is diagnostic only.

### Strict exactness

`VERIFIED_EXACT` requires the selected canonical proposition to be **byte-for-byte equal** to accepted verified `exact` evidence. Merely citing an exact block cannot make a paraphrase exact.

```text
exact source + paraphrase -> VERIFIED + NOT_CLAIMED
byte-equal exact source   -> VERIFIED + VERIFIED_EXACT
```

`FALSE_EXACT=0` remains mandatory.

## Fixed Carrier R2

[`docs/FIXED_CARRIER_R2.md`](docs/FIXED_CARRIER_R2.md) and [`spec/FIXED_CARRIER_R2.json`](spec/FIXED_CARRIER_R2.json) define the frozen visual control plane:

```text
canvas      640 x 640
PNG bytes   exactly 8192
hard max    512000
```

Bootstrap:

```text
T0 plaintext BOOT
T1 Rosetta + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

The fixed carrier binds to the external canonical memory `StoreRoot`; corpus growth does not grow the image. OCR is not BOOT. The visual control plane does not grant semantic or exact authority by itself.

## Perception / Native boundary

Perception remains distinct from semantic verification:

```text
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
```

A well-formed VLM/OHF observation can be stored as addressable perception evidence but enters with `Verified=false` until a separate resolution/verification path promotes it.

The existing Perception Lab P0-P6, degradation, attention, swarm, visual-family and Hybrid/Native assets are preserved. Native model promotion remains strict: the tested model receives the universal Master Prompt, one image and the user query—not evaluator ground truth or private registry/manifest data.

Hybrid support and Native visual support are separate claims. Hybrid may eventually be promoted from real BOOT/probe/tool-loop evidence without requiring Native T3 promotion.

## Memory Scale Lab R0

[`docs/MEMORY_SCALE_LAB_R0.md`](docs/MEMORY_SCALE_LAB_R0.md) keeps the canonical deterministic ladder:

```text
1 -> 10 -> 100 -> 1,000 carriers
```

with a fixed 4,000 token-equivalent active budget. It reports routing work, metadata loading and model-facing exposure separately and measures carrier recall, target/evidence accuracy, negative UNKNOWN, Useful Context Efficiency and Knowledge Scale Degradation.

This lab does not by itself establish LLM quality or Native visual readability.

## CLI

Build:

```bash
make build
```

Neutral source ingest:

```bash
./bin/origami-ingest -in paper.pdf -out paper.origami.json
```

Connected semantic/memory interface:

```bash
./bin/origami-system -bundle paper.origami.json -op QUERY -query 'What evidence supports X?' -budget 4000
./bin/origami-system -bundle paper.origami.json -op EXPAND -address 'ohf://paper/page-000001/block-000003' -fidelity evidence -budget 800
./bin/origami-system -bundle paper.origami.json -op REDUCE -proposals proposals.json
./bin/origami-system -bundle paper.origami.json -op PROMOTE -proposals proposals.json -write-bundle paper.semantic.origami.json
```

Existing tools remain available:

```text
origami-memory
origami-memory-scale
origami-reduce
origami-fixed-carrier
origami-hybrid-carrier
origami-hybrid-tool
ohf-lab
ohf-glyphcalc
```

## Hard invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
source plane != semantic plane
literal transport != semantic Fold
screenshot != Origami semantic representation
S_(t+1) = F(S_t,C_t,R)
observation != transition
ABSENT != UNKNOWN != INHIBITED != CANCELLED
Fold preserves unresolved alternatives
Selective Unfold touches declared closure only
address != CID
active model interface != total memory
no implicit global exact scan
visual navigation != exactness authority
swarm proposal != canonical state
agent confidence != evidence authority
duplicate evidence counts once
exact claim requires byte-equal exact evidence
perception != resolution != execution != verification
Fixed Carrier != corpus
nominal visual capacity != SAFE perceptual capacity
carrier target <= 500 KB
model specialization is optional
```

## Evidence still pending

Alpha.7 does not claim:

- built-in OCR for arbitrary scans/images inside Origami;
- implementation of every perceptual-channel operation;
- held-out Native visual navigation;
- cross-model perceptual promotion;
- universal Hybrid support;
- optimized large combinatorial exploration without equivalence to the exhaustive oracle.

These are now promotion/runtime experiments on top of a connected semantic reference rather than gaps hidden by a carrier transport shortcut.

## Development

```bash
go test ./...
go vet ./...
make build
```

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/SEMANTIC_SPINE_R1.md
docs/ARCHITECTURE.md
docs/VIRTUAL_MEMORY_R0.md
docs/EVIDENCE_REDUCTION_R0.md
docs/FIXED_CARRIER_R2.md
spec/SEMANTIC_SPINE_R1.json
spec/VIRTUAL_MEMORY_R0.json
spec/EVIDENCE_REDUCTION_R0.json
spec/FIXED_CARRIER_R2.json
spec/VISUAL_MEMORY_PROFILE_R0.json
spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json
spec/MEMORY_SCALE_LAB_R0.json
changes/
```

## Version

`6.0.0-alpha.7`

Alpha.7 restores and connects Origami's semantic spine while preserving Virtual Memory, Evidence Reduction, Fixed Carrier R2, the OHF/Perception Lab research track and the existing external Tlaloc R2 orchestration/data plane.
