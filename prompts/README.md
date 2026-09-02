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

## Historical prompts

The old OHF R2 master prompt is no longer carried in the active tree. Git history remains sufficient for archaeology/reproducibility; current work uses `generated/MASTER_PROMPT.md` plus carrier-local BOOT/ROSETTA/PROGRAM semantics.
