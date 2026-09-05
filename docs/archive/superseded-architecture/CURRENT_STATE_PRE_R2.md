# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.12  
**Status:** experimental

## Project role

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory system.

Origami owns its own semantics, canonical visual grammar, ROSETTA, Writer/Reader contracts, Master Prompt and version/profile releases. Tlaloc, Blueprint Framework and future systems may help develop it; Tonal may optionally compose/pin those tools. None is required for Origami's portable baseline.

## Alpha.12: Native Semantic Navigation R0

A real external multimodal trial established a failure that alpha.11 did not prevent:

```text
BOOT readable                      yes
requested semantic index           yes
actual semantic index recovered    no
binary/file decoder requested      yes
unverified mechanical claims       yes
```

The failure is now stored as `experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json` and is a permanent regression.

The corrected Native path is:

```text
semantic question
 -> T0
 -> T1 / ROSETTA
 -> T2 actual semantic superindex
 -> smallest relevant semantic region
 -> answer
```

Binary/exact decoding is **not** a prerequisite for identity, index, overview or topic-location questions.

## Current capability state

<!-- BEGIN GENERATED CLAIMS TABLE: do not edit; run python3 tools/claims.py generate -->
| Claim | Statement | Status | Evidence | Version introduced | Last checked | Notes |
|---|---|---|---|---|---|---|
| `ORIGAMI.ADDRESSING.PAGE_RESOLVE` | Resolve a requested logical book page to the correct physical source offset without assuming a constant printed-page offset. | `designed` | — | `unknown` | 2026-09-01 | A historical failure artifact exists but does not satisfy the normalized run_id requirement for evidenced_failing; see ORIGAMI-DISCREPANCY-004. |
| `ORIGAMI.BOOK.CARRIER_RECOVERY` | Recover the referenced book corpus from its compact carrier and verify the recovered bytes by hash in a normalized model run. | `designed` | — | `unknown` | 2026-09-01 | The evidence plan reports a 7.55 MB to approximately 474 KB success, but published main has no matching normalized run_id or directly traceable artifact; see ORIGAMI-DISCREPANCY-002. |
| `ORIGAMI.CANDIDATE.EXPERIMENTAL_BUILDER` | Compile supported experimental representation mutations while preserving the parent temporal program and frozen carrier envelope. | `implemented` | `test:internal/fixedcarrier:TestTemporalCandidateBuildPreservesExactProgram`<br>`test:internal/fixedcarrier:TestTemporalCandidateBuildSupportsTemporalAndNumericSignals`<br>`test:internal/fixedcarrier:TestTemporalCandidateRejectsUnsupportedMutation`<br>`test:internal/fixedcarrier:TestTemporalCandidateRejectsParentMismatch` | `6.0.0-alpha.15` | 2026-09-01 | A successful build remains experimental and does not establish model improvement or promotion. |
| `ORIGAMI.CAPACITY.BOUNDED_NAVIGATION` | Enforce bounded superindex navigation and return unknown instead of falling back to a global scan. | `implemented` | `test:internal/lab/capacity:TestSuperIndexDoesNotGlobalScan`<br>`test:internal/lab/capacity:TestWorkingWindowIsEnforced`<br>`test:internal/lab/capacity:TestUnknownDoesNotFallbackToScan` | `6.0.0-alpha.4` | 2026-09-01 | The bounded-navigation harness does not establish Native VLM addressing performance. |
| `ORIGAMI.CARRIER.FIXED_R2` | Encode and decode a self-contained 640 by 640 PNG with an exactly 8192-byte frozen envelope and corpus-independent carrier size. | `implemented` | `test:internal/fixedcarrier:TestR2RoundTripFrozenSizeAndBootBinding`<br>`test:internal/fixedcarrier:TestCarrierSizeDoesNotGrowWithCorpusOrEntropy`<br>`test:internal/fixedcarrier:TestSemanticIndexComesFromActualCarrierHint`<br>`test:internal/fixedcarrier:TestSemanticIndexDoesNotTreatBinaryNoiseAsMeaning` | `6.0.0-alpha.5` | 2026-09-01 | This is a deterministic carrier contract; it does not claim recovery by an external model. |
| `ORIGAMI.CARRIER.PROFILE_3` | Render and decode fixed-carrier profile 3 without regressing profile 1 or profile 2 decoding. | `implemented` | `test:internal/fixedcarrier:TestProfile3RoundTripFrozenSize`<br>`test:internal/fixedcarrier:TestDecodeAnyAcceptsProfile2AndProfile3`<br>`test:internal/fixedcarrier:TestProfile3DeclaresRequiredZones` | `6.0.0-alpha.13` | 2026-09-01 | Native S2 and E2 model evidence remains pending. |
| `ORIGAMI.CONTEXT_SIMD.SAFE_WIDTH` | Establish Context SIMD safe width from explicit perceptual orthogonality evidence rather than nominal codebook capacity. | `designed` | — | `unknown` | 2026-09-01 | Bundle and checking code exists, but native-model safe-width evidence remains pending. |
| `ORIGAMI.CORE.SELECTIVE_UNFOLD` | Preserve unresolved alternatives during folding and selectively unfold only the declared dependency closure. | `implemented` | `test:internal/core:TestADV006FoldPreservesUnresolvedAlternatives`<br>`test:internal/core:TestADV007SelectiveUnfoldReadsOnlyDependencyClosure`<br>`test:internal/memory:TestManagerUsesBoundedSelectiveUnfolding` | `6.0.0-alpha.7` | 2026-09-01 | Selective unfolding is tested independently from claims about perceptual access to a carrier. |
| `ORIGAMI.CORE.STATE_SEMANTICS` | Execute semantic state transitions with explicit context, declarative rules, distinct statuses, and observation separated from transition authority. | `implemented` | `test:internal/core:TestADV002ContextParticipatesInTransition`<br>`test:internal/core:TestADV003RulesAreDataWithStableProvenance`<br>`test:internal/core:TestADV004SemanticStatusesRemainDistinct`<br>`test:internal/core:TestADV005ObservationDoesNotRewriteTransition` | `6.0.0-alpha.6` | 2026-09-01 | The claim covers deterministic engine behavior, not Native VLM interpretation. |
| `ORIGAMI.DIMENSIONAL.DVR_PSS` | Validate the dimensional visual register and require explicit pair evidence before declaring perceptual orthogonality or safe promoted capacity. | `implemented` | `test:internal/lab/dimensional:TestDefaultRegistryAndNominalCapacity`<br>`test:internal/lab/dimensional:TestValidateDVR`<br>`test:internal/lab/dimensional:TestOrthogonalityRequiresExplicitPass`<br>`test:internal/lab/dimensional:TestPromotedSafeCapacityStillNeedsPairEvidence` | `unknown` | 2026-09-01 | Nominal codebook capacity is not promoted safe capacity. Introduction predates the current VERSION history. |
| `ORIGAMI.DOCUMENT.EXACT_SOURCE_IDENTITY` | Ingest source blocks with hash-bound exact identity while keeping layout geometry separate from semantic authority. | `implemented` | `test:internal/document:TestSourceBlockExactnessIsHashBound`<br>`test:internal/document:TestLayoutJSONPreservesImageOrPDFGeometryWithoutGrantingSemantics`<br>`test:internal/integration:TestEndToEndExactSourceToCanonicalSemanticMemory` | `6.0.0-alpha.7` | 2026-09-01 | The exact plane is independently verifiable and does not grant semantic meaning to layout metadata. |
| `ORIGAMI.EVIDENCE.DETERMINISTIC_REDUCTION` | Reduce addressed evidence deterministically while preserving verified conflicts and refusing unsupported or inexact conclusions. | `implemented` | `test:internal/evidence:TestReducerVerifiesClaimFromAddressedEvidence`<br>`test:internal/evidence:TestReducerIsOrderDeterministic`<br>`test:internal/evidence:TestReducerPreservesVerifiedConflict`<br>`test:internal/evidence:TestReducerReturnsUnknownForUnverifiedEvidence`<br>`test:internal/evidence:TestReducerClaimsExactnessOnlyForByteEqualExactEvidence` | `6.0.0-alpha.6` | 2026-09-01 | Evidence reduction keeps exactness and semantic support as separate gates. |
| `ORIGAMI.GLYPH_CALCULUS.LOCAL_RUNTIME` | Execute the Glyph Calculus pipeline through a tested local runtime and record reproducible reconstruction evidence. | `designed` | — | `unknown` | 2026-09-01 | Historical local artifacts exist, but internal/lab/glyphcalc has no go test and its run folder is not a normalized Run Record. |
| `ORIGAMI.HYBRID.BOUNDED_RUNTIME` | Query, expand, follow, and verify a hybrid carrier through bounded runtime operations without scanning on a missing lookup. | `implemented` | `test:internal/lab/hybridruntime:TestRuntimeVirtualQueryExpandAndSignatureFromCompositePNG`<br>`test:internal/lab/hybridruntime:TestRuntimeLookupFollowAndVerifyFromPNG`<br>`test:internal/lab/hybridruntime:TestRuntimeMissingLookupReturnsUnknownWithoutScan` | `6.0.0-alpha.4` | 2026-09-01 | This claim is limited to the deterministic runtime path. |
| `ORIGAMI.LAB.ARTIFACT_HASHING` | Store experiment artifacts by verified SHA-256 identity and retain their provenance in a normalized run record. | `designed` | — | `unknown` | 2026-09-01 | The package and design exist, but the complete claim lacks an associated go test; see ORIGAMI-DISCREPANCY-003. |
| `ORIGAMI.LAB.EXPERIMENT_IDENTITY` | Derive canonical content-addressed ExperimentID and RunID values from a captured experiment specification and environment. | `designed` | — | `unknown` | 2026-09-01 | Implementation code is present, but no go test demonstrates this contract; see ORIGAMI-DISCREPANCY-003. |
| `ORIGAMI.LAB.FAILURE_TO_REGRESSION` | Convert a captured model failure into a reproducible regression fixture linked to its originating run. | `designed` | — | `unknown` | 2026-09-01 | The R3.10-LAB documentation describes this path, but no go test establishes the end-to-end contract. |
| `ORIGAMI.MEMORY.SCALE_HARNESS` | Run deterministic memory-scale scenarios while bounding metadata loading to selected carriers. | `implemented` | `test:internal/lab/memoryscale:TestRunPreservesAccuracyWhileKnowledgeScales`<br>`test:internal/lab/memoryscale:TestMetadataLoadingIsBoundedBySelectedCarriers`<br>`test:internal/lab/memoryscale:TestConfigRejectsUnsafeOrAmbiguousRanges` | `6.0.0-alpha.4` | 2026-09-01 | Synthetic scale-harness results are implementation evidence, not a real-corpus benchmark. |
| `ORIGAMI.MEMORY.VIRTUAL_CONTEXT` | Provide bounded virtual semantic memory with selective expansion, fidelity downgrade, federation, and token-budget eviction. | `implemented` | `test:internal/memory:TestManagerUsesBoundedSelectiveUnfolding`<br>`test:internal/memory:TestExpandDowngradesFidelityBeforeExceedingBudget`<br>`test:internal/memory:TestFederationRanksCarrierSignaturesBeforeLocalSearch`<br>`test:internal/memory:TestWorkingSetEvictsLeastRecentlyUsedByTokenBudget` | `6.0.0-alpha.4` | 2026-09-01 | The implementation provides bounded access mechanics; it does not prove a remote model can use them natively. |
| `ORIGAMI.MICROISA.SAFE_PROFILE` | Promote a SAFE_MICRO_ISA profile only from measured native-model evidence under the declared protocol. | `designed` | — | `unknown` | 2026-09-01 | The profiler and protocol are present, but native evidence and an associated go test are absent. |
| `ORIGAMI.PERCEPTION.NATIVE_CHANNEL_EXECUTION` | Demonstrate held-out Native VLM interpretation of the declared Origami perceptual channels without undeclared decoder or filesystem access. | `designed` | — | `6.0.0-alpha.2` | 2026-09-01 | Mechanical renderers and evaluators are implemented, but they do not constitute native perceptual evidence. |
| `ORIGAMI.PROMOTION.EVIDENCE_GATES` | Gate promotion on declared evidence and prevent mock, incomplete, or wrong-probe observations from becoming verified semantic claims. | `implemented` | `test:internal/promotion:TestMockCanPassMechanicallyButNeverPromote`<br>`test:internal/promotion:TestRealModelNativeT3TrialRequiresAllT3Fields`<br>`test:internal/promotion:TestWrongProbeProvesImageChallengeFailed`<br>`test:internal/integration:TestPerceptionObservationCannotDirectlyBecomeVerifiedSemanticEvidence` | `6.0.0-alpha.8` | 2026-09-01 | Promotion policy is deterministic; successful real-model campaigns must still be recorded separately. |
| `ORIGAMI.PROTOCOL.CODEC_REGISTRY` | Negotiate declared semantic codecs and perform the required S2 and E2 semantic roundtrip without requiring exact-plane capability. | `implemented` | `test:internal/codec:TestReferenceRegistryValid`<br>`test:internal/codec:TestSemanticCodecDoesNotRequireExactCapability`<br>`test:internal/codec:TestSelectS2WithSemanticReadOnly`<br>`test:internal/codec:TestS2E2RoundTrip` | `6.0.0-alpha.13` | 2026-09-01 | Registry and reference codecs are implemented; cross-model interoperability is not asserted. |
| `ORIGAMI.RECEIVER.REFERENCE_IMPORT` | Validate receiver artifacts, preserve reference-only candidates before promotion, and import Tlaloc candidates without granting physical-symbol authority. | `implemented` | `test:internal/receiver:TestPromotedArtifactRequiresCleanEvidence`<br>`test:internal/receiver:TestReferenceCandidateMayExistBeforePromotionEvidence`<br>`test:internal/receiver:TestImportTlalocCandidateDoesNotBindPhysicalSymbols`<br>`test:internal/receiver:TestImportTlalocCandidateRejectsShortTraceHash` | `6.0.0-alpha.3` | 2026-09-01 | Origami retains authority over representation and promotion. |
| `ORIGAMI.REFERENCE.DETERMINISTIC_ENGINE` | Execute the bounded reference semantics deterministically, enumerate the declared class, and reject contradictions or broken dependency folds. | `implemented` | `test:internal/reference:TestDeterministic`<br>`test:internal/reference:TestContradiction`<br>`test:internal/reference:TestEnumeration`<br>`test:internal/reference:TestDependencyFoldRejectsBrokenPath` | `6.0.0-alpha.3` | 2026-09-01 | Unsupported advanced semantics fail explicitly in the reference phase instead of being inferred. |
| `ORIGAMI.SEMANTIC_SPINE.R1` | Expose canonical semantic graph and trajectory queries without requiring LLM reconstruction and preserve reference relation semantics. | `implemented` | `test:internal/integration:TestCoreSemanticGraphAndTrajectoryAreDirectlyQueryableWithoutLLMReconstruction`<br>`test:internal/integration:TestCoreR1MatchesReferenceR0ForAtomicRequiresRelation`<br>`test:internal/integration:TestCoreR1MatchesReferenceR0ForAtomicExcludesRelation` | `6.0.0-alpha.7` | 2026-09-01 | Direct queryability is established for the Go runtime, not for prompt-only image interpretation. |
| `ORIGAMI.TEMPORAL.AUTOMATON_RUNTIME` | Execute deterministic semantic graph automata and temporal programs with trace generation, checkpoint seek, and replay verification. | `implemented` | `test:internal/automaton:TestPropagation`<br>`test:internal/automaton:TestPriorityThenRuleID`<br>`test:internal/temporal:TestRunSeekReplay`<br>`test:internal/temporal:TestSeekDetectsCorruptDelta`<br>`test:internal/codec:TestTimelineRoundtripReplay` | `6.0.0-alpha.14` | 2026-09-01 | The runtime represents semantic temporal state, not literal video frames. |
| `ORIGAMI.TEMPORAL.CAUSAL_TRACE` | Preserve parallel causal frontiers and reject cyclic or unknown-parent causal traces. | `implemented` | `test:internal/temporal:TestCausalTraceFrontiersPreserveParallelBranches`<br>`test:internal/temporal:TestCausalTraceRejectsCycle`<br>`test:internal/temporal:TestCausalTraceRejectsUnknownParent` | `6.0.0-alpha.15` | 2026-09-01 | This records executable causal trace semantics, not cross-model temporal interpretation. |
| `ORIGAMI.TEMPORAL.SELF_CONTAINED_CARRIER` | Roundtrip an exact temporal program through the fixed self-contained carrier and reject programs that exceed its frozen envelope. | `implemented` | `test:internal/fixedcarrier:TestTemporalCarrierSelfContainedRoundtrip`<br>`test:internal/fixedcarrier:TestTemporalCarrierRejectsOversizeProgram`<br>`test:internal/codec:TestAutomatonRoundtrip` | `6.0.0-alpha.14` | 2026-09-01 | Arbitrary unstructured temporal data is outside the fixed-envelope claim. |
| `ORIGAMI.WRITER.CANONICAL_PLANNING` | Plan semantic visual output deterministically against one canonical profile and reject unsupported or ambiguous intents. | `implemented` | `test:internal/writer:TestPlanUsesOneCanonicalProfile`<br>`test:internal/writer:TestPlanDeterministicAcrossMapOrder`<br>`test:internal/writer:TestUnknownHintAddressRejected`<br>`test:internal/writer:TestArbitraryIntentRejected` | `6.0.0-alpha.9` | 2026-09-01 | Temporal intent currently uses the tested static fallback until its channel is promoted. |
<!-- END GENERATED CLAIMS TABLE -->

Machine-readable capability state: `state/CLAIMS.json`.

Historical component state remains in `state/ORIGAMI_STATE.json`; it is not the authority for capability status.

## Portable compatibility baseline

The Master Prompt R3 assumes only:

```text
Master Prompt
+ explicit user input
+ Origami image/carrier when image input exists
```

It does not assume Tlaloc, Tonal, sandbox, Go/Python, external tools, hidden file access or a binary decoder.

Tool-assisted behavior and prompt-only/native behavior are separate capability claims.

## Semantic visual plane vs exact/control plane

Alpha.12 makes this separation explicit:

```text
SEMANTIC VISUAL PLANE
T0 BOOT
T1 ROSETTA / duplicated visual probe
T2 actual semantic superindex
visible/selective semantic state

EXACT / CONTROL PLANE
T3 deterministic record
CID / root / hash / Merkle / exact verification
```

The exact plane is preserved. It may require deterministic tools/runtime. Its unavailability cannot force a semantic index query to fail mechanically, and it cannot justify invented byte/hash/compression claims.

## Fixed Carrier profile-2

Physical invariants remain unchanged:

```text
640 x 640
8192 PNG bytes exactly
hard maximum 512000 bytes
```

Profile-2 changes semantics, not the envelope:

```text
T0  semantic-first BOOT
T1  ROSETTA + duplicated probe
T2  bounded actual semantic superindex
T3  deterministic control/exact record
VERIFY
```

The integration layer now constructs bounded **valid JSON** for the carrier graph hint and prioritizes actual high-value/top-level canonical labels. The old blind `body[:256]` JSON truncation is removed.

The deterministic decoder still recognizes profile-1 records for backward compatibility.

## Native query routing

```text
identity          -> T0/T1
index             -> T2
overview          -> T2 + visible semantic state
locate topic      -> T2 -> selective semantic expand when available
explain topic     -> T2 -> smallest relevant semantic region
exact quote/hash  -> exact plane / declared tool when available
```

If only a bounded T2 top-level index is visible, the answer is `PARTIAL`, not a fabricated full table of contents.

If T2 cannot support the answer, return `UNKNOWN`/`NOT_VERIFIED`; do not decode T3 pixels speculatively.

## Prohibited unverified mechanical claims

Without an actually executed declared exact path, a receiver must not claim:

```text
payload byte count
header size
compression algorithm
compressed residual size
hidden prompt byte length
exact hash
archive contents
exact page ranges allegedly recovered from binary
```

This prohibition is now in the Master Prompt R3 and Native Semantic Navigation contract.

## WRITE path

```text
source
 -> Semantic IR
 -> visual intents
 -> canonical profile
 -> ROSETTA
 -> PROGRAM / T2 / MEMORY / VERIFICATION
 -> deterministic compiler when available
 -> carrier
 -> roundtrip verification
```

Writer should deliberately create a useful T2 because common semantic questions must not require whole-corpus or exact-plane decoding.

Without a compiler, WRITE remains `CONSTRUCTION_SPEC_ONLY`.

## Canonical visual rule

Origami has one canonical functional aesthetic per visual-profile version. ROSETTA is mandatory and declares profile/version, active dimensions, visual primitive -> semantic role, approved carrier parameters and reveal procedures for promoted non-static channels.

Experimental dimensions remain COLOR, NUMERIC_STRUCTURE, INTERFERENCE, DEPTH, TEMPORAL and EMERGENT. Candidate operations include MOIRE, PHASE_SHIFT, STEREO_BIND, PARALLAX_RESOLVE, KINETIC_REVEAL, TEMPORAL_INTEGRATE and TEMPORAL_DECAY.

A promoted non-static channel must declare both meaning and reveal procedure. Failure to reveal returns `UNKNOWN`, not `ABSENT`.

## Semantic core and memory

The semantic core remains:

```text
S_(t+1) = F(S_t, C_t, R)
```

with distinct PRESENT, ABSENT, UNKNOWN, INHIBITED and CANCELLED states.

Semantic Fold preserves dependencies/alternatives. Selective Unfold opens only declared closure. Virtual Memory keeps the active model-facing window around 4000 token-equivalent while total external memory may be much larger.

Evidence Reduction accepts external proposals but Origami itself verifies evidence identity/fidelity and reduces to VERIFIED, REJECTED, CONFLICT or UNKNOWN. `VERIFIED_EXACT` requires byte-equal accepted exact evidence.

## Development lifecycle

A development tool may use the failed trial as an experiment/regression, search prompt/profile candidates and produce evidence. Origami decides whether a candidate becomes canonical.

```text
Origami profile/version N
     ↓
external development experiment
     ↓
candidate + evidence
     ↓
Origami validation
     ↓
Origami decides N+1
     ↓
optional Tonal composition/pin
```

Tlaloc alpha.15 adds corresponding Native index/semantic/mechanical-dependency metrics, but that remains external development machinery rather than Origami runtime authority.

## Hard invariants

```text
ORIGAMI_OWNS_ORIGAMI_RELEASES
MASTER_PROMPT_IS_PORTABLE_BASELINE
SEMANTIC_NAVIGATION != MECHANICAL_DECODE
INDEX_QUERY_ROUTES_TO_T2
T2_CONTAINS_ACTUAL_SEMANTIC_ENTRIES
EXACT_PLANE_OPTIONAL_FOR_SEMANTIC_NAVIGATION
NO_UNVERIFIED_BYTE_HASH_COMPRESSION_CLAIMS
FAILED_NATIVE_TRIAL_BECOMES_REGRESSION
TOOL_ASSISTED_SUCCESS != PROMPT_ONLY_SUCCESS
FALSE_EXACT = 0
UNKNOWN > invented exactness
ONE_CANONICAL_AESTHETIC_PER_PROFILE_VERSION
ROSETTA_ALWAYS_PRESENT
PROMOTED_PERCEPTUAL_CHANNEL_REQUIRES_REVEAL_PROCEDURE
FAILED_REVEAL != SEMANTIC_ABSENCE
SOURCE_SCREENSHOT != ORIGAMI
CONSTRUCTION_SPEC != COMPILED_VERIFIED_CARRIER
ROUNDTRIP_REQUIRED
source plane != semantic plane
literal transport != semantic Fold
observation != transition
address != CID
active model interface != total memory
no implicit global exact scan
perception != resolution != execution != verification
Fixed Carrier != corpus
```

## Evidence boundary

Alpha.12 fixes the architecture and deterministic regression harness. It does **not** establish that a real VLM now succeeds. Native promotion remains blocked until held-out models can recover T2/index and other semantic answers without undeclared mechanical decoding.
