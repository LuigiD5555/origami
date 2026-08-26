# Origami HyperFold (OHF)

**Origami HyperFold is an experimental protocol for encoding data into a compact, addressable carrier that can be navigated and reconstructed by an LLM and/or deterministic software.**

The long-term goal is not simply to make another compressor. OHF explores whether structured information can be represented through reusable rules, graphs, references, transforms and residual data so that a receiver can **unfold only what it needs** instead of expanding the complete source for every query.

The repository currently includes the **R2 exact visual-carrier baseline** and the scientific/architectural documentation for the larger HyperFold system.

> OHF is experimental research. Extreme ratios such as 1:10,000 are a target to investigate on highly structured/generative data, not a universal compression guarantee.

## What R2 does today

R2 provides a simple exact baseline:

```text
SOURCE FILE / PROJECT
        ↓
       ZIP
        ↓
R2 binary envelope
        ↓
CRC32 + SHA-256
        ↓
binary visual grid
        ↓
      PNG
```

The PNG contains a visible **BOOT** header and a deterministic binary payload. The BOOT header tells a decoder how to locate and interpret the grid.

When the receiving environment has pixel-processing/code execution available, it can reconstruct the embedded ZIP and verify it cryptographically before calling the result exact.

R2 is intentionally much simpler than the final OHF architecture. It gives the project a deterministic transport and exactness baseline while Folding, SuperIndex, Perception Lab and Native decoding continue to evolve.

## Quick start — test OHF with another LLM

The easiest experiment uses exactly two inputs:

1. `prompts/OHF_R2_MASTER_PROMPT.txt`
2. an `.ohf-r2.png` carrier from `examples/`

Start a **fresh conversation** with the target LLM.

### Step 1 — send the Master Prompt

Copy the complete contents of:

```text
prompts/OHF_R2_MASTER_PROMPT.txt
```

Do not add a decoder manifest, ABI file or other helper metadata.

### Step 2 — send the OHF image

Upload the original PNG without resizing, screenshotting or recompressing it.

The image's visible BOOT section describes the grid geometry and verification information required by R2.

### Step 3 — ask a question

For a project carrier, examples include:

```text
List the files contained in this OHF carrier.
```

```text
Show me exactly README.md.
```

```text
Which file contains the main entry point?
```

```text
Reconstruct the ZIP and verify its SHA-256.
```

For a document stored inside a carrier you can similarly ask for its contents after the ZIP has been verified.

### Step 4 — inspect the fidelity reported by the model

OHF distinguishes:

| Fidelity | Meaning |
|---|---|
| `SEMANTIC` | interpreted/summarized information |
| `STRUCTURAL` | claims grounded in verified paths/structure |
| `EXACT_CONTENT` | exact requested content from a verified member |
| `EXACT_SOURCE` | reconstructed ZIP matches the SHA-256 stored by the encoder |

A model must not claim exactness merely because its visual interpretation looks plausible.

If deterministic decoding is unavailable or verification fails, the correct answer is `UNKNOWN` / `NOT VERIFIED`.

## Why the Master Prompt is reusable

R2 is designed around a **universal prompt + self-describing carrier** model.

The prompt defines the decoding contract:

```text
how to read BOOT
how to interpret binary cells
how to verify the envelope
how to report fidelity
```

The image supplies carrier-specific values such as grid geometry, payload length and hashes.

Therefore a new source normally requires a **new carrier image, not a new Master Prompt**, as long as it uses the same R2 protocol version.

A protocol-version change may require an updated prompt.

## Why BOOT matters

The decoder must not depend on magic coordinates hard-coded for one image.

BOOT declares the geometry required to locate the binary payload. The Master Prompt explicitly instructs the receiver to use those values instead of assuming them.

The envelope is:

```text
magic[8]
version:u16 big-endian
payload_length:u32 big-endian
crc32:u32 big-endian
sha256[32]
zip[payload_length]
```

The R2 magic value is:

```text
OHFR2ZIP
```

## Exact decoding path

For an exact request the intended process is:

```text
OHF PNG
  ↓
read BOOT
  ↓
locate GRID
  ↓
sample cell centers
  ↓
bits → bytes
  ↓
parse R2 envelope
  ↓
extract ZIP
  ↓
CRC32 check
  ↓
SHA-256 check
  ↓
VERIFIED ZIP
```

Only after verification should the receiver navigate or extract exact members.

## Selective use

Once the ZIP has been reconstructed and verified, a narrow question should not cause every member to be semantically analyzed.

For example:

```text
"Show me exactly config.yaml"
```

should conceptually become:

```text
verify carrier once
      ↓
read ZIP directory
      ↓
locate config.yaml
      ↓
extract requested member
      ↓
return content
```

This is an early baseline for OHF's broader **Selective Unfolding** principle.

## Current limitations of R2

R2 should not be confused with the final research target.

Currently:

- the visual grid primarily acts as an exact binary transport;
- exact decoding works best when the receiving environment can inspect pixels programmatically;
- a VLM without a sandbox may not be able to reconstruct a large binary carrier exactly;
- R2 does not yet provide the planned full generative Folding system;
- R2 does not prove a 1:10,000 compression ratio;
- R2 exactness comes from deterministic reconstruction + hashes, not from LLM consensus.

These limitations are precisely why OHF also contains Perception Lab and the larger R3+/1.0 architecture.

## Where OHF is going

The complete architecture is intended to evolve toward:

```text
SOURCE
  ↓
Canonical Source Model
  ↓
Dedup / Grammar / Graph / Motif / Transform discovery
  ↓
Representation Tournament
  ↓
Generative IR + Residual
  ↓
SuperIndex
  ↓
Attention Router + Window Sliders
  ↓
Visual Compiler
  ↓
OHF Carrier
```

On reception:

```text
QUERY
  ↓
SuperIndex
  ↓
small dependency window
  ↓
Perception / Resolution
  ↓
Deterministic Execution
  ↓
Residual
  ↓
Verification
  ↓
ANSWER
```

The represented object may be large while the active working context remains small.

## Repository layout

```text
origami/
├── README.md
├── docs/
│   └── SCIENTIFIC_FOUNDATIONS.md
├── prompts/
│   └── OHF_R2_MASTER_PROMPT.txt
└── examples/
    └── ... OHF R2 carrier images
```

`README.md` is the practical introduction and usage guide.

`docs/SCIENTIFIC_FOUNDATIONS.md` contains the detailed scientific motivation, architecture, assumptions, experimental status and theoretical foundations.

`prompts/OHF_R2_MASTER_PROMPT.txt` is the reusable R2 receiver prompt.

`examples/` contains carriers intended for cross-model experiments.

## Scientific documentation

For the detailed explanation of information theory, MDL/Kolmogorov motivation, grammar folding, graph motifs, Representation Tournament, SuperIndex, Window Sliders, Perception Lab, swarms and Verification Spine, see:

**`docs/SCIENTIFIC_FOUNDATIONS.md`**

## Core rule

OHF uses one rule above all others:

> **UNKNOWN is preferable to false exactness.**

Semantic interpretation may be probabilistic. Exact recovery is deterministic and independently verified.

## Status

OHF is under active experimental development. R2 is a working exact-carrier baseline; the broader HyperFold architecture remains a staged research and engineering effort toward a generative, addressable and selectively unfoldable protocol.
