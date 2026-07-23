---
name: estimator
description: Estimator — per-tier EWMA seconds-per-point with spike rejection, a windowed breach signal that fires once, and mutex-guarded state.
kind: component
sources:
  - estimate.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Estimator — the live seconds-per-point

`Estimator` in `estimate.go` is the live seconds-per-point estimate for one
tier — the only stateful piece of landfall. It is an EWMA with spike
rejection, fed by successes only, seeded pessimistically, with a breach
signal that recommends recalibration when the recent window says the baseline
is lying.

## How it works

`NewEstimator(seed float64, cfg EstimatorConfig)` — the seed is the
uncalibrated starting estimate; doctrine says seed pessimistically
([[calibration-doctrine]]). `EstimatorConfig` fields: `Alpha` (EWMA weight
for a fresh sample), `SpikeFactor` (a sample beyond
`SpikeFactor × estimate` is a spike), `WindowSize` (ring of recent samples
the breach signal watches), `BreachRate` (spike rate over a full window that
triggers the signal). `DefaultEstimatorConfig()` carries promptworld's tuned
values: Alpha 0.2, SpikeFactor 3.0, WindowSize 20, BreachRate 0.3. The zero
config is not usable.

`Sample(secPerPoint float64) (recalibrate bool)` per successful call:

- Spike test first: a sample `> SpikeFactor × estimate` is excluded from the
  EWMA but counted — one-shot lag never poisons the estimate while systemic
  drift is still followed (`TestEstimatorFollowsDriftRejectsSpikes`).
- Non-spikes update `estimate = (1−Alpha)×estimate + Alpha×sample`.
- Every sample marks its slot in the ring window. No breach verdicts fire
  until the window is full; the signal returns true exactly once when the
  windowed spike rate first exceeds `BreachRate`, then re-arms only after the
  rate falls back under (`TestEstimatorBreachSignalFiresOnceAndRearms`).

`Estimate()` returns the current value; `Stats()` returns estimate, rolling
spike rate, and lifetime sample/spike counts. All state is guarded by a
`sync.Mutex`, so one estimator can be shared across goroutines.

## Connections

Its output is the `secondsPerPoint` input to [[route-gate]] and the
`SecondsPerPoint` on each tier of [[route-tiered]] — one estimator per tier
([[points-latency-model]]). What may feed it and who moves the recorded
baseline are [[calibration-doctrine]]; wiring samples and acting on the
breach signal are host work ([[host-obligations]]).

## Operational notes

Process-lifetime only: persist calibration under a human's hand and re-seed
on restart, so the recorded baseline never drifts silently. Failures must
not be sampled — a failed call's duration says nothing about how long a
landed answer takes.
