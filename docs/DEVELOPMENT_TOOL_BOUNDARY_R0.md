# Development Tool Boundary R0

Origami is the project being built and versioned. Development kits may help improve it, but they are not part of the minimum Origami runtime contract.

## Core rule

```text
Origami
  owns Origami semantics
  owns the canonical visual grammar
  owns ROSETTA semantics
  owns the Master Prompt release
  owns Writer/Reader contracts
  owns Origami version/profile promotion
```

External development systems may propose improvements:

```text
Tlaloc
Blueprint Framework
future alternative/complementary tools
```

but a proposal remains external until Origami adopts it.

## Tlaloc relationship

Tlaloc is especially useful for discovering behavior through bounded Tlaloque swarms and then distilling that demonstrated behavior into portable prompt candidates.

For Origami this can mean experimenting with:

```text
Master Prompt wording
BOOT/ROSETTA behavior
visual grammar changes
color/numeric/interference/depth/temporal channels
Writer/Reader procedures
perception/navigation behavior
```

Tlaloc supplies experiments, reference traces, candidates and evidence. Origami decides whether any result becomes canonical Origami.

## Prompt portability baseline

The universal Origami Master Prompt must remain useful to a compatible model even when the model has no:

```text
Tlaloc
Tonal
sandbox
Go runtime
Python runtime
tools
file access beyond explicitly supplied inputs
```

This is the portable baseline.

Tools and runtimes are optional enhancements. They can improve exact recovery, navigation, compilation and ingestion, but their absence must not be disguised.

## READ without tools

A model may use the Master Prompt plus the supplied carrier/image and question.

If a required mechanical/exact operation cannot be performed from what is genuinely available, the model must return `UNKNOWN` or `NOT_VERIFIED` rather than pretending it executed a tool.

## WRITE without compiler

A model can still derive Semantic IR, visual intents, ROSETTA and a construction specification.

If it cannot invoke a deterministic Origami compiler or otherwise produce and roundtrip-verify the requested carrier, it must report:

```text
WRITE_STATUS: CONSTRUCTION_SPEC_ONLY
```

It must not claim `COMPILED_VERIFIED` merely because it described how the image should look.

## Richer modes

The following may be added explicitly when available:

```text
prompt + declarative IR/context
prompt + Origami tools
prompt + compiler/runtime
specialized Origami-aware model
```

These are enhancements, not retroactive requirements of the universal prompt.

## Tonal relationship

Tonal may compose exact versions of development tools and target projects for reproducibility, for example:

```text
Tonal
├── Tlaloc
├── Blueprint Framework
└── Origami revision
```

Tonal does not decide what Origami means and is not required to read or write Origami.

## Promotion lifecycle

```text
development tool experiment
        ↓
candidate + evidence
        ↓
Origami validation
        ↓
Origami release/profile decision
        ↓
optional Tonal composition/pin
```

## Invariants

```text
ORIGAMI_OWNS_ORIGAMI_RELEASES
MASTER_PROMPT_IS_PORTABLE_BASELINE
DEVELOPMENT_TOOL != RUNTIME_REQUIREMENT
TLALOC_IS_OPTIONAL
TONAL_IS_OPTIONAL
TOOLS_ARE_OPTIONAL_ENHANCEMENTS
FALSE_EXACT=0
```
