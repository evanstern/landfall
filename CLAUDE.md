# landfall

Staleness-budgeted admission control for slow oracles, generalized from
promptworld's cognition horizon (`~/projects/promptworld/internal/cognition`).
Read DESIGN.md before changing anything — the invariants section defines what
makes this the pattern and not a knockoff.

Non-negotiables:

- `Route`, `RouteTiered`, and `Lease.Land` stay pure: no clock reads, no
  randomness, no model, no hidden state. Suppression debt is *passed in* and
  recorded in the verdict, never accrued inside the gate.
- Any logged Verdict must be reproducible from its own recorded fields
  (`TestVerdictReproducible` guards this — extend it when adding fields).
- The Estimator learns from successes only; failures belong to the host's
  breaker. Pessimistic seeds: uncalibrated fails toward the degrade floor.
- Budgets and points are doctrine — reviewed constants, no runtime tuning
  API. Do not add one.
- This package stays a leaf: no deps beyond stdlib. Velocity estimation,
  persistence, breakers, and transport belong to hosts.

Run `go test -race ./...` before considering any change done.
