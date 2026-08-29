# Hybrid Receiver R0 — practical quickstart

Status: **EXPERIMENTAL / REFERENCE CANDIDATE / NOT PROMOTED**

This is the shortest reproducible path from the current repositories to a real `carrier.png`, deterministic Origami tools and an automatic OpenAI-compatible model↔tool loop coordinated by Tlaloc.

## 1. Build the Origami Hybrid tools

From the Origami repository root:

```bash
make build
```

The experimental Hybrid binaries are written to:

```text
bin/origami-hybrid-carrier
bin/origami-hybrid-tool
```

They are intentionally built from the feature tree and are not claimed as globally promoted/installable release commands yet.

## 2. Generate the self-contained synthetic carrier

```bash
./bin/origami-hybrid-carrier
```

Default input:

```text
experiments/hybrid-carrier-synthetic-r0.json
```

Default output:

```text
runs/hybrid-carrier-synthetic-r0/
├── public/
│   ├── MASTER_PROMPT.md
│   ├── carrier.png
│   └── model_packet.json
└── private/
    ├── carrier.envelope.json
    ├── carrier.meta.json
    └── source.json
```

Only `public/` belongs on the answering side. `private/` exists for evaluator/debug evidence and sharing it with the answering model contaminates the experiment.

Before returning, the generator verifies that the logical carrier can be reconstructed directly from `public/carrier.png`; it does not use the private envelope to perform that check.

## 3. Inspect BOOT from the image-backed runtime

```bash
./bin/origami-hybrid-tool \
  -op BOOT
```

The tool:

1. reads `public/model_packet.json`;
2. verifies the SHA-256 identity of `public/carrier.png`;
3. decodes the Glyph Calculus transport directly from the PNG;
4. verifies the embedded memory commitment;
5. exposes the BOOT structure.

The defaults point to `runs/hybrid-carrier-synthetic-r0/public/`.

## 4. Query one address/key

```bash
./bin/origami-hybrid-tool \
  -op LOOKUP \
  -query K7F91
```

Expected synthetic value:

```text
VIOLET-48271
```

## 5. Follow the dependency path

```bash
./bin/origami-hybrid-tool \
  -op FOLLOW \
  -query K7F91 \
  -relation depends \
  -depth 2
```

The synthetic path is:

```text
K7F91 / VIOLET-48271
  -> P1Z83 / DELTA-77102
  -> R9C20 / AMBER-10593
```

The tool reports per-operation touched/exposed metrics separately from carrier-open/index-build metrics. Parsing/index construction therefore cannot be silently presented as zero-cost selective access.

## 6. Verify the carrier memory

```bash
./bin/origami-hybrid-tool -op VERIFY
```

`VERIFY` succeeds only after the image-backed runtime has already checked the embedded memory SHA-256 during carrier open.

## 7. Run the automatic Hybrid loop with a local model

The Tlaloc feature branch adds an OpenAI-compatible multimodal tool loop. The model receives the Master Prompt and carrier image, may request declared Origami functions, and Tlaloc executes those requests through `origami-hybrid-tool` before returning the tool results to the model.

Assuming sibling clones such as:

```text
projects/
├── origami/
└── tlaloc/
```

start LM Studio's OpenAI-compatible local server, load a model that supports the image/tool behavior being tested, then from `tlaloc/behavior-lab/` run:

```bash
go run ./cmd/behaviorlab receiver-run \
  -endpoint http://127.0.0.1:1234/v1 \
  -model <YOUR_MODEL_ID> \
  -prompt ../../origami/runs/hybrid-carrier-synthetic-r0/public/MASTER_PROMPT.md \
  -carrier ../../origami/runs/hybrid-carrier-synthetic-r0/public/carrier.png \
  -packet ../../origami/runs/hybrid-carrier-synthetic-r0/public/model_packet.json \
  -origami-tool ../../origami/bin/origami-hybrid-tool \
  -question 'What is the value of the second-order depends dependency of K7F91?'
```

Adjust relative paths if the repositories are located elsewhere.

Tlaloc sends the carrier image to the model and exposes only the declared Origami functions. The Origami tool itself re-validates `carrier.png` against `model_packet.json` on every invocation and reads the logical memory from the image, not from `private/source.json` or `private/carrier.envelope.json`.

The JSON result records the final answer, model/tool turns, number of tool calls, tool-output bytes/token-equivalent and model-reported token usage when the endpoint supplies it.

## 8. First held-out question

The first synthetic question is:

```text
What is the value of the second-order `depends` dependency of K7F91?
```

Ground truth is evaluator-side only:

```text
AMBER-10593
```

A valid Hybrid route may resemble:

```text
BOOT
-> LOOKUP K7F91
-> FOLLOW depends depth=2
-> VERIFY when exactness is required
-> ANSWER
```

The model is not required to choose exactly this sequence if another declared selective route is semantically equivalent and evidence-correct. It must not receive `private/source.json` or `private/carrier.envelope.json`.

## 9. Native diagnostic

To isolate visual self-boot behavior, start a fresh VLM session and provide only:

```text
System prompt: public/MASTER_PROMPT.md
Image:         public/carrier.png
Question:      held-out question
```

No Origami runtime tools are exposed. This tests VLM BOOT/ROSETTA readability and must not be confused with the Hybrid result.

## What this implementation proves and does not prove

The deterministic/transport integration can prove:

- carrier payload is actually integrated into the PNG;
- PNG -> glyphs -> bytes -> BOOT/ROSETTA/PROGRAM/INDEX/MEMORY is reversible for the generated transport;
- memory commitment detects mutation;
- queries use the image-backed index/runtime rather than a hidden source sidecar;
- missing keys fail as UNKNOWN without a semantic global-scan fallback;
- carrier target `<= 500 KB` is enforced by the generator;
- Tlaloc can execute the OpenAI-compatible image + declared-tool conversation protocol without importing Origami as a code dependency.

It does **not** yet prove:

- a real VLM successfully discovers BOOT from this carrier;
- a VLM correctly learns arbitrary carrier-local Rosetta assignments;
- the current Master Prompt is optimal;
- Tlaloc's swarm has already converged on the best receiver behavior;
- cross-model Hybrid behavior is reliable;
- a real PDF can already be transferred through the same carrier at useful scale.

Those require actual held-out model-facing campaigns. They remain promotion gates rather than assumptions hidden inside the deterministic implementation.
