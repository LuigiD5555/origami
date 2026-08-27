# Current OHF state

**Status date:** 2026-08-27  
**Protocol reference:** R3.10  
**Executable laboratory:** R3.10-LAB  
**Canonical implementation language:** Go

## Production/reference state versus laboratory state

OHF distinguishes protocol reference from experimental implementation.

### R3.10 reference

R3.10 is the current protocol reference used by this project. The exact historical `ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip` artifact is not present in this repository.

Therefore:

```text
historical_baseline_bytes = MISSING_R310_SDK
```

The laboratory must report this as blocked. It must not reconstruct the historical SDK from prose documentation and then call the reconstruction the baseline.

### R3.10-LAB

The executable repository contains a reproducible development workbench through the F1 evaluator foundation plus the newer Glyph Calculus / SAFE_MICRO_ISA experiments.

## Implemented foundation

| Area | State |
|---|---|
| isolated LAB workspace | PASS |
| ExperimentSpec canonicalization | PASS |
| content-addressed ExperimentID / RunID | PASS |
| deterministic namespaced seeds | PASS |
| truth-first synthetic source generation | PASS |
| SHA-256 artifact store | PASS |
| deterministic LAB-only geometry renderer | PASS |
| PNG save/reopen and pixel hash verification | PASS |
| hard carrier-size gate | PASS |
| Native blind three-input boundary | PASS |
| ObservationEnvelope | PASS |
| external scorer | PASS |
| failure -> regression pipeline | PASS |
| Glyph Calculus deterministic runtime | EXPERIMENTAL / implemented |
| SAFE_MICRO_ISA profiler | EXPERIMENTAL / implemented |
| Context SIMD bundles/checker | EXPERIMENTAL / implemented |
| Macro-Gestalt | DESIGN/EXPERIMENT ACTIVE; not promoted |
| historical R3.10 SDK freeze | BLOCKED: artifact missing |

## Current hard invariants

```text
FALSE_EXACT = 0
UNKNOWN > invented exactness
GO = canonical deterministic runtime
Native = Master Prompt + one image + query only
Perception != Resolution != Execution != Verification
carrier target <= 500 KB
real failures become regressions
experimental dimensions require promotion evidence
```

## Current visual codebook reference

R3.10's documented physical byte model uses four four-state physical dimensions:

```text
outer_shape
holes
direction_beacon
frame_count
```

and the deterministic mapping:

```text
byte = shape + 4*holes + 16*direction + 64*frame
```

The LAB is investigating whether a VLM should output the byte directly or emit smaller physical state tuples and let Go perform the arithmetic.

## Active research: Glyph Calculus

Goal: reduce LLM work to small, local, mechanical perception.

Conceptual boundary:

```text
IMAGE
  ↓
VLM: observe small physical states only
  ↓
PERCEPTION WALL
  ↓
Go: packing, references, runs, sequence expansion, verification
```

Current bounded experimental operations include `HALT`, `LIT`, `RUN`, `REF` and `SEQ`.

The compiler is allowed to fall back to raw transport when structure does not reduce the description cost.

## Active research: SAFE_MICRO_ISA

The project does not assume that an operation such as counting, adding or decoding is cheap for every model. Instead it profiles candidate operations using model-specific evidence.

Current candidates:

- `OBSERVE4`
- `TUPLE4`
- `DECODE8`

Selection priority:

1. `false_known == 0`;
2. highest exact-known rate;
3. fewer retries;
4. fewer output tokens per glyph;
5. lower latency per glyph;
6. lower measured monetary cost per glyph, if available.

The result is a profile, not a universal law about all VLMs.

## Active research: Context SIMD

Independent glyph operations may be grouped into lanes in a single instruction context.

This is **logical batching**, not a claim that the underlying model executes hardware SIMD.

Lane width is widened adaptively:

```text
1 -> 4 -> 8 -> 16
```

and testing stops at the first unsafe boundary.

## Macro-Gestalt status

Macro-Gestalt explores whether global recognizable composition can become an additional carrier dimension through existing OHF visual families such as motif, symmetry, spatial zone and macro-region.

It is not currently promoted. No everyday visual concept (tree, flower, butterfly, etc.) has fixed protocol meaning. Physical rendering and carrier-local semantics must remain separate.

## External testing policy

Local generation may be exhaustive when inexpensive. External tests are progressive and conservative:

```text
max concurrency: 1 by default
small batches
cache completed probes
resume without repeating
respect 429 / Retry-After
bounded retries
exponential backoff
circuit breaker
stop early on structural failure
```

The objective is both experimental efficiency and respectful use of external services.

## Known blocker

The principal historical blocker remains:

```text
MISSING_R310_SDK
```

This does not block LAB-only research, but it blocks claims that the repository contains a byte-identical historical R3.10 baseline.

## Source of truth

Implementation state must be recoverable from the repository, not from conversation history.

See:

```text
state/ORIGAMI_STATE.json
state/components/
changes/
```
