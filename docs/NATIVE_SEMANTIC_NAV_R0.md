# Native Semantic Navigation R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED_EVIDENCE_PENDING`

This contract originates from `FAILED_TRIAL_001`, where a multimodal model could read the visible Origami bootstrap but could not answer **“What is the index?”**. It treated the carrier as a binary archive requiring pixel extraction/decompression and produced unverified byte/hash claims.

## Alpha.12 correction preserved

```text
SEMANTIC VISUAL PLANE
T0 -> T1/ROSETTA -> T2 semantic superindex -> selective semantic regions

EXACT / CONTROL PLANE
T3 -> CID / hash / root / residual / exact verification
```

T3 remains useful and independent. It is not a prerequisite for semantic navigation.

## Protocol refinement

The correction does **not** mean that Native Origami must have no decoder. It means that semantic decoding must be self-declared by Origami and must not secretly require an external exact/mechanical decoder.

ROSETTA/Codec Registry therefore expose semantic decoder entrypoints:

```text
S0 READ_IDENTITY
S1 READ_HIERARCHY
S2 READ_SUPERINDEX
S3 LOCATE_TOPIC
S4 READ_RELATION
S5 EXPAND_CLUSTER
S6 SEMANTIC_UNFOLD
```

For the index benchmark:

```text
BOOT -> ROSETTA -> discover S2 -> T2 -> answer
```

not:

```text
BOOT -> T3 bits -> bytes -> decompression -> reconstruct corpus -> answer
```

## Decoder boundary

Allowed:

- a decoder procedure declared by ROSETTA/Codec Registry;
- conceptual/perceptual execution that the receiving model can genuinely perform;
- a declared fallback when the preferred channel is unsupported;
- an exact decoder/tool for an exact question when genuinely available.

Forbidden:

- undeclared external decoder dependency for semantic navigation;
- semantic codec that silently calls an exact codec;
- filesystem/sandbox requirement hidden inside a Native semantic claim;
- semantic-to-exact escalation without need;
- invented exact mechanics when the receiver cannot execute them.

## Query routing

```text
identity             -> T0/T1 -> S0
index                -> S2 -> T2
overview             -> S2 -> T2 + visible semantic state
locate topic         -> T2 -> S3
explain topic        -> S3 -> S5/S6
exact quote/hash     -> exact codec/tool when available
```

## T2 profile-2 baseline

Alpha.12 profile-2 remains the current reference Fixed Carrier. It renders bounded actual semantic entries instead of generic PAGE/GRAPH/SOURCE categories and keeps the frozen 640x640 / 8192-byte envelope.

## Profile-3 candidate

The protocol work introduces `origami.fixed-carrier.r2.profile-3` only as a construction candidate. Its first goal is to bind ROSETTA to compact semantic codec IDs and capability/fallback declarations while preserving the useful alpha.12 separation.

Profile-3 is not promoted by this document.

## Native benchmark

```text
Q: What is the index?
Expected route: T0 -> T1 -> S2 -> T2
Exact codec allowed: NO
```

Metrics now distinguish:

- T2 index recovery;
- semantic decoder discovery;
- undeclared external decoder dependency;
- unnecessary semantic-to-exact escalation;
- unverified mechanical claims;
- false exactness.

## Hard interpretation

```text
SELF_DECLARED_SEMANTIC_DECODER = ALLOWED
UNDECLARED_EXTERNAL_DECODER_DEPENDENCY = FORBIDDEN
SEMANTIC_NAVIGATION_MUST_NOT_REQUIRE_EXACT_CODEC
FALSE_EXACT = 0
```

Held-out model evidence is still required before Native semantic capability is promoted.
