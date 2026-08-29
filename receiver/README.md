# Origami Receiver Artifacts

This directory is the Origami-owned registry for receiver artifacts that are candidates for, or have passed, Origami semantic promotion.

Tlaloc may discover and distill receiver behavior, but it does not become the canonical storage authority for Origami receiver artifacts.

## Lifecycle

```text
Tlaloc swarm / Behavior Lab
        ↓
receiver candidate package
        ↓
Origami validation
        ↓
receiver/registry/<revision>/manifest.json
        + generated/MASTER_PROMPT.md
        + carrier/runtime artifacts
```

Allowed manifest states:

- `REFERENCE_CANDIDATE` — stored for reproducibility; not promoted.
- `PROMOTED` — passed the required Origami semantic, isolation and model-facing gates.
- `REJECTED` — retained only when useful as negative/provenance evidence.

## Promotion invariants

A `PROMOTED` receiver must have:

- `FALSE_EXACT = 0`;
- no contaminated trials;
- peak active model-facing token-equivalent within the declared working window;
- explicit validation evidence;
- repository-relative immutable artifact paths/hashes;
- compatibility with the declared Origami receiver contract.

The registry stores provenance; it does not replace the carrier's own BOOT/ROSETTA/PROGRAM structures.

## R0 status

`registry/r0/manifest.json` records the current universal bootstrap as `REFERENCE_CANDIDATE`. It is intentionally **not promoted** until Hybrid model-facing and cross-model symbol-permutation tests exist.
