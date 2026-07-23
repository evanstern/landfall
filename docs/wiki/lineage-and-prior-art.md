---
name: lineage-and-prior-art
description: Where each element comes from — promptworld's cognition horizon, the sampling-theory optimality proofs, AoII, Neely queues, OCC — and the defensive disclosure.
kind: concept
sources:
  - DESIGN.md
  - PATTERN.md
  - PRIOR-ART.md
  - README.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Lineage and prior art

landfall claims packaging, not theory. It was generalized from promptworld's
cognition horizon (`internal/cognition`) — a deterministic gate deciding
whether an LLM call is allowed at the current sim speed — and every element
has an explicit lineage that DESIGN.md cites rather than presenting the
arithmetic as folk wisdom.

## How it works

The lineages DESIGN.md claims per element:

- **Problem statement** — AI metareasoning: Boddy & Dean's time-dependent
  planning; Russell & Wefald's value of computation. The gate is the
  budget-form linearization of "compute only if the utility gain exceeds the
  time cost."
- **Gate formula** — Sun, Polyanskiy & Uysal (arXiv:1701.06734) prove the
  optimal sampling policy for a drifting signal over a random-delay channel
  is a threshold on drift-during-service — structurally
  `velocity × time-to-land ≤ budget` ([[route-gate]]). Ornee & Sun extend to
  Ornstein–Uhlenbeck; self-triggered control is the same ex-ante stance.
- **Drift unit** — the Age-of-Incorrect-Information stance: age accrues only
  while the picture is wrong ([[drift-unit]]).
- **Suppression debt** — Neely's drift-plus-penalty virtual queues, one per
  constraint, with `Budget × (1 + debt × DebtFactor)` as the inflation term
  ([[suppression-debt]], [[per-tier-debt]]).
- **Lease** — optimistic concurrency control (Kung & Robinson), TL2's global
  version clock as the generation-counter discipline, Gray & Cheriton for the
  term "lease" ([[lease-lifecycle]]).

The honest gap DESIGN.md records: the theory *derives* optimal thresholds
from process statistics; landfall's budgets are doctrine constants, carrying
regret but buying auditability — reviewed, diffable, human-meaningful. The
calibration log monitors the tradeoff ([[calibration-doctrine]]).

PATTERN.md publishes the pattern as a **defensive disclosure** (first
published 2026-07-23): five elements (ex-ante gate, points × sec/pt latency
factoring, degrade taxonomy, suppression debt, generation lease +
verdict-carries-arithmetic), named and placed in the public record with no
patent sought. PRIOR-ART.md records the publish-don't-patent decision and the
survey behind it: every element individually formalized somewhere, no single
reference packaging all five as one auditable protocol.

## Connections

What landfall claims as its own is the packaging in
[[invariants-and-purity]]: the pure gate with no model in the enforcement
path, the [[verdict]] carrying its arithmetic, and degradation forced to
exist at registration ([[degrade-modes]]).

## Operational notes

Full citations live in the research vault at
`research/Staleness-Budgeted-Admission-Control/`. The prior-art survey is a
keyword search with an 18-month application blind spot, not a professional
clearance search — PRIOR-ART.md says so explicitly.
