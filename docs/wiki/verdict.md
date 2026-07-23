---
name: verdict
description: Verdict — the gate's decision plus the full arithmetic that produced it; reproducible from its own recorded fields.
kind: component
sources:
  - route.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Verdict — the decision that carries its arithmetic

`Verdict` in `route.go` is the gate's output: not just allow/refuse, but every
number that produced the decision plus the equation as a human-readable
string. Every verdict — especially a suppression — is meant to be recorded by
the host, making each gate decision auditable and replayable.

## How it works

Fields: `Allow`, `Class`, `Tier` (chosen tier name; empty for plain `Route`
or when no tier fits), `Points`, `SecondsPerPoint`, `Velocity` (as fed by the
host), `PredictedWallSec` (Points × SecondsPerPoint), `PredictedDrift`
(PredictedWallSec × Velocity), `Budget` (the class's declared budget),
`EffectiveBudget` (debt-inflated), `Debt`, `Degrade` (what the caller does
when `!Allow`), and `Arithmetic` (the equation, human-readable).

The defining property: a logged verdict is reproducible from its own recorded
fields. `TestVerdictReproducible` reconstructs a `Class` from the verdict's
recorded values, re-runs `Route` with the recorded velocity,
seconds-per-point, and debt, and requires an identical struct. The same holds
for tiered verdicts (`TestRouteTieredVerdictReproducible`) — the recorded
seconds-per-point and debt are the chosen tier's, so `Route` over them
regenerates everything but the tier name.

## Connections

Produced by [[route-gate]] and [[route-tiered]]. An allowed verdict's
`EffectiveBudget` and `PredictedDrift` are frozen into the lease by
`Checkout` ([[lease-lifecycle]]), which is how the landing log inherits the
same reproducibility ([[landing-outcomes]]). The recorded `Debt` documents
the [[suppression-debt]] the gate was told about; the `Degrade` field echoes
the class's declared floor ([[degrade-modes]]).

## Operational notes

Recording is the host's job ([[host-obligations]]): the doctrine is that
every route decision, especially suppressions, is persisted — the verdict log
is half of the gate's calibration story (the landing log is the other half,
see [[drift-unit]]). The `Arithmetic` string is for humans; the numeric
fields are the replayable record. Extending `Verdict` requires extending
`TestVerdictReproducible` (see [[test-suite]]).
