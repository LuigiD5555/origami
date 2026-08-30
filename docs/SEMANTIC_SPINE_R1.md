# Origami Semantic Spine R1

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`

Semantic Spine R1 restores the executable path that the alpha.3 intent audit described but the repository did not yet implement as one connected system.

The governing rule is:

> A visual carrier, literal payload, OCR result, model interpretation or swarm vote is not Origami semantic authority by itself.

## Canonical path

```text
source bytes / extracted text / layout observations
        |
        v
EXACT SOURCE PLANE
address + source hash + page/block/layout structure
        |
        +---------------------- exact/evidence reopening
        |
        v
SEMANTIC CORE R1
S_t + C_t + declared R
relations / alternatives / higher-order participants
        |
        v
SEMANTIC FOLD
commitment + address index + dependencies + alternatives
        |
        v
SELECTIVE UNFOLD / VIRTUAL MEMORY
GraphSignature -> local candidates -> declared closure
        |
        v
ContextPacket <= active budget
        |
        v
model / Tlaloc swarm when interpretation is needed
        |
        v
SUPPORT / OPPOSE / UNKNOWN proposals + addresses
        |
        v
EVIDENCE REDUCTION
address / CID / source hash / fidelity / verification
        |
        v
VERIFIED / REJECTED / CONFLICT / UNKNOWN
        |
        v
CANONICAL SEMANTIC MEMORY
        |
        +---------------------- future bounded queries
        |
        v
optional Fixed Carrier / visual navigation control plane
```

The image is therefore an interface/projection of the machine and memory topology. It is not allowed to replace the machine by drawing all source text into pixels.

## 1. First-class state, context and rules

`internal/core` implements the formal equation:

```text
S_(t+1) = F(S_t, C_t, R)
```

`Context` is an explicit transition input and every `Rule` is finite data with a stable ID. Transition evidence records the context hash, state hashes, fired rule IDs, participant sets and resolutions.

The runtime does not ask a language model to apply rules that the Go core claims to implement.

## 2. Semantic statuses

The executable core preserves:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

They are not aliases. `CANCELLED` can arise from signed contributions whose sum is zero while retaining the contributing rules. `INHIBITED` is also distinct from absence and uncertainty.

## 3. Contradiction and higher-order interaction

Contradictions are declared as rules with IDs and reasons. A matching contradiction terminates before ordinary mutation and remains in transition evidence.

Rules can require any finite participant set. A triple interaction can therefore exist without pretending that pairwise evidence proves it.

`internal/core/explore.go` keeps exhaustive Cartesian enumeration as the small reference oracle. Optimized exploration may be added later only with equivalence evidence.

## 4. Observation does not mutate state

`internal/core/observation.go` evaluates finite Observation Contracts over an already committed trajectory.

Implemented reference observations:

```text
state_predicate
temporal_sequence
rule_fired
participant_set
```

Observation results are:

```text
PASS
FAIL
INVALID_CONTRACT
UNSUPPORTED
```

An observation budget is finite. Failing to reveal a requested temporal/emergent property inside that budget is a real failure, not permission to invent it.

## 5. Temporal latent properties

A temporal property may be true only over an ordered trajectory.

The regression suite demonstrates a property that exists for:

```text
A -> B -> C
```

and disappears for:

```text
A -> C -> B
```

No state stores the answer string `ABC`. The property is produced by the declared observation over the trajectory.

## 6. Semantic Fold and selective Unfold

`internal/core/fold.go` folds a semantic graph into:

```text
commitment
address index
CID per semantic node
dependency references
unresolved alternatives
```

The canonical graph remains in a separate store. Selective Unfold reads only the requested address plus its declared dependency/alternative closure.

This matters because a fake implementation could load the complete graph and merely return two nodes. The new store tracks touched addresses, and the adversarial test requires an unrelated region to remain untouched.

Fold never resolves an unresolved alternative merely to make the representation smaller.

## 7. Exact source plane is not semantic Fold

`internal/document` introduces `origami.document.r0`.

A raw text/PDF/OCR block can expose:

```text
abstract  deterministic routing metadata only
evidence  exact extracted block
exact     exact extracted block
```

It deliberately receives **no semantic summary automatically**.

For a digital PDF, `origami-ingest -mode pdf` uses `pdftotext`. If text cannot be extracted, it fails instead of rasterizing the PDF and calling the screenshot a semantic Origami representation.

Scans/images can enter through a canonical layout JSON adapter that preserves page/block/type/bounding-box structure. An OCR/VLM adapter is still evidence-pending; its output does not automatically become semantic truth.

## 8. Bounded model access

Virtual Memory R0 remains the active memory interface:

```text
GraphSignature
 -> metadata routing
 -> declared graph traversal
 -> selective fidelity unfold
 -> ContextPacket
```

The ~4000 token-equivalent default is the simultaneous model-facing window, not storage capacity.

Raw exact blocks are reopened by address. Deep exact payload is not silently globally indexed.

## 9. Swarm interpretation and deterministic reduction

Tlaloc or another external system may run a swarm when semantic interpretation is genuinely model-dependent.

Agents produce candidates, not facts:

```text
SUPPORT
OPPOSE
UNKNOWN
+ evidence addresses
```

Origami reopens those addresses and checks identity and verification. Duplicate evidence counts once. Verified opposition is preserved as `CONFLICT` rather than hidden by majority vote.

A verified exact source block does **not** make a paraphrase exact. `VERIFIED_EXACT` now requires the selected canonical proposition to be byte-for-byte equal to accepted `exact` evidence.

## 10. Canonical semantic memory

After evidence reduction, `internal/integration` can promote the canonical result back into Virtual Memory as a `canonical-claim` node.

That node contains:

- a semantic summary;
- evidence audit;
- links to supporting/opposing source addresses;
- reduction digest and canonical claim ID;
- exact fidelity only when strict exactness is earned.

Subsequent questions can query this semantic memory directly instead of rerunning the entire source interpretation.

## 11. Core-to-memory projection

Semantic graphs and deterministic core trajectories can be projected directly into Virtual Memory without an LLM.

This preserves:

```text
semantic node status/data
relations/dependencies/alternatives
Fold commitment
logical-time states
transition evidence
context IDs
rule provenance
trajectory hash
```

The model is not required to reconstruct deterministic machine state from prose.

## 12. Perception boundary

Existing OHF/native perception observations can enter the shared memory plane through `PerceptionNode`.

They are deliberately stored with:

```text
Verified = false
VerificationRef = PERCEPTION_ONLY_NOT_SEMANTIC_AUTHORITY
```

A well-formed perception record therefore cannot satisfy Evidence Reduction by itself. A separate resolution/verification step is required.

This enforces:

```text
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
```

## 13. Fixed Carrier boundary

Fixed Carrier R2 remains a fixed 640x640, 8192-byte visual **control plane**.

Semantic Spine R1 binds it to the canonical external memory using:

```text
StoreRoot
SourceSHA256
compact GraphSignature
object/page/block counts
```

The corpus does not grow the image. The fixed carrier points to the memory universe and teaches navigation/bootstrap; it is not a screenshot or corpus dump.

The older Hybrid glyph-grid exact transport remains useful for deterministic transport experiments, but exact transport alone has no semantic authority under R1 accounting.

## 14. Literal/residual anti-cheat accounting

`internal/core/accounting.go` separates two claims:

```text
SEMANTIC
EXACT
```

A representation made only of literal/residual bytes may support explicit exact recovery accounting, but it cannot earn a semantic-computation claim.

Semantic authority requires represented semantic material plus machine-derived and verified semantic claims. Removing the exact transport must not destroy the semantic state being claimed.

## 15. CLI path

Digital PDF or text:

```bash
./bin/origami-ingest -in paper.pdf -out paper.origami.json
```

Bounded query:

```bash
./bin/origami-system \
  -bundle paper.origami.json \
  -op QUERY \
  -query 'entropy dynamic patches' \
  -budget 4000
```

Targeted reopening:

```bash
./bin/origami-system \
  -bundle paper.origami.json \
  -op EXPAND \
  -address 'ohf://paper/page-000001/block-000003' \
  -fidelity evidence \
  -budget 800
```

External swarm output can then be reduced:

```bash
./bin/origami-system \
  -bundle paper.origami.json \
  -op REDUCE \
  -proposals proposals.json
```

and promoted back into the bundle:

```bash
./bin/origami-system \
  -bundle paper.origami.json \
  -op PROMOTE \
  -proposals proposals.json \
  -write-bundle paper.semantic.origami.json
```

The promoted bundle can be queried again with the same bounded interface.

## 16. What remains evidence-pending

R1 does not claim completion of everything the research direction envisions. These remain explicit:

- built-in scanned-PDF/image OCR adapter;
- general executable implementations for every perceptual-channel operation (moire, phase, stereo, etc.);
- held-out Native visual navigation evidence;
- cross-model perceptual promotion;
- large optimized combinatorial exploration proven equivalent to the exhaustive oracle;
- proof that a specialized small Origami model provides superior navigation/interpretation.

Those gaps are now downstream of a connected semantic spine instead of being silently replaced by an image transport trick.
