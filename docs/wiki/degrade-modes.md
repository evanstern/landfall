---
name: degrade-modes
description: The Degrade enum — skip | reflex | template | faster-tier; a legal degrade mode is mandatory at class registration.
kind: concept
sources:
  - class.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Degrade modes — the declared floor

`Degrade` in `class.go` names what the caller does when a class is
suppressed. Installing the gate forces the degraded mode to exist: a `Class`
without a legal `Degrade` fails `Validate`, so every slow-oracle dependency
acquires a declared, always-available fallback — graceful degradation by
construction.

## How it works

`Degrade` is a string type with four legal values, checked against the
package-private `legalDegrade` map:

- `DegradeSkip` ("skip") — nothing replaces the answer. Recorded, never
  silent.
- `DegradeReflex` ("reflex") — a deterministic fallback policy covers.
- `DegradeTemplate` ("template") — pre-authored output stands in.
- `DegradeFasterTier` ("faster-tier") — retry the route against a faster
  tier. [[route-tiered]] does this implicitly; a plain [[route-gate]] caller
  handles it explicitly.

The reflex floor is a spectrum, not a single mandatory callback — DESIGN.md
records this as one of the lessons taken from reading promptworld's real
implementation, which declared a faster-tier degrade but never wired it;
landfall wires it as `RouteTiered`.

## Connections

Declared per [[decision-class]] and enforced by its `Validate`. Every
[[verdict]] carries the class's degrade mode — including suppressions, so the
caller always knows its floor (`TestRouteUnmeasurableVelocitySuppresses`
checks this). Mandatory degrade is invariant 2 of [[invariants-and-purity]].
Executing the degrade path is the host's job ([[host-obligations]]).

## Operational notes

Skip is not silence: a skipped decision still produces a recorded suppression
verdict. The choice of mode is doctrine per class — a persistently hot world
running at the floor is the starvation scenario that [[suppression-debt]]
exists to relieve.
