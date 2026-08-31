# Origami Temporal Codec — R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTATION`

Temporal behavior extends the existing Origami semantic codec registry. It does not replace `S*/E*` and it is not required for Core compatibility.

## Decoder family

```text
ST0 READ_AUTOMATON
ST1 READ_CELL
ST2 READ_TIMELINE
ST3 LOCATE_EVENT
ST4 READ_TRANSITION
ST5 UNFOLD_TEMPORAL_REGION
ST6 SIMULATE_DECLARED_STEP
```

## Encoder family

```text
ET0 ENCODE_AUTOMATON
ET1 ENCODE_CELL
ET2 ENCODE_TIMELINE
ET3 ENCODE_EVENT
ET4 ENCODE_TRANSITION
ET5 FOLD_TEMPORAL_REGION
ET6 ENCODE_CHECKPOINT
```

A temporal read codec requires both `SEMANTIC_READ` and `TEMPORAL_READ`. A temporal write codec requires `SEMANTIC_WRITE` and `TEMPORAL_WRITE`.

A model without temporal capability remains a valid lower-level Origami receiver and returns `UNKNOWN` when a temporal semantic operation is required.

## First executable pairs

```text
ET0 <-> ST0  automaton construction roundtrip
ET2 <-> ST2  temporal program construction roundtrip
```

The reference Go implementation also verifies deterministic replay after the ET2/ST2 roundtrip.

## Step semantics

`ST6` may execute only the declared automaton state, local rules, graph and declared context. It may not invent a plausible next step.

## Exact boundary

The semantic temporal family is separate from byte-exact trace/recovery codecs. Exact temporal evidence may be added later without making it a prerequisite for semantic timeline navigation.
