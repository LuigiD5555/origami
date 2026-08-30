# Origami changelog

## 6.0.0-alpha.9 — Canonical Visual Grammar R0 + Writer R0

### Added

- `origami.canonical-visual-grammar.r0` defining one canonical Origami aesthetic per promoted profile version;
- explicit distinction between functional visual encoding and arbitrary per-document styling;
- canonical visual dimensions for geometry, fill/contrast, position, topology, enclosure, scale, repetition, density and limited declarative text;
- color, numeric/mathematical structure and temporal channels as experimental candidate dimensions rather than automatic semantic authority;
- canonical primitive roles matching the current Fixed Carrier visual language;
- mandatory ROSETTA as self-description of profile/version, active dimensions and concrete bindings without permitting arbitrary aesthetic drift;
- `origami.writer.r0` READ/WRITE construction contract;
- deterministic `internal/writer` planner mapping Semantic IR intents onto the canonical visual grammar;
- `origami-write-plan` CLI;
- semantic roundtrip requirement `DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR`;
- Master Prompt R1 with READ and WRITE modes;
- explicit PDF/source -> Semantic IR -> visual intent -> canonical grammar -> ROSETTA -> compiler -> roundtrip verification workflow.

### Corrected

- replaced the earlier interpretation that every carrier may invent a completely different visual dialect;
- profile-bound symbols now have versioned canonical meaning while ROSETTA remains mandatory for self-description and compatibility;
- a model may not directly paint an arbitrary PNG and call it a valid Origami carrier;
- a PDF screenshot or text poster remains forbidden as an Origami semantic representation.

### Evolution boundary

- Tlaloc owns experimental search over candidate prompt/channel/layout/color/shape/numeric/temporal strategies;
- Origami owns validation and canonical visual-profile semantics;
- Tonal owns aggregate stack/profile promotion;
- aesthetic changes are promoted as new canonical profile versions, not private styles for individual documents;
- old carriers remain readable through embedded profile/version + ROSETTA.

## 6.0.0-alpha.8 — Perception Promotion R1

### Added

- deterministic `origami.perception-promotion.r1` evaluator;
- T0 BOOT comparison against actual Fixed Carrier ground truth;
- duplicated top/bottom visual probe checks and agreement gate;
- exact Tool Protocol and Address ABI checks;
- optional T3 machine-record comparison against the actual decoded carrier;
- explicit `MOCK` vs `REAL_MODEL` evidence kinds;
- independent Hybrid and Native-T3 per-trial eligibility;
- `origami-perception-eval` CLI.

### Promotion boundary

- a mock may pass every mechanical check but can never become empirical promotion evidence;
- a real single trial is only one eligible trial, not a support claim;
- Hybrid trial success does not require Native T3;
- Tlaloc owns cross-model campaigns/transport variants/real tool loops;
- Tonal owns final aggregate stack promotion;
- real external VLM evidence remains pending.

## 6.0.0-alpha.7 — Semantic Spine R1

### Restored / connected

- executable `S_(t+1)=F(S_t,C_t,R)` core with first-class Context and declarative Rule IDs;
- distinct `PRESENT`, `ABSENT`, `UNKNOWN`, `INHIBITED`, `CANCELLED` states;
- declarative contradictions, signed contribution/cancellation evidence and higher-order participant-set rules;
- finite Observation Contracts separated from transition;
- temporal-sequence observations over ordered trajectories;
- exhaustive finite combinatorial exploration oracle without silent pruning;
- semantic graph Fold with commitment, dependencies, relations and unresolved alternatives;
- selective Unfold with touched-address accounting so hidden whole-graph loading cannot masquerade as selective access;
- neutral exact source plane and source->Virtual Memory projection without automatically promoting extracted text to semantic summaries;
- semantic graphs and deterministic trajectories projected directly into Virtual Memory;
- Evidence Reduction results promoted back into memory as canonical claims;
- perception observations accepted as addressable `PERCEPTION_ONLY` records but never semantic authority by themselves;
- Fixed Carrier R2 bound explicitly to the external canonical memory root rather than treated as the corpus.

### Anti-shortcut corrections

- formalized `SCREENSHOT != ORIGAMI_SEMANTIC_REPRESENTATION`;
- formalized `LITERAL_TRANSPORT != SEMANTIC_FOLD`;
- a cited exact-fidelity block no longer grants exactness to a paraphrase;
- `VERIFIED_EXACT` now requires the selected proposition to be byte-for-byte equal to accepted verified exact evidence;
- raw PDF/OCR text does not receive an automatic semantic `summary` merely because it was extracted.

### Preserved work

- Virtual Memory R0, Memory Scale Lab R0, Evidence Reduction R0 and Fixed Carrier R2 remain intact;
- OHF R3.10-LAB, Perception Lab, visual-family experiments, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA and Context SIMD remain active evidence-gated research assets;
- Tlaloc remains the existing Canonical Document IR/OCR, Tlaloque, Merkle exact plane, CanonicalState, External Recursive Attention H0-H5 and `tlaloc.origami-tools.r2` orchestration/data-plane implementation; alpha.7 does not fork that pipeline into Origami.

### New tools/contracts

- `origami.semantic-spine.r1`;
- `origami.document.r0` neutral source format;
- `internal/core`, `internal/document`, `internal/integration`;
- `origami-ingest`;
- `origami-system`.

### Evidence boundary

Alpha.7 provides a connected deterministic semantic reference. Built-in arbitrary scan/image OCR inside Origami, complete alpha.2 perceptual operations, held-out Native visual navigation, cross-model Hybrid/Native promotion and optimized large combinatorial search remain evidence-pending.

## 6.0.0-alpha.6 — Evidence Reduction R0

### Added

- generalized `origami.evidence-reduction.r0` contract;
- structured external swarm proposals with `SUPPORT`, `OPPOSE` and `UNKNOWN` positions;
- Origami-owned evidence resolution against address, CID, source hash and requested fidelity;
- verified-evidence-only reduction authority;
- deterministic evidence deduplication so repeated agents cannot manufacture consensus from one source;
- independent-source thresholds;
- canonical claim states `VERIFIED`, `REJECTED`, `CONFLICT`, `UNKNOWN`;
- explicit preservation of verified conflicts instead of majority-vote erasure;
- order-deterministic canonical claim IDs and whole-reduction digest;
- separate exactness state from claim verification;
- model-neutral `origami-reduce` CLI;
- deterministic tests for order invariance, duplicate collapse, identity mismatch, unverified evidence, conflict preservation and independent sources.

### Boundary

- Tlaloc continues to own swarm/Tlaloque search, prompts, exploration and proposal generation;
- Origami owns addressable evidence identity, verification, deterministic reduction and exactness state;
- agent confidence is diagnostic and never substitutes for verified evidence.

### Invariants

- `SWARM_PROPOSAL != CANONICAL_STATE`;
- `AGENT_CONFIDENCE != EVIDENCE_AUTHORITY`;
- `DUPLICATE_EVIDENCE_COUNTS_ONCE`;
- `UNKNOWN > unsupported invention`;
- `FALSE_EXACT=0`;
- reduction semantics are independent of proposal arrival order.

### Evidence boundary

Alpha.6 implements the deterministic reducer. Alpha.7 later tightens exactness from “verified exact reference” to byte-equality between the canonical proposition and accepted exact evidence.

## 6.0.0-alpha.5 — Fixed Carrier R2

### Added

- frozen `origami.fixed-carrier.r2` visual control plane;
- fixed 640×640 canvas and exactly 8192-byte PNG profile under the existing 500 KB hard ceiling;
- T0 plaintext BOOT, T1 Rosetta/duplicated visual challenge, T2 root index/graph navigation, T3 deterministic machine record and VERIFY stage;
- `tlaloc.origami-tools.r2` and address ABI binding;
- deterministic fixed-carrier encoder/decoder and CLI;
- reversible installer integration for the fixed-carrier binary.

### Clarified

- corpus growth occurs in the external canonical memory plane rather than by growing the visual canvas;
- OCR is optional and is not BOOT authority;
- Native VLM readability remains unpromoted until held-out evidence exists.

## 6.0.0-alpha.4 — Virtual Memory R0

### Added

- model-agnostic `origami.virtual-memory.r0` contract;
- canonical address ABI with address/CID separation;
- multi-resolution node fidelity: `label -> abstract -> summary -> detail -> evidence -> exact`;
- explicit `ContextBudget` and portable `ContextPacket`;
- metadata-only semantic indexing that avoids implicit exact-payload scans;
- query-intent routing for causal, evidence, contradiction and exact-access needs;
- utility-per-token selection and fidelity downgrade before budget violation;
- bounded token-aware working-set cache;
- compact per-carrier `GraphSignature` and multi-carrier federation;
- `origami-memory` model-agnostic CLI;
- Hybrid runtime `QUERY`, `EXPAND` and `SIGNATURE` operations while preserving legacy `LOOKUP/FOLLOW/TRACE`;
- experimental macro/meso/micro visual memory navigation projection;
- one-PNG memory carrier combining a redundant navigation band with the unchanged exact self-decodable transport;
- machine-readable `origami.visual-memory-profile.r0` mapping all F01–F42 families to memory/navigation roles without promoting them as independent reliable channels;
- `origami.virtual-memory-nav-eval.r0` for carrier/cluster selection, graph following, address reopening, distractor rejection, resize robustness and fixed-context scale tests.

### Compatibility

- `origami.hybrid-carrier.r0` remains the carrier schema;
- legacy `Address/Value/Links` memory remains readable and queryable;
- legacy `Render/DecodePNG` remains available;
- rich virtual-memory carriers use optional memory-profile metadata and `RenderMemory/DecodeMemoryPNG`;
- the original Master Prompt remained receiver-only at this stage.

### Invariants

- `FALSE_EXACT=0` remains mandatory;
- active model-facing context is not total Origami memory;
- address is not content identity;
- no missing query/relation may fall back to an undeclared global exact scan;
- visual navigation is redundant and never exactness authority;
- the 42 visual families remain evidence-gated generator-tree candidates;
- model specialization is optional: any model can consume a `ContextPacket`.

### Evidence boundary

Deterministic runtime, PNG round-trip and budget tests are part of this release candidate. Native VLM navigation of the macro/meso/micro projection is **not promoted** until held-out cross-model Perception Lab evidence exists.

## 6.0.0-alpha.3 — Repository hierarchy and state reconciliation

### Corrected

- reconciled the executable OHF R3.10-LAB workbench with the Origami 6.x project boundary;
- made the repository hierarchy explicit: Origami is the representation/state-machine language; OHF is a research carrier/protocol track inside Origami;
- replaced the stale OHF-only `state/ORIGAMI_STATE.json` root model with a project-level state that tracks semantic contracts and the OHF substate separately;
- repaired README and documentation that described R3.10/OHF as the current state or architecture of all Origami;
- clarified that `install.sh` / `uninstall.sh` manage OHF laboratory CLIs, not the complete future Origami runtime;
- preserved the detailed pre-reconciliation OHF R3.10-LAB architecture/state/development/installation/roadmap documents under `docs/ohf/`;
- preserved R3.10-LAB, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible-installation work as active experimental assets rather than discarding them.

### Preserved

- `origami.perceptual-channels.r0` semantics from alpha.2;
- coherent-state semantics from alpha.1;
- Tlaloc/Origami ownership boundary;
- historical R2/R3.10 evidence and change records;
- existing Go experimental behavior.

## 6.0.0-alpha.2 — Perceptual Channels R0

### Added

- perceptual-channel contract `origami.perceptual-channels.r0`;
- spatial, interference, depth, temporal and emergent channel families;
- explicit-vs-emergent channel distinction;
- Temporal Latent Image (TLI) semantic object;
- perceptual availability conditions: `latent`, `masked`, `motion-bound`, `phase-bound`, `stereo-bound`, `emergent`, `revealed`, `decayed`;
- operations `MOIRE`, `PHASE_SHIFT`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE`, `TEMPORAL_DECAY`;
- machine-readable `spec/PERCEPTUAL_CHANNELS_R0.json`;
- invariants preventing static-frame assumptions and false absence claims.

### Clarified

- coherent-state kinds and perceptual conditions are separate semantic layers;
- a valid Origami percept may exist only across a declared temporal trajectory;
- an emergent percept belongs to the interaction that produces it, not automatically to one source layer;
- runtime support for these operations is incomplete and unsupported behavior must be explicit.

## 6.0.0-alpha.1 — State Semantics R0

### Added / formalized

- `determinate`, `superposed`, `coupled` and `observed` states;
- complex branch amplitudes for the quantum-inspired profile;
- `TRANSFORM` without implicit observation;
- `INTERFERE` and cancellation;
- `CONSTRAIN`;
- explicit `OBSERVE`;
- `FOLD` / `UNFOLD` relationship to state resolution/expansion.

### Removed from Origami ownership

Behavior compilation, PromptIR, Tlaloque coordination/training, target-model orchestration and behavioral-artifact promotion belong to Tlaloc.

## Historical OHF / R3.10 laboratory milestones

### 2026-08-27 — Reversible installation lifecycle

- added tracked SHA-256 restore semantics for OHF laboratory binaries;
- preserved the original pre-install rollback point across reinstalls;
- added conflict detection and opt-in legacy `.bak` removal support;
- kept normal smoke runs ephemeral.

### 2026-08-27 — Dimensional Visual Register / PSS formalization

- added hierarchical MICRO/MESO/MACRO visual-dimension registry;
- added executable DVR validation and scope-aware capacity reporting;
- separated nominal PSS from evidence-gated `SAFE_PSS`;
- added Perceptual Orthogonality Law and pair-evidence model;
- integrated Macro-Gestalt as a MACRO-scale experimental candidate;
- added `CHG-DVR-0001` and local inspection gates.

### 2026-08-27 — R3.10-LAB repository modernization

- moved the public OHF research track from R2 documentation-only status to an executable Go laboratory;
- added reproducible experiment infrastructure, Native isolation, observations/scoring and failure-to-regression plumbing;
- added Glyph Calculus, SAFE_MICRO_ISA and Context SIMD experiments;
- added conservative external-execution policy, installer, Makefile and CI;
- changed the Go module path to `github.com/LuigiD5555/origami`;
- retained R2 material as explicit legacy evidence;
- added machine-readable OHF/LAB state and change records.

### 2026-08-26 — R2 public documentation baseline

- added scientific foundations document;
- added R2 Master Prompt;
- added practical R2 README.
