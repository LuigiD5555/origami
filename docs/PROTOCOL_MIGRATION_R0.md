# Protocol migration R0 — alpha.12 baseline to read/write protocol candidate

Status: `IMPLEMENTATION_IN_PROGRESS_FIRST_HALF`

This migration preserves the alpha.12 Native semantic correction and adds the missing communication-protocol symmetry.

## Before the failure

```text
image -> BOOT -> binary pixels -> bytes -> decompression -> understand
```

This made a plain VLM treat Origami as an archive requiring external mechanical decoding.

## Alpha.12 baseline

```text
SEMANTIC PLANE                 EXACT / CONTROL PLANE
T0 -> T1 -> T2 -> answer       T3 -> CID/hash/Merkle/exact
```

T2 became an actual bounded semantic superindex and semantic questions no longer require T3.

## Protocol candidate

```text
T0 BOOT
  -> T1 ROSETTA
       -> semantic grammar
       -> Decoder Registry S*
       -> Encoder Registry E*
       -> capability/fallback declaration
  -> T2 / semantic memory
  -> optional exact/control codecs X*/Q*
  -> VERIFY
```

The change does not mean “no decoder”. It means:

```text
NO UNDECLARED EXTERNAL DECODER DEPENDENCY
```

and similarly for writing:

```text
NO UNDECLARED EXTERNAL ENCODER DEPENDENCY
```

Origami declares how compatible models should read and write its semantic structures.

## First paired codec

The first executable gate is the superindex pair:

```text
E2 ENCODE_SUPERINDEX
S2 READ_SUPERINDEX

S2(E2(INDEX)) ~= INDEX
```

This is tested deterministically in `internal/codec` before any Native VLM claim is made.

## Profile-3 boundary

`origami.fixed-carrier.r2.profile-3` is introduced only as an experimental construction contract in this phase. It is not the default renderer yet.

This preserves profile-1/profile-2 compatibility and prevents a protocol design change from silently changing the frozen current carrier before evidence exists.

## Master Prompt R4

R4 becomes a protocol handshake:

```text
BOOT -> ROSETTA -> discover codecs -> negotiate capability
     -> choose smallest sufficient READ/WRITE codec
```

It no longer treats “decoder” as inherently bad. The forbidden behavior is hidden/undeclared external decoding or unnecessary semantic-to-exact escalation.

## Completed in first half

- Protocol R0 contract;
- Codec Registry R0;
- Capability Negotiation R0;
- S0-S6 / E0-E6 reference registry;
- exact/control X0-X5 / Q0-Q5 separation;
- deterministic S2/E2 roundtrip;
- profile-3 construction contract;
- Writer/ROSETTA/Native-navigation contract integration;
- Master Prompt R4 candidate.

## Deferred to second half

- actual profile-3 visual renderer promotion;
- Tlaloc alpha.16 cross-model codec evaluation;
- Native S2 held-out trials;
- Native E2 write trials;
- A -> B -> C semantic-drift benchmark;
- Tonal composition update for accepted Origami/Tlaloc revisions.
