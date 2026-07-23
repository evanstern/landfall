---
name: host-obligations
description: Everything landfall deliberately leaves outside — velocity measurement, generation bumps, debt bookkeeping, persistence, breakers, calibration files.
kind: concept
sources:
  - DESIGN.md
  - README.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Host obligations — the deliberate leaf boundary

landfall is deliberately a leaf package: it decides and records, the host
does everything else. DESIGN.md enumerates the host's side of the contract;
keeping these outside is what keeps the gate pure and separately testable.

## How it works

The host owns:

- **Measuring velocity** — estimate it outside the gate and pass a number in;
  the [[verdict]] records what it was fed ([[velocity-semantics]]).
- **Bumping the world generation** on salient events — the counter that
  [[lease-lifecycle]] freezes and [[landing-outcomes]] compares.
- **Accruing and resetting suppression debt.** Plain `Route`: one queue per
  class — accrue on suppression, reset on landing ([[suppression-debt]]).
  `RouteTiered`: per-tier queues — on every verdict increment each tier whose
  quality exceeds the landed tier's (all tiers when nothing lands); a landing
  at tier T resets T and every lower-quality tier ([[per-tier-debt]]).
- **Persisting verdicts and outcomes** — the record-everything invariant
  ([[invariants-and-purity]]) names the artifacts; the host provides the log.
- **Calibration files** — seeding each tier's [[estimator]], persisting
  baselines under a human's hand, re-seeding on restart, and acting on the
  breach signal ([[calibration-doctrine]]).
- **The circuit breaker** — failures, retries, pooling, and transport;
  landfall only insists failures never feed the estimator.
- **Running the degrade path** when a verdict suppresses
  ([[degrade-modes]]), and validating every [[decision-class]] at startup.

## Connections

This boundary is the third distinguishing property in [[overview]]'s framing:
the gate that governs the model contains no model — because everything that
observes, learns (beyond the estimator), or persists lives on the host side
of this line.

## Operational notes

Velocity-estimator helpers (salient-event EWMA, state-hash churn/sec) are an
acknowledged open question in DESIGN.md — candidates for optional shipped
helpers, not yet present at this pin. Until then every host writes its own.
