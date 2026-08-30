# Fixed Carrier R2 profile-3 candidate

Status: `EXPERIMENTAL_RENDERER_IMPLEMENTED_EVIDENCE_PENDING`

`origami.fixed-carrier.r2.profile-3` is the first rendered carrier candidate designed explicitly for Origami Protocol R0.

It does **not** replace profile-2 yet. Profile-1 and profile-2 remain compatible, and the profile-2 renderer remains the default.

## Physical target

```text
640 x 640
8192 PNG bytes target
<= 512000 hard maximum
```

Profile-3 is rendered in parallel by `internal/fixedcarrier.RenderProfile3` and the `origami-profile3-carrier` CLI. Deterministic roundtrip/size tests are required before this branch can merge.

## Logical layout

```text
T0 BOOT
  protocol/profile identity
  classify request before selecting a codec

T1 ROSETTA + CODEC REGISTRY
  semantic grammar
  S0/E0 identity
  S1/E1 hierarchy
  S2/E2 superindex
  core capability declaration

T2 ACTUAL SEMANTIC SUPERINDEX
  bounded real semantic entries

PROGRAM
  self-declared semantic codec procedures
  S2: T2 -> index
  E2: index -> T2 Construction IR
  declared fallbacks only

T3 / EXACT CONTROL
  X*/Q* exact codec family
  CID / hashes / Merkle / residual
  explicitly not required for T2 semantic navigation

VERIFY
  S2(E2(index)) ~= index
  FALSE_EXACT=0
```

## Read route

For `What is the index?`:

```text
T0 -> T1 -> discover S2 -> T2 -> answer
```

The receiver must not escalate to X*/T3 merely because the exact plane exists.

## Write route

For a bounded semantic index:

```text
Semantic IR -> discover E2 -> T2 Construction IR
```

A model without a deterministic image compiler may legitimately stop at `CONSTRUCTION_SPEC_ONLY`. A compiler-capable receiver can render profile-3 and verify the result.

## Decoder/encoder rule

The regression rule is **not** `NO DECODER`.

It is:

```text
SELF_DECLARED_SEMANTIC_CODEC = ALLOWED
UNDECLARED_EXTERNAL_CODEC_DEPENDENCY = FORBIDDEN
SEMANTIC_TO_EXACT_ESCALATION_WITHOUT_NEED = FORBIDDEN
```

## Promotion gates

- rendered PNG remains 8192 bytes or an explicit profile-size change is approved;
- profile-1/profile-2 decoding is not regressed;
- deterministic profile-3 render/decode roundtrip passes;
- S2 index recovery meets held-out Native evidence gate;
- E2 construction meets held-out write gate;
- cross-model roundtrip evidence exists;
- no undeclared external codec dependency;
- no unnecessary semantic-to-exact escalation;
- `FALSE_EXACT=0`.

## Current evidence boundary

Renderer implementation is not model evidence. Until held-out model trials are imported and scored, profile-3 remains experimental and must not be described as Native-promoted.
