# Origami Capacity Lab R0

Status: EXPERIMENT DESIGN

## Purpose

Origami experiments are not demonstrations that a chosen capacity exists. They are instruments for discovering where representation, addressability, selective access and semantic computation begin to fail.

The central question is:

> How much external semantic state can Origami make correctly usable while the active working window remains bounded?

R0 uses an equivalent working-window budget of 4,000 model tokens as the primary reference point. This is a budget, not a claim that Origami state is token-based.

## Separation of concerns

1. Correctness gates (ADV): prove claimed semantics and reject shortcuts.
2. Capacity experiments (CAP): increase represented state until failure.
3. Access experiments (ACC): hold active window constant while represented state grows.
4. Complexity experiments (CMP): increase dependency depth, branching and interaction order.

## Primary variables

- M: total represented semantic units.
- W: active working-window budget, reported in bytes, semantic units, and token-equivalent estimate.
- C(q): dependency closure unfolded for query q.
- D: maximum dependency depth reached.
- B: dependency branching factor.
- H: highest participating-set/interference order required by the query.
- T: temporal span required by the observation.

## Required measurements per query

- correctness outcome;
- addresses requested;
- addresses visited;
- semantic units unfolded;
- unrelated units unfolded;
- closure size;
- closure depth;
- active bytes exposed;
- token-equivalent active budget;
- lookup/unfold/derive latency separately;
- provenance/evidence IDs;
- literal and residual contribution;
- whether the query exceeded W;
- failure reason when unsuccessful.

A PASS without access evidence is INVALID for capacity claims.

## Experiment families

### CAP-001 Memory scale under fixed W

Keep W ~= 4,000 token-equivalent. Increase M through a geometric ladder, initially:

1K -> 10K -> 100K -> 1M -> 10M semantic units.

At each scale use deterministic held-out queries whose answers are not stored as dedicated literals. Measure correctness and accessed fraction.

Goal: find the first scale where address lookup, closure construction, active budget, latency or correctness degrades.

### ACC-001 Window sweep under fixed M

Hold M constant and test W at approximately:

256, 512, 1K, 2K, 4K, 8K token-equivalent.

Goal: estimate the minimum active budget needed for each query class.

### ACC-002 Selectivity

For each query compute:

selectivity = |C(q)| / M

and separately measure unrelated material exposed. Reading all M and hiding it afterward is a FAIL.

### CMP-001 Dependency depth and branching

Hold M and W stable. Increase D and B independently until the closure no longer fits or lookup/derivation fails.

### CMP-002 Higher-order interaction

Increase H from pairwise to triples and beyond. Pairwise success never promotes higher-order capacity.

### TMP-001 Temporal span

Increase T while ensuring the derived property is absent from individual frames. Measure how much trajectory evidence must enter W.

## Anti-cheat requirements

Capacity results are invalid if any of the following occurs:

- expected answer is stored directly;
- source is copied into a large literal or residual;
- an external model performs the claimed core derivation;
- Selective Unfold reads the entire semantic source before filtering;
- exact byte recovery is presented as semantic computation;
- queries used to construct the representation are reused as hidden answer hints.

## Baselines

Every capacity campaign must include at least:

- exhaustive/full-unfold baseline where feasible;
- trivial full-scan baseline;
- Origami selective route.

Origami only earns an access advantage when it produces equivalent correct semantics while touching materially less represented state.

## Failure is data

The lab must report the first limiting factor, not hide it. Expected classes include:

- ADDRESS_LOOKUP_LIMIT
- CLOSURE_EXPLOSION
- WINDOW_EXCEEDED
- DEPENDENCY_DEPTH_LIMIT
- HIGHER_ORDER_LIMIT
- TEMPORAL_SPAN_LIMIT
- LATENCY_LIMIT
- SEMANTIC_ERROR
- ROUTE_INTEGRITY_FAILURE
- UNSUPPORTED

The purpose is to map the operating envelope, not maximize PASS counts.
