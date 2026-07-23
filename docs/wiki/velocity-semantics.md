---
name: velocity-semantics
description: Velocity edge doctrine — 0 means paused and is always allowed; negative/NaN/Inf are unmeasurable and suppress.
kind: concept
sources:
  - route.go
  - DESIGN.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Velocity semantics — zero is paused, unmeasurable is refused

Velocity is the measured rate of decision-relevant change, in drift-units per
wall-second, fed to the gate by the host. Two edges are doctrine, and they
deliberately invert the source implementation's encoding.

## How it works

In `Route` (`route.go`):

- **Velocity 0 is a paused world.** Zero staleness: predicted drift is
  `anything × 0 = 0`, which always fits the budget. A paused world is always
  allowed — the pause sandbox is one of the composable starvation escapes
  named in DESIGN.md.
- **Negative, NaN, or ±Inf velocity is unmeasurable.** Checked explicitly
  before the drift multiplication; the verdict is a suppression whose
  Arithmetic string records the unmeasurable value. Guarded by
  `TestRouteUnmeasurableVelocitySuppresses`, which also requires the
  suppressed verdict to carry the class's degrade mode.

The inversion: promptworld encoded "uncapped speed" as `ticksPerSecond ≤ 0`
and suppressed on it. landfall's velocity is the measured drift rate, so 0
legitimately means "nothing is changing" — and the pessimism is preserved by
routing the genuinely unknowable encodings (NaN/Inf/negative) to suppression
instead. An uncalibrated world fails toward the degrade floor, never toward
stale action.

## Connections

Applied inside [[route-gate]] and therefore per-tier inside [[route-tiered]].
The unit velocity is measured in is [[drift-unit]]. Measuring velocity is
explicitly the host's obligation ([[host-obligations]]) — the gate records
what it was fed in the [[verdict]] rather than estimating anything itself.
The fail-pessimistic stance is shared with the estimator's seeding rules
([[calibration-doctrine]]).

## Operational notes

Hosts wanting a "pause" affordance get it for free: report velocity 0 while
paused and every class is admitted with zero drift. Hosts whose velocity
estimator can fail should let it fail to NaN rather than 0 — the two encode
opposite decisions. Velocity-estimator helpers (salient-event EWMA,
state-hash churn/sec) are an open question in DESIGN.md, deliberately not
shipped at this pin.
