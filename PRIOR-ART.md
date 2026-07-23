# Prior art & formalization survey

Deep-research sweep run 2026-07-23 (web + Google Patents) asking two
questions: has this pattern been formalized, and is it patentable? Research
context, **not legal advice**; patent keyword search is not a professional
clearance search, and applications under 18 months old are invisible.

For analysis the pattern decomposes into five claimed elements:

- **E1** ex-ante gate: `predicted drift = velocity × predicted latency ≤ per-class budget`
- **E2** latency factored as shape points × calibrated seconds-per-point per tier
- **E3** declared degrade taxonomy (skip / reflex / template / faster-tier)
- **E4** suppression debt inflating the effective budget (anti-starvation)
- **E5** generation lease + landing validation (Landed / Superseded / Stale), verdict-carries-arithmetic

## Verdict in one paragraph

Every element is individually formalized somewhere, most of it maturely. The
E1 arithmetic is not just known — it is *proven optimal*: Sun, Polyanskiy &
Uysal (2017) show that for sampling a drifting signal over a random-delay
channel, the optimal policy is a threshold on how much the process varies
during the delay — mathematically the same object as our gate
([arxiv 1701.06734](https://arxiv.org/abs/1701.06734), extended to
Ornstein–Uhlenbeck by Ornee & Sun). No single reference packages all five
elements as one auditable protocol for LLM oracles — the packaging appears
novel as a *named engineering pattern* — but the LLM-agent literature is
actively circling it (2025–2026), so independent formalization of the full
combination looks imminent.

## Per-element closest prior art

| element | closest formalization | closeness |
|---|---|---|
| E1 gate | Sun et al. 2017 remote estimation (threshold on drift-during-delay, optimal); self-triggered control (Heemels/Johansson/Tabuada); Age of Information & Query-AoI; Russell & Wefald value-of-computation; Boddy & Dean deliberation scheduling 1989/94 ("world changes while you think", fully anticipated) | direct (math) |
| E2 shape×host | standard LLM serving: TTFT + tokens × TPOT, per-hardware calibration (ELIS, LENS, learning-to-rank scheduling) | direct as capability |
| E3 degrade taxonomy | imprecise computation (Liu — mandatory + optional parts), stale-while-revalidate (RFC 5861), model cascades/routing (FrugalGPT, RouteLLM), fail-open/fail-closed policy | anticipated; the *declared per-class* discipline is engineering hygiene, not a new mechanism |
| E4 debt | Neely's Lyapunov virtual queues — a debt counter grows while a constraint is violated and inflates priority until served ([arxiv 1008.3519](https://arxiv.org/pdf/1008.3519)); OS aging | direct mechanism shape, different setting |
| E5 lease/landing | OCC (Kung & Robinson 1981), TL2 global version-clock validation (Dice/Shalev/Shavit 2006), Gray & Cheriton leases (1989); robotics plan-validity re-checks after planner latency | direct, verbatim mechanism in another domain |

## Closest overall

1. **Sun et al. 2017 / Ornee & Sun** — E1's gate proven optimal, framed as sampling control.
2. **Boddy & Dean 1989/1994** — the whole problem statement, 35 years early.
3. **RRARA** ([2506.07223](https://arxiv.org/pdf/2506.07223)) — LLM agent with rule-based reflex + async LLM refinement + latency converted into simulation frames; closest applied system.
4. **ASSCG** ([2606.25509](https://arxiv.org/pdf/2606.25509)) — learned Query/Cache/Drop gate over a slow LLM planner; same decision space, learned rather than deterministic-auditable.
5. **TL2 / OCC / Leases** — E5 verbatim.

Also converging: "Win Fast or Lose Slow" ([2505.19481](https://arxiv.org/abs/2505.19481)),
Rational Metareasoning for LLMs ([2410.05563](https://arxiv.org/abs/2410.05563)),
Dynamic Speculative Agent Planning ([2509.01920](https://arxiv.org/pdf/2509.01920)),
Faramesh ([2601.17744](https://arxiv.org/pdf/2601.17744) — oracle latency
budgets + audited fallback policy). No literature hit for "cognition horizon"
in this sense.

## Patent landscape (nearest found)

US 12,452,345 (latency-budgeted routing of AI inference to external models);
US 11,257,002 / 11,461,300 (accuracy/latency-based ML model selection);
US 2022/0004929 A1 (on-device ML with stale-context expiration before
inference); US 11,717,748 B2 (latency compensation via prediction, games);
US 9,679,003 / 8,396,831 (OCC validation); US 8,972,306 (value-of-information
sensor tasking); EP 1,813,065 B1 (event-triggered communication). **Nothing
found claiming ex-ante world-drift-priced admission of an LLM call with
post-hoc generation validation.**

## Patentability read

- **§102 novelty**: a claim requiring all five elements together would likely
  survive strict anticipation — no single reference discloses the stack.
- **§103 obviousness: high risk.** Each element is textbook in an adjacent
  field and the motivation to combine is documented *in the field itself* —
  a predictable combination of known elements performing their known
  functions (KSR).
- **§101 eligibility: nontrivial risk** — the gate is arithmetic on abstract
  quantities (Alice); claims would need concrete technical framing.
- **Practical read**: weak candidate for a broad patent; a narrow five-element
  claim might issue but would be easy to design around. The stronger play is
  **publication / defensive disclosure** — treat landfall as a named pattern
  à la circuit-breaker, and name it before someone else does.

## Consequences for the design

- We are in good company: the gate's arithmetic has optimality proofs behind
  it (Sun/Ornee). Worth citing rather than re-deriving.
- Our differentiators to keep sharp: **deterministic-auditable** (vs. the
  learned gates in ASSCG et al.), verdict reproducibility, and the unified
  ask→flight→land lifecycle. That is the part nobody else has packaged.
- DESIGN.md's "budget as expected salient events" maps onto content-aware
  AoI; the Stale-outcome calibration loop maps onto signal-aware sampling.
