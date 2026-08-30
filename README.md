# Origami 6.0.0-alpha.12

Origami is an experimental **visual/computational representation, state-machine language and model-agnostic virtual memory** for complex state, relations, dynamics, selectively accessible knowledge and perceptual emergence.

Origami is independently usable. Development tools can improve it, but they are not hidden runtime requirements.

## What alpha.12 fixes

A real multimodal trial exposed a concrete failure: the model could read the visible Origami bootstrap but could not answer the simple question **“What is the index?”**. Instead it treated the carrier as a binary archive, requested file/decoder access and produced unverified byte, compression and hash claims.

Alpha.12 turns that failure into a permanent regression and changes the Native path from **binary-first** to **semantic-first**:

```text
SEMANTIC QUESTION
 -> T0 BOOT
 -> T1 ROSETTA / visual probe
 -> T2 ACTUAL SEMANTIC SUPERINDEX
 -> smallest visible/selective semantic region
 -> answer

EXACT QUESTION
 -> declared exact/control plane or tool when available
 -> verification
```

The exact/control plane remains useful. It is no longer an implicit prerequisite for semantic navigation.

See `docs/NATIVE_SEMANTIC_NAV_R0.md`, `spec/NATIVE_SEMANTIC_NAV_R0.json` and `experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json`.

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

## Portable Master Prompt baseline

`generated/MASTER_PROMPT.md` is **Universal Read / Write Bootstrap R3 — semantic-first**.

Its minimum compatibility assumption is intentionally small:

```text
Master Prompt
+ explicit user input
+ Origami image/carrier when image input is available
```

It does **not** assume Tlaloc, Tonal, sandbox, Go/Python, filesystem access, a decoder or external tools.

For semantic questions R3 explicitly forbids reconstructing bytes, guessing compression/archive structure or requesting a decoder merely to reach the index/overview. For exact questions, unavailable mechanical verification returns `NOT_VERIFIED` rather than invented exactness.

The registered R3 prompt remains `REFERENCE_CANDIDATE`; architecture and deterministic tests are implemented, but held-out prompt-only VLM evidence is still required.

## READ / WRITE

```text
READ
carrier/image
 -> BOOT
 -> ROSETTA
 -> classify semantic vs exact question
 -> T2 / visible semantic structure
 -> optional selective expand
 -> semantic answer

WRITE
source semantics
 -> Semantic IR
 -> Visual Intent
 -> canonical visual profile
 -> ROSETTA
 -> PROGRAM / T2 / MEMORY / VERIFICATION
 -> compiler when available
 -> roundtrip verify
```

Without tools, READ continues from what the model can genuinely perceive. An unavailable exact/mechanical operation returns `UNKNOWN` or `NOT_VERIFIED`; it is never pretended to have run.

Without a compiler, WRITE may produce a construction specification but must report:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

A textual description or screenshot is not a compiled Origami.

## Native semantic query routing

The reference routing contract is now explicit:

```text
What is this?             -> T0 + T1
What is the index?        -> T2
What is it about?         -> T2 + visible semantic structure
Where is topic X?         -> T2 -> selective semantic expansion when available
Explain topic X           -> T2 -> smallest relevant semantic region
Quote/hash exact X        -> exact/control plane or declared tool
```

`T2` must contain **actual semantic entry points**, not only generic labels such as PAGE/GRAPH/SOURCE. The current Fixed Carrier profile renders a bounded top-level semantic superindex from canonical memory. If only that bounded level is visible, the answer must be marked partial/top-level rather than inventing lower levels.

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
       |             ROSETTA + PROGRAM + T2 +
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

and preserves distinct `PRESENT`, `ABSENT`, `UNKNOWN`, `INHIBITED`, `CANCELLED` states. It implements declarative contradictions, signed cancellation, higher-order participant rules, finite observation contracts, temporal trajectory observations, semantic Fold and selective Unfold. Observation remains separate from transition.

## Canonical Visual Grammar R0

The current reference profile is:

```text
origami.canonical-aesthetic.r0
```

Origami has **one canonical functional aesthetic per profile version**. Current canonical dimensions include geometry, fill/contrast, position, topology, enclosure, scale, repetition, density and limited declarative text.

Experimental profile candidates include color, numeric/prime/modular/factorization structure, interference/moiré/phase, stereo/parallax/depth, temporal/motion-bound structure and emergent multi-layer or multi-instant percepts.

These are potential semantic channels, not decoration. They become authoritative only if a later Origami profile adopts them after evidence.

## ROSETTA and reveal procedures

ROSETTA always declares profile ID/version, active dimensions, visual primitive -> semantic-role mappings and approved local parameters/extensions.

For non-static channels it additionally declares a reveal procedure. Experimental operations already represented by Origami contracts include:

```text
MOIRE
PHASE_SHIFT
STEREO_BIND
PARALLAX_RESOLVE
KINETIC_REVEAL
TEMPORAL_INTEGRATE
```

If the declared reveal condition cannot be established, the result is `UNKNOWN`, not `ABSENT`.

## Writer R0

Writer turns semantic intention into a deterministic construction plan before pixels:

```text
Semantic IR
 -> Visual Intent Plan
 -> canonical profile binding
 -> ROSETTA
 -> PROGRAM / T2 / MEMORY / VERIFICATION plan
 -> compiler
```

A useful T2 semantic superindex is now an explicit Writer objective because common semantic questions must be cheap to answer without exact-plane decoding.

Reference planner:

```bash
./bin/origami-write-plan -in request.json -out plan.json
```

The roundtrip rule remains:

```text
DECODE(COMPILE(SEMANTIC_IR)) ~= SEMANTIC_IR
```

with exact equality for declared identity/address/hash/proof fields.

## Virtual Memory R0

Origami memory can be much larger than active model context. The model-facing working budget remains approximately 4000 token-equivalent; this is an active interface budget, not total storage capacity.

```text
large / multi-carrier memory
 -> GraphSignature routing
 -> local metadata graph
 -> selective fidelity unfold
 -> ContextPacket
```

Fidelity order remains:

```text
label -> abstract -> summary -> detail -> evidence -> exact
```

Address is location; CID is content identity. Deep exact payload is reopened selectively rather than silently scanned globally.

## Evidence Reduction R0

External models/development systems may propose `SUPPORT`, `OPPOSE`, `UNKNOWN`. Origami independently resolves evidence addresses/CIDs/source hashes/fidelity and reduces deterministically to `VERIFIED`, `REJECTED`, `CONFLICT`, `UNKNOWN`.

Agent/model confidence is never evidence authority. `VERIFIED_EXACT` requires byte-equal accepted exact evidence.

## Fixed Carrier R2 profile-2

Physical envelope is unchanged:

```text
640 x 640
8192 PNG bytes exactly
hard current-profile ceiling: 512000 bytes
```

Bootstrap is now functionally separated:

```text
T0  plaintext semantic-first BOOT
T1  ROSETTA + duplicated visual probe
T2  actual bounded semantic superindex
T3  deterministic control/exact record
VERIFY
```

The carrier is not the corpus. Corpus growth occurs in canonical/addressable memory, not by enlarging the PNG.

The integration layer no longer blindly truncates routing JSON at 256 bytes. It produces bounded **valid JSON** and prioritizes actual top-level semantic labels for T2.

Profile-1 deterministic decoding remains accepted for backward compatibility. Profile-2 changes the Native semantic route without discarding the exact/control record.

The existing `tlaloc.origami-tools.r2` ABI remains an optional external tool-provider ABI. Its existence does not make Tlaloc mandatory for Origami.

## Perception Promotion R1

Origami can deterministically compare a model's claimed carrier perception against actual carrier ground truth for T0, duplicated probes, Tool Protocol, Address ABI and optional T3. MOCK can validate plumbing but cannot provide empirical model promotion evidence.

The new Native semantic contract adds a different requirement: perceiving BOOT/probes is insufficient if the model still cannot answer through T2.

## Failed-trial discipline

The failed index experiment is preserved rather than overwritten by the new design. It establishes a regression:

```text
BOOT perceived              = yes
actual T2 index recovered   = no
mechanical decoder required = yes
unverified exact claims     = yes
```

A future profile may be smaller, denser or visually more sophisticated, but it cannot be promoted as a better Native semantic profile if it recreates this failure.

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

## Hard invariants

```text
ORIGAMI OWNS ORIGAMI RELEASES
MASTER PROMPT IS PORTABLE BASELINE
SEMANTIC NAVIGATION != MECHANICAL DECODE
INDEX QUERY ROUTES TO T2
T2 CONTAINS ACTUAL SEMANTIC ENTRIES
EXACT PLANE OPTIONAL FOR SEMANTIC NAVIGATION
NO UNVERIFIED BYTE / HASH / COMPRESSION CLAIMS
FAILED REAL TRIAL -> REGRESSION
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

Alpha.12 does **not** claim that the failed experiment already passes after this change. It establishes the corrected architecture and deterministic gates. Still pending:

- repeat the index question on the new profile-2 carrier with held-out real VLMs;
- measure Native T2 index recovery across multiple models/trials;
- verify semantic overview/topic-location behavior without undeclared mechanical decoding;
- test transport degradation of T2 readability;
- continue Hybrid/exact-plane tests separately.

## Source of truth

```text
VERSION
state/ORIGAMI_STATE.json
PROJECT_BOUNDARY.md
generated/MASTER_PROMPT.md
docs/DEVELOPMENT_TOOL_BOUNDARY_R0.md
docs/NATIVE_SEMANTIC_NAV_R0.md
docs/CURRENT_STATE.md
docs/CANONICAL_VISUAL_GRAMMAR_R0.md
docs/WRITER_R0.md
docs/SEMANTIC_SPINE_R1.md
spec/DEVELOPMENT_TOOL_BOUNDARY_R0.json
spec/NATIVE_SEMANTIC_NAV_R0.json
spec/FIXED_CARRIER_R2.json
spec/CANONICAL_VISUAL_GRAMMAR_R0.json
spec/WRITER_R0.json
experiments/native-semantic-nav-r0/FAILED_TRIAL_001.json
```

## Version

`6.0.0-alpha.12`
