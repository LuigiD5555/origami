# Origami Fixed Carrier R2

Fixed Carrier R2 is Origami's frozen visual **control plane**. It does not try to embed an arbitrarily growing corpus in pixels. The carrier teaches a generic visual model how to boot, proves that the image was actually seen, identifies the canonical memory universe, exposes root navigation spaces and binds that control state to Tlaloc's verified data plane.

## Frozen physical profile

```text
profile     origami.fixed-carrier.r2.profile-1
canvas      640 x 640
PNG bytes   exactly 8192
hard max    512000 bytes
```

Corpus growth is forbidden from changing these dimensions or bytes. If a future profile needs a different frozen envelope it requires a new version/profile and its own Perception Lab evidence.

## Bootstrap levels

```text
T0  plaintext bootloader
T1  Rosetta + duplicated visual probe
T2  root index / graph navigation
T3  Origami machine record
VERIFY  root binding / exactness state
```

T0 deliberately uses ordinary text so a generic VLM can enter Origami from a representation it already understands. OCR is optional: OCR failure is not a carrier failure. T1 then proves the model actually saw this carrier instance by requiring both visual challenge rows. T2 exposes PAGE / CONCEPT / SOURCE / GRAPH / VERIFY root spaces. T3 contains the deterministic 512-byte machine record and visual state structure.

## Tool/data-plane boundary

Fixed Carrier R2 binds `tlaloc.origami-tools.r2` and `ohf-address.r2`. Exact external corpus retrieval requires Tlaloc or another implementation of that ABI. The carrier contains control/routing/binding state, not arbitrary corpus bytes.

```text
model + origami.png
  -> T0/T1/T2/T3 BOOT
  -> origami_boot
  -> origami_query / origami_expand / origami_verify
  -> verified ContextPacket
```

## Exactness

A visual interpretation never grants exactness. Exact results require CID/Merkle verification against the store root bound into the carrier. `FALSE_EXACT=0` and `UNKNOWN > invention` remain mandatory. Native VLM-only decoding of T2/T3 is not promoted until held-out cross-model evidence exists.
