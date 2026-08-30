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

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels / TLI | EXPERIMENTAL_CONTRACT — alpha.2 |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| Virtual Memory R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.4 |
| Memory Scale Lab R0 | EXPERIMENTAL_DETERMINISTIC_LAB |
| Fixed Carrier R2 profile-2 | EXPERIMENTAL_REFERENCE_IMPLEMENTED / NATIVE SEMANTIC EVIDENCE PENDING — alpha.12 |
| Evidence Reduction R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.6 |
| Semantic Spine R1 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.7 |
| Perception Promotion R1 | EXPERIMENTAL_REFERENCE_IMPLEMENTED / REAL MODEL EVIDENCE PENDING — alpha.8 |
| Canonical Visual Grammar R0 | EXPERIMENTAL_REFERENCE_PROFILE — alpha.9 |
| Writer R0 | EXPERIMENTAL_REFERENCE_IMPLEMENTED — alpha.9 |
| perceptual profile evolution / reveal semantics | EXPERIMENTAL_CONTRACT_EXTENSION — alpha.10 |
| Development Tool Boundary R0 | EXPERIMENTAL_REFERENCE_CONTRACT — alpha.11 |
| Native Semantic Navigation R0 | **EXPERIMENTAL_REFERENCE_IMPLEMENTED / HELD-OUT VLM EVIDENCE PENDING — alpha.12** |
| Master Prompt | **R3 SEMANTIC-FIRST READ/WRITE REFERENCE_CANDIDATE** |
| generic Writer plan -> fully profile-generic renderer | PARTIAL |

Machine-readable state: `state/ORIGAMI_STATE.json`.

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
