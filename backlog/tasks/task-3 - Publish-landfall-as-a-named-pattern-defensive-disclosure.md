---
id: TASK-3
title: Publish landfall as a named pattern (defensive disclosure)
status: In Progress
assignee: []
created_date: '2026-07-23 04:28'
updated_date: '2026-07-23 04:57'
labels: []
dependencies:
  - TASK-1
  - TASK-2
references:
  - >-
    research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md
priority: high
ordinal: 3000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Derived from Analysis-Using-The-Prior-Art. The corpus shows every element individually anticipated but the packaged ask-gate -> flight -> land-validate protocol unnamed, with 2025-2026 LLM-agent literature converging (ASSCG, RRARA, plan caching). Publish a citable write-up (blog post, arXiv note, or README grown into one) naming the pattern, before someone else formalizes it. Deliberately forgoes broad-patent path per the survey's 103/101 read. Lean urgency argument on RRARA + Win Fast or Lose Slow (NeurIPS 2025) + older literature, not the fragile 2026 arXiv entries.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 A public, citable write-up exists naming the pattern and its five elements
- [ ] #2 Write-up cites the prior-art lineage and states the deterministic-auditable differentiator
- [ ] #3 Publish/no-patent decision is recorded in the repo (PRIOR-ART.md or DESIGN.md)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Write PATTERN.md: self-contained, citable defensive disclosure naming the pattern (Landfall — staleness-budgeted admission control) and its five elements (E1 ex-ante drift gate, E2 shape×host latency factoring, E3 declared degrade taxonomy, E4 suppression debt, E5 lease + landing validation). Cites the prior-art lineage (Boddy & Dean, Russell & Wefald, Sun/Polyanskiy/Uysal, Ornee & Sun, AoII, Neely, OCC/TL2/leases) and states the deterministic-auditable differentiator vs learned gates. Urgency/convergence argument leans on RRARA + Win Fast or Lose Slow (NeurIPS 2025) + plan caching + older literature per task guidance; fragile 2026 arXiv entries mentioned only with caveats. Includes explicit defensive-publication statement and a how-to-cite block (public GitHub repo = public + citable).
2. Record the publish/no-patent decision as a dated Decision section in PRIOR-ART.md.
3. Link PATTERN.md from README.
4. go test -race ./... (doctrine), one branch task-3-defensive-disclosure, one PR.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Write-up complete and pushed on branch task-3-defensive-disclosure: PATTERN.md (names the pattern Landfall; five elements E1–E5; ask-gate → flight → land-validate lifecycle; per-element prior-art lineage; deterministic-auditable differentiator; explicit claimed/not-claimed section; plain + BibTeX citation block). PRIOR-ART.md gains a dated 'Decision: publish, don't patent (2026-07-23)' section; README links both. Urgency argument leans on RRARA + Win Fast or Lose Slow (NeurIPS 2025) + plan caching + older literature per task guidance; ASSCG cited only with an unverified-review caveat; Faramesh not leaned on. go test -race passes. PR #2 open: https://github.com/evanstern/landfall/pull/2 — merge blocked by session permissions, awaiting human merge. ACs to be ticked once the write-up is on main, since the citable URL points at main.
<!-- SECTION:NOTES:END -->
