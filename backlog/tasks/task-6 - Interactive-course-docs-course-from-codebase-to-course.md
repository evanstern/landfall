---
id: TASK-6
title: 'Interactive course: docs/course from codebase-to-course'
status: To Do
assignee: []
created_date: '2026-07-23 05:54'
labels: []
dependencies: []
priority: medium
ordinal: 6000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Generate the interactive single-page HTML course (praxisflux codebase-to-course) teaching the landfall codebase to non-technical readers, at the PDLC standing location docs/course/. Built from raw source + DESIGN.md (no docs/wiki corpus exists yet; grounding wiki is a separate task).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Course directory docs/course/ exists with prebuilt chrome (styles.css, main.js v2 stamp) copied from plugin references
- [ ] #2 Five modules covering gate formula, actors, lease lifecycle, estimator, debt/tiers/purity — all code snippets verbatim from source
- [ ] #3 build.sh validation passes and index.html is assembled
- [ ] #4 codebase-to-course output gate passes (self-contained, mandatory interactive elements present)
<!-- AC:END -->
