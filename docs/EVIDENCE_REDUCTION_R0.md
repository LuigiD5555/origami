# Evidence Reduction R0

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`  
Contract: `origami.evidence-reduction.r0`

Evidence Reduction R0 turns probabilistic model/swarm interpretations into a reproducible Origami state without pretending that the model itself is deterministic.

The central rule is:

> Models propose. Origami resolves evidence and reduces the proposals deterministically.

## Project boundary

Origami does **not** absorb Tlaloc's swarm orchestration. Tlaloc or another model system may decide how many agents to run, what prompts to use and how to explore ambiguity. Origami owns the representation and verification boundary:

```text
addressable Origami memory
        |
        v
bounded ContextPacket
        |
        v
Tlaloc / external swarm
  probabilistic proposals
        |
        v
ClaimProposal[]
        |
        v
Origami evidence resolver
        |
        v
CID/hash/fidelity/verified checks
        |
        v
deterministic reducer
        |
        v
CanonicalClaim[]
```

This preserves the existing project rule: Tlaloc coordinates work; Origami defines how memory, state and exactness are interpreted and verified.

## Proposal contract

Each agent emits a small structured proposal instead of free-form authority:

```json
{
  "id": "reader-a-17",
  "claim_key": "blt.dynamic-patching",
  "proposition": "BLT can use entropy-driven dynamic patches.",
  "agent_id": "reader-a",
  "position": "SUPPORT",
  "confidence": 0.91,
  "evidence": [
    {
      "address": "ohf://paper/evidence/blt-support",
      "fidelity": "evidence",
      "cid": "sha256:...",
      "source_hash": "sha256:..."
    }
  ]
}
```

Allowed positions are:

```text
SUPPORT
OPPOSE
UNKNOWN
```

Agent confidence is retained as diagnostic information only. It never substitutes for evidence.

## Evidence acceptance

An evidence reference counts only when Origami can prove all applicable conditions:

1. the address exists;
2. the declared CID matches the addressed node when supplied;
3. the declared source hash matches provenance when supplied;
4. the requested fidelity exists;
5. that content level is explicitly marked `Verified=true`.

There is no silent fallback from an unavailable evidence fidelity to a weaker summary. A proposal that asks for evidence must actually point to evidence.

Duplicate agents do not manufacture consensus. Evidence is deduplicated by content identity/fidelity/source identity, so ten agents repeating the same source still contribute one evidence anchor.

## Canonical states

For each `claim_key`, Origami counts unique verified evidence anchors and independent sources under a finite policy.

```text
support passes, opposition does not -> VERIFIED
opposition passes, support does not -> REJECTED
both pass                         -> CONFLICT
neither passes                    -> UNKNOWN
```

`CONFLICT` is preserved instead of being hidden by majority voting. The caller may then spend more compute, unfold more context or ask additional agents to inspect the conflict.

`UNKNOWN` is a valid state and is preferred to an invented answer.

## Exactness

`VERIFIED` does not automatically mean exact reproduction. Exactness is earned separately.

```text
verified evidence-level source -> VERIFIED + NOT_CLAIMED
verified exact-level source    -> VERIFIED + VERIFIED_EXACT
```

Conflicted, rejected or unknown claims never retain an exactness claim. `FALSE_EXACT=0` remains mandatory.

## Determinism

The reducer sorts claim keys, proposal IDs and evidence references before reduction. The same set of proposals and the same memory state therefore produce the same:

- canonical status;
- accepted/rejected evidence sets;
- canonical claim ID;
- reduction digest.

Proposal arrival order is not semantically meaningful.

## CLI

`cmd/origami-reduce` is a model-neutral bridge. It accepts a JSON object with the carrier memory nodes and external proposals:

```json
{
  "carrier_id": "book-a",
  "nodes": [],
  "policy": {
    "min_verified_evidence": 1,
    "min_independent_sources": 1
  },
  "proposals": []
}
```

Run:

```bash
go run ./cmd/origami-reduce -in input.json -out result.json
```

or stream through stdin/stdout:

```bash
cat input.json | go run ./cmd/origami-reduce
```

This is intentionally simple enough for a local 4K-context model, Tlaloc, shell scripts or another agent framework to consume.

## User-facing meaning

A user can ask a question against a very large document universe while the model only sees a bounded `ContextPacket`. The model/swarm explores candidate interpretations. Origami then reopens the cited addresses and verifies what the candidates actually relied on.

The user-facing result can therefore say:

```text
Claim: BLT can use entropy-driven dynamic patches.
State: VERIFIED
Evidence: ohf://paper/evidence/blt-support
Exactness: NOT_CLAIMED
```

If two verified regions genuinely disagree, the result is `CONFLICT`, not a fabricated certainty.

## What R0 does not claim

R0 does not prove that an arbitrary LLM will find the right addresses, understand every PDF/image correctly or answer every question. It establishes the deterministic authority boundary required to test those model-facing behaviors safely.

Direct PDF extraction, OCR/VLM perception and Tlaloc swarm search remain separate ingestion/orchestration concerns. Their outputs can be converted into addressable Origami nodes and then evaluated by this reducer.
