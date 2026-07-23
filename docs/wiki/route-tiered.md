---
name: route-tiered
description: RouteTiered — the comparative gate over a quality-ordered tier ladder; best tier whose drift fits its own debt-inflated budget.
kind: component
sources:
  - route.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# RouteTiered — the comparative gate

`RouteTiered` in `route.go` is the comparative form of the gate: route to the
best-quality tier whose predicted drift fits that tier's debt-inflated
budget — the big slow oracle when the world is calm, a faster one when it is
hot, the degrade floor when nothing lands. It is the wired form of the
faster-tier degrade mode the source implementation declared but never built.

## How it works

Signature: `RouteTiered(c Class, velocity float64, tiers []Tier) Verdict`.

`Tier` carries `Name`, `SecondsPerPoint` (the host's current pace estimate
for that oracle), `Quality` (orders tiers best-first in the caller's
judgment), and `Debt` (that tier's own suppression-debt queue,
[[per-tier-debt]]).

Selection: run plain `Route` ([[route-gate]]) once per tier with that tier's
pace and debt; among the tiers that allow, choose the highest quality. The
final verdict is `Route` re-run on the chosen tier with `Tier` set to its
name. Edge cases:

- **No tiers offered** — a suppression with Arithmetic
  `no tiers offered - suppressed`.
- **No tier fits** — a suppression carrying the best-quality tier's
  arithmetic, prefixed `best of N tiers (<name>): …`, with `Tier` left empty.

`TestRouteTiered` pins the calm/hot/frantic behavior: quality 3 cloud tier at
low velocity, quality 1 local tier when hot, suppression when nothing fits.

## Connections

Each per-tier check is exactly [[route-gate]] — same purity, same
[[velocity-semantics]]. The chosen tier's pace and debt are recorded in the
[[verdict]], keeping tiered verdicts reproducible
(`TestRouteTieredVerdictReproducible`). Per-tier pace estimates come from one
[[estimator]] per tier ([[points-latency-model]]); per-tier debt discipline
is [[per-tier-debt]]. The faster-tier degrade mode ([[degrade-modes]])
describes this component's behavior as a class-level declaration.

## Operational notes

Quality is the caller's judgment, not measured; ties in quality resolve to
the earlier tier in the slice (strict `>` comparison). One tier's debt never
inflates another tier's budget. A suppressed tiered verdict still carries the
class degrade mode, so hosts fall through to the declared floor.
