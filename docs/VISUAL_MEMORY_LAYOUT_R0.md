# Visual Memory Layout R0

Virtual Memory R0 separates logical memory from its visual projection. The visual projection helps a receiver navigate; it is not the sole storage/exactness authority.

## One image, two complementary layers

`RenderMemory` produces one PNG:

```text
┌─────────────────────────────────────────────┐
│ redundant visual navigation band            │
│ MACRO            MESO             MICRO     │
│ roots/domains -> clusters/subgraphs -> nodes│
├─────────────────────────────────────────────┤
│ exact self-decodable Hybrid/Glyph transport │
│ BOOT / ROSETTA / PROGRAM / INDEX / MEMORY   │
│ / VERIFICATION                              │
└─────────────────────────────────────────────┘
```

Computational mode strips the declared navigation-band height and deterministically decodes the exact payload. Native/Hybrid perception may use the band as a map. A corrupted or misread navigation mark cannot by itself authorize `EXACT`.

## Macro / meso / micro

### MACRO

Answers: *Which broad region/root is relevant?*

Root nodes are projected into the macro zone. In a future multi-carrier visual composition, graph signatures provide the corresponding carrier-level macro map.

### MESO

Answers: *Which cluster/subgraph should be entered?*

Clusters are separate anchors. Grouping, containment and hierarchy are preferred over stuffing more scalar symbol states into one glyph.

### MICRO

Answers: *Which address should be unfolded?*

Addressable nodes are projected as local graph objects. Relation edges are redundant visual copies of canonical `Links`/`Edges`; the logical graph remains authoritative.

## Family allocation

The complete mapping lives in `spec/VISUAL_MEMORY_PROFILE_R0.json`.

The current reference projection actually draws a conservative structural subset:

```text
F01 node kind
F21 adjacency
F22 containment
F24 connection
F25 edge direction
F29 grouping
F32 hierarchy depth
F34 graph degree
F39 spatial zone
F40 macro region
F41 meso/micro transition
F42 verification redundancy
```

The remaining registered families still have proposed memory roles, but they are not automatically consumed as independent capacity. Color, texture, fill, stroke, rings, offsets, pair geometry, motif, symmetry and repetition remain evidence-gated candidates.

This is deliberate. Origami must not become a code bar whose nominal capacity comes from multiplying theoretical channels a VLM cannot reliably separate.

## Visual semantics by class

The reference renderer currently uses shape class for node kind, spatial zones for resolution, cluster boxes/grouping for meso structure, edges/arrows for relations/direction, node size as a weak importance redundancy, local degree marks as graph-degree redundancy and an extra ring for verified evidence/exact-level availability.

None of these marks replace the logical node fields. They provide redundant navigation cues that Perception Lab can test and promote independently.

## Perception Lab requirements

Visual-memory promotion is not established by exact PNG round-trip. It requires held-out model evidence for:

1. carrier/root selection;
2. cluster selection;
3. macro -> meso -> micro descent;
4. edge following;
5. random address reopening;
6. irrelevant cluster rejection;
7. resize/degradation robustness;
8. symbol permutation through carrier-local Rosetta;
9. composition/higher-order family interactions;
10. cross-model behavior.

See `spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json`.

## Why exact transport remains below the map

During this stage, the map and exact transport intentionally coexist. That gives us clean failure attribution:

- if `ContextPacket` logic fails, the memory architecture failed;
- if deterministic PNG recovery fails, the transport failed;
- if a VLM cannot navigate the band, perception failed;
- if verification cannot reopen evidence, exactness failed.

Only after those dimensions are independently understood should more logical burden be moved into emergent/visual-only channels.

## Future family experiments

Additional perceptual families such as hue/saturation/value, gradients, texture, frequency, negative space, interference, depth and time can increase useful visual bandwidth only after they pass the same evidence discipline. Perceptual Channels R0 permits latent/emergent/temporal information, but contract availability is not runtime promotion.
