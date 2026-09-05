# Anti-prior tests for Origami carrier experiments

## Purpose

Compact visual or multimodal transport can appear stronger than it is when a receiving model reconstructs plausible content from linguistic or semantic priors.

Origami carrier claims should therefore include tests that separate genuine representation recovery from prior-assisted completion.

## Test families

Candidate controls include:

1. arbitrary identifiers with no natural-language semantics;
2. randomized symbol-to-meaning assignments;
3. assignments permuted independently between runs;
4. semantically corrupted but mechanically valid sequences;
5. exact roundtrip reconstruction after model recovery;
6. downstream deterministic execution from recovered semantics.

## Example

Instead of testing only predictable language, test structures such as:

```text
REF_83
CALL_12
ARG_42
ARG_91
VERIFY_3
RETURN_6
```

The model should not receive hidden ground-truth mapping artifacts.

## Required reporting

Report separately:

- exact recovered structure;
- unknown/abstained structure;
- false-known claims;
- downstream execution success;
- representation size/cost;
- performance change after semantic priors are removed.

A large drop under anti-prior controls is a finding, not a reason to hide the control.
