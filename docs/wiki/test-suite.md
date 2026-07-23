---
name: test-suite
description: What each test file guards — boundary and doctrine tests, the reproducibility properties, and the go test -race definition of done.
kind: concept
sources:
  - route_test.go
  - lease_test.go
  - estimate_test.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Test suite — the guards on the doctrine

The tests are doctrine made executable: each named invariant or semantic edge
has a test that pins it. `go test -race ./...` is the definition of done for
any change (CLAUDE.md). The shared fixture is `planner` — a class with
3 points, budget 10, reflex degrade, DebtFactor 0.5.

## How it works

`route_test.go`:

- `TestRouteBoundary` — paused world allowed, drift exactly at budget
  allowed (inclusive comparison), over-budget suppressed ([[route-gate]]).
- `TestRouteUnmeasurableVelocitySuppresses` — negative/NaN/Inf suppress and
  the suppression still carries the degrade mode ([[velocity-semantics]]).
- `TestRouteDebtRelievesStarvation` — debt admits a previously suppressed
  call; DebtFactor 0 never inflates ([[suppression-debt]]).
- `TestVerdictReproducible` — the pattern's defining property: rebuild the
  class from the verdict's own fields, re-run `Route`, require an identical
  struct; also asserts `Checkout` freezes the verdict's arithmetic
  ([[verdict]], [[lease-lifecycle]]).
- `TestRouteTiered` / `TestRouteTieredPerTierDebt` /
  `TestRouteTieredVerdictReproducible` — calm/hot/frantic tier selection,
  per-tier queue isolation and repayment doctrine, tiered reproducibility
  ([[route-tiered]], [[per-tier-debt]]).
- `TestClassValidate` — all six rejection cases ([[decision-class]]).

`lease_test.go`:

- `TestLeaseLand` — the outcome switch including supersession trumping
  staleness and drift-at-budget landing ([[landing-outcomes]]).
- `TestLandingCarriesCalibrationPair` — a Stale landing carries predicted vs.
  actual drift in its own record, no join ([[drift-unit]]).
- `TestLandingReproducible` — the reproducibility property extended to the
  ex-post record.

`estimate_test.go`:

- `TestEstimatorFollowsDriftRejectsSpikes` — a one-shot spike leaves the
  estimate untouched; systemic drift converges ([[estimator]]).
- `TestEstimatorBreachSignalFiresOnceAndRearms` — no breach until the window
  fills, exactly one firing per breach, re-arm on recovery.

## Connections

The reproducibility tests are the executable form of invariant 1
([[invariants-and-purity]]); CLAUDE.md requires extending
`TestVerdictReproducible` when adding verdict fields. Calibration behavior
under test mirrors [[calibration-doctrine]].

## Operational notes

All tests are table-driven or scenario-driven over the pure functions — no
mocks, no clocks, no goroutines needed except what `-race` exercises through
the estimator's mutex.
