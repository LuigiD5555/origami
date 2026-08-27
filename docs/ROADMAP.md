# OHF roadmap

This roadmap tracks the research sequence, not release marketing.

## Completed laboratory foundation

### F1 — reproducible workbench

Implemented through the evaluator/failure-registry foundation:

- deterministic experiment identity;
- seed namespaces;
- truth-first ground truth;
- artifact manifests;
- deterministic LAB renderer;
- save/reopen verification;
- Native isolation;
- structured observations;
- external scoring;
- regression capture.

## Active experiments


### Dimensional Visual Register / Perceptual State Space

Implemented locally as an experimental formalization. Next evidence steps:

1. controlled state-visibility probes per leaf dimension;
2. one-dimension mutation tests;
3. pairwise orthogonality matrix;
4. correlated-failure analysis;
5. higher-order compound collision tests;
6. only then calculate non-zero `SAFE_PSS`.

The optimization target is safe perceptual bits per visual unit, not raw cartesian cardinality.

### Glyph Calculus

Reduce open-ended LLM work to local glyph observation and deterministic Go execution.

Next evidence target: demonstrate model-specific safe perceptual operations on real VLMs without allowing false-known outputs.

### SAFE_MICRO_ISA

Profile primitive operations rather than guessing which ones are cheap or reliable.

Current progression:

```text
OBSERVE4 / TUPLE4 / DECODE8 @ width 1
       ↓
select safe winner
       ↓
width 4 -> 8 -> 16
       ↓
confirm safe/fail boundary
```

Future candidate operations may include small equality, selection, copy, reference and lookup operations, but none should enter the safe profile without evidence.

### Context SIMD

Test whether independent lanes reduce per-glyph interaction cost without increasing false-known outputs.

## Next visual-language phase

### Macro-Gestalt

Goals:

- represent global topology/composition as measurable physical state;
- reuse existing OHF visual families rather than invent semantics casually;
- keep physical appearance separate from carrier-local meaning;
- test primitives before compounds;
- require zero deterministic rendering collisions before external tests;
- measure cross-model perception before promotion.

Candidate physical axes include topology, symmetry, orientation and macro-region structure.

Macro-Gestalt remains experimental and unpromoted.

## Evolution control

Expand `state/components/` until every major OHF component has explicit status, latest evidence and promotion state. Generate human-readable continuity from machine state rather than reconstructing it manually.

## Historical baseline recovery

Obtain and verify the exact historical R3.10 SDK artifact. Until then:

```text
MISSING_R310_SDK
```

remains an explicit blocker for byte-identical baseline claims.

## Longer-term protocol work

- Representation Tournament implementation;
- graph/grammar/transform folding;
- SuperIndex and dependency manifests;
- multidimensional attention routing;
- macro/meso/micro visual compilation;
- restoration-gated Perception Lab;
- model-specific safe visual profiles;
- generated Master Prompt from canonical Go protocol state;
- cross-model Native and Hybrid evaluation;
- promotion to a post-R3.10 protocol candidate only after evidence supports it.
