# Master Prompt Generation R0

`generated/MASTER_PROMPT.md` is an operational artifact, not an independently authored source of truth.

## Rule

Until the Go prompt generator is implemented/promoted, the checked-in R0 prompt is a **bootstrap reference artifact**. Any semantic change to the receiver contract must first be reflected in the canonical Origami protocol/specification and tests. Once `cmd/promptgen` exists, CI must regenerate `generated/MASTER_PROMPT.md` and fail on divergence.

Target flow:

```text
internal protocol/spec
  -> tests
  -> cmd/promptgen
  -> generated/MASTER_PROMPT.md
  -> sha256 manifest
```

Do not maintain multiple hand-edited operational Master Prompts for different models. Model-specific adapters may wrap transport details, but the Origami receiver semantics remain common.

## Required generated sections

The generator must cover:

- authority and anti-contamination;
- W as active interface budget, not storage cap;
- Native receiver rules;
- Computational tool semantics;
- selective navigation policy;
- Dependency Folding guards/proofs;
- context recycling;
- exactness/UNKNOWN policy;
- evidence/status response contract.

## Gate

A future PR that changes receiver semantics without regenerating the prompt is invalid once generation is promoted to required CI.
