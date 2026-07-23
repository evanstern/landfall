---
name: decision-class
description: Class — the registered category of oracle-reaching decision; points, budget, degrade mode, DebtFactor, and Validate.
kind: component
sources:
  - class.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# Class — the registered decision category

`Class` in `class.go` is one registered category of oracle-reaching decision.
Every slow-oracle call a host makes must belong to a class; the doctrine is
"no unclassified oracle calls" — hosts validate every class at startup and
refuse to run otherwise.

## How it works

Fields:

- `Name` — identifies the class in every verdict and landing record.
- `Points` — the cost of the prompt shape, host-independent, Fibonacci by
  convention ([[points-latency-model]]).
- `Budget` — the answer's half-life in drift-units ([[drift-unit]]). Budgets
  are doctrine: reviewed constants, never runtime-tuned.
- `Degrade` — the declared fallback when the gate suppresses
  ([[degrade-modes]]).
- `DebtFactor` — governs starvation relief: the effective budget is
  `Budget × (1 + debt × DebtFactor)`, so persistent suppression eventually
  admits a call; zero declares starvation acceptable (e.g. purely cosmetic
  decisions). See [[suppression-debt]].

`Validate` enforces the invariants: non-empty name, positive points, a
positive finite budget (rejecting 0, NaN, ±Inf), a legal degrade mode, and a
non-negative DebtFactor. Each failure returns a distinct error naming the
class. `TestClassValidate` covers all six rejection cases.

## Connections

Consumed by [[route-gate]] and [[route-tiered]]; the class's name, points,
budget, and degrade mode are echoed into every [[verdict]]. Registration
without a legal degrade mode failing is invariant 2 of
[[invariants-and-purity]] — installing the gate forces the degraded mode to
exist. Hosts own the registry and the startup validation pass
([[host-obligations]]).

## Operational notes

A new class needs no latency observations: assign points and it inherits the
target tier's calibration from the [[estimator]]. Changing a budget is a
reviewed constant change, prompted by the calibration evidence in the landing
log ([[landing-outcomes]]), never by an automatic tuner.
