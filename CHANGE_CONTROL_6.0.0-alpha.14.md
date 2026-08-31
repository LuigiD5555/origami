# Origami 6.0.0-alpha.14 — Temporal Graph Automata R0

## Scope

This release adds the first executable temporal/automaton layer on top of the alpha.13 self-describing read/write protocol.

Implemented:

- `origami.automaton-model.r0` for deterministic semantic graph automata;
- `origami.temporal-program.r0` for generative temporal state using rules, deltas and checkpoints rather than literal video frames;
- `origami.temporal-graph.r0` for semantic-temporal T2 navigation;
- `origami.temporal-codec.r0` with declared `ST*` read/simulate and `ET*` write/fold temporal codec families;
- deterministic Go runtimes for cell/rule stepping, trace generation, checkpoint seek and replay verification;
- `origami-automaton`, `origami-temporal` and `origami-temporal-carrier` managed CLIs;
- `origami.temporal-carrier.r0.profile-1`, a 640x640 exactly-8192-byte self-contained PNG carrying a visible semantic graph/timeline plus an exact compressed TemporalProgram payload;
- exact embedded-program SHA/CRC verification and explicit oversize failure instead of hidden sidecars;
- CI contracts and canonical A -> B -> C temporal reference experiment.

## Hard boundaries

- `SEMANTIC_FILM != LITERAL_VIDEO_FRAMES`;
- temporal semantic navigation must not require global exact replay;
- unsupported temporal capability returns `UNKNOWN` / `NOT_VERIFIED` rather than invention;
- exact temporal recovery remains separately verifiable;
- arbitrary unstructured temporal data is not claimed to fit the fixed carrier;
- Tlaloc may generate/test candidates, but Origami remains authoritative for its representation and pixel compilation;
- deterministic implementation does not constitute held-out Native VLM temporal evidence.

## Evidence status

Deterministic contracts, runtime stepping, replay, seek, carrier size, exact program roundtrip and managed lifecycle are implemented and CI-gated. Real Native VLM temporal interpretation and cross-model temporal interoperability remain evidence pending.
