# Origami Perceptual Channels — R0

Version target: `6.0.0-alpha.2`
Contract ID: `origami.perceptual-channels.r0`

Origami may encode meaning not only through explicit marks in a static frame, but also through properties that emerge from interference, depth, motion, temporal integration, and relations among layers.

> A represented structure may exist as a perceptual phenomenon even when no single static component contains the complete visible structure.

This contract defines perceptual channels. It does **not** replace the coherent-state semantics in `STATE_SEMANTICS_R0.md`.

The falsifiability rules for non-direct perceptual claims are defined by `origami.observation-contract.r0` in `spec/OBSERVATION_CONTRACT_R0.json` and the formal rationale in `docs/FORMAL_CORE_R0.md`.

## 1. Channel families

### Spatial
Information directly available from a static projection: geometry, position, orientation, texture, contrast, density, topology, and explicit marks.

### Interference
Information that emerges from relative alignment or spatial frequency relationships among patterns. Typical mechanisms include moiré-like fields, reinforcement, cancellation, and phase-sensitive overlays.

### Depth
Information that depends on disparity, parallax, occlusion, or multi-view relationships. Depth may encode hierarchy, separation, dependency, or layer membership.

### Temporal
Information whose percept requires change over time. A temporal channel may depend on trajectory, velocity range, phase progression, persistence, or integration across multiple instants.

### Emergent
Information that is not owned by any one layer in isolation and becomes available only through interaction among layers, channels, or observations.

## 2. Explicit and emergent channels

An **explicit channel** stores a property directly in a component.

An **emergent channel** stores the conditions from which a property becomes perceptible.

For example, two individually simple line fields may produce a moiré structure that is not present as an explicit mark in either field.

## 3. Temporal Latent Image

A **Temporal Latent Image (TLI)** is a structure whose intended percept is unavailable or incomplete in an authorized static observation and becomes available only under a declared temporal transformation.

```text
carrier + mask/layer relation + trajectory + temporal integration
                              -> percept
```

The trajectory is part of the representation contract. A TLI must also provide a finite Observation Contract before its reveal claim can be promoted as verified.

A TLI may declare carrier pattern/layers, reveal transform, trajectory, accepted velocity/rate, phase constraints, integration window, visibility condition, and decay/persistence rule.

A static frame that does not reveal the intended percept MUST NOT by itself be interpreted as evidence that the represented structure is absent. Conversely, once the declared valid observation procedure has exhausted its finite budget without satisfying the expected predicate, the reveal claim is `FAIL` for that contract.

## 4. Perceptual conditions

- `latent`: encoded conditions exist but the intended percept is not currently available.
- `masked`: a declared occlusion or masking relation prevents direct access.
- `motion-bound`: availability requires an authorized temporal transformation.
- `phase-bound`: availability requires a declared relative phase/alignment.
- `stereo-bound`: availability requires binocular or multi-view disparity information.
- `emergent`: the percept belongs to an interaction, not one component.
- `revealed`: the declared reveal conditions are satisfied.
- `decayed`: a previously available temporal percept no longer satisfies its persistence/integration condition.

## 5. Operations

- `MOIRE`: construct/evaluate an interference channel from relative spatial frequency/orientation/alignment.
- `PHASE_SHIFT`: change a declared relative phase/alignment.
- `STEREO_BIND`: bind compatible views into a depth-bearing relation.
- `PARALLAX_RESOLVE`: derive depth/relational signal from authorized viewpoint or motion displacement.
- `KINETIC_REVEAL`: apply a declared trajectory to a motion-bound representation.
- `TEMPORAL_INTEGRATE`: integrate authorized observations across a declared temporal window.
- `TEMPORAL_DECAY`: update perceptual availability when persistence/integration conditions cease.

## 6. Core invariants

1. `NO_STATIC_EQUIVALENCE_ASSUMPTION`: a static projection is not assumed to contain every valid percept.
2. `EMERGENT_NOT_OWNED_BY_SINGLE_LAYER`: an emergent percept belongs to the declared interaction.
3. `MOTION_BOUND_VISIBILITY_EXPLICIT`: motion-dependent visibility declares its reveal condition.
4. `TEMPORAL_REVEAL_IS_TRAJECTORY_DEPENDENT`: reveal depends on the declared trajectory.
5. `DEPTH_REQUIRES_RELATIONAL_EVIDENCE`: depth semantics require a declared relational cue.
6. `PHASE_IS_RELATIONAL`: phase/alignment is relative to another layer, field, or reference.
7. `LATENT_IS_NOT_ABSENT`: an insufficient or unauthorized observation cannot prove absence.
8. `LATENT_IS_FALSIFIABLE`: a valid finite Observation Contract that exhausts its budget without satisfying the expected predicate returns `FAIL`.
9. `DECAY_IS_NOT_DELETION`: perceptual decay does not automatically erase representation state.
10. `OBSERVATION_CONTRACT_REQUIRED`: non-direct perceptual claims require an authorized finite observation contract before verification.

## 7. Relationship to coherent-state semantics

The coherent-state profile and perceptual-channel contract are orthogonal layers. State semantics are not perceptual availability.

## 8. Runtime status

R0 is a semantic and machine-readable contract. Origami does not yet claim a production renderer, detector, or decoder for every operation above. Implementations must report unsupported operations explicitly rather than silently approximating them.
