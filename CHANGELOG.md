# Changelog

## Unreleased — Self-Contained Book R0 lab (experimental branch)

- adds `origami.self-contained-book.r0-lab.profile-2` as a failure-first, non-promoted experiment;
- embeds the exact canonical `MASTER_PROMPT.md` UTF-8 bytes inside the same PNG as the book carrier;
- adds visible `S0 SELF BOOT` and `MP0 MASTER PROMPT` panels around the existing Fixed Carrier R2 control plane;
- binds Master Prompt length and SHA-256 in the self-contained payload header;
- adds `-mode prompt` so the Master Prompt can be recovered from the PNG alone;
- makes `decode` recover and verify both the Master Prompt and declared PDF extracted-text plane;
- changes model-facing required inputs for the mechanical experiment from `MASTER_PROMPT.md + PNG + question` to `PNG + question`;
- keeps `SELF_CONTAINED_NATIVE_MODEL` explicitly unpromoted: exact mechanical extraction does not prove a clean VLM can natively bootstrap the prompt from pixels;
- keeps the book body residual classified as literal transport, not semantic Fold;
- preserves `FALSE_EXACT = 0`, `ROSETTA_ALWAYS_PRESENT`, and independent 512 KB size accounting.

