# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.9  
**Status:** experimental

## Root project state

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory layer. Tlaloc remains a separate optional work/orchestration and experimentation system.

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels / TLI | EXPERIMENTAL_CONTRACT — alpha.2 |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| Virtual Memory R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.4 |
| Memory Scale Lab R0 | EXPERIMENTAL_DETERMINISTIC_LAB |
| Fixed Carrier R2 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.5 |
| Evidence Reduction R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.6 |
| Semantic Spine R1 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.7 |
| Perception Promotion R1 | EXPERIMENTAL_REFERENCE_IMPLEMENTED / REAL MODEL EVIDENCE PENDING — alpha.8 |
| Canonical Visual Grammar R0 | EXPERIMENTAL_REFERENCE_PROFILE — alpha.9 |
| Writer R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.9 |
| Master Prompt | READ/WRITE REFERENCE_CANDIDATE |
| generic Writer plan -> final visual renderer | PARTIAL / next compiler integration step |
| OHF R3.10-LAB / Perception Lab | EXPERIMENTAL / executable / promotion evidence gated |

The machine-readable root state is `state/ORIGAMI_STATE.json`.

## Canonical visual rule

Origami now has **one canonical functional aesthetic per visual profile version**.

```text
profile version
    |
    v
canonical aesthetic
    |
    +-- geometry
    +-- contrast / fill
    +-- position
    +-- topology
    +-- enclosure
    +-- scale
    +-- repetition
    +-- density
    +-- limited text
    |
    v
many carriers with different semantics
```

ROSETTA remains mandatory. It declares profile/version, active dimensions and concrete visual-token/semantic-role mappings.

ROSETTA does not authorize each PDF or model to invent a private aesthetic.

Old carriers remain interpretable because BOOT + profile version + ROSETTA travel with the carrier.

## Current aesthetic R0

The active reference profile is:

```text
origami.canonical-aesthetic.r0
```

It is high-contrast and primarily monochrome.

Current canonical visual roles include:

```text
ringed square       root / verification anchor
black-center square page / major container
white square        concept
white diamond       source / evidence origin
small black square  graph / relation cluster
line / arrow        relation / dependency
box / enclosure     scope / address space
black/white cells   visual probe bits
macro/meso/micro    representation level
repetition          redundancy / agreement
```

Color, mathematical/numeric structures and temporal/phase channels remain experimental candidate dimensions, not primary authority.

A pattern based on primes, modular arithmetic or another mathematical structure is therefore a valid **candidate representation strategy** for Tlaloc to test. It does not become Origami merely by looking structured.

## Writer R0

Origami now supports a formal WRITE direction in addition to READ.

```text
READ
carrier -> BOOT -> ROSETTA -> semantic state

WRITE
source
 -> document/source ingestion
 -> Semantic IR
 -> visual intent plan
 -> canonical visual profile
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> compiler
 -> carrier
 -> semantic roundtrip verification
```

`internal/writer` implements the deterministic construction-plan reference.

The CLI is:

```bash
./bin/origami-write-plan -in request.json -out plan.json
```

The plan contains:

```text
profile_id
source identity
semantic commitment
visual intents
ROSETTA
required carrier sections
verification plan
canonical semantic graph
```

Visual intent classes include:

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

The model declares function; the canonical profile chooses visual encoding.

## PDF/image -> Origami

A PDF or image is not converted by rasterizing the source page into the carrier.

The intended composed path is:

```text
PDF / IMAGE
   |
   v
Tlaloc Canonical Document IR / OCR / source plane
   |
   v
Semantic IR
   |
   v
Origami Writer R0
   |
   v
canonical visual construction plan
   |
   v
Origami compiler / renderer
   |
   v
carrier.png
```

Tlaloc's existing document ingestion pipeline remains external; Origami does not duplicate it.

The current alpha.9 writer produces the canonical construction plan. Full generic renderer integration for every Writer plan is still a next implementation step; the existing Fixed Carrier renderer remains deterministic and unchanged.

## Master Prompt R1

`generated/MASTER_PROMPT.md` now teaches both READ and WRITE modes.

READ mode preserves the existing BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION process.

WRITE mode instructs a compatible model to:

1. ingest or use declared source/document structure;
2. produce Semantic IR;
3. declare visual intents;
4. bind them to the canonical profile;
5. generate ROSETTA;
6. construct PROGRAM/INDEX/MEMORY/VERIFICATION;
7. invoke deterministic compilation when available;
8. roundtrip verify the resulting carrier.

If no compiler exists, the model may emit a construction specification but must not claim that an arbitrary painted image is a valid Origami.

## Roundtrip gate

Writer R0 requires:

```text
DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR
```

Identity, addresses, hashes and proof fields require exact equality.

Semantic fields require declared semantic equivalence.

Failure state:

```text
INVALID_CARRIER
```

## Semantic Spine R1

The semantic core remains:

```text
S_(t+1) = F(S_t, C_t, R)
```

with first-class Context, declarative Rules, transition provenance and distinct:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

Contradictions, signed cancellation, higher-order participant sets, observation contracts and temporal observations remain implemented.

## Observation and temporal state

Observation is separate from transition.

Implemented reference observation kinds are:

```text
state_predicate
temporal_sequence
rule_fired
participant_set
```

The complete experimental perceptual operation family remains incomplete. Moire, phase, stereo and other advanced channels remain evidence/runtime work rather than assumed capacity.

## Semantic Fold / Unfold

Semantic Fold preserves semantic node identity, dependencies, unresolved alternatives, relations and evidence references.

Selective Unfold opens only the requested address plus declared dependency/alternative closure and records touched addresses.

A hidden whole-graph load cannot count as selective access merely because the returned packet is small.

## Virtual Memory R0

The active model-facing working window defaults to approximately:

```text
4000 token-equivalent
```

This is an active interface budget, not total storage capacity.

The memory path remains:

```text
large / multi-carrier memory
 -> GraphSignature routing
 -> local metadata graph
 -> selective fidelity unfold
 -> ContextPacket
```

Fidelity order:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Address is location; CID is content identity. Deep exact payload is reopened by address rather than silently globally scanned.

## Evidence Reduction R0

External models/Tlaloques emit:

```text
SUPPORT
OPPOSE
UNKNOWN
```

with evidence addresses.

Origami verifies address/CID/source-hash/fidelity/Verified state, deduplicates repeated evidence and reduces deterministically to:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

Verified disagreement remains `CONFLICT`.

`VERIFIED_EXACT` requires byte-equality with accepted verified exact evidence. Exact citation alone does not make a paraphrase exact.

## Fixed Carrier R2

The current frozen visual control plane remains:

```text
640 x 640
8192 PNG bytes exactly
hard maximum 512000 bytes
```

Bootstrap:

```text
T0 plaintext BOOT
T1 ROSETTA + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

The carrier binds external canonical memory rather than growing with the corpus.

Fixed Carrier R2 is now treated as the current concrete reference instance of the canonical aesthetic R0, not as a universal limitation on future promoted visual profiles.

## Perception Promotion R1

Origami deterministically checks a model's claimed perception against the actual carrier:

```text
T0 BOOT
probe top
probe bottom
probe agreement
Tool Protocol
Address ABI
optional T3
```

MOCK may validate mechanics but cannot promote.

Hybrid and Native T3 remain separate trial claims.

Real cross-model aggregation belongs to Tlaloc; stack promotion belongs to Tonal.

## How the aesthetic evolves

The visible aesthetic may change over time, but only as a promoted profile revision:

```text
current Origami profile
    |
    v
Tlaloc representation/prompt search
    |
    +-- shape candidates
    +-- layout candidates
    +-- color candidates
    +-- redundancy candidates
    +-- numeric / prime / modular candidates
    +-- temporal candidates
    +-- Master Prompt candidates
    |
    v
measured candidate evidence
    |
    v
Origami contract validation
    |
    v
Tonal promotion
    |
    v
new canonical profile version
```

A Tlaloc experiment does not directly mutate canonical Origami semantics or pixels.

## Hard invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
ONE_CANONICAL_AESTHETIC_PER_PROFILE_VERSION
PROFILE_BOUND_SYMBOLS_HAVE_VERSIONED_MEANING
ROSETTA_ALWAYS_PRESENT
ROSETTA != PER_DOCUMENT_AESTHETIC_DRIFT
AESTHETIC_IS_FUNCTIONAL_ENCODING
SOURCE_SCREENSHOT != ORIGAMI
SEMANTIC_IR_PRECEDES_VISUAL_COMPILATION
CONSTRUCTION_SPEC_PRECEDES_COMPILATION
ROUNDTRIP_REQUIRED
TLALOC_PROPOSES / ORIGAMI_VALIDATES / TONAL_PROMOTES
source plane != semantic plane
literal transport != semantic Fold
S_(t+1)=F(S_t,C_t,R)
observation != transition
ABSENT != UNKNOWN != INHIBITED != CANCELLED
Fold preserves unresolved alternatives
Selective Unfold touches declared closure only
address != CID
active model interface != total memory
no implicit global exact scan
perception != resolution != execution != verification
Fixed Carrier != corpus
carrier target <= 500 KB for current profile
model specialization is optional
```

## Evidence boundary

Alpha.9 establishes the canonical READ/WRITE contract and deterministic construction-plan front end.

It does **not** establish that:

- arbitrary scans/images are ingested natively by Origami without external adapters;
- every experimental perceptual channel works;
- any color/numeric/prime-derived/temporal candidate is superior to current R0;
- Native visual interpretation is cross-model promoted;
- a generic Writer R0 plan can already be rendered by a fully profile-generic compiler.

Those remain explicit experiments/implementation steps.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
README.md
PROJECT_BOUNDARY.md
generated/MASTER_PROMPT.md
docs/CURRENT_STATE.md
docs/CANONICAL_VISUAL_GRAMMAR_R0.md
docs/WRITER_R0.md
docs/SEMANTIC_SPINE_R1.md
docs/PERCEPTION_PROMOTION_R1.md
spec/CANONICAL_VISUAL_GRAMMAR_R0.json
spec/WRITER_R0.json
spec/HYBRID_RECEIVER_R0.json
changes/CHG-ORIGAMI-0011.json
```
