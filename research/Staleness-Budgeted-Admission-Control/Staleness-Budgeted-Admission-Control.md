---
title: Staleness-Budgeted Admission Control
aliases: [landfall pattern, cognition horizon prior art]
tags: [prior-art, admission-control, slow-oracles, llm]
type: moc
created: 2026-07-23
updated: 2026-07-23
related: []
---

# Staleness-Budgeted Admission Control

Prior art and surrounding literature for the landfall pattern: admitting a slow oracle (e.g. an
LLM) call only if the world's predicted drift during the call's latency fits a per-class
staleness budget, with declared degrades, suppression debt, and post-hoc landing validation.
Seeded from the repo's `PRIOR-ART.md` survey (2026-07-23) and expanded/verified by web research.

## Scope

**In:** the five claimed elements (E1 drift-vs-budget gate, E2 shape×rate latency estimation,
E3 degrade taxonomy, E4 suppression debt, E5 lease + landing validation), their closest
formalizations in adjacent fields, the converging 2025–2026 LLM-agent literature, and the
patent landscape as surveyed. **Out:** the landfall Go implementation itself, legal advice, and
any verdict on strategy (publication vs. patent) — that belongs in an analysis note.

## What is known

- E1's arithmetic is not merely known but proven optimal in remote-estimation theory, and the
  problem statement dates to 1988 metareasoning — [[Drift-Thresholded-Sampling-Theory]].
- E2's shape×rate latency factoring is standard LLM serving practice with production
  calibration systems — [[LLM-Latency-Estimation]].
- E3's degrade arms each exist as mature mechanisms (imprecise computation, RFC 5861,
  cascades/routing); the per-class declaration discipline is the survey's claimed residual —
  [[Degradation-and-Fallback-Mechanisms]].
- E4's debt counter matches Neely's Lyapunov virtual queues and OS aging —
  [[Debt-and-Anti-Starvation-Mechanisms]].
- E5 is a direct mechanism transfer from OCC / TL2 / leases —
  [[Validation-and-Lease-Mechanisms]].
- The 2025–2026 LLM-agent literature (RRARA, ASSCG, Win Fast or Lose Slow, plan caching) is
  independently converging on the same decision space, mostly with learned rather than
  deterministic-auditable gates — [[Converging-LLM-Agent-Literature]].
- The nearest patents cover fragments (latency-budgeted routing, model selection, stale-context
  expiration); the seed survey found no grant combining ex-ante drift-priced admission with
  post-hoc generation validation — [[Patent-Landscape]].
- The request's constraints and the assumptions made are in [[Brief-and-Assumptions]].

## Notes

- [[Brief-and-Assumptions]] — the request restated, assumptions, flagged ambiguities
- [[Drift-Thresholded-Sampling-Theory]] — E1: remote estimation optimality, AoI/AoII, self-triggered control, metareasoning
- [[LLM-Latency-Estimation]] — E2: TTFT/TPOT factoring and calibrated predictors in LLM serving
- [[Degradation-and-Fallback-Mechanisms]] — E3: imprecise computation, serve-stale, cascades/routing
- [[Debt-and-Anti-Starvation-Mechanisms]] — E4: Lyapunov virtual queues, aging
- [[Validation-and-Lease-Mechanisms]] — E5: OCC, TL2 version clocks, leases
- [[Converging-LLM-Agent-Literature]] — 2025–2026 systems circling the same decision space
- [[Patent-Landscape]] — nearest patents and the survey's §102/§103/§101 framing

## Analyses

- [[Analysis-Using-The-Prior-Art]] — given the landfall codebase, how should the project use these findings (task-derivation oriented)

## Open questions

- ~~Does Faramesh's full text support the seed survey's "oracle latency budgets + audited
  fallback" characterization?~~ Resolved 2026-07-23: it does not — Faramesh is audit-plane
  only ([[Converging-LLM-Agent-Literature]]).
- ~~The seed survey's patent list beyond US 12,452,345 is unverified here.~~ Re-verified
  against actual claims 2026-07-23 ([[Patent-Landscape]]); a clearance search remains out of
  scope.
- Peer-review status of the 2026 arXiv entries (ASSCG, Faramesh) is unknown.
- Is there published work on *reproducible-from-recorded-fields* admission verdicts (the
  verdict-carries-arithmetic discipline) in any adjacent field? Nothing surfaced in this pass.

## Grounding

- [[_grounding]] — the research pass this branch is built on
- [PRIOR-ART.md](../../PRIOR-ART.md) — seed survey at the landfall repo root
