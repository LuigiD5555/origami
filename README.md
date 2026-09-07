# Origami 6.0.0-alpha.15

Origami is an experimental **representation, transport, addressing and virtual-memory substrate** for structured semantic state and selectively accessible memory.

Origami remains independently testable and independently usable. Within the Tonal ecosystem it may carry, index, address or selectively unfold structures such as Shponglese programs, but it is not required for Tonal correctness and it does not own Shponglese semantics.

## Architecture R2 role

```text
Shponglese
  semantic operational IR
  what the structure means
          │
          ├── JSON / text / binary
          └── Origami
                representation / carrier / memory
                how the structure is transported,
                addressed or selectively unfolded
```

Tonal is the complete heterogeneous runtime/research system. Tlaloc is the capability foundry and Behavior Lab. Origami supplies an optional representation/memory substrate whose usefulness must be demonstrated independently.

## Core evidence boundaries

Origami preserves the experimental discipline established by earlier alpha work:

```text
FALSE_EXACT = 0
PRESENT != ABSENT != UNKNOWN
PERCEPTION != EXECUTION != VERIFICATION
SOURCE PLANE != SEMANTIC PLANE
NOMINAL CAPACITY != DEMONSTRATED SAFE CAPACITY
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
```

A deterministic runtime roundtrip does not prove native multimodal-model interpretation. A model's plausible reconstruction does not prove exact carrier recovery.

## Current technical areas

The repository contains active work on:

- semantic state-machine execution;
- canonical semantic graph/trajectory access;
- bounded virtual memory and selective unfolding;
- exact-source identity and evidence reduction;
- fixed self-contained carriers;
- codec negotiation and semantic roundtrip;
- temporal programs and causal traces;
- perceptual-channel experiments;
- SAFE_MICRO_ISA / Context SIMD;
- candidate profile construction and promotion gates.

`docs/CURRENT_STATE.md` and the capability-status/evidence documents remain the place to determine which claims are implemented, merely designed, or still awaiting real-model evidence.

## SAFE_MICRO_ISA / Context SIMD

The existing MICRO-ISA work remains active research under Architecture R2.

It asks which minimal operations and perceptual widths are actually reliable, using gates such as:

- begin with narrow probes;
- widen only demonstrated winners;
- stop on false-known/exactness failure;
- measure retries, tokens, latency and cost;
- distinguish nominal capacity from safe measured capacity;
- avoid promotion without appropriate repeated evidence.

Under R2 this work may inform later Shponglese carrier experiments, but an Origami visual primitive does not automatically become a universal Shponglese semantic primitive.

See `docs/research/MICRO_ISA_DIRECTION_R2.md`.

## Anti-prior carrier testing

Future Shponglese/Origami carrier comparisons must distinguish genuine recovery from language-model completion based on priors.

Useful controls include randomized symbol-to-meaning mappings, arbitrary held-out identifiers, permuted assignments and exact downstream execution from the recovered semantic program.

See:

- `docs/SHPONGLESE_CARRIER.md`
- `docs/research/ANTI_PRIOR_TESTS.md`
- `docs/research/CODEC_BASELINES.md`

## Fair codec comparisons

Origami should compete against conventional representations while holding underlying semantics constant:

```text
same semantic program
  ├── JSON
  ├── compact text
  ├── compact binary
  └── Origami
```

Measure exact semantic recovery, execution success, bytes/tokens, latency, selective-access behavior and false-known errors rather than visual novelty alone.

## Standalone boundary

Tonal integration does not transfer semantic authority to Tonal, and Origami integration does not make Origami a hidden requirement of Tonal.

Origami continues to own its own representation/profile semantics and Origami-specific promotion evidence.

## Documentation authority

Start with:

1. `CLAUDE.md`
2. this `README.md`
3. `docs/CURRENT_STATE.md`
4. `docs/ROLE_IN_TONAL.md`
5. current architecture/evidence documents
6. the active experiment specification

Anything under `docs/archive/` is historical and does not override current Architecture R2 documentation.

## Development checks

Use the repository's current Makefile/Go test workflow and the experiment-specific gates declared by the active protocol before promoting claims.

> Origami should earn its place in Tonal through measured representation, transport, addressing or memory advantages—not because the architecture assumes it must be useful.
