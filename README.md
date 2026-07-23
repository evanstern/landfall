# landfall

Staleness-budgeted admission control for slow oracles.

A fast loop advances the world. A slow, expensive oracle — an LLM, a heavy
planner, a human approver — answers with latency. Every oracle query is a bet
that the world at landing time will still resemble the world at ask time.
landfall prices that bet **before the call is spent**, with no model in the
enforcement path, and re-validates it when the answer lands:

```
predicted drift = world velocity × predicted time-to-land
allow iff predicted drift ≤ the decision class's budget
```

```go
c := landfall.Class{Name: "planner", Points: 3, Budget: 10, Degrade: landfall.DegradeReflex, DebtFactor: 0.5}

v := landfall.Route(c, velocity, tier.Estimate(), debt)
if !v.Allow {
    record(v)          // every suppression carries its own arithmetic
    runDegrade(v)      // the declared floor: skip | reflex | template | faster-tier
    return
}
lease := landfall.Lease{Gen: world.Gen(), Class: c.Name, Budget: v.EffectiveBudget}
answer := oracle.Ask(...)                       // slow; world keeps moving
switch lease.Land(world.Gen(), world.DriftSince(lease.Gen)) {
case landfall.Landed:     apply(answer)
case landfall.Superseded: warmCache(answer)     // salient event won mid-flight
case landfall.Stale:      recordLoudly(...)     // the gate's prediction was wrong
}
```

The gate (`Route`, `RouteTiered`), the landing check (`Lease.Land`), and the
class registry are pure and deterministic — same inputs, same verdict, and
any logged verdict is reproducible from its own recorded fields (this is
property-tested). The only stateful piece is the per-tier latency
`Estimator`: EWMA with spike rejection, fed by successes only, seeded
pessimistically so an uncalibrated system fails toward its floor, never
toward stale action.

Lineage: generalized from the cognition horizon in
[promptworld](https://github.com/evanstern/promptworld)
(`internal/cognition`) — a deterministic gate deciding whether an LLM call is
allowed at the current sim speed, based on how stale its answer will be when
it lands. See [DESIGN.md](DESIGN.md) for the full design, what the source
implementation taught, invariants, and open questions.

Status: planted 2026-07-23. Core gate + lease + estimator implemented and
tested (`go test -race ./...`); velocity-estimator helpers and worked
examples not yet.
