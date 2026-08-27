# Development and change control

OHF evolves quickly enough that repository state must be more trustworthy than conversation history.

## Installation

```bash
./install.sh --check-only
./install.sh
```

Development checks:

```bash
make test
make vet
make check
```

## Change workflow

Every material change should follow this sequence:

1. open a Change Record under `changes/`;
2. name the primary component;
3. record the reason/hypothesis;
4. determine downstream impact;
5. preserve the before-state identity;
6. modify only the declared scope;
7. run cheap local gates first;
8. run affected regressions;
9. produce evidence;
10. compare before/after;
11. classify the result;
12. create regressions for reproducible new failures;
13. update component state;
14. promote only through a separate promotion decision.

## Allowed states

Component/change status should use explicit values rather than informal prose:

```text
UNIMPLEMENTED
EXPERIMENTAL
BLOCKED
FAILED
EXPERIMENTAL_VALIDATED
PROMOTION_CANDIDATE
PROMOTED
REJECTED
SUPERSEDED
```

Experiment results:

```text
PASS
FAIL
INCONCLUSIVE
BLOCKED
```

## Impact closure

Do not rerun the entire system when a component cannot affect unrelated upstream layers.

Example: changing a visual renderer normally requires renderer/perception/visual regressions, but does not automatically require canonical-source or grammar-discovery tests.

A Generative IR change has a much larger downstream closure and should trigger execution/verification/restoration tests.

The goal is rigorous **and** efficient testing.

## Repository hygiene

Do not commit:

- generated `runs/` contents;
- external-model raw responses unless they are converted into a curated regression fixture;
- historical R3.10 SDK bytes until their provenance is known;
- local binaries;
- secrets or API keys.

Do commit:

- experiment specifications;
- deterministic source code;
- minimal fixtures;
- component states;
- change records;
- reproducible regression definitions;
- summarized evidence needed to justify a protocol claim.

## Module path

The canonical Go module for this repository is:

```text
github.com/LuigiD5555/origami
```
