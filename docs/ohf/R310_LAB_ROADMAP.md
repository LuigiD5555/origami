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

### EXP-002 — addressable PDF carrier R0

`ohf-lab` can ingest a bounded PDF page range as canonical extracted text, materialize a flat `page:N -> byte range` index, render page tiles in a deterministic visual carrier, and unfold only pages selected by deterministic query-token overlap.

This is the first executable step toward SuperIndex and selective unfolding. R0 deliberately has narrow claims:

- the index is flat and page-level, without grammar, motifs, graphs, or hierarchy;
- retrieval uses normalized token overlap and selects at most three matching pages, not semantic search or ML;
- the query answer is the selected extracted text, without an LLM;
- source truth covers canonical extracted text, not byte-exact recovery of the original PDF;
- exact PDF recovery is deferred to EXP-002-B;
- Native evaluation remains external and can use the existing `native bundle` command.

The reference experiment uses PDF pages 17–20 so that it runs against the currently staged `Thinking Swarms.pdf`; several earlier pages contain malformed content streams rejected by the pure-Go parser. Larger clean page ranges are supported and are preferable for measuring closure ratios. The local source path may need adjustment on another workstation.

The first five deterministic smoke queries on that four-page range unfolded 36.7%, 26.5%, 63.3%, 82.5%, and 80.7% of canonical bytes. These measurements demonstrate selective unfolding but do not meet the ideal `<20%` target; the fixture is too small for that claim. No exact-recovery claim is made.

### EXP-003 — preregistered SuperIndex capacity tournament

The capacity question is evaluated per declared representation scope rather than with a universal `fits_carrier` flag. `experiments/EXP-003-superindex-capacity-001.json` fixes the fixture, seed-independent physical parameters, candidates and PASS/FAIL rules before the tournament run.

Run `93ae706b66dc7e21` measured three candidates. Direct JSON failed both the surface and PNG budgets. Deterministic zlib reduced the index from 11,860 to 3,822 bytes and passed the PNG budget, but still exceeded the 2,420-glyph surface. A 250-byte referenced manifest passed both budgets while declaring its external SHA-256 payload dependency; this is evidence for a referenced routing carrier, not self-contained compression. No candidate is promoted as perceptually safe.

`superindex-carrier@v1` materializes that referenced candidate as a separate renderer. Preregistered run `2f68a4155acf3c67` produced a 24,687-byte PNG; `run verify` independently reconstructed the canonical SuperIndex hash and passed the referenced-payload gate. Native readability and retrieval superiority remain untested, so the renderer is not promoted.


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
