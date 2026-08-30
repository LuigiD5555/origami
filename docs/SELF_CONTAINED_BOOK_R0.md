# Self-Contained Book R0 Lab

Status: `EXPERIMENTAL_LAB_NOT_PROMOTED`

`origami.self-contained-book.r0-lab` is a failure-first experiment for the stronger claim:

```text
PDF -> one self-contained Origami image -> exact extracted-text recovery
```

It does **not** replace Fixed Carrier R2 and it does not change the canonical visual profile.

## Model-facing package

The carrier image and the Master Prompt are different artifacts.

```text
MASTER_PROMPT.md
  = portable instructions that teach a compatible model how to READ/WRITE Origami

book.png
  = the concrete carrier: BOOT + ROSETTA + INDEX + MEMORY + VERIFY

model_packet.json
  = manifest only; it binds the exact carrier and Master Prompt by SHA-256 and declares read order
```

The JSON is **not** the Master Prompt.

A build now emits all three in the output directory:

```text
runs/self-contained-book-r0/
├── book.png
├── MASTER_PROMPT.md
└── model_packet.json
```

The self-contained carrier header embeds `master_prompt_sha256`. The copied Markdown prompt must hash to exactly that value. Therefore a model-facing package cannot silently substitute another prompt without breaking the declared binding.

Portable baseline input is:

```text
MASTER_PROMPT.md
+ book.png
+ user question
```

The stronger future `IMAGE_ONLY_BOOTSTRAP` experiment—where the complete portable Master Prompt itself is recoverable from the image—is deliberately separate because it consumes carrier capacity and changes the deployment claim.

## Why this lab exists

Fixed Carrier R2 intentionally keeps the PNG as a frozen visual control plane and keeps the corpus outside the image. This lab asks a different question: how far can a complete book be carried inside one image while still preserving Origami BOOT/ROSETTA/root navigation and exact verification?

The first baseline deliberately separates two planes:

```text
SEMANTIC / CONTROL PLANE
  Fixed Carrier R2
  T0 BOOT
  T1 ROSETTA
  T2 root index (PAGE / CONCEPT / SOURCE / GRAPH / VERIFY)
  T3 deterministic machine record / GraphSignature
  VERIFY

EXACT RESIDUAL PLANE
  pdftotext -layout UTF-8 source
  -> bzip2 -9
  -> self-describing binary header
  -> black/white row-major pixel matrix
```

The residual plane is **literal transport**. It is present to establish an exact-recovery ceiling and must never be reported as semantic Fold.

## Hard interpretation rule

```text
PASS_ROUNDTRIP != PASS_512K
LITERAL_RESIDUAL != SEMANTIC_FOLD
ROSETTA_ALWAYS_PRESENT
MASTER_PROMPT_SHA256_BOUND
FALSE_EXACT = 0
```

## Source plane

R0 embeds the UTF-8 text produced by:

```bash
pdftotext -layout -enc UTF-8 book.pdf -
```

Therefore exact recovery means byte equality with this declared extracted-text plane, **not** byte equality with the original PDF file and not recovery of figures that exist only as raster/vector images.

A later multimodal document adapter must account for those objects explicitly.

## Build

```bash
make build

./bin/origami-self-contained-book \
  -mode build \
  -in /path/to/book.pdf \
  -prompt generated/MASTER_PROMPT.md \
  -out runs/self-contained-book-r0/book.png
```

Dependencies:

```text
pdftotext
bzip2
```

The command writes the PNG, copies the exact Master Prompt, writes `model_packet.json`, and prints a JSON report containing source/residual sizes, hashes, prompt binding, page/block counts, PNG size, 512 KB delta and round-trip result.

## Decode

```bash
./bin/origami-self-contained-book \
  -mode decode \
  -in runs/self-contained-book-r0/book.png \
  -out runs/self-contained-book-r0/recovered.txt
```

Then compare:

```bash
pdftotext -layout -enc UTF-8 /path/to/book.pdf expected.txt
cmp expected.txt runs/self-contained-book-r0/recovered.txt
```

## Current visual contract

The upper control region is the existing 640x640 Fixed Carrier R2 visual language. It visibly includes T0 BOOT, T1 ROSETTA, T2 root navigation, T3 machine state and VERIFY.

The lower region is the exact residual bit matrix. The matrix has no authority to redefine ROSETTA. It is decoded mechanically as one bit per pixel, row-major, black=1, white=0.

The current T2 root index acts only as the **root SuperIndex**. Book-specific chapter/page-group routing is the next experiment and must be represented independently of the exact residual before it can be credited as semantic Fold.

## Success ladder

1. ROSETTA visible: required.
2. Master Prompt copied and SHA-bound: required for the portable package.
3. exact extracted-text round-trip: required.
4. `FALSE_EXACT=0`: required.
5. report the PNG byte size honestly.
6. if PNG > 512000 bytes, status remains `PASS_ROUNDTRIP_FAIL_512K`.
7. replace residual bytes with generative/addressable Fold structures only when deterministic Unfold can recover the same declared information.
8. measure each reduction against this baseline.

Target:

```text
BOOT + ROSETTA + SUPERINDEX + FOLD + residual + VERIFY <= 512000 bytes
```

without silently moving required reconstruction state outside the image.
