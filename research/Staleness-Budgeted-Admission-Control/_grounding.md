---
title: Staleness-Budgeted Admission Control — Grounding
aliases: []
tags: [grounding]
type: source
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Staleness-Budgeted Admission Control — Grounding

> Source-of-truth artifact. This is the raw, cited output of a research pass (the `deep-research`
> skill, or a direct web-search fan-out). Keep it close to verbatim — do not editorialize, prune,
> or draw conclusions here. Knowledge notes and analyses cite *into* this file.

**Research question:** What prior art, formalizations, and converging literature exist for the
five elements of the landfall pattern (staleness-budgeted admission control for slow oracles):
E1 ex-ante drift-vs-budget gate, E2 shape×host latency factoring, E3 declared degrade taxonomy,
E4 suppression debt, E5 generation lease + landing validation?
**Method:** web-search fan-out (16 parallel searches + 2 arXiv abstract fetches) · 2026-07-23
**Seed document:** `PRIOR-ART.md` at the landfall repo root (deep-research sweep run 2026-07-23).

---

## E1 — Ex-ante drift-vs-budget gate

### Remote estimation / sampling theory (optimality proofs)

- Sun, Polyanskiy & Uysal study sampling a Wiener process forwarded to a remote estimator over
  a channel modeled as a queue with random delay, minimizing mean-square estimation error under
  a sampling-rate constraint. **The optimal sampling strategy is a threshold policy, and the
  optimal threshold is determined by how much the Wiener process varies during the random
  service time** and the maximum allowed sampling rate. If sampling times are independent of the
  observed process, the problem reduces to minimizing Age of Information (AoI).
  [arXiv:1701.06734](https://arxiv.org/abs/1701.06734) (ISIT 2017),
  [arXiv:1707.02531](https://arxiv.org/abs/1707.02531) (journal version).
- Ornee & Sun extend this to the Ornstein–Uhlenbeck (Gauss-Markov) process: the optimal sampling
  policy is again **a threshold policy on instantaneous estimation error**, for both stable and
  unstable OU cases, recovering the Wiener result as a special case.
  [IEEE/ACM ToN 2021](https://ieeexplore.ieee.org/abstract/document/9437348).
- Follow-on work handles **unknown delay statistics** by learning the optimal threshold online
  via stochastic approximation and the virtual-queue method.
  [arXiv:2308.15401](https://arxiv.org/abs/2308.15401).

### Age of Information family (freshness metrics)

- AoI quantifies the information time lag at a monitor; the age process grows linearly and
  resets to 0 on update delivery.
- **Age of Incorrect Information (AoII)** (Maatouk et al.) merges freshness and accuracy:
  it "extends the notion of fresh updates to that of fresh *informative* updates" — if the
  source state has not changed, age does not accrue. Positioned as fixing shortcomings of both
  AoI and conventional error penalties. [arXiv:1907.06604](https://arxiv.org/abs/1907.06604).
- Query-based variants exist ("perceived age" in query-based update systems,
  [arXiv:2601.14075](https://arxiv.org/pdf/2601.14075)); a 2025 survey covers the whole
  "X of Information" continuum ([arXiv:2507.19657](https://arxiv.org/pdf/2507.19657)).

### Self-triggered / event-triggered control

- Heemels, Johansson & Tabuada, "An introduction to event-triggered and self-triggered control"
  (IEEE CDC 2012, pp. 3270–3285): **event-triggered control is reactive** — sample/actuate when
  the plant state deviates more than a threshold; **self-triggered control is proactive** —
  compute the next sampling instant ahead of time from a model.
  [Springer encyclopedia entry](https://link.springer.com/rwe/10.1007/978-3-030-44184-5_97),
  [author page](https://heemels.tue.nl/research/event-triggered-and-self-triggered-control).

### Metareasoning / deliberation scheduling

- Boddy & Dean: "An Analysis of Time-Dependent Planning" (AAAI-88) and "Deliberation scheduling
  for problem solving in time-constrained environments" (Artificial Intelligence 67, 1994,
  245–285). Frames planning where **the world changes while the agent deliberates**; anytime
  algorithms can be interrupted at any point returning a result whose utility is a function of
  computation time; gives algorithms for optimally scheduling deliberation.
  [IJCAI-89 paper](https://www.ijcai.org/Proceedings/89-2/Papers/021.pdf),
  [Semantic Scholar](https://www.semanticscholar.org/paper/31585092a9fe9d15e0e0c1ff030adca9a023cb2f).
- Russell & Wefald, *Do the Right Thing: Studies in Limited Rationality* (1991) and "Principles
  of metareasoning": the **value of computation (VOC)** balances a computation's expected utility
  gain against its cost (time). "Object-level computations are actions with costs (the passage
  of time) and benefits (improvements in decision quality)."
  [Principles of metareasoning](https://www.sciencedirect.com/science/article/abs/pii/000437029190015C),
  [Russell, Rationality and Intelligence](https://people.eecs.berkeley.edu/~russell/papers/aij-cnt.pdf).

## E2 — Latency factored as shape points × calibrated seconds-per-point

- Standard LLM serving decomposition: **total latency = TTFT + TPOT × number of output tokens**
  (TTFT = queueing + prefill before first token; TPOT = steady-state per-token decode time).
  [BentoML LLM inference metrics](https://bentoml.com/llm/llm-inference-basics/llm-inference-metrics),
  [AWS Neuron benchmarking guide](https://awsdocs-neuron.readthedocs-hosted.com/en/latest/libraries/nxd-inference/developer_guides/llm-inference-benchmarking-guide.html).
- Per-deployment calibration is practiced in production schedulers: llm-d trains a lightweight
  XGBoost regressor on request/server features (prompt length, prefix-cache hit rate, queue
  depth, KV utilization) against observed TTFT/TPOT, continuously retrained on a sliding window.
  [llm-d: Predicted-latency based scheduling](https://llm-d.ai/blog/predicted-latency-based-scheduling-for-llms).
- **ELIS** trains a response-length predictor (BGE encoder) and schedules by iterative shortest
  remaining time first. [arXiv:2505.09142](https://arxiv.org/abs/2505.09142).
- **LENS/PRISM**: SLO-aware engine scheduling plus predictive state-driven routing.
  [arXiv:2509.23384](https://arxiv.org/html/2509.23384v1).

## E3 — Declared degrade taxonomy

- **Imprecise computation** (Jane W.S. Liu et al., real-time systems): each task decomposes into
  a **mandatory part** (must complete for an acceptable result) and an **optional part** (refines
  the result); under overload the optional part is dropped, giving **graceful degradation** by
  design. [Scheduling imprecise computations to minimize total error](https://www.sciencedirect.com/science/article/abs/pii/0165607489901464),
  [Use of Imprecise Computation to Enhance Dependability](https://link.springer.com/chapter/10.1007/978-0-585-27316-7_6);
  applied to deep-learning services in [arXiv:2011.01112](https://arxiv.org/pdf/2011.01112).
- **RFC 5861** defines `stale-while-revalidate` (serve the stale response immediately,
  revalidate in the background — hides latency) and `stale-if-error` (serve stale on origin
  error — improves availability). [RFC 5861](https://httpwg.org/specs/rfc5861.html).
- **Model cascades / routing**: FrugalGPT (prompt adaptation, LLM approximation, small→large LLM
  cascade with early termination on confidence); RouteLLM (routers learned from preference data
  choosing weaker vs. stronger models). Active 2025–2026 follow-ons: MixLLM
  ([arXiv:2502.18482](https://arxiv.org/pdf/2502.18482)), MTRouter, cost-aware cascade serving
  ([arXiv:2606.27457](https://arxiv.org/pdf/2606.27457)).

## E4 — Suppression debt / anti-starvation

- **Neely's Lyapunov drift-plus-penalty with virtual queues**: each time-average constraint gets
  a virtual queue; **the queue (a debt counter) grows while the constraint is violated and
  inflates that constraint's priority until served**; mean-rate stability of the virtual queue
  implies constraint satisfaction. A parameter V trades utility against backlog, with explicit
  bounds. [arXiv:1008.3519](https://arxiv.org/pdf/1008.3519) (queue stability & convergence via
  Lyapunov optimization); [Drift-plus-penalty for finite-capacity queues](https://ieeexplore.ieee.org/document/9152999/);
  [Low-power dynamic scheduling](https://arxiv.org/pdf/1112.2797).
- OS **aging** (priority escalation for long-waiting tasks) is the classical anti-starvation
  mechanism in scheduling (textbook material; named in the seed survey).

## E5 — Generation lease + landing validation

- **Kung & Robinson 1981, "On Optimistic Methods for Concurrency Control"**: transactions run in
  three phases — read (private workspace), **validation** (serializability check against
  concurrent transactions), write (apply if validation succeeds, else abort/restart).
  "Optimistic" = rely on backup/retry, hoping conflicts are rare.
  [ACM TODS](https://dl.acm.org/doi/10.1145/319566.319567).
- **TL2** (Dice, Shalev & Shavit, DISC 2006): software transactional memory using commit-time
  locking and a **global version-clock**; transactions validate their read-set against the clock
  at commit, then stamp written locations with the new clock value.
  [Transactional Locking II](https://dl.acm.org/doi/10.1007/11864219_14).
- **Gray & Cheriton 1989, "Leases"**: a **time-based grant** providing efficient consistent
  access to cached data in distributed systems; the term "lease" originates here.
  [ACM SOSP/OSR](https://dl.acm.org/doi/10.1145/74851.74870),
  [PDF](https://web.eecs.umich.edu/~mosharaf/Readings/Leases.pdf).

## Converging LLM-agent literature (2025–2026)

- **RRARA** — "LLM-Enhanced Rapid-Reflex Async-Reflect Embodied Agent for Real-Time
  Decision-Making in Dynamically Changing Environments"
  ([arXiv:2506.07223](https://arxiv.org/abs/2506.07223)): hybrid agent — a rule-based reflex
  policy acts immediately while an **asynchronous LLM Reflector refines actions in situ**; a
  Time Conversion Mechanism translates inference delay into equivalent simulation frames;
  introduces latency-aware metrics (Respond Latency, Latency-to-Action Ratio) on HAZARD.
- **ASSCG** — "Just-Right Gating over Chattering for Fast-Slow LLM Planning in Autonomous
  Driving" ([arXiv:2606.25509](https://arxiv.org/abs/2606.25509), abstract verified 2026-07-23):
  makes **frame-level Query/Cache/Drop decisions to refresh, reuse, or suppress slow guidance**
  from an LLM planner; RWKV backbone; supervised fine-tuning + compute-aware RL; motivated by
  hand-designed rules that over-invoke or mis-invoke the slow system. Reported: −60% inference
  latency on one benchmark, +25% average speed on another.
- **Win Fast or Lose Slow** ([arXiv:2505.19481](https://arxiv.org/abs/2505.19481), NeurIPS
  2025): first systematic study of the latency–quality tradeoff for LLM agents; HFTBench and
  StreetFighter benchmarks; FPX framework adaptively selects model size and quantization;
  sacrificing quality for latency can improve downstream reward (up to 80% win-rate gain,
  +26.52% daily yield).
- **Rational Metareasoning for LLMs** ([arXiv:2410.05563](https://arxiv.org/pdf/2410.05563)):
  applies Russell & Wefald's value-of-computation to LLM inference-time reasoning.
- **Dynamic Speculative Agent Planning**
  ([arXiv:2509.01920](https://arxiv.org/pdf/2509.01920)): speculative execution for agent
  planning (named by the seed survey as converging work).
- **Faramesh** — "A Protocol-Agnostic Execution Control Plane for Autonomous Agent Systems"
  ([arXiv:2601.17744](https://arxiv.org/abs/2601.17744), full text v1 HTML verified 2026-07-23):
  mandatory Action Authorization Boundary; canonicalizes agent actions into a Canonical Action
  Representation, evaluates them deterministically against policy and state, issues binding
  PERMIT/DEFER/DENY before execution; **append-only provenance logging keyed by canonical action
  hashes** for auditability and deterministic replay; fail-closed semantics (any authorization
  failure resolves to deny/defer). Full-text findings: no latency budgets, deadlines, or
  latency-aware admission anywhere — latency appears only as a measured metric (p50/p95
  end-to-end decision latency, 2.24/9.61 ms); DEFER means pause-for-approval, not an alternative
  execution path — no fallback or degradation policy; staleness/drift between decision time and
  execution time is acknowledged as a limitation, not addressed. **Discrepancy resolved: the
  seed survey's "oracle latency budgets + audited fallback policy" characterization is not
  supported by the full text** (verified 2026-07-23, TASK-2).
- Adjacent caching/scheduling systems named in searches: Agentic Plan Caching
  ([arXiv:2506.14852](https://arxiv.org/abs/2506.14852)), AgenticCache
  ([arXiv:2604.24039](https://arxiv.org/abs/2604.24039)) — background Cache Updater
  asynchronously validates/refines cached plans; Continuum (KV-cache TTL for agent scheduling,
  [arXiv:2511.02230](https://arxiv.org/html/2511.02230v5)).
- No search hit found for "cognition horizon" used in landfall's sense (consistent with the
  seed survey's claim).

## Patent landscape (as recorded in the seed survey; one verified)

- **US 12,452,345** — "Managing artificial intelligence inference requests that are directed to
  an AI model external to a distributed cloud computing network"
  ([USPTO PDF](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/12452345),
  surfaced independently in this fan-out): assigns a time-based budget/ceiling of the latency
  requirement for reaching another data center; routing chooses among data centers with a
  qualified compute server within the time budget.
- Seed-survey list re-verified against actual claims 2026-07-23 via Google Patents (TASK-2):
  - **US 11,257,002** (Amazon, granted 2022-02-22) — "Dynamic accuracy-based deployment and
    monitoring of machine learning models in provider networks": accuracy-based traffic
    redirection among candidate ML models; latency appears only as a monitored model metric.
    Seed's "accuracy/latency-based ML model selection" holds for accuracy, weak on latency
    ([Google Patents](https://patents.google.com/patent/US11257002B2/en)).
  - **US 11,461,300** (SAP, granted 2022-10-04) — "Dynamic model server for multi-model machine
    learning inference services": claim 1 selects a model *server* based on the prediction
    request (consistent hashing for cache locality, LRU cache management); nothing accuracy- or
    latency-based. **Seed characterization not supported**
    ([Google Patents](https://patents.google.com/patent/US11461300B2/en)).
  - **US 2022/0004929 A1** (Google, published 2022-01-06, status pending as fetched) —
    "On-Device Machine Learning Platform": spec assigns context features expiration periods,
    with expired values deleted from stored training examples; claim text not retrievable from
    the Google Patents page, so "stale-context expiration *before inference*" is unconfirmed at
    claim level ([Google Patents](https://patents.google.com/patent/US20220004929A1/en)).
  - **US 11,717,748 B2** (Valve, granted 2023-08-08) — "Latency compensation using
    machine-learned prediction of user input": ML predicts a player's input so game control data
    is generated ahead of receipt "to compensate for latency". **Seed characterization
    confirmed** ([Google Patents](https://patents.google.com/patent/US11717748B2/en)).
  - **US 9,679,003** (IBM, granted 2017-06-13) — "Rendezvous-based optimistic concurrency
    control": OCC with the validate phase interleaved asynchronously with read/compute/write via
    read subscriptions / write publications. **OCC characterization confirmed**
    (concurrent-validation variant) ([Google Patents](https://patents.google.com/patent/US9679003B2/en)).
  - **US 8,396,831** (Microsoft, granted 2013-03-12) — "Optimistic serializable snapshot
    isolation": end-of-transaction read-set validation and phantom detection by re-scan over
    timestamp-versioned records. **OCC-validation characterization confirmed**
    ([Google Patents](https://patents.google.com/patent/US8396831B2/en)).
  - **US 8,972,306** (Raytheon, granted 2015-03-03) — "System and method for sensor tasking":
    claim 1 selects and schedules sensors from the state of a fuzzy cognitive map of
    environmental conditions; value-based ranking (Real Options Theory + market-based auction)
    appears only in dependent claim 7. Seed's "value-of-information sensor tasking" only
    partially supported ([Google Patents](https://patents.google.com/patent/US8972306B2/en)).
  - **EP 1,813,065 B1** (NXP, granted 2015-12-23) — "Device and method for event-triggered
    communication between and among a plurality of nodes": event-notification messaging with
    acknowledgment tracking between network nodes — event-triggered *networking*, not
    event-triggered control/sampling in the Heemels/Tabuada sense
    ([Google Patents](https://patents.google.com/patent/EP1813065B1/en)).
- Seed survey's conclusion (recorded verbatim as its claim, not endorsed here): "Nothing found
  claiming ex-ante world-drift-priced admission of an LLM call with post-hoc generation
  validation." Caveats stated in the seed: keyword search is not a professional clearance
  search; applications under 18 months old are invisible.

## Patentability doctrine (as framed in the seed survey)

- §102 novelty: a claim requiring all five elements together would likely survive strict
  anticipation — no single reference discloses the stack.
- §103 obviousness: high risk — each element is textbook in an adjacent field and motivation to
  combine is documented in the field itself (KSR: predictable combination of known elements
  performing known functions).
- §101 eligibility: nontrivial risk — the gate is arithmetic on abstract quantities (Alice).
- Seed survey's practical read: weak candidate for a broad patent; suggests publication /
  defensive disclosure as the alternative. (This is the seed's judgment; recorded as such.)

## Sources

1. Sun, Polyanskiy & Uysal — Remote Estimation of the Wiener Process over a Channel with Random Delay — https://arxiv.org/abs/1701.06734 · https://arxiv.org/abs/1707.02531
2. Ornee & Sun — Sampling and Remote Estimation for the Ornstein-Uhlenbeck Process Through Queues — https://ieeexplore.ieee.org/abstract/document/9437348 · https://arxiv.org/abs/2308.15401
3. Maatouk et al. — The Age of Incorrect Information — https://arxiv.org/abs/1907.06604
4. "X of Information" continuum survey — https://arxiv.org/pdf/2507.19657
5. Perceived age in query-based update systems — https://arxiv.org/pdf/2601.14075
6. Heemels, Johansson & Tabuada — Event-Triggered and Self-Triggered Control — https://link.springer.com/rwe/10.1007/978-3-030-44184-5_97 · https://heemels.tue.nl/research/event-triggered-and-self-triggered-control
7. Boddy & Dean — Solving Time-Dependent Planning Problems (IJCAI-89) — https://www.ijcai.org/Proceedings/89-2/Papers/021.pdf ; Deliberation Scheduling — https://www.semanticscholar.org/paper/31585092a9fe9d15e0e0c1ff030adca9a023cb2f
8. Russell & Wefald — Principles of Metareasoning — https://www.sciencedirect.com/science/article/abs/pii/000437029190015C ; Russell — Rationality and Intelligence — https://people.eecs.berkeley.edu/~russell/papers/aij-cnt.pdf
9. BentoML — Key metrics for LLM inference — https://bentoml.com/llm/llm-inference-basics/llm-inference-metrics
10. AWS Neuron — LLM inference benchmarking guide — https://awsdocs-neuron.readthedocs-hosted.com/en/latest/libraries/nxd-inference/developer_guides/llm-inference-benchmarking-guide.html
11. llm-d — Predicted-Latency Based Scheduling for LLMs — https://llm-d.ai/blog/predicted-latency-based-scheduling-for-llms
12. ELIS — https://arxiv.org/abs/2505.09142
13. LENS/PRISM two-layer scheduling — https://arxiv.org/html/2509.23384v1
14. Liu et al. — Scheduling imprecise computations — https://www.sciencedirect.com/science/article/abs/pii/0165607489901464 · https://link.springer.com/chapter/10.1007/978-0-585-27316-7_6 · https://arxiv.org/pdf/2011.01112
15. RFC 5861 — HTTP Cache-Control Extensions for Stale Content — https://httpwg.org/specs/rfc5861.html
16. MixLLM — https://arxiv.org/pdf/2502.18482 ; Cluster, Route, Escalate — https://arxiv.org/pdf/2606.27457
17. Neely — Queue Stability and Probability 1 Convergence via Lyapunov Optimization — https://arxiv.org/pdf/1008.3519 · https://ieeexplore.ieee.org/document/9152999/ · https://arxiv.org/pdf/1112.2797
18. Kung & Robinson — On Optimistic Methods for Concurrency Control — https://dl.acm.org/doi/10.1145/319566.319567
19. Dice, Shalev & Shavit — Transactional Locking II — https://dl.acm.org/doi/10.1007/11864219_14
20. Gray & Cheriton — Leases — https://dl.acm.org/doi/10.1145/74851.74870 · https://web.eecs.umich.edu/~mosharaf/Readings/Leases.pdf
21. RRARA — https://arxiv.org/abs/2506.07223
22. ASSCG — https://arxiv.org/abs/2606.25509 (abstract fetched 2026-07-23)
23. Win Fast or Lose Slow — https://arxiv.org/abs/2505.19481 · https://github.com/HaoKang-Timmy/LatencySensitiveBench
24. Rational Metareasoning for LLMs — https://arxiv.org/pdf/2410.05563
25. Dynamic Speculative Agent Planning — https://arxiv.org/pdf/2509.01920
26. Faramesh — https://arxiv.org/abs/2601.17744 · full text https://arxiv.org/html/2601.17744v1 (full text verified 2026-07-23)
27. Agentic Plan Caching — https://arxiv.org/abs/2506.14852 ; AgenticCache — https://arxiv.org/abs/2604.24039 ; Continuum — https://arxiv.org/html/2511.02230v5
28. US Patent 12,452,345 — https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/12452345
29. Seed document — PRIOR-ART.md, landfall repo root (deep-research sweep of 2026-07-23)
30. Patent re-verification (claims fetched 2026-07-23) — https://patents.google.com/patent/US11257002B2/en · https://patents.google.com/patent/US11461300B2/en · https://patents.google.com/patent/US20220004929A1/en · https://patents.google.com/patent/US11717748B2/en · https://patents.google.com/patent/US9679003B2/en · https://patents.google.com/patent/US8396831B2/en · https://patents.google.com/patent/US8972306B2/en · https://patents.google.com/patent/EP1813065B1/en
