---
name: invariants-and-purity
description: The six invariants that make an implementation the pattern — purity, mandatory degrade, record everything, successes-only, shared unit, fail pessimistic.
kind: concept
sources:
  - DESIGN.md
  - CLAUDE.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Invariants — what makes it the pattern, not a knockoff

DESIGN.md enumerates six invariants; CLAUDE.md restates the load-bearing ones
as non-negotiables for contributors. They are the contract any conforming
implementation must keep, and the reason the gate can be trusted to govern a
model while containing none.

## How it works

1. **Purity.** `Route`, `RouteTiered`, and `Land` are pure — no clock reads,
   no randomness, no model, no hidden state. Property-tested: any logged
   verdict or landing is reproducible from its own recorded fields
   ([[test-suite]]).
2. **Mandatory degrade.** Every class has a legal degrade mode; registration
   without one fails `Validate` ([[degrade-modes]], [[decision-class]]).
3. **Record everything.** Every route decision — especially suppressions —
   and every landing outcome — especially non-landings — is recorded
   ([[verdict]], [[landing-outcomes]]).
4. **Successes only.** The latency estimator learns from successes; failures
   feed a breaker, never the estimator ([[calibration-doctrine]]).
5. **Shared unit.** Predicted and actual drift share a unit; their paired log
   is the gate's calibration report ([[drift-unit]]).
6. **Fail pessimistic.** Uncalibrated fails toward the floor, never toward
   stale action ([[velocity-semantics]], [[calibration-doctrine]]).

Supporting doctrine from CLAUDE.md: suppression debt is passed in and
recorded, never accrued inside the gate ([[suppression-debt]]); budgets and
points are reviewed constants with no runtime tuning API; the package stays
a stdlib-only leaf — velocity estimation, persistence, breakers, and
transport belong to hosts ([[host-obligations]]).

## Connections

Purity is why the gate is unit-testable, replayable, immune to prompt
injection, and cannot hallucinate itself open — the second of the three
properties DESIGN.md uses to distinguish the gate from a rate limiter
([[overview]]). The published pattern write-up claims exactly this packaging
([[lineage-and-prior-art]]).

## Operational notes

Extending `Verdict` or `Landing` requires extending the corresponding
reproducibility test — `TestVerdictReproducible` is named in CLAUDE.md as the
guard to grow. `go test -race ./...` is the definition of done for any
change.
