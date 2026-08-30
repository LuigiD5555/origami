# Current Origami state

**Status date:** 2026-08-29  
**Project version:** 6.0.0-alpha.4  
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
| receiver Master Prompt | REFERENCE_CANDIDATE / unchanged |
| OHF R3.10-LAB | EXPERIMENTAL / executable |

The machine-readable root state is `state/ORIGAMI_STATE.json`.

## Virtual Memory R0

The implemented reference contract is `origami.virtual-memory.r0`.

The active model-facing working window defaults to:

```text
W = 4000 token-equivalent
```

`W` is not storage capacity. The runtime now implements:

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

Each rich carrier exposes a compact `GraphSignature`. Federation ranks signatures before mounting/searching local nodes. Federation now reports global routing and local metadata preparation separately through:

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

The fixture deliberately gives all carriers a shared domain/evidence vocabulary plus query-specific topics. This makes routing non-trivial while keeping exact truth deterministic.

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

The JSONL trace is the stable case/evidence contract for the next phase: small local LLMs, larger baseline LLMs, Hybrid model-driven access and eventually Native visual navigation can all consume the same cases.

The deterministic lab does **not** establish LLM answer quality or visual readability by itself.

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

## Visual Memory Layout R0

A memory carrier can now be one PNG containing:

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

Carrier-local physical symbols still have no global semantic meaning. The Master Prompt remains `REFERENCE_CANDIDATE` and was not manually modified by alpha.4.

Hybrid remains the preferred receiver target; Native and Computational remain diagnostic baselines.

## Evidence boundary

Alpha.4 implements the deterministic Virtual Memory reference. Memory Scale Lab R0 tests whether its routing/unfolding behavior remains bounded as available memory grows. It does **not** promote model-facing answer quality or the visual Native profile without held-out evidence.

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
42 registered families != 42 safe independent channels
perception != resolution != execution != verification
carrier target <= 500 KB
model specialization is optional
```

## OHF track

OHF remains a nested carrier/protocol laboratory. The historical R3.10 SDK/codebook artifact is still missing, so byte-identical historical baseline claims remain blocked. That does not block independent Origami 6.x contracts or current deterministic laboratory work.

Implemented/experimental OHF assets now include Accessible Memory/SuperIndex evidence, Dependency Folding, Hybrid image-backed receiver, Virtual Memory runtime/federation, Memory Scale Lab, visual memory map, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD, DVR/PSS and reversible lab-tool installation.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
README.md
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/VIRTUAL_MEMORY_R0.md
docs/VISUAL_MEMORY_LAYOUT_R0.md
docs/MEMORY_SCALE_LAB_R0.md
spec/VIRTUAL_MEMORY_R0.json
spec/VISUAL_MEMORY_PROFILE_R0.json
spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json
spec/MEMORY_SCALE_LAB_R0.json
spec/HYBRID_RECEIVER_R0.json
changes/CHG-ORIGAMI-0006.json
```
