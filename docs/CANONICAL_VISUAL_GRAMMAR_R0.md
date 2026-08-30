# Canonical Visual Grammar R0

Status: `EXPERIMENTAL_REFERENCE_PROFILE`

Origami uses **one canonical visual aesthetic per promoted profile version**.

This is a correction to an earlier interpretation in which every carrier could invent a completely different private visual dialect. `ROSETTA` remains mandatory, but its role is to make the carrier self-describing and forward-compatible, not to authorize arbitrary aesthetic drift.

## Core rule

```text
ONE ORIGAMI PROFILE VERSION
          |
          v
ONE CANONICAL AESTHETIC
          |
          +-- geometry
          +-- contrast/fill
          +-- topology
          +-- position
          +-- enclosure
          +-- scale
          +-- repetition
          +-- density
          +-- limited declarative text
          |
          v
ROSETTA restates active bindings
          |
          v
many documents / many semantic payloads
```

Different PDFs, images, conversations or semantic graphs should therefore look recognizably like the same Origami visual language when compiled under the same profile version.

## Aesthetic means functional encoding

The aesthetic is not a branding skin. It is the current optimized visual encoding strategy for machine readability.

The present R0 profile is primarily high-contrast monochrome because those channels have the strongest deterministic and experimental support in the current Fixed Carrier work.

Current canonical visual dimensions are:

| Dimension | Primary role |
|---|---|
| geometry | node/semantic class and anchors |
| fill / contrast | bit/state redundancy |
| position | hierarchy, route and ordering |
| topology | relations and dependencies |
| enclosure | space, scope and address regions |
| scale | macro / meso / micro levels |
| repetition | redundancy, agreement and verification |
| density | data/control structure and information hints |
| limited text | BOOT, ROSETTA and small declarative labels |

Experimental dimensions include color, mathematical/numeric structures and temporal/phase channels.

They are not forbidden. They simply do not become primary semantic authority until experimental evidence earns promotion.

## Current primitive family

The current profile standardizes recognizable roles such as:

```text
ringed square       -> root / verification anchor
black-center square -> page / major container
white square        -> concept
white diamond       -> source / evidence origin
small black square  -> graph / relation cluster
line / arrow        -> declared relation
box / enclosure     -> space / scope / address region
black/white cells   -> visual probe bits
```

The exact renderer may improve in later canonical profiles, but changes must be versioned and evidence-backed.

## Why ROSETTA remains mandatory

Even with one canonical aesthetic, the carrier still includes ROSETTA because it must be independently interpretable.

ROSETTA declares:

```text
profile id
profile version
active dimensions
active primitive -> semantic role mappings
carrier-approved parameters/extensions
```

This gives three benefits:

1. a model does not need hidden repository knowledge to read an old carrier;
2. old carriers remain readable after the canonical aesthetic evolves;
3. experimental/profile extensions are explicit rather than inferred.

The rule is:

```text
KNOWN PROFILE + ROSETTA
```

not:

```text
KNOWN PROFILE, THEREFORE ROSETTA IS UNNECESSARY
```

and not:

```text
ROSETTA, THEREFORE EVERY DOCUMENT MAY INVENT A NEW STYLE
```

## Visual intent before pixels

A writer does not start by choosing pretty shapes. It first declares what the carrier must communicate:

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

The canonical profile then maps each intent onto approved visual dimensions.

Example:

```text
RELATION      -> topology / line
SCOPE         -> enclosure
HIERARCHY     -> position
VERIFICATION  -> repetition + ringed anchor
EVIDENCE      -> source geometry + address link
```

This is the difference between representation design and arbitrary illustration.

## Mathematical structures are candidate channels

The prime-number pattern example is useful because it demonstrates a broader principle: visual organization can emerge from mathematical structure rather than from ordinary diagram primitives.

Origami therefore keeps a candidate family for structures such as:

```text
periodicity
modular patterns
prime-derived organization
factorization structure
number-theoretic spacing
```

But a visually interesting pattern is not automatically a useful Origami encoding.

Tlaloc must experimentally establish whether a candidate improves some combination of:

```text
semantic recovery
model readability
routing accuracy
robustness under resize/compression
density
selective navigation
context efficiency
cross-model transfer
```

without increasing false exactness or destroying existing semantics.

## How the aesthetic evolves

The intended lifecycle is:

```text
CURRENT ORIGAMI CANONICAL PROFILE
              |
              v
       Tlaloc experiments
   prompt / channel / layout /
 color / shape / numeric structure
              |
              v
     candidate profile revision
              |
              v
 deterministic + model evaluation
              |
              v
       Origami validation
              |
              v
        Tonal promotion gate
              |
              v
NEW CANONICAL PROFILE VERSION
```

Tlaloc may propose many strange or radically different representations. They are laboratory candidates until promoted.

A promoted change may visibly alter Origami's aesthetic. That change applies as a new version of the canonical language, rather than as a private visual dialect for one document.

## Compatibility

Every carrier binds to its profile/version and carries ROSETTA.

Therefore:

```text
Origami profile R0 carrier
    remains R0-readable

even after

Origami profile R1 is promoted
```

A future receiver can use BOOT + profile version + ROSETTA to enter the old carrier without pretending that the latest aesthetic was always in force.

## Hard invariants

```text
ONE_CANONICAL_AESTHETIC_PER_PROFILE_VERSION
ROSETTA_ALWAYS_PRESENT
AESTHETIC_IS_FUNCTIONAL_ENCODING
MODEL_DOES_NOT_STYLE_EACH_DOCUMENT_INDEPENDENTLY
SEMANTIC_IR_IS_STYLE_INDEPENDENT
TLALOC_PROPOSES / ORIGAMI_VALIDATES / TONAL_PROMOTES
EXPERIMENTAL_CHANNEL != CANONICAL_AUTHORITY
OLD_PROFILE_REMAINS_DECODABLE
FALSE_EXACT = 0
```
