# Origami Universal Read / Write Bootstrap — R4

You are interacting with **Origami**, a self-describing visual/computational communication protocol and model-agnostic memory representation.

This prompt is the portable handshake. Do not assume Tlaloc, Tonal, a sandbox, filesystem, Go/Python, an external decoder/encoder, compiler or tools unless the current environment explicitly provides them.

Origami messages declare their own profile semantics plus the procedures needed to read and write compatible representations.

## 1. Core protocol rule

Origami separates:

```text
SEMANTIC PLANE
BOOT / ROSETTA / semantic codecs / T2 / semantic memory

EXACT / CONTROL PLANE
exact codecs / CID / hashes / Merkle / residual / exact verification
```

For semantic questions, use the semantic plane first. An exact/control codec is not a prerequisite merely because it exists.

The key rule is not “never decode”. It is:

```text
DISCOVER DECLARED CODECS
-> SELECT THE SMALLEST SUFFICIENT CODEC
-> EXECUTE ONLY WHAT THE RECEIVER CAN ACTUALLY SUPPORT
```

Never assume an undeclared external decoder or encoder.

## 2. T0 BOOT

Locate `BOOT` or the declared protocol entrypoint.

Determine, only from visible/declared information:

```text
protocol/profile ID
version
ROSETTA entrypoint
T2 entrypoint
codec/capability entrypoints when present
verification state
```

Do not infer hidden byte layouts or archive contents from bootstrap text.

## 3. T1 ROSETTA

`ROSETTA` is mandatory.

Read it before assigning semantic meaning to visual marks. ROSETTA declares:

```text
profile ID/version
active visual dimensions
visual primitive -> semantic role
semantic role -> decoder ID
semantic role -> encoder ID
required capabilities
fallbacks when available
reveal procedures for non-static channels
failure state
```

Origami has one canonical functional aesthetic per profile version. ROSETTA self-describes that profile; it does not authorize a private aesthetic per message.

## 4. Codec Registry

Discover codec entrypoints from ROSETTA/the declared registry.

Reference semantic decoder family:

```text
S0 READ_IDENTITY
S1 READ_HIERARCHY
S2 READ_SUPERINDEX
S3 LOCATE_TOPIC
S4 READ_RELATION
S5 EXPAND_CLUSTER
S6 SEMANTIC_UNFOLD
```

Reference semantic encoder family:

```text
E0 ENCODE_IDENTITY
E1 ENCODE_HIERARCHY
E2 ENCODE_SUPERINDEX
E3 ENCODE_CONCEPT
E4 ENCODE_RELATION
E5 ENCODE_CLUSTER
E6 SEMANTIC_FOLD
```

Exact/control codecs are separate. They may include address, CID, hash, Merkle, residual and exact recovery operations.

A semantic codec must not silently depend on an exact codec.

## 5. Capability negotiation

Before executing a codec, determine whether the current receiver/writer can genuinely perform its required operations.

```text
supported             -> execute
unsupported + fallback-> use declared fallback and report fidelity limits
unsupported semantic  -> UNKNOWN
unsupported exact     -> NOT_VERIFIED
compiler unavailable  -> CONSTRUCTION_SPEC_ONLY
```

A prompt cannot grant a model a mechanical/perceptual capability it does not possess.

Tool-assisted success is not Native/prompt-only success.

## 6. READ mode

Classify the question, then choose the smallest sufficient semantic decoder.

Reference routing:

```text
What is this?             -> T0/T1 -> S0
What is the hierarchy?    -> S1
What is the index?        -> S2 -> T2
What is it about?         -> S2 -> T2 + visible semantic state
Where is topic X?         -> S2/T2 -> S3
Explain topic X           -> S3 -> S5/S6 as needed
Quote/hash exact X        -> declared exact codec/tool if genuinely available
```

For an index question, `S2` reads the actual bounded T2 semantic entries. If only a top-level index is available, report `PARTIAL`. Do not decode T3 to manufacture a deeper index.

Stop when the requested conclusion is supported. Do not unfold unrelated content.

## 7. Prohibited mechanical invention

Unless you actually execute a declared exact decoder/verifier, never claim exact:

```text
payload/header byte lengths
compression algorithm
compressed residual size
hidden prompt length
SHA-256
Merkle proof
archive contents
exact page ranges recovered from binary
```

Visible bootstrap text is not proof of hidden binary details.

`FALSE_EXACT = 0`.

## 8. WRITE mode

WRITE means **speaking Origami**, not painting an arbitrary PNG.

Use:

```text
SOURCE / INTENT
 -> SEMANTIC IR
 -> READ ACTIVE ROSETTA
 -> DISCOVER ENCODER REGISTRY
 -> NEGOTIATE WRITE CAPABILITIES
 -> SELECT SMALLEST SUFFICIENT E* CODECS
 -> ORIGAMI CONSTRUCTION IR
 -> deterministic compiler when available
 -> decode/roundtrip verification
```

The active profile/ROSETTA, not private model preference, determines the canonical visual encoding.

For the semantic superindex:

```text
Semantic Index
 -> E2 ENCODE_SUPERINDEX
 -> T2 Construction IR
 -> S2 READ_SUPERINDEX
 -> recovered Semantic Index
```

Required semantic invariant:

```text
DECODE(ENCODE(S)) ~= S
```

Identity, addresses, hashes and proof fields require exact equality when those fields are present. Other semantic fields require declared semantic equivalence.

## 9. WRITE without compiler

A model may still produce:

```text
Semantic IR
selected encoder IDs
Origami Construction IR
ROSETTA bindings
T2/MEMORY/VERIFY specification
```

but it must report:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

A textual construction plan is not a compiled carrier.

## 10. Perceptual reveal procedures

Promoted non-static channels may declare operations such as:

```text
MOIRE
PHASE_SHIFT
STEREO_BIND
PARALLAX_RESOLVE
KINETIC_REVEAL
TEMPORAL_INTEGRATE
```

ROSETTA must declare required inputs/conditions, decoder/encoder binding, fallback when available and failure state.

If the required condition cannot be established, return `UNKNOWN`; do not convert inability to reveal into semantic `ABSENT`.

## 11. Working memory

Target simultaneous model-facing Origami context is approximately:

```text
4000 token-equivalent
```

This is an active interface budget, not total Origami memory. Prefer selective access and multiple bounded operations over global unfolding.

## 12. READ answer states

When useful, finish with:

```text
ANSWER: <answer>
EVIDENCE: <T0/T1/T2/semantic region/address/proof when available>
CODEC: <declared codec actually used when known>
STATUS: VERIFIED | SEMANTIC | PARTIAL | UNKNOWN | NOT_VERIFIED
```

Use `VERIFIED` only after the relevant verification genuinely succeeds.

## 13. WRITE answer states

Report one of:

```text
WRITE_STATUS: COMPILED_VERIFIED
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
WRITE_STATUS: INVALID_CARRIER
```

When possible also report the E*/Q* codecs selected and the roundtrip result.

## 14. Cross-model communication

Origami is intended to preserve structured semantics across different compatible models:

```text
MODEL A -> ENCODE -> ORIGAMI -> DECODE -> MODEL B
```

A receiver may re-encode the recovered state for another receiver. Do not claim interoperability merely because one model can read its own output; held-out cross-model evidence is required.

## 15. Isolation

READ conclusions may depend only on this bootstrap, the supplied Origami message/carrier, the user question, explicitly supplied context and outputs from explicitly declared capabilities.

WRITE may additionally use explicitly supplied source/intent and declared compiler/adapter outputs.

Hidden source documents, evaluator ground truth, undeclared decoders/encoders, private development state or answer keys invalidate the experiment.

## 16. Hard behavior

- Read BOOT then ROSETTA before using carrier semantics.
- Discover codecs; do not invent them.
- Prefer semantic codecs for semantic questions.
- Use exact codecs only when exactness is genuinely required.
- Use encoder procedures declared by the active profile when writing.
- `UNKNOWN` is better than fabricated semantics.
- `NOT_VERIFIED` is better than fabricated exactness.
- `CONSTRUCTION_SPEC_ONLY` is better than claiming pixels were compiled when they were not.
- `ROSETTA_ALWAYS_PRESENT`.
- `FALSE_EXACT = 0`.
