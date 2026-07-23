---
name: per-tier-debt
description: Debt under RouteTiered — one virtual queue per tier; a landed local call never repays the cloud tier's queue.
kind: concept
sources:
  - route.go
  - DESIGN.md
  - route_test.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Per-tier debt — one queue per constraint

Under `RouteTiered`, suppression debt follows Neely's virtual-queue
discipline: one queue per constraint, and each tier is its own constraint —
"consulted at quality ≥ this tier's, often enough." This is doctrine recorded
in DESIGN.md and carried on the `Tier.Debt` field in `route.go`.

## How it works

Each `Tier` carries its own `Debt`; the gate tests each tier against its
*own* debt-inflated budget, and one tier's queue never inflates another
tier's budget. Consequences, all pinned by `TestRouteTieredPerTierDebt`:

- **A landed local call does not repay the cloud tier's debt.** In a
  persistently hot world the class keeps being served locally while the cloud
  queue grows, until the cloud tier's own inflated budget admits a cloud
  call — relief from *quality* starvation, not just consultation starvation.
  The test shows cloud debt 16 (× factor 0.5) pulling routing back up-tier.
- **Isolation both ways** — arbitrarily large local debt never admits or
  affects the cloud tier.
- **A tier's own debt admits it** even in a frantic world where nothing fits
  debtless.
- **DebtFactor 0** disables inflation for every tier's queue at once — the
  class-level declaration that starvation is acceptable.

Host discipline (from DESIGN.md's obligations section): on every verdict,
increment the queue of each tier whose quality exceeds the landed tier's (all
tiers when nothing lands); when a call lands at tier T, reset the queues of T
and every lower-quality tier — they were served at ≥ their quality. The
lowest tier's queue still grows under total suppression, preserving the
plain eventual-consultation guarantee.

## Connections

Specializes [[suppression-debt]] for [[route-tiered]]; the chosen tier's debt
is recorded in the [[verdict]], preserving reproducibility. Accrual and reset
are host work ([[host-obligations]]); the gate stays pure
([[invariants-and-purity]]). The Lyapunov virtual-queue lineage is in
[[lineage-and-prior-art]].

## Operational notes

The doctrine trades simplicity for a stronger guarantee: hosts must maintain
N queues per class instead of one, but a hot world can no longer pin a class
to its cheapest tier forever.
