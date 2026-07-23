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

What the seed survey found in the patent literature; every entry has since been re-verified
against its actual claims (2026-07-23, TASK-2). Research context, not legal advice; no
professional clearance search has been done, and applications under 18 months old are
invisible.

## Nearest grants and applications

- **US 12,452,345** — latency-budgeted routing of AI inference requests to an external model
  across a distributed cloud network: a time-based latency budget/ceiling qualifies which data
  centers may serve the request. Independently surfaced in this pass
  ([USPTO](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/12452345);
  [[_grounding]] §patents).
- Seed-survey list re-verified against actual claims 2026-07-23 ([[_grounding]] §patents):
  - **Confirmed:** US 11,717,748 B2 (Valve — latency compensation via ML-predicted user input,
    games); US 9,679,003 (IBM — rendezvous OCC, validation interleaved with read/compute/write);
    US 8,396,831 (Microsoft — optimistic serializable snapshot isolation, read-set validation +
    phantom re-scan).
  - **Weaker than surveyed:** US 11,257,002 (Amazon — accuracy-based model deployment/routing;
    latency only a monitored metric, not a selection basis); US 8,972,306 (Raytheon — claim 1 is
    fuzzy-cognitive-map sensor selection; value-of-information enters only in dependent
    claim 7); EP 1,813,065 B1 (NXP — event-notification messaging between network nodes:
    event-triggered *networking*, not event-triggered control).
  - **Not supported:** US 11,461,300 (SAP — claim 1 is consistent-hash model-*server* selection
    for cache locality; nothing accuracy- or latency-based).
  - **Unconfirmed at claim level:** US 2022/0004929 A1 (Google — spec applies context-feature
    expiration to stored training examples; claims not retrievable, application pending as
    fetched — "stale-context expiration before inference" unverified).
- The seed survey states it found nothing claiming ex-ante world-drift-priced admission of an
  LLM call combined with post-hoc generation validation ([[_grounding]] §patents); the
  re-verification weakened several nearest neighbors, leaving that null result intact.

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
