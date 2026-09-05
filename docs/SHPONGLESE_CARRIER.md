# Shponglese carrier experiments in Origami

## Separation of concerns

Shponglese is a semantic operational IR. Origami is a possible representation/carrier/memory substrate.

This separation is required so that representation experiments can fail without invalidating the semantic IR, and so that a representation advantage can be measured against fair baselines.

## Initial comparison contract

Given one canonical Shponglese program or trace, produce semantically equivalent forms in multiple codecs:

```text
canonical semantics
   ├── JSON
   ├── compact text
   ├── binary
   └── Origami
```

Measure the same downstream task for every codec.

Potential metrics include:

- exact semantic recovery;
- successful deterministic execution after recovery;
- bytes transported;
- model/input tokens where measurable;
- latency;
- number of selective reads/unfolds;
- error rate under corruption;
- false-known rate;
- performance under randomized symbol mappings.

## Anti-prior controls

At least one test family should use arbitrary identifiers or randomized symbol mappings that the receiving model cannot reconstruct from normal linguistic expectations.

Example semantic sequence:

```text
REF_83
CALL_12
ARG_42
ARG_91
VERIFY_3
RETURN_6
```

Mappings may be permuted between runs while preserving a declared machine-readable ground truth unavailable to the tested model.

Success on ordinary natural-language payloads does not by itself establish exact carrier capacity.

## Promotion rule

An Origami codec or profile may be promoted for Shponglese transport only from measured evidence against declared baselines. Visual novelty, nominal information density or model-assisted plausibility are insufficient.
