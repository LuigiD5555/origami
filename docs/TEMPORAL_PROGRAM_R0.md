# Origami Temporal Program — R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTATION`

## Semantic film

Origami does not define a temporal behavior as a literal stack of video frames. The default representation is generative:

```text
initial state
+ graph
+ declared rules
+ schedule
+ deltas
+ checkpoints
= semantic film
```

A complete trajectory may therefore be represented without storing every full state repeatedly when the trajectory has reusable structure.

## Modes

R0 recognizes:

```text
PROCEDURAL
DELTA
CHECKPOINT_PLUS_DELTA
```

The reference runtime uses `CHECKPOINT_PLUS_DELTA`.

A frame records only state changes. A checkpoint stores a complete state at a declared step so random temporal access can start there rather than replaying from t=0.

## Core operations

```text
RUN(program) -> trace
SEEK(trace,t) -> state_t
VERIFY_REPLAY(program,trace)
MEASURE(trace)
```

Reference equations:

```text
TRACE = REPLAY(PROGRAM)
STATE_t = SEEK(TRACE,t)
REPLAY(FOLD(TRACE)) ~= TRACE
```

## Exactness

`Replay equivalence` is deterministic for the declared R0 machine. An implementation must not call a lossy summary an exact temporal trace.

A failed delta precondition invalidates replay rather than silently patching state.

## Compression boundary

Temporal Fold exploits repeated rules, sparse changes and checkpoints. It does not imply that an arbitrary incompressible trajectory can fit in a fixed carrier.

```text
COMPRESSIBILITY_DEPENDS_ON_STRUCTURE
```

## Reference implementation

```text
internal/temporal
cmd/origami-temporal
```

Examples:

```bash
origami-temporal -mode run -in program.json -out trace.json
origami-temporal -mode seek -in program.json -trace trace.json -step 12
origami-temporal -mode verify -in program.json -trace trace.json
origami-temporal -mode metrics -trace trace.json
```

## Metrics

R0 reports:

```text
steps
changes
checkpoints
explicit_snapshot_bytes
folded_trace_bytes
temporal_fold_ratio
```

`temporal_fold_ratio` is a transport/representation metric. It is not by itself evidence that a VLM can understand the representation.
