# Origami generative machine model

Status: recovery contract for the Origami 6.x core.

## Purpose

Origami is not defined by a screenshot, raster image, PDF page render, or any other single carrier. Origami is a relational, generative, addressable and executable Fold/Unfold machine. Visual projection is one possible representation surface of an Origami machine state.

The core question is not "can a model read this image?". The core question is whether a source can be transformed into an explicit Origami representation whose structure can be addressed, evolved, observed, selectively unfolded and verified.

## Five core verbs

### FOLD

`FOLD` transforms source structure into an Origami representation. A valid fold must identify and encode reusable structure rather than merely rasterizing the source.

Candidate representation mechanisms include literals, references, concatenation, repetition, transforms, rules, graphs, residuals and verification records.

### RELATE

`RELATE` makes dependencies, hierarchy, references, constraints and other declared relations explicit and addressable. Relations belong to the machine representation; they are not inferred retrospectively from a screenshot.

### EVOLVE

`EVOLVE` applies declared operations and rules to machine state. The R0 deterministic transition law remains:

```text
S_(t+1) = F(S_t, C_t, R)
```

where state, context and rules are explicit inputs.

### OBSERVE

`OBSERVE` is separate from transition. Non-direct claims require finite observation contracts with declared preconditions, observer/procedure, authorized trajectory, budget, predicate and failure semantics.

Observation outcomes remain distinct: `PASS`, `FAIL`, `INVALID_CONTRACT`, and `UNSUPPORTED`.

### UNFOLD

`UNFOLD` materializes a requested part of a representation. The preferred path is selective: resolve an address or query, compute the minimum dependency closure, unfold only the required structure, and verify the result. Exact full unfolding is a separate operation used when the selected profile declares exact recovery.

## Core representation model

```text
SOURCE
  ↓
canonical structure
  ↓
FOLD
  ↓
ORIGAMI REPRESENTATION
├── state
├── relations
├── rules
├── graph
├── grammar
├── references
├── transforms
├── residual
├── verification
└── index
  ↓
RELATE / EVOLVE / OBSERVE
  ↓
address resolution
  ↓
selective UNFOLD
  ↓
result + provenance + evidence
```

A representation may later be projected to JSON, a visual carrier, a temporal sequence, an OHF carrier, or another surface. Those projections do not become the semantic source of truth merely because they are visible.

## Generative leverage

Origami must be able to distinguish literal storage from generative representation. For example, storing `ABC` one hundred thousand times is not equivalent to representing a literal `ABC` plus a deterministic repetition rule.

Generative leverage is therefore an engineering property to measure, not an excuse to invent unavailable information. Every generated structure must be justified by the representation and its source-derived contracts. Residual information is retained when exact recovery requires data that cannot safely be generated from rules or references.

## Addressability and selective unfolding

A representation is useful only if relevant structure can be located without reconstructing the entire source by default. The core therefore requires stable addresses and an index capable of mapping a query or declared identifier to a bounded dependency closure.

Useful measurements include:

```text
total logical expansion
representation nodes
nodes visited per query
nodes unfolded per query
bytes materialized per query
residual bytes
verification result
```

## Raster baseline rule

`RASTER_IS_NOT_FOLD` is a core recovery invariant.

Rendering a PDF page, HTML page, text document or other source to pixels may be useful as:

- an ingestion artifact;
- a visual projection;
- a perception fixture;
- an OHF experiment input;
- a control/baseline.

Rasterization alone must never satisfy a gate claiming that Origami generative Fold occurred. A raster-only path has not demonstrated generative structure, semantic addressability, selective unfolding, relational execution or exact generative recovery.

## OHF boundary

OHF remains a carrier/protocol and laboratory research track inside Origami. OHF may test how an Origami representation is projected into a compact perceptual carrier and how safely that carrier can be observed or resolved. OHF carrier success does not replace the requirement to establish the underlying Origami representation.

The OHF `<=500 KB` target, `FALSE_EXACT = 0`, Native isolation and progressive external-testing rules remain OHF track requirements, not universal definitions of the Origami language.

## Recovery acceptance

The recovered core is not considered demonstrated merely because a VLM can answer a question from a rendered page. A minimum recovery demonstration must show:

1. source canonicalization;
2. a non-raster Origami Fold;
3. explicit generative/relational representation nodes;
4. stable addressing;
5. selective unfolding;
6. deterministic verification;
7. exact recovery when the selected profile claims exactness;
8. provenance/evidence sufficient to distinguish generated structure from residual source material.

The first implementation experiments should use small deterministic text/graph fixtures before PDF ingestion or visual carriers. PDF is an ingestion format, not the definition of Origami.
