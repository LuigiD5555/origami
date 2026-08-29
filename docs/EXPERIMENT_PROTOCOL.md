# OHF experiment protocol

The purpose of the laboratory is not to produce attractive examples. It is to produce repeatable evidence about which representations are safe enough to use.

## Non-negotiable procedure

Every material experiment follows this order:

```text
1. define ExperimentSpec
2. canonicalize spec
3. derive deterministic seeds
4. create ground truth
5. render candidate
6. save artifact
7. reopen artifact
8. verify file/pixel identity
9. create blind Native bundle when needed
10. collect ObservationEnvelope
11. score outside the Native boundary
12. register reproducible failures
13. replay regression
14. decide reject / continue / promotion candidate
```

Ground truth is created **before** the image. A model never defines what the expected answer was.

## Identity

`ExperimentID` identifies a canonical experiment specification.

`RunID` additionally binds the materialized source and implementation identities needed to reproduce one concrete run.

The two identifiers are intentionally different.

## Seeds

Randomness is namespaced. A master seed derives independent sub-seeds for concerns such as layout, motifs, degradation and sampling.

Adding a new random consumer must not shift every unrelated random sequence.

The seed is a reproducibility mechanism. Native decoding must not depend on a hidden seed.

## Visual identity

The LAB records both:

```text
PNG file SHA-256
canonical decoded-pixel SHA-256
```

This distinguishes a different PNG byte stream representing the same pixels from an actual visual change.

## Native isolation

Every Native request uses a physical `blind/` directory containing exactly:

```text
MASTER_PROMPT.md
carrier.png
query.txt
```

A fourth file inside the blind directory is a protocol violation.

## Observation model

Native output is normalized to an `ObservationEnvelope`.

Allowed epistemic states:

```text
KNOWN
AMBIGUOUS
UNKNOWN
INVALID
```

The scorer is a separate component with access to expected truth.

## Failure policy

A reproducible real failure becomes a regression. Regressions preserve sufficient evidence to recreate the failure without exposing ground truth to the original Native request.

## Local before external

Always execute cheap deterministic gates before a VLM request:

```text
schema
build/tests
determinism
pixel uniqueness
carrier size
manifest verification
```

If a local structural gate fails, external testing is forbidden because the request cannot answer the relevant engineering question.

## Reference-engine transitions

The deterministic relational reference engine applies transition phases in this order:

```text
requires -> excludes -> couples -> toggles
```

`requires`, `excludes`, and `couples` retain their R0 behavior. The explicit `toggles` relation adds bounded oscillation: when `from` is `active`, `to` alternates between `active` and `inactive`; a non-binary target becomes `active`. A toggle is not intrinsically contradictory. Contradictions continue to be evaluated before each transition from violated `requires` and `excludes` relations.

Experiments that depend on an observed outcome declare minimum counts in `protocol.required_classifications`. The evidence gate must reject a complete, deterministic run when a required classification is absent.


## Dimensional visual experiments

Before an external perception request, dimension work follows an additional cheap-to-expensive progression:

```text
state visibility
  ↓
controlled one-dimension mutation
  ↓
deterministic collision check
  ↓
pairwise orthogonality
  ↓
correlated failures
  ↓
higher-order compounds
  ↓
invariance / degradation / OOD
  ↓
transfer / restoration where applicable
```

A pair missing from the orthogonality matrix is `UNKNOWN`, never implicitly PASS. A pairwise-compatible set is only a candidate for compound testing; it is not a promoted alphabet.

Capacity reports must keep three scopes separate: MICRO per cell, MESO per group and MACRO per carrier.

## Progressive external testing

Do not send a full combinatorial campaign immediately.

Preferred progression:

```text
Stage A: tiny sentinel set
Stage B: pairwise/covering sample
Stage C: hard negatives / OOD
Stage D: larger campaign only when justified
```

SAFE_MICRO_ISA is even more restrictive: width-1 candidate comparison first, then widen only the winner.

## Service-friendly policy

External execution is intentionally conservative:

- default concurrency `1`;
- configurable delay between requests;
- honor `429` and `Retry-After`;
- exponential backoff with bounded retries;
- response cache keyed by the effective request identity;
- resume without repeating completed probes;
- circuit breaker on repeated errors;
- stop early when the experiment has already failed structurally.

This is both respectful and cheaper.

## Promotion law

A successful example is not a promoted capability.

Where applicable, perceptual capabilities should progress through:

```text
Recognition
Invariance
OOD Generalization
Composition
Transfer
```

and include hard negatives, prior-collision testing, degradation and real restoration gates.

`FALSE_EXACT = 0` remains mandatory at every stage.
