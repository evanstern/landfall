---
id: TASK-4
title: 'Doctrine: debt semantics under RouteTiered via virtual-queue framing'
status: To Do
assignee: []
created_date: '2026-07-23 04:28'
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
- [ ] #1 Doctrine decided and written into DESIGN.md (open question removed)
- [ ] #2 Route/RouteTiered tests cover the chosen cross-tier debt semantics
- [ ] #3 go test -race ./... passes
<!-- AC:END -->
