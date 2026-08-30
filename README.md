# Origami 6.0.0-alpha.9

Origami is an experimental **visual/computational representation, state-machine language and model-agnostic virtual memory** for complex state, relations, dynamics, selectively accessible knowledge and perceptual emergence.

Origami is independent of **Tlaloc**. Tlaloc owns work/orchestration, document acquisition, swarms, model-facing experiments and candidate search. Origami owns representation semantics, state evolution, Fold/Unfold, memory/evidence authority, canonical visual grammar, carrier semantics and exactness boundaries.

## What changed in alpha.9

Origami now has an explicit **single canonical visual aesthetic per profile version** and a WRITE path.

The rule is:

```text
many documents
     |
     v
one Semantic IR language
     |
     v
one canonical Origami visual profile/version
     |
     v
many carriers that share the same functional aesthetic
```

`ROSETTA` remains mandatory, but it self-describes the active canonical profile, dimensions and concrete bindings. It does **not** allow every document/model to invent a private aesthetic.

Tlaloc may experiment with radically different candidate prompts, colors, shapes, layouts, mathematical structures or perceptual dimensions. A successful candidate changes Origami's aesthetic only through a **versioned profile promotion**.

## Canonical architecture

```text
PDF / IMAGE / TEXT / CONVERSATION / STATE
                    |
                    v
        SOURCE / DOCUMENT INGESTION
     (Tlaloc or another declared adapter)
                    |
                    v
              SEMANTIC IR
 concepts / relations / hierarchy / state /
 uncertainty / evidence / addresses / time
                    |
          +---------+----------+
          |                    |
          v                    v
 SEMANTIC SPINE R1       WRITER R0
 S_t + C_t + R           visual intents
 Fold / Unfold                 |
 evidence                      v
          |        CANONICAL VISUAL GRAMMAR R0
          |                    |
          |                    v
          |       ROSETTA + PROGRAM + INDEX +
          |          MEMORY + VERIFICATION
          |                    |
          |                    v
          |          deterministic compiler
          |                    |
          |                    v
          |               carrier.png
          |                    |
          +----------+---------+
                     |
                     v
              VIRTUAL MEMORY
                     |
                     v
         bounded ContextPacket (~4000)
                     |
                     v
             model / Tlaloc
                     |
                     v
        evidence-backed proposals
                     |
                     v
          EVIDENCE REDUCTION R0
                     |
                     v
 VERIFIED / REJECTED / CONFLICT / UNKNOWN
```

A source screenshot is not Origami. A text poster is not Origami. Literal recovery is not semantic Fold. The visual carrier is a compiled representation/control surface over canonical semantics and addressable memory.

## Canonical Visual Grammar R0

[`docs/CANONICAL_VISUAL_GRAMMAR_R0.md`](docs/CANONICAL_VISUAL_GRAMMAR_R0.md) and [`spec/CANONICAL_VISUAL_GRAMMAR_R0.json`](spec/CANONICAL_VISUAL_GRAMMAR_R0.json) define the current aesthetic profile:

```text
origami.canonical-aesthetic.r0
```

The aesthetic is functional encoding, not decoration.

Current canonical dimensions are:

```text
geometry
fill / contrast
position
topology
enclosure
scale
repetition
density
limited declarative text
```

The current Fixed Carrier family is intentionally high-contrast and mostly monochrome because those channels currently have the strongest deterministic support.

Experimental candidate dimensions include:

```text
color
numeric / mathematical structure
temporal / phase structure
```

These may become canonical later if Tlaloc experiments demonstrate measurable gains without semantic/exactness regressions.

For example, prime-derived or modular visual organization is a legitimate candidate strategy. It is not automatically useful merely because the resulting image looks complex or interesting.

## ROSETTA

ROSETTA is still a core invariant.

It declares:

```text
profile id
profile version
active dimensions
visual primitive -> semantic role mappings
profile-approved extensions/parameters
```

This lets a receiver bootstrap a carrier without hidden repository knowledge and keeps older profile versions readable after the aesthetic evolves.

The rule is:

```text
ONE CANONICAL PROFILE + SELF-DESCRIBING ROSETTA
```

not:

```text
EVERY CARRIER INVENTS A NEW VISUAL LANGUAGE
```

## Writer R0 — READ and WRITE

[`docs/WRITER_R0.md`](docs/WRITER_R0.md) and [`spec/WRITER_R0.json`](spec/WRITER_R0.json) add the missing write direction.

```text
READ
carrier -> BOOT -> ROSETTA -> semantic state

WRITE
source -> Semantic IR -> visual intent -> canonical grammar
       -> ROSETTA -> construction plan -> compiler -> verify
```

The universal Master Prompt now teaches both modes.

When asked to convert a PDF into Origami, a compatible model should not paint PDF text into an image. It should preserve/document source structure, construct Semantic IR, declare visual intents and use the canonical compiler path.

The reference planner is:

```bash
./bin/origami-write-plan -in request.json -out plan.json
```

The planner deterministically maps semantic intents such as:

```text
IDENTITY
HIERARCHY
RELATION
DEPENDENCY
SCOPE
IMPORTANCE
STATE
UNCERTAINTY
EVIDENCE
VERIFICATION
ADDRESSABILITY
TEMPORAL_ORDER
REDUNDANCY
ROUTING
```

onto the approved canonical dimensions.

If no compiler is available, a model may emit a construction specification, but it must not claim that an arbitrary image is a valid carrier.

## Roundtrip requirement

A compiled carrier must satisfy:

```text
DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR
```

with exact equality for identity, addresses, hashes and verification/proof fields.

Failure means `INVALID_CARRIER`, not “visually close enough.”

## Semantic Spine R1

[`docs/SEMANTIC_SPINE_R1.md`](docs/SEMANTIC_SPINE_R1.md) and [`spec/SEMANTIC_SPINE_R1.json`](spec/SEMANTIC_SPINE_R1.json) implement the connected semantic reference path.

The executable core uses:

```text
S_(t+1) = F(S_t, C_t, R)
```

with first-class Context and declarative Rule IDs.

Semantic statuses remain distinct:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

The core also preserves contradiction provenance, signed cancellation, higher-order participant rules, finite observation contracts and temporal observations over ordered trajectories.

The exhaustive finite state enumerator remains the reference oracle for small declared state spaces; optimized search cannot silently replace it without equivalence evidence.

## Observation is not transition

Reference observations include:

```text
state_predicate
temporal_sequence
rule_fired
participant_set
```

A property may exist only over an ordered trajectory. Reordering the same states can destroy it. Observation does not mutate the committed trajectory.

The complete alpha.2 perceptual operation family is still not universally implemented; moire/phase/stereo/temporal candidates remain evidence-gated research capabilities.

## Semantic Fold / selective Unfold

Semantic Fold commits to the canonical graph while preserving dependencies, relations and unresolved alternatives.

Selective Unfold opens only the requested address plus its declared dependency/alternative closure and records touched addresses.

Therefore loading the entire graph and merely returning a small answer does not count as selective access.

## Virtual Memory R0

[`docs/VIRTUAL_MEMORY_R0.md`](docs/VIRTUAL_MEMORY_R0.md) formalizes memory much larger than the active model context:

```text
large / multi-carrier memory
  -> GraphSignature routing
  -> local metadata graph
  -> selective fidelity unfold
  -> ContextPacket <= active budget
  -> model
```

The default model-facing working window is approximately:

```text
4000 token-equivalent
```

This is not total memory capacity and not a cumulative task limit.

Fidelity order remains:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Canonical addresses use:

```text
ohf://<carrier>/<cluster>/<node>
```

Address is location; CID is content identity.

Deep evidence/exact payload is reopened by address rather than silently globally indexed.

## Evidence Reduction R0

External models/Tlaloques produce candidates:

```text
SUPPORT
OPPOSE
UNKNOWN
```

Origami reopens cited evidence and verifies address, CID, source hash, fidelity and `Verified` state.

Duplicate agents citing the same evidence do not manufacture consensus.

Canonical states are:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

Verified disagreement remains `CONFLICT` rather than being hidden by majority vote.

Agent confidence is diagnostic only.

### Strict exactness

`VERIFIED_EXACT` requires the selected canonical proposition to be byte-for-byte equal to accepted verified `exact` evidence.

```text
exact source + paraphrase -> VERIFIED + NOT_CLAIMED
byte-equal exact source   -> VERIFIED + VERIFIED_EXACT
```

`FALSE_EXACT=0` remains mandatory.

## Fixed Carrier R2

Fixed Carrier remains the frozen visual/control profile currently used by experiments:

```text
canvas      640 x 640
PNG bytes   exactly 8192
hard max    512000
```

Bootstrap remains:

```text
T0 plaintext BOOT
T1 Rosetta + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

The carrier binds to external canonical memory through StoreRoot/source identity rather than containing the corpus as a screenshot.

The current image therefore already represents the beginning of the canonical aesthetic R0: high contrast, geometry, topology, macro/meso/micro hierarchy, address regions, redundant probes and explicit verification anchors.

## Perception Promotion R1

[`docs/PERCEPTION_PROMOTION_R1.md`](docs/PERCEPTION_PROMOTION_R1.md) supplies the deterministic per-trial evaluator.

It independently checks what a model claims to have seen against the real carrier:

```text
T0 BOOT
upper visual probe
lower visual probe
probe agreement
Tool Protocol
Address ABI
optional T3 machine record
```

`MOCK` can exercise the complete mechanical path but can never become promotion evidence.

Hybrid and Native T3 remain separate claims:

```text
HYBRID trial
  T0 + probes + protocol/address negotiation

NATIVE T3 trial
  Hybrid mechanical entry + T3 machine record
```

Tlaloc owns cross-model campaigns, transport variants and real tool loops. Tonal owns aggregate stack promotion.

## Tlaloc-driven visual evolution

Origami intentionally does **not** contain an optimizer that changes its own visual language on every run.

Instead:

```text
current canonical Origami profile
          |
          v
Tlaloc representation experiments
 prompt / shape / layout / color /
 redundancy / mathematical structure /
 temporal channel candidates
          |
          v
candidate profile + evidence
          |
          v
Origami contract validation
          |
          v
Tonal promotion gate
          |
          v
new canonical Origami profile version
```

This allows Origami's visible aesthetic to evolve as better representations are discovered while preserving one coherent language at each profile version.

## Memory Scale Lab R0

The deterministic scale ladder remains:

```text
1 -> 10 -> 100 -> 1,000 carriers
```

under a fixed 4,000 token-equivalent active budget.

It reports routing/metadata work separately from model-facing exposure and measures routing recall, target/evidence accuracy, negative UNKNOWN, Useful Context Efficiency and Knowledge Scale Degradation.

The deterministic lab does not by itself establish VLM quality or Native visual readability.

## Model-facing CLIs

Build everything:

```bash
make build
```

Core tools include:

```text
origami-ingest
origami-system
origami-memory
origami-memory-scale
origami-reduce
origami-write-plan
origami-fixed-carrier
origami-perception-eval
origami-hybrid-carrier
origami-hybrid-tool
ohf-lab
ohf-glyphcalc
```

Example source/semantic flow:

```bash
./bin/origami-ingest -in paper.pdf -out paper.origami.json
./bin/origami-system -bundle paper.origami.json -op QUERY -query 'What supports X?' -budget 4000
./bin/origami-write-plan -in construction-request.json -out construction-plan.json
```

The current Writer R0 produces the deterministic canonical construction plan. Full integration from that generic plan into a next-generation visual renderer remains an explicit following step; the existing Fixed Carrier renderer remains intact and deterministic.

## Hard invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
source plane != semantic plane
literal transport != semantic Fold
screenshot != Origami semantic representation
one canonical aesthetic per profile version
ROSETTA always present
ROSETTA != permission for per-document style drift
profile-bound symbols have versioned meaning
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
construction spec precedes visual compilation
roundtrip required
Tlaloc proposes / Origami validates / Tonal promotes
model specialization is optional
carrier target <= 500 KB
```

## Evidence still pending

Alpha.9 does not claim:

- built-in arbitrary scan/image OCR inside Origami;
- implementation of every experimental perceptual channel;
- held-out Native visual navigation;
- cross-model Hybrid/Native promotion;
- that color/prime-derived/numeric/temporal representation candidates improve current R0;
- a promoted successor visual profile;
- fully generic rendering of every Writer R0 construction plan into pixels.

Those are now explicit experimentation/promotion tasks rather than hidden architectural gaps.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
generated/MASTER_PROMPT.md
docs/CURRENT_STATE.md
docs/CANONICAL_VISUAL_GRAMMAR_R0.md
docs/WRITER_R0.md
docs/SEMANTIC_SPINE_R1.md
docs/PERCEPTION_PROMOTION_R1.md
docs/VIRTUAL_MEMORY_R0.md
docs/EVIDENCE_REDUCTION_R0.md
docs/FIXED_CARRIER_R2.md
spec/CANONICAL_VISUAL_GRAMMAR_R0.json
spec/WRITER_R0.json
spec/SEMANTIC_SPINE_R1.json
spec/PERCEPTION_PROMOTION_R1.json
spec/VIRTUAL_MEMORY_R0.json
spec/EVIDENCE_REDUCTION_R0.json
spec/FIXED_CARRIER_R2.json
changes/
```

## Version

`6.0.0-alpha.9`
