# Origami project boundary

```text
TLALOC                                      ORIGAMI
work/orchestration system                   representation/state-machine language
Canonical Document IR / OCR         ----->  source/representation contracts
swarm/Tlaloque search               ----->  receiver / claim candidate inputs
behavior compilation + distillation ----->  prompt/boot/micro-agent candidates
ERA H0-H5 / tool loop               <---->  address/memory/evidence contracts
model-facing evaluation             <-----  semantic/reference contracts
promotion recommendation            ----->  candidate evidence
                                             |
                                             +-- semantic validation
                                             +-- evidence address resolution
                                             +-- deterministic claim semantics
                                             +-- carrier execution
                                             +-- promoted artifact storage
```

Origami may be used independently of Tlaloc. Tlaloc may operate without Origami or consume Origami as one representation provider.

The stable rule is:

> Tlaloc discovers, ingests, coordinates and distills work; Origami defines and executes representation semantics, including how state evolves, how Fold/Unfold behaves, how memory is addressed, how evidence earns authority, how a carrier bootstraps and how exactness is claimed.

## Preserved Tlaloc R2 implementation

Tlaloc `6.0.0-alpha.11` already contains the production-oriented R2 orchestration/data-plane work and is not replaced by Origami alpha.7:

- layout-preserving Canonical Document IR for digital PDFs with OCR fallback for text-sparse/raster pages;
- stable document/page/region/block addressing and CIDs;
- figures/geometry/canonical text/original source preservation;
- proposal-only Tlaloque candidate generation;
- deterministic CanonicalState, conflicts, uncertainty/verification queues and state ledger;
- Merkle-bound exact plane with `FALSE_EXACT=0`;
- External Recursive Attention H0-H5 under a bounded active context;
- `tlaloc.origami-tools.r2` BOOT/QUERY/EXPAND/VERIFY provider;
- OpenAI-compatible multimodal function loop and plain-text tool bridge.

Those components remain Tlaloc-owned orchestration/runtime assets. Origami Semantic Spine R1 supplies the semantic/evidence rules that a composed system must satisfy; it does not fork Tlaloc's PDF/OCR implementation into a second competing pipeline.

Where Tlaloc carries a deterministic reducer for its R2 runtime, that implementation is treated as an orchestration-side implementation/mirror. Its promoted outputs must remain compatible with Origami-owned evidence/exactness semantics rather than creating an independent definition of truth.

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
- Fixed Carrier R2 visual/control-plane semantics and StoreRoot binding;
- carrier-local symbol bindings;
- deterministic micro-agent execution semantics where generalized into Origami;
- promoted receiver artifacts and their provenance;
- evidence/exactness boundaries.

## Tlaloc owns

- document acquisition/canonicalization workflow and OCR execution;
- swarm/Tlaloque exploration;
- behavior-compilation lifecycle;
- candidate prompt search/mutation;
- candidate BOOT/Rosetta strategy search;
- claim proposal generation and model-side ambiguity exploration;
- uncertainty-driven verification scheduling;
- External Recursive Attention orchestration;
- model/tool-loop transports;
- distillation of complex swarm traces into simple bounded micro-agent candidate rules;
- cross-model evaluation/tournaments;
- promotion recommendations and regression coordination.

A swarm proposal is never a canonical Origami fact merely because many agents repeat it. Origami evidence semantics reopen cited addresses, check identity/verification state, deduplicate repeated evidence and preserve verified disagreement. Agent confidence remains diagnostic rather than authoritative.

## Visual/perception boundary

Within Origami, **OHF (Origami HyperFold)** remains a carrier/protocol and executable research track. R3.10-LAB, Perception Lab, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and visual memory experiments remain active evidence-gated assets.

The preferred end-to-end receiver is Hybrid:

```text
model perception
  -> BOOT / Rosetta / coarse route
Tlaloc / Origami tools
  -> bounded query / expand / verify
Origami semantics
  -> representation / state / Fold-Unfold / evidence authority
model
  -> integrate / ask next bounded question / answer
```

Native and Computational remain separate diagnostic/promotion modes. Native visual success must not be inferred from Hybrid tool success, and Hybrid support must not require Native T3 promotion.

## Anti-shortcut rule

A screenshot, text poster, glyph payload or exact literal transport can be useful as transport evidence, but none of them is semantic Fold merely because a decoder can recover bytes.

```text
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
```

Origami is the authority for its semantic contracts. Tlaloc and Tonal may build adapters, fixtures, evaluators, mirrors and promotion gates around those contracts, but they must preserve the distinction between implemented machinery, empirical evidence and promoted support.
