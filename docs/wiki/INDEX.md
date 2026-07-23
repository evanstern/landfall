# landfall — grounding wiki

Code-grounded corpus for the landfall package: staleness-budgeted admission
control for slow oracles. One note per concept, each pinned to the commit it
was verified against. Start at [[overview]].

## Shape

- [[overview]] — the whole machine: pure gate → lease → estimator, and the leaf-package doctrine

## The gate (ex-ante)

- [[route-gate]] — `Route`: the pure admission decision and its arithmetic
- [[verdict]] — the decision record that carries its own arithmetic
- [[velocity-semantics]] — 0 is paused (always allowed); unmeasurable suppresses
- [[drift-unit]] — expected salient events during flight; budget as half-life
- [[route-tiered]] — `RouteTiered`: the comparative gate over a quality ladder of tiers

## Classes & degradation

- [[decision-class]] — `Class`: the registered, validated decision category
- [[degrade-modes]] — skip | reflex | template | faster-tier, forced to exist at registration
- [[points-latency-model]] — latency factored as shape points × per-tier seconds-per-point

## Starvation relief

- [[suppression-debt]] — debt-inflated effective budget guarantees eventual consultation
- [[per-tier-debt]] — virtual-queue doctrine under `RouteTiered`: one queue per tier

## The landing (ex-post)

- [[lease-lifecycle]] — `Checkout`: freezing the verdict's bet into an OCC read set
- [[landing-outcomes]] — `Land`: Landed | Superseded | Stale, supersession first

## Estimation & calibration

- [[estimator]] — the EWMA seconds-per-point estimator with spike rejection and breach signal
- [[calibration-doctrine]] — successes only, pessimistic seeds, human-owned baselines

## Doctrine & context

- [[invariants-and-purity]] — the six invariants that make an implementation the pattern
- [[host-obligations]] — everything landfall deliberately leaves to the host
- [[test-suite]] — what each test file guards, including the reproducibility properties
- [[lineage-and-prior-art]] — promptworld lineage, the theory behind each element, the defensive disclosure
