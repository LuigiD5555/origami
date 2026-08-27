# OHF architecture

This document describes the current architectural contract. It separates long-term protocol structure from experimental renderers and model-specific profiles.

## Sender

```text
SOURCE
  ↓
INGESTION
  ↓
CANONICAL SOURCE MODEL
  ↓
STRUCTURE / SEMANTICS / ENTROPY ANALYSIS
  ↓
SOURCE GRAPH
  ↓
PATTERN DISCOVERY
  ↓
GRAMMAR / GRAPH / TRANSFORMS / MOTIFS / HIERARCHY
  ↓
FOLD CANDIDATES
  ↓
REPRESENTATION TOURNAMENT
  ↓
GLOBAL OPTIMIZATION
  ↓
GENERATIVE IR + RESIDUAL
  ↓
SUPERINDEX + DEPENDENCY MANIFESTS
  ↓
ATTENTION / WINDOW PLAN
  ↓
VISUAL COMPILER / PERCEPTUAL ABI
  ↓
ROSETTA
  ↓
OHF CARRIER
```

## Receiver

```text
OHF CARRIER
  ↓
BOOT
  ↓
ROSETTA
  ↓
SUPERINDEX
  ↓
QUERY PLAN
  ↓
ATTENTION ROUTER / WINDOW SLIDERS
  ↓
PERCEPTION
  ↓
EVIDENCE
==== PERCEPTION WALL ====
  ↓
RESOLUTION
==== EXECUTION GATE ====
  ↓
BOUNDED EXECUTION
  ↓
RESIDUAL
  ↓
VERIFICATION SPINE
==== COMMIT GATE ====
  ↓
ANSWER / EXACT DATA
```

## Epistemic separation

### Perception

Perception is allowed to produce only evidence-like states:

```text
KNOWN
AMBIGUOUS
UNKNOWN
INVALID
```

Perception does not create verified truth.

### Resolution

Resolution receives candidate sets, evidence and constraints. It may eliminate impossible states and collapse to one candidate only when the remaining state is uniquely determined.

Resolution does not inspect pixels.

### Execution

Execution is deterministic and bounded. The project deliberately avoids turning the carrier into a general-purpose VM.

Representative operations at the full IR level include:

```text
LITERAL
REF
CONCAT
REPEAT
SLICE
PATCH
TRANSFORM
RULE
GRAPH_EXPAND
MOTIF_EXPAND
DEFAULT
OVERRIDE
RESIDUAL
VERIFY
```

Glyph Calculus investigates a much smaller local instruction set for visually transported data.

### Verification

Only verification can promote data to exact/verified status.

Possible checks include:

- length;
- CRC;
- local/block hash;
- Merkle proof;
- root hash;
- source hash.

`FALSE_EXACT = 0` is the governing invariant.

## Addressability and selective unfolding

OHF does not target the workflow:

```text
carrier -> expand everything -> search
```

The intended workflow is:

```text
carrier
  ↓
SuperIndex
  ↓
query-specific logical addresses
  ↓
small dependency closure
  ↓
selective unfold
```

Windowing can be multidimensional rather than purely spatial: graph radius, hierarchy depth, dependency depth, order and verification scope may all be adjusted independently.

## Visual families

OHF treats visual capabilities as generative families with subdimensions, states and parametric renderings rather than one flat enum or a requirement to use every dimension simultaneously.

Visual capacity may come from:

```text
state
position
relation
topology
containment
order
graph
hierarchy
composition
```

A compound is not considered reliable merely because its primitive dimensions work separately.

## Macro -> meso -> micro

The current research direction explicitly allows multiple perceptual scales:

```text
MACRO
  global composition / topology / symmetry / regions

MESO
  groups / branches / motifs / routes / hierarchies

MICRO
  glyph physical states / residual / verification marks
```

Macro-Gestalt is experimental and has no promoted semantic codebook yet.

## Native versus computational modes

### Native

The blind Native boundary receives only:

```text
Master Prompt
one carrier image
query
```

No ground truth, seed, manifest, ABI sidecar or deterministic decoder is allowed through that boundary.

### Computational

The deterministic evaluator may use Go, ground truth, hashes and pixel access.

### Hybrid

Hybrid combines visual perception with deterministic execution. Glyph Calculus and SAFE_MICRO_ISA are especially relevant to this mode because they aim to keep the VLM's work local and mechanical.
