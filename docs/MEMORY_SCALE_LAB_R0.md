# Memory Scale Lab R0

Memory Scale Lab R0 measures the central Virtual Memory claim under a fixed active model-facing budget:

```text
available Origami knowledge grows
active ContextPacket budget stays fixed
```

The first deterministic ladder is:

```text
1 carrier
10 carriers
100 carriers
1,000 carriers
```

with 64 nodes per carrier and a 4,000 token-equivalent active budget by default.

## What this phase is testing

This phase does **not** ask whether a particular LLM is intelligent enough to answer the final question. It asks whether Origami itself can keep locating and unfolding the right evidence without increasing the model-facing window as the available memory grows.

The deterministic fixture gives every carrier:

- one topic root;
- one verified evidence node;
- shared-domain distractor nodes;
- a compact carrier GraphSignature;
- an addressable local graph.

All carriers share the broad `systems` domain and evidence class. The target carrier therefore has to win by the query-specific topic signal rather than by being the only carrier sharing any query term.

## Why work accounting is split

A result such as `the model saw only 700 tokens` is not enough. The system could have secretly scanned millions of nodes first.

Memory Scale Lab reports three different costs.

### Routing work

```text
signatures_scanned
candidate_carriers
selected_carriers
```

This answers how much global routing work happened before local memory was opened.

### Metadata work

```text
metadata_nodes_loaded
```

This counts nodes whose searchable metadata had to be mounted/indexed inside the selected carriers. It prevents local preprocessing from disappearing from the accounting.

### Unfolding / model-facing work

```text
nodes_touched
graph_hops
nodes_unfolded
tokens_exposed
irrelevant_tokens_exposed
```

This is the work closest to the active context window.

## Quality metrics

### Carrier Top-1 Accuracy

Did global routing put the target carrier first?

### Carrier Recall@K

Was the target carrier inside the bounded selected carrier set?

### Target Hit Accuracy

Did the resulting ContextPacket contain material from the target carrier?

### Evidence Hit Accuracy

Did it contain the target evidence node?

### Negative UNKNOWN Accuracy

A query for material that does not exist should return `UNKNOWN` while exposing zero payload tokens.

### Useful Context Efficiency (UCE)

```text
UCE = useful target-carrier tokens / all exposed tokens
```

A correct answer with 90% irrelevant context is worse than the same answer with 5% irrelevant context.

### Knowledge Scale Degradation (KSD)

```text
KSD = max(0, 1 - target_hit_accuracy(scale) / target_hit_accuracy(baseline))
```

The ideal value is zero as memory grows.

### Routing Selectivity

```text
average selected carriers / total available carriers
```

The desired direction is downward with scale while retrieval accuracy remains stable.

### Exposure Fraction

```text
average unfolded nodes / total available nodes
```

Again, this should shrink as the universe grows if selective unfolding works.

## Reproducible run

Build:

```bash
make build
```

Run the canonical experiment:

```bash
./bin/origami-memory-scale \
  -config experiments/memory-scale-r0/config.json \
  -out runs/memory-scale-r0/report.json \
  -trace runs/memory-scale-r0/trace.jsonl
```

`report.json` contains aggregate metrics per scale. `trace.jsonl` contains one record per positive/negative query with its expected carrier, ranked carriers, federation work, ContextPacket metrics and correctness flags.

That trace is intentionally model-neutral. The next model-facing phase can reuse exactly the same cases rather than inventing new questions for each model.

## Expected interpretation

A strong deterministic result is **not** merely:

```text
accuracy = 1.0
```

It must simultaneously show:

```text
FALSE_EXACT = 0
budget violations = 0
negative UNKNOWN accuracy = 1.0
stable target/evidence accuracy
bounded selected carriers
bounded metadata work after routing
a shrinking fraction of total knowledge unfolded
```

The global signature scan may still grow linearly with carrier count in R0. That is not hidden: `signatures_scanned` exposes it explicitly. If that becomes the dominant cost, the next optimization is a hierarchical/global signature index rather than pretending the cost does not exist.

## What comes after this deterministic lab

After the deterministic ladder is stable, use the same trace/cases for:

1. small local LLMs consuming ContextPackets;
2. larger baseline LLMs consuming the same packets;
3. Hybrid model-driven routing/unfold requests;
4. Native visual macro/meso/micro navigation;
5. real PDF/book corpora instead of synthetic fixtures.

A specialized Origami LLM remains later. It should be compared against these baselines using the same memory interface and the same accounting rules.
