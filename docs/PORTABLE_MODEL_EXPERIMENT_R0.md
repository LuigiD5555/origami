# Portable Hybrid Model Experiment R0

## Goal

Make an Origami carrier runnable in an external/local model with a reproducible **self-boot + hybrid execution** contract.

The primary target is Hybrid. Native and Computational remain diagnostic baselines so failures can be localized to perception/bootstrap versus deterministic execution/navigation.

## Primary Hybrid packet

The answering model receives:

```text
public/MASTER_PROMPT.md
public/carrier.png
public/model_packet.json
question
+ only the Origami tools explicitly declared by model_packet.json
```

The model must discover carrier-local semantics from the image's own `BOOT -> ROSETTA -> PROGRAM` path. The external prompt must not reveal the carrier's private symbol mapping.

Hybrid division of labor:

```text
carrier perception / BOOT discovery
              ↓
        model / VLM
              ↓
carrier-local ROSETTA + high-level route
              ↓
    Origami deterministic runtime
  micro-agents / addressing / computation
              ↓
small verifiable result or next region
              ↓
             model
```

### Native baseline

```text
MASTER_PROMPT.md
carrier.png
question
```

No external decoder or Origami runtime. This isolates self-boot/perceptual readability.

### Computational baseline

```text
MASTER_PROMPT.md
model_packet.json
question
Origami runtime/tools
```

No hidden source access. This isolates addressability, micro-machine execution, selective unfolding and context recycling from visual perception.

## Experiment directory contract

```text
runs/<run-id>/
  public/
    MASTER_PROMPT.md
    model_packet.json
    carrier.png
    questions.json
  private/
    manifest.json
    oracle.json
    source/
    carrier_semantics.json      # evaluator truth; NEVER public
  results/
    responses.jsonl
    traces.jsonl
    metrics.json
    report.md
```

Only `public/` plus the declared runtime interface may reach the answering model. `private/` is evaluator-only.

## Carrier section requirement

A Hybrid R0 carrier should expose, directly or through its boot-discovered references:

```text
BOOT
ROSETTA
PROGRAM
INDEX
MEMORY
VERIFICATION
[RESIDUAL]
```

- `BOOT` says where/how to begin.
- `ROSETTA` gives carrier-local symbol semantics.
- `PROGRAM` contains deterministic micro-agent/automata behavior.
- `INDEX` supports selective addressing.
- `MEMORY` contains represented knowledge/behavior.
- `VERIFICATION` supports proof/exactness checks.
- `RESIDUAL` is optional and must remain explicit/measured.

The carrier is not a screenshot/poster of the source.

## model_packet.json minimum fields

```json
{
  "schema": "origami.model-packet.r0",
  "receiver_contract": "origami.hybrid-receiver.r0",
  "mode": "hybrid",
  "working_window_token_eq": 4000,
  "master_prompt_sha256": "...",
  "carrier_sha256": "...",
  "allowed_tools": [],
  "anti_contamination": true
}
```

The packet describes transport/interface capability only. It must not contain hidden source facts, answer keys or the carrier-local Rosetta mapping.

## Symbol permutation test

A mandatory anti-memorization test uses at least two equivalent carriers with different physical symbol assignments:

```text
Carrier A: triangle=OPEN, circle=VALUE
Carrier B: circle=OPEN, square=VALUE
```

The same Master Prompt is used for both. Questions and underlying semantics are equivalent. Success requires discovering each local Rosetta rather than relying on globally memorized glyph semantics.

## Questions

Freeze the synthetic source and carrier before selecting/generating held-out questions. Each public question has an opaque id. Expected answers and full carrier semantics remain in `private/`.

## Local manual protocol

1. Start a fresh session with no source/PDF in context.
2. Use a vision-capable local model for Hybrid/Native.
3. Set `public/MASTER_PROMPT.md` as system prompt when supported.
4. Attach only `public/carrier.png`.
5. For Hybrid, expose only tools named in `public/model_packet.json`.
6. Ask one held-out question.
7. Save answer and tool trace verbatim.
8. Evaluate externally against `private/oracle.json` and the expected route/evidence.
9. Repeat with a symbol-permuted carrier using the same Master Prompt.

## Metrics

At minimum record:

- BOOT_FOUND;
- ROSETTA_RESOLVED;
- PROGRAM_INITIALIZED;
- ANSWER_CORRECT;
- ROUTE_INTEGRITY;
- EVIDENCE_INTEGRITY;
- UNKNOWN_CORRECT;
- FALSE_EXACT;
- TOOL_CALLS;
- INTERNAL_TOUCHED;
- PEAK_ACTIVE_TOKEN_EQ;
- CUMULATIVE_EXPOSED_TOKEN_EQ;
- RETAINED_STATE_TOKEN_EQ;
- CONTAMINATED.

## Contamination failures

Mark `CONTAMINATED` if the answering model can access the original source, oracle, private carrier semantics, a hidden full-memory dump, previous oracle-bearing answers or an undeclared decoder/retrieval channel.

A contaminated correct answer is not a PASS.

## Promotion sequence

1. Deterministic Go receiver contract + symbol-permutation tests.
2. Synthetic Computational baseline.
3. Synthetic Native BOOT/ROSETTA baseline.
4. Synthetic Hybrid end-to-end.
5. Symbol-permuted cross-model Hybrid replication.
6. Frozen real-source Hybrid.
7. PDF/image source transfer only after the synthetic gates are clean.
