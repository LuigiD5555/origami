# Perception Promotion R1

Status: `EXPERIMENTAL_REFERENCE_IMPLEMENTED`

Perception Promotion R1 reconstructs the missing deterministic evaluator from the earlier promotion plan and places it on top of Semantic Spine R1.

Its job is deliberately narrow:

> Compare what a target model says it perceived against the real Fixed Carrier. Do not let the model see evaluator ground truth, and do not turn a mock or one successful trial into a support claim.

## Why this exists

Fixed Carrier R2 already has deterministic ground truth for:

```text
T0 plaintext BOOT
T1 duplicated visual probe
Tool protocol
Address ABI
T3 deterministic machine record
StoreRoot / source identity / object counts
```

A VLM, however, is probabilistic. The evaluator therefore receives two different inputs:

```text
canonical carrier.png
model observation JSON
```

and compares them deterministically.

The model does not receive the decoded carrier record used by the evaluator.

## Strict model boundary

For strict Native trials, the target model receives only:

```text
universal Master Prompt
one Origami carrier image
user query
```

It must not receive:

```text
private evaluator manifest
expected probe bits
canonical T3 decode
private registry
ground-truth answer packet
```

For Hybrid trials the model may additionally access only the declared tool ABI. Tool results are part of Hybrid operation, not proof that T3 was read natively.

## Observation packet

A model runner serializes its interpretation as:

```json
{
  "schema": "origami.perception-observation.r1",
  "model": "example-vlm",
  "trial": 1,
  "transport": "original",
  "evidence_kind": "REAL_MODEL",
  "boot_text": ["..."],
  "probe_top": "10101010",
  "probe_bottom": "10101010",
  "tool_protocol": "tlaloc.origami-tools.r2",
  "address_abi": "ohf-address.r2",
  "t3": null
}
```

`MOCK` is also allowed so the harness can be tested, but mock evidence is mechanically useful only.

## Deterministic evaluator

`origami-perception-eval` independently decodes the canonical carrier and checks:

```text
T0_BOOT
T1_PROBE_TOP
T1_PROBE_BOTTOM
T1_PROBE_DUPLICATION
TOOL_PROTOCOL
ADDRESS_ABI
T3_MACHINE_RECORD (optional)
```

T0 comparison tolerates case/whitespace normalization but not semantic paraphrase. The duplicated visual probe must be reproduced from the image in both positions and both readings must agree.

Tool protocol and address ABI are exact protocol gates.

## Hybrid and Native are separate

A Hybrid trial can mechanically pass when the model demonstrates:

```text
T0
+
upper visual probe
+
lower visual probe
+
probe agreement
+
correct tool protocol
+
correct address ABI
```

T3 is intentionally not required.

This means:

```text
T3 Native failure
      !=
Hybrid failure
```

A Native T3 trial additionally requires the target model to reproduce the deterministic T3 machine record fields bound into the carrier:

```text
carrier id
StoreRoot
source SHA
page count
block count
document count
object count
```

That trial still does not by itself establish universal Native visual support.

## MOCK is never promotion evidence

A mock can produce a perfect packet and therefore validate the entire mechanical path:

```text
image
 -> observation packet
 -> evaluator
 -> report
```

The report can say:

```text
hybrid_mechanical_pass = true
native_t3_mechanical_pass = true
```

while still saying:

```text
hybrid_trial_promotion_eligible = false
native_t3_trial_promotion_eligible = false
```

because `evidence_kind = MOCK`.

That distinction is intentional and non-negotiable.

## REAL_MODEL still means one trial

A real-model observation that passes is only one eligible trial. The evaluator report carries:

```text
promotion_boundary = SINGLE_TRIAL_ONLY_NOT_STACK_SUPPORT
```

Cross-model thresholds, repeated trials, transport variants, real tool-loop success and final stack promotion are not decided by this package.

Ownership remains:

```text
Origami
  deterministic per-trial evaluator
  semantic/perceptual validity

Tlaloc
  model campaign runner
  transport variants
  repeated trials
  real tool loops
  routing benchmarks

Tonal
  immutable component composition
  aggregate promotion gates
  HYBRID_SUPPORTED / NATIVE_VISUAL_SUPPORTED
```

## Relation to Semantic Spine

Perception Promotion sits outside semantic authority:

```text
carrier image
   ↓
model perception
   ↓
Origami deterministic perception evaluator
   ↓
perception evidence
   ↓
resolution / tool use / source reopening
   ↓
Semantic Spine / Evidence Reduction
```

A correct probe reading proves that the image was perceived correctly for that challenge. It does not prove that every semantic claim made by the model is correct.

This preserves:

```text
PERCEPTION != RESOLUTION != EXECUTION != VERIFICATION
```

## CLI

```bash
./bin/origami-perception-eval \
  -carrier origami.png \
  -observation observation.json \
  -out perception-report.json
```

The report is suitable for Tlaloc campaign aggregation and Tonal promotion gates.

## What this implementation proves

The deterministic reference implementation proves that:

- carrier ground truth can be checked independently of model output;
- mock evidence cannot accidentally promote;
- Hybrid trial success is independent of Native T3;
- incorrect duplicated probes fail;
- wrong tool/address protocols fail;
- T3 claims are checked against the actual carrier record.

It does not prove that any external VLM can pass those gates. That evidence must come from real external campaigns.
