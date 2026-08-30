# Self-Contained Book R0 experimental changelog

This file records only the non-promoted `experiment/self-contained-book-r0` work so the canonical repository `CHANGELOG.md` remains byte-identical to `main`.

## Profile 2 — image-only Master Prompt bootstrap

- embeds exact canonical `MASTER_PROMPT.md` UTF-8 bytes inside the same PNG as the book carrier;
- adds visible `S0 SELF BOOT` and `MP0 MASTER PROMPT` panels around the Fixed Carrier R2 control plane;
- binds Master Prompt length and SHA-256 inside the exact self-contained header;
- adds `-mode prompt` for `PNG -> MASTER_PROMPT.md` extraction;
- makes decode recover and verify both prompt and declared PDF extracted-text plane;
- mechanical required inputs become `PNG + user question`;
- native clean-model image-only bootstrap remains an empirical, unpromoted claim;
- book residual remains literal transport, not semantic Fold;
- 512 KB accounting remains independent of exact round-trip success.
