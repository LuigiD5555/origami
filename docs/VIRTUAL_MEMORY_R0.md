# Virtual Memory R0

`origami.virtual-memory.r0` turns the Accessible Memory research direction into an executable, model-agnostic contract.

The central constraint is:

```text
available Origami knowledge  >>  active model context
active model context         <=  explicit budget (4000 token-eq by default)
```

The budget is not a storage limit and not a cumulative task limit. It is the maximum model-facing working set for one bounded access step.

## Memory unit

The addressable unit is a `Node`. Hybrid Carrier R0's original `Address`, `Value` and `Links` remain valid; richer nodes add:

- `CID`: content identity, separate from location;
- `carrier_id`, `cluster_id`, `kind`, `semantic_class`, `label`;
- parent/children plus typed weighted edges;
- importance/confidence and per-fidelity cost hints;
- provenance and verification references;
- multi-resolution content.

Canonical address form is:

```text
ohf://<carrier>/<cluster>/<node>
```

Legacy carrier-local addresses remain valid inside one carrier and are qualified when a federation is built. Multiple addresses may share one CID. Context selection deduplicates by CID before address.

## Fidelity ladder

A node may expose:

```text
label
  -> abstract
  -> summary
  -> detail
  -> evidence
  -> exact
```

A higher-level view never becomes a replacement for the lower-level evidence. The address is retained so a model/runtime can reopen the same node at greater fidelity.

If requested fidelity does not fit the active context budget, the runtime degrades fidelity before exceeding the budget. It must not silently truncate an exact claim and call it exact.

## Metadata index, not payload scan

Search indexing is deliberately restricted to:

- label;
- kind / semantic class / cluster;
- abstract;
- summary.

Legacy `Value` contributes only a bounded prefix when no richer metadata exists. `detail`, `evidence` and `exact` are not silently indexed. This is critical: a fast-looking query is invalid if it achieved speed by preprocessing/scanning the entire exact knowledge space without accounting for it.

## Query path

```text
question
  -> information-needs analysis
  -> metadata candidates
  -> declared graph traversal only when useful
  -> utility / token-cost ranking
  -> fidelity selection
  -> CID/address deduplication
  -> ContextPacket
```

Intent changes traversal. Causal questions prefer causal/dependency edges; evidence questions prefer evidence/support/source edges; contradiction questions expose contradiction edges separately; exact requests start at exact fidelity and degrade only if the active budget cannot hold it.

Missing edges/addresses return `UNKNOWN`; they never trigger an undeclared global scan.

## Universal ContextPacket

The output given to a model is independent of the model family:

```text
query
plan
items[]
relations[]
contradictions[]
uncertainties[]
provenance[]
expandable_refs[]
budget
metrics
unknown/reason
```

Therefore a model does not need visual understanding of Origami. Computational mode can decode the carrier and hand the same packet to a small local LLM, a larger LLM or a non-LLM consumer. A later Origami-specialized model may learn to request better packets, but specialization is not required by the format.

## Multi-Origami federation

Each carrier publishes a compact `GraphSignature` containing root terms, semantic classes, relation types, clusters and metadata cost. Federation follows:

```text
question
  -> rank GraphSignatures
  -> select small carrier set
  -> qualify carrier-local addresses
  -> local metadata search
  -> local graph traversal
  -> selective unfold
```

This is the first implemented form of the global-graph/local-graph distinction. The global layer decides *where to enter*; local graph addressing decides *what to unfold*.

## Working-set cache

`WorkingSet` is token bounded and LRU-like. Previously unfolded items can remain active without being fetched again; obsolete items leave the active window while their addresses remain reopenable.

## Runtime / CLI

`internal/memory` implements the generic core. `internal/lab/hybridruntime` exposes it over an image-backed Hybrid carrier.

```bash
make build

./bin/origami-memory \
  -carrier runs/.../public/carrier.png \
  -op QUERY \
  -query 'What evidence supports X?' \
  -budget 4000

./bin/origami-memory \
  -carrier runs/.../public/carrier.png \
  -op EXPAND \
  -address 'ohf://carrier/cluster/node' \
  -fidelity evidence \
  -budget 800
```

Hybrid model packets additionally permit `QUERY`, `EXPAND` and `SIGNATURE` while preserving `LOOKUP`, `FOLLOW` and `TRACE`.

## Metrics

At minimum, experiments should retain:

- candidates considered;
- nodes touched;
- graph hops;
- nodes unfolded;
- tokens exposed;
- deduplicated nodes;
- irrelevant candidates skipped;
- carrier/cluster selection accuracy;
- answer/evidence accuracy;
- Useful Context Efficiency;
- Knowledge Scale Degradation;
- `FalseExact`.

The decisive scaling experiment increases the number/size of Origami carriers while holding the active context budget fixed.

## Non-negotiable invariants

- `FALSE_EXACT = 0`.
- `UNKNOWN` is preferable to invented recovery.
- Address is not CID.
- Active context is not total memory.
- No implicit global exact scan.
- Summaries retain a path to deeper evidence.
- Model specialization is optional.
- Perception, resolution, execution and verification remain separate claims.
