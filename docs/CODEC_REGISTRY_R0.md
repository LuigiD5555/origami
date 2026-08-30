# Codec Registry R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`

The Codec Registry is the self-description layer that tells a compatible model both **how to read Origami** and **how to produce compatible Origami**.

ROSETTA owns semantic bindings; the Codec Registry owns executable/conceptual procedures attached to those bindings.

## Semantic decoder ABI

```text
S0 READ_IDENTITY
S1 READ_HIERARCHY
S2 READ_SUPERINDEX
S3 LOCATE_TOPIC
S4 READ_RELATION
S5 EXPAND_CLUSTER
S6 SEMANTIC_UNFOLD
```

Semantic decoders must operate on the semantic plane and must not silently require exact/control decoding.

## Semantic encoder ABI

```text
E0 ENCODE_IDENTITY
E1 ENCODE_HIERARCHY
E2 ENCODE_SUPERINDEX
E3 ENCODE_CONCEPT
E4 ENCODE_RELATION
E5 ENCODE_CLUSTER
E6 SEMANTIC_FOLD
```

A text-only model may execute an encoder by emitting an Origami Construction IR. If no deterministic compiler is available, the result remains `CONSTRUCTION_SPEC_ONLY`.

## Exact/control codecs

Exact codecs are declared separately because they may require deterministic runtime capabilities.

Reference decoder roles:

```text
X0 RESOLVE_ADDRESS
X1 READ_CID
X2 VERIFY_HASH
X3 VERIFY_MERKLE
X4 DECODE_RESIDUAL
X5 EXACT_RECOVERY
```

Reference encoder roles:

```text
Q0 ASSIGN_ADDRESS
Q1 CREATE_CID
Q2 HASH_CONTENT
Q3 BUILD_MERKLE
Q4 ENCODE_RESIDUAL
Q5 BUILD_VERIFICATION
```

The existence of X/Q codecs never makes them prerequisites for S/E semantic operations.

## Registry entry

Each codec declares at least:

```text
id
family
operation
input semantic class
output semantic class
required capabilities
optional fallbacks
exactness class
procedure steps
paired codec when applicable
```

## Roundtrip

Paired semantic codecs are tested as:

```text
S0(E0(identity)) ~= identity
S1(E1(hierarchy)) ~= hierarchy
S2(E2(index)) ~= index
...
```

The first mandatory Native pair is S2/E2 because index recovery is the current regression boundary.

## Selection rule

The receiver/writer chooses the smallest sufficient codec. Semantic questions must prefer S* over X*. Semantic writing must prefer E* over Q* unless exact/control material is genuinely required.

## Hard invariants

- `CODEC_IS_DECLARED_BEFORE_USE`;
- `SEMANTIC_CODEC_NE_EXACT_CODEC`;
- `S_CODEC_DOES_NOT_HIDE_X_DEPENDENCY`;
- `E_CODEC_DOES_NOT_HIDE_Q_DEPENDENCY`;
- `PAIRED_SEMANTIC_CODEC_HAS_ROUNDTRIP_GATE`;
- `UNSUPPORTED_CODEC_RETURNS_DECLARED_FAILURE_STATE`;
- `FALSE_EXACT=0`.