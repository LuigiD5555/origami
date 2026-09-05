# Origami current state — Architecture R2

**Status date:** 2026-09-05  
**Repository version:** 6.0.0-alpha.15  
**Status:** experimental

## Current role

Origami is an independently testable representation, transport, addressing and virtual-memory substrate.

Within the Tonal ecosystem, Origami may carry, index, address or selectively unfold semantic structures such as Shponglese programs. Tonal is the complete heterogeneous runtime/research system; Tlaloc is the capability foundry and Behavior Lab. Origami does not own Tonal runtime decisions or Shponglese semantic meaning.

Origami remains independently usable. Tonal integration is a use case, not a replacement for Origami-specific evidence.

## Evidence authority

The machine-readable current capability claims are in:

```text
state/CLAIMS.json
```

Those claims distinguish `implemented`, `designed` and evidence-backed boundaries. Architecture R2 does not upgrade any claim merely because the ecosystem role changed.

The detailed pre-R2 `CURRENT_STATE.md` snapshot, including its generated claims table and alpha-era narrative, is preserved under:

```text
docs/archive/superseded-architecture/CURRENT_STATE_PRE_R2.md
```

Use it for provenance, not current ecosystem authority.

## Technical capabilities retained

Current implementation/evidence includes work on:

- deterministic semantic state transitions;
- selective unfolding and bounded virtual memory;
- exact-source identity and deterministic evidence reduction;
- fixed/self-contained carriers;
- semantic codec negotiation;
- temporal programs and causal traces;
- candidate build and evidence-gated promotion;
- reference semantics;
- receiver/import boundaries;
- experimental perceptual/MICRO-ISA tooling.

Refer to `state/CLAIMS.json` for the exact status of each claim and its evidence references.

## SAFE_MICRO_ISA / Context SIMD

SAFE_MICRO_ISA and Context SIMD remain active research, now surfaced under `docs/research/`.

They are relevant to the broader Tonal Primitive Swarm question because they measure which minimal operations or perceptual widths are actually safe. They do **not** define the universal Shponglese primitive set automatically.

Key discipline remains:

```text
false-known / false-exact errors are hard failures
nominal capacity != safe measured capacity
start narrow
widen only demonstrated winners
measure cost/tokens/latency
promote only with evidence
```

## Shponglese relationship

Shponglese is semantic operational IR. Origami is one possible representation/carrier/memory substrate for that IR.

A fair experiment holds semantics constant:

```text
same Shponglese program
  ├── JSON
  ├── compact text
  ├── binary
  └── Origami
```

and compares recovery, execution, transport size, latency, selective access and false-known behavior.

## Anti-prior requirement

Future compact visual carrier claims must distinguish real recovery from language-model completion based on priors. Randomized identifiers, permuted symbol mappings and exact downstream execution are planned controls.

See:

- `docs/ROLE_IN_TONAL.md`
- `docs/SHPONGLESE_CARRIER.md`
- `docs/research/ANTI_PRIOR_TESTS.md`
- `docs/research/CODEC_BASELINES.md`
- `docs/research/MICRO_ISA_DIRECTION_R2.md`

## Current R2 migration work

Completed in the R2 branch so far:

- current ecosystem role documented;
- pre-R2 project boundary removed from root and archived;
- root README aligned to repository VERSION and R2 role;
- change-control history moved out of root;
- MICRO-ISA/Glyph Calculus material moved into active `docs/research/`;
- Shponglese carrier and anti-prior research boundaries documented.

No existing empirical claim has been promoted by these documentation changes.
