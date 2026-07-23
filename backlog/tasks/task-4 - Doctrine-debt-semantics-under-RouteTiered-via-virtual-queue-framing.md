---
id: TASK-4
title: 'Doctrine: debt semantics under RouteTiered via virtual-queue framing'
status: Done
assignee: []
created_date: '2026-07-23 04:28'
updated_date: '2026-07-23 05:09'
labels: []
dependencies: []
references:
  - >-
    research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md
priority: medium
ordinal: 4000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Derived from Analysis-Using-The-Prior-Art; resolves the DESIGN.md open question 'does a landed local call repay debt owed to the suppressed cloud class?'. Neely's virtual-queue formalism suggests per-class (per-constraint) debt queues: a landed local call does NOT repay the cloud class's debt. Decide the doctrine, document it in DESIGN.md, and cover it with tests. No runtime tuning API; debt stays passed-in and recorded in the verdict.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Doctrine decided and written into DESIGN.md (open question removed)
- [x] #2 Route/RouteTiered tests cover the chosen cross-tier debt semantics
- [x] #3 go test -race ./... passes
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Adopt doctrine: per-constraint (per-tier) debt queues per Neely — a landing at tier T repays T and lower-quality tiers only; higher-quality tiers' debt survives and keeps inflating until that tier itself is served.
2. Code: Tier gains a Debt field; RouteTiered drops the class-wide debt param and gates each tier against its own debt-inflated budget; Verdict records the chosen tier's debt (purity preserved).
3. DESIGN.md: doctrine in Starvation + host obligations; remove the open question. PATTERN.md E4: one-paragraph tiered-debt refinement.
4. Tests: cross-tier isolation (one tier's debt never inflates another), debt pulls routing back up-tier, tiered-verdict reproducibility.
5. go test -race ./..., PR, merge, board sync.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Doctrine implemented on branch task-4-tiered-debt-semantics; PR #3 open (https://github.com/evanstern/landfall/pull/3). Per-tier debt queues per Neely: Tier carries Debt, RouteTiered gates each tier against its own inflated budget, landed lower tier does not repay higher tier's queue. go test -race passes.

Merge of PR #3 blocked by session permissions — awaiting user merge; task closes Done after merge.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Doctrine decided and shipped in PR #3 (merged as 0c5075d): per-constraint (per-tier) debt queues per Neely's virtual-queue discipline — a landed local call does NOT repay the suppressed cloud tier's debt; a landing repays only its own tier and lower-quality tiers, bounding quality starvation as well as consultation starvation. Tier carries its own Debt, RouteTiered gates each tier against its own debt-inflated budget, gate stays pure with the chosen tier's debt recorded in the verdict. DESIGN.md open question removed (doctrine in Starvation + host obligations), PATTERN.md E4 refined, tests cover cross-tier isolation, up-tier pull, DebtFactor 0, and tiered-verdict reproducibility. go test -race ./... passes.
<!-- SECTION:FINAL_SUMMARY:END -->
