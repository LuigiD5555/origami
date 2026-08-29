# Contributing to Origami

Origami is experimental representation/state-machine research. The repository also contains the OHF carrier/laboratory track.

Before opening a pull request, identify the ownership scope of the change:

- **core Origami** — state semantics, perceptual channels, machine/dynamics, general projection contracts;
- **OHF track** — R3.10-LAB carrier tooling, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and related experiments.

Run the smallest relevant verification closure. For current Go changes, start with:

```bash
make check
```

For OHF visual/Glyph/Native work, also run the relevant reproducibility scripts documented in `docs/`.

Every material change should record before/after state, evidence, affected downstream components, regressions and promotion status. Root changes update `state/ORIGAMI_STATE.json`; track/component changes update the corresponding `state/components/*.json` and `changes/*.json` records.

Do not describe an experimental capability as promoted unless its evidence and state record support that claim.
