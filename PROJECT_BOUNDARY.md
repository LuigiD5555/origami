# Origami project boundary

```text
TLALOC                          ORIGAMI
work system                     representation
orchestration                   state language
behavior compilation   ----->   semantics contract
Tlaloque coordination  ----->   reference semantics engine/API
verification                     dynamics / projection
```

Origami may be used independently of Tlaloc. Tlaloc may also operate without Origami or with another representation provider.

The stable architectural rule is:

> Tlaloc coordinates work; Tlaloque perform bounded specialist work; Origami defines how a class of states is represented and transformed.

Origami supplies semantics and deterministic reference behavior. Tlaloc may turn those materials into expected-state fixtures and evaluation campaigns.
