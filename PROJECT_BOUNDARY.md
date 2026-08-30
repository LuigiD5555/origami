# Origami project boundary

Origami, Tlaloc and Tonal occupy different levels and must not be treated as one system.

```text
                           TONAL
          optional development-tool composition layer
        /                 |                    \
   TLALOC          Blueprint Framework      other tools
 development kit      development tool      future tools
        |
        | experiments / candidates / ingestion / evaluation
        v
                           ORIGAMI
          representation language + machine + memory
          owns its own canonical semantic/visual versions
```

## Stable hierarchy

**Origami** is the thing being built and evolved: the representation/state-machine language, semantic memory, Fold/Unfold system, evidence rules, ROSETTA, visual grammar, Writer and carrier compiler/runtime.

**Tlaloc** is a development kit and experimental work system. It can be used to improve Origami, but it is not defined only by Origami. The same orchestration, swarm, behavior-compilation, experiment and evaluation machinery may also be used to develop other tools.

**Tonal** is optional and sits above development tools. It exists so multiple independent or alternative development systems can be composed, pinned and verified together—for example Tlaloc, Blueprint Framework and future tooling. Tonal does not own Origami semantics and does not decide when an Origami profile becomes canonical.

Origami may be used without Tlaloc. Tlaloc may be used without Origami. Tonal is needed only when a reproducible composition of multiple tools/components is useful.

## Canonical visual ownership

Origami owns **one canonical functional visual aesthetic per promoted profile version**.

Tlaloc may experimentally mutate and benchmark candidate changes to:

```text
Master Prompt
shape / primitive roles
layout
channel assignments
redundancy
color usage
numeric / mathematical structures
moire / interference / phase
stereo / parallax / depth
temporal / motion-bound structures
emergent multi-layer or multi-instant percepts
```

A Tlaloc result is candidate evidence, not a canonical Origami change.

The lifecycle is:

```text
Origami canonical profile N
        |
        v
Tlaloc or another development tool
search / experiment / benchmark
        |
        v
candidate profile N+1 + evidence
        |
        v
Origami semantic + visual validation
        |
        v
Origami promotes profile N+1
        |
        +------------------------------+
        |                              |
        v                              v
standalone Origami              optional Tonal composition
                                with Tlaloc / Blueprint /
                                other development tools
```

A successful experiment does not directly mutate the canonical renderer, grammar or Master Prompt. Origami remains the authority over its own versions.

## ROSETTA ownership

ROSETTA remains Origami-owned and mandatory.

For a canonical visual profile it declares:

```text
profile ID/version
active dimensions
canonical visual primitive -> semantic role mappings
carrier-approved parameters/extensions
reveal procedure for any promoted non-static perceptual channel
```

For interference, depth, temporal or emergent channels, ROSETTA must describe how the percept is revealed as well as what it means. If the receiver cannot satisfy the declared reveal condition, the result is `UNKNOWN`, not semantic absence.

ROSETTA is a self-description/compatibility mechanism, not permission for arbitrary per-document visual dialects. Old carriers remain readable because profile/version + ROSETTA travel with them.

## Preserved Tlaloc R2 implementation

Tlaloc's existing orchestration/data-plane remains external and is not replaced by Origami Writer:

- layout-preserving Canonical Document IR for digital PDFs with OCR fallback for text-sparse/raster pages;
- stable document/page/region/block addressing and CIDs;
- figures/geometry/canonical text/original source preservation;
- proposal-only Tlaloque candidate generation;
- deterministic CanonicalState, conflicts, uncertainty/verification queues and state ledger;
- Merkle-bound exact plane with `FALSE_EXACT=0`;
- External Recursive Attention H0-H5 under a bounded active context;
- `tlaloc.origami-tools.r2` BOOT/QUERY/EXPAND/VERIFY provider;
- OpenAI-compatible multimodal function loop and plain-text tool bridge;
- Perception Promotion Campaign machinery;
- canonical visual-profile search and experimental representation tournaments.

These are Tlaloc development/runtime capabilities. They may support Origami development but do not define Origami semantics.

Origami Writer consumes source/semantic structures and defines how they become a canonical visual construction plan; it does not fork Tlaloc's document ingestion pipeline.

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
- exactness boundaries and `FALSE_EXACT=0`;
- self-boot receiver contract;
- universal Master Prompt READ/WRITE semantics;
- canonical visual profile and versioned aesthetic;
- geometry/color/numeric/interference/depth/temporal/emergent dimension semantics when promoted;
- ROSETTA semantics and reveal procedures;
- Writer construction-plan semantics and semantic roundtrip requirements;
- Fixed Carrier/control-plane semantics and StoreRoot binding;
- profile validation and promotion of new Origami versions;
- promoted receiver/visual artifacts and provenance.

## Tlaloc owns

As a development kit/work system, Tlaloc may provide:

- document acquisition/canonicalization and OCR;
- swarm/Tlaloque exploration;
- behavior compilation and distillation;
- Master Prompt search/mutation;
- BOOT/ROSETTA strategy experiments constrained by the target contract;
- visual primitive/channel/layout/color/numeric/interference/depth/temporal search;
- claim proposal generation and ambiguity exploration;
- uncertainty-driven verification scheduling;
- External Recursive Attention orchestration;
- model/tool-loop transports;
- cross-model perception campaigns and degradation experiments;
- candidate representation tournaments and fitness measurement;
- recommendations and regression coordination.

These capabilities can target Origami or other development targets.

## Tonal owns

Tonal is not the semantic or profile authority for Origami. Its role is ecosystem/toolchain composition:

- pin exact revisions of multiple independent development tools/components;
- express compatible/reproducible compositions;
- run stack-level integration verification;
- build snapshots/distributions of a chosen toolchain;
- allow Tlaloc, Blueprint Framework and future alternatives/complements to coexist without pretending they are one repository;
- record which exact toolchain was used to produce or validate an artifact when provenance matters.

Tonal may record that an Origami version is supported by a given composed toolchain, but **Origami itself decides which Origami version/profile is canonical**.

## Profile-evolution objective

A new Origami version is useful only when it improves measurable properties while preserving semantic correctness.

Tlaloc or another development kit may optimize:

```text
maximize:
  recoverable semantic capacity per byte
  semantic roundtrip
  VLM/readability success
  routing and verified-evidence accuracy
  transport robustness
  perceptual reveal reliability

minimize:
  carrier bytes
  recognition latency
  BOOT steps
  decode/unfold steps
  model-facing context
```

Hard invariants such as `FALSE_EXACT=0`, evidence authority and `UNKNOWN` discipline cannot be traded away for a better score.

## Visual/perception boundary

Within Origami, OHF/Perception Lab remains a carrier/protocol and perceptual research track. Interference, depth, temporal and emergent channels are valid candidates because Origami already models operations such as MOIRE, PHASE_SHIFT, STEREO_BIND, PARALLAX_RESOLVE, KINETIC_REVEAL and TEMPORAL_INTEGRATE.

The preferred Hybrid receiver path remains:

```text
model perception
  -> BOOT / ROSETTA / coarse route
external tool provider when available
  -> bounded query / expand / verify
Origami semantics
  -> state / Fold-Unfold / evidence / canonical visual contract
model
  -> integrate / ask next bounded question / answer
```

Tlaloc is one strong provider of those development/runtime tools, not the definition of the interface itself.

## READ / WRITE boundary

```text
READ
carrier -> Origami profile + ROSETTA -> semantic interpretation

WRITE
source -> external ingestion/semantic extraction when needed
       -> Origami Writer -> canonical visual plan
       -> Origami compiler -> carrier -> roundtrip verify
```

A language model may propose Semantic IR and visual intents. It is not the pixel-level carrier authority.

## Anti-shortcut rule

```text
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PER-DOCUMENT PRIVATE STYLE != CANONICAL ORIGAMI PROFILE
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
DEVELOPMENT TOOL != ORIGAMI SEMANTIC AUTHORITY
TONAL COMPOSITION != ORIGAMI PROFILE PROMOTION
```
