# Self-Contained Book R0 Lab

Status: `EXPERIMENTAL_LAB_NOT_PROMOTED`

`origami.self-contained-book.r0-lab` is a failure-first experiment for the stronger claim:

```text
PDF -> one self-contained Origami image -> exact extracted-text recovery
```

It does **not** replace Fixed Carrier R2 and it does not change the canonical visual profile.

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

A run can pass exact recovery and still fail the Origami size objective.

```text
PASS_ROUNDTRIP != PASS_512K
LITERAL_RESIDUAL != SEMANTIC_FOLD
ROSETTA_ALWAYS_PRESENT
FALSE_EXACT = 0
```

The lab reports both conditions independently.

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
  -out runs/self-contained-book-r0/book.png
```

Dependencies for this R0 experiment:

```text
pdftotext
bzip2
```

The command prints a JSON report containing source/residual sizes, hashes, page/block counts, PNG size, 512 KB delta and round-trip result.

## Decode

```bash
./bin/origami-self-contained-book \
  -mode decode \
  -in runs/self-contained-book-r0/book.png \
  -out runs/self-contained-book-r0/recovered.txt
```

Then compare against the declared source plane:

```bash
pdftotext -layout -enc UTF-8 /path/to/book.pdf expected.txt
cmp expected.txt runs/self-contained-book-r0/recovered.txt
```

## Current visual contract

The upper control region is the existing 640x640 Fixed Carrier R2 visual language. It therefore visibly includes T0 BOOT, T1 ROSETTA, T2 root navigation, T3 machine state and VERIFY.

The lower region is the exact residual bit matrix. The matrix has no authority to redefine ROSETTA. It is decoded mechanically as one bit per pixel, row-major, black=1, white=0.

The current T2 root index acts only as the **root SuperIndex**. Book-specific chapter/page-group routing is the next experiment and must be represented independently of the exact residual before it can be credited as semantic Fold.

## Success ladder

R0 establishes the baseline:

1. ROSETTA visible: required.
2. exact extracted-text round-trip: required.
3. `FALSE_EXACT=0`: required.
4. report the PNG byte size honestly.
5. if PNG > 512000 bytes, status remains `PASS_ROUNDTRIP_FAIL_512K`.
6. replace residual bytes with generative/addressable Fold structures only when deterministic Unfold can recover the same declared information.
7. measure each reduction against this baseline.

The target is eventually:

```text
BOOT + ROSETTA + SUPERINDEX + FOLD + residual + VERIFY <= 512000 bytes
```

without silently moving required reconstruction state outside the image.
