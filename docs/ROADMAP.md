# Origami roadmap

This roadmap separates **core Origami language/runtime work** from the **OHF carrier/laboratory research track** and restores generative Fold/Unfold as a prerequisite for document/carrier demonstrations.

## Recovery sequence — Core Origami 6.x

Current priorities, in dependency order:

1. **Generative representation contract** — preserve `docs/FORMAL_CORE_R0.md` while making `FOLD`, `RELATE`, `EVOLVE`, `OBSERVE` and `UNFOLD` explicit core responsibilities. Enforce `RASTER_IS_NOT_FOLD`.
2. **Reference machine contracts** — executable finite state/context/rule/operation/observation result types with deterministic default semantics and distinct `PASS`, `FAIL`, `INVALID_CONTRACT`, `UNSUPPORTED` outcomes.
3. **Generative IR R0** — implement and test a minimal explicit IR including `LITERAL`, `REF`, `CONCAT`, `REPEAT`, `TRANSFORM`, `RULE`, `GRAPH`, `RESIDUAL`, `VERIFY` or the smallest justified subset that preserves these responsibilities.
4. **Exact deterministic Fold/Unfold fixture** — begin with tiny text/pattern fixtures. Demonstrate source hash == exact-unfold hash when exact recovery is declared. Do not use PDF or VLM as the first proof.
5. **Relations and addressing** — stable node addresses, dependency graph and deterministic relation queries.
6. **Selective unfolding** — resolve an address/query to a bounded dependency closure and measure nodes/bytes materialized instead of reconstructing everything by default.
7. **Representation tournament instrumentation** — compare literal, reference, repeat, graph, rule, transform and residual candidates using measured compactness, generative leverage, addressability, verification and unfold cost. Do not invent scoring weights before evidence exists.
8. **Machine evolution** — connect the generative representation to deterministic state transitions and test fixed points, cycles, constraints and bounded dynamics.
9. **Observation runtime** — execute finite Observation Contracts while preserving `LATENT_IS_NOT_ABSENT` without making latent claims unfalsifiable.
10. **Projection boundary** — support multiple projections of the same semantic state and prove that visual projection is downstream of representation.
11. **Raster baseline classification** — existing page-render/screenshot paths remain valid perception/control baselines but cannot emit a successful generative-Fold claim.
12. **Document IR** — only after the preceding gates, introduce canonical structured-document ingestion with provenance.
13. **PDF ingestion** — parse PDF into Document IR, then Fold that structure into Origami. PDF is an input format, not an Origami representation.
14. **Document query experiment** — compare RAW, RASTER and ORIGAMI paths on correctness, provenance, exact recovery where applicable, representation size, generative leverage, materialized bytes/query, nodes unfolded/query and false-exact behavior.
15. **Perceptual-channel execution** — integrate coherent-state and perceptual-channel contracts into the executable reference machine without changing declared semantics; add deterministic fixtures and implement `MOIRE`, `PHASE_SHIFT`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE`, `TEMPORAL_DECAY` incrementally with explicit unsupported-operation behavior.
16. expose stable versioned contracts that Tlaloc or other consumers can evaluate against.

Throughout the sequence, keep state semantics separate from perceptual availability and do not claim physical quantum computation.

## OHF research track

Continue the executable R3.10-LAB workbench as an Origami profile/carrier research line, but place carrier experiments downstream of an established Origami representation when the experiment claims Origami Fold:

- DVR / PSS evidence and orthogonality gates;
- Glyph Calculus;
- SAFE_MICRO_ISA;
- Context SIMD;
- Macro-Gestalt experiments;
- representation tournament / graph / grammar / transform folding;
- SuperIndex and selective unfolding;
- model-specific safe perceptual profiles;
- conservative Native/Hybrid evaluation;
- exact historical R3.10 SDK recovery when the authentic bytes become available.

OHF-specific `<=500 KB`, `FALSE_EXACT=0` and Native-boundary invariants remain track requirements unless explicitly revised by an OHF change record. They are not automatically universal constraints on every future Origami representation profile.

## Cross-project integration

Tlaloc integration should consume versioned Origami contracts. Tlaloc may train/compile behavior for a model to obey Origami semantics, but Origami remains independently versioned and authoritative for representation laws.

Do not migrate Tlaloc's current Origami-specific reference implementation until the recovered Origami reference API has executable fixtures and a stable contract. Adapter first, migration second, duplicate removal last.

## Promotion rule

Do not promote a generative mechanism, perceptual dimension, compound, operation or profile from plausibility or one attractive example. Promotion requires deterministic validity where applicable, declared failure semantics, controlled-mutation evidence, relevant perception/evaluation evidence and regression coverage for the declared scope.

A screenshot/page render can pass raster/perception gates while still failing or not participating in the generative-Fold gate.
