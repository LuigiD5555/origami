# Origami Protocol R0

Status: `EXPERIMENTAL_REFERENCE_CANDIDATE`

Origami is a self-describing visual/computational communication protocol for exchanging structured state, knowledge and memory between models.

The protocol is distinct from any particular carrier. PNG is the current reference carrier family; future visual, temporal or multilayer carriers may implement the same protocol semantics.

## Communication model

```text
MODEL A
  -> Semantic IR
  -> Origami Encoder
  -> ORIGAMI MESSAGE
  -> Origami Decoder
  -> Semantic IR
  -> MODEL B
```

A compatible message declares enough information for a receiver to discover:

1. protocol/profile identity;
2. active semantic grammar;
3. decoding procedures;
4. encoding procedures;
5. required/optional capabilities and fallbacks;
6. semantic entry points;
7. exact/control mechanisms;
8. failure behavior.

## Required logical sections

```text
T0 BOOT
T1 ROSETTA
CODEC REGISTRY
CAPABILITY DECLARATION
T2 SEMANTIC SUPERINDEX
SEMANTIC MEMORY / PROGRAM
T3 EXACT / CONTROL
VERIFY
```

The physical layout is profile-specific.

## Semantic-first routing

A semantic query selects the smallest semantic decoder sufficient for the task:

```text
S0 READ_IDENTITY
S1 READ_HIERARCHY
S2 READ_SUPERINDEX
S3 LOCATE_TOPIC
S4 READ_RELATION
S5 EXPAND_CLUSTER
S6 SEMANTIC_UNFOLD
```

Exact/control decoding is a separate capability family. Semantic navigation must not escalate to exact decoding merely because an exact plane exists.

## Read/write symmetry

Origami is bidirectional:

```text
DECODE(ENCODE(S)) ~= S
```

The Core interoperability profile requires paired semantic read/write procedures.

## Capability negotiation

A receiver may support only a subset of declared codecs or perceptual channels.

```text
supported codec      -> execute
unsupported optional -> use declared fallback if available
unsupported required -> UNKNOWN / NOT_VERIFIED / CONSTRUCTION_SPEC_ONLY
```

Unsupported capability never authorizes invention.

## Carrier boundary

`PROTOCOL != CARRIER`.

A compact image is not automatically a successful Origami message if receivers cannot recover its declared semantics.

## Exact/control boundary

CID, hashes, Merkle proofs, residuals and byte-exact recovery remain valid functions of the exact/control plane. They are not prerequisites for semantic identity, index, overview or topic routing.

## Cross-model objective

Reference interoperability eventually measures:

```text
A encode -> B decode
A encode -> B decode/re-encode -> C decode
```

with entity, relation, hierarchy, evidence, uncertainty and state preservation plus semantic drift per hop.

## Hard invariants

- `PROTOCOL_NE_CARRIER`;
- `ROSETTA_ALWAYS_PRESENT`;
- `ROSETTA_DECLARES_DECODING_SEMANTICS`;
- `ROSETTA_DECLARES_ENCODING_SEMANTICS`;
- `NO_UNDECLARED_EXTERNAL_DECODER_DEPENDENCY`;
- `NO_UNDECLARED_EXTERNAL_ENCODER_DEPENDENCY`;
- `SEMANTIC_NAVIGATION_MUST_NOT_REQUIRE_EXACT_CODEC`;
- `UNSUPPORTED_CAPABILITY_NE_INVENTION`;
- `DECODE_ENCODE_ROUNDTRIP_REQUIRED`;
- `FALSE_EXACT=0`.

## Evidence boundary

This contract defines the protocol candidate. It does not prove held-out cross-model read/write interoperability; that remains evidence-gated.