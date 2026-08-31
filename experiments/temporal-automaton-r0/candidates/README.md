# Temporal carrier experimental candidates

These files are experimental build inputs for `origami-candidate-build`. They do not change the canonical `origami.temporal-carrier.r0.profile-1` renderer.

## t2-temporal-grammar-visible-r1

Evidence source: manual external real-VLM trial recorded by Tlaloc against DeepSeek using the `signal-chain-r0` carrier.

Observed frontier:

```text
TEMPORAL_ROUTE
  -> TEMPORAL_RULE_AMBIGUOUS
  -> T_TEMPORAL
  -> TEMPORAL_GRAMMAR
```

DeepSeek recovered BOOT/ROSETTA, the cells and initial states, but did not reliably recover the declared causal rules or the synchronous final-state execution. The exact DeepSeek model variant was not captured and remains unspecified in the evidence record.

The candidate applies exactly one representation mutation:

```text
kind   = TEMPORAL_STRUCTURE
target = T2_SEMANTIC_TEMPORAL_SUPERGRAPH
value  = VISIBLE_RULE_MICROGRAMMAR_R1
```

For one-row temporal carriers, the mutation uses the empty band between the T2 cell graph and the timeline to expose:

- that rules, rather than graph topology alone, drive state changes;
- that all rules in a step are evaluated against the same pre-step snapshot;
- each visible rule's requirement, target cell, from-state and to-state.

Dense carriers with more than four cells fail over to a bounded marker plus the existing checkpoint emphasis rather than overwriting the second T2 cell row.

Hard invariant:

```text
candidate TemporalProgram SHA-256 == parent TemporalProgram SHA-256
```

A successful build is not evidence of better VLM understanding. The candidate must be tested against the same blind prompt/model conditions and then evaluated by Tlaloc before any promotion can be considered.
