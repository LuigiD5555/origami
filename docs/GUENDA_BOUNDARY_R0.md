# Guenda ↔ Origami boundary R0

Guenda may use Origami as structured IR / protocol / external memory in Parrot experiments, but Origami remains independently versioned and owns its own semantics.

## Roles

```text
GUENDA
  experimental R&D system
  discovers and evaluates Parrot candidates
  compares learned-inside vs externalized capacity

ORIGAMI
  representation / protocol / memory target
  owns semantic contracts, carriers, ROSETTA, profiles and releases
```

## Integration direction

```text
human input
  -> Parrot candidate under test
  -> structured meaning
  -> optional Origami encoding / memory interaction
  -> result
  -> Guenda evaluator
```

Guenda may benchmark whether Origami improves quality, memory efficiency, robustness or externalization. Such evidence is an experimental result only; it does not redefine Origami.

## R0 experiment discipline

Origami is not required for Guenda P0/bootstrap. Add it only when the experimental protocol reaches the Origami phases and compare it against an explicit baseline.

Guenda must record the resource tradeoff, not only task quality:

```text
quality
resident model bytes
active compute
state / memory bytes
external calls
latency
```

This allows Guenda to test whether a capability should remain neural or be externalized through Origami.

## Ownership rules

```text
ORIGAMI_OWNS_ORIGAMI_SEMANTICS
GUENDA_EVIDENCE_NE_ORIGAMI_PROMOTION
GUENDA_MAY_PROPOSE_NOT_REDEFINE
ORIGAMI_NE_REQUIRED_FOR_GUENDA_BOOTSTRAP
EXTERNALIZED_CAPACITY_MUST_BE_ACCOUNTED
PUBLIC_CONTRACT_OVER_PRIVATE_IMPORTS
```

A future Parrot artifact may declare Origami capabilities or compatibility, but canonical Origami contracts continue to live in the Origami repository.
