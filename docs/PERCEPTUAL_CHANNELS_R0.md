# Origami Perceptual Channels — R0

Version target: `6.0.0-alpha.2`
Contract ID: `origami.perceptual-channels.r0`

Origami may encode meaning not only through explicit marks in a static frame, but also through properties that emerge from interference, depth, motion, temporal integration, and relations among layers.

> A represented structure may exist as a perceptual phenomenon even when no single static component contains the complete visible structure.

This contract defines perceptual channels. It does **not** replace the coherent-state semantics in `STATE_SEMANTICS_R0.md`.

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

A TLI is not equivalent to ordinary frame-by-frame animation. Its semantic object is the relationship:

```text
carrier + mask/layer relation + trajectory + temporal integration
                              -> percept
```

The trajectory is therefore part of the representation contract.

A TLI may declare:

- carrier pattern or layer set;
- reveal transform;
- trajectory or motion family;
- accepted velocity/rate range when relevant;
- phase/alignment constraints;
- temporal integration window;
- visibility condition;
- decay or persistence rule.

A static frame that does not reveal the intended percept MUST NOT be interpreted as evidence that the represented structure is absent.

## 4. Perceptual conditions

These conditions annotate perceptual availability; they are not replacements for coherent-state kinds.

- `latent`: encoded conditions exist but the intended percept is not currently available.
- `masked`: a declared occlusion or masking relation prevents direct access.
- `motion-bound`: availability requires an authorized temporal transformation.
- `phase-bound`: availability requires a declared relative phase/alignment.
- `stereo-bound`: availability requires binocular or multi-view disparity information.
- `emergent`: the percept belongs to an interaction, not one component.
- `revealed`: the declared reveal conditions are satisfied.
- `decayed`: a previously available temporal percept no longer satisfies its persistence/integration condition.

## 5. Operations

### `MOIRE`

Construct or evaluate an interference channel whose output depends on relative spatial frequency/orientation/alignment of two or more patterns.

### `PHASE_SHIFT`

Change a declared relative phase or alignment without implying semantic resolution of unrelated state branches.

### `STEREO_BIND`

Bind compatible views into a depth-bearing perceptual relation. The output is a joint depth representation, not merely two independent images.

### `PARALLAX_RESOLVE`

Derive a depth or relational signal from authorized viewpoint/motion displacement.

### `KINETIC_REVEAL`

Apply a declared trajectory to a motion-bound representation so that a Temporal Latent Image may become perceptually available.

### `TEMPORAL_INTEGRATE`

Integrate authorized observations across a declared temporal window into one perceptual result.

### `TEMPORAL_DECAY`

Update perceptual availability when the temporal persistence/integration condition is no longer satisfied.

## 6. Core invariants

1. `NO_STATIC_EQUIVALENCE_ASSUMPTION`: a static projection is not assumed to contain every valid percept of the representation.
2. `EMERGENT_NOT_OWNED_BY_SINGLE_LAYER`: an emergent percept must not be reassigned to one component unless an explicit derivation does so.
3. `MOTION_BOUND_VISIBILITY_EXPLICIT`: motion-dependent visibility must declare its reveal condition.
4. `TEMPORAL_REVEAL_IS_TRAJECTORY_DEPENDENT`: a valid reveal depends on the declared transformation/trajectory, not merely the existence of multiple frames.
5. `DEPTH_REQUIRES_RELATIONAL_EVIDENCE`: depth semantics require disparity, parallax, occlusion, or another declared relational cue.
6. `PHASE_IS_RELATIONAL`: phase/alignment is defined relative to another layer, field, or reference.
7. `LATENT_IS_NOT_ABSENT`: failure to reveal under the current observation does not imply nonexistence.
8. `DECAY_IS_NOT_DELETION`: temporal percept decay changes perceptual availability; it does not automatically erase the underlying representation.
9. `OBSERVATION_CONTRACT_REQUIRED`: a receiver must know which observation modes are authorized before claiming a perceptual result.

## 7. Relationship to coherent-state semantics

The coherent-state profile and perceptual-channel contract are orthogonal layers.

A state can be `superposed` while one of its projections is `motion-bound`; a `coupled` state can be projected through stereo disparity; an `observed` state can still contain a temporal latent visual encoding.

Therefore:

```text
state semantics != perceptual availability
```

## 8. Runtime status

R0 is a semantic and machine-readable contract. Origami `6.0.0-alpha.2` does not yet claim a production renderer, detector, or decoder for every operation above.

Implementations must report unsupported operations explicitly rather than silently approximating them.
