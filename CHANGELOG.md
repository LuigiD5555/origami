# Changelog

## Unreleased

### Reversible installation lifecycle

- added `uninstall.sh` with tracked restore semantics;
- `install.sh` now records SHA-256 installation state and preserves the binaries that existed before the first tracked Origami install;
- reinstalls preserve the original rollback point instead of backing up a previous Origami binary as if it were user state;
- uninstall refuses to overwrite/remove binaries changed after installation unless `--force` is explicit;
- added opt-in legacy uninstall support for the previous `.bak`-based installer;
- normal smoke tests are now ephemeral and leave no persistent `runs/install-smoke` directory.


This file summarizes major repository-level milestones. Detailed experimental state belongs in `state/` and `changes/`.


## 2026-08-27 — Dimensional Visual Register / PSS formalization

- added hierarchical MICRO/MESO/MACRO visual-dimension registry;
- added executable DVR validation and scope-aware capacity reporting;
- separated nominal PSS from evidence-gated `SAFE_PSS`;
- added Perceptual Orthogonality Law and pair-evidence model;
- integrated Macro-Gestalt as the MACRO scale of the same dimensional system;
- added machine-readable state and `CHG-DVR-0001`;
- added local inspection script with zero external requests.

## 2026-08-27 — R3.10-LAB repository modernization

- moved the public repository from R2 documentation-only status to an executable Go laboratory;
- added reproducible F1 laboratory infrastructure;
- added Native isolation, observations, scoring and failure-regression plumbing;
- added Glyph Calculus;
- added SAFE_MICRO_ISA and Context SIMD experiments;
- added conservative external-execution policy;
- added installer, Makefile and CI;
- changed the Go module path to `github.com/LuigiD5555/origami`;
- replaced the R2-centric README with current R3.10-LAB documentation;
- retained R2 prompt as explicit legacy material;
- added machine-readable state/change records and current-state documentation.

## 2026-08-26 — R2 public documentation baseline

- added scientific foundations document;
- added R2 Master Prompt;
- added practical R2 README.
