# Origami role in Tonal — Architecture R2

## Purpose

Origami is an independent representation, carrier, addressing and virtual-memory subsystem that Tonal may use when its measured properties justify it.

Origami is not Tonal's runtime and is not required for Tonal correctness.

## Origami owns

- representation/carrier formats;
- canonical Origami profile semantics;
- addressing and selective unfolding mechanisms;
- virtual-memory behavior;
- exactness/fidelity boundaries;
- transport and roundtrip contracts;
- visual/perceptual channel experiments;
- Origami-specific promotion evidence.

## Origami does not own

- Tonal goal intake or workflow execution;
- Tonal routing/scheduling;
- Tlaloc capability promotion;
- Parrot behavior;
- Shponglese semantic meaning;
- system-level truth.

## Relationship with Shponglese

Shponglese is semantic operational IR. Origami is one possible codec/carrier/memory substrate for that IR.

```text
same Shponglese semantics
        │
        ├── JSON
        ├── compact text
        ├── binary
        └── Origami
```

Experiments should compare codecs on the same underlying semantics.

A successful Origami result should demonstrate a measurable advantage such as lower transport cost, bounded access, better selective retrieval, useful visual capacity or another preregistered metric without silently changing the task.

## Anti-prior requirement

Visual decoding tests must distinguish genuine carrier recovery from language-model completion based on priors.

Where appropriate, include controls such as:

- randomized symbol-to-meaning mappings;
- held-out arbitrary identifiers;
- permuted assignments between runs;
- semantically corrupted but mechanically valid payloads;
- exact roundtrip checks.

If performance collapses only when semantic priors are removed, report that limitation rather than treating prior-assisted completion as exact carrier capacity.

## MICRO-ISA connection

SAFE_MICRO_ISA and Context SIMD experiments are relevant to the Tonal research program because they ask which minimal operations or perceptual widths are actually reliable.

They may inform later Shponglese primitive encoding, but representation-specific observations do not automatically define Shponglese semantics.
