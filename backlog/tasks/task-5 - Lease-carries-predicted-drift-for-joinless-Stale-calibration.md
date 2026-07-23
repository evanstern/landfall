---
id: TASK-5
title: Lease carries predicted drift for joinless Stale calibration
status: Done
assignee: []
created_date: '2026-07-23 04:28'
updated_date: '2026-07-23 05:17'
labels: []
dependencies: []
references:
  - >-
    research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md
priority: medium
ordinal: 5000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Derived from Analysis-Using-The-Prior-Art; resolves the DESIGN.md open question 'should Lease carry predicted drift?'. Invariant 5 makes predicted-vs-actual drift the calibration report; recording the admission-time prediction in the lease lets Stale outcomes log the comparison without a join (OCC read-set-into-validation shape). Extend TestVerdictReproducible for any new recorded fields; Lease.Land stays pure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Lease records predicted drift at checkout; Land outcomes expose predicted vs actual without a join
- [x] #2 TestVerdictReproducible extended to the new field(s)
- [x] #3 go test -race ./... passes
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add PredictedDrift to Lease; add pure Checkout(Verdict, gen) freezing Class, EffectiveBudget, PredictedDrift at launch.
2. Land returns a Landing record (Outcome + Gen/CurrentGen/Budget/PredictedDrift/ActualDrift) — self-contained, reproducible, pairs predicted vs actual with no join (OCC read-set-into-validation shape).
3. Extend TestVerdictReproducible through Checkout; add TestLandingReproducible; update TestLeaseLand.
4. Update DESIGN.md lifecycle + resolve the open question.
5. go test -race ./...; one branch, one PR (task-5-lease-predicted-drift).
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented on branch task-5-lease-predicted-drift; merged as PR #4 (https://github.com/evanstern/landfall/pull/4). Lease gains PredictedDrift via pure Checkout(verdict, gen); Land returns a self-contained Landing record pairing predicted vs actual drift. TestLandingReproducible added; TestVerdictReproducible extended through Checkout. DESIGN.md open question resolved; README example updated. go test -race ./... passes on main.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Lease now carries the gate's predicted drift, resolving the DESIGN.md open question. Checkout(verdict, gen) freezes class, effective budget, and predicted drift into the lease (the OCC read set); Land returns a Landing record pairing predicted vs actual drift, so Stale calibration needs no join back to the verdict log. Landing is reproducible from its own fields (TestLandingReproducible); TestVerdictReproducible extended through Checkout; Land stays pure. Merged in PR #4.
<!-- SECTION:FINAL_SUMMARY:END -->
