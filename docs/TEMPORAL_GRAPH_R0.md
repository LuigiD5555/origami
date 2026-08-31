# Origami Temporal SuperGraph — R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTATION`

T2 remains semantic-first, but gains temporal routes so a receiver can jump directly to the relevant agent, rule, event, phase or checkpoint instead of replaying a complete trajectory.

## Root routes

```text
ROOT
├── CONTENT
├── CONCEPTS
├── AGENTS
├── RULES
├── TIMELINE
├── EVENTS
├── CHECKPOINTS
├── EVIDENCE
└── VERIFY
```

Typical queries route as follows:

```text
What agents exist?       -> AGENTS
What does A do?          -> AGENTS/A -> RULES
When did X happen?       -> EVENTS/X -> TIMELINE
Why did X happen?        -> EVENTS/X -> RULES -> INPUTS
Replay phase P           -> TIMELINE/P -> CHECKPOINT -> bounded temporal unfold
```

The graph is an address map, not proof authority. Address and content identity remain distinct.

## Native semantic rule

Temporal navigation must not automatically escalate to exact trace recovery. A semantic question should use the smallest sufficient semantic-temporal route.

```text
TEMPORAL_NAVIGATION != GLOBAL_REPLAY
```

If the requested temporal region cannot be decoded with declared receiver capabilities, the result is `UNKNOWN`, not a fabricated event history.
