# Dependency Folding R0

Status: experimental semantic optimization.

## Purpose

Dependency Folding reduces cognitive/exposure depth without pretending that canonical computational depth disappeared.

```text
Canonical graph
  -> Dependency Folding
  -> hierarchical compiled graph
  -> SuperIndex / selective unfold
```

A compiled edge is a derived view, never a replacement for canonical state or rules.

## Invariants

1. Canonical derivations remain addressable.
2. A compiled edge preserves every condition, precondition and exception that can change its result.
3. Every compiled edge carries a deterministic proof commitment to its ordered canonical derivation.
4. `VERIFY` can reopen the canonical derivation and validate the commitment.
5. Folding must not turn unknown or unresolved dynamic state into a determinate result.
6. Undeclared rules still cannot affect execution.
7. Selective execution may use a compiled edge only when its guards are satisfied or explicitly exposed as unresolved.
8. A semantic result obtained from folding must equal canonical execution for the same inputs.

## Three views

- **Canonical** — exact derivation and source of truth.
- **Navigation** — hierarchical/skip structure used to locate relevant dependencies cheaply.
- **Compiled** — a folded relation used when its semantic guards permit direct evaluation.

Physical depth and cognitive depth are therefore intentionally different quantities.

## Verification

`proof_ref = SHA256(canonical ordered steps + rules + normalized guards)`.

The initial reference implementation stores the canonical steps with the derived edge for testing. Production storage may replace that field with independently addressable canonical ranges, provided the proof commitment and exact verification semantics remain equivalent.

## Working-window accounting

Experiments must distinguish at least:

- REQUESTED
- TOUCHED / VISITED
- UNFOLDED
- EXPOSED
- active bytes
- token-equivalent exposure
- physical dependency depth
- cognitive/folded depth
- verification depth

A successful fold is not permission to count only the final answer. Guards and proof material actually exposed to the answering model count toward the working window.

## Failure rule

If a compiled edge cannot preserve a dependency's semantics, do not fold across that boundary. Expose or recursively refine only the necessary region. Never silently load the whole graph and report it as selective unfolding.
