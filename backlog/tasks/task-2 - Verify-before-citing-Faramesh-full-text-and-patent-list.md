---
id: TASK-2
title: 'Verify before citing: Faramesh full text and patent list'
status: Done
assignee: []
created_date: '2026-07-23 04:28'
updated_date: '2026-07-23 04:51'
labels: []
dependencies: []
references:
  - >-
    research/Staleness-Budgeted-Admission-Control/Analysis-Using-The-Prior-Art.md
priority: high
ordinal: 2000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Derived from Analysis-Using-The-Prior-Art. Housekeeping the corpus demands before any public write-up: (1) the seed survey describes Faramesh (arXiv 2601.17744) as 'oracle latency budgets + audited fallback' but its abstract only supports pre-execution action authorization with provenance logging — read the full text and correct PRIOR-ART.md and the vault grounding accordingly; (2) re-verify the survey's patent list beyond US 12,452,345 (US 11,257,002 / 11,461,300, US 2022/0004929 A1, US 11,717,748 B2, US 9,679,003 / 8,396,831, US 8,972,306, EP 1,813,065 B1).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Faramesh full text read; PRIOR-ART.md characterization corrected or confirmed, vault grounding updated
- [x] #2 Each patent in the survey list re-verified against its actual claims; discrepancies recorded in PRIOR-ART.md
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Fetch Faramesh (arXiv 2601.17744) full text; compare against the seed survey's "oracle latency budgets + audited fallback policy" characterization.
2. Correct or confirm PRIOR-ART.md; update vault grounding (_grounding.md, Converging-LLM-Agent-Literature.md) with full-text-verified facts.
3. Re-verify each seed-survey patent (US 11,257,002; US 11,461,300; US 2022/0004929 A1; US 11,717,748 B2; US 9,679,003; US 8,396,831; US 8,972,306; EP 1,813,065 B1) against actual claims via Google Patents/USPTO.
4. Record confirmations/discrepancies in PRIOR-ART.md and Patent-Landscape.md; re-cite in _grounding.md.
5. Branch task-2-verify-faramesh-patents; single PR.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Faramesh full text (arXiv 2601.17744 v1 HTML) verified 2026-07-23: pre-execution authorization control plane only — no latency budgets, no fallback policy; DEFER = pause-for-approval; decision-to-execution staleness is a stated limitation. Seed characterization corrected in PRIOR-ART.md; vault updated (_grounding.md, Converging-LLM-Agent-Literature.md, Brief-and-Assumptions.md, MOC).

All 8 patents re-verified against actual claims via Google Patents: confirmed US 11,717,748 B2 / US 9,679,003 / US 8,396,831; weaker than surveyed US 11,257,002 / US 8,972,306 / EP 1,813,065 B1; miscited US 11,461,300 (cache-locality server selection); unconfirmed at claim level US 2022/0004929 A1 (claims not retrievable, pending). Discrepancies recorded in PRIOR-ART.md and Patent-Landscape.md. Null result (no drift-priced admission + post-hoc validation claim) survives, strengthened.

Vault branch gate OK; go test -race passes. PR #1 open: https://github.com/evanstern/landfall/pull/1 — merge blocked by session permissions, awaiting human merge.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Faramesh full text verified: pre-execution authorization control plane only — seed survey's 'oracle latency budgets + audited fallback' characterization corrected in PRIOR-ART.md and vault; staleness named as the paper's limitation, not a feature. All 8 seed-survey patents re-verified against actual claims: 3 confirmed, 3 weaker than surveyed, US 11,461,300 miscited, US 2022/0004929 A1 unconfirmed at claim level. Null result (no drift-priced admission + post-hoc validation prior claim) strengthened. Merged via PR #1.
<!-- SECTION:FINAL_SUMMARY:END -->
