---
name: suppression-debt
description: Starvation relief — host-accrued debt inflates the effective budget (Budget × (1 + debt × DebtFactor)) until a call is admitted.
kind: concept
sources:
  - route.go
  - class.go
  - DESIGN.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Suppression debt — the anti-starvation queue

At persistently high velocity the oracle is never consulted and the system
flatlines at floor quality. Suppression debt is landfall's primary relief:
the host accrues debt per suppressed call, the gate inflates the effective
budget by it, and a landed call resets it — guaranteeing eventual
consultation without breaking the gate's purity.

## How it works

`Route` computes `EffectiveBudget = Budget × (1 + debt × DebtFactor)`
(`route.go`), where `debt` is passed in by the host and `DebtFactor` is
declared on the [[decision-class]]. The verdict records both the base budget
and the debt, so suppressed-then-admitted decisions remain reproducible.

- `DebtFactor: 0` declares starvation acceptable for a class — the budget
  never inflates regardless of accrued debt (`TestRouteDebtRelievesStarvation`
  covers both the admission-by-debt case and the zero-factor case).
- Debt is *passed in and recorded*, never accrued inside the gate — purity
  and reproducibility survive ([[invariants-and-purity]]).

Accrual and reset discipline is the host's ([[host-obligations]]): plain
`Route` keeps a single queue per class — accrue on suppression, reset on
landing. Under `RouteTiered` the discipline is per-tier
([[per-tier-debt]]).

DESIGN.md names two composable escapes alongside debt: speculative warm calls
(run, don't act, warm a cache — the natural fate of a Superseded answer, see
[[landing-outcomes]]) and the zero-velocity pause sandbox
([[velocity-semantics]]).

## Connections

The inflation term appears in every [[verdict]] as `EffectiveBudget` and
`Debt`; `Checkout` freezes the *effective* budget into the lease, so a
debt-admitted call is landing-validated against the same inflated tolerance
it was admitted under ([[lease-lifecycle]]). The theoretical frame — Neely's
Lyapunov virtual queues — is in [[lineage-and-prior-art]].

## Operational notes

Debt units are counts of suppressions, weighted by the class-declared
`DebtFactor`. The guarantee is eventual consultation, not bounded delay; a
class that must never starve needs a nonzero factor sized against its budget.
