# Fixed Carrier R2 profile-3 candidate

Status: `EXPERIMENTAL_CONSTRUCTION_SPEC_ONLY`

`origami.fixed-carrier.r2.profile-3` is the first carrier candidate designed explicitly for Origami Protocol R0.

It does **not** replace profile-2 yet. Profile-1 and profile-2 remain the currently decodable historical/reference profiles.

## Goal

Preserve the existing Fixed Carrier envelope while making the protocol read/write path self-described:

```text
640 x 640
8192 PNG bytes target
<= 512000 hard maximum
```

## Logical layout

```text
T0 BOOT
  protocol/profile identity
  semantic-first route

T1 ROSETTA
  semantic grammar
  codec registry entry points
  capability/fallback declaration

T2 ACTUAL SEMANTIC SUPERINDEX
  bounded real semantic entries

PROGRAM / SEMANTIC MEMORY
  semantic decoder/encoder procedures and addresses

T3 EXACT / CONTROL
  CID / hashes / Merkle / residual / exact verification

VERIFY
```

## Compact codec declaration

The visual carrier must not become a large textual manual. Profile-3 therefore references stable codec IDs such as `S2`/`E2`, while ROSETTA binds those IDs to the active visual profile and any local parameters.

The universal Master Prompt teaches the receiver how to discover and invoke the registry. The carrier declares which procedures are active and how they bind to this message.

## First Native requirement

```text
Question: What is the index?
Route: T0 -> T1 -> S2 -> T2
Exact codec: forbidden unless the question changes to exact verification
```

## First write requirement

```text
Input: bounded semantic index
Route: Semantic IR -> E2 -> T2 Construction IR
Roundtrip: S2(E2(index)) ~= index
```

## Promotion gates

- 8192-byte envelope remains satisfied or an explicit profile-size change is approved;
- profile-1/profile-2 decoding is not regressed;
- S2 index recovery meets Native evidence gate;
- E2 index construction meets write gate;
- S2/E2 roundtrip is non-regressive;
- no undeclared external codec dependency;
- no unnecessary semantic-to-exact escalation;
- `FALSE_EXACT=0`.

## Current boundary

This phase defines and validates the profile construction contract. It does not yet promote profile-3 as the default renderer and does not claim held-out VLM evidence.