# Origami Intent Restoration Audit — R0

Status: AUDIT / NOT YET A FROZEN LANGUAGE SPECIFICATION
Baseline: `7de2598cc050dfbf4c0649225a1b7b64f9b0f88f` (`6.0.0-alpha.3`)
Branch: `audit/origami-intent-restoration`

## Purpose

This document is the restoration ledger for Origami. It exists to prevent implementation convenience, an experimental carrier, or a successful narrow benchmark from silently redefining the project.

The audit records each recovered intent as:

`INTENT -> MOTIVATION -> MECHANISM -> CURRENT IMPLEMENTATION -> GAP -> RESTORATION -> REQUIRED TEST`

No item becomes a frozen Origami invariant merely by appearing here. Promotion requires evidence and an explicit decision after the historical and implementation audit is complete.

## Project boundary recovered so far

Origami is a visual/computational representation and state-machine language. Its core subject is relational state and its transformation through logical time. Visual projection is a possible observation/projection mechanism, not the identity of the machine. Tlaloc is external and owns work/orchestration behavior. OHF is an Origami research profile/carrier and laboratory, not the whole language.

## Audit matrix R0

| ID | Intent | Motivation | Current implementation | Gap / risk | Restoration direction | Required evidence |
|---|---|---|---|---|---|---|
| ORI-I001 | Relational state is first-class | Meaning can belong to relations, not isolated symbols | Formal Core R0 and EXP-001 | Executable generality is still narrow | Preserve and expand independently of carrier | Higher-order relational experiments |
| ORI-I002 | State evolves through logical time | Origami is a machine, not a static encoding | `S_(t+1)=F(S_t,C_t,R)` formal contract | Runtime coverage is incomplete | Reference transition runtime with explicit rules/context | Deterministic trajectory/replay tests |
| ORI-I003 | Observation is distinct from transition | Observing/resolving must not silently mutate semantics | Formal Core separates `F` and `G`; `OBSERVE` is explicit | Some profile tooling can blur perception/resolution/execution | Enforce observation authority at core boundary | Negative tests proving TRANSFORM cannot resolve |
| ORI-I004 | Absence, unknown, inhibition and cancellation differ | Missing information can itself change meaning | State semantics explicitly distinguishes them | Project state only freezes `LATENT_IS_NOT_ABSENT`; absence semantics are under-specified | Promote explicit absence-family algebra only after experiments | Presence/absence counterfactual matrix |
| ORI-I005 | Combinatorial space is first-class | Higher-order combinations can carry emergent meaning; the goal is to structure complexity rather than erase it | OHF has bounded multidimensional/PSS work; core has superposed/coupled states | No general core contract currently prevents optimization from collapsing meaningful alternatives/interactions | Define bounded combinatorial exploration and preservation semantics | Tests where pairwise evidence is insufficient but higher-order interaction is meaningful |
| ORI-I006 | Fold/Unfold are semantic operations | Origami needs selective construction/constraining of alternatives, not only encode/decode | STATE_SEMANTICS_R0 defines both | Semantics remain minimal and OHF historically carries much of selective-unfold machinery | Generalize Fold/Unfold independently of OHF carrier | Round-trip, selective-unfold and policy-bound tests |
| ORI-I007 | Emergence may belong to interaction/trajectory | Some information exists only across interference, depth, motion or time | Perceptual Channels R0 contract | Several operations remain unsupported | Keep contract; implement one channel at a time with evidence | Finite observation contracts with PASS/FAIL |
| ORI-I008 | Visual projection is optional | The machine must not collapse into an image format | Formal Core says visual projection optional | OHF is the main executable track and can dominate project identity | Build carrier-independent reference core | Core suite must run without renderer/VLM/OHF |
| ORI-I009 | Selective unfolding is fundamental | A representation should expose only the dependency closure needed for an observation/query | Strongly present in OHF architecture | Not yet a general Origami core primitive | Recover addressability/dependency closure from OHF into a carrier-neutral contract | Query requiring partial rather than full expansion |
| ORI-I010 | Exactness must be earned | Plausibility must never be reported as exact reconstruction | OHF `FALSE_EXACT=0`, SHA-256 verification, Formal Core outcomes | Residual/literal storage can become an escape route that wins tests without meaningful folding | Add accounting for literal/residual burden and forbid it as hidden success criterion | Adversarial tests and explicit residual ratio/evidence |
| ORI-I011 | Perception, resolution, execution and verification are different layers | Success in one layer does not prove success in another | OHF invariant explicitly separates them | Core/project boundary does not yet encode the full separation mechanically | Generalize useful separation without importing OHF-specific machinery wholesale | Layer-specific failure injection tests |
| ORI-I012 | Mathematical capacity differs from safe/evidenced capacity | Nominal combinatorics are not proof of usable independent dimensions | OHF `NOMINAL_CAPACITY_NE_SAFE_CAPACITY`, orthogonality gates | Mostly OHF-specific today | Retain as experimental discipline; generalize only when core multidimensional semantics are defined | Capacity curve with failure boundary and higher-order gates |

## Known implementation asymmetry

The repository currently has a stronger semantic specification than general-purpose execution support. The coherent-state and perceptual-channel layers are experimental contracts. OHF R3.10-LAB is the broadest executable research track, but that does not grant OHF mechanisms universal-language status.

This asymmetry is intentional during restoration: unsupported operations must remain `UNSUPPORTED`; they must not be simulated by a shortcut and reported as implemented.

## Anti-shortcut audit rules

During restoration, an experiment is not evidence for the intended Origami mechanism when its success can be explained primarily by one of these substitutions:

1. storing the answer directly;
2. copying most of the source into `LITERAL` or equivalent storage;
3. moving most information into an unrestricted residual;
4. using a screenshot/render as if that alone were a semantic Fold;
5. asking an LLM/VLM to perform the core state transition or recover relations that the Origami machine itself did not represent;
6. collapsing semantically distinct absence/unknown/inhibited/cancelled states;
7. pruning alternatives before a declared Fold/Observe policy authorizes resolution;
8. claiming higher-order capacity from pairwise-only evidence;
9. reporting an unsupported perceptual operation as a negative semantic result rather than `UNSUPPORTED`.

Literal and residual mechanisms are not forbidden. They must be visible, measured, justified, and prevented from becoming an unbounded escape hatch.

## Layer map to preserve during restoration

```text
ORIGAMI CORE
  state + relations + context + rules
  transition + Fold/Unfold + observation contracts
  addressability + dependency closure
  bounded combinatorial interaction
            |
            v
PROJECTIONS / OBSERVATION CHANNELS
  spatial / interference / depth / temporal / emergent
            |
            v
PROFILES / CARRIERS
  OHF and future profiles

TLALOC (external)
  orchestration / behavior / Tlaloque coordination / model-facing work
```

## Promotion rule

Do not freeze the language from this R0 audit. Before freezing:

- finish historical recovery and attach provenance/reasoning to each candidate invariant;
- map candidate invariants to actual code/spec/tests;
- identify contradictory or convenience-driven behavior;
- design adversarial tests before implementation changes;
- classify each item `KEEP`, `GENERALIZE`, `RESTORE`, `BUILD`, `CONSTRAIN`, `REMOVE`, or `REJECT`;
- record downstream impact on OHF, Tlaloc and Tonal;
- promote only evidence-backed invariants.

## Current decision

`AUDIT_CONTINUES`.

The existing `main` remains the promoted implementation baseline. This branch is evidence-gathering/restoration work and does not yet redefine canonical Origami semantics.
