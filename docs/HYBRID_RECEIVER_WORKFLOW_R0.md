# Origami Hybrid Receiver Workflow R0

## Purpose

Origami receives a candidate produced by Tlaloc, preserves its distilled behavior, binds that behavior to a concrete self-boot carrier, stores all generated receiver artifacts, and only promotes them after Origami-owned validation.

## Import from Tlaloc

Tlaloc produces:

```text
tlaloc.origami-hybrid-artifact-set.r0
  - universal_prompt
  - boot_strategy
  - rosetta_constraints
  - micro_program
  - source_trace_sha256
  - working_window_token_eq
```

Import it from the Origami repository root:

```bash
go run ./cmd/origami-receiver-import \
  -in /path/to/tlaloc/behavior-lab/generated/origami-hybrid-receiver-r0.artifact-set.json \
  -root receiver/candidates
```

Origami stores:

```text
receiver/candidates/<candidate-id>/
  MASTER_PROMPT.md
  candidate.json
  manifest.json
```

The imported manifest is always `REFERENCE_CANDIDATE`. Import is not promotion.

## Authority split

Tlaloc proposes semantic receiver behavior. Origami owns the concrete carrier.

Origami must still create/bind:

```text
BOOT
ROSETTA physical bindings
PROGRAM
INDEX
MEMORY
VERIFICATION
optional RESIDUAL
```

The physical Rosetta mapping is carrier-local and must not be taken from a global glyph dictionary hidden in the external prompt.

## Behavior preservation

The Tlaloc -> Origami import preserves:

```text
state
token
action
next_state
emit
```

Losing any of these silently would invalidate the distillation claim.

## Intended compiled form

```text
Tlaloc rich swarm
  -> semantic trace
  -> distilled micro-program
  -> Origami import
  -> concrete carrier-local BOOT/ROSETTA binding
  -> Go deterministic micro-runtime
  -> visual/computational carrier
```

The external `MASTER_PROMPT.md` remains a universal receiver bootstrap. Carrier-specific knowledge belongs in the carrier/runtime, not in a growing model-specific prompt.

## Hybrid execution

```text
model
  -> locate BOOT
  -> learn carrier-local Rosetta
  -> high-level navigation / ambiguity resolution

Origami runtime
  -> micro-agent transitions
  -> addressability
  -> selective unfold
  -> Fold / Dependency Folding
  -> FILTER/PROJECT/AGGREGATE/etc.
  -> verification

model
  -> retain compact semantic state
  -> request next bounded access or answer
```

The target `W≈4000 token-equivalent` constrains simultaneous model-facing Origami exposure, not memory size and not cumulative session traffic.

## Promotion gates

Before a receiver candidate becomes `PROMOTED`, Origami must have evidence for:

- deterministic semantic equivalence for the compiled micro-program;
- carrier BOOT discovery;
- Rosetta correctness with carrier-local symbol permutation;
- fail-closed UNKNOWN behavior;
- no hidden source/oracle contamination;
- no false exactness;
- bounded active model exposure;
- evidence/proof reopening;
- model-facing Hybrid tests on held-out synthetic carriers.

Real PDF/source transfer comes after these synthetic gates.
