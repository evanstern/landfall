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

This threshold shape is not folk engineering: in remote estimation it is the
provably optimal sampling policy (Sun, Polyanskiy & Uysal; Ornee & Sun — see
[Prior art & theory](#prior-art--theory)).

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

- Checked out at launch (`Checkout(verdict, gen)`): world generation, the
  effective budget allowed against, and the predicted drift the gate bet on —
  the OCC read set, frozen into the lease.
- Salient events bump the host's generation.
- `Land(currentGen, actualDrift)` → a `Landing` record: `Landed | Superseded
  | Stale` plus the numbers that produced it, reproducible from its own
  fields like a `Verdict`.
- Supersession trumps staleness. Every outcome is recorded — a dropped
  answer without a trace is a mystery bug factory.
- `Stale` means the gate's prediction was wrong. A pattern of them indicts
  the velocity estimate or the tier calibration — and because the lease
  carries the prediction to the landing, every `Landing` pairs predicted
  vs. actual drift in one record: the calibration report is the landing log
  itself, no join back to the verdict log.

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

**Debt under `RouteTiered` — per-constraint queues (doctrine).** Debt queues
follow Neely's virtual-queue discipline: one queue per constraint, and under
`RouteTiered` each tier is its own constraint — "consulted at quality ≥ this
tier's, often enough." Each `Tier` carries its own `Debt`; the gate tests each
tier against its *own* debt-inflated budget, and one tier's queue never
inflates another tier's budget. Consequently **a landed local call does not
repay the debt owed to the suppressed cloud tier**: in a persistently hot
world the class keeps being served locally while the cloud queue keeps
growing, until the cloud tier's own inflated budget admits a cloud call —
relief from quality starvation, not just consultation starvation. A landing
at tier T repays the queues of T and every lower-quality tier (they were
served at ≥ their quality); queues of higher-quality tiers survive. The
lowest tier's queue still grows under total suppression, so the plain
eventual-consultation guarantee is preserved. Accrual and reset are the
host's (see obligations below); the gate stays pure — per-tier debt is
passed in on the `Tier` and the chosen tier's debt is recorded in the
verdict.

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

Debt accrual under `RouteTiered` follows the per-constraint doctrine: on
every verdict, the host increments the queue of each tier whose quality
exceeds the landed tier's (all tiers, when nothing lands); when a call lands
at tier T, the host resets the queues of T and every lower-quality tier.
Plain `Route` keeps a single queue per class: accrue on suppression, reset on
landing.

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

## Prior art & theory

landfall is the engineering packaging of proven results, not a new theory.
Every element has a lineage, and the design claims that lineage explicitly
rather than presenting the arithmetic as folk wisdom. (Grounding with full
citations: `research/Staleness-Budgeted-Admission-Control/`.)

- **The problem statement is AI metareasoning.** Boddy & Dean (1988–1994)
  framed time-dependent planning: the world changes while the agent
  deliberates, so deliberation itself must be scheduled against that drift.
  Russell & Wefald (1991) formalized the value of computation — run a
  computation only if its expected utility gain exceeds its time cost. The
  gate is the budget-form linearization of exactly that test.
- **The gate formula has optimality proofs behind it.** Sun, Polyanskiy &
  Uysal ([arXiv:1701.06734](https://arxiv.org/abs/1701.06734)) prove that
  when sampling a Wiener process observed through a random-delay channel,
  the policy minimizing mean-square estimation error is a threshold on how
  much the process drifts during the service time — structurally the same
  object as `velocity × time-to-land ≤ budget`. Ornee & Sun extend the
  threshold to mean-reverting (Ornstein–Uhlenbeck) processes, and later
  work learns it online under unknown delay statistics
  ([arXiv:2308.15401](https://arxiv.org/abs/2308.15401)). Self-triggered
  control (Heemels, Johansson & Tabuada, IEEE CDC 2012) is the same
  ex-ante stance in control theory: compute the next sampling instant
  ahead of time from a model, rather than reacting after deviation.
- **The drift unit is the Age-of-Incorrect-Information stance.** Plain Age
  of Information ages from delivery regardless of content; AoII
  ([arXiv:1907.06604](https://arxiv.org/abs/1907.06604)) accrues age only
  while the monitor's picture is actually *wrong*. "Expected salient events
  during flight" is that content-aware refinement: the budget prices
  decision-relevant change, not elapsed time, which is what lets the same
  unit serve as both admission tolerance and landing measurement.
- **Suppression debt is a Lyapunov virtual queue.** In Neely's
  drift-plus-penalty framework
  ([arXiv:1008.3519](https://arxiv.org/abs/1008.3519)), each time-average
  constraint gets a virtual queue that grows while the constraint is
  violated and inflates its priority until it is served, with proven
  stability bounds. Debt is that counter with
  `Budget × (1 + debt × DebtFactor)` as the inflation term — accumulate
  while suppressed, spend by landing. OS scheduler aging is the same
  bounded-starvation invariant with simpler arithmetic. The per-tier debt
  doctrine under `RouteTiered` (see Starvation) is this discipline applied
  literally: one queue per constraint, each tier a constraint.
- **The lease is optimistic concurrency control.** Kung & Robinson (1981):
  work in private, validate against what committed meanwhile, apply only if
  validation passes. Ask-gate → flight → land-validate is that three-phase
  shape — the oracle call is the read phase, landing is validation,
  `Superseded`/`Stale` are the aborts. TL2's global version clock (Dice,
  Shalev & Shavit 2006) is the concrete generation-counter discipline the
  `Lease` uses, and the term itself is Gray & Cheriton (1989): a
  time-bounded grant whose expiry bounds how stale a holder's view can be.

**The honest gap: doctrine constants vs. derived optimal thresholds.** The
optimality results above *derive* their thresholds from process statistics;
landfall's budgets are doctrine — reviewed constants, deliberately not
optimized (see "Learned budgets" under out-of-scope). Measured against
Sun/Ornee, doctrine constants carry regret: a derived threshold would track
the process better. They are kept anyway because they buy what the derived
threshold can't: budgets that are auditable, diffable in review, and
meaningful to a human ("this decision survives 5 salient events"). The
calibration log (invariant 5) is the monitor on that tradeoff — when
predicted-vs-actual drift shows systematic error, the constant is wrong and
a human changes it. What landfall claims as its own is not the theory but
the packaging: the pure gate with no model in the enforcement path, the
verdict that carries its own arithmetic, and degradation forced to exist at
registration.

## Open questions

- Velocity estimators worth shipping as optional helpers: salient-event EWMA,
  state-hash churn/sec.

Resolved: *should `Lease` carry predicted drift?* — yes. The lease is the OCC
read set; recording the admission-time prediction in it lets every landing
pair predicted vs. actual drift without a join (invariant 5). See the
lifecycle section.
