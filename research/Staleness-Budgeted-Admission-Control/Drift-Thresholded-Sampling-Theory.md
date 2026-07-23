---
title: Drift-Thresholded Sampling Theory
aliases: [E1 prior art]
tags: [e1, sampling, age-of-information, control-theory]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Drift-Thresholded Sampling Theory

The mathematical core of landfall's E1 gate — admit a slow observation only if the world's
predicted drift during the observation's latency stays within a budget — exists in several
mature literatures, with optimality proofs in one of them.

## Remote estimation: the threshold is proven optimal

Sun, Polyanskiy & Uysal (2017) study sampling a Wiener process sent to a remote estimator over
a random-delay channel. The optimal sampling policy minimizing mean-square estimation error is
**a threshold on how much the process varies during the random service time** — structurally
the same object as a `velocity × latency ≤ budget` gate. When sampling ignores the signal, the
problem collapses to Age of Information minimization ([[_grounding]] §E1;
[arXiv:1701.06734](https://arxiv.org/abs/1701.06734)).

Ornee & Sun extend the result to Ornstein–Uhlenbeck processes (threshold on instantaneous
estimation error, stable and unstable cases), and later work learns the threshold online when
delay statistics are unknown ([arXiv:2308.15401](https://arxiv.org/abs/2308.15401)).

## Freshness metrics: AoI and its content-aware refinements

Age of Information measures the time lag of a monitor's knowledge; it resets on delivery and
grows linearly otherwise. **Age of Incorrect Information** (AoII) refines it so age accrues only
while the monitor's picture is *wrong* — "fresh informative updates" — merging freshness with
accuracy ([arXiv:1907.06604](https://arxiv.org/abs/1907.06604)). Query-aware variants
("perceived age") measure freshness at query time rather than continuously. DESIGN.md's "budget
as expected salient events" corresponds to this content-aware family ([[_grounding]] §E1).

## Control theory: when to sample, decided ahead of time

Event-triggered control samples when the state deviates past a threshold (reactive);
**self-triggered control computes the next sampling instant ahead of time from a model**
(proactive) — the ex-ante stance of E1. Canonical reference: Heemels, Johansson & Tabuada,
IEEE CDC 2012 ([[_grounding]] §E1).

## AI metareasoning: the problem statement, decades earlier

Boddy & Dean (1988–1994) framed **time-dependent planning**: the world changes while the agent
deliberates, so deliberation itself must be scheduled against that drift. Russell & Wefald
(1991) formalized the **value of computation**: a computation is worth running only if its
expected utility gain exceeds its time cost. Both anticipate the E1 tradeoff as a general
problem, without the specific drift-arithmetic gate ([[_grounding]] §E1).

## Grounding

- [[_grounding]] — §E1 (remote estimation, AoI/AoII, self-triggered control, metareasoning)
