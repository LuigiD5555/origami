# Codec baselines for Origami experiments

Origami representation claims should be compared against conventional codecs on identical underlying semantics.

## Minimum baseline set

Depending on the experiment, compare against at least:

- canonical JSON;
- compact text;
- compact binary representation;
- Origami representation.

## Hold constant

- semantic program/graph;
- task/question;
- verifier;
- allowed external tools;
- evaluation metric;
- model/profile where applicable.

## Report separately

- semantic recovery;
- exact transport recovery;
- execution success;
- representation size;
- model/input token cost where measurable;
- latency;
- selective-access behavior;
- failures/UNKNOWN;
- false-known errors.

A result is not an Origami advantage if the Origami condition receives additional semantic hints unavailable to the baseline.
