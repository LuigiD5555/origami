# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.6  
**Status:** experimental

## Root project state

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory layer. Tlaloc remains a separate optional work/orchestration system.

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels / TLI | EXPERIMENTAL_CONTRACT — alpha.2 |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| Dependency Folding R0 | EXPERIMENTAL / reference implementation |
| Hybrid Receiver R0 | EXPERIMENTAL_REFERENCE / model evidence pending |
| Virtual Memory R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.4 |
| multi-carrier GraphSignature federation | IMPLEMENTED / deterministic reference |
| universal ContextPacket interface | IMPLEMENTED / deterministic reference |
| Memory Scale Lab R0 | EXPERIMENTAL_DETERMINISTIC_LAB / fixed-context scaling |
| visual macro/meso/micro memory projection | IMPLEMENTED / Native perception NOT PROMOTED |
| Fixed Carrier R2 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.5 |
| Evidence Reduction R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.6 |
| receiver Master Prompt | REFERENCE_CANDIDATE / unchanged |
| OHF R3.10-LAB | EXPERIMENTAL / executable |

The machine-readable root state is `state/ORIGAMI_STATE.json`.

## Virtual Memory R0

The implemented reference contract is `origami.virtual-memory.r0`.

The active model-facing working window defaults to:

```text
W = 4000 token-equivalent
```

`W` is not storage capacity. The runtime implements:

```text
large/multi-carrier memory
  -> compact GraphSignatures
  -> carrier selection
  -> local metadata index
  -> declared graph traversal
  -> selective fidelity unfold
  -> ContextPacket <= active budget
```

Memory nodes preserve Hybrid R0 compatibility (`Address`, `Value`, `Links`) while optionally adding CID, carrier/cluster identity, semantic type, multi-resolution levels, parent/children, typed edges, importance/cost, provenance and verification references.

Fidelity order is:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

The runtime degrades fidelity rather than exceed budget. Search never silently indexes/opens the complete exact payload.

## Address and federation state

Canonical federated address form:

```text
ohf://<carrier>/<cluster>/<node>
```

Address identifies location; CID identifies content. Multiple addresses may share one CID.

Each rich carrier exposes a compact `GraphSignature`. Federation ranks signatures before mounting/searching local nodes. Federation reports global routing and local metadata preparation separately through:

```text
signatures_scanned
candidate_carriers
selected_carriers
metadata_nodes_loaded
```

This prevents a bounded ContextPacket from hiding an uncontrolled preparation scan.

## Memory Scale Lab R0

`origami.memory-scale-lab.r0` is the deterministic scale campaign for Virtual Memory.

Canonical configuration:

```text
carriers: 1 -> 10 -> 100 -> 1,000
nodes per carrier: 64
active model-facing budget: 4,000 token-equivalent
selected carrier limit: 4
```

It measures:

```text
Carrier Top-1 Accuracy
Carrier Recall@K
Target Hit Accuracy
Evidence Hit Accuracy
Negative UNKNOWN Accuracy
Useful Context Efficiency (UCE)
Knowledge Scale Degradation (KSD)
Routing Selectivity
Exposure Fraction
```

and keeps routing/metadata work separate from model-facing unfolding work.

`cmd/origami-memory-scale` writes:

```text
runs/memory-scale-r0/report.json
runs/memory-scale-r0/trace.jsonl
```

The deterministic lab does **not** establish LLM answer quality or visual readability by itself.

## Evidence Reduction R0

`origami.evidence-reduction.r0` closes the authority gap between probabilistic model/swarm reasoning and deterministic Origami state.

The boundary is:

```text
Origami ContextPacket
  -> Tlaloc/external swarm proposals
  -> Origami evidence address resolution
  -> CID/hash/fidelity/Verified checks
  -> deterministic deduplication
  -> canonical reduction
```

External agents emit structured proposals with a `claim_key`, one of `SUPPORT`, `OPPOSE`, `UNKNOWN`, and zero or more addressable evidence references. Agent confidence is diagnostic only and is never truth authority.

Accepted evidence must exist at the requested address/fidelity, match declared CID/source hash when present, and have `Verified=true`. Repeated agents citing the same evidence do not manufacture additional support.

Canonical states are:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

`CONFLICT` preserves verified disagreement instead of hiding it through majority voting. `UNKNOWN` is preferred when neither side reaches the evidence policy.

Exactness is separate from claim verification. Only accepted verified `exact` fidelity on a `VERIFIED` claim can produce `VERIFIED_EXACT`; all other outcomes use `NOT_CLAIMED`.

`cmd/origami-reduce` exposes this contract to Tlaloc, local models or other agent systems without embedding swarm orchestration inside Origami.

## Fixed Carrier R2

`origami.fixed-carrier.r2` is the frozen visual control plane introduced in alpha.5.

Canonical physical profile:

```text
canvas: 640 x 640
frozen PNG bytes: 8192
hard maximum: 512000 bytes
```

Bootstrap:

```text
T0 plaintext BOOT
T1 Rosetta + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

Corpus growth occurs in the canonical memory plane and does not grow the fixed canvas. OCR is not BOOT and Native visual readability remains unpromoted without held-out evidence.

## Model-facing interface

`ContextPacket` is model neutral. A model receives selected context plus relations, contradictions, provenance, expandable addresses, budget and access metrics. A model does not need Native Origami vision to use this interface.

`cmd/origami-memory` exposes:

```text
QUERY
EXPAND
SIGNATURE
VERIFY
LOOKUP
```

The Hybrid runtime/model packet also preserves:

```text
BOOT
LOOKUP
FOLLOW
TRACE
VERIFY
STOP
```

and adds `QUERY`, `EXPAND`, `SIGNATURE` for new carriers.

The new deterministic answer boundary is intentionally separate:

```text
model/swarm proposes -> cmd/origami-reduce -> canonical claim state
```

## Visual Memory Layout R0

A memory carrier can be one PNG containing:

```text
experimental redundant MACRO/MESO/MICRO navigation band
+
verified exact self-decodable Hybrid/Glyph transport
```

Deterministic `DecodeMemoryPNG` removes only the declared navigation band and recovers the exact logical carrier. That proves transport separation, not VLM readability.

The reference band uses structural navigation cues derived from F01, F21, F22, F24, F25, F29, F32, F34, F39, F40, F41 and F42. `spec/VISUAL_MEMORY_PROFILE_R0.json` records possible memory roles for all F01–F42 families, but they remain evidence-gated generator-tree candidates.

Native visual navigation is **not promoted**. Required tests are specified in `spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json` and include carrier/cluster selection, macro→meso→micro descent, relation following, random address reopening, distractor rejection, resize/degradation and symbol permutation.

## Hybrid receiver state

`origami.hybrid-receiver.r0` remains:

```text
BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
```

Carrier-local physical symbols still have no global semantic meaning. The Master Prompt remains `REFERENCE_CANDIDATE` and was not manually modified by alpha.6.

Hybrid remains the preferred receiver target; Native and Computational remain diagnostic baselines.

## Evidence boundary

Alpha.6 implements deterministic evidence reduction over the existing Virtual Memory authority. It does **not** establish that a model can find the right address, that OCR/VLM perception of arbitrary documents is correct, or that Native visual navigation is ready for promotion.

Important invariants:

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
address != CID
active model interface != total memory
no implicit global exact scan
routing work must be reported
metadata loading must be reported
fidelity fallback before budget violation
visual navigation != exactness authority
swarm proposal != canonical state
agent confidence != evidence authority
duplicate evidence counts once
canonical reduction is order-deterministic
perception != resolution != execution != verification
carrier target <= 500 KB
model specialization is optional
```

Direct PDF extraction, OCR/image perception and Tlaloc swarm search remain separate ingestion/orchestration concerns. Their outputs can be converted into addressable Origami memory and passed through Evidence Reduction R0 without weakening the verification boundary.

## OHF track

OHF remains a nested carrier/protocol laboratory. The historical R3.10 SDK/codebook artifact is still missing, so byte-identical historical baseline claims remain blocked. That does not block independent Origami 6.x contracts or current deterministic laboratory work.

Implemented/experimental OHF assets include Accessible Memory/SuperIndex evidence, Dependency Folding, Hybrid image-backed receiver, Virtual Memory runtime/federation, Memory Scale Lab, visual memory map, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD, DVR/PSS and reversible lab-tool installation.

Evidence Reduction R0 is a generalized Origami core contract rather than an OHF-only mechanism.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
README.md
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/VIRTUAL_MEMORY_R0.md
docs/EVIDENCE_REDUCTION_R0.md
docs/FIXED_CARRIER_R2.md
spec/VIRTUAL_MEMORY_R0.json
spec/EVIDENCE_REDUCTION_R0.json
spec/FIXED_CARRIER_R2.json
spec/VISUAL_MEMORY_PROFILE_R0.json
spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json
spec/MEMORY_SCALE_LAB_R0.json
changes/CHG-ORIGAMI-0008.json
```
