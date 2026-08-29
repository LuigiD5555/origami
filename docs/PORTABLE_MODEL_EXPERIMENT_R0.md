# Portable Model Experiment R0

## Goal

Make an Origami experiment runnable in an external/local model with an unambiguous receiver contract.

Two modes are intentionally separate.

### Native

Public answering packet:

```text
generated/MASTER_PROMPT.md
carrier.png
question
```

The answering VLM receives no source, oracle, private manifest or external decoder.

### Computational

Public answering packet:

```text
generated/MASTER_PROMPT.md
model_packet.json
question
Origami tool endpoint/runtime
```

The answering model can use bounded Origami operations but cannot read the underlying source/memory directly.

## Experiment directory contract

A prepared portable run should have:

```text
runs/<run-id>/
  public/
    MASTER_PROMPT.md
    model_packet.json
    carrier.png              # required for Native/Hybrid
    questions.json
  private/
    manifest.json
    oracle.json
    source/                   # synthetic or real canonical source
  results/
    responses.jsonl
    metrics.json
    report.md
```

`public/` is the only directory that may be shared with the answering model. `private/` is evaluator-only.

## model_packet.json minimum fields

```json
{
  "schema": "origami.model-packet.r0",
  "mode": "native|computational|hybrid",
  "working_window_token_eq": 4000,
  "master_prompt_sha256": "...",
  "carrier_sha256": "...",
  "allowed_tools": [],
  "anti_contamination": true
}
```

The packet describes the interface; it must not contain the answer key or hidden source facts.

## questions.json

Questions must be generated or selected after the source representation/carrier is frozen whenever the experiment is intended to test retrieval rather than memorized benchmark behavior.

Each public question has an opaque id and question text. Expected answers remain only in `private/oracle.json`.

## Native carrier requirement

The carrier is not a screenshot of the source text. It is an Origami visual/computational representation with enough boot/Rosetta/index structure for the receiver to identify how to navigate it.

Native carrier target remains <= 500 KB. Failure to encode irreducible exact information within that target must be reported rather than hidden behind a poster, external file or undeclared side channel.

## Local manual protocol

1. Start a fresh model session with no source document in context.
2. Use a vision-capable model for Native/Hybrid.
3. Set `public/MASTER_PROMPT.md` as the system prompt when the client supports system prompts; otherwise paste it before any question.
4. Attach only `public/carrier.png` for Native.
5. Ask one entry from `public/questions.json`.
6. Save the model's answer verbatim.
7. Repeat in a fresh session when independence between questions is required.
8. Evaluate answers against `private/oracle.json` outside the answering model.

## Contamination failures

Mark a trial CONTAMINATED when the answering model can access any of:

- original source document;
- private manifest/oracle;
- hidden full-memory dump;
- previous answer containing the oracle;
- another retrieval channel not declared in `model_packet.json`.

A contaminated correct answer is not a PASS.

## Promotion sequence

1. Synthetic Native boot/readability.
2. Synthetic Computational navigation.
3. Synthetic Hybrid.
4. Frozen real-source Computational.
5. Frozen real-source Native/Hybrid.
6. Cross-model replication.

Do not use a real PDF as the first proof that the receiver protocol works.
