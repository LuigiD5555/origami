# Origami project instructions — Architecture R2

## Role

Origami is an independently testable representation, transport, addressing and virtual-memory substrate.

Within the Tonal ecosystem, Origami may carry, index, address or selectively unfold semantic structures such as Shponglese programs. Origami does not own Tonal runtime decisions and does not own Shponglese semantics.

Origami must remain independently measurable: Tonal integration is a use case, not a shortcut around Origami's own evidence requirements.

## Shared North Star

The Tonal ecosystem investigates whether complex, reliable behavior can emerge from composition of small, bounded, verifiable and reusable capabilities instead of requiring one increasingly complex general-purpose model.

Origami's contribution is representation efficiency, bounded access, transport, addressability and selective unfolding while preserving explicit semantic/exactness boundaries.

## Canonical terms

- **Tonal**: complete heterogeneous runtime and research system.
- **Tlaloc**: capability foundry and Behavior Lab.
- **Tlaloque**: bounded typed measurable capability.
- **Parrot**: one probabilistic Tlaloque; no system-level authority.
- **Shponglese**: semantic operational IR whose meaning must be codec-independent.
- **Origami**: representation/carrier/memory substrate that may encode Shponglese or other semantic structures.

## Shponglese boundary

Shponglese answers: **what does the operational structure mean?**

Origami answers: **how can that structure be represented, transported, addressed or selectively unfolded?**

Do not bind Shponglese semantics to a specific visual primitive, image profile or Origami codec.

The same Shponglese program should remain comparable across text, JSON, binary and Origami representations.

## Evidence discipline

Preserve existing Origami safety principles including:

- `FALSE_EXACT = 0`;
- PRESENT / ABSENT / UNKNOWN distinctions;
- source-plane vs semantic-plane separation;
- perception != execution != verification;
- nominal capacity != demonstrated safe capacity;
- bounded/selective access rather than undeclared global scans;
- explicit provenance and roundtrip evidence.

Do not claim visual/carrier success when a model could succeed from linguistic priors alone. Future carrier experiments should include anti-prior controls using randomized, permuted or semantically corrupted symbol assignments where appropriate.

## MICRO-ISA direction

Existing SAFE_MICRO_ISA / Context SIMD work remains active research, not historical debris. It should be interpreted as evidence about minimal reliable operations and perceptual capacity, and may later support Shponglese/Origami codec experiments.

Do not promote a primitive or width merely because nominal codebook capacity permits it. Follow measured false-known/exactness gates.

## Document authority

Read in this order:

1. `CLAUDE.md`
2. `README.md`
3. `docs/CURRENT_STATE.md`
4. `docs/ARCHITECTURE.md`
5. `docs/ROLE_IN_TONAL.md`
6. active experiment specification

Anything under `docs/archive/` is historical and not current architectural authority.

If archived documentation conflicts with Architecture R2 current documentation, current documentation wins. Frozen experiment artifacts remain authoritative for their experiment.

## Change discipline

- Do not make Origami mandatory for Tonal correctness.
- Do not grant Origami authority over Shponglese semantics.
- Do not use model priors as proof of carrier fidelity.
- Keep deterministic and native-model evidence clearly separated.
- Preserve portable/standalone Origami testing where meaningful.
- Archive superseded architecture instead of deleting project history.
