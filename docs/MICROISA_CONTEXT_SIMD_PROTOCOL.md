# OHF SAFE_MICRO_ISA + CONTEXT_SIMD v0.1

Status: **EXPERIMENTAL / LAB-ONLY / NOT PROMOTED**

## Hypothesis

Minimize VLM work by discovering an empirically safe set of local operations and batching independent operations into context lanes. Black-box internal FLOPs are not observable, so cost is ranked with observable proxies: exact-known/false-known rates, latency, output tokens, retries and monetary cost when available.

The hard invariant remains `FALSE_EXACT = 0`.

## Candidate micro-ops

- `OBSERVE4`: emit four physical state integers; no byte arithmetic.
- `TUPLE4`: emit four base-4 digits such as `2031`; no byte arithmetic.
- `DECODE8`: observe the states and apply the deterministic byte formula.

An operation/width is safe only if `false_known == 0`, exact-known rate meets the threshold, lanes remain independent, ambiguity is surfaced as UNKNOWN, and evidence belongs to the exact prompt+carrier+model profile tested.

## Context SIMD

Context SIMD means independent-lane batching, not a claim about hardware scheduling inside a model. The widening sequence is:

```text
1 -> 4 -> 8 -> 16
```

The first unsafe boundary stops widening.

## Adaptive external testing

Stage A uses exactly one width-1 case for each candidate op. Select the safe winner; Stage B widens only that operation. Stage C confirms only the last safe and first failing boundary when needed.

The command:

```bash
go run ./cmd/ohf-glyphcalc microbench-next --root runs/microbench --min-exact 0.98
```

returns exactly the next allowed probe or STOP.

## Service etiquette

Default: concurrency 1, no burst, respect `429`/`Retry-After`, bounded retries, exponential backoff, cache completed probes, resume rather than repeat and stop on systematic failure.

The supplied kit never contacts an external service by itself.

## Receive path under test

```text
IMAGE
  -> local tuple observations
  -> deterministic tuple-to-byte packing
  -> deterministic Glyph Calculus
  -> exact verification
```

rather than asking the VLM to understand, decompress, execute and verify the carrier itself.
