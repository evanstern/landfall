---
name: route-gate
description: Route — the pure ex-ante admission gate; predicted drift vs. debt-inflated budget, with unmeasurable velocity suppressing.
kind: component
sources:
  - route.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Route — the ex-ante gate

`Route` in `route.go` decides whether a decision may go to the oracle. It is
pure arithmetic over the class's registered values, the host's
seconds-per-point for the target tier, the current world velocity, and the
class's accrued suppression debt: no model, no randomness, no wall-clock
reads. The same inputs always produce the same verdict.

## How it works

Signature: `Route(c Class, velocity, secondsPerPoint, debt float64) Verdict`.

The computation, in order:

1. `PredictedWallSec = Points × secondsPerPoint` — the latency prediction
   ([[points-latency-model]]).
2. `EffectiveBudget = Budget × (1 + debt × DebtFactor)` — the debt-inflated
   tolerance ([[suppression-debt]]).
3. If velocity is negative, NaN, or ±Inf, the verdict is a suppression with an
   Arithmetic string noting the unmeasurable velocity — computed *before* any
   drift math ([[velocity-semantics]]).
4. Otherwise `PredictedDrift = PredictedWallSec × velocity`, and
   `Allow = PredictedDrift <= EffectiveBudget`. The comparison is inclusive:
   drift exactly at budget is allowed (guarded by `TestRouteBoundary`).

Every path fills the `Arithmetic` field with the human-readable equation,
e.g. `3pt x 2.0s/pt x 2/s = 12.0 drift > budget 10.0 (base 10.0, debt 0.0)`.
The verdict also records every input it was fed ([[verdict]]), which is what
makes `TestVerdictReproducible` possible.

## Connections

Consumes a validated [[decision-class]] and a seconds-per-point number the
host typically reads from an [[estimator]]. Suppression debt is passed in by
the host and never accrued inside the gate ([[host-obligations]]). An allowed
verdict is frozen into a lease by `Checkout` ([[lease-lifecycle]]);
[[route-tiered]] runs this same function once per tier. The drift unit that
makes the comparison meaningful is [[drift-unit]].

## Operational notes

A suppressed verdict still carries the class's `Degrade` mode
([[degrade-modes]]) so the caller knows its declared floor. Purity is a hard
invariant ([[invariants-and-purity]]): suppression debt is passed in and
recorded, never stored; there is no runtime budget-tuning API by doctrine.
