# Retest contract — t2-temporal-grammar-visible-r1

Use the candidate PNG and the same blind prompt/model conditions as the baseline manual DeepSeek trial.

Primary expected improvement is isolated to temporal understanding:

```text
Q3 causal trigger for B becoming ACTIVE
Q4 consequences after B becomes ACTIVE
Q7 final synchronous states A=DONE, B=DONE, C=ACTIVE
```

Non-regression expectations:

```text
BOOT/ROSETTA discovery must not regress
A/B/C perception must not regress
initial A=ACTIVE recovery must not regress
FALSE_EXACT must remain 0
exact TemporalProgram SHA-256 must remain unchanged
PNG envelope must remain exactly 8192 bytes
```

This is a hypothesis contract, not a claimed result. A model response recorded after the candidate is built must supply the evidence.
