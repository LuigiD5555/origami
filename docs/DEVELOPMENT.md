# Development and change control

Origami evolves through multiple semantic and research tracks. Repository state must be more trustworthy than conversation history.

## Basic local checks

```bash
go test ./...
go vet ./...
```

The current Go code primarily exercises the OHF laboratory track. A passing OHF suite does not by itself prove runtime support for every Origami 6.x semantic operation.

## Change workflow

For every material change:

1. identify whether it belongs to core Origami semantics/runtime or a specific profile/research track such as OHF;
2. open/update a change record under `changes/`;
3. record component, reason, before/after state and downstream impact;
4. run the smallest valid impact closure first;
5. add regressions for reproducible failures;
6. update `state/ORIGAMI_STATE.json` for root changes and `state/components/` for OHF/component changes;
7. promote only through an explicit decision.

## Semantic ownership

Do not put Tlaloc responsibilities into Origami. Behavior compilation, Tlaloque coordination, target-model orchestration and behavioral-artifact promotion are Tlaloc-owned.

Do not silently turn OHF-specific constraints into universal Origami laws. If a capability is being generalized from OHF into core Origami, record that ownership transition explicitly.

## Dimensional/perceptual changes

Record visual scope, physical states, controlled render mutation, collision results, orthogonality evidence and the difference between nominal capacity and evidence-gated safe capacity. Pairwise success does not establish higher-order independence.

Perceptual availability must remain distinct from coherent-state kind. Failure to reveal a latent/motion-bound percept in one static frame is not proof of absence.

## Repository hygiene

Do not commit generated runs, raw external-model responses, local binaries, secrets or unverified historical SDK bytes. Commit deterministic source, minimal fixtures, semantic contracts, state/change records and evidence summaries required to justify claims.

External campaigns should be incremental, cached, low-concurrency and respectful of rate limits/Retry-After.
