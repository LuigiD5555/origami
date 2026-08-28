# Origami 6.0.0-alpha.1 — State Semantics R0

## Why 6.0 is retained

The 6.0 line introduced a semantic expansion beyond dynamic Ω(t): Origami can represent determinate, superposed, coupled and observed states, with explicit interference/cancellation and resolution boundaries.

## Added / formalized

- `determinate` state;
- `superposed` state;
- `coupled` joint state;
- `observed` state;
- complex branch amplitudes for the quantum-inspired profile;
- `TRANSFORM` without implicit observation;
- `INTERFERE` and cancellation;
- `CONSTRAIN`;
- explicit `OBSERVE`;
- `FOLD` / `UNFOLD` relationship to state resolution/expansion.

## Removed from Origami ownership

The following concepts now belong to Tlaloc and must not be versioned as Origami features:

- Behavior Compiler;
- PromptIR;
- Tlaloque coordination/training (Tlaloc-owned bounded-agent layer);
- target-model orchestration;
- behavior-prompt promotion arbiter.

