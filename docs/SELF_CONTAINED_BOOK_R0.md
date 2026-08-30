# Self-Contained Book R0 Lab

Status: `EXPERIMENTAL_LAB_NOT_PROMOTED`

`origami.self-contained-book.r0-lab` tests the stronger claim:

```text
PDF + canonical Master Prompt
 -> one self-contained Origami image
 -> recover Master Prompt + exact extracted-text plane
```

It does **not** replace Fixed Carrier R2 and it does not promote native model-only readability.

## One-image deployment claim

The required model-facing payload is intended to be:

```text
book.png
+ user question
```

The complete canonical `generated/MASTER_PROMPT.md` UTF-8 bytes are physically embedded inside the PNG payload and SHA-256 verified when recovered.

A build may still export:

```text
MASTER_PROMPT.extracted.md
model_packet.json
```

but those are convenience/debug artifacts. They are **not required reconstruction state**.

`model_packet.json` remains a manifest only. It is not the Master Prompt.

## Self bootstrap order

The experimental wrapper adds an explicit visible bootstrap around the existing Fixed Carrier R2 control plane:

```text
S0  SELF BOOT
    declares that the image contains MP0
    declares prompt byte count + SHA-256
    instructs receiver to load MP0 before entering the ordinary carrier

MP0 MASTER PROMPT
    exact UTF-8 bytes embedded in the pixel payload
    SHA-256 verified before use

T0  Fixed Carrier BOOT
T1  ROSETTA + visual probe
T2  root SuperIndex
T3  deterministic machine record / memory binding
VERIFY
```

The visible S0/MP0 panels are bootstrap metadata. The exact embedded MP0 bytes are authoritative.

## Important capability boundary

There are two different claims:

```text
SELF_CONTAINED_MECHANICAL
  one PNG contains enough bytes to recover the exact Master Prompt and declared book source plane

SELF_CONTAINED_NATIVE_MODEL
  a previously unprepared compatible model can perceive/bootstrap/extract/use MP0 from the image without an external decoder
```

This experiment implements and tests the first claim.

The second remains an **empirical perception/promotion experiment**. Do not report mechanical recovery as proof that every VLM can natively extract the embedded prompt.

## Payload layout

```text
VISIBLE CONTROL / BOOTSTRAP
  S0 self boot panels
  Fixed Carrier R2
    T0 BOOT
    T1 ROSETTA
    T2 root index
    T3 machine record
    VERIFY

EXACT PIXEL PAYLOAD
  self-contained header
    source length/hash
    Master Prompt length/hash
    residual length/hash
    page/block counts
  MP0 exact UTF-8 Master Prompt bytes
  exact book residual
    pdftotext -layout UTF-8
    -> bzip2 -9
```

The book residual is still **literal transport**. It establishes an exact-recovery baseline and must never be reported as semantic Fold.

Embedding the Master Prompt also does not make the Master Prompt itself a semantic compression victory; it is required self-bootstrap state.

## Hard interpretation rules

```text
ONE_IMAGE_CONTAINS_MASTER_PROMPT
EXTERNAL_MASTER_PROMPT_REQUIRED = false
MASTER_PROMPT_HASH_VERIFIED
PASS_SELF_CONTAINED_ROUNDTRIP != PASS_512K
MECHANICAL_SELF_BOOT != NATIVE_MODEL_PROMOTION
LITERAL_RESIDUAL != SEMANTIC_FOLD
ROSETTA_ALWAYS_PRESENT
FALSE_EXACT = 0
```

## Source plane

R0 embeds the UTF-8 text produced by:

```bash
pdftotext -layout -enc UTF-8 book.pdf -
```

Therefore exact book recovery means byte equality with this declared extracted-text plane, **not** byte equality with the original PDF and not recovery of figures that exist only as raster/vector objects.

A later multimodal adapter must account for those objects explicitly.

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

Build must prove both round trips from the generated PNG:

```text
PNG -> MP0 -> exact Master Prompt bytes
PNG -> residual -> exact declared book text bytes
```

## Recover only the Master Prompt from the image

```bash
./bin/origami-self-contained-book \
  -mode prompt \
  -in runs/self-contained-book-r0/book.png \
  -out runs/self-contained-book-r0/MASTER_PROMPT.recovered.md
```

To print it directly:

```bash
./bin/origami-self-contained-book \
  -mode prompt \
  -in runs/self-contained-book-r0/book.png \
  -out -
```

## Decode both prompt and book plane

```bash
./bin/origami-self-contained-book \
  -mode decode \
  -in runs/self-contained-book-r0/book.png \
  -out runs/self-contained-book-r0/recovered.txt \
  -prompt-out runs/self-contained-book-r0/MASTER_PROMPT.recovered.md
```

Then verify independently:

```bash
cmp generated/MASTER_PROMPT.md runs/self-contained-book-r0/MASTER_PROMPT.recovered.md
pdftotext -layout -enc UTF-8 /path/to/book.pdf expected.txt
cmp expected.txt runs/self-contained-book-r0/recovered.txt
```

## Visual contract

With the default wide laboratory canvas, the upper area contains three visible regions:

```text
left:   S0 SELF BOOT + MP0 hash/order
center: existing 640x640 Fixed Carrier R2 with T0/T1/T2/T3/VERIFY
right:  MP0 MASTER PROMPT status/purpose panel
```

Below them is the exact pixel payload.

This makes the Master Prompt discoverable as a named part of the carrier instead of silently hiding it behind a JSON manifest.

## Success ladder

1. S0 self-bootstrap visible.
2. ROSETTA visible.
3. complete Master Prompt bytes physically inside the PNG.
4. Master Prompt SHA-256 exact recovery.
5. exact declared book-text round-trip.
6. external Master Prompt not required for mechanical decode.
7. `FALSE_EXACT=0`.
8. report PNG size honestly.
9. if PNG > 512000 bytes, size objective remains failed.
10. test native clean models separately before claiming image-only native bootstrap.
11. replace book residual bytes with generative/addressable Fold structures only when deterministic Unfold preserves the declared information.

Target:

```text
S0 + MP0 + BOOT + ROSETTA + SUPERINDEX + FOLD + residual + VERIFY <= 512000 bytes
```

without silently moving required reconstruction state outside the image.
