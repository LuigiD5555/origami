# Hybrid Receiver R0 — practical quickstart

Status: **EXPERIMENTAL / REFERENCE CANDIDATE / NOT PROMOTED**

This is the shortest reproducible path from the current repository to a real `carrier.png` plus deterministic Hybrid tools.

## 1. Generate the self-contained synthetic carrier

From the Origami repository root:

```bash
go run ./cmd/origami-hybrid-carrier
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

The generator verifies before returning that the logical carrier can be reconstructed directly from `public/carrier.png`; it does not use the private envelope to perform that check.

## 2. Inspect BOOT from the image-backed runtime

```bash
go run ./cmd/origami-hybrid-tool \
  -op BOOT
```

The tool:

1. reads `public/model_packet.json`;
2. verifies the SHA-256 identity of `public/carrier.png`;
3. decodes the Glyph Calculus transport directly from the PNG;
4. verifies the embedded memory commitment;
5. exposes the BOOT structure.

## 3. Query one address/key

```bash
go run ./cmd/origami-hybrid-tool \
  -op LOOKUP \
  -query K7F91
```

Expected synthetic value:

```text
VIOLET-48271
```

## 4. Follow the dependency path

```bash
go run ./cmd/origami-hybrid-tool \
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

## 5. Verify the carrier memory

```bash
go run ./cmd/origami-hybrid-tool -op VERIFY
```

`VERIFY` succeeds only after the image-backed runtime has already checked the embedded memory SHA-256 during carrier open.

## 6. Current local-model use

For a fresh VLM session:

```text
System prompt: public/MASTER_PROMPT.md
Image:         public/carrier.png
Interface:     public/model_packet.json
```

Do not provide anything from `private/`.

The current repository now provides the deterministic tool CLI, but it does **not yet automatically register these CLI operations as OpenAI-compatible function calls inside LM Studio**. Until the model↔tool harness is added, a manual Hybrid probe can relay a requested declared operation to `origami-hybrid-tool` and return only its JSON output to the model.

A direct Native diagnostic may omit the tools and give the model only the Master Prompt + image + question; that tests VLM BOOT/ROSETTA readability and must not be confused with the Hybrid runtime result.

## 7. First held-out question

Ask the model:

```text
What is the value of the second-order `depends` dependency of K7F91?
```

Ground truth is evaluator-side only:

```text
AMBER-10593
```

The desired Hybrid route is conceptually:

```text
BOOT
-> LOOKUP K7F91
-> FOLLOW depends depth=2
-> VERIFY when required
-> ANSWER
```

The answering model must not receive `private/source.json` or `private/carrier.envelope.json`.

## What this proves and does not prove

The deterministic Go path can prove:

- carrier payload is actually integrated into the PNG;
- PNG -> glyphs -> bytes -> BOOT/ROSETTA/PROGRAM/INDEX/MEMORY is reversible for the generated transport;
- memory commitment detects mutation;
- queries use the image-backed index/runtime rather than a hidden source sidecar;
- missing keys fail as UNKNOWN without a semantic global-scan fallback;
- carrier target `<= 500 KB` is enforced by the generator.

It does **not** by itself prove:

- a VLM can visually discover BOOT unaided;
- a VLM correctly learns arbitrary carrier-local Rosetta assignments;
- Tlaloc has found the optimal Master Prompt;
- cross-model Hybrid behavior is reliable;
- a real PDF can already be transferred through the same carrier at useful scale.

Those are the next model-facing promotion gates, not assumptions hidden inside this deterministic implementation.
