# Origami Implementation Gap Map — R0

Status: AUDIT
Baseline main: `7de2598cc050dfbf4c0649225a1b7b64f9b0f88f`
Audit parent: `315ff67bc32ee0d7455a1b7b1b3da96a49aca4bf`

This document maps the candidate intents in `ORIGAMI_INTENT_AUDIT_R0.md` to executable repository mechanisms. It deliberately distinguishes specification, executable core, OHF profile implementation, and missing machinery.

## Executable core discovered

The carrier-independent reference engine currently lives in `internal/reference/engine.go`.

It provides:

- finite entity sets;
- a finite domain;
- relations `requires`, `excludes`, `couples`;
- deterministic transition steps;
- contradiction detection;
- fixed-point, cycle and budget-exhaustion classification;
- exhaustive Cartesian enumeration of initial states;
- deterministic evidence counts.

This is genuine executable relational-state machinery. It is not merely documentation. However, it is intentionally a small R0 engine and must not be treated as the complete Origami machine.

## Code-level audit matrix

| Intent | Executable evidence | Classification | Gap |
|---|---|---|---|
| ORI-I001 relational state | `internal/reference/engine.go`: `Relation`, `State`, `step`, `contradiction` | KEEP + EXPAND | Only three relation types; relations are pair-shaped and state values are strings |
| ORI-I002 logical-time evolution | `Execute` repeatedly applies deterministic `step` and retains `Trace` | KEEP + EXPAND | No explicit `C_t`; rules are encoded as relation-type behavior; no general rule object |
| ORI-I003 observation != transition | Formal contract exists; reference engine emits classification/evidence after execution | RESTORE/BUILD | No first-class Observation Contract executor in the state transition engine; observation authority is not a core runtime type |
| ORI-I004 absence/unknown/inhibited/cancelled distinct | Semantic spec exists | BUILD | Reference engine domain does not enforce semantic distinctions; `active()` recognizes only literal `active`; cancellation is absent |
| ORI-I005 combinatorial space first-class | `Enumerate` performs full Cartesian enumeration; OHF dimensional package computes product-state capacity | KEEP EXPERIMENTALLY + GENERALIZE | Enumeration is brute-force and detached from higher-order relational semantics; dimensional implementation is OHF/profile-specific |
| ORI-I006 Fold/Unfold | Semantic specification exists | BUILD | No carrier-neutral Fold/Unfold implementation found in reference engine |
| ORI-I007 emergent trajectory/percept | Perceptual contracts and OHF laboratory mechanisms exist | KEEP CONTRACT + BUILD | General reference runtime does not execute temporal/emergent observation channels |
| ORI-I008 visual projection optional | Reference engine runs without visual/OHF dependencies | KEEP | This property is already demonstrated by architecture, but general core remains very small |
| ORI-I009 selective unfolding | Historical/OHF mechanisms and architecture support the idea | RESTORE + GENERALIZE | No core address graph/dependency-closure API found in `internal/reference` |
| ORI-I010 earned exactness | Deterministic evidence engine + OHF artifact verification | KEEP + CONSTRAIN | Core lacks explicit literal/residual accounting because those mechanisms currently live mainly in profile/carrier work |
| ORI-I011 layer separation | OHF lab has dedicated observation/evaluator/pipeline/failure/etc. packages | GENERALIZE CAREFULLY | Useful architecture exists but is profile-shaped; do not copy OHF implementation into core wholesale |
| ORI-I012 nominal != evidenced capacity | `internal/lab/dimensional`: nominal/safe reports, pairwise evidence, higher-order warning | KEEP AS EVIDENCE DISCIPLINE | Safe-set construction is pairwise; code explicitly warns higher-order gates remain required; not a general Origami capacity model |

## Important finding 1 — the current reference engine does confront combinatorics, but only at R0 scale

`Enumerate` constructs the Cartesian product of every declared entity and domain value. This is not a shortcut that avoids the state space. For `N` entities and domain cardinality `D`, it explicitly creates `D^N` initial states.

That behavior is aligned with the recovered principle that combinatorial state should not be silently discarded. It is also deliberately unscalable. The next design must not replace it with premature pruning that changes semantics. Instead it needs explicit bounded exploration, equivalence/closure rules and evidence that any reduction is semantics-preserving.

Decision: `KEEP AS ORACLE/REFERENCE`, not `KEEP AS PRODUCTION ALGORITHM`.

## Important finding 2 — the executable transition algebra is much poorer than the semantic language

Current executable transition behavior is hard-coded around:

- `requires`;
- `excludes`;
- `couples`;
- string states such as `active` / `inactive`.

The semantic documents already discuss determinate/superposed/coupled/observed states, interference/cancellation, Fold/Unfold, explicit observation and distinct absence/unknown/inhibited conditions. Therefore the executable reference engine is behind the specification.

Decision: `BUILD`, without pretending the missing semantics already execute.

## Important finding 3 — context and rules are not actually first-class yet

The Formal Core states:

`S_(t+1) = F(S_t, C_t, R)`

The reference implementation effectively executes:

`S_(t+1) = step(S_t, relations)`

There is no explicit runtime `Context` corresponding to `C_t`, and relation behavior is implemented by Go switch/if logic rather than a declared finite rule representation `R`.

This is a real implementation gap, not a documentation issue.

Decision: `RESTORE/BUILD` first-class Context and Rule structures before calling the engine a reference implementation of the full Formal Core equation.

## Important finding 4 — OHF already contains useful machinery, but it must be mined rather than promoted wholesale

`internal/lab/` is structurally richer than `internal/reference/`. It contains separate packages for dimensional state, observation, evaluation, pipeline, failure/regression, artifact handling, geometry, glyph calculus, MicroISA, native bundles and other OHF laboratory concerns.

The dimensional implementation is especially relevant:

- MICRO / MESO / MACRO scales;
- explicit dimensions and state alphabets;
- nominal product-state capacity;
- evidence statuses `UNKNOWN`, `PASS`, `FAIL`;
- pairwise orthogonality evidence;
- distinction between nominal and safe capacity;
- an explicit warning that higher-order compound/restoration gates are still required.

This code is valuable evidence for future core design, but its schema is explicitly `ohf.dimensional-registry.v1`. It must not silently become universal Origami semantics.

Decision: `MINE -> GENERALIZE -> TEST`, not `MOVE -> RENAME`.

## Important finding 5 — current pairwise capacity logic is intentionally insufficient for the recovered goal

`MaxPairwiseCompatibleSet` can identify a set whose every pair has PASS evidence. `Capacity` can calculate nominal and safe-looking capacity from promoted dimensions. But the resulting claim itself says higher-order compound/restoration gates remain required.

That is exactly where the restoration work should attack next: a set can be pairwise compatible while a triple or larger combination creates interference, ambiguity or a new emergent property.

Decision: create a core experiment in which all pairs are insufficient to determine the behavior of a higher-order combination.

## First implementation restoration sequence

Do not begin by rewriting OHF. The safest order is:

1. preserve the current R0 reference engine as a deterministic oracle;
2. add adversarial tests that expose missing semantics before changing the engine;
3. introduce first-class `Context` and declared `Rule` objects;
4. introduce explicit semantic values for absence/unknown/inhibited/cancelled rather than treating them as arbitrary strings;
5. implement carrier-neutral Fold/Unfold and explicit observation authority;
6. add addressability/dependency closure;
7. add higher-order relation/interaction support;
8. only then generalize selected dimensional/temporal mechanisms from OHF;
9. keep OHF as a profile that consumes the core rather than defining it.

## No-promotion decision

No code has been reclassified as production-ready by this audit. `main` remains unchanged. The audit branch now has enough implementation evidence to begin designing failure-first tests.

Next gate: `ADVERSARIAL_CORE_TEST_PLAN_R0`.
