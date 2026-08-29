# OHF Glyph Calculus v0.1 — Experimental Protocol

## Status

**EXPERIMENTAL / LAB-ONLY / NOT PROMOTED TO R3.10**

This experiment tests a specific hypothesis: reduce the model's job to local visual perception and move all unfolding/computation into deterministic Go code.

## Core split

```text
IMAGE
  ↓
VLM: observe each glyph only
  ↓
(shape, holes, direction, frames)
  ↓
byte = shape + 4*holes + 16*direction + 64*frames
  ↓
BYTE STREAM
  ↓
Go bounded glyph calculus
  ↓
reconstructed bytes
  ↓
SHA-256 / exact verification
```

The VLM is explicitly **not** asked to infer semantics, decompress, execute a program, follow a long chain of thought, decide exactness, or assign semantic roles to physical shapes.

## Physical alphabet

The existing R3.10 physical byte mapping is reused:

- F01 outer shape: circle=0, square=1, diamond=2, triangle=3
- F02 topology holes: 0..3
- F03 direction beacon: top=0, right=1, bottom=2, left=3
- F04 frame count: 0..3

Formula:

```text
byte = shape + 4*holes + 16*direction + 64*frames
```

## Transport modes

### PROGRAM

Header `A0`, followed by bounded micro-operations:

| Opcode | Name | Mechanical effect |
|---:|---|---|
| `00` | HALT | stop |
| `01` | LIT | emit one literal byte |
| `02` | RUN | emit one byte repeatedly |
| `03` | REF | copy a bounded previous region |
| `04` | SEQ | emit arithmetic byte sequence modulo 256 |

Unknown opcodes, invalid references and missing HALT fail closed.

### RAW

Header `A1 + uint16 length + raw bytes`.

This is the entropy fallback. OHF does not force a generative representation when structure does not reduce cost.

## Why this is calculus-like

The analogy is procedural compression: a short instruction such as `RUN CC 80` mechanically produces many bytes. Intelligence is spent at encode time discovering/choosing the representation; receive-time execution becomes algorithmic.

## Current local evidence

For deterministic 4,096-byte fixtures, structured patterns locally use substantially fewer transport glyphs while reconstructing exactly. Pseudo-random control falls back to RAW. These results do **not** demonstrate Native visual readability.

## Promotion requirements

At minimum: Recognition, Invariance, OOD, sequence integrity, random addressing, UNKNOWN on ambiguity, exact restoration through Go, zero false exact, prior-collision testing and sufficient cross-model evidence for the target profile.
