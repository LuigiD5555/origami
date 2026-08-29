# Dimensional Visual Register (DVR) and Perceptual State Space (PSS)

**Status:** experimental design implemented in R3.10-LAB  
**Change:** `CHG-DVR-0001`  
**Promotion:** not promoted

## Purpose

This layer formalizes a central OHF idea: visual families and subdimensions are not flat enums. They are observable variables that modify a visual carrier at different scales. A receiver should only need to identify small physical states; deterministic software can perform packing, execution and verification.

The model deliberately separates mathematical capacity from perceptually demonstrated capacity.

```text
DIMENSION TREE
    ↓
leaf observable variables
    ↓
Dimensional Visual Register (DVR)
    ↓
cartesian candidate state space
    ↓
Perception Lab filters
    ↓
Perceptual State Space (PSS)
    ↓
orthogonality + compound + restoration gates
    ↓
SAFE_PSS
```

## DVR

A DVR is a vector of physical state assignments:

```text
DVR = (d1, d2, ... dn)
```

Each `di` belongs to one leaf visual dimension. Parent dimensions may contain subdimensions but do not contribute independent capacity by themselves.

Example:

```text
micro.outer_shape      = triangle
micro.topology_holes   = holes_2
micro.direction_beacon = east
micro.frame_count      = frames_1

meso.ordering          = clockwise
meso.grouping          = cluster

macro.topology         = branching
macro.symmetry         = bilateral
macro.orientation      = north
macro.region_structure = single_core
```

This is a physical register, not a semantic statement. `branching` does not universally mean knowledge, hierarchy or any other logical concept.

## Scale and scope

The implementation prevents a misleading multiplication of all dimensions as though every property existed independently in every glyph.

Three scopes are modeled:

### MICRO — per visual cell / glyph

Current reference leaves:

- `F01 outer_shape`
- `F02 topology_holes`
- `F03 direction_beacon`
- `F04 frame_count`

Each currently has four physical states. Therefore the nominal MICRO alphabet is:

```text
4 × 4 × 4 × 4 = 256 states = 8 nominal bits / glyph
```

This is a physical/codebook reference, not a claim that Native VLM perception has already demonstrated all 256 compounds.

### MESO — per group / motif

Current LAB candidates:

- `F28 ordering`
- `F29 grouping`

Each currently has four candidate states in the registry:

```text
4 × 4 = 16 states = 4 nominal bits / group
```

These states are experimental and unpromoted.

### MACRO — per carrier

Current LAB candidates:

- topology via `F35/F40`
- symmetry via `F37`
- orientation via `F03/F39`
- region structure via `F40/F41`

Current candidate product:

```text
4 × 4 × 4 × 4 = 256 states = 8 nominal bits / carrier
```

Again, this is an upper bound for the current candidate registry, not safe capacity.

## Carrier capacity

For:

- `C` micro cells,
- `G` meso groups,
- one macro carrier state,

nominal capacity is computed as:

```text
NominalBits = C * B_micro + G * B_meso + B_macro
```

The CLI computes this without claiming perceptual safety:

```bash
go run ./cmd/ohf-lab dimensional capacity 64 8
```

With the current registry this yields:

```text
MICRO  64 × 8 = 512 bits
MESO    8 × 4 =  32 bits
MACRO   1 × 8 =   8 bits
-------------------------
NOMINAL            552 bits
```

`SAFE_PSS` remains zero until evidence promotes dimensions.

## Perceptual Orthogonality Law

A dimension may contribute independent capacity only when changing that dimension does not make another relevant dimension ambiguous or destroy its observable state.

> **Perceptual Orthogonality Law:** a visual dimension contributes independent OHF capacity only if its controlled mutation remains discriminable while preserving the relevant states of the other active dimensions.

This converts previously informal rendering bugs into explicit interaction failures.

Example:

```text
orientation = east
symmetry    = rotational_6
```

If the chosen symmetry makes east and west visually equivalent, orientation is not independent in that compound. The nominal product therefore overestimates effective capacity.

## Orthogonality evidence

Pair evidence has four important properties:

1. absent evidence means `UNKNOWN`, never PASS;
2. PASS is explicit and evidence-linked;
3. FAIL prevents the pair from entering a candidate compatible set;
4. pairwise PASS is necessary but not sufficient for a high-dimensional compound.

Example matrix entry:

```json
{
  "a": "macro.topology",
  "b": "macro.symmetry",
  "status": "PASS",
  "evidence_ids": ["RUN-..."]
}
```

The implementation can find a maximum pairwise-compatible candidate set. This is a search aid only. It does not promote the compound.

## PSS and SAFE_PSS

### Known inherited interaction failure

The pre-DVR Macro-Gestalt experiment found a deterministic collision between orientation and rotational symmetry (`M0333` vs `M0313`): the symmetry erased the east/west distinction. The initial orthogonality matrix therefore records `macro.orientation × macro.symmetry = FAIL` until a corrected renderer passes controlled-mutation and perception gates.

Three quantities must remain distinct:

```text
THEORETICAL / NOMINAL PSS
all combinations allowed by the registry

OBSERVED PSS
states actually exercised by Perception Lab

SAFE_PSS
states allowed after promotion gates
```

A dimension can enter `SAFE_PSS` only after its state marks it `PROMOTED_SAFE` and required orthogonality evidence is present. Even then, the capacity report warns that higher-order compound and restoration gates are still required.

The current project deliberately reports:

```text
SAFE_PSS = UNPROVEN
```

for the default registry.

## Multi-dimensional bounded machine interpretation

The visual carrier can be treated as a finite multidimensional tape/register space while Go acts as the bounded deterministic machine.

```text
VISUAL TAPE / REGISTERS
    Micro: cell state
    Meso: group/motif state
    Macro: carrier gestalt state
          ↓
ATTENTION HEAD
    spatial / graph / hierarchy / dependency windows
          ↓
VLM PERCEPTION
    observe physical states only
          ↓
PERCEPTION WALL
          ↓
SAFE_MICRO_ISA / Glyph Calculus
    deterministic bounded execution
          ↓
VERIFICATION
```

OHF does not need or want unrestricted Turing completeness in the carrier. The execution layer remains bounded by explicit budgets and verification.

## Perception Lab progression for a new dimension

A new dimension or subdimension should pass the following sequence:

```text
D0 state visibility
D1 one-dimension controlled mutation
D2 pairwise orthogonality
D3 correlated-failure analysis
D4 compound collision tests
D5 invariance / degradation
D6 OOD
D7 composition
D8 transfer
D9 restoration gate where exactness applies
```

Only then can it become a promotion candidate.

## Optimization objective

The project should optimize:

```text
safe perceptual bits / visual unit
```

not raw cartesian-state count.

Adding a state that increases nominal capacity but causes ambiguity, retries, attention expansion or false-known outputs is a regression, not a win.

## Files

```text
internal/lab/dimensional/
profiles/dimensional_registry_v1.json
profiles/orthogonality_matrix_v1.json
experiments/dvr-example-001.json
state/components/dimensional_visual_register.json
state/components/perceptual_state_space.json
state/components/perceptual_orthogonality.json
```

## CLI

```bash
# print the current built-in registry
go run ./cmd/ohf-lab dimensional registry

# inspect nominal vs safe capacity for 64 cells and 8 meso groups
go run ./cmd/ohf-lab dimensional capacity 64 8

# validate a DVR assignment
go run ./cmd/ohf-lab dimensional validate-register experiments/dvr-example-001.json

# analyze a supplied orthogonality matrix
go run ./cmd/ohf-lab dimensional orthogonality profiles/orthogonality_matrix_v1.json
```

No command above performs external requests.
