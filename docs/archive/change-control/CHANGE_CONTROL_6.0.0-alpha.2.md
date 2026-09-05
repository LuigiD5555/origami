# Change control — Origami 6.0.0-alpha.2

Date: 2026-08-28
Status: `EXPERIMENTAL_CONTRACT`

## Component changed

Origami representation semantics — perceptual channel layer.

## Before

`6.0.0-alpha.1` defined the project split and coherent-state profile but treated visual projection mostly as a general concept without a formal contract for interference, depth, temporal-only perception, or emergent visual effects.

## After

`6.0.0-alpha.2` adds `origami.perceptual-channels.r0`, defining:

- channel families: spatial, interference, depth, temporal, emergent;
- Temporal Latent Image;
- explicit vs emergent representation;
- perceptual conditions distinct from coherent-state kinds;
- seven perceptual operations;
- nine invariants;
- machine-readable contract.

## Key invariant added

> A static projection is not assumed to contain every valid percept of an Origami representation.

In particular, failure to reveal a motion-bound structure in one static frame does not prove absence.

## Files changed

- `VERSION`
- `README.md`
- `CHANGELOG.md`
- `docs/STATE_SEMANTICS_R0.md`
- added `docs/PERCEPTUAL_CHANNELS_R0.md`
- added `spec/PERCEPTUAL_CHANNELS_R0.json`
- added this change-control record.

## Evidence / verification

- machine-readable contract parses as valid JSON;
- contract ID/version consistency checked;
- required `KINETIC_REVEAL` operation present;
- `LATENT_IS_NOT_ABSENT` invariant present;
- terminology scan keeps Tlaloc ownership outside Origami;
- no claim of complete runtime implementation is made.

## Regressions

No coherent-state semantic rule was removed. Perceptual conditions are explicitly orthogonal to coherent-state kinds to avoid ontology drift.

## Downstream impact

Tlaloc `6.0.0-alpha.6` remains architecturally compatible because Origami is optional and accessed through contracts. However, its current Origami-facing skill/reference profile does not yet generate or evaluate the new perceptual operations. That integration is deferred to a Tlaloc-side revision rather than hidden inside this Origami release.

## Promotion decision

`EXPERIMENTAL_CONTRACT`.

Promotion to an implemented semantic profile requires at minimum:

1. exact-base integration with the pre-split Origami implementation;
2. reference fixtures for interference, depth and Temporal Latent Image behavior;
3. deterministic tests for supported operations;
4. explicit unsupported-operation behavior for everything not implemented.
