# Capability Negotiation R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`

Origami cannot assume that every receiver or writer can execute every declared codec or perceptual channel. Capability negotiation defines deterministic behavior when capabilities differ.

## Capability classes

```text
SEMANTIC_READ
SEMANTIC_WRITE
EXACT_READ
EXACT_WRITE
PERCEPTUAL_CHANNEL
COMPILER
VERIFIER
```

A receiver may declare support at operation granularity, for example:

```text
S0 S1 S2 S3 S4 supported
S5 S6 unsupported
X0 supported
X1-X5 unsupported
MOIRE unsupported
```

## Negotiation rule

For every requested operation:

1. discover the codec from ROSETTA/Codec Registry;
2. inspect required capabilities;
3. execute if supported;
4. otherwise try a declared fallback;
5. otherwise return the declared failure state.

Never infer that a capability exists merely because the protocol describes it.

## Fallbacks

Fallbacks must preserve the semantic class they claim to preserve. An advanced perceptual channel may fall back to canonical topology/geometry only when the profile explicitly declares that fallback as semantically sufficient.

Fallbacks may lose fidelity; that loss must be reported rather than hidden.

## Required failure behavior

```text
semantic read unsupported -> UNKNOWN
exact read unsupported    -> NOT_VERIFIED
semantic write unsupported-> CONSTRUCTION_SPEC_ONLY or UNKNOWN
compiler unsupported      -> CONSTRUCTION_SPEC_ONLY
verification unsupported  -> NOT_VERIFIED
```

## Core compatibility

The initial Core target is deliberately small:

```text
T0 / BOOT
T1 / ROSETTA
S0/E0 identity
S1/E1 hierarchy
S2/E2 superindex
UNKNOWN / NOT_VERIFIED
```

A model can therefore participate in Origami communication without supporting advanced visual channels or exact recovery.

## Hard invariants

- `CAPABILITY_IS_DECLARED_OR_OBSERVED_NOT_ASSUMED`;
- `UNSUPPORTED_CAPABILITY_NE_INVENTION`;
- `FALLBACK_IS_EXPLICIT`;
- `FALLBACK_FIDELITY_IS_DECLARED`;
- `CORE_COMPATIBILITY_DOES_NOT_REQUIRE_EXACT_PLANE`;
- `TOOL_ASSISTED_NE_NATIVE_CAPABILITY`;
- `FALSE_EXACT=0`.