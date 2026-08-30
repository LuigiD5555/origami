# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.11  
**Status:** experimental

## Project role

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory system.

Origami owns its own semantics, canonical visual grammar, ROSETTA, Writer/Reader contracts, Master Prompt and version/profile releases.

External systems may help develop it:

```text
Tlaloc
Blueprint Framework
future development tools
```

Tonal may optionally compose/pin several of those tools and exact target revisions. None is required for Origami's portable baseline.

## Portable compatibility baseline — alpha.11

The universal Master Prompt R2 is now explicitly designed for a compatible target model that may have only:

```text
Master Prompt
+ explicit user input
+ Origami image/carrier when image input exists
```

It does not assume:

```text
Tlaloc
Tonal
sandbox
Go / Python
external tools
hidden file access
```

When richer capabilities exist, they are explicit enhancements:

```text
prompt + declarative context/IR
prompt + Origami tools
prompt + ingestion adapter
prompt + compiler/runtime
specialized model
```

Tool-assisted behavior and prompt-only/native behavior are separate capability claims.

## Current capability state

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
| perceptual profile evolution / reveal semantics | EXPERIMENTAL_CONTRACT_EXTENSION — alpha.10 |
| Development Tool Boundary R0 | EXPERIMENTAL_REFERENCE_CONTRACT — alpha.11 |
| Master Prompt | R2 PORTABLE READ/WRITE REFERENCE_CANDIDATE |
| generic Writer plan -> fully profile-generic renderer | PARTIAL |

Machine-readable state: `state/ORIGAMI_STATE.json`.

## READ without tools

Prompt-only/native READ can bootstrap BOOT/ROSETTA/profile information and interpret visual structure that is genuinely available to the model.

If an exact lookup, hash check, query operation or deterministic execution is required but no such capability is available, the correct result is:

```text
UNKNOWN
```

or, when semantic support exists but requested verification cannot be executed:

```text
NOT_VERIFIED
```

Tool absence is a capability boundary, not evidence of semantic absence.

## WRITE without compiler

The prompt may still produce:

```text
Semantic IR
Visual Intent Plan
ROSETTA
PROGRAM / INDEX / MEMORY / VERIFICATION specification
```

If no deterministic compiler or equivalent verifiable image-generation path exists, it must not claim a compiled carrier.

Correct status:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

## Canonical visual rule

Origami has one canonical functional aesthetic per visual-profile version.

ROSETTA is mandatory and declares:

```text
profile ID/version
active dimensions
visual primitive -> semantic role
approved carrier parameters
reveal procedure for promoted non-static perceptual channels
```

ROSETTA does not license a private aesthetic per PDF/model.

## Perceptual channels

Experimental dimensions include:

```text
COLOR
NUMERIC_STRUCTURE
INTERFERENCE
DEPTH
TEMPORAL
EMERGENT
```

Candidate operations include:

```text
MOIRE
PHASE_SHIFT
STEREO_BIND
PARALLAX_RESOLVE
KINETIC_REVEAL
TEMPORAL_INTEGRATE
TEMPORAL_DECAY
```

A promoted non-static channel must declare both its meaning and its reveal procedure. Failure to satisfy the reveal condition returns `UNKNOWN`, not `ABSENT`.

## Writer R0

```text
READ
carrier -> BOOT -> ROSETTA -> semantic interpretation

WRITE
source
 -> Semantic IR
 -> visual intents
 -> canonical profile
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> deterministic compiler when available
 -> carrier
 -> roundtrip verification
```

The model may construct semantic/visual intent. A verified compiler/rendering path remains carrier authority for `COMPILED_VERIFIED` status.

## Development lifecycle

Tlaloc is currently the strongest behavior-experiment kit around Origami, but it is external.

It may use bounded Tlaloque swarms to discover a working READ/WRITE/representation procedure, distill that behavior into a prompt candidate and test it on clean target models.

Other tools may contribute different kinds of development evidence.

Lifecycle:

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

## Semantic core and memory

The semantic core remains:

```text
S_(t+1) = F(S_t, C_t, R)
```

with distinct:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

Semantic Fold preserves dependencies/alternatives. Selective Unfold opens only declared closure.

Virtual Memory keeps the active model-facing window around 4000 token-equivalent while total external memory may be much larger.

Evidence Reduction accepts external proposals but Origami itself verifies evidence identity/fidelity and reduces to:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

`VERIFIED_EXACT` requires byte-equal accepted exact evidence.

## Fixed Carrier and optional tools

Current reference Fixed Carrier:

```text
640 x 640
8192 PNG bytes exactly
hard current-profile maximum 512000 bytes
```

with T0/T1/T2/T3/VERIFY bootstrap.

The carrier is not the corpus.

`tlaloc.origami-tools.r2` remains a useful existing optional tool ABI. It does not make Tlaloc mandatory for Origami or for the Master Prompt.

## Perception Promotion

Hybrid and Native-T3 claims remain separate. MOCK validates mechanics only.

A real campaign may be orchestrated by Tlaloc or another development system. Origami owns the interpretation/promotion of Origami capability/profile evidence.

## Profile-evolution objective

External experiments may search:

```text
prompt
geometry / primitives
layout
redundancy
color
numeric / prime / modular structure
moire / phase / interference
stereo / parallax / depth
temporal / motion-bound structure
emergent percepts
```

A newer profile should improve measurable representation while preserving hard semantics.

Typical objectives:

```text
maximize:
  recoverable semantic capacity
  semantic roundtrip
  readability/reveal reliability
  routing/evidence accuracy
  transport robustness

minimize:
  carrier bytes
  recognition latency
  BOOT steps
  decode/unfold steps
  model-facing context
```

## Hard invariants

```text
ORIGAMI_OWNS_ORIGAMI_RELEASES
MASTER_PROMPT_IS_PORTABLE_BASELINE
TLALOC_IS_OPTIONAL
TONAL_IS_OPTIONAL
DEVELOPMENT_TOOL != RUNTIME_REQUIREMENT
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
model specialization is optional
```

## Evidence boundary

Alpha.11 defines the development-tool boundary and portable prompt baseline. It does **not** claim that prompt-only operation already matches Hybrid/tool operation on every task, nor that every model can natively decode/write every Origami channel.
