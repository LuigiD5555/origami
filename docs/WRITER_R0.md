# Origami Writer R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`

Writer R0 adds the missing symmetry to the universal Origami bootstrap:

```text
READ  : carrier -> semantic state
WRITE : semantic state -> canonical carrier plan -> compiler -> carrier
```

The key restriction is that WRITE mode uses the **single canonical Origami visual profile**. A model does not invent a new visual style for every PDF.

## PDF / image / text to Origami

The intended path is:

```text
PDF / IMAGE / TEXT / CONVERSATION
            |
            v
Tlaloc or declared source ingestion
Canonical Document IR / source structure
            |
            v
SEMANTIC IR
concepts / relations / hierarchy / states /
uncertainty / evidence / verification / addresses
            |
            v
ORIGAMI WRITE PLAN
visual intents
            |
            v
CANONICAL VISUAL GRAMMAR
            |
            v
ROSETTA + PROGRAM + INDEX + MEMORY + VERIFICATION
            |
            v
DETERMINISTIC ORIGAMI COMPILER
            |
            v
carrier.png
            |
            v
DECODE / ROUNDTRIP VERIFY
```

A screenshot of the PDF is not this process.

## Construction specification

`origami-write-plan` takes a semantic graph and produces a deterministic construction plan.

The plan includes:

```text
profile_id
source_identity
semantic_commitment
visual_intents
ROSETTA
required carrier sections
verification plan
canonical semantic graph
hard invariants
```

The reference profile is:

```text
origami.canonical-aesthetic.r0
```

## Visual intents

The model/semantic pipeline describes **what needs to be visible**, not how it wants the image to look aesthetically.

Writer R0 supports intent classes such as:

```text
IDENTITY
HIERARCHY
RELATION
DEPENDENCY
SCOPE
IMPORTANCE
STATE
UNCERTAINTY
EVIDENCE
VERIFICATION
ADDRESSABILITY
TEMPORAL_ORDER
REDUNDANCY
ROUTING
```

The canonical profile maps them to approved dimensions.

For example:

```text
IDENTITY       -> geometry
RELATION       -> topology
DEPENDENCY     -> topology
SCOPE          -> enclosure
STATE          -> fill / contrast
VERIFICATION   -> repetition / verify marks
ADDRESSABILITY -> enclosure / address region
TEMPORAL_ORDER -> static positional sequence in R0
```

The last example matters: the current profile does not silently use an unpromoted temporal visual channel simply because one exists as a research idea.

## ROSETTA generation

WRITE mode generates ROSETTA as part of the carrier construction plan.

ROSETTA records:

```text
canonical profile ID
profile version
active dimensions
visual token -> semantic role mappings
explicit profile-approved extensions
```

This means another model can receive the finished carrier later and bootstrap it without knowing the source PDF.

## Model versus compiler authority

A language model may:

- interpret source semantics;
- propose visual intents;
- emit a construction request/specification;
- use declared compiler tools.

It is **not** the pixel-level carrier authority.

If no compiler is available, the correct output is:

```text
CONSTRUCTION_SPEC_ONLY
```

not a claim that an arbitrary generated image is a valid Origami carrier.

This prevents a visually plausible diagram from bypassing deterministic structure and verification.

## Roundtrip gate

Every compiled carrier must satisfy the conceptual gate:

```text
DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR
```

Identity/address/hash/proof fields require exact equality.

Semantic fields require equivalence under declared canonicalization.

Failure means:

```text
INVALID_CARRIER
```

not “close enough.”

## Example input

```json
{
  "schema": "origami.writer.r0.request",
  "source_identity": "paper.pdf",
  "semantic_graph": {
    "nodes": {
      "root": {
        "address": "root",
        "kind": "root",
        "label": "Document",
        "value": {"status": "PRESENT", "data": "paper"},
        "dependencies": ["concept"],
        "verified": true
      },
      "concept": {
        "address": "concept",
        "kind": "concept",
        "label": "Dynamic programming",
        "value": {"status": "UNKNOWN"}
      }
    }
  }
}
```

Run:

```bash
./bin/origami-write-plan -in request.json -out plan.json
```

The output does not yet mean “PNG exists.” It is the deterministic, canonical front-end plan that a renderer/compiler consumes.

## Relationship with Tlaloc

Tlaloc is allowed to experiment aggressively with candidate visual strategies and prompt variations.

It may ask questions such as:

```text
Would color improve state discrimination?
Would another node shape reduce VLM errors?
Would more spatial redundancy survive JPEG better?
Would a prime/modular pattern encode a graph property more efficiently?
Would a different Master Prompt reduce BOOT failures?
```

Those experiments produce **candidate profile revisions**.

They do not let each individual WRITE call use whichever aesthetic a model happens to prefer.

The lifecycle is:

```text
Tlaloc candidate
  -> benchmark
  -> evidence
  -> Origami contract validation
  -> Tonal promotion
  -> new canonical Origami profile version
```

## Relationship with the Master Prompt

The Master Prompt now contains both READ and WRITE procedures.

In WRITE mode it instructs a model to:

1. ingest/understand the source;
2. create Semantic IR;
3. declare visual intents;
4. bind to the canonical profile;
5. generate ROSETTA;
6. compile through deterministic tools;
7. perform roundtrip verification.

This means the same universal prompt can teach a compatible model how to use Origami in both directions without embedding the source document or a private visual style in the prompt itself.

## Hard invariants

```text
WRITE_USES_CANONICAL_VISUAL_PROFILE
ROSETTA_REQUIRED
ONE_DOCUMENT != ONE_NEW_AESTHETIC
SOURCE_SCREENSHOT != ORIGAMI
SEMANTIC_PLAN_PRECEDES_PIXELS
COMPILER_IS_CARRIER_AUTHORITY
ROUNDTRIP_REQUIRED
FALSE_EXACT = 0
```
