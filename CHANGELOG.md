# Origami changelog

## 6.0.0-alpha.13 — Self-Describing Read/Write Protocol R0 + profile-3 renderer

### Added

- formalizes `origami.protocol.r0`: Origami is a bidirectional, self-describing communication protocol rather than only a visual memory representation;
- adds `origami.codec-registry.r0` with semantic decoder family `S0..S6`, semantic encoder family `E0..E6`, exact decoder family `X0..X5` and exact encoder family `Q0..Q5`;
- adds `origami.capability-negotiation.r0` so receivers may declare supported codecs/channels and fail with `UNKNOWN` / `NOT_VERIFIED` instead of invention;
- makes ROSETTA declare decoding semantics **and** encoding semantics;
- adds deterministic reference Go registry under `internal/codec` and `origami-codec` CLI;
- implements the first paired semantic roundtrip gate `S2(E2(INDEX)) ~= INDEX`;
- upgrades Master Prompt to R4 protocol/codec-aware handshake, SHA-256 `78cbcee8d0f71e1b06e2902c9b9600779e3cd0148251ad9d50623d07969c2a56`;
- adds `origami.fixed-carrier.r2.profile-3` as a parallel rendered candidate with visible T0 protocol boot, T1 ROSETTA + codec entry points, T2 semantic superindex, PROGRAM codec hints, separate T3 exact/control and VERIFY;
- adds `origami-profile3-carrier` build/decode CLI;
- profile-3 deterministic CI preserves the 640x640 canvas, exactly 8192 PNG bytes and profile-1/profile-2 decode compatibility.

### Corrected interpretation

- decoder use itself is not a failure;
- a self-declared semantic decoder such as `S2` is part of Origami and may be used natively;
- undeclared external decoder/file/binary dependency remains forbidden for Native semantic navigation;
- a semantic query must not unnecessarily escalate to exact/binary codecs;
- encoder instructions are first-class so a compatible model can construct a new Origami Construction IR even when no deterministic image compiler is available.

### Evidence boundary

Profile-3 renderer/roundtrip/size compatibility is deterministically implemented and green. Held-out Native S2 index recovery, held-out E2 write behavior and real A -> B -> C cross-model interoperability remain external-model evidence pending. Profile-2 remains the default renderer and profile-3 is not promoted.

## 6.0.0-alpha.12 — Native Semantic Navigation R0

### Failure-driven correction

- preserves `FAILED_TRIAL_001`, where a real multimodal model read BOOT but failed the index question, requested binary/file decoding and emitted corrupted/unverified mechanical claims;
- changes the Native path from binary-first to semantic-first;
- makes Fixed Carrier profile-2 T2 render actual bounded semantic entry labels from canonical memory instead of generic PAGE/CONCEPT/SOURCE/GRAPH/VERIFY categories;
- keeps T3 as deterministic exact/control plane rather than a prerequisite for semantic index/overview questions;
- adds `origami.native-semantic-navigation.r0` and the first hard Native benchmark `What is the index?`;
- updates the integration layer so the compact GraphSignature remains structurally valid JSON instead of blindly truncating JSON at byte 256;
- adds Master Prompt R3 semantic-first behavior and blocks invented byte counts, compression algorithms, hidden hashes/archive claims when exact verification was not executed;
- preserves profile-1 deterministic decode compatibility and the 640x640 / 8192-byte envelope.

### Evidence boundary

Alpha.12 corrects the architecture and deterministic regression harness. It does not by itself prove that held-out VLMs now recover T2 natively.

## 6.0.0-alpha.11 — Portable Development-Tool Boundary R0

### Added / corrected

- adds `origami.development-tool-boundary.r0`;
- makes the universal Master Prompt R2 the portable compatibility baseline;
- explicitly requires no Tlaloc, Tonal, sandbox, Go/Python runtime or external tools for baseline prompt compatibility;
- distinguishes prompt/native capability from optional tool-assisted/Hybrid capability;
- READ without tools now returns `UNKNOWN` or `NOT_VERIFIED` when an unavailable mechanical/exact operation is necessary;
- WRITE without a compiler may emit Semantic IR, visual intent, ROSETTA and construction specification but must report `CONSTRUCTION_SPEC_ONLY` rather than claiming a compiled carrier;
- generalizes Evidence Reduction and Semantic Spine external proposal/orchestration boundaries beyond Tlaloc;
- makes Tlaloc, Blueprint Framework and future systems optional development tools while Origami owns its own semantics, Master Prompt, canonical profile and releases;
- clarifies Tonal as optional multi-tool composition/reproducibility rather than Origami profile authority;
- preserves the existing `tlaloc.origami-tools.r2` path as an optional external provider ABI.

### Invariants

- `ORIGAMI_OWNS_ORIGAMI_RELEASES`;
- `MASTER_PROMPT_IS_PORTABLE_BASELINE`;
- `DEVELOPMENT_TOOL != RUNTIME_REQUIREMENT`;
- `TOOL_ASSISTED_SUCCESS != PROMPT_ONLY_SUCCESS`;
- `FALSE_EXACT=0`.

### Evidence boundary

Alpha.11 defines the portable/no-hidden-runtime contract. It does not claim that prompt-only operation already matches Hybrid/tool-assisted operation on every task or that every target model can natively read/write every channel.

## 6.0.0-alpha.10 — Perceptual Profile Evolution + ROSETTA Reveal Semantics

### Added / connected

- connects the alpha.2 interference, depth, temporal and emergent perceptual-channel contracts to canonical visual-profile evolution;
- adds experimental canonical dimension families for interference, depth and emergent perception alongside color, numeric and temporal channels;
- extends Writer/ROSETTA planning with reveal procedures for non-static channels;
- reveal procedures can declare operations such as `MOIRE`, `PHASE_SHIFT`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL` and `TEMPORAL_INTEGRATE` plus required inputs/conditions;
- failed reveal remains `UNKNOWN`, never semantic absence;
- formalizes profile optimization toward more recoverable semantics/readability/robustness with fewer bytes, recognition time, BOOT/decode steps and model-facing context;
- corrects project authority: external development tools propose and measure; Origami validates/adopts its own profile versions; Tonal is optional composition.

### Evidence boundary

No moiré, stereo, temporal, numeric or color candidate is promoted merely by this release. Advanced channels remain evidence-gated.

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

### Evolution boundary at alpha.9

- Tlaloc was the primary experimental search implementation for prompt/channel/layout/color/shape/numeric/temporal strategies;
- Origami owned validation and canonical visual-profile semantics;
- later releases clarified that Origami itself owns profile promotion and Tonal is only optional toolchain composition;
- aesthetic changes are versioned profile changes, not private styles for individual documents;
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
- Tlaloc was the primary cross-model campaign/transport/tool-loop implementation;
- Origami remains capability/profile authority; later releases clarify Tonal as optional composition only;
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
- Tlaloc remains an available external Canonical Document IR/OCR, Tlaloque, Merkle exact plane, CanonicalState, External Recursive Attention H0-H5 and `tlaloc.origami-tools.r2` implementation; alpha.11 later generalizes this as optional external development/tool infrastructure.

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
- structured external proposals with `SUPPORT`, `OPPOSE` and `UNKNOWN` positions;
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

- external model/development systems may generate proposals;
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
- Native VLM readability remains unpromoted until held-out evidence exists;
- later releases classify the Tlaloc tool ABI as an optional provider rather than a universal Origami dependency.

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

Behavior compilation, PromptIR, Tlaloque coordination/training and target-model experiment orchestration belong to external development tooling such as Tlaloc. Origami owns Origami semantics/releases.

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
