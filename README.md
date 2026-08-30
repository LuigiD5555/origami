# Origami 6.0.0-alpha.11

Origami is an experimental **visual/computational representation, state-machine language and model-agnostic virtual memory** for complex state, relations, dynamics, selectively accessible knowledge and perceptual emergence.

Origami is independently usable. Development tools can improve it, but they are not hidden runtime requirements.

## Ecosystem boundary

```text
Tlaloc             Blueprint Framework           future tools
 behavioral dev        structural dev                 ...
      \                     |                       /
       \------ candidates / experiments / evidence -/
                              |
                              v
                           ORIGAMI
          semantic + visual language / machine / memory
          owns Master Prompt, ROSETTA, profiles and releases
                              |
                              v
                    standalone portable use

Optional: Tonal may pin/compose the exact development toolchain.
```

Tlaloc is a particularly useful behavioral development kit: it can discover a working procedure through bounded Tlaloque swarms and distill that behavior into prompt candidates. Origami is one possible Tlaloc target, not a Tlaloc subsystem.

See `docs/DEVELOPMENT_TOOL_BOUNDARY_R0.md` and `spec/DEVELOPMENT_TOOL_BOUNDARY_R0.json`.

## Portable Master Prompt baseline

`generated/MASTER_PROMPT.md` is now **Universal Read / Write Bootstrap R2**.

Its minimum compatibility assumption is intentionally small:

```text
Master Prompt
+ explicit user input
+ Origami image/carrier when image input is available
```

It does **not** assume:

```text
Tlaloc
Tonal
sandbox
Go / Python
external tools
hidden file access
```

When tools, declarative IR, ingestion adapters or an Origami compiler exist they are explicit enhancements, not retroactive requirements.

This gives Origami a prompt-first retrocompatibility path for models that have no runtime except their normal text/image interface.

## READ / WRITE

```text
READ
carrier/image
 -> BOOT
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> semantic answer

WRITE
source semantics
 -> Semantic IR
 -> Visual Intent
 -> canonical visual profile
 -> ROSETTA
 -> PROGRAM / INDEX / MEMORY / VERIFICATION
 -> compiler when available
 -> roundtrip verify
```

Without tools, READ continues from what the model can genuinely perceive/derive. An unavailable exact/mechanical operation returns `UNKNOWN` or `NOT_VERIFIED`; it is never pretended to have run.

Without a compiler, WRITE may produce a construction specification but must report:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

A textual description or screenshot is not a compiled Origami.

## Canonical architecture

```text
PDF / IMAGE / TEXT / CONVERSATION / STATE
                    |
                    v
       declared source/ingestion adapter
          (optional external tooling)
                    |
                    v
              SEMANTIC IR
                    |
       +------------+-------------+
       |                          |
       v                          v
SEMANTIC SPINE R1             WRITER R0
S_t + C_t + R                visual intents
Fold / Unfold                     |
evidence                          v
       |              CANONICAL VISUAL GRAMMAR R0
       |                          |
       |                          v
       |              ROSETTA + PROGRAM + INDEX +
       |                 MEMORY + VERIFICATION
       |                          |
       |                          v
       |                 deterministic compiler
       |                    when available
       |                          |
       +------------+-------------+
                    |
                    v
               VIRTUAL MEMORY
                    |
                    v
          bounded ContextPacket
                    |
                    v
            compatible model
```

External development systems can construct experiments around this path without becoming its semantic authority.

## Semantic Spine R1

The deterministic semantic core uses:

```text
S_(t+1) = F(S_t, C_t, R)
```

and preserves distinct:

```text
PRESENT
ABSENT
UNKNOWN
INHIBITED
CANCELLED
```

It implements declarative contradictions, signed cancellation, higher-order participant rules, finite observation contracts, temporal trajectory observations, semantic Fold and selective Unfold.

Observation remains separate from transition.

## Canonical Visual Grammar R0

The current reference profile is:

```text
origami.canonical-aesthetic.r0
```

Origami has **one canonical functional aesthetic per profile version**. Many documents may produce different carriers, but they share the same active profile rules.

Current canonical dimensions include:

```text
geometry
fill / contrast
position
topology
enclosure
scale
repetition
density
limited declarative text
```

Experimental profile candidates include:

```text
color
numeric / prime / modular / factorization structure
interference / moiré / phase
stereo / parallax / depth
temporal / motion-bound structure
emergent multi-layer or multi-instant percepts
```

These are potential semantic channels, not decoration. They become authoritative only if a later Origami profile adopts them after evidence.

## ROSETTA and reveal procedures

ROSETTA always declares:

```text
profile ID/version
active dimensions
visual primitive -> semantic role mappings
approved local parameters/extensions
```

For non-static channels it additionally declares a **reveal procedure**: how a receiver should establish the percept.

Examples of experimental operations already represented by Origami contracts include:

```text
MOIRE
PHASE_SHIFT
STEREO_BIND
PARALLAX_RESOLVE
KINETIC_REVEAL
TEMPORAL_INTEGRATE
```

If the declared reveal condition cannot be established:

```text
UNKNOWN
```

not `ABSENT`.

## Writer R0

Writer turns semantic intention into a deterministic construction plan before pixels:

```text
Semantic IR
 -> Visual Intent Plan
 -> canonical profile binding
 -> ROSETTA
 -> carrier sections / verification plan
 -> compiler
```

Reference planner:

```bash
./bin/origami-write-plan -in request.json -out plan.json
```

Visual intents include identity, hierarchy, relation, dependency, scope, importance, state, uncertainty, evidence, verification, addressability, temporal order, redundancy and routing.

The roundtrip rule remains:

```text
DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR
```

with exact equality for declared identity/address/hash/proof fields.

## Virtual Memory R0

Origami memory can be much larger than active model context.

The model-facing working budget remains approximately:

```text
4000 token-equivalent
```

This is an active interface budget, not total storage capacity.

```text
large / multi-carrier memory
 -> GraphSignature routing
 -> local metadata graph
 -> selective fidelity unfold
 -> ContextPacket
```

Fidelity order:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Address is location; CID is content identity. Deep exact payload is reopened selectively rather than silently scanned globally.

## Evidence Reduction R0

External models/development systems may propose:

```text
SUPPORT
OPPOSE
UNKNOWN
```

Origami independently resolves evidence addresses/CIDs/source hashes/fidelity and reduces deterministically to:

```text
VERIFIED
REJECTED
CONFLICT
UNKNOWN
```

Agent/model confidence is never evidence authority.

`VERIFIED_EXACT` requires byte-equal accepted exact evidence. A paraphrase that cites exact evidence is not byte-exact.

## Fixed Carrier R2

Current fixed visual control-plane reference:

```text
640 x 640
8192 PNG bytes exactly
hard current-profile ceiling: 512000 bytes
```

Bootstrap:

```text
T0 plaintext BOOT
T1 ROSETTA + duplicated visual probe
T2 root index / graph navigation
T3 deterministic machine record
VERIFY
```

The carrier is not the corpus. Corpus growth occurs in canonical/addressable memory, not by enlarging the PNG.

The existing `tlaloc.origami-tools.r2` ABI is an **optional external tool-provider ABI**. Its existence does not make Tlaloc mandatory for Origami.

## Perception Promotion R1

Origami can deterministically compare a model's claimed carrier perception against actual carrier ground truth for:

```text
T0
duplicated visual probes
Tool Protocol
Address ABI
optional T3
```

MOCK can validate plumbing but cannot provide empirical model promotion evidence. Hybrid and Native-T3 claims remain separate.

## How Origami evolves

A development tool may propose changes, but Origami owns adoption:

```text
Origami profile N
 -> external experiment
 -> candidate prompt/profile + evidence
 -> Origami semantic/visual validation
 -> Origami decides profile/version N+1
```

Tlaloc can use its swarm-to-prompt playground to search better Master Prompts and representation behaviors. Blueprint Framework or another tool may contribute different development guarantees. Tonal may optionally record which exact toolchain/revisions were used.

None of them replaces Origami's release authority.

## Hard invariants

```text
ORIGAMI OWNS ORIGAMI RELEASES
MASTER PROMPT IS PORTABLE BASELINE
TLALOC IS OPTIONAL
TONAL IS OPTIONAL
DEVELOPMENT TOOL != RUNTIME REQUIREMENT
TOOL-ASSISTED SUCCESS != PROMPT-ONLY SUCCESS
ROSETTA ALWAYS PRESENT
ONE CANONICAL AESTHETIC PER PROFILE VERSION
FAILED REVEAL != SEMANTIC ABSENCE
SCREENSHOT != ORIGAMI SEMANTIC REPRESENTATION
LITERAL TRANSPORT != SEMANTIC FOLD
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
CONSTRUCTION SPEC != COMPILED VERIFIED CARRIER
ACTIVE MODEL INTERFACE != TOTAL MEMORY
NO IMPLICIT GLOBAL EXACT SCAN
FALSE_EXACT = 0
UNKNOWN > INVENTED EXACTNESS
```

## Evidence still pending

Alpha.11 does not claim that:

- every model can natively read/write the current carrier;
- every experimental perceptual channel is reliable;
- prompt-only mode equals Hybrid/tool mode for every task;
- generic Writer plans already have a complete profile-generic renderer for all visual channels;
- any prime/moiré/stereo/temporal candidate is better than R0.

These remain experiments, not assumptions.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
generated/MASTER_PROMPT.md
docs/DEVELOPMENT_TOOL_BOUNDARY_R0.md
docs/CURRENT_STATE.md
docs/CANONICAL_VISUAL_GRAMMAR_R0.md
docs/WRITER_R0.md
docs/SEMANTIC_SPINE_R1.md
docs/PERCEPTION_PROMOTION_R1.md
spec/DEVELOPMENT_TOOL_BOUNDARY_R0.json
spec/CANONICAL_VISUAL_GRAMMAR_R0.json
spec/WRITER_R0.json
```

## Version

`6.0.0-alpha.11`
