---
name: landing-outcomes
description: Lease.Land and the Landing record — Landed | Superseded | Stale, supersession first, every outcome recorded and reproducible.
kind: component
sources:
  - lease.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Landing outcomes — the second clearance

`Lease.Land(currentGen uint64, actualDrift float64) Landing` in `lease.go`
validates the lease against the world at landing time. It is pure, like
`Route`, and returns a `Landing` — the ex-post sibling of `Verdict`, carrying
the outcome plus every number that produced it.

## How it works

The ruling is a three-case switch in strict priority order:

1. `currentGen != l.Gen` → **Superseded** — a salient event bumped the
   generation mid-flight. Supersession trumps staleness: it wins regardless
   of how little drift accrued (`TestLeaseLand` pins the gen-bump +
   over-budget case to Superseded).
2. `actualDrift > l.Budget` → **Stale** — no supersession, but actual drift
   exceeded the budget the gate allowed against. The gate's prediction was
   wrong; the comparison is strict, so drift exactly at budget still lands.
3. otherwise → **Landed** — the world still matches; apply the answer.

`Outcome` is an int enum with a `String()` form (`landed` / `superseded` /
`stale` / `unknown`). `Landing` records `Outcome`, `Class`, `Gen`,
`CurrentGen`, `Budget`, `PredictedDrift` (the gate's frozen ex-ante bet), and
`ActualDrift` — so any logged Landing is reproducible from its own fields
(`TestLandingReproducible`).

Doctrine attached to the outcomes: every outcome — not just Landed — must be
recorded; a dropped answer without a trace is a mystery-bug factory. A
Superseded answer must not be applied as-is, but hosts may warm a cache with
it. A *pattern* of Stale outcomes indicts the velocity estimate or the tier
calibration.

## Connections

Validates the read set frozen by [[lease-lifecycle]]. The
predicted-vs-actual pair in one record is the calibration report described in
[[drift-unit]] — no join to the [[verdict]] log needed. Recording outcomes,
bumping generations, and acting on Stale patterns are host obligations
([[host-obligations]]); reproducibility is invariant 1 of
[[invariants-and-purity]]. The warm-cache fate of Superseded answers is one
of the starvation escapes noted in [[suppression-debt]].

## Operational notes

`Budget` in the landing is the *effective* (possibly debt-inflated) budget
from checkout — Stale is judged against what the gate actually allowed, not
the base constant. Extending `Landing` requires extending
`TestLandingReproducible` ([[test-suite]]).
