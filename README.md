# Origami 6.0.0-alpha.15

Origami is an experimental **self-describing visual/computational communication protocol, state-machine language and model-agnostic virtual memory** for structured state, relations, knowledge and selectively accessible memory.

Origami is independently usable. Tlaloc, Blueprint Framework, Tonal and other development systems may improve, test or compose it, but they are not hidden runtime requirements.

Current maintained surface: [`docs/ACTIVE_SURFACE.md`](docs/ACTIVE_SURFACE.md).

## What alpha.15 means now

Alpha.13 established the self-describing read/write protocol and codec registry. Alpha.14 added executable temporal graph automata and a self-contained temporal carrier. Alpha.15 added the Origami-owned experimental candidate compiler and the current causal-trace foundations used for Tlaloc interoperability. The active direction is now also testing whether tiny VLMs can execute **visible precompiled microinstructions** rather than inventing a plan at runtime; see [`docs/BEHAVIOR_IN_IMAGE_R0.md`](docs/BEHAVIOR_IN_IMAGE_R0.md).

The durable protocol principle remains:

> Origami carries not only the grammar needed to understand a message, but also declared procedures for reading and writing compatible Origami messages.

The core protocol path is now:

```text
                         ORIGAMI PROTOCOL
                                |
                               T0
                              BOOT
                                |
                               T1
                             ROSETTA
                                |
              +-----------------+-----------------+
              |                                   |
       DECODER REGISTRY                    ENCODER REGISTRY
          S0 ... S6                           E0 ... E6
              |                                   |
             READ                                WRITE
              |                                   |
              v                                   v
       T2 / semantic state                 Construction IR
              |                                   |
              +-----------------+-----------------+
                                |
                         EXACT / CONTROL
                            X* / Q*
                                |
                               T3
                                |
                             VERIFY
```

The reference equation is:

```text
DECODE(ENCODE(S)) ~= S
```

The first implemented semantic pair is the superindex path:

```text
E2 ENCODE_SUPERINDEX
S2 READ_SUPERINDEX
S2(E2(INDEX)) ~= INDEX
```

See `docs/ORIGAMI_PROTOCOL_R0.md`, `docs/CODEC_REGISTRY_R0.md` and `docs/CAPABILITY_NEGOTIATION_R0.md`.

## The decoder rule

Alpha.13 does **not** define “decoder = bad”.

The actual rule is:

```text
SELF_DECLARED_SEMANTIC_CODEC = ALLOWED
UNDECLARED_EXTERNAL_CODEC_DEPENDENCY = FORBIDDEN
SEMANTIC_NAVIGATION_MUST_NOT_REQUIRE_EXACT_CODEC
```

For example, the question:

```text
What is the index?
```

routes through:

```text
T0 -> T1 -> discover S2 -> T2 -> answer
```

It must not require:

```text
extract every pixel
 -> bits
 -> bytes
 -> decompression
 -> reconstruct whole source
 -> finally find the index
```

Exact/control codecs remain available when the query actually requires exactness.

## The encoder rule

A compatible model can also discover how to express new information in Origami.

For a semantic index:

```text
Semantic IR
 -> discover E2
 -> T2 Construction IR
 -> deterministic compiler when available
 -> Origami carrier
 -> S2 roundtrip verification
```

A text-only model without an image compiler may stop honestly at:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

It must not claim that it produced a verified carrier when it did not actually compile one.

## Master Prompt R4

`generated/MASTER_PROMPT.md` is **Universal Read / Write Bootstrap R4 — protocol/codec aware**.

SHA-256:

```text
78cbcee8d0f71e1b06e2902c9b9600779e3cd0148251ad9d50623d07969c2a56
```

R4 is intentionally a small universal handshake. Its role is to teach a model to:

```text
1. locate BOOT;
2. read ROSETTA;
3. discover declared codecs;
4. inspect available capabilities;
5. classify the request;
6. select the smallest sufficient codec;
7. prefer semantic codecs before exact codecs for semantic work;
8. READ or WRITE;
9. return UNKNOWN / NOT_VERIFIED when the required capability is unavailable.
```

The concrete carrier/profile supplies its own bindings and codec entry points. The Master Prompt does not need to hard-code every future visual channel.

R4 remains `REFERENCE_CANDIDATE` until held-out real-model evidence supports promotion.

## Capability negotiation

Different models may support different portions of Origami.

A receiver can conceptually expose capabilities such as:

```text
S0 identity       supported
S1 hierarchy      supported
S2 superindex     supported
E0 identity       supported
E1 hierarchy      supported
E2 superindex     supported
X4 residual       unavailable
MOIRE             unavailable
```

Failure behavior is explicit:

```text
unsupported semantic operation -> UNKNOWN
unsupported exact operation    -> NOT_VERIFIED
compiler unavailable            -> CONSTRUCTION_SPEC_ONLY
```

A model is not required to pretend it can execute every advanced channel.

## Fixed Carrier R2

The current default Fixed Carrier remains profile-2. Profile-1 remains decodable for backward compatibility.

```text
profile-1  legacy deterministic decode
profile-2  current default semantic-first renderer
profile-3  protocol-aware experimental renderer
```

The physical envelope remains:

```text
640 x 640
8192 PNG bytes exactly
hard ceiling: 512000 bytes
```

The carrier is not the corpus. Corpus growth occurs in canonical/addressable memory rather than by growing the PNG.

### Profile-3

`origami.fixed-carrier.r2.profile-3` is now a real parallel renderer rather than only a construction specification.

It visibly carries:

```text
T0  protocol BOOT
T1  ROSETTA + S*/E* codec entry points + capability hints
T2  actual bounded semantic superindex
PROGRAM  semantic codec procedure hints
T3 / EXACT  X*/Q* control/exact plane
VERIFY  S2(E2(INDEX)) ~= INDEX / FALSE_EXACT=0
```

Deterministic CI currently proves:

```text
profile-3 render/decode roundtrip       PASS
640 x 640                               PASS
8192-byte frozen PNG envelope           PASS
profile-1/profile-2 decode regression   PASS
```

It does **not** yet prove:

```text
held-out VLM Native S2 index recovery
held-out VLM E2 write behavior
cross-model A -> B -> C interoperability
```

Therefore profile-3 status is:

```text
EXPERIMENTAL_RENDERER_IMPLEMENTED_EVIDENCE_PENDING
```

and profile-2 remains the default renderer.

See `docs/FIXED_CARRIER_PROFILE_3_R0.md` and `spec/FIXED_CARRIER_PROFILE_3_R0.json`.

## Native semantic query routing

Reference routes are now codec-aware:

```text
What is this?             -> T0/T1 -> S0
What is the index?        -> T0/T1 -> S2 -> T2
What is it about?         -> T2 + visible semantic structure
Where is topic X?         -> S3 / T2 -> selective semantic expand
Explain topic X           -> smallest sufficient semantic region
Quote/hash exact X        -> declared X* exact/control path when available
```

`T2` must contain **actual semantic entry points**, not merely generic PAGE/GRAPH/SOURCE categories.

## The alpha.12 failure remains a regression

The failure is preserved at:

`experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json`.

It established:

```text
BOOT perceived                    yes
actual semantic index recovered   no
external mechanical path sought   yes
unverified exact claims emitted   yes
```

Alpha.13 refines the lesson: the problem was not the existence of a decoder. The model was being driven toward the **wrong decoder family**. Semantic navigation now has self-declared semantic codecs; exact mechanics remain separate.

## READ / WRITE portability

Minimum compatibility assumption remains intentionally small:

```text
Master Prompt
+ explicit user input
+ Origami carrier when image input exists
```

No Tlaloc, Tonal, sandbox, Go/Python, filesystem or hidden runtime is assumed.

A model may use an explicitly available compiler/tool/runtime, but tool-assisted success must not be reported as proof of prompt-only or Native success.

## Canonical architecture

```text
PDF / IMAGE / TEXT / CONVERSATION / STATE
                    |
                    v
       declared source/ingestion adapter
                    |
                    v
              SEMANTIC IR
                    |
       +------------+-------------+
       |                          |
       v                          v
SEMANTIC SPINE R1             WRITER R0
Fold / Unfold                  E* encoders
       |                          |
       |                          v
       |             CANONICAL VISUAL GRAMMAR R0
       |                          |
       |                          v
       |           ROSETTA + CODEC REGISTRY + T2
       |                + MEMORY + VERIFICATION
       |                          |
       |                          v
       |                 deterministic compiler
       |                    when available
       +------------+-------------+
                    |
                    v
               VIRTUAL MEMORY
                    |
                    v
          bounded ContextPacket
                    |
                    v
            compatible model
```

External development systems can construct experiments around this path without becoming semantic authority.

## Semantic Spine R1

The deterministic semantic core uses:

```text
S_(t+1) = F(S_t, C_t, R)
```

and preserves distinct `PRESENT`, `ABSENT`, `UNKNOWN`, `INHIBITED`, `CANCELLED` states. Observation remains separate from transition. Fold preserves unresolved alternatives; selective Unfold records touched addresses.

## Canonical Visual Grammar R0

Origami keeps one canonical functional aesthetic per profile version. Current canonical dimensions include geometry, fill/contrast, position, topology, enclosure, scale, repetition, density and limited declarative text.

Experimental candidates include color, numeric structure, interference/moiré/phase, stereo/parallax/depth, temporal/motion-bound structure and emergent multi-layer/multi-instant percepts.

ROSETTA is always present. In alpha.13 it additionally declares or binds:

```text
semantic roles
active dimensions
reveal procedures
Decoder Registry
Encoder Registry
capability/fallback semantics
```

Advanced channels are not promoted merely because they are representable.

## Virtual Memory R0

Origami memory can be much larger than active model context. The model-facing working budget remains approximately 4000 token-equivalent; this is an active interface budget, not total storage capacity.

```text
large / multi-carrier memory
 -> GraphSignature routing
 -> local metadata graph
 -> selective fidelity unfold
 -> ContextPacket
```

Fidelity order:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Address is location; CID is content identity. Deep exact payload is reopened selectively rather than silently scanned globally.

## Evidence Reduction R0

External systems may propose `SUPPORT`, `OPPOSE`, `UNKNOWN`. Origami independently resolves evidence addresses/CIDs/source hashes/fidelity and reduces deterministically to `VERIFIED`, `REJECTED`, `CONFLICT`, `UNKNOWN`.

Agent/model confidence is never evidence authority. `VERIFIED_EXACT` requires byte-equal accepted exact evidence.

## Ecosystem boundary

```text
Tlaloc             Blueprint Framework           future tools
 behavioral dev        structural dev                 ...
      \                     |                       /
       \------ candidates / experiments / evidence -/
                              |
                              v
                           ORIGAMI
       owns protocol semantics, ROSETTA, codecs, profiles, releases
                              |
                              v
                    standalone portable use

Optional: Tonal may pin/compose exact revisions.
```

Tlaloc can search/evaluate candidate behavior and cross-model interoperability. Origami decides what is canonical. Tonal records composition/provenance; it does not promote model capability.

## Hard invariants

```text
ORIGAMI OWNS ORIGAMI RELEASES
PROTOCOL != CARRIER
MASTER PROMPT IS PORTABLE BASELINE
ROSETTA ALWAYS PRESENT
ROSETTA DECLARES DECODING SEMANTICS
ROSETTA DECLARES ENCODING SEMANTICS
SELF-DECLARED SEMANTIC CODEC IS ALLOWED
NO UNDECLARED EXTERNAL DECODER DEPENDENCY
SEMANTIC NAVIGATION MUST NOT REQUIRE EXACT CODEC
INDEX QUERY ROUTES THROUGH S2/T2
T2 CONTAINS ACTUAL SEMANTIC ENTRIES
EXACT PLANE OPTIONAL FOR SEMANTIC NAVIGATION
NO UNVERIFIED BYTE / HASH / COMPRESSION CLAIMS
FAILED REAL TRIAL -> REGRESSION
TOOL-ASSISTED SUCCESS != PROMPT-ONLY SUCCESS
ONE CANONICAL AESTHETIC PER PROFILE VERSION
FAILED REVEAL != SEMANTIC ABSENCE
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
ACTIVE MODEL INTERFACE != TOTAL MEMORY
NO IMPLICIT GLOBAL EXACT SCAN
FALSE_EXACT = 0
UNKNOWN > INVENTED EXACTNESS
```

## Evidence still pending

Alpha.13 deliberately does **not** claim universal LLM interoperability. Still required:

- repeat `What is the index?` with profile-3 + R4 on held-out real multimodal models;
- measure Native `S2` recovery over multiple models/trials;
- test `E2` write/construction behavior on clean models;
- test A -> B -> C semantic preservation and drift;
- verify transport degradation of T1/T2 readability;
- continue Hybrid/exact-plane tests separately;
- only then consider profile-3/default-profile promotion.

## Useful commands

```bash
# semantic codec registry / S2-E2 roundtrip
go run ./cmd/origami-codec -mode registry -out -

go run ./cmd/origami-codec -mode roundtrip-index -in index.json -out -

# experimental profile-3 carrier
go run ./cmd/origami-profile3-carrier -mode build -in metadata.json -out origami-profile3.png

go run ./cmd/origami-profile3-carrier -mode decode -in origami-profile3.png
```

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
docs/ACTIVE_SURFACE.md
PROJECT_BOUNDARY.md
generated/MASTER_PROMPT.md
docs/ORIGAMI_PROTOCOL_R0.md
docs/CODEC_REGISTRY_R0.md
docs/CAPABILITY_NEGOTIATION_R0.md
docs/FIXED_CARRIER_PROFILE_3_R0.md
docs/NATIVE_SEMANTIC_NAV_R0.md
docs/CANONICAL_VISUAL_GRAMMAR_R0.md
docs/WRITER_R0.md
docs/BEHAVIOR_IN_IMAGE_R0.md
spec/ORIGAMI_PROTOCOL_R0.json
spec/CODEC_REGISTRY_R0.json
spec/CAPABILITY_NEGOTIATION_R0.json
spec/FIXED_CARRIER_PROFILE_3_R0.json
experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json
changes/CHG-ORIGAMI-0015.json
```

## Version

`6.0.0-alpha.15`
