# OHF documentation map

The documentation is split deliberately so that architectural ideas, current implementation status and historical material do not become mixed together.

## Read first

1. [`CURRENT_STATE.md`](CURRENT_STATE.md) — current executable state, blockers and active research.
2. [`ARCHITECTURE.md`](ARCHITECTURE.md) — current protocol architecture and trust boundaries.
3. [`EXPERIMENT_PROTOCOL.md`](EXPERIMENT_PROTOCOL.md) — reproducible experimental procedure.
4. [`DEVELOPMENT.md`](DEVELOPMENT.md) — change control and repository workflow.
5. [`ROADMAP.md`](ROADMAP.md) — next legal research steps.

## Specialized documents

- [`DIMENSIONAL_VISUAL_REGISTER.md`](DIMENSIONAL_VISUAL_REGISTER.md) — DVR/PSS, scope-aware capacity and perceptual orthogonality.
- [`DVR_PSS_EXECUTION_REPORT.md`](DVR_PSS_EXECUTION_REPORT.md) — local validation evidence for `CHG-DVR-0001`.
- [`GLYPH_CALCULUS_PROTOCOL.md`](GLYPH_CALCULUS_PROTOCOL.md)
- [`GLYPH_CALCULUS_LOCAL_TESTS.md`](GLYPH_CALCULUS_LOCAL_TESTS.md)
- [`GLYPH_CALCULUS_RESULTS.md`](GLYPH_CALCULUS_RESULTS.md)
- [`MICROISA_CONTEXT_SIMD_PROTOCOL.md`](MICROISA_CONTEXT_SIMD_PROTOCOL.md)

## Foundation document

[`SCIENTIFIC_FOUNDATIONS.md`](SCIENTIFIC_FOUNDATIONS.md) is the longer scientific/architectural research document. It remains useful for information theory, MDL, grammar/graph folding, SuperIndex, selective unfolding and verification concepts. Its historical status statements should be interpreted through `CURRENT_STATE.md`, which is the canonical human-readable status document.

## Legacy

[`LEGACY_R2.md`](LEGACY_R2.md) explains the older R2 binary-grid ZIP baseline and why it remains in the repository only for reproducibility.

## Machine-readable truth

Human documentation is explanatory. Current machine-readable implementation state is stored under:

```text
state/ORIGAMI_STATE.json
state/components/*.json
changes/*.json
```

- `INSTALLATION_LIFECYCLE.md` — reversible install/uninstall contract and safety rules.
- `INSTALL_UNINSTALL_EXECUTION_REPORT.md` — local gates proving reversible installation behavior.
