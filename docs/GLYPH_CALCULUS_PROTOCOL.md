# OHF Glyph Calculus v0.1 — Experimental Protocol

## Status

**EXPERIMENTAL / LAB-ONLY / NOT PROMOTED TO R3.10**

This experiment tests a specific hypothesis: reduce the model's job to local visual perception and move low-level unfolding/computation into deterministic Go code.

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

The VLM is explicitly **not** asked to infer payload semantics, decompress, execute the low-level byte program, follow a long mechanical chain, decide exactness, or assign undeclared semantic roles to physical shapes.

## Physical alphabet

The existing R3.10 experimental physical byte mapping is reused **inside this profile**:

- F01 outer shape: circle=0, square=1, diamond=2, triangle=3
- F02 topology holes: 0..3
- F03 direction beacon: top=0, right=1, bottom=2, left=3
- F04 frame count: 0..3

Formula:

```text
byte = shape + 4*holes + 16*direction + 64*frames
```

This fixed mapping is a Glyph Calculus transport convention. It is **not** a universal semantic meaning for all Origami carrier glyphs. Under `origami.hybrid-receiver.r0`, carrier-level semantic roles are discovered from that carrier's own BOOT/ROSETTA contract. A Rosetta entry may point to or declare a Glyph Calculus transport region without making its shapes globally semantic.

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

The analogy is procedural compression: a short instruction such as `RUN CC 80` mechanically produces many bytes. Intelligence is spent at encode/search time discovering/choosing the representation; receive-time execution becomes algorithmic.

## Relation to Hybrid Receiver R0

Glyph Calculus can act as one **low-level deterministic execution/transport substrate** underneath the Hybrid receiver:

```text
model finds BOOT / interprets carrier-local ROSETTA
                 ↓
       BOOT selects a program/region
                 ↓
        Glyph Calculus / MicroISA
        performs mechanical work
                 ↓
        bounded value / proof ref
                 ↓
               model
```

This does not require the model to semantically understand every Glyph Calculus byte. Conversely, a successful Glyph Calculus reconstruction alone does not prove that a model can self-bootstrap an Origami carrier. The Hybrid campaign measures those layers separately.

## Current local evidence

For deterministic 4,096-byte fixtures, structured patterns locally use substantially fewer transport glyphs while reconstructing exactly. Pseudo-random control falls back to RAW. These results do **not** demonstrate Native BOOT/ROSETTA readability or end-to-end Hybrid navigation.

## Promotion requirements

At minimum: Recognition, Invariance, OOD, sequence integrity, random addressing, UNKNOWN on ambiguity, exact restoration through Go, zero false exact, prior-collision testing and sufficient cross-model evidence for the target profile.
