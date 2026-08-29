# Change control — Origami 6.0.0-alpha.3

Date: 2026-08-28  
Status: `PROMOTED_REPOSITORY_HYGIENE`

## Component changed

Project hierarchy, current-state documentation, machine-readable root state and ownership labeling after integration of the OHF R3.10-LAB workbench.

## Before

Origami `6.0.0-alpha.2` correctly defined the Tlaloc/Origami split and perceptual-channel contracts. A later merge brought the executable OHF R3.10-LAB tree into `main`, but several root/current files again described OHF/R3.10 as the state or architecture of all Origami. `state/ORIGAMI_STATE.json` still used `ohf.origami-state.v1` and had no knowledge of the 6.x semantic contracts.

## After

- Origami advances to `6.0.0-alpha.3` without changing alpha.1/alpha.2 semantic laws.
- root project state uses `origami.state.v2` and separates semantic contracts from the OHF research-track substate;
- README, current state, architecture, roadmap, boundary and development/contribution docs consistently represent `Origami > OHF`;
- OHF R3.10-LAB code, DVR/PSS, Glyph Calculus, SAFE_MICRO_ISA, Context SIMD and reversible tooling lifecycle are preserved;
- the installer is documented accurately as OHF laboratory tooling rather than the complete Origami runtime;
- historical evidence remains intact.

## Evidence / tests

Impact closure is documentation/evolution-control only; no Go runtime source is modified by this revision.

- `VERSION` / README / root-state version consistency checked;
- `state/ORIGAMI_STATE.json` and `changes/CHG-ORIGAMI-0001.json` parse as valid JSON;
- project hierarchy scan confirms OHF is labeled as a nested research track in active root/current docs;
- Tlaloc ownership remains outside Origami;
- existing semantic contract files are unchanged;
- no historical R2/R3.10 change record is deleted.

## Regression risk

Low for executable code; medium for consumers that parsed the old root-state schema directly. Search found no code references to `ohf.origami-state.v1`, so the root-state schema can advance to `origami.state.v2` without a known executable consumer break.

## Promotion decision

Promoted as repository/state reconciliation. Runtime semantic promotion is explicitly out of scope.
