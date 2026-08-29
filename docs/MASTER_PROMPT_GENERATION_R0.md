# Master Prompt Generation and Promotion R0

`generated/MASTER_PROMPT.md` is a promoted/exported receiver artifact. It is not the source of Origami semantics and it is not expected to carry the carrier's private symbol system.

## Split of responsibility

```text
ORIGAMI semantics + receiver contract
              |
              v
TLALOC swarm / Tlaloque experiments
  generate + mutate + compare candidates
              |
              v
receiver candidate package
  prompt + boot strategy + Rosetta strategy + micro-agent IR
              |
              v
ORIGAMI validation / promotion
              |
              +--> generated/MASTER_PROMPT.md
              +--> carrier BOOT/ROSETTA/PROGRAM artifacts
              +--> provenance + hashes
```

Tlaloc may discover a better way to make models bootstrap and operate Origami. Tlaloc does not redefine what Origami states, symbols, transitions or verification claims mean. Origami validates a candidate against its own contracts and stores only promoted receiver artifacts.

## Universal prompt rule

The Master Prompt should remain as small and carrier-independent as practical. It may explain the universal bootstrap discipline, isolation, active-context budget, UNKNOWN policy and stage separation. It must **not** become a second copy of the carrier's Rosetta table or micro-program.

A carrier must be able to change its local glyph/symbol assignments without requiring the external Master Prompt to be rewritten.

## Current bootstrap status

Until the automated cross-repository promotion pipeline is implemented, `generated/MASTER_PROMPT.md` is the checked-in R0 reference candidate. Changes to it must be accompanied by receiver-contract tests and provenance in the relevant PR.

The target deterministic export flow inside Origami is:

```text
promoted receiver artifact
  -> validate against spec/HYBRID_RECEIVER_R0.json
  -> deterministic exporter
  -> generated/MASTER_PROMPT.md
  -> manifest + sha256
```

The *search/optimization* of the prompt belongs to Tlaloc. The *deterministic storage/export* of the promoted artifact belongs to Origami.

## Candidate fitness

Tlaloc should rank candidates with more than answer correctness. At minimum the campaign should measure:

- bootstrap success;
- correct use of carrier-local Rosetta mappings;
- navigation correctness;
- answer correctness;
- evidence integrity;
- UNKNOWN accuracy;
- false-exact count;
- peak active token-equivalent;
- model/tool cost;
- contamination.

Any candidate with false exactness or contamination is ineligible for promotion regardless of aggregate score.

## Required cross-model regression

Keep the Master Prompt fixed and permute the physical symbol assignments in independent carriers. A promoted receiver must follow each carrier's `BOOT -> ROSETTA -> PROGRAM` path rather than relying on memorized glyph meanings.

## Gate

Once automated export is promoted to required CI, a receiver-contract change that does not regenerate and re-hash the stored receiver artifacts is invalid.
