---
title: Degradation and Fallback Mechanisms
aliases: [E3 prior art]
tags: [e3, graceful-degradation, caching, model-routing]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Degradation and Fallback Mechanisms

Landfall's E3 — a declared per-class degrade taxonomy (skip / reflex / template / faster-tier)
— assembles mechanisms that each exist independently in adjacent fields.

## Imprecise computation (real-time systems)

Jane W.S. Liu's imprecise-computation model splits each time-critical task into a **mandatory
part** (must complete for an acceptable result) and an **optional part** (refines it). Under
overload the optional part is shed, so the system degrades gracefully by construction rather
than by failure ([[_grounding]] §E3). The model has been applied to scheduling deep-learning
inference as imprecise computation ([arXiv:2011.01112](https://arxiv.org/pdf/2011.01112)).

## Serve-stale policies (HTTP caching)

RFC 5861 standardizes two degrade responses to a slow or failing origin:
`stale-while-revalidate` (answer from the stale copy now, refresh off the request path) and
`stale-if-error` (serve last-known-good on origin failure)
([RFC 5861](https://httpwg.org/specs/rfc5861.html)). These correspond to reflex/template-style
fallbacks: a cheaper, staler answer chosen deliberately over waiting ([[_grounding]] §E3).

## Model cascades and routing (LLM cost/quality)

FrugalGPT chains LLMs small→large with early exit on confidence; RouteLLM learns to route
between weaker and stronger models from preference data; 2025–2026 successors (MixLLM,
cost-aware cascade serving) continue the line. These implement the "faster-tier" degrade arm as
a learned or confidence-driven choice ([[_grounding]] §E3).

## Relation to the pattern

Each mechanism is established. What the seed survey identifies as landfall-specific is the
*declared, per-class* discipline — a reviewed taxonomy naming which degrade arm each request
class takes — which it characterizes as engineering hygiene layered on known mechanisms rather
than a new mechanism ([[_grounding]] §E3).

## Grounding

- [[_grounding]] — §E3 (imprecise computation, RFC 5861, cascades/routing)
