---
title: Brief and Assumptions
aliases: []
tags: [brief]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Brief and Assumptions

## The request (restated)

"Please read PRIOR-ART.md and find what information you can on this subject."

`PRIOR-ART.md` (landfall repo root) is a prior-art and patentability survey for the **landfall
pattern**: staleness-budgeted admission control for slow oracles. It decomposes the pattern into
five claimed elements:

- **E1** ex-ante gate: `predicted drift = velocity × predicted latency ≤ per-class budget`
- **E2** latency factored as shape points × calibrated seconds-per-point per tier
- **E3** declared degrade taxonomy (skip / reflex / template / faster-tier)
- **E4** suppression debt inflating the effective budget (anti-starvation)
- **E5** generation lease + landing validation (Landed / Superseded / Stale), verdict-carries-arithmetic

## Assumptions made

- "This subject" = the prior-art landscape around those five elements plus the converging
  LLM-agent literature and patent landscape the survey names — not the landfall Go
  implementation itself.
- The survey's citations were taken as leads to verify and expand, not as ground truth; where a
  fetched abstract disagreed with the survey's characterization, the discrepancy is recorded in
  [[_grounding]] (see Faramesh).
- The `deep-research` skill was unavailable in the research session; the documented fallback
  (parallel web-search fan-out, 16 searches + 2 arXiv abstract fetches) was used instead.

## Open questions / flagged ambiguities

- The seed survey's patent list (other than US 12,452,345) was not independently re-verified;
  a professional clearance search remains out of scope.
- Whether Faramesh's full text (beyond the abstract) supports the survey's "oracle latency
  budgets + audited fallback policy" description is unresolved.
- arXiv IDs dated 2026 (ASSCG, Faramesh, others) are recent; their peer-review status is
  unknown.

## Grounding

- [[_grounding]] — the research pass this brief scoped
