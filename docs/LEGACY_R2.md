# Legacy R2 baseline

The repository originally centered on an R2 carrier that embedded a ZIP envelope into a binary visual grid.

Its reusable prompt remains at:

```text
prompts/OHF_R2_MASTER_PROMPT.txt
```

R2 is preserved because it documents an important exactness lesson:

```text
visual plausibility != exact recovery
```

The R2 exact path relied on deterministic pixel decoding, CRC32 and SHA-256 verification.

## Why it is no longer the main README

R2 is not the current OHF architecture because:

- its dominant visual language is effectively a binary grid;
- Native VLMs may struggle to reconstruct large grids exactly;
- it does not implement the current Perception/Resolution/Execution/Verification separation;
- it does not implement the current Glyph Calculus / SAFE_MICRO_ISA research;
- it predates the current R3.10-LAB reproducibility infrastructure.

Do not use the R2 prompt as the Master Prompt for current R3.10-LAB experiments.
