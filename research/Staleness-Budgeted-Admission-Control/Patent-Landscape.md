---
title: Patent Landscape
aliases: [patentability]
tags: [patents, prior-art, ip]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Patent Landscape

What the seed survey found in the patent literature, plus the one grant independently
re-surfaced in this research pass. Research context, not legal advice; no professional
clearance search has been done, and applications under 18 months old are invisible.

## Nearest grants and applications

- **US 12,452,345** — latency-budgeted routing of AI inference requests to an external model
  across a distributed cloud network: a time-based latency budget/ceiling qualifies which data
  centers may serve the request. Independently surfaced in this pass
  ([USPTO](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/12452345);
  [[_grounding]] §patents).
- Named by the seed survey (not re-verified here): US 11,257,002 / 11,461,300
  (accuracy/latency-based ML model selection); US 2022/0004929 A1 (on-device ML with
  stale-context expiration before inference); US 11,717,748 B2 (latency compensation via
  prediction, games); US 9,679,003 / 8,396,831 (OCC validation); US 8,972,306
  (value-of-information sensor tasking); EP 1,813,065 B1 (event-triggered communication).
- The seed survey states it found nothing claiming ex-ante world-drift-priced admission of an
  LLM call combined with post-hoc generation validation ([[_grounding]] §patents).

## The survey's patentability framing

As recorded in the seed survey (its judgment, reproduced descriptively):

- **§102 novelty** — a five-element combination claim would likely survive strict anticipation;
  no single reference discloses the full stack.
- **§103 obviousness** — high risk: each element is textbook in an adjacent field, and the
  motivation to combine is documented in the field itself (KSR standard: predictable
  combination of known elements performing known functions).
- **§101 eligibility** — nontrivial risk: the gate is arithmetic on abstract quantities
  (Alice); claims would need concrete technical framing.
- The survey's practical read is that the pattern is a weak candidate for a broad patent and
  suggests publication / defensive disclosure instead ([[_grounding]] §patentability).

## Grounding

- [[_grounding]] — §Patent landscape, §Patentability doctrine
