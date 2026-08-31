# Origami 6.0.0-alpha.15 — Experimental Candidate Builder R0

## Scope

Alpha.15 adds an Origami-owned compiler for bounded experimental representation mutations requested by external development systems such as Tlaloc.

The development path is:

```text
failure evidence / adaptive search
        -> declarative CandidateSpec
        -> Origami candidate compiler
        -> experimental PNG
        -> exact-plane verification
        -> held-out VLM trials
```

## Implemented

- `origami.experimental-candidate.r0` declarative mutation input;
- `origami.experimental-candidate-build.r0` contract;
- `origami-candidate-build` CLI and managed install/uninstall lifecycle;
- direct compatibility with Tlaloc's `TLALOC_*` candidate build hook;
- deterministic parent PNG SHA binding when supplied;
- parent temporal carrier mechanical decode before mutation;
- visual-only mutation on the semantic/temporal plane;
- post-build mechanical decode and exact `TemporalProgram` SHA equality requirement;
- frozen 640x640 / exactly-8192-byte output envelope;
- explicit failure if mutation output cannot fit the frozen carrier;
- capability negotiation through `origami-candidate-build capabilities`;
- isolated R0 support for `PROMPT`, `CHANNEL_ROLE`, `PRIMITIVE`, `LAYOUT`, `REDUNDANCY`, `NUMERIC_STRUCTURE`, and `TEMPORAL_STRUCTURE`;
- explicit `UNSUPPORTED_MUTATION_FOR_PROFILE` for current-profile color/interference/depth/emergent requests;
- CI regression reproducing Tlaloc environment-variable invocation and proving parent/candidate exact-program equality;
- lifecycle regression proving the builder installs, reports capabilities, and uninstalls reversibly.

## Authority boundary

```text
Tlaloc
  owns failure analysis, experiment selection, model trials and evidence

Origami
  owns how supported mutation intents are compiled into pixels

Tonal
  may later pin a verified composition
```

Tlaloc does not gain canonical pixel authority by calling this builder.

## Hard invariants

```text
CANDIDATE_MUTATION = EXPERIMENTAL
PARENT_TEMPORAL_CARRIER_MUST_VERIFY
EXACT_TEMPORAL_PROGRAM_MUST_NOT_CHANGE
OUTPUT_PNG_BYTES = 8192
OUTPUT_DIMENSIONS = 640x640
NO_HIDDEN_SIDECAR
UNSUPPORTED_MUTATION -> EXPLICIT FAILURE
BUILD_SUCCESS != MODEL_IMPROVEMENT
TLALOC_RECOMMENDATION != ORIGAMI_PROMOTION
FALSE_EXACT = 0
```

## Evidence status

Mechanical candidate construction, exact-program preservation, fixed-envelope compliance and lifecycle integration are implemented and CI-gated. No candidate is promoted because it compiles. Native VLM improvement requires held-out model evidence from the experiment system.
