# Origami architecture

Origami separates semantic state, virtual memory, perceptual availability, deterministic execution, receiver bootstrap and carrier/profile projection.

```text
SOURCE / STATE / RELATIONS
          ↓
ORIGAMI SEMANTICS
  ├── coherent-state profile
  ├── Fold / Unfold
  ├── address != CID
  └── perceptual-channel contracts
          ↓
VIRTUAL MEMORY R0
  ├── multi-resolution nodes
  ├── local semantic/dependency/provenance graph
  ├── metadata SuperIndex
  ├── GraphSignature federation
  ├── ContextBudget / WorkingSet
  └── ContextPacket
          ↓
MACHINE / DYNAMICS
          ↓
SELF-BOOT RECEIVER
  BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
          ↓
VISUAL + COMPUTATIONAL PROJECTION
  ├── redundant MACRO/MESO/MICRO map
  └── exact self-decodable transport
          ↓
MODEL / TOOL CONSUMER
```

## 1. Semantic state is not perceptual availability

Coherent-state kinds (`determinate`, `superposed`, `coupled`, `observed`) remain separate from perceptual conditions (`latent`, `motion-bound`, `phase-bound`, `stereo-bound`, `revealed`, `decayed`). A valid semantic state may require a declared trajectory/interference/depth relation to become perceptually available.

## 2. Virtual Memory R0

The core architecture now treats Origami memory as a potentially much larger address space behind a bounded model-facing working set.

```text
available knowledge >> active context
```

The default experimental active interface is ~4000 token-equivalent, but that number does not limit total carrier/federated memory.

### Addressable node

A node can hold identity/structure without forcing all payload into active context:

```text
address          location identity
CID              content identity
carrier/cluster  routing scope
kind/class/label searchable metadata
parent/children  hierarchy
links/edges      graph
importance/cost  selection hints
provenance       source reopening
levels           label/abstract/summary/detail/evidence/exact
```

Address and CID are deliberately separate. Duplicate content can occur at multiple locations/provenance paths and be deduplicated at exposure time.

### Metadata index

The metadata index is intentionally shallow. It indexes labels/classes/clusters plus abstract/summary views. Deeper evidence/exact payload is opened by address. This prevents preprocessing from hiding a whole-memory scan behind an apparently small result.

### Local graph

Local graph traversal follows only declared links/typed edges reachable from relevant candidates. Missing structure fails closed as `UNKNOWN` instead of falling back to a full scan.

### Federation / global graph

Each carrier publishes a compact `GraphSignature`:

```text
root terms
semantic classes
relation types
clusters
node count
metadata cost
```

Multi-carrier federation ranks those signatures first, selects a small carrier set, qualifies carrier-local addresses, then performs local query/unfold. This is the implemented global-graph -> local-graph path.

### Context budget

`ContextBudget` tracks maximum, reserve, used and remaining token-equivalent cost. Selection prioritizes expected relevance/importance/confidence per token and degrades fidelity before budget violation.

### ContextPacket

`ContextPacket` is the model-neutral interface. It carries selected items, relations, contradictions, provenance, expandable references, budget and access metrics. The consumer does not need to decode visual families.

## 3. Working-set cycle

```text
question
  -> query plan
  -> carrier signatures
  -> local metadata candidates
  -> bounded graph traversal
  -> selective unfold
  -> ContextPacket
  -> model
  -> retain useful working-set items / reopen address / answer
```

Cumulative task information may exceed the active window. What must remain bounded is each model-facing working set, with touched/unfolded/exposed work measured separately.

## 4. Self-boot receiver

`origami.hybrid-receiver.r0` remains the bootstrap contract:

```text
BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
```

Carrier-local physical marks do not gain universal semantics. BOOT finds the receiver entry, ROSETTA maps local physical states, PROGRAM exposes bounded micro-rules, INDEX/MEMORY expose the address space, and VERIFICATION preserves exactness authority.

## 5. Visual memory projection

`RenderMemory` creates one carrier PNG with two complementary regions:

```text
redundant navigation band
  MACRO -> MESO -> MICRO
exact self-decodable transport
```

The navigation band represents hierarchy/grouping/graph structure so a visual receiver can decide where to look before spending context. The exact transport remains present in the same PNG so deterministic recovery and visual navigation can be tested separately.

The current reference projection uses a conservative structural subset of F01–F42. `spec/VISUAL_MEMORY_PROFILE_R0.json` assigns potential memory roles to all 42 registered families, but no family becomes a reliable channel without Perception Lab promotion.

## 6. Three receiver modes

### Computational

```text
carrier.png -> deterministic Origami runtime -> ContextPacket -> any model
```

No visual competence is required from the model.

### Native

```text
Master Prompt + carrier.png + question -> VLM
```

No external decoder/runtime. This remains the strictest perceptual diagnostic and requires held-out evidence.

### Hybrid

```text
model perception -> high-level route / bootstrap
Origami runtime  -> graph/address/unfold/compute/verify
model            -> integrate ContextPacket / decide next access
```

Hybrid remains the preferred end-to-end direction.

## 7. Visual-family discipline

Families F01–F42 are generator-tree dimensions, not an instruction to use a full Cartesian product. The production/evidence rule remains:

```text
candidate family/composition
 -> single-family evidence
 -> pair/composition evidence
 -> OOD/permutation
 -> resize/degradation
 -> higher-order navigation
 -> cross-model evidence
 -> promotion
```

Nominal visual capacity must never be reported as SAFE perceptual capacity without this evidence.

## 8. Exactness and verification

Semantic usefulness, structural correctness and byte/source exactness are separate claims. A summary can answer a semantic question while remaining insufficient for exact reconstruction. `FALSE_EXACT = 0` is mandatory; verified evidence/source reopening remains the authority for exact claims.

## 9. OHF research track

OHF remains the principal executable carrier/protocol laboratory inside Origami. Glyph Calculus, SAFE_MICRO_ISA, DVR/PSS, Context SIMD, Dependency Folding and the Hybrid carrier are implementation/research substrates; they do not redefine all Origami semantics.

## 10. Tlaloc boundary

Tlaloc may search/distill better receiver behavior and specialized model interaction. Origami owns the memory representation, addresses, Fold/Unfold, carrier-local semantics, runtime operations and verification boundaries. A future Origami-specialized LLM is an optional consumer, not a prerequisite of Virtual Memory R0.
