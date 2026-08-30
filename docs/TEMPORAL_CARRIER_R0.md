# Origami Temporal Carrier — R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTATION`

The Temporal Carrier is a bounded proof that an Origami image can carry both:

```text
SEMANTIC TEMPORAL PLANE
  visible cells / transition graph / timeline / checkpoints

EXACT PROGRAM PLANE
  exact compressed Temporal Program + SHA-256 + CRC32
```

in one image.

It is deliberately **not** a video container. The film is the generative program:

```text
initial state + cells + rules + graph + max steps + checkpoints
```

The visible plane is intended for semantic navigation. The exact payload exists for deterministic recovery/verification and is not a prerequisite for answering semantic questions.

## Envelope

Reference profile:

```text
origami.temporal-carrier.r0.profile-1
640 x 640
8192 PNG bytes exactly
512-byte embedded program record
462 bytes maximum compressed program payload
```

The program record stores canonical Temporal Program JSON compressed with zlib, its uncompressed length, SHA-256 and CRC32.

If the compressed program exceeds the declared capacity, build fails:

```text
OVERSIZE_PROGRAM -> FAIL_BUILD
```

It does not create a sidecar and does not silently enlarge the carrier.

## Visible semantics

The current experimental grammar shows:

```text
BOX   = automaton cell
ARROW = declared transition/dependency
RING  = checkpoint
X axis = temporal order
```

Up to eight cells/rules are projected as a bounded top-level semantic graph in R0. More complex programs may still be invalid for the self-contained R0 envelope even if another future Origami memory profile could reference them selectively.

## CLI

```bash
origami-temporal-carrier \
  -mode build \
  -in program.json \
  -out film.png

origami-temporal-carrier \
  -mode decode \
  -in film.png \
  -out recovered-program.json
```

The canonical signal-chain regression verifies:

```text
program -> 8192-byte PNG -> exact program
                           -> deterministic replay
```

## Scientific boundary

This experiment demonstrates self-containment for programs that actually fit. It does not claim that arbitrary high-entropy trajectories or arbitrary automata fit in 8192 bytes.

```text
COMPRESSIBILITY_DEPENDS_ON_STRUCTURE
```

Native VLM recognition of the visible graph/timeline also remains an empirical gate separate from deterministic roundtrip.
