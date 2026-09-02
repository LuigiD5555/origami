# Behavior-in-Image R0

Status: **experimental laboratory surface; not promoted protocol semantics**.

## Question

Can a small multimodal model execute a useful complex task when the carrier itself exposes a bounded sequence of very small instructions, instead of asking the model to invent a plan at runtime?

The R0 hypothesis is:

```text
simple visible instructions
+ bounded capabilities
+ fail-closed routing
+ external verification
= complex behavior by composition
```

This is deliberately compatible with tiny models. A step may be as small as `LOCATE_REGION`, `READ_SHORT_TEXT`, `EXTRACT_NUMBER`, `COMPARE_TWO_VALUES`, or `RETURN A/B`.

## Contract

`origami.behavior-recipe.r0` is a forward-only recipe with at most 8 steps. Each step declares:

```text
id
capability
instruction
input_ref
output_ref
on_success
on_unknown
halt
```

R0 forbids backward jumps and undeclared transitions. `UNKNOWN` is a first-class route; failure must not be guessed through.

The recipe is data, not authority. A visible instruction does not make its output verified.

## Renderer

`cmd/origami-behavior-image` renders a recipe as deterministic visible bootstrap text using the same tiny font family as the fixed-carrier BOOT surface.

Example:

```sh
go run ./cmd/origami-behavior-image \
  -recipe experiments/behavior-in-image-r0/recipes/find-number-r0.json \
  -out /tmp/find-number-r0.png
```

The renderer writes a PNG plus a manifest containing recipe and PNG SHA-256 values.

## Experimental comparison

The intended Tlaloc experiment compares the same task under three conditions:

```text
C0  no recipe; model must infer the workflow
C1  equivalent recipe supplied outside the image
C2  recipe visibly precompiled into the Origami image
```

Measure semantic accuracy, contract success, routing errors, calls, latency, tokens, and bytes exposed. C2 is useful only if it preserves or improves task performance without smuggling the answer into the image.

## Boundary

This package proves only deterministic recipe validation/rendering. It does **not** prove that LFM2-VL or any other VLM can read or execute the recipe. That requires held-out external-model evidence before any canonical Origami promotion.
