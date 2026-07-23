---
id: TASK-5
title: Lease carries predicted drift for joinless Stale calibration
status: To Do
assignee: []
created_date: '2026-07-23 04:28'
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
- [ ] #1 Lease records predicted drift at checkout; Land outcomes expose predicted vs actual without a join
- [ ] #2 TestVerdictReproducible extended to the new field(s)
- [ ] #3 go test -race ./... passes
<!-- AC:END -->
