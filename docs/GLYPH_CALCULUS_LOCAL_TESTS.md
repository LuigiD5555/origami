# Glyph Calculus — Exact Local Test Instructions

Run from repository root. Do not continue past a failed gate.

```bash
go test ./...
go vet ./...
```

Inspect one physical byte:

```bash
go run ./cmd/ohf-glyphcalc glyph 0xAA
```

Expected tuple is `shape=2 holes=2 direction=2 frames=2`.

Run the local matrix:

```bash
./scripts/glyphcalc_local_matrix.sh
```

Required local gates include exact reconstruction, structured reduction, pseudo-random RAW fallback, deterministic replay and tamper detection.

Prepare SAFE_MICRO_ISA locally:

```bash
./scripts/microisa_prepare_local.sh
```

This creates local blind bundles and makes zero external requests.

For real VLM testing, supply only the files under a generated `blind/` directory. Never supply ground truth, expected responses, source bytes or transport bytes.

External testing should begin with one small probe, proceed adaptively, use concurrency 1 by default, obey Retry-After, cache completed probes and stop early on systematic failure.

Keep the trust stages separate:

```text
PERCEPTION != EXECUTION != VERIFICATION
```
