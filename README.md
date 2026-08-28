# Origami 6.0.0-alpha.1

Origami is an experimental **visual/computational representation and state-machine language** for expressing complex states, relations and evolution.

Origami is **not** the name of the complete work system. The broader orchestration, behavior compilation, Tlaloque coordination and model-facing workflow now belong to **Tlaloc**.

## Project boundary

```text
TLALOC                          ORIGAMI
work system                     representation
orchestration                   state language
behavior compilation   ----->   semantics contract
Tlaloque coordination  ----->   reference semantics engine/API
verification                     dynamics / projection
```

Origami may be used independently of Tlaloc. Tlaloc may also operate without Origami or with another representation provider.

> Tlaloc coordinates work; Tlaloque perform bounded specialist work; Origami defines how a class of states is represented and transformed.

## Current semantic profile

The 6.0 line formalizes a quantum-inspired coherent-state profile as one representation mode. It is a computational analogy, **not a claim of physical quantum computation** and not the entire scope of Origami.

Current state kinds:

- `determinate` — one active branch;
- `superposed` — multiple coherent branches with complex amplitudes;
- `coupled` — a joint state whose members must not be evaluated independently;
- `observed` — an explicitly resolved state.

Core laws include:

1. `TRANSFORM` evolves a valid state but cannot select a branch.
2. `OBSERVE` is an explicit resolution boundary.
3. `INTERFERE` combines complex amplitudes; equal/opposite paths cancel.
4. Zero amplitude means cancellation, not unknown.
5. Absence, unknown and inhibited are distinct semantic conditions.
6. A coupled state remains a joint object until an explicit decomposition operation is applied.

See [`docs/STATE_SEMANTICS_R0.md`](docs/STATE_SEMANTICS_R0.md).

## Fold / Unfold

In the current semantic contract:

- `UNFOLD` constructs or expands alternatives under a declared representation rule.
- `FOLD` constrains and/or explicitly resolves according to a declared policy.

These are Origami operations. Model training, prompt compilation and behavior enforcement are Tlaloc responsibilities.

## Historical continuity

Origami keeps the existing 6.x numbering. `6.0.0-alpha.1` is the representation-side branch point created when the former unified project separated into **Tlaloc + Origami**.

The repository previously described **Origami HyperFold (OHF) R2** as the complete project. Those materials are retained as historical research inputs:

- `docs/SCIENTIFIC_FOUNDATIONS.md`
- `prompts/OHF_R2_MASTER_PROMPT.txt`

They document an earlier carrier/compression direction and **do not define the current project boundary**. See [`docs/LEGACY_OHF_R2.md`](docs/LEGACY_OHF_R2.md).

## Repository status

This repository currently contains the semantic contract and project-boundary documentation for Origami 6.0.0-alpha.1.

The exact pre-split Origami 5.1 implementation remains the authoritative implementation base until the 6.0 semantic overlay is merged and regression-tested against it. Therefore this release remains **experimental**.

## Files

```text
origami/
├── README.md
├── VERSION
├── CHANGELOG.md
├── PROJECT_BOUNDARY.md
├── CHANGE_CONTROL_6.0.0-alpha.1.md
├── docs/
│   ├── STATE_SEMANTICS_R0.md
│   ├── LEGACY_OHF_R2.md
│   └── SCIENTIFIC_FOUNDATIONS.md      # historical OHF material
└── prompts/
    └── OHF_R2_MASTER_PROMPT.txt        # historical OHF R2 artifact
```

## Version

`6.0.0-alpha.1`
