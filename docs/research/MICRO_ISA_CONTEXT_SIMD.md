# Run this exactly — SAFE_MICRO_ISA / CONTEXT_SIMD

This is the minimal operator procedure for testing on your machine.

## 0. Enter repository

```bash
cd origami
```

## 1. Prove local code is healthy

```bash
go test ./...
go vet ./...
```

Both must pass. If either fails, do not run external probes.

## 2. Rebuild the small Glyph Calculus carrier

```bash
go run ./cmd/ohf-glyphcalc demo \
  --pattern mixed \
  --length 1024 \
  --seed 817392 \
  --out runs/glyphcalc-native-small

go run ./cmd/ohf-glyphcalc verify \
  --run runs/glyphcalc-native-small
```

Expected: `STATUS=PASS` and exact reconstruction.

## 3. Generate micro-op probes locally

```bash
go run ./cmd/ohf-glyphcalc microbench-prepare \
  --out runs/microbench \
  --widths 1,4,8,16 \
  --cases 2 \
  --seed 817392
```

This performs **zero external requests**.

Each case contains:

```text
probe.json                 <- ground truth, DO NOT SEND
response.template.json     <- local helper
blind/
  MASTER_PROMPT.md         <- SEND
  carrier.png              <- SEND
  query.txt                <- SEND
```

Never send `probe.json` to the model.

## 4. Ask the state machine for the next allowed request

Run:

```bash
go run ./cmd/ohf-glyphcalc microbench-next \
  --root runs/microbench \
  --min-exact 0.98
```

It will return exactly one `case_dir` or `STOP`. Do not skip ahead.

The first three requests will be the Stage-A width-1 comparison.

## 5. External Stage A: exactly 3 requests

Use only the `case_dir` returned by `microbench-next`. The Stage-A set is:

```text
runs/microbench/observe4/w01/c00/blind/
runs/microbench/tuple4/w01/c00/blind/
runs/microbench/decode8/w01/c00/blind/
```

For each, supply only the three files in `blind/` to the model.

Save its raw JSON as `response.json` in the parent case directory.

If the service reports latency/token usage, put those numbers in the response `metrics` object. If unavailable, leave zeros.

Score immediately:

```bash
go run ./cmd/ohf-glyphcalc microbench-score \
  --case runs/microbench/tuple4/w01/c00 \
  --response runs/microbench/tuple4/w01/c00/response.json
```

If the command exits non-zero because `false_known > 0`, that candidate is unsafe. Do not widen it.

After scoring **every response**, immediately run `microbench-next` again.

## 6. Select one operation

Compare the three Stage-A `score.json` files.

Priority:

```text
1. false_known must be 0
2. highest exact_known_rate
3. lower retries
4. lower output tokens per glyph
5. lower latency per glyph
6. lower measured cost per glyph
```

Do not choose DECODE8 merely because it produces a byte directly. If TUPLE4 is safer, let Go pack the tuple.

## 7. External Stage B: widen only the winner

Example if `TUPLE4` wins:

```text
w04/c00
w08/c00
w16/c00
```

One request at a time. Score immediately after each.

Stop at the first:

```text
false_known > 0
or
exact_known_rate < 0.98
```

Do not send the wider probes after the stop condition.

## 8. Confirm the boundary only

Use `c01` only for:

- the largest passing width;
- the first failing width, if any.

Do not repeat every width.

## 9. Build profile

```bash
go run ./cmd/ohf-glyphcalc microbench-profile \
  --root runs/microbench \
  --min-exact 0.98
```

This writes:

```text
runs/microbench/profile.json
```

Treat the result as model/profile-specific evidence, not a universal truth about all LLMs.

## 10. Create a Context SIMD bundle for the real 82-glyph probe

Replace `8` with the safe TUPLE4 width discovered above:

```bash
go run ./cmd/ohf-glyphcalc simd-prepare \
  --run runs/glyphcalc-native-small \
  --width 8 \
  --out runs/glyphcalc-native-small/context-simd-w08
```

Only send:

```text
context-simd-w08/blind/MASTER_PROMPT.md
context-simd-w08/blind/carrier.png
context-simd-w08/blind/query.txt
```

Never send:

```text
plan.json
expected.response.json
source.bin
transport.bin
glyphs.json
```

Save the model response outside `blind/`, e.g.:

```text
context-simd-w08/response.json
```

Check it:

```bash
go run ./cmd/ohf-glyphcalc simd-check \
  --bundle runs/glyphcalc-native-small/context-simd-w08 \
  --response runs/glyphcalc-native-small/context-simd-w08/response.json
```

Required safety result:

```text
false_exact_zero = true
```

Ideal result:

```text
exact_transport_recovered = true
```

Only after that should the recovered transport be executed by the deterministic Glyph Calculus runtime.

## 11. Do not promote anything yet

The experiment remains LAB-only until repeated cross-render / cross-model Perception Lab gates pass.

## Quick installation

```bash
chmod +x install.sh
./install.sh
```

The installer uses no `sudo`, downloads no packages and makes no external requests. It installs `ohf-lab` and `ohf-glyphcalc` under `~/.local/bin` by default.
