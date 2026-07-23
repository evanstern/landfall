# Landfall: staleness-budgeted admission control for slow oracles

**A named engineering pattern, published as a defensive disclosure.**

- Author: Evan Stern
- First published: 2026-07-23, in the public repository
  <https://github.com/evanstern/landfall>
- Status: defensive publication. This document is placed in the public record
  deliberately, as prior art. No patent will be sought on the pattern
  described here (decision recorded in [PRIOR-ART.md](PRIOR-ART.md)).
- Reference implementation: this repository (Go, stdlib-only). Design
  rationale and invariants: [DESIGN.md](DESIGN.md).

This write-up is self-contained: it names the pattern, defines its five
elements and their lifecycle, cites the lineage of each element, and states
what is and is not claimed as new.

## Abstract

A fast loop advances the world. A slow, expensive oracle — an LLM, a heavy
planner, a remote model, a human approver — answers with latency. Every
oracle query is implicitly a bet that the world at *landing* time will still
resemble the world at *ask* time closely enough for the answer to apply.
**Landfall** is the pattern that prices this bet before the call is spent and
re-validates it when the answer lands, using only deterministic local
arithmetic — no model, no learned policy, no clock reads inside the gate —
so that every decision the gate makes is reproducible from its own logged
fields. The individual mechanisms are all established results in adjacent
fields, several with optimality proofs; what this document names and fixes is
their packaging as one auditable protocol: **ask-gate → flight →
land-validate**.

## The problem

Timeouts, retries, and circuit breakers cope with a slow dependency *after*
the cost is paid, and they protect against a dependency that is *down* or
*late* — not against one whose answer is *no longer true*. When the caller's
world drifts while the oracle thinks, an answer can arrive on time, well-formed,
and useless. Worse, stale calls congest the oracle, which raises latency,
which makes the next call staler: a latency spiral. Classical admission
control gates on cost or capacity. The missing gate is one that refuses a
call because its answer is predicted to be **stale on arrival**.

## The pattern: five elements

A conforming implementation packages all five. Each element is individually
established in an adjacent literature (see [Lineage](#lineage)); the claim
here is only to the named combination.

**E1 — The ex-ante drift gate.** Before the call is issued:

```
predicted drift = world velocity × predicted time-to-land
allow  iff  predicted drift ≤ budget(decision class)
```

*Velocity* is the measured rate of decision-relevant change, in drift-units
per wall-second. The *budget* is a per-class constant declaring how much
drift the class's answers survive — its half-life. The canonical drift unit
is **expected salient events during flight**: budget 5 means "this decision
survives 5 salient events." The gate is a pure function: same inputs, same
verdict. Suppression is a first-class outcome, recorded with the arithmetic
that produced it.

**E2 — Latency factored as shape × host.** Predicted time-to-land is
`points × seconds-per-point`, where *points* measure the call's shape (a
property of the prompt/request, host-independent) and *seconds-per-point* is
per-tier calibration maintained by a live estimator. The estimator learns
from successes only — failures belong to the host's circuit breaker — and is
seeded pessimistically, so an uncalibrated system fails toward its degraded
floor, never toward stale action.

**E3 — A declared degrade taxonomy.** No decision class may be registered
without a legal degraded mode: *skip* (recorded, never silent), *reflex* (a
fast deterministic answer), *template* (a canned answer), or *faster-tier*
(route the same ask to a quicker, lower-quality oracle). Degradation exists
by construction, not by incident retrospective. The comparative form of the
gate — route to the best-quality tier whose predicted drift fits the budget —
falls out of E1 + E3.

**E4 — Suppression debt.** Each suppressed call accrues debt for its class;
the gate inflates the effective budget by `budget × (1 + debt × factor)`, and
a landed call resets the debt. This bounds starvation: at persistently high
velocity the oracle is still eventually consulted. Debt is passed *into* the
gate and recorded in the verdict — never accrued inside it — so purity
survives.

**E5 — The generation lease and landing validation.** An allowed call checks
out a lease: the world's generation counter at launch plus the effective
budget it was admitted against. Salient events bump the generation. When the
answer lands, validation yields one of three recorded outcomes: **Landed**
(apply it), **Superseded** (a salient event won mid-flight; don't act on it —
at most warm a cache), or **Stale** (actual drift exceeded the budget: the
gate's prediction was wrong). Predicted and actual drift share a unit, so the
paired log is the gate's calibration report, and any logged verdict is
reproducible from its own recorded fields — *verdict-carries-arithmetic*.

## The lifecycle

```
            ┌── E1+E2: gate ──┐
   ask ─────┤ allow           ├──── flight (oracle thinks, world moves) ──── E5: land
            │ suppress → E3   │                                              │
            │ (accrues E4)    │                                    Landed / Superseded / Stale
            └─────────────────┘                                    (all three recorded)
```

Ask-gate → flight → land-validate. The gate is ex-ante; even an allowed call
can be invalidated mid-flight, which is why E5 exists. The three-phase shape
is optimistic concurrency control with the oracle call as the read phase.

## The differentiator: deterministic and auditable

The gate that governs the model contains no model. `Route` and `Land` are
pure functions — no randomness, no wall-clock reads, no learned policy, no
hidden state. Consequences:

- Every verdict is **reproducible from its own logged fields** — the audit
  log is not a narrative about the decision, it *is* the decision.
- The gate is unit-testable and replayable, cannot be prompt-injected, and
  cannot hallucinate itself open.
- Budgets are **doctrine**: reviewed constants, diffable in code review,
  meaningful to a human — deliberately not optimized at runtime. Measured
  against the optimal-threshold results below, doctrine constants carry
  regret; they are kept because they buy auditability, and the E5
  calibration log is the standing monitor on that tradeoff.

The contemporary LLM-agent systems converging on this decision space use
*learned* gates. Landfall's position is the opposite corner, on purpose: the
determinism is not an implementation shortcut but the load-bearing property.

## Lineage

None of the five elements is new in isolation. The pattern claims its
lineage explicitly (full survey with per-element closeness:
[PRIOR-ART.md](PRIOR-ART.md); cited grounding:
`research/Staleness-Budgeted-Admission-Control/`).

- **Problem statement — AI metareasoning.** Boddy & Dean (1988–1994) framed
  time-dependent planning: the world changes while the agent deliberates, so
  deliberation must be scheduled against that drift. Russell & Wefald (1991)
  formalized the value of computation. E1 is the budget-form linearization
  of that test.
- **E1's arithmetic is provably optimal in remote estimation.** Sun,
  Polyanskiy & Uysal ([arXiv:1701.06734](https://arxiv.org/abs/1701.06734))
  show that when sampling a Wiener process through a random-delay channel,
  the error-minimizing policy is a threshold on how much the process drifts
  during the delay — structurally the same object as
  `velocity × time-to-land ≤ budget`. Ornee & Sun extend this to
  mean-reverting processes; later work learns the threshold online under
  unknown delay ([arXiv:2308.15401](https://arxiv.org/abs/2308.15401)).
  Self-triggered control (Heemels, Johansson & Tabuada, CDC 2012) is the
  same ex-ante stance in control theory.
- **The drift unit is the Age-of-Incorrect-Information stance.** AoII
  ([arXiv:1907.06604](https://arxiv.org/abs/1907.06604)) accrues age only
  while the monitor's picture is actually wrong; "expected salient events
  during flight" is that content-aware refinement, which is what lets one
  unit serve as both admission tolerance (E1) and landing measurement (E5).
- **E2 is standard LLM-serving practice** — TTFT + tokens × TPOT with
  per-hardware calibration — reshaped into host-independent points times
  per-tier seconds.
- **E3 descends from imprecise computation** (Liu et al.: mandatory +
  optional parts), stale-while-revalidate
  ([RFC 5861](https://www.rfc-editor.org/rfc/rfc5861)), and model
  cascades/routing (FrugalGPT, RouteLLM). The *declared per-class*
  discipline is engineering hygiene applied to known mechanisms.
- **E4 is a Lyapunov virtual queue.** In Neely's drift-plus-penalty
  framework ([arXiv:1008.3519](https://arxiv.org/abs/1008.3519)), each
  constraint gets a queue that grows while violated and inflates its
  priority until served, with proven stability bounds. OS scheduler aging is
  the same invariant with simpler arithmetic.
- **E5 is optimistic concurrency control.** Kung & Robinson (1981): work in
  private, validate against what committed meanwhile. TL2's global version
  clock (Dice, Shalev & Shavit 2006) is the generation-counter discipline;
  the term *lease* is Gray & Cheriton (1989). Robotics re-checks plan
  validity after planner latency in the same shape.

## What is claimed, and what is not

**Not claimed:** novelty of any individual element, or of the threshold
arithmetic — the opposite: E1 has optimality proofs behind it, and citing
them is the point.

**Claimed:** the *packaging* — five specific elements composed as one
deterministic, auditable admission protocol for slow oracles, under one
name, with verdict-carries-arithmetic as the audit discipline and a shared
drift unit closing the loop between admission prediction (E1) and landing
measurement (E5). At the time of the survey recorded in
[PRIOR-ART.md](PRIOR-ART.md) (2026-07-23), no single reference was found
that packages all five elements, and no prior use of a name for the combined
protocol was found.

## Why name it now

The LLM-agent literature of 2025 is visibly converging on this decision
space from several directions at once:

- **RRARA** ([arXiv:2506.07223](https://arxiv.org/abs/2506.07223)) pairs a
  rule-based reflex with an asynchronous LLM refiner and explicitly prices
  inference delay in simulation frames — the reflex-fallback arm (E3) plus
  latency pricing, without an ex-ante drift gate.
- **Win Fast or Lose Slow** ([arXiv:2505.19481](https://arxiv.org/abs/2505.19481),
  NeurIPS 2025) is the first systematic study of the latency–quality
  tradeoff for LLM agents, showing that trading quality for latency can
  raise downstream reward — the economic premise of tiered routing (E1+E3).
- **Rational Metareasoning for LLMs**
  ([arXiv:2410.05563](https://arxiv.org/abs/2410.05563)) applies Russell &
  Wefald's value-of-computation directly to LLM inference.
- **Plan caching and speculative planning**
  ([arXiv:2506.14852](https://arxiv.org/abs/2506.14852),
  [arXiv:2509.01920](https://arxiv.org/abs/2509.01920)) institutionalize
  "act on possibly-stale guidance, refresh off the critical path" — flight
  and supersession (E5) without the priced admission (E1).

More recent preprints go further into the same territory — at least one 2026
preprint occupies the exact query/cache/drop decision space over a slow LLM
planner with a *learned* gate (ASSCG,
[arXiv:2606.25509](https://arxiv.org/abs/2606.25509); preprint of unverified
review status, cited with that caveat). Each of these systems holds a piece
of the protocol; none packages it, and none is deterministic-auditable. The
precedent this publication follows is the *circuit breaker*: a combination
of known mechanisms that became reusable the day it got a name. Landfall is
offered as that name for the ask-gate → flight → land-validate protocol.

## Conformance

An implementation is the pattern — not a knockoff — iff it upholds the
invariants (stated normatively in [DESIGN.md](DESIGN.md)): the gate and
landing check are pure; every class has a legal degrade mode at
registration; every suppression and every non-landed outcome is recorded;
the latency estimator learns from successes only; predicted and actual drift
share a unit; and an uncalibrated system fails toward its floor. A learned
or runtime-tuned gate may well be useful — the converging literature above
is building exactly that — but it is a different pattern, because it
surrenders the property the name stands for: any logged verdict is
reproducible from its own recorded fields.

## How to cite

> Evan Stern. *Landfall: staleness-budgeted admission control for slow
> oracles.* Defensive publication, 2026-07-23.
> https://github.com/evanstern/landfall/blob/main/PATTERN.md

```bibtex
@misc{stern2026landfall,
  author       = {Stern, Evan},
  title        = {Landfall: staleness-budgeted admission control for slow oracles},
  year         = {2026},
  month        = {7},
  howpublished = {\url{https://github.com/evanstern/landfall/blob/main/PATTERN.md}},
  note         = {Defensive publication. Named engineering pattern:
                  ask-gate → flight → land-validate.}
}
```

The git history of this repository fixes the publication date; cite a commit
hash for a byte-exact version.
