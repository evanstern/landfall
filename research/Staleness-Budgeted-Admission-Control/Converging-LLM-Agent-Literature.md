---
title: Converging LLM-Agent Literature
aliases: [2025-2026 convergence]
tags: [llm-agents, latency, convergence]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Converging LLM-Agent Literature

The 2025–2026 LLM-agent literature is independently circling the pattern's decision space:
when is a slow oracle call worth its latency, and what happens to its answer when it lands late.

## Closest applied systems

- **RRARA** ([arXiv:2506.07223](https://arxiv.org/abs/2506.07223)) — rule-based reflex acts
  immediately; an asynchronous LLM Reflector refines in situ; a Time Conversion Mechanism prices
  inference delay in simulation frames; adds latency-aware metrics (Respond Latency,
  Latency-to-Action Ratio). Covers the reflex-fallback arm and latency pricing, without an
  ex-ante drift gate ([[_grounding]] §converging).
- **ASSCG** ([arXiv:2606.25509](https://arxiv.org/abs/2606.25509), abstract verified
  2026-07-23) — frame-level **Query / Cache / Drop** decisions over a slow LLM planner in
  autonomous driving; the same decision space as the pattern's admit / reuse / suppress, but the
  gate is *learned* (SFT + compute-aware RL on an RWKV backbone) rather than deterministic and
  auditable ([[_grounding]] §converging).

## Latency–quality as a studied tradeoff

- **Win Fast or Lose Slow** ([arXiv:2505.19481](https://arxiv.org/abs/2505.19481), NeurIPS 2025)
  — first systematic study of the latency–quality tradeoff for LLM agents (HFTBench,
  StreetFighter); FPX adaptively picks model size/quantization; trading quality for latency can
  raise downstream reward substantially.
- **Rational Metareasoning for LLMs** ([arXiv:2410.05563](https://arxiv.org/pdf/2410.05563)) —
  applies Russell & Wefald's value-of-computation directly to LLM inference.
- **Dynamic Speculative Agent Planning** ([arXiv:2509.01920](https://arxiv.org/pdf/2509.01920))
  — speculative execution against planner latency.

## Async plan caching

**AgenticCache** ([arXiv:2604.24039](https://arxiv.org/abs/2604.24039)) reuses cached plans
while a background updater asynchronously validates/refines entries via the LLM; **Agentic Plan
Caching** ([arXiv:2506.14852](https://arxiv.org/abs/2506.14852)) stores and adapts plan
templates at test time. Both institutionalize "act on possibly-stale guidance, refresh off the
critical path" ([[_grounding]] §converging).

## Audit / authorization plane

**Faramesh** ([arXiv:2601.17744](https://arxiv.org/abs/2601.17744), full text verified
2026-07-23) — a pre-execution Action Authorization Boundary issuing PERMIT/DEFER/DENY, with
append-only provenance logs keyed by canonical action hashes for deterministic replay.
Full-text verification resolved the seed survey's discrepancy: the paper contains no latency
budgets and no fallback policy — latency appears only as a measured metric, DEFER pauses for
approval rather than triggering an alternative path, and decision-to-execution staleness is
named as a limitation, not addressed. Faramesh is adjacent on the audit/authorization plane
only ([[_grounding]] §converging).

## Naming

No search hit was found for "cognition horizon" in landfall's sense, consistent with the seed
survey ([[_grounding]] §converging).

## Grounding

- [[_grounding]] — §Converging LLM-agent literature
