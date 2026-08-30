# Origami project boundary

Origami is the thing being built and versioned. Tlaloc, Blueprint Framework and Tonal are external development/composition systems and must not become hidden requirements of the language.

```text
                 optional development ecosystem

        Tlaloc       Blueprint Framework       future tools
           \               |                    /
            \------ experiments/evidence ------/
                         |
                         v
                      ORIGAMI
       representation language + machine + memory
       owns its semantics, Master Prompt, visual grammar,
       ROSETTA, Writer/Reader contracts and releases
                         |
                         v
              standalone portable use

Optional:
Tonal may pin/compose the exact development tools + target revisions.
```

## Stable hierarchy

**Origami** owns Origami.

**Tlaloc** is a behavioral development kit. It can decompose a requested behavior into bounded Tlaloque steps, prove a reference swarm can perform it, distill that behavior into prompt candidates and test candidate improvements. Origami is one possible Tlaloc target, not Tlaloc's definition.

**Blueprint Framework** or future systems may contribute other development methods.

**Tonal** is an optional composition/reproducibility layer for combining those independent tools. Tonal does not decide what Origami means or which Origami profile is canonical.

## Portable Origami baseline

Origami's universal Master Prompt is designed to remain useful when a compatible model has only:

```text
Master Prompt
+ explicit user input
+ Origami image/carrier when image input is available
```

The baseline must not require:

```text
Tlaloc
Tonal
sandbox
Go
Python
external tools
hidden file access
```

Richer capabilities are explicit enhancements:

```text
+ declarative context / Semantic IR
+ Origami tools
+ ingestion adapters
+ compiler/runtime
+ specialized model
```

Tool-assisted success must not be presented as proof that prompt-only/native operation succeeded.

## Canonical visual ownership

Origami owns **one canonical functional visual aesthetic per promoted profile version**.

External development tools may experimentally propose changes to:

```text
Master Prompt
shape / primitive roles
layout
channel assignments
redundancy
color
numeric / mathematical structure
moire / interference / phase
stereo / parallax / depth
temporal / motion-bound structure
emergent multi-layer or multi-instant percepts
```

Those remain candidates until Origami adopts them.

Lifecycle:

```text
Origami profile N
      |
      v
external development experiment
      |
      v
candidate + evidence
      |
      v
Origami semantic/visual validation
      |
      v
Origami decides profile N+1
      |
      +--> standalone use
      +--> optional Tonal toolchain composition/pin
```

## ROSETTA ownership

ROSETTA is Origami-owned and mandatory.

It declares:

```text
profile ID/version
active dimensions
visual primitive -> semantic role mappings
approved parameters/extensions
reveal procedures for promoted non-static channels
```

For interference, depth, temporal or emergent channels, ROSETTA must describe both **what the channel means** and **how it is revealed**. If the reveal condition cannot be established, the receiver returns `UNKNOWN`; inability to reveal is not semantic absence.

ROSETTA does not authorize per-document private dialects.

## Origami owns

- `S_(t+1)=F(S_t,C_t,R)` state/context/rule semantics;
- PRESENT / ABSENT / UNKNOWN / INHIBITED / CANCELLED distinctions;
- contradictions, cancellation and higher-order interactions;
- Observation Contracts and temporal observations;
- semantic Fold/Unfold and unresolved alternatives;
- Virtual Memory/addressability/ContextPacket contracts;
- source-plane vs semantic-plane separation;
- evidence/CID/source-hash/fidelity authority;
- `VERIFIED / REJECTED / CONFLICT / UNKNOWN` claim semantics;
- `FALSE_EXACT=0` and exactness boundaries;
- BOOT/ROSETTA/PROGRAM/INDEX/MEMORY/VERIFICATION semantics;
- universal Master Prompt READ/WRITE behavior;
- canonical visual profile and versioned aesthetic;
- promoted perceptual-channel semantics and reveal procedures;
- Writer construction-plan and roundtrip semantics;
- Fixed Carrier/control-plane semantics;
- Origami version/profile promotion and artifact provenance.

## External development tools may provide

Tlaloc or another tool may provide, when useful:

- swarm/reference experiments;
- behavior/prompt distillation;
- prompt mutation/search;
- document ingestion/OCR adapters;
- representation/profile candidate search;
- model campaigns;
- degradation/transport experiments;
- perception/reveal benchmarks;
- regression suites;
- recommendations with evidence.

These are development services. They do not redefine Origami and must not silently become requirements of the portable Master Prompt.

## READ boundary

```text
portable:
Master Prompt + carrier/image + question
 -> BOOT
 -> ROSETTA
 -> semantic interpretation
 -> UNKNOWN/NOT_VERIFIED where unavailable operations are required
```

Optional tools may provide bounded lookup/verify/execute acceleration.

## WRITE boundary

```text
source
 -> Semantic IR
 -> visual intent
 -> canonical profile
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> compiler when available
 -> roundtrip verify
```

Without a compiler/verified image-generation path, a model may produce a construction specification but must report:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

A description of pixels is not a compiled verified Origami.

## Tonal boundary

Tonal may represent a reproducible development environment such as:

```text
Tonal
├── Tlaloc @ exact SHA
├── Blueprint Framework @ exact SHA
└── Origami @ exact SHA
```

That composition records provenance/compatibility. It does not promote Origami semantics or profiles.

## Anti-shortcut rules

```text
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PER-DOCUMENT PRIVATE STYLE != CANONICAL ORIGAMI PROFILE
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
DEVELOPMENT TOOL != ORIGAMI RUNTIME REQUIREMENT
DEVELOPMENT TOOL != ORIGAMI SEMANTIC AUTHORITY
TONAL COMPOSITION != ORIGAMI PROFILE PROMOTION
TOOL-ASSISTED SUCCESS != PROMPT-ONLY SUCCESS
FALSE_EXACT = 0
```
