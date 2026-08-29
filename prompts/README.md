# Prompt directory

Prompt files here belong to specific Origami research profiles or historical protocol generations.

Origami does **not** define one giant universal prompt containing all carrier semantics. The experimental Hybrid Receiver instead uses a small **universal bootstrap discipline** at `generated/MASTER_PROMPT.md`. Carrier-specific symbol meaning, micro-program behavior, index structure and verification instructions remain self-described by each carrier through `BOOT -> ROSETTA -> PROGRAM -> INDEX -> MEMORY -> VERIFICATION`.

This distinction matters:

```text
universal bootstrap != universal carrier semantics
```

## Current experimental receiver prompt

For `origami.hybrid-receiver.r0`, the current stored reference candidate is:

```text
generated/MASTER_PROMPT.md
```

Its provenance/status is recorded under:

```text
receiver/registry/r0/manifest.json
```

It is `REFERENCE_CANDIDATE`, not yet PROMOTED. Tlaloc may discover improved candidates through receiver distillation; Origami remains responsible for validating and storing promoted artifacts.

## `OHF_R2_MASTER_PROMPT.txt`

**Status: LEGACY / OHF R2 ONLY.**

It is retained for reproducibility of the older binary-grid ZIP carrier and must not be treated as the current Origami 6.x semantic contract, the current OHF R3.10-LAB Native profile, or the Hybrid Receiver bootstrap.

Current OHF-LAB Native fixtures are generated/tested through laboratory commands and blind-bundle boundaries. Generated receiver prompts remain derived artifacts rather than an alternate source of Origami semantics.
