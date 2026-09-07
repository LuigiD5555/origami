# Origami project boundary — archived pre-R2

This file preserves the pre-Architecture-R2 project boundary for historical provenance.

It is **not current architectural authority**. Current authority begins with `CLAUDE.md`, `docs/ROLE_IN_TONAL.md`, current architecture/state documents, and active experiment specifications.

The pre-R2 boundary defined Origami as the thing being built and versioned, Tlaloc/Blueprint/Tonal as external development/composition systems, and Tonal as an optional composition/reproducibility layer. That historical separation was useful for protecting Origami's standalone evidence boundaries, but its definition of Tonal is superseded by Architecture R2.

Key historical invariants that remain useful include:

- Origami owns its own representation/profile semantics and promotion evidence;
- portable/prompt-only claims must not silently rely on undeclared tools;
- perception, resolution, execution and verification are distinct;
- `FALSE_EXACT = 0`;
- inability to reveal a channel is `UNKNOWN`, not semantic absence;
- construction specification is not a compiled verified carrier;
- development-tool success is not prompt-only/native success.

For the exact historical version, use Git history before Architecture R2 migration.
