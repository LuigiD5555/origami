# Accessible Memory Synthetic Campaign — R0 Results

Status: deterministic reference/simulation campaign. These results validate semantics and scaling hypotheses; they do **not** yet constitute an isolated LLM-navigation or real-source proof.

Primary interface target: `W = 4000 token-equivalent` peak active Origami exposure.

## AM-0 — interface/context recycling

1000 sequential exchanges, each exposing ~120 token-eq and retaining 40 token-eq semantic state:

- peak active: 160 token-eq
- cumulative exposure: 120,000 token-eq
- final retained state: 40 token-eq
- released: ~119,960 token-eq
- result: PASS

This demonstrates the accounting distinction: cumulative traffic can greatly exceed W while simultaneous active exposure remains bounded.

## AM-1 — addressable memory scale

Materialized random memories with direct address indexes. 10,000 random queries where possible.

| Units | Correct | Peak output token-eq | Mean lookup us |
|---:|:---:|---:|---:|
| 1K | yes | 9 | 0.78 |
| 10K | yes | 9 | 0.65 |
| 100K | yes | 9 | 1.01 |
| 1M | yes | 9 | 1.56 |

Result: PASS through materialized 1M. Output/exposure remained constant. 10M+ was not materialized in this run, so no storage/index claim is made above 1M.

## AM-2 — dependency depth

Using the calibrated ~13 token-eq per raw semantic node, canonical exposure grows linearly. Hierarchical folding uses block size 64 and localized binary proof paths.

| Physical depth | Canonical token-eq | Folded cognitive levels | Folded exposure est. | Proof est. |
|---:|---:|---:|---:|---:|
| 64 | 845 | 2 | 60 | 182 |
| 512 | 6,669 | 3 | 90 | 260 |
| 1K | 13,013 | 3 | 90 | 286 |
| 10K | 130,013 | 4 | 120 | 390 |
| 100K | 1,300,013 | 4 | 120 | 468 |
| 1M | 13,000,013 | 5 | 150 | 546 |
| 10M | 130,000,013 | 5 | 150 | 650 |

Result: PASS for the synthetic composable-chain model. This does not imply arbitrary dependency graphs fold this well.

## AM-3 — branching

Depth fixed at 6 while indexed branch factor increased from 2 to 1024. A query requiring exactly one indexed branch exposed only the selected 7-node path (~91 token-eq) independent of total theoretical tree size.

Result: PASS for selective-one-branch queries.

Important failure boundary: with raw exposure ~13 token-eq/unit, only ~307 simultaneously exposed leaf units fit in W=4000. Queries requiring all irreducible leaves therefore fail the interface budget unless a valid internal aggregate/summary exists.

## AM-4 — internal computation

Materialized 1,000,000 synthetic records in 1000 categories. Built deterministic count/sum summaries and queried 10,000 random categories.

- summary build: ~0.143 s in this Python reference run
- mean summary lookup: ~1.07 us
- peak result exposure: 5 token-eq
- count conservation: PASS
- sum conservation: PASS

Result: PASS for precomputed COUNT/SUM. This is intentionally not called a selective scan: the one-time build touched the full dataset and is reported as precomputation cost.

## AM-5 — conditions/exceptions/higher-order guards

100,000 randomized guarded/exception trials were compared against the canonical oracle. Compiled evaluation preserved the same guards and exception semantics.

- matches: 100,000 / 100,000
- result: PASS for tested guarded rules

Arbitrary irreducible high-order truth tables are not proven compressible and remain a known boundary.

## AM-6 — temporal memory

Materialized 1,000,000 synthetic events (10,000 entities × 100 ordered events). Indexed range resolution used binary search rather than exposing full histories.

- mean query time: ~1.33 us in reference run
- approximate internal index probes: 14/query
- raw history exposure: not required for summarized range location
- result: PASS for indexed temporal range location

Complex temporal reasoning by an isolated model remains pending.

## AM-7 — multi-step navigation

A deterministic navigation policy traversed a random 100K-node dependency graph for 10,000 trials with depths 1..100.

- maximum tool-like calls: 101
- mean calls: 51.34
- maximum single output: 8 token-eq
- maximum cumulative exposure observed: 608 token-eq
- result: PASS for deterministic navigator

This is **not yet an LLM-agent result**. It validates the tool/protocol budget only.

## AM-8 — context recycling

With each step exposing ~100 token-eq and retaining 30 token-eq, task lengths from 10 to 100,000 steps retained constant peak active exposure:

- peak active: ~130 token-eq
- retained state: 30 token-eq
- cumulative traffic at 100K steps: 10,000,000 token-eq

Result: PASS for ideal deterministic state compaction. Real model semantic-loss accumulation is not tested here and is expected to be a major later breakpoint.

## AM-9 — negative knowledge

10,000 guaranteed-missing random addresses against a 100K-node materialized index:

- false positives: 0
- result exposure: ~3 token-eq (`UNKNOWN`)
- no exhaustive scan required
- result: PASS

## AM-10 — real-source transfer

NOT RUN / BLOCKED for rigorous evaluation in this environment.

Reason: the answering ChatGPT session can have independent access to uploaded/project source material. That violates the campaign isolation invariant (`ANSWER may depend only on EXPOSED + RETAINED`). A real-source PASS would therefore be contaminated unless the answerer is isolated from the original source.

The same issue prevents labeling AM-7 as an LLM-agent proof in this run. A clean harness must place the full memory behind Origami tools and give the answerer only the query, bounded tool outputs and retained state.

## Conclusions from R0

1. `W=4000` behaves correctly as an interface budget rather than a storage budget in the deterministic harness.
2. Direct addressing remains approximately constant in exposure through materialized 1M memory units.
3. Dependency Folding removes linear active-context growth for composable chains while retaining proof paths.
4. Selective branching scales when the relevant branch can be addressed; irreducible all-branch queries expose the real combinatorial boundary.
5. Precomputation/aggregation can make million-record answers extremely cheap to expose, but build/update cost must remain visible.
6. Context recycling can keep peak exposure constant in an ideal deterministic system even when cumulative traffic is enormous.
7. The next decisive experiment is no longer another deterministic scale test. It is an **isolated model-driven tool-navigation harness**, followed by isolated real-source transfer.

No claim is made that Origami currently provides unlimited storage, that arbitrary computations collapse to O(1), or that an LLM has already navigated these scales correctly.
