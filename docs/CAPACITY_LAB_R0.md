# Origami Accessible Memory Lab R0

Primary question: **how much external knowledge can Origami make accurately usable by a model while keeping the Origami↔model interface bounded to about `W = 4000 token-equivalent` at any one time?**

`W` is an interface/active-context budget. It is **not** an Origami storage limit and it is **not** a cumulative session-token limit.

The experiment therefore holds `W ~= 4000` as the principal target while memory size, structural depth, branching, temporal span and task complexity grow.

## Operational model

```text
large Origami memory
  -> SuperIndex / address resolution
  -> cheap semantic tools
  -> selective unfold / dependency folding / internal computation
  -> bounded exposure <= W
  -> model decision or intermediate state
  -> next tool call (if needed)
  -> ...
  -> answer + evidence
```

A task may use multiple tool exchanges. Old evidence may be released after the model has retained the minimal semantic state required for later steps. The important constraint is the **maximum simultaneously active Origami exposure**, not the sum of every byte/token ever exchanged during the task.

## Core hypothesis

Useful accessible memory can grow much larger than model context when Origami can:

1. address the needed region without a global scan;
2. compute/filter/fold data before exposure;
3. return compact typed results instead of raw regions;
4. preserve evidence/provenance for verification on demand;
5. retain small intermediate semantic state while releasing obsolete evidence;
6. repeat this process without exceeding the active interface budget.

## Required counters

Every trial records at least:

- TOTAL_MEMORY_UNITS / TOTAL_MEMORY_BYTES
- REQUESTED
- TOUCHED / VISITED
- UNFOLDED
- EXPOSED
- ACTIVE_BYTES_CURRENT
- ACTIVE_TOKEN_EQ_CURRENT
- PEAK_ACTIVE_BYTES
- PEAK_ACTIVE_TOKEN_EQ
- CUMULATIVE_EXPOSED_BYTES (diagnostic only; not the W gate)
- CUMULATIVE_TOKEN_EQ (diagnostic only; not the W gate)
- TOOL_CALLS
- TOOL_OUTPUT_TOKEN_EQ
- RETAINED_STATE_TOKEN_EQ
- RELEASED_TOKEN_EQ
- PHYSICAL_DEPTH
- COGNITIVE_DEPTH
- ANSWER_CORRECT
- ROUTE_INTEGRITY
- EVIDENCE_INTEGRITY
- GLOBAL_SCAN
- LATENCY

TOUCHED, UNFOLDED and EXPOSED remain distinct. Internal work is not counted as model exposure, but it must be measured so optimization cannot hide a global scan or unbounded computation.

The answering model may depend only on EXPOSED material plus explicitly RETAINED intermediate state. Access to the complete source outside those channels contaminates the trial.

## Tool surface under test

### Address/navigation

- `LOOKUP(address)`
- `FOLLOW(relation)`
- `EXPAND(address, depth/resolution)`
- `UNFOLD(address, resolution)`
- `TRACE(entity)`
- `VERIFY(evidence)`
- `BACKTRACK()`
- `STOP()`

### Cheap internal computation

The campaign must progressively test operations that reduce expensive raw exposure:

- `FILTER(predicate)`
- `PROJECT(fields)`
- `COUNT(predicate)`
- `AGGREGATE(operation, field)`
- `COMPARE(a, b)`
- `INTERSECT(a, b)`
- `RANGE(time/start/end)`
- `RESOLVE(rule)`
- `FOLD(result)`

These operations may inspect internal Origami structures, but their cost and touched regions must be observable. A small output produced by a hidden full-memory scan is not a selective-access success.

## Gates

A scale is promoted only when all applicable gates pass:

- answer correctness
- route integrity
- evidence integrity
- `PEAK_ACTIVE_TOKEN_EQ <= W`
- no hidden global scan
- deterministic execution where semantics require determinism
- UNKNOWN for genuinely absent/unresolved information

Efficiency is reported separately from correctness. A correct result that requires pathological tool-call count, exposure, latency or internal scanning is not considered an efficient success.

## Revised campaign

### AM-0 — Interface baseline

Freeze `W = 4000` as the primary interface target. Verify accounting for current/peak/cumulative exposure, retained state and release. Smaller/larger W values are diagnostic controls, not the main capacity axis.

### AM-1 — Addressable memory scale

Hold task complexity approximately constant and grow total synthetic memory:

`1K -> 10K -> 100K -> 1M -> 10M -> higher if feasible`.

Question: does lookup/navigation cost remain tied to the requested semantic region rather than total memory size?

### AM-2 — Dependency depth

Increase physical dependency depth while `W ~= 4000` remains fixed. Compare canonical traversal, hierarchical navigation and Dependency Folding.

Measure physical depth separately from cognitive depth and peak exposure.

### AM-3 — Branching and selective resolution

Increase branching factors and required branch selectivity. Determine whether Origami can resolve only relevant branches and where irreducible branch interaction forces larger active state.

### AM-4 — Internal computation vs raw exposure

Test FILTER/PROJECT/COUNT/AGGREGATE/COMPARE/INTERSECT. Construct tasks where exposing all matching records would exceed W but an internally computed, verifiable result is small.

Success requires compact correct output plus evidence reference without hidden global-scan cheating.

### AM-5 — Higher-order interactions

Answers depend on multiple relations, conditions and exceptions. Determine which interactions can be compiled/folded and which are irreducible.

### AM-6 — Temporal memory

Increase event history and temporal span. Test RANGE, state-at-time, transition tracing and temporal conditions without exposing the complete history.

### AM-7 — Multi-turn tool navigation

The model receives only the query, bounded tool outputs and retained semantic state. It must decide which Origami operation to call next.

Measure tool-call count, peak active exposure, cumulative exposure, retained-state size, backtracking and answer/evidence correctness.

### AM-8 — Context recycling

Explicitly test:

`EXTRACT -> COMPUTE/FOLD -> RETAIN minimal state -> RELEASE old evidence -> NEXT EXTRACT`.

Grow task length until retained state, navigation error or tool overhead becomes the bottleneck. This is the main test of Origami as external cognitive memory rather than one-shot retrieval.

### AM-9 — Negative knowledge

Absent facts, missing addresses and unresolved guards must produce `UNKNOWN` without exhaustive search or hallucinated completion.

### AM-10 — Real-source transfer

Only after synthetic gates pass: freeze a real source representation, then query it through the same bounded interface. The answerer must not have independent access to the complete source. PDF/image carrier experiments come after semantic accessible-memory behavior is demonstrated.

## Scale strategy

Synthetic memory size is not capped by the campaign definition. `1K, 10K, 100K, 1M, 10M` are checkpoints, not an intended maximum.

Promote until one of these occurs:

1. correctness breaks;
2. evidence/route integrity breaks;
3. peak active exposure exceeds W;
4. tool-call/retained-state growth becomes impractical;
5. internal access degenerates toward global scanning;
6. storage/index/precomputation/update cost becomes the dominant bottleneck.

The failure point is a result, not something to conceal by increasing W.

## Dependency Folding accounting

A folded result does not erase physical depth or internal computation. Reports carry physical depth, cognitive depth, touched regions, proof references and peak model exposure.

Conditions, exceptions and verification material actually exposed to the answerer count against W. If folding cannot preserve semantics across a boundary, refine only that region.

## Primary report

The campaign report must answer:

- How large was Origami memory?
- How much of it was touched internally?
- How much was unfolded?
- What was the largest simultaneous exposure to the model?
- How many tool calls were required?
- How much semantic state had to be retained between calls?
- Was the final answer correct?
- Can its evidence be reopened and verified?
- Where did cost grow: storage, indexing, navigation, internal computation, exposure, retained state, model reasoning, or verification?

The headline metric is **accessible external memory under a bounded model interface**, not "how much information fits inside 4000 tokens".
