# Origami project boundary

```text
TLALOC                                      ORIGAMI
work/orchestration system                   representation/state-machine language
Canonical Document IR / OCR         ----->  source/representation contracts
swarm/Tlaloque search               ----->  receiver / claim candidate inputs
behavior compilation + distillation ----->  prompt/boot/micro-agent candidates
representation/prompt search         ----->  candidate visual-profile revisions
ERA H0-H5 / tool loop               <---->  address/memory/evidence contracts
model-facing evaluation             <-----  semantic/visual/reference contracts
promotion recommendation            ----->  candidate evidence
                                             |
                                             +-- semantic validation
                                             +-- canonical visual-profile validation
                                             +-- evidence address resolution
                                             +-- deterministic claim semantics
                                             +-- carrier execution
                                             +-- promoted artifact storage
```

Origami may be used independently of Tlaloc. Tlaloc may operate without Origami or consume Origami as one representation provider.

The stable rule is:

> Tlaloc discovers, ingests, coordinates, searches and benchmarks candidate behavior/representation strategies; Origami defines and executes the canonical representation semantics, including state evolution, Fold/Unfold, memory, evidence, visual grammar, ROSETTA, carrier compilation and exactness.

## Canonical visual ownership

Origami owns **one canonical functional visual aesthetic per promoted profile version**.

Tlaloc does not choose a private aesthetic for each document. Instead it may experimentally mutate and benchmark candidate changes to:

```text
Master Prompt
shape / primitive roles
layout
channel assignments
redundancy
color usage
numeric / mathematical structures
phase / temporal structures
other perceptual dimensions
```

A Tlaloc candidate is only a recommendation with evidence.

The intended lifecycle is:

```text
Origami canonical profile N
        |
        v
Tlaloc search / experiments
        |
        v
candidate profile N+1 + evidence
        |
        v
Origami semantic/visual validation
        |
        v
Tonal composition/promotion gate
        |
        v
Origami canonical profile N+1
```

A successful Tlaloc experiment does not directly mutate the canonical renderer, grammar or Master Prompt in production.

## ROSETTA ownership

ROSETTA remains Origami-owned and mandatory.

With the canonical visual grammar, ROSETTA means:

```text
profile ID/version
active dimensions
canonical visual primitive -> semantic role mappings
carrier-approved parameters/extensions
```

It is a self-description/compatibility mechanism, not permission for arbitrary per-document visual dialects.

Profile-bound symbols may have stable versioned meanings. Old carriers remain readable because profile/version + ROSETTA travel with them.

## Preserved Tlaloc R2 implementation

Tlaloc's R2 orchestration/data-plane remains external and is not replaced by Origami Writer R0:

- layout-preserving Canonical Document IR for digital PDFs with OCR fallback for text-sparse/raster pages;
- stable document/page/region/block addressing and CIDs;
- figures/geometry/canonical text/original source preservation;
- proposal-only Tlaloque candidate generation;
- deterministic CanonicalState, conflicts, uncertainty/verification queues and state ledger;
- Merkle-bound exact plane with `FALSE_EXACT=0`;
- External Recursive Attention H0-H5 under a bounded active context;
- `tlaloc.origami-tools.r2` BOOT/QUERY/EXPAND/VERIFY provider;
- OpenAI-compatible multimodal function loop and plain-text tool bridge;
- Perception Promotion Campaign machinery with real-model, transport, tool-loop and routing gates.

Those remain Tlaloc-owned orchestration/runtime assets.

Origami Writer R0 consumes source/semantic structures and defines how they become a canonical visual construction plan; it does not fork Tlaloc's document ingestion pipeline.

Where Tlaloc carries a deterministic reducer for its R2 runtime, that implementation is an orchestration-side implementation/mirror. Promoted outputs must remain compatible with Origami-owned evidence/exactness semantics rather than creating an independent definition of truth.

## Origami owns

- first-class semantic state, context and declared rule semantics (`S_(t+1)=F(S_t,C_t,R)`);
- PRESENT / ABSENT / UNKNOWN / INHIBITED / CANCELLED distinctions;
- contradiction/cancellation/higher-order interaction semantics;
- finite Observation Contracts and temporal/higher-order observations;
- semantic Fold/Unfold and unresolved-alternative preservation;
- addressability/Virtual Memory contracts and bounded ContextPacket semantics;
- separation of exact source plane from semantic authority;
- evidence address/CID/source-hash/fidelity authority;
- canonical `VERIFIED / REJECTED / CONFLICT / UNKNOWN` evidence semantics;
- exactness rule: an exact reference alone is insufficient; `VERIFIED_EXACT` requires byte-equal accepted exact evidence;
- self-boot receiver contract (`BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION`);
- universal Master Prompt READ/WRITE semantics;
- canonical visual profile and versioned aesthetic;
- visual intent -> canonical dimension/primitive mapping;
- Writer construction-plan semantics and semantic roundtrip requirements;
- Fixed Carrier visual/control-plane semantics and StoreRoot binding;
- deterministic micro-agent execution semantics where generalized into Origami;
- promoted receiver/visual artifacts and provenance;
- evidence/exactness boundaries.

## Tlaloc owns

- document acquisition/canonicalization workflow and OCR execution;
- swarm/Tlaloque exploration;
- behavior-compilation lifecycle;
- candidate Master Prompt search/mutation;
- candidate BOOT/ROSETTA strategy experiments constrained by Origami's canonical profile contract;
- candidate visual primitive/channel/layout/color/numeric/temporal search;
- claim proposal generation and model-side ambiguity exploration;
- uncertainty-driven verification scheduling;
- External Recursive Attention orchestration;
- model/tool-loop transports;
- cross-model perception campaigns and degradation/transport experiments;
- distillation of complex swarm traces into simple bounded micro-agent candidate rules;
- candidate representation tournaments and fitness measurement;
- promotion recommendations and regression coordination.

## Tonal owns

- immutable composition of exact Tlaloc + Origami revisions;
- stack-level Gatekeeper policy;
- aggregate support-state promotion;
- distinction among implementation, candidate evidence and empirically supported profiles;
- locking exact component SHAs only after promotion criteria are met.

## Candidate representation examples

Tlaloc is explicitly allowed to test unconventional visual strategies.

Examples include:

```text
new color semantics
alternative shapes
spatial hierarchy changes
redundancy changes
prime-derived patterns
modular/factorization structures
periodic density patterns
phase or temporal encodings
```

These are experiments, not new per-document Origami dialects.

A candidate only matters if it improves measurable properties such as:

```text
semantic roundtrip
VLM readability
routing accuracy
transport robustness
context efficiency
selective access
carrier density
cross-model transfer
```

while preserving `FALSE_EXACT=0` and other hard semantic invariants.

## Visual/perception boundary

Within Origami, **OHF (Origami HyperFold)** remains a carrier/protocol and executable research track. R3.10-LAB, Perception Lab, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and visual-memory experiments remain evidence-gated assets.

The preferred end-to-end receiver is Hybrid:

```text
model perception
  -> BOOT / ROSETTA / coarse route
Tlaloc / Origami tools
  -> bounded query / expand / verify
Origami semantics
  -> state / Fold-Unfold / evidence / canonical visual contract
model
  -> integrate / ask next bounded question / answer
```

Native and Computational remain separate diagnostic/promotion modes. Native visual success must not be inferred from Hybrid tool success, and Hybrid support must not require Native T3 promotion.

## READ / WRITE boundary

```text
READ
carrier -> Origami profile + ROSETTA -> semantic interpretation

WRITE
source -> Tlaloc/document adapter -> Semantic IR
       -> Origami Writer -> canonical visual plan
       -> Origami compiler -> carrier -> roundtrip verify
```

A language model may propose Semantic IR and visual intents. It is not the pixel-level carrier authority.

## Anti-shortcut rule

A screenshot, text poster, glyph payload or exact literal transport can be useful as transport evidence, but none of them is semantic Fold merely because a decoder can recover bytes.

```text
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PER-DOCUMENT PRIVATE STYLE != CANONICAL ORIGAMI PROFILE
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
```

Origami is the authority for its semantic and canonical visual contracts. Tlaloc and Tonal may build adapters, search systems, fixtures, evaluators, mirrors and promotion gates around those contracts, but they must preserve the distinction between implemented machinery, candidate strategies, empirical evidence and promoted support.
