# Origami Formal Core R0

Status: experimental specification.

## Purpose

Origami R0 treats the language as a system for representing and transforming relational state. Visual projection is optional. A claim about an Origami state is not accepted merely because it is plausible or visually suggestive; non-direct claims require a finite observation contract capable of returning PASS or FAIL.

## Minimal model

At logical time `t`, an Origami machine has finite state `S_t`, declared context `C_t`, and rules `R`.

```text
S_(t+1) = F(S_t, C_t, R)
```

The default R0 transition model is deterministic. A future profile may explicitly define other semantics, but nondeterminism must never be introduced accidentally by an implementation detail.

Observation is separate from transition:

```text
P = G(S_[0:n], O, tau)
```

where `O` is a declared observer/procedure and `tau` is a finite authorized trajectory.

## Falsifiability rule

A non-direct property must provide an Observation Contract before execution. The contract declares preconditions, observer, trajectory, finite budget, expected predicate, and failure semantics.

`LATENT_IS_NOT_ABSENT` therefore means only that failure to reveal a property in an unauthorized or insufficient observation is not proof of absence. It does **not** make a latent claim immune to failure. If the declared observation is valid, its finite budget is exhausted, and the expected predicate has not been satisfied, the result is `FAIL`.

## Outcomes

- `PASS` — the declared predicate was satisfied within budget.
- `FAIL` — valid experiment, exhausted budget, expected predicate not satisfied.
- `INVALID_CONTRACT` — experiment cannot establish the claim because its contract is malformed or preconditions do not hold.
- `UNSUPPORTED` — semantics are declared but the selected runtime does not implement the requested operation.

These states must not be collapsed. In particular, `UNSUPPORTED` is not `FAIL`, and `INVALID_CONTRACT` is not evidence against the semantic claim.

## Experiment discipline

Experiments declare their hypothesis and failure conditions before execution. Evidence is stored separately from the hypothesis. A result may promote a claim only when the declared gate passes.

The first experiment is `experiments/EXP-001-relational-state/experiment.json`. It deliberately starts without visual channels, VLMs, or OHF carrier constraints so the relational machine can be tested independently of presentation.

## Non-goals

R0 does not claim physical quantum computation, biological equivalence, superiority to conventional solvers, or that every Origami semantic operation has executable runtime support.
