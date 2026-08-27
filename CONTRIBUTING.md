# Contributing to Origami HyperFold

OHF is experimental protocol research. A change is useful only if its effect can be identified and replayed.

Before opening a pull request:

```bash
make check
```

For changes affecting Glyph Calculus or Native evaluation, also run the relevant reproducibility scripts documented in `docs/EXPERIMENT_PROTOCOL.md`.

Please keep generated `runs/` and raw external-model responses out of commits unless they have been curated into a minimal reproducible regression.

Every material protocol/lab change should update or add a record under `changes/` and update the corresponding `state/components/*.json` file.

Do not describe an experimental capability as promoted unless its component state says so.
