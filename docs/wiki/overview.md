---
name: overview
description: The whole landfall machine — pure ex-ante gate, ex-post lease validation, and the one stateful estimator — and the leaf-package doctrine.
kind: concept
sources:
  - README.md
  - DESIGN.md
  - class.go
  - route.go
  - lease.go
  - estimate.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Overview

landfall is a stdlib-only Go package implementing staleness-budgeted admission
control for slow oracles. A fast loop advances the world; a slow, expensive
oracle (LLM, heavy planner, remote model, human approver) answers with
latency. Every oracle query is a bet that the world at landing time will still
resemble the world at ask time; landfall prices that bet before the call is
spent and re-validates it when the answer lands.

## How it works

The package is four files, one lifecycle:

- `class.go` — [[decision-class]] registration: every oracle-reaching decision
  is a validated `Class` with points, budget, a legal degrade mode
  ([[degrade-modes]]), and a `DebtFactor`.
- `route.go` — the ex-ante gate: `Route` ([[route-gate]]) computes
  `predicted drift = velocity × points × secondsPerPoint` and allows iff it
  fits the debt-inflated budget ([[suppression-debt]]); `RouteTiered`
  ([[route-tiered]]) is the comparative form over a tier ladder. Both return a
  [[verdict]] carrying the full arithmetic.
- `lease.go` — the ex-post complement: `Checkout` freezes the allowed
  verdict's bet into a `Lease` ([[lease-lifecycle]]); `Lease.Land` validates
  it against the world at landing time and returns a `Landing`
  ([[landing-outcomes]]).
- `estimate.go` — the only stateful piece: a per-tier EWMA seconds-per-point
  [[estimator]] with spike rejection and a breach signal, governed by
  [[calibration-doctrine]].

`Route`, `RouteTiered`, and `Lease.Land` are pure: no clock reads, no
randomness, no model, no hidden state. Any logged Verdict or Landing is
reproducible from its own recorded fields — property-tested (see
[[test-suite]] and [[invariants-and-purity]]).

## Connections

Everything the gate consumes is measured outside and passed in: velocity,
seconds-per-point, debt, world generation — see [[host-obligations]]. The
drift unit shared by prediction and measurement is defined in [[drift-unit]];
velocity edge cases in [[velocity-semantics]]. The design's lineage and the
published pattern write-up are covered in [[lineage-and-prior-art]].

## Operational notes

The package is deliberately a leaf: no dependencies beyond stdlib. Budgets and
points are doctrine — reviewed constants with no runtime tuning API. Status at
this pin: core gate + lease + estimator implemented and tested
(`go test -race ./...`); velocity-estimator helpers and worked examples not
yet shipped.
