# Origami project boundary

```text
TLALOC                                      ORIGAMI
work/orchestration system                   representation/state-machine language
swarm/Tlaloque search               ----->  receiver candidate inputs
behavior compilation + distillation ----->  prompt/boot/micro-agent candidates
model-facing evaluation             <-----  semantic/reference contracts
promotion recommendation            ----->  candidate evidence
                                             |
                                             +-- semantic validation
                                             +-- carrier execution
                                             +-- promoted artifact storage
```

Origami may be used independently of Tlaloc. Tlaloc may operate without Origami or consume Origami as one representation provider.

The stable rule is:

> Tlaloc discovers, coordinates and distills behavior; Origami defines and executes representation semantics, including how a carrier bootstraps, maps local symbols, evolves state and exposes/verifies memory.

Tlaloc may optimize **how** a target model is instructed to bootstrap and operate Origami, and may distill successful swarm behavior into candidate micro-agent rules. It must not redefine Origami state semantics, symbol meaning for a specific carrier, transition authority, exactness rules or verification claims.

Origami owns:

- semantic state/relations/dynamics;
- Fold/Unfold and addressability contracts;
- self-boot receiver contract (`BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION`);
- carrier-local symbol bindings;
- deterministic micro-agent execution semantics;
- carrier generation/projection;
- promoted receiver artifacts and their provenance;
- evidence/exactness boundaries.

Tlaloc owns:

- swarm/Tlaloque exploration;
- behavior-compilation lifecycle;
- candidate prompt search/mutation;
- candidate BOOT/Rosetta strategy search;
- distillation of complex swarm traces into simple bounded micro-agent candidate rules;
- cross-model evaluation/tournaments;
- promotion recommendations and regression coordination.

Within Origami, **OHF (Origami HyperFold)** is a carrier/protocol and executable research track. R3.10-LAB, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA and Context SIMD belong to that track unless and until a capability is deliberately generalized into the core Origami language.

The preferred end-to-end receiver is Hybrid: model perception handles bootstrap/high-level navigation and deterministic Origami execution handles mechanical local work, addressing, computation and verification. Native and Computational remain diagnostic modes.

Origami owns its semantic contracts. Tlaloc may build fixtures, evaluators and receiver candidates around those contracts, but Origami is the authority that decides whether a candidate is semantically valid and stores it only after promotion.
