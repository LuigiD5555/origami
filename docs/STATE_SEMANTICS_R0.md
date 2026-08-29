# Origami state semantics — R0 coherent-state profile

This profile is quantum-inspired as a computational analogy; it is **not** a claim of physical quantum computation and it does not define the entire scope of Origami. Origami remains a general visual/computational representation language.

## State kinds

- `determinate`: one active branch.
- `superposed`: multiple coherent branches with complex amplitudes.
- `coupled`: joint state whose members must not be evaluated as independent states.
- `observed`: explicitly resolved state.

## Core laws

1. `TRANSFORM` evolves a valid state but has no authority to select a branch.
2. `OBSERVE` is an explicit resolution boundary.
3. `INTERFERE` combines complex amplitudes; equal/opposite paths cancel.
4. zero amplitude means cancellation, not unknown.
5. absence, unknown and inhibited are distinct semantic states.
6. a coupled state is a joint object; decomposing it requires an explicit operation.

## Relation to Fold / Unfold

R0 uses these as Origami semantic operations, not as claims about physical collapse:

- `UNFOLD`: construct or expand a set of alternatives under a declared representation rule.
- `FOLD`: constrain and/or explicitly resolve according to a declared policy.

## Relationship to perceptual channels

Origami `6.0.0-alpha.2` adds a separate perceptual-channel contract in `PERCEPTUAL_CHANNELS_R0.md`.

The two layers are intentionally orthogonal:

```text
coherent-state kind != perceptual availability condition
```

For example, a `superposed` state may project to a `motion-bound` Temporal Latent Image, and an `observed` state may still require a stereo/parallax observation mode to expose one of its encoded relations.

Perceptual conditions therefore do not add new coherent-state kinds.

This coherent-state profile is one representation mode that Tlaloc can learn to enforce in a target model.
