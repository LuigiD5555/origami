# Origami Semantic Graph Automaton — R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTATION`

## Purpose

Origami can represent behavior, not only static knowledge. R0 models a bounded behavior as a graph of local semantic transition units.

A cell is intentionally more general than a classical cellular-automaton cell: its neighborhood is an addressable graph rather than a mandatory regular lattice.

```text
Automaton A = (C, G, R, S0)

C  cells
G  directed semantic graph
R  declared local transition rules
S0 initial state
```

The reference transition is:

```text
S_(t+1) = F(S_t, ACTIVE_RULES_t, G, R)
```

## Cell

A cell declares:

```text
id
kind
initial_state
neighbors
```

A Tlaloc Tlaloque may be distilled into one or more cells, but Tlaloc is not an Origami runtime dependency.

## Rule

A rule declares:

```text
id
target_cell
from_state
to_state
requires[]
priority
```

R0 uses synchronous updates. Every rule is evaluated against the same input state. Changes become visible together in the next state.

When multiple applicable rules target one cell:

```text
higher priority wins
then lexicographically smaller rule id
```

This is a deterministic rule, not a model judgment.

## Graph

Edges may describe `NEIGHBOR`, `DEPENDENCY`, `EMISSION`, `CAUSAL` or `ROUTE` relations. The runtime does not infer undeclared edges.

## Why this is not a swarm runtime

The development process may use many agents. The portable Origami artifact carries the distilled state/rules/graph required to reproduce the behavior. The artifact does not require the original swarm trace or Tlaloc installation.

## Reference implementation

```text
internal/automaton
cmd/origami-automaton
```

Example:

```bash
origami-automaton -mode inspect -in automaton.json
origami-automaton -mode step -in automaton.json
```

## Invariants

```text
LOCAL_RULES_ARE_DECLARED
SYNCHRONOUS_STEP_IS_DETERMINISTIC
RULE_CONFLICT_RESOLUTION_IS_DETERMINISTIC
GRAPH_ADDRESSABILITY_REQUIRED
UNKNOWN > INVENTION
FALSE_EXACT = 0
```
