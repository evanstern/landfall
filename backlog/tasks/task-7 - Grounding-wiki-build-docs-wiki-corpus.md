---
id: TASK-7
title: 'Grounding wiki: build docs/wiki corpus'
status: In Progress
assignee: []
created_date: '2026-07-23 05:54'
updated_date: '2026-07-23 06:02'
labels: []
dependencies: []
priority: medium
ordinal: 7000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Build the code-grounded corpus (grounding-wiki plugin) at docs/wiki/: per-concept notes pinned to the verified commit, interlinked and indexed, passing the freshness gate. This re-grounds the codebase per the PDLC loop and becomes the primary input for future course updates.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 docs/wiki/ exists with INDEX.md and per-concept notes pinned to the current commit
- [x] #2 Notes are interlinked and cover the core concepts (gate/route, class/degrade, lease/landing, estimator, debt/tiers)
- [x] #3 grounding-wiki freshness gate passes
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Verify preconditions (git repo, no existing docs/wiki) and take the pin from HEAD.
2. Survey codebase (all four source files, tests, README/DESIGN/PATTERN/PRIOR-ART) and draft the note list in INDEX.md.
3. Write overview + 18 concept/component notes per templates/note.md, each with sources and verified_against pin, interlinked.
4. Run the grounding-wiki freshness gate until green.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Corpus built at docs/wiki/: INDEX.md + 19 notes pinned to a9b5701, grouped by gate/classes/starvation/landing/estimation/doctrine. Freshness gate: OK, 19 notes fresh against pinned sources.
<!-- SECTION:NOTES:END -->
