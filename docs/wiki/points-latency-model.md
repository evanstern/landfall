---
name: points-latency-model
description: Latency factored as shape × host — class Points (Fibonacci, host-independent) times per-tier calibrated seconds-per-point.
kind: concept
sources:
  - class.go
  - DESIGN.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# The points latency model — shape × host

landfall predicts a call's wall time as `Points × secondsPerPoint`: the
product of a host-independent measure of the prompt shape and a per-tier,
live-calibrated pace. DESIGN.md records this factoring as the first of three
ideas the initial generalization missed and the source implementation taught.

## How it works

- `Points` (on [[decision-class]], `class.go`) is a property of the prompt
  shape — how big the ask is — Fibonacci by convention and independent of
  which oracle serves it. `Validate` requires it positive.
- `secondsPerPoint` is a property of the tier: calibrated, live-estimated by
  the [[estimator]], and supplied to the gate at route time (as an argument
  to `Route`, or on each `Tier` for `RouteTiered`).

Consequences DESIGN.md calls out:

- One estimator per tier instead of one per (oracle, class) pair.
- A brand-new class gets latency predictions with zero observations: assign
  points, inherit the tier's calibration.
- For a multi-call tool-use loop, seconds-per-point covers the *whole loop* —
  the same unit the estimator observes from completed calls — so the
  arithmetic stays truthful.

## Connections

The product becomes `PredictedWallSec` in the [[verdict]], which
[[route-gate]] multiplies by velocity to get predicted drift in the shared
[[drift-unit]]. Per-tier paces are the distinguishing input of
[[route-tiered]]. Feeding and seeding the pace estimate is governed by
[[calibration-doctrine]]; supplying the number at route time is a host
obligation ([[host-obligations]]).

## Operational notes

Points are doctrine like budgets: reviewed constants on the class, no runtime
tuning. If a class's calls systematically run long relative to prediction,
the landing log surfaces it ([[landing-outcomes]]) — the fix is a human
revising points or recalibrating the tier, not an automatic adjustment.
