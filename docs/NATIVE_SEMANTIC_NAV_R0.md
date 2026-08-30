# Native Semantic Navigation R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED_EVIDENCE_PENDING`

This change is derived from a failed external multimodal trial in which a model could read the visible Origami bootstrap but could not answer the simple question **“What is the index?”**.

Instead, the model treated the image as a binary archive that first needed pixel extraction, byte reconstruction and decompression. It then claimed byte counts, a compression format, an incorrect SHA-256 and page ranges that were not verified.

That result is recorded in:

```text
experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json
```

## Diagnosis

The old Fixed Carrier T2 visually rendered generic categories:

```text
PAGE | CONCEPT | SOURCE | GRAPH | VERIFY
```

rather than the actual semantic entry points of the represented document/memory.

At the same time, T3 displayed a dense black/white deterministic record. A model capable of reading the bootstrap but not executing a pixel decoder therefore had a strong incentive to conclude that semantic navigation required mechanical decoding.

That was the wrong dependency.

## Corrected architecture

```text
SEMANTIC VISUAL PLANE
T0 -> T1/ROSETTA -> T2 semantic superindex -> visible/selective semantic regions

EXACT / CONTROL PLANE
T3 deterministic record -> CID / hash / root / exact verification
```

The exact plane remains valuable and is not removed. It is simply no longer a prerequisite for semantic navigation.

## Query routing

```text
identity             -> T0/T1
index                -> T2
overview             -> T2 + visible semantic state
locate topic         -> T2 -> selective semantic expansion
explain topic        -> T2 -> smallest relevant semantic region
exact quote/hash     -> exact plane / declared tool when available
```

A model asking for the index must not have to reconstruct the entire corpus, decompress an archive or execute a binary grid decoder.

## T2 profile-2

`origami.fixed-carrier.r2.profile-2` changes T2 from a generic navigation legend into a bounded **actual semantic superindex**.

The carrier integration layer projects high-value/top-level canonical node labels into the bounded graph hint. The JSON hint is always kept structurally valid; it is no longer blindly truncated at 256 bytes.

The renderer gives those actual labels visual priority in T2.

If there is no usable semantic index hint, T2 explicitly displays an UNKNOWN state rather than encouraging a binary fallback.

## Backward compatibility

Profile-2 remains inside the Fixed Carrier R2 envelope:

```text
640 x 640
8192 PNG bytes
<= 512000 hard maximum
```

The deterministic decoder still accepts the legacy profile-1 profile/BOOT digest.

## Master Prompt R3

The portable prompt now classifies a question before navigation.

For semantic questions it explicitly forbids whole-image binary extraction as a prerequisite. It also forbids claims about hidden byte lengths, compression algorithms, hashes or archive contents unless an exact decoder/verification path was genuinely executed.

The prompt remains a `REFERENCE_CANDIDATE` until held-out real-model evidence passes.

## Native benchmark

The first mandatory Native semantic regression is now:

```text
Q: What is the index?
Expected route: T0 -> T1 -> T2
Binary decode allowed: NO
```

Further cases cover identity, overview, topic location and exact quotation.

The key invariant is:

```text
SEMANTIC_NAVIGATION != MECHANICAL_DECODE
```

Exact recovery may still require a deterministic runtime. Semantic usability may not silently depend on one.
