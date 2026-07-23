---
name: drift-unit
description: The canonical drift unit — expected salient events during flight; budget as answer half-life; one unit shared by prediction and measurement.
kind: concept
sources:
  - DESIGN.md
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# The drift unit — salient events during flight

Budgets, predicted drift, and actual drift all share one canonical unit:
**expected salient events during flight**. A budget of 5 reads "this decision
survives 5 salient events." The unit is what makes the gate's arithmetic
meaningful, domain-portable, and self-calibrating.

## How it works

DESIGN.md's vocabulary table maps the source implementation's terms onto the
generalized ones: sim speed × tick rate → **velocity** (drift-units per
wall-second of decision-relevant change); points × sec/pt → **time-to-land**;
staleness in ticks → **predicted drift** (velocity × time-to-land);
BudgetTicks → **budget** (drift-units the answer survives — its half-life).

The critical property is that the ex-ante prediction and the ex-post
measurement count the *same thing*. The gate predicts drift in salient
events; the landing check receives actual drift in salient events. Because a
[[lease-lifecycle]] carries the prediction to the landing, every `Landing`
pairs predicted vs. actual drift in one record — the landing log **is** the
gate's calibration report, with no join back to the verdict log
(`TestLandingCarriesCalibrationPair`).

DESIGN.md grounds the stance in Age-of-Incorrect-Information: plain Age of
Information ages from delivery regardless of content, while AoII accrues age
only while the monitor's picture is actually wrong. "Expected salient events
during flight" is that content-aware refinement — the budget prices
decision-relevant change, not elapsed time.

## Connections

Consumed by [[route-gate]] (the comparison), declared per class in
[[decision-class]] (budgets are doctrine), interpreted at landing in
[[landing-outcomes]]. Velocity edge cases live in [[velocity-semantics]].
The AoII lineage is detailed in [[lineage-and-prior-art]]; a pattern of
predicted-vs-actual mismatch indicts velocity or tier calibration
([[calibration-doctrine]]).

## Operational notes

Domain examples from DESIGN.md: an editor assistant's velocity is edit
events/sec; trading velocity is realized volatility with budget as edge
half-life; robotics budgets are trajectory divergence. Hosts choose what
counts as "salient" — the gate only requires that prediction and measurement
agree on it.
