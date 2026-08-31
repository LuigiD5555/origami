# Experimental Candidate Build R0

Origami owns how experimental representation mutations become pixels. Tlaloc may propose and prioritize mutations, but it does not directly edit an Origami carrier.

The R0 flow is:

```text
Tlaloc failure frontier / adaptive plan
        -> declarative mutation list
        -> origami-candidate-build
        -> parent temporal carrier verification
        -> visual-only experimental mutation
        -> exact TemporalProgram re-verification
        -> 640x640 / exactly 8192-byte candidate PNG
        -> Tlaloc held-out model trials
```

## Why this boundary exists

A Tlaloc candidate such as:

```json
{
  "kind": "LAYOUT",
  "target": "T1_TO_T2_ENTRY_ROUTE",
  "value": "EXPLICIT_DIRECTIONAL_ANCHOR",
  "experimental": true
}
```

expresses an experimental intent. It does not grant Tlaloc authority to define canonical Origami pixels. `origami-candidate-build` is the deterministic Origami-owned compiler for the currently supported subset.

## Exact-program preservation

R0 accepts only `origami.temporal-carrier.r0.profile-1` parents. Before editing the visible semantic plane it mechanically decodes and verifies the embedded TemporalProgram. After rendering the mutation, it decodes the new PNG again and requires:

```text
candidate.program_sha256 == parent.program_sha256
```

The builder never mutates the exact-program pixel region. If the candidate no longer decodes, if the exact SHA changes, or if the PNG no longer fits the frozen 8192-byte envelope, the build fails.

## Supported mutation families

R0 supports bounded visual changes for:

- `PROMPT`
- `CHANNEL_ROLE`
- `PRIMITIVE`
- `LAYOUT`
- `REDUNDANCY`
- `NUMERIC_STRUCTURE`
- `TEMPORAL_STRUCTURE`

The current grayscale temporal profile does not claim native support for `COLOR_USAGE`, `INTERFERENCE_STRUCTURE`, `DEPTH_STRUCTURE`, or `EMERGENT_STRUCTURE`; those requests fail as unsupported rather than being approximated silently.

## Tlaloc closed-loop hook

Tlaloc alpha.19 already invokes explicit external candidate builders with these environment values:

```text
TLALOC_CANDIDATE_ID
TLALOC_OUTPUT_PNG
TLALOC_MUTATIONS_JSON
TLALOC_PARENT_SPECIMEN_ID
TLALOC_PARENT_PNG
```

Therefore a candidate may use:

```json
"build_command": ["origami-candidate-build"]
```

No wrapper script is required. The same command can also be called directly with `build` flags.

## Capabilities

```bash
origami-candidate-build capabilities
```

returns the supported parent profile and mutation families so an adaptive experiment runner can filter suggestions before spending model trials.

## Evidence boundary

A successful build proves only:

- the parent is a valid temporal carrier;
- the requested mutation belongs to the implemented experimental subset;
- the candidate still contains the same exact TemporalProgram;
- the output remains within the frozen carrier envelope.

It does **not** prove that a VLM understands the candidate better. Only held-out model trials may establish that.
