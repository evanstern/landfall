---
name: calibration-doctrine
description: Calibration rules — successes-only feeding, pessimistic seeds, breach as a recalibration signal, recorded baselines moved only by humans.
kind: concept
sources:
  - estimate.go
  - DESIGN.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Calibration doctrine

DESIGN.md fixes four rules for how landfall's one learning component is fed,
seeded, and corrected. Together they implement the package-wide pessimism
principle: an uncalibrated system fails toward the degrade floor, never
toward stale action.

## How it works

- **Successes only.** The estimator learns from successful calls exclusively;
  failures belong to the host's circuit breaker. Rationale (recorded on
  `Estimator` in `estimate.go`): a failed call's duration says nothing about
  how long a landed answer takes. This is invariant 4 of
  [[invariants-and-purity]].
- **Pessimistic seeds.** Bootstrap defaults deliberately overestimate
  seconds-per-point, so an uncalibrated tier predicts long flights, fails the
  gate, and lands on the declared floor ([[degrade-modes]]) until real
  samples arrive. The same stance routes unmeasurable velocity to
  suppression ([[velocity-semantics]]).
- **Spikes excluded but counted; breach recommends, humans act.** One-shot
  lag is kept out of the EWMA; a windowed spike rate above `BreachRate`
  raises a one-time recalibration signal ([[estimator]]). The signal is
  advisory — it tells the host the baseline is lying.
- **The recorded baseline moves only under a human's hand.** The live
  estimate adapts freely at runtime; the persisted calibration file is
  edited by a person and re-seeded on restart. The daemon never writes it.

The monitor on all of this is the landing log: because every Landing pairs
predicted vs. actual drift ([[drift-unit]]), systematic error in velocity or
pace shows up as a pattern of Stale outcomes ([[landing-outcomes]]), and the
response is a human revising a constant — the same reasoning that keeps
budgets free of any runtime tuning API.

## Connections

Applies to the [[estimator]]; the numbers it produces feed
[[points-latency-model]] and thence [[route-gate]]. Persistence, breakers,
and the calibration file are host territory ([[host-obligations]]). The
doctrine-constants-vs-derived-thresholds tradeoff is argued in
[[lineage-and-prior-art]].

## Operational notes

Pair the gate with a busy-≠-down breaker; landfall only insists failures
don't feed the estimator. Cost-aware gating and learned budgets are
deliberately out of scope — the budget form is the practical linearization,
and closing the calibration loop automatically would put a model back in the
enforcement path.
