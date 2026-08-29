# Glyph Calculus v0.1 — Local Execution Report

Date: 2026-08-27

Scope: deterministic local execution only; no external VLM was contacted.

Local gates passed in the development kit: byte↔glyph roundtrip, bounded execution, exact reconstruction, structured glyph reduction, pseudo-random RAW fallback, deterministic replay and tamper detection.

For 4,096-byte synthetic structured fixtures, local transport used substantially fewer glyph bytes than literal byte-per-glyph transport; the pseudo-random control selected RAW automatically. All reconstructed sources matched byte-for-byte.

A small 1,024-byte mixed fixture compiled to a much smaller transport in the development run. This is only local representation evidence and must not be presented as Native VLM evidence.

Current interpretation:

> A small bounded rewrite language can represent many structured output bytes with fewer transport glyphs than literal byte-per-glyph representation, and Go can execute it exactly.

Native perceptual reliability remains a separate Perception Lab gate.
