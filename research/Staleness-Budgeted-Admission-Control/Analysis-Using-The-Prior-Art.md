---
title: Analysis — How Should Landfall Use These Prior-Art Findings
aliases: [prior-art action plan]
tags: [analysis, prior-art, strategy]
type: analysis
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Analysis — How Should Landfall Use These Prior-Art Findings

Given the landfall codebase and its function, how should the project best use this branch's
prior-art findings? Written to be decomposed into actionable tasks.

## Verdict

Use the grounding in three ways, in this priority order:

1. **Position, don't re-derive: cite the theory in the project's own documents.** The gate's
   arithmetic has optimality proofs behind it; landfall should claim that lineage explicitly
   instead of presenting the formula as folk engineering.
2. **Publish the pattern as a named engineering pattern, soon — and drop any broad-patent
   ambition.** The corpus shows every element is individually anticipated (§103 exposure) while
   the packaged protocol is not yet named by anyone — and the 2025–2026 agent literature is
   visibly converging on it. The window in which landfall gets to name the pattern is closing;
   defensive disclosure both fits the evidence and is time-sensitive.
3. **Mine the corpus for design-hardening work on landfall's two documented open questions**
   (debt under `RouteTiered`; `Lease` carrying predicted drift) — the literature supplies
   ready-made machinery for both — while explicitly *rejecting* the learned-gate direction the
   converging papers take, because determinism-with-audit is landfall's one genuinely
   unoccupied differentiator.

## Reasoning

**Why citation is the cheapest high-value move.** [[Drift-Thresholded-Sampling-Theory]] shows
E1's gate is mathematically the same object as Sun/Polyanskiy/Uysal's provably optimal
threshold-on-drift-during-delay, with Ornee & Sun extending it to mean-reverting processes.
DESIGN.md currently derives the formula from promptworld experience alone. A short "prior art &
theory" section (DESIGN.md and/or README) citing Sun/Ornee for the gate, Boddy & Dean and
Russell & Wefald for the problem statement, AoII for the drift unit, Neely for debt, and
OCC/TL2/leases for the lease converts "we invented this" into "this is the engineering
packaging of proven results" — a strictly stronger position, for an afternoon of work. The
drift unit deserves a specific callout: "expected salient events during flight" is
content-aware freshness, i.e. the Age of *Incorrect* Information stance rather than plain AoI
([[Drift-Thresholded-Sampling-Theory]]), and naming that kinship makes the budget's semantics
legible to anyone from that literature.

**Why publication beats patent, and why now.** [[Patent-Landscape]] records the survey's read:
§102 likely survivable but §103 high-risk and §101 nontrivial — a weak, easily-designed-around
patent at best. Meanwhile [[Converging-LLM-Agent-Literature]] shows ASSCG already occupying the
exact Query/Cache/Drop decision space with a *learned* gate, RRARA pricing latency in
world-drift units, and plan-caching systems institutionalizing act-on-stale-refresh-async. No
one has yet named the unified ask-gate → flight → land-validate protocol; the seed survey
judges independent formalization "imminent." When the write-up exists, the pattern's name is
whoever published first — the circuit-breaker precedent the survey itself invokes. This makes a
public, citable write-up (blog post, arXiv note, or the README grown into one) the single most
consequential task derivable from this branch, and one with a real deadline.

**Why the differentiator dictates what *not* to build.** Across
[[Converging-LLM-Agent-Literature]], the field's gates are learned (ASSCG: SFT + compute-aware
RL). Landfall's non-negotiables — pure `Route`, no model in the enforcement path, verdicts
reproducible from their own fields — are exactly what none of the converging systems offer, and
[[Patent-Landscape]] plus the grounding's null result for "verdict-carries-arithmetic" prior
art suggest that discipline is the unoccupied ground. Consequence: any task that would soften
determinism (learned budgets, runtime tuning, adaptive thresholds inside the gate) doesn't just
violate repo doctrine, it spends the only differentiator. The corpus turns an internal rule
into a market position.

**Design-hardening tasks the literature directly funds:**

- *Debt semantics under `RouteTiered`* (DESIGN.md open question): Neely's virtual-queue
  formalism ([[Debt-and-Anti-Starvation-Mechanisms]]) is precisely a per-constraint debt
  counter with proven bounded-violation guarantees. Framing each class's suppression debt as a
  virtual queue answers the cross-tier repayment question in a principled way (each constraint
  keeps its own queue → a landed local call should *not* repay the cloud class's debt) and
  gives the `DebtFactor` arithmetic a citable stability story.
- *`Lease` carrying predicted drift* (DESIGN.md open question): invariant 5 says
  predicted-vs-actual drift is the calibration report; OCC's read-set-carried-into-validation
  shape ([[Validation-and-Lease-Mechanisms]]) supports recording the admission-time prediction
  in the lease so `Stale` outcomes log the comparison without a join.
- *Calibration loop*: the Stale-outcome feedback maps onto signal-aware sampling and online
  threshold learning ([[Drift-Thresholded-Sampling-Theory]]) — useful as documentation framing
  now; a host-side helper later, never in-gate.
- *External validation*: Win Fast or Lose Slow's latency-sensitive benchmarks
  ([[Converging-LLM-Agent-Literature]]) are a ready-made arena to demonstrate `RouteTiered`
  against learned baselines, if the project ever wants empirical evidence rather than argument.

**Housekeeping the corpus demands:** before any public write-up cites the survey's related-work
list, resolve the Faramesh discrepancy — the abstract does not support "oracle latency budgets
+ audited fallback" ([[Converging-LLM-Agent-Literature]]) — and re-verify the patent list
beyond US 12,452,345. Publishing an incorrect characterization in a positioning document would
undercut the credibility the document exists to build.

## Tensions & tradeoffs

- **Publication burns patent optionality.** Disclosure starts statutory clocks and ends any
  broad-claim path. The corpus says that path was weak anyway ([[Patent-Landscape]]), but that
  judgment rests on a keyword search with an explicit 18-month blind spot; a real clearance
  search could in principle change the calculus. The verdict accepts this risk knowingly.
- **Citing optimality invites an unflattering comparison.** Sun/Ornee prove *optimal
  thresholds*; landfall's budgets are reviewed constants, deliberately not optimized
  ([[Drift-Thresholded-Sampling-Theory]] vs. DESIGN.md's "budgets are doctrine"). The write-up
  must own this honestly — doctrine constants trade regret for auditability — or a
  theory-literate reader will spot the gap.
- **The 2026 references are fragile.** ASSCG and Faramesh are recent arXiv entries of unknown
  peer-review status ([[Brief-and-Assumptions]]); leaning the "convergence is imminent" urgency
  argument on them carries citation risk. RRARA, Win Fast or Lose Slow (NeurIPS 2025), and the
  older literature carry the argument adequately on their own.
- **Task appetite.** The full derivable list is larger than the project's leaf-package scope;
  the benchmark task in particular drags in hosts, transports, and baselines landfall
  deliberately excludes. It's flagged as optional, not core.

## Confidence & open questions

**High confidence** in the cite-don't-re-derive and publish-don't-patent recommendations — both
follow from multiply-sourced, mutually consistent findings, and the second merely operationalizes
the seed survey's own conclusion with added urgency from the convergence evidence.
**Moderate confidence** in the virtual-queue answer to the `RouteTiered` debt question — the
mechanism mapping is direct, but nothing in the corpus proves the per-class-queue choice is
right for this setting; it's principled, not derived. **What would change my mind:** a
professional patent clearance search finding blocking or highly-valuable claims (flips the
publish/patent balance); Faramesh's full text actually containing latency-budgeted admission
(would mean the packaging is already partially formalized and the naming window is tighter than
assumed — publication becomes *more* urgent, the novelty claim narrower).

Unresolved and worth carrying into tasks: the Faramesh full-text check, patent-list
re-verification, and whether a "prior art" section belongs in DESIGN.md, README, or a new
standalone document.

## Basis

- [[_grounding]] — all cited facts
- [[Drift-Thresholded-Sampling-Theory]] — E1 optimality, AoII kinship, metareasoning lineage
- [[Patent-Landscape]] — §102/§103/§101 exposure, publish-vs-patent evidence
- [[Converging-LLM-Agent-Literature]] — convergence urgency, learned-gate contrast, Faramesh discrepancy
- [[Debt-and-Anti-Starvation-Mechanisms]] — virtual-queue machinery for the debt question
- [[Validation-and-Lease-Mechanisms]] — lease/validation shape for the predicted-drift question
- [[LLM-Latency-Estimation]], [[Degradation-and-Fallback-Mechanisms]] — E2/E3 are established practice; no positioning weight placed on them
- [[Brief-and-Assumptions]] — verification caveats
