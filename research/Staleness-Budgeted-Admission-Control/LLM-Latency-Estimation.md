---
title: LLM Latency Estimation
aliases: [E2 prior art]
tags: [e2, llm-serving, latency, calibration]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# LLM Latency Estimation

Landfall's E2 — latency modeled as shape points × calibrated seconds-per-point per tier — has a
direct counterpart in standard LLM serving practice.

## The standard decomposition

LLM request latency is conventionally split as **total = TTFT + TPOT × output tokens**: time to
first token (queueing + prefill) plus time per output token (steady-state decode). This is the
industry-standard linear "shape × rate" factoring ([[_grounding]] §E2;
[BentoML metrics guide](https://bentoml.com/llm/llm-inference-basics/llm-inference-metrics)).

## Per-deployment calibration in production

- **llm-d** trains a lightweight XGBoost regressor mapping request/server features (prompt
  length, prefix-cache hit rate, queue depth, KV utilization) to observed TTFT/TPOT, retrained
  continuously on a sliding window — live calibration of the rate term.
- **ELIS** trains a response-length predictor (BGE encoder) to estimate the shape term, feeding
  an iterative shortest-remaining-time-first scheduler
  ([arXiv:2505.09142](https://arxiv.org/abs/2505.09142)).
- **LENS/PRISM** pairs an SLO-aware engine scheduler with a predictive router
  ([arXiv:2509.23384](https://arxiv.org/html/2509.23384v1)).

## Relation to the pattern

The capability (predicting an LLM call's wall-clock from shape features and calibrated
per-deployment rates) is established practice; the serving literature uses it for scheduling
and SLO attainment rather than as an input to a drift-vs-budget admission gate
([[_grounding]] §E2).

## Grounding

- [[_grounding]] — §E2 (TTFT/TPOT decomposition, llm-d, ELIS, LENS)
