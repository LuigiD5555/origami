# Current Origami state

**Status date:** 2026-08-30  
**Project version:** 6.0.0-alpha.10  
**Status:** experimental

## Project role

Origami is the visual/computational representation, state-machine language and model-agnostic virtual-memory system.

The surrounding tools are separate:

```text
                           TONAL
          optional multi-tool composition layer
        /                 |                    \
   TLALOC          Blueprint Framework      other tools
 development kit      development tool      future tools
        |
        | may improve / test / build
        v
                           ORIGAMI
          owns its own canonical versions/profiles
```

Tlaloc is the primary current development kit for Origami, but Origami does not depend conceptually on Tlaloc. Tlaloc may also be used to develop other tools. Tonal exists only when several independent tools/components need a reproducible composition; Tonal does not promote Origami semantics or visual profiles.

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
| Master Prompt | READ/WRITE REFERENCE_CANDIDATE |
| generic Writer plan -> fully profile-generic renderer | PARTIAL |

The machine-readable source is `state/ORIGAMI_STATE.json`.

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

## Perceptual depth is semantic, not decoration

Origami's alpha.2 perceptual contract already defines interference, depth, temporal and emergent channels. Alpha.10 connects those ideas to Writer/profile evolution.

Experimental dimensions now include:

```text
D10 COLOR
D11 NUMERIC_STRUCTURE
D12 TEMPORAL
D13 INTERFERENCE
D14 DEPTH
D15 EMERGENT
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

If a future Origami profile promotes one of these channels, ROSETTA must declare both:

```text
what it means
+
how to reveal/observe it
```

For example, a moiré relation may require two patterns and a declared relative alignment; a stereo channel may require compatible views; a temporal latent percept may require a trajectory and integration window.

If the receiver cannot satisfy the reveal condition:

```text
UNKNOWN
```

not `ABSENT` and not an invented interpretation.

## Writer R0

Origami supports formal READ and WRITE directions.

```text
READ
carrier -> BOOT -> ROSETTA -> semantic interpretation

WRITE
source / Document IR
 -> Semantic IR
 -> visual intents
 -> canonical profile
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> deterministic compiler
 -> carrier
 -> roundtrip verification
```

`internal/writer` produces a deterministic construction plan. Alpha.10 extends each ROSETTA entry so promoted non-static visual channels can include a `RevealProcedure` containing operation, required inputs, observation condition, trajectory/phase, integration window and failure state.

The model may produce Semantic IR or construction intent; the compiler remains carrier authority.

## Visual-profile evolution

A new Origami version should not win because it is prettier. It should demonstrably improve useful representation.

Tlaloc or another development tool may search candidates over:

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
emergent multi-layer percepts
```

The optimization direction is:

```text
maximize:
  recoverable semantic capacity per byte
  semantic roundtrip
  readability / perceptual reveal
  routing accuracy
  verified-evidence accuracy
  transport robustness

minimize:
  carrier bytes
  recognition latency
  BOOT steps
  decode/unfold steps
  model-facing context
```

Hard semantic/evidence invariants cannot be traded away for score.

The lifecycle is:

```text
Origami profile N
     ↓
development-tool search (Tlaloc today; alternatives allowed)
     ↓
candidate + measured evidence
     ↓
Origami validation
     ↓
Origami promotes profile N+1
     ↓
optional Tonal composition with chosen development toolchain
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

Semantic Fold preserves dependencies/alternatives and Selective Unfold opens only declared closure.

Virtual Memory keeps the active model-facing window around 4000 token-equivalent while total external memory may be much larger.

Evidence Reduction accepts external candidate claims but Origami itself verifies addresses/CIDs/source hashes/fidelity and reduces to:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

`VERIFIED_EXACT` requires byte-equal accepted exact evidence.

## Fixed Carrier and perception

The current reference Fixed Carrier remains:

```text
640 x 640
8192 PNG bytes exactly
hard maximum 512000 bytes
```

with T0/T1/T2/T3/VERIFY bootstrap.

It is the current concrete R0 aesthetic instance, not a permanent ban on future promoted visual profiles.

Perception Promotion still keeps Hybrid and Native-T3 claims separate. A mock can validate mechanics but never empirical model capability.

Cross-model campaign orchestration may be performed by Tlaloc or another development tool. Origami owns its own capability/profile authority. Tonal may record a reproducible composition but is not Origami's semantic authority.

## Hard invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
ONE_CANONICAL_AESTHETIC_PER_PROFILE_VERSION
ROSETTA_ALWAYS_PRESENT
ROSETTA != PER_DOCUMENT_AESTHETIC_DRIFT
PROMOTED_PERCEPTUAL_CHANNEL_REQUIRES_REVEAL_PROCEDURE
FAILED_REVEAL != SEMANTIC_ABSENCE
SOURCE_SCREENSHOT != ORIGAMI
SEMANTIC_IR_PRECEDES_VISUAL_COMPILATION
CONSTRUCTION_SPEC_PRECEDES_COMPILATION
ROUNDTRIP_REQUIRED
DEVELOPMENT_TOOL_PROPOSES / ORIGAMI PROMOTES
TONAL_COMPOSITION != ORIGAMI_PROFILE_PROMOTION
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

Alpha.10 does not claim that moiré, stereo, temporal, color, prime-derived or other advanced candidates already outperform the current profile. It defines how such a candidate can be represented, revealed, measured and promoted without confusing visual novelty with demonstrated capability.
