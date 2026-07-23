---
id: TASK-1
title: 'Cite prior art: add theory & positioning section to DESIGN.md'
status: Done
assignee: []
created_date: '2026-07-23 04:28'
updated_date: '2026-07-23 04:37'
labels: []
dependencies: []
references:
  - >-
    research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md
priority: medium
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Derived from research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md. Position the pattern on its proven lineage instead of re-deriving: cite Sun/Polyanskiy/Uysal + Ornee & Sun for the E1 gate's optimality, Boddy & Dean and Russell & Wefald for the problem statement, Age of Incorrect Information for the salient-events drift unit, Neely's virtual queues for suppression debt, and OCC/TL2/Gray-Cheriton leases for the lease. Must honestly own the doctrine-constants vs optimal-thresholds gap (budgets trade regret for auditability).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 DESIGN.md (or a referenced standalone section) cites the E1 optimality results (Sun 1701.06734 / Ornee & Sun) where the gate formula is introduced
- [x] #2 Drift unit's kinship to Age of Incorrect Information is stated
- [x] #3 Debt is framed against Neely virtual queues; lease against OCC/TL2/leases
- [x] #4 The doctrine-constants vs derived-optimal-thresholds tradeoff is stated explicitly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Read vault notes (Drift-Thresholded-Sampling-Theory, Debt-and-Anti-Starvation, Validation-and-Lease) for exact citations
2. Branch task-1-prior-art-design
3. Add 'Prior art & theory' section to DESIGN.md: Boddy&Dean/Russell&Wefald problem statement; Sun/Polyanskiy/Uysal 1701.06734 + Ornee&Sun for E1 optimality (cited inline at the gate formula too); AoII 1907.06604 for the drift unit; Neely 1008.3519 virtual queues for debt; OCC/TL2/Gray-Cheriton for the lease
4. Explicitly own doctrine-constants vs optimal-thresholds gap
5. Cite only established theory — no ASSCG/Faramesh until TASK-2 verifies them
6. go test -race ./..., commit on branch, merge --no-ff to main (no remote configured)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Added 'Prior art & theory' section to DESIGN.md plus a one-line optimality pointer directly under the gate formula. Cites: Sun/Polyanskiy/Uysal 1701.06734 + Ornee & Sun (E1 optimality, incl. online-threshold 2308.15401), Boddy & Dean + Russell & Wefald (problem statement), AoII 1907.06604 (drift unit kinship), Neely 1008.3519 virtual queues (debt) + OS aging, Kung & Robinson OCC / TL2 / Gray & Cheriton (lease). Deliberately cites only established theory — ASSCG/Faramesh/2026 convergence papers held back pending TASK-2 verification. Doctrine-vs-optimal gap owned in a dedicated closing paragraph tied to invariant 5 and the learned-budgets out-of-scope entry.

Remote origin (git@github.com:evanstern/landfall.git) added after the local merge; pushed main (includes merge commit of task-1-prior-art-design) and the task branch to origin. Future tasks (TASK-2+) should use the real branch→PR flow.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
DESIGN.md now carries a Prior art & theory section positioning landfall as the engineering packaging of proven results: E1 gate optimality (Sun/Polyanskiy/Uysal, Ornee & Sun) cited where the formula is introduced, drift unit framed as AoII, debt as a Neely virtual queue, lease as OCC/TL2/Gray-Cheriton, and the doctrine-constants vs derived-optimal-thresholds regret/auditability tradeoff stated explicitly. Landed as branch task-1-prior-art-design, merged --no-ff to main (commit 87edf52; no remote configured, so local merge stands in for the PR). go test -race passed.
<!-- SECTION:FINAL_SUMMARY:END -->
