# landfall — design

Staleness-budgeted admission control for slow oracles. Generalized from
promptworld's cognition horizon (`internal/cognition`), with the pattern's
sibling steals folded in — verdict-carries-its-arithmetic, the landing
lease, successes-only estimator feeding — because they turn out to be one
machine, not four patterns.

## The pattern

A fast loop advances the world. A slow, expensive oracle (LLM, heavy planner,
remote model, human) produces answers with latency. Every oracle query is
implicitly a bet: *the world at landing time will still resemble the world at
ask time closely enough for the answer to apply.* landfall prices that bet
**before spending the call**, using only local deterministic signals:

```
predicted drift = world velocity × predicted time-to-land
allow iff predicted drift ≤ tolerance(decision class)
```

Three properties distinguish it from a rate limiter:

1. **Ex-ante, not ex-post.** Timeouts and fallbacks cope after you've paid.
   The gate refuses bets it can predict will lose — and thereby breaks latency
   spirals (stale calls congest the pool → latency rises → more calls stale).
   Classical admission control gates on cost or capacity; this gates on
   whether the answer will still be *true*.
2. **The gate that governs the model contains no model.** `Route` is a pure
   function: no randomness, no wall-clock reads, no model. Unit-testable,
   replayable, immune to prompt injection, cannot hallucinate itself open.
3. **It forces the degraded mode to exist.** A class cannot be registered
   without a legal `Degrade`; every slow-oracle dependency acquires a declared
   always-available fallback. Graceful degradation by construction.

## Vocabulary

| promptworld | landfall | definition |
|---|---|---|
| sim speed × tick rate | **velocity** | drift-units per wall-second of decision-relevant change |
| points × sec/pt | **time-to-land** | predicted wall seconds for the call |
| staleness in ticks | **predicted drift** | velocity × time-to-land |
| BudgetTicks | **budget** | drift-units the answer survives; its half-life |
| degrade floor | **Degrade** | skip \| reflex \| template \| faster-tier |
| cog.* trail | **Verdict** | the decision plus the arithmetic that produced it |

**Canonical drift unit:** expected salient events during flight. Budget 5 =
"this decision survives 5 salient events." Domain-portable, and identical to
the unit the landing check counts *actually* — so the verdict log doubles as
the gate's calibration report.

## What the source implementation taught (and landfall keeps)

Reading promptworld's real `internal/cognition` surfaced three ideas the
first-draft generalization missed:

1. **Factor latency as shape × host.** `Points` (Fibonacci, a property of the
   prompt shape, host-independent) × `secondsPerPoint` (a property of the
   tier, calibrated and live-estimated). One estimator per tier instead of
   per (oracle, class), and a brand-new class gets latency predictions with
   zero observations — assign points, inherit the tier's calibration. For a
   multi-call tool-use loop, seconds-per-point covers the *whole loop*, the
   same unit the estimator observes, so the arithmetic stays truthful.
2. **The reflex floor is a spectrum.** `Degrade` is an enum — skip (recorded,
   never silent), reflex, template, faster-tier — not a single mandatory
   callback. promptworld declared faster-tier and never wired it; landfall
   wires it as `RouteTiered`, the comparative form of the gate: route to the
   best-quality tier whose drift fits, big-slow when calm, small-fast when
   hot, floor when nothing lands.
3. **Calibration doctrine.** Bootstrap defaults are deliberately pessimistic —
   an uncalibrated world fails toward the floor, never toward stale action.
   The live estimator is EWMA with spike rejection (one-shot lag excluded but
   counted; systemic drift followed) plus a breach signal when the spike rate
   over a window says the baseline is lying. The *recorded* baseline moves
   only under a human's hand; the daemon never writes it.

One deliberate inversion: promptworld encodes "uncapped speed" as
`ticksPerSecond ≤ 0` and suppresses. landfall's velocity is the measured
drift rate, so **0 means paused — zero staleness, always allowed** — and an
unmeasurable rate (NaN/Inf/negative) suppresses, preserving the pessimism.

## The lifecycle: ask-gate → flight → land-validate

The gate is ex-ante; even an allowed call can be invalidated mid-flight.
`Lease` is the ex-post complement (promptworld's landing ladder +
generation-bump supersession, distilled):

- Checked out at launch: world generation + the effective budget allowed against.
- Salient events bump the host's generation.
- `Land(currentGen, actualDrift)` → `Landed | Superseded | Stale`.
- Supersession trumps staleness. Every outcome is recorded — a dropped
  answer without a trace is a mystery bug factory.
- `Stale` means the gate's prediction was wrong. A pattern of them indicts
  the velocity estimate or the tier calibration; predicted vs. actual drift
  is one query over the logs.

## Starvation

At persistently high velocity the oracle is never consulted and the system
flatlines at floor quality — promptworld hit this as the horizon-vs-iteration
tension. landfall's primary relief is **suppression debt**: the host accrues
debt per suppressed call, `Route` inflates the effective budget by
`Budget × (1 + debt × DebtFactor)`, and a landed call resets it. Guarantees
eventual consultation; `DebtFactor: 0` declares starvation acceptable for a
class. Debt is passed in and recorded in the verdict, so purity and
reproducibility survive. Composable escapes: speculative warm calls (run,
don't act, warm a cache — the natural fate of a `Superseded` answer) and the
zero-velocity pause sandbox (pause is zero staleness by doctrine).

## Where it applies

- **LLM agents on live state** — assistant on an actively-edited buffer:
  velocity = edit events/sec; is the 20s deep pass worth it or will the user
  type through the diff?
- **Trading** — velocity = realized volatility; budget = edge half-life. The
  formula is identical.
- **Robotics** — heavy planner vs. reactive controller; budget in trajectory
  divergence.
- **Human-in-the-loop** — an approver is an oracle with hours of latency;
  gate whether the approval will still be actionable, else route to
  policy-default. Approval queues everywhere are full of stale asks nobody
  priced.
- **Model-tier routing** — `RouteTiered`, above.

## Invariants (what makes an implementation the pattern, not a knockoff)

1. `Route`/`RouteTiered`/`Land` are pure — property-tested: any logged
   verdict is reproducible from its own recorded fields.
2. Every class has a legal degrade mode; registration without one fails.
3. Every route decision — **especially suppressions** — and every landing
   outcome — **especially non-landings** — is recorded.
4. The latency estimator learns from successes only; failures feed a breaker,
   never the estimator.
5. Predicted drift and actual drift share a unit; their paired log is the
   gate's calibration report.
6. Uncalibrated fails pessimistic: toward the floor, never toward stale action.

## The host's obligations

landfall is deliberately a leaf. The host owns: measuring velocity (estimate
it *outside* the gate, pass a number in — the estimator is separately
testable and the verdict records what it was fed), bumping the world
generation on salient events, accruing/resetting debt, persisting verdicts
and outcomes, calibration files, and the circuit breaker.

## Out of scope, deliberately

- **Cost-aware gating** (`E[value] − floor_value > $cost`) — the budget form
  is the practical linearization; add money if a domain demands it.
- **Learned budgets.** Budgets are doctrine: reviewed constants. The
  calibration log shows when declarations are wrong; closing that loop
  automatically would be a model creeping back into the enforcement path
  through the side door. Resist until the log proves systematic error.
- **Transport concerns** — retries, breakers, pooling. Pair with a
  busy-≠-down breaker; landfall only insists failures don't feed the
  estimator.

## Open questions

- Debt semantics under `RouteTiered`: does a landed *local* call repay debt
  owed to the suppressed *cloud* class? Currently the host decides; may
  deserve doctrine.
- Should `Lease` carry predicted drift too, so `Stale` outcomes can log
  predicted-vs-actual without a join?
- Velocity estimators worth shipping as optional helpers: salient-event EWMA,
  state-hash churn/sec.
