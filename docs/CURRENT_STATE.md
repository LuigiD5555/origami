# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.7  
**Status:** experimental

## Root project state

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory layer. Tlaloc remains a separate optional work/orchestration system.

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels / TLI | EXPERIMENTAL_CONTRACT — alpha.2 |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| Dependency Folding R0 | EXPERIMENTAL / historical reference implementation |
| Hybrid Receiver R0 | EXPERIMENTAL_REFERENCE / model evidence pending |
| Virtual Memory R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.4 |
| multi-carrier GraphSignature federation | IMPLEMENTED / deterministic reference |
| universal ContextPacket interface | IMPLEMENTED / deterministic reference |
| Memory Scale Lab R0 | EXPERIMENTAL_DETERMINISTIC_LAB / fixed-context scaling |
| visual macro/meso/micro memory projection | IMPLEMENTED / Native perception NOT PROMOTED |
| Fixed Carrier R2 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.5 |
| Evidence Reduction R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.6 |
| Semantic Spine R1 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.7 |
| receiver Master Prompt | REFERENCE_CANDIDATE / unchanged |
| OHF R3.10-LAB / Perception Lab | EXPERIMENTAL / executable / promotion evidence gated |

The machine-readable root state is `state/ORIGAMI_STATE.json`.

## Semantic Spine R1

`origami.semantic-spine.r1` reconnects the semantic intent that had remained fragmented across early contracts, reference code, Virtual Memory and OHF experiments.

The canonical path is now:

```text
exact source / canonical observations
  -> addressable source memory
  -> S_t + C_t + declared R
  -> deterministic semantic transition
  -> semantic Fold
  -> selective Unfold / Virtual Memory
  -> bounded ContextPacket
  -> external model/Tlaloc swarm when interpretation is needed
  -> evidence-backed proposals
  -> deterministic Evidence Reduction
  -> canonical semantic memory
  -> optional visual/Fixed Carrier projection
```

The executable semantic core implements:

```text
S_(t+1) = F(S_t, C_t, R)
```

with declarative rules carrying stable IDs and transition provenance. It preserves distinct semantic statuses:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

Context is a first-class transition input. Contradictions are declarative and terminate before ordinary mutation. Signed contributions retain cancellation provenance. Higher-order participant sets can be required directly; pairwise evidence does not promote a three-way interaction.

The exhaustive finite state enumerator remains the small reference oracle. Optimized exploration is not promoted unless it can be checked against that oracle.

## Observation and temporal state

Observation is separate from transition. `internal/core/observation.go` evaluates finite contracts over an already committed trajectory.

Implemented reference observation kinds are:

```text
state_predicate
temporal_sequence
rule_fired
participant_set
```

A temporal property may exist only over an ordered trajectory. Reordering the same frames can invalidate it. The observer does not mutate the trajectory or smuggle a temporal answer into one frame.

The complete alpha.2 perceptual operation family is still not implemented. Moire/phase/stereo and other perceptual-channel promotion remain explicit future evidence/runtime work.

## Semantic Fold / Unfold

Semantic Fold is no longer equivalent to compressing literal bytes or collapsing a dependency chain.

`internal/core/fold.go` produces a commitment plus an address index that preserves:

```text
semantic node identity
dependencies
unresolved alternatives
relations
evidence references
```

Selective Unfold reads only the requested address plus its declared dependency/alternative closure. The reference store records touched addresses so an implementation cannot load the full graph and merely return a small answer while claiming selective access.

Fold does not resolve an unresolved alternative without an explicit later decision/observation path.

## Exact source plane

`origami.document.r0` is a neutral exact-source representation/diagnostic path. It is deliberately separate from semantic authority.

Raw blocks may expose:

```text
abstract  deterministic routing metadata
evidence  extracted source block
exact     extracted source block
```

They do **not** receive a semantic `summary` merely because text was extracted.

`origami-ingest` can ingest plain text or a digitally extractable PDF through `pdftotext`. A textless/scanned PDF fails rather than falling back to a raster screenshot and calling it semantic Origami. Canonical layout JSON is accepted for external OCR/VLM adapters.

This does not replace Tlaloc's existing R2 document pipeline. Tlaloc 6.0.0-alpha.11 remains the external Canonical Document IR/OCR, Tlaloque orchestration, Merkle exact-plane, CanonicalState, ERA H0-H5 and `tlaloc.origami-tools.r2` provider. Origami alpha.7 defines the representation/evidence semantics that such a composition must respect.

## Virtual Memory R0

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

Fidelity order remains:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

The runtime degrades fidelity rather than exceed budget. Search never silently indexes/opens the complete exact payload.

Canonical federated address form remains:

```text
ohf://<carrier>/<cluster>/<node>
```

Address identifies location; CID identifies content. Multiple addresses may share one CID.

Federation reports global routing and local preparation separately through:

```text
signatures_scanned
candidate_carriers
selected_carriers
metadata_nodes_loaded
```

so a small ContextPacket cannot hide an uncontrolled preparation scan.

## Core/trajectory projection into memory

Semantic graphs and deterministic trajectories can now be projected directly into Virtual Memory without requiring an LLM to reconstruct machine state from prose.

Projection preserves semantic state, relations, alternatives, Fold commitment, logical time, context IDs, fired rules, transition evidence and trajectory hashes.

Evidence Reduction results can also be promoted as `canonical-claim` memory nodes with reopenable support/opposition links.

## Evidence Reduction R0

External agents emit structured proposals using:

```text
SUPPORT
OPPOSE
UNKNOWN
```

with addressable evidence references. Agent confidence is diagnostic only.

Origami accepts evidence only when the requested address/fidelity exists, declared CID/source hash matches when supplied and the content level is verified. Duplicate agents citing the same evidence do not manufacture additional support.

Canonical states are:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

`CONFLICT` preserves verified disagreement. `UNKNOWN` is preferred when neither side reaches the evidence policy.

Exactness is stricter in alpha.7: opening or citing an `exact` block is not enough. `VERIFIED_EXACT` requires the selected canonical proposition to be byte-for-byte equal to accepted verified exact-fidelity evidence. A supported paraphrase may be `VERIFIED` while exactness remains `NOT_CLAIMED`.

## Fixed Carrier R2

Fixed Carrier remains the frozen visual control plane:

```text
canvas: 640 x 640
frozen PNG bytes: 8192
hard maximum: 512000 bytes
```

Bootstrap remains:

```text
T0 plaintext BOOT
T1 Rosetta + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

Alpha.7 binds the fixed control plane to the external canonical memory root rather than treating the pixels as the corpus. Corpus growth occurs outside the frozen canvas. OCR is not BOOT and Native visual readability remains unpromoted without held-out evidence.

Existing Hybrid/glyph exact transport remains useful as a transport experiment but literal recovery alone is not semantic Fold authority.

## Perception boundary

Existing OHF/native perception observations may enter the shared addressable plane, but they remain:

```text
PERCEPTION_ONLY
Verified = false
```

until a separate resolution/verification path promotes them. A well-formed VLM observation therefore cannot directly satisfy semantic Evidence Reduction.

The governing invariant remains:

```text
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
```

Perception Lab P0-P6 and visual-family experiments are preserved. Native model promotion remains evidence gated; the strict Native setup continues to mean a universal Master Prompt plus one carrier image plus the user query, without leaking private evaluator manifests/registries to the model.

## Memory Scale Lab R0

The deterministic scale ladder remains:

```text
carriers: 1 -> 10 -> 100 -> 1,000
nodes per carrier: 64
active model-facing budget: 4,000 token-equivalent
selected carrier limit: 4
```

It measures carrier routing, target/evidence accuracy, negative UNKNOWN, Useful Context Efficiency, Knowledge Scale Degradation, routing selectivity and exposure fraction while separating routing/metadata work from model-facing unfolding.

The deterministic lab does **not** establish LLM answer quality or Native visual readability by itself.

## Model-facing interfaces

Current core CLIs include:

```text
origami-memory       QUERY / EXPAND / SIGNATURE / VERIFY / LOOKUP
origami-reduce       deterministic proposal/evidence reduction
origami-ingest       neutral source-plane ingest
origami-system       QUERY / EXPAND / LOOKUP / SIGNATURE / VERIFY / REDUCE / PROMOTE
origami-fixed-carrier frozen R2 control plane
```

The Hybrid receiver packet also preserves BOOT/LOOKUP/FOLLOW/TRACE/VERIFY/STOP semantics where declared.

## Evidence boundary

Alpha.7 establishes an executable connected semantic reference. It does **not** claim:

- built-in OCR for arbitrary scans/images inside Origami;
- complete runtime support for every alpha.2 perceptual operation;
- held-out Native visual navigation;
- cross-model perceptual promotion;
- universal Hybrid support;
- large optimized combinatorial exploration without equivalence to the exhaustive oracle.

Those remain experimental promotion work, not missing semantic authority.

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
routing work must be reported
metadata loading must be reported
visual navigation != exactness authority
swarm proposal != canonical state
agent confidence != evidence authority
duplicate evidence counts once
exact claim requires byte-equal exact evidence
perception != resolution != execution != verification
Fixed Carrier != corpus
carrier target <= 500 KB
model specialization is optional
```

## OHF track

OHF remains a nested carrier/protocol and perceptual laboratory. Glyph Calculus, SAFE_MICRO_ISA, Context SIMD, DVR/PSS, visual memory projection, Perception Lab and Hybrid receiver assets are retained. They no longer stand in for the semantic core merely because they can transport or display information.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
README.md
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
changes/CHG-ORIGAMI-0009.json
```
