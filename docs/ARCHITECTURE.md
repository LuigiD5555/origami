# Origami architecture

Origami is a visual/computational representation and state-machine language. Its architecture separates semantic state, perceptual availability, deterministic dynamics, receiver bootstrap and profile/carrier implementations.

```text
STATE / RELATIONS
      ↓
ORIGAMI SEMANTICS
  ├── coherent-state profile
  ├── Fold / Unfold / addressability
  └── perceptual-channel contracts
      ↓
MACHINE / DYNAMICS
      ↓
SELF-BOOT RECEIVER CONTRACT
  BOOT -> ROSETTA -> PROGRAM
      ↓
INDEX / MEMORY / VERIFICATION
      ↓
VISUAL + COMPUTATIONAL PROJECTION
      ↓
PROFILE OR CARRIER
  └── OHF research track
```

## Semantic layers

Coherent-state kinds (`determinate`, `superposed`, `coupled`, `observed`) are distinct from perceptual availability conditions (`latent`, `motion-bound`, `phase-bound`, `stereo-bound`, `revealed`, `decayed`).

A representation may therefore have a valid semantic state while its intended percept is available only through interference, depth, motion or temporal integration. Meaning may also belong to the trajectory of a machine rather than to one static frame.

## Self-boot receiver contract

`origami.hybrid-receiver.r0` defines a carrier-facing bootstrap structure. It does not define every Origami semantic operation; it defines how a receiver can discover the semantics needed for one carrier without prior knowledge of that carrier's private symbol assignments.

Required sections are:

```text
BOOT
  ↓
ROSETTA
  ↓
PROGRAM
  ↓
INDEX
  ↓
MEMORY
  ↓
VERIFICATION
```

`RESIDUAL` is optional and must remain explicit/measured.

### BOOT

Declares entry point, carrier version, start state and references to the remaining receiver sections. BOOT exists so the model/runtime knows where to begin rather than relying on a globally memorized layout.

### ROSETTA

Maps physical carrier-local marks to semantic tokens. Physical glyphs have no universal meaning merely because a previous carrier used them. Symbol assignments may be permuted between carriers while preserving semantics.

### PROGRAM

Contains bounded deterministic micro-agent/automata rules. Individual rules should be locally simple. Complex behavior may emerge from their composition and temporal evolution.

### INDEX / MEMORY

Supports selective navigation over memory that can be much larger than the active model-facing interface. The working-window target is an interface constraint, not a storage-capacity claim.

### VERIFICATION

Preserves exactness boundaries, proof references and the ability to reopen canonical evidence. `FALSE_EXACT = 0` remains mandatory.

## Hybrid receiver path

Hybrid is the primary target mode because it keeps model perception where it is useful while moving mechanical work into deterministic Origami execution.

```text
carrier.png
    ↓
model: find BOOT / interpret local Rosetta / choose semantic region
    ↓
Origami runtime: micro-agents + addressing + Fold/Unfold + computation
    ↓
bounded result / next region / proof ref
    ↓
model: integrate retained state and decide next access
    ↓
verification / answer
```

This is distinct from two diagnostic baselines:

- **Native:** Master Prompt + one carrier image + question; no external decoder/runtime.
- **Computational:** Master Prompt + declared Origami tools + question; no visual bootstrap requirement.

A Hybrid PASS requires both sides to cooperate through the declared contract. A correct answer obtained by giving the model the hidden source is contaminated, not successful.

## Accessible memory

The central capacity target is not "fit all knowledge inside 4000 tokens". It is:

```text
large Origami memory
  -> selective navigation / computation / folding
  -> <= ~4000 token-equivalent active model-facing context
  -> retained compact semantic state
  -> release obsolete evidence
  -> next bounded access
```

Cumulative information used during a task may exceed the active-window target. Internal work and model-facing exposure must be measured separately so a compact answer cannot hide an uncontrolled global scan.

## Perceptual channels

The current R0 contract defines spatial, interference, depth, temporal and emergent channels. Emergent meaning may belong to an interaction rather than a single source layer. A Temporal Latent Image may require a declared trajectory to become perceptually available.

The contract defines meaning and validity conditions; it does not imply every operation already has a production renderer/detector.

## OHF research profile

OHF remains the main executable carrier/laboratory track currently in this repository. It explores compact addressable carriers, self-bootstrapping visual structures, selective unfolding, deterministic verification and bounded model perception.

Representative direction:

```text
source
  ↓
canonical semantic representation
  ↓
graph / grammar / transforms / Dependency Folding
  ↓
SuperIndex + receiver sections
  BOOT + ROSETTA + PROGRAM + INDEX + MEMORY + VERIFY
  ↓
visual compiler / perceptual ABI
  ↓
OHF carrier
  ↓
Hybrid receiver path
```

DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA and Context SIMD are experimental OHF mechanisms. They are useful building blocks for the receiver path but do not automatically become universal Origami language semantics.

## Tlaloc boundary

Tlaloc owns work orchestration, behavior compilation, Tlaloque coordination, model-facing experimentation and the search/distillation process that can discover better receiver behavior.

Tlaloc may produce a **candidate package** containing a universal receiver prompt candidate, BOOT/Rosetta strategy and micro-agent IR derived from swarm behavior. Origami remains authoritative for its semantics, validates the candidate against Origami contracts and stores/exports only promoted receiver artifacts.

```text
Tlaloc swarm / search
      ↓
receiver candidate
      ↓
Origami semantic validation
      ↓
promoted prompt + carrier machinery
```

Tlaloc may optimize *how* a model operates Origami; it must not redefine *what Origami means*.

## Verification discipline

For representation/receiver changes, record before/after state, evidence, regressions, downstream impact and promotion decision. Mathematical capacity must remain separate from evidence-gated perceptual capacity. Unsupported behavior must be explicit rather than guessed.
