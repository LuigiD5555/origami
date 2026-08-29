# Current Origami state

**Status date:** 2026-08-29  
**Project version:** 6.0.0-alpha.3  
**Status:** experimental

## Root project state

Origami is the visual/computational representation and state-machine language. Tlaloc is a separate optional work/orchestration system.

Current semantic/research layers:

| Area | State |
|---|---|
| coherent-state semantics | EXPERIMENTAL_CONTRACT — alpha.1 |
| perceptual channels | EXPERIMENTAL_CONTRACT — alpha.2 |
| Temporal Latent Image semantics | EXPERIMENTAL_CONTRACT |
| complete perceptual runtime for all alpha.2 operations | NOT IMPLEMENTED |
| Dependency Folding R0 | EXPERIMENTAL / reference implementation |
| Accessible Memory Lab R0 | EXPERIMENTAL / synthetic evidence |
| self-boot Hybrid Receiver R0 | EXPERIMENTAL_NOT_PROMOTED |
| receiver artifact registry | IMPLEMENTED / reference candidate only |
| Hybrid public isolation bundle | IMPLEMENTED |
| OHF R3.10-LAB research track | EXPERIMENTAL / executable |

The root machine-readable source is `state/ORIGAMI_STATE.json`. OHF-specific component state remains under `state/components/`.

## Self-boot Hybrid Receiver R0

The experimental contract `origami.hybrid-receiver.r0` defines the carrier-facing bootstrap path:

```text
BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION
```

`RESIDUAL` remains optional and explicitly measured when used.

The intended division of labor is Hybrid:

```text
model perception
  -> find BOOT
  -> interpret carrier-local ROSETTA
  -> choose high-level semantic region
Origami deterministic runtime
  -> execute micro-agent rules
  -> address / Fold / Unfold / compute / verify
model
  -> retain compact semantic state
  -> choose next bounded access or answer
```

Native and Computational are diagnostic baselines, not separate final architectures.

Carrier-local physical symbols have no assumed global semantic meaning. The mandatory anti-memorization direction is to keep the external Master Prompt fixed while permuting equivalent carriers' physical symbol assignments and requiring the receiver to follow each carrier's own Rosetta mapping.

Current deterministic reference support includes:

- BOOT/Rosetta/program validation;
- carrier-local symbol translation;
- bounded deterministic micro-rule execution;
- fail-closed UNKNOWN on unknown symbols/missing transitions;
- symbol-permutation equivalence tests;
- receiver artifact promotion gates;
- stored Master Prompt SHA-256 identity check;
- isolated Hybrid public bundle.

The current `generated/MASTER_PROMPT.md` is `REFERENCE_CANDIDATE`, **not PROMOTED**. Model-facing synthetic Hybrid and cross-model permutation evidence are still required.

## Accessible Memory Lab R0

The working target remains approximately:

```text
W = 4000 token-equivalent active model-facing context
```

W is not the amount Origami may store and not the cumulative amount of information a task may consume. The intended cycle is:

```text
large Origami memory
  -> selective access / internal computation / Dependency Folding
  -> bounded exposure
  -> compact retained semantic state
  -> release obsolete evidence
  -> next access
```

Synthetic evidence supports the addressing/computation hypothesis but does not yet establish end-to-end model-driven Hybrid access to a materialized visual carrier.

## Tlaloc relationship

Tlaloc may use rich swarm/Tlaloque behavior to search for receiver strategies and distill successful behavior into candidate universal bootstrap prompts, BOOT/Rosetta strategies and deterministic micro-agent candidate rules.

Origami remains authoritative for its semantics and for storing/promoting receiver artifacts. A Tlaloc candidate is not automatically an Origami promoted artifact.

## OHF research track

OHF is a nested carrier/protocol track, not the complete project identity.

**Protocol reference:** R3.10  
**Executable laboratory:** R3.10-LAB  
**Canonical deterministic laboratory language:** Go

The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` is still missing. Therefore byte-identical historical R3.10 baseline claims remain blocked. This blocker is specific to historical OHF baseline recovery; it does not invalidate the Origami 6.x semantic contracts.

Current implemented/experimental OHF assets include deterministic experiment identity, truth-first fixtures, SHA-256 artifact verification, Native blind bundles, Hybrid public bundles, observations/scoring, failure-to-regression plumbing, Dependency Folding, Accessible Memory/SuperIndex work, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD, DVR/PSS, perceptual-orthogonality evidence and reversible laboratory-tool installation.

Glyph Calculus and SAFE_MICRO_ISA are currently treated as possible low-level deterministic substrates under the Hybrid receiver; their profile-specific physical transport encodings are not universal carrier semantic meanings.

Important OHF invariants remain:

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
carrier target <= 500 KB
nominal PSS != SAFE_PSS
pairwise PASS != higher-order promotion
perception != resolution != execution != verification
external tests are progressive/cached/rate-limited
```

## Origami 6.x perceptual state

`origami.perceptual-channels.r0` defines spatial, interference, depth, temporal and emergent channels. It includes Temporal Latent Images and the invariant that a motion-bound or otherwise latent percept is not absent merely because one static observation does not reveal it.

Contract awareness is not runtime implementation. Individual operations such as `MOIRE`, `STEREO_BIND`, `PARALLAX_RESOLVE`, `KINETIC_REVEAL`, `TEMPORAL_INTEGRATE` and `TEMPORAL_DECAY` require explicit implementation/evidence before being reported as executable.

## Source of truth

Use:

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
docs/CURRENT_STATE.md
docs/ARCHITECTURE.md
docs/STATE_SEMANTICS_R0.md
docs/PERCEPTUAL_CHANNELS_R0.md
spec/PERCEPTUAL_CHANNELS_R0.json
spec/HYBRID_RECEIVER_R0.json
receiver/registry/
changes/
state/components/
```

Historical R2/R3.10 materials are evidence and continuity inputs, not authority over the current project boundary.
