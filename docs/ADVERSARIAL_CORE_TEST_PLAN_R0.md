# Origami Adversarial Core Test Plan — R0

Status: FAILURE-FIRST TEST DESIGN
Baseline main: `7de2598cc050dfbf4c0649225a1b7b64f9b0f88f`
Audit lineage:
- intent audit: `315ff67bc32ee0d7455a1b7b1b3da96a49aca4bf`
- implementation gap map: `27b045e2ec8347795059e43d9adfd7ae72accfe2`

Purpose: define tests before implementing the missing core semantics. A restored feature is not accepted merely because code exists; it must make the corresponding adversarial test pass without invalidating the R0 deterministic oracle.

## Test philosophy

Origami must not win by deleting the difficult part of the state space, hiding source information in an opaque literal/residual channel, moving semantic work into an observer, or declaring pairwise evidence sufficient for a higher-order claim.

The existing exhaustive R0 engine remains an oracle for the finite subset it already models. New machinery must agree with it on overlapping semantics.

## Gates

### ADV-001 — Higher-order interaction is not reducible to pairwise evidence

Construct entities A, B, C such that every pair is individually admissible, but the joint configuration A+B+C activates an interaction that does not exist in A+B, A+C or B+C.

Required result:
- the triple is represented and evaluated as a distinct higher-order condition;
- pairwise PASS evidence cannot promote the triple automatically;
- the evidence identifies the exact participating set;
- removing any one member removes the higher-order result.

Current expected status: FAIL/UNSUPPORTED.

Protects: ORI-I001, I005, I007, I012.

### ADV-002 — Context changes evolution without changing initial state

Run the same `S_0` and declared rules under contexts `C_alpha` and `C_beta`.

Required result:
- `S_0` is byte-for-byte identical in both runs;
- only `C_t` differs;
- traces diverge deterministically because context participates in `F(S_t, C_t, R)`;
- evidence records which context caused each transition.

Current expected status: FAIL/UNSUPPORTED because `Context` is not first-class in the reference engine.

Protects: ORI-I002.

### ADV-003 — Rules are data, not hidden Go control flow

Define two finite rule sets over the same state/context without recompiling the engine.

Required result:
- rules have stable IDs and finite declared predicates/effects;
- changing R changes the trace while code remains unchanged;
- every transition cites the rule(s) that fired;
- an undeclared rule cannot affect the run.

Current expected status: FAIL/UNSUPPORTED.

Protects: ORI-I001, I002, I010.

### ADV-004 — Absence, unknown, inhibited and cancelled are observably different

Create four otherwise equivalent cases where a target is respectively ABSENT, UNKNOWN, INHIBITED and CANCELLED.

Required result:
- they cannot collapse to one string/value or one observation result;
- UNKNOWN cannot be treated as false/absent;
- INHIBITED retains the fact that a potential transition was blocked;
- CANCELLED records competing contributions and cancellation evidence;
- observation reports the semantic distinction without mutating the pre-observation trace.

Current expected status: FAIL/UNSUPPORTED.

Protects: ORI-I003, I004, I007.

### ADV-005 — Observation is not transition

Execute a trajectory once, then apply two different Observation Contracts to the same immutable pre-observation evidence.

Required result:
- transition trace hash is identical before both observations;
- observations may expose/classify different aspects;
- neither observer may rewrite prior states;
- observation authority and observation result are recorded separately from transition provenance.

Current expected status: FAIL/PARTIAL.

Protects: ORI-I003, I007, I010.

### ADV-006 — Fold preserves unresolved alternatives

Create a finite state containing multiple unresolved alternatives and Fold it.

Required result:
- Fold does not choose one alternative merely to reduce representation size;
- Unfold restores the same alternative set and relations;
- Fold carries a verifiable commitment/hash to the folded semantic content;
- observation after Fold must obey the same authority rules as observation before Fold.

Current expected status: FAIL/UNSUPPORTED.

Protects: ORI-I005, I006, I010.

### ADV-007 — Selective Unfold is addressable and dependency-complete

Fold a graph with at least two mostly independent regions plus one cross-region dependency. Request one address.

Required result:
- only the requested semantic closure is unfolded;
- all dependencies required to interpret it are included;
- unrelated regions remain folded;
- selective result is equivalent to selecting the same region after full Unfold;
- evidence reports addresses and dependency closure.

Current expected status: FAIL/UNSUPPORTED in carrier-neutral core.

Protects: ORI-I006, I009, I010.

### ADV-008 — Temporal latent property exists only across a trajectory

Construct a trajectory where no individual state/frame contains property P, but an Observation Contract over the ordered sequence detects P from temporal relations.

Required result:
- inspecting any single state cannot produce P;
- reordering the same states changes/destroys P;
- the temporal observer cites the sequence/time relations that establish P;
- P is not pre-stored as a literal label in any state.

Current expected status: FAIL/UNSUPPORTED in general core.

Protects: ORI-I002, I003, I007.

### ADV-009 — Interference/cancellation cannot be replaced by priority selection

Create two independently valid contributions whose combination cancels or inhibits an outcome.

Required result:
- both contributions remain in evidence;
- the result records cancellation/interference rather than selecting a winner by arbitrary ordering;
- reversing declaration order does not change the semantic result;
- removing either contribution changes the result.

Current expected status: FAIL/UNSUPPORTED.

Protects: ORI-I001, I004, I007.

### ADV-010 — No semantic answer smuggling

For every semantic-core adversarial fixture, scan the fixture and intermediate representation for the expected final answer/classification encoded directly as a dedicated literal field whose only purpose is to make the test pass.

Required result:
- expected observation must be derivable from state, context, rules, trajectory and observation contract;
- no fixture field may simply contain the answer as an oracle hint;
- opaque residual/literal payloads are not accepted as proof of semantic computation.

Current expected status: test harness to be built before feature implementation.

Protects: ORI-I007, I010.

### ADV-011 — Exact recovery is measured separately from semantic computation

Use a fixture with an exact byte payload plus a semantic state representation.

Required result:
- exact recovery may use explicitly budgeted residual/literal mechanisms when the profile permits them;
- semantic queries must succeed with the original byte source unavailable;
- removing residual bytes may break exact byte recovery but must not silently change already-represented semantic state;
- evidence reports literal/residual contribution separately from machine-derived semantic contribution.

Current expected status: FAIL/UNSUPPORTED at core integration level.

Protects: ORI-I010 and prevents conflating compression/recovery with Origami computation.

### ADV-012 — Reduced exploration must prove equivalence to exhaustive oracle

For state spaces small enough for R0 exhaustive enumeration, run both exhaustive enumeration and any future optimized explorer.

Required result:
- same reachable semantic classes;
- same contradiction/fixed-point/cycle outcomes under overlapping semantics;
- any merged/pruned states carry a proof/evidence reason;
- no optimization is accepted solely because it is faster.

Current expected status: BLOCKED until optimized exploration exists; this is a permanent future regression gate.

Protects: ORI-I005, I010.

## Anti-cheat invariants

The adversarial suite must reject a claimed PASS when any of these occurs:

1. a final semantic answer is stored directly in the fixture instead of derived;
2. an observer/VLM/LLM performs a state transition that the core claims to have performed;
3. unresolved alternatives are discarded before an authorized observation;
4. UNKNOWN is silently coerced to ABSENT/FALSE;
5. higher-order behavior is inferred solely from pairwise PASS evidence;
6. selective unfolding reads the entire source and merely hides unrelated output;
7. exact recovery success is presented as semantic-computation success;
8. a residual/literal channel is effectively unbounded and becomes the primary semantic representation;
9. declaration order changes an order-independent semantic interaction;
10. an optimized explorer differs from the exhaustive oracle on a finite overlap without explicit semantic justification.

## Implementation order forced by the tests

Phase A — test harness only:
- encode ADV-001 through ADV-005 as failing/unsupported fixtures;
- preserve current reference-engine tests unchanged;
- add explicit expected-failure/unsupported reporting rather than weakening assertions.

Phase B — Formal Core restoration:
- first-class Context;
- first-class finite Rules with provenance;
- semantic value/status model;
- explicit Observation Contract runtime.

Phase C — compositional state:
- higher-order interactions;
- interference/cancellation;
- carrier-neutral Fold/Unfold;
- address graph and selective dependency closure.

Phase D — temporal/emergent semantics:
- trajectory observations;
- temporal latent properties;
- higher-order compound evidence.

Phase E — optimization without semantic loss:
- bounded/optimized exploration;
- equivalence gate against exhaustive R0 oracle;
- only after equivalence evidence may optimized exploration replace exhaustive enumeration for larger spaces.

## Promotion rule

A feature cannot move from `BUILD/RESTORE` to `KEEP` merely because its unit tests pass. Promotion requires:

- corresponding ADV gate PASS;
- deterministic replay PASS;
- provenance/evidence PASS;
- no anti-cheat invariant violation;
- existing R0 oracle tests PASS;
- explicit change-control record.

No semantic implementation is changed by this document. The next action is to materialize Phase A fixtures/tests and confirm that the current engine fails for the intended reasons.
