---
title: Debt and Anti-Starvation Mechanisms
aliases: [E4 prior art]
tags: [e4, lyapunov, virtual-queues, starvation]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Debt and Anti-Starvation Mechanisms

Landfall's E4 — suppression debt that inflates the effective budget so a repeatedly-suppressed
class eventually gets through — matches a well-studied mechanism shape in stochastic network
optimization.

## Lyapunov virtual queues (Neely)

In Neely's drift-plus-penalty framework, every time-average constraint is paired with a
**virtual queue**: the queue length grows while the constraint is being violated and enters the
controller's objective, **inflating that constraint's priority until it is served**; keeping the
virtual queue mean-rate stable is equivalent to satisfying the constraint. A tunable parameter
V trades objective value against backlog with explicit bounds ([[_grounding]] §E4;
[arXiv:1008.3519](https://arxiv.org/pdf/1008.3519)).

This is the same debt-counter shape as suppression debt — accumulate while denied, spend by
being admitted — proven in a different setting (network scheduling rather than oracle
admission). Notably, the same virtual-queue machinery appears inside the E1 literature itself:
Ornee & Sun's online threshold-learning uses it ([[_grounding]] §E1).

## OS aging

Classical schedulers prevent starvation by **aging**: a waiting task's priority rises with time
queued until it must be scheduled. Same invariant (bounded starvation), simpler arithmetic
([[_grounding]] §E4).

## Grounding

- [[_grounding]] — §E4 (Lyapunov virtual queues, aging)
