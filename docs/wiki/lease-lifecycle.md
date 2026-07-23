---
name: lease-lifecycle
description: Lease and Checkout — freezing the allowed verdict's generation, effective budget, and predicted drift into an OCC read set at launch.
kind: pipeline
sources:
  - lease.go
verified_against: a9b5701e58fbb844058ed08ba23d95082fa77362
---

# The lease — ask-gate → flight → land-validate

The gate is ex-ante; even an allowed call can be invalidated mid-flight. The
`Lease` in `lease.go` is the ex-post complement: checked out when an allowed
call launches, validated when the answer lands. The shape is optimistic
concurrency control — the oracle call is the read phase, landing is
validation.

## How it works

`Lease` fields: `Gen` (world generation at launch — the host bumps its
generation on salient events), `Class`, `Budget` (the *effective* budget the
gate allowed against), and `PredictedDrift` (the drift the gate bet on). All
four are frozen at launch: the read set the landing validates against.

`Checkout(v Verdict, gen uint64) Lease` copies `v.Class`,
`v.EffectiveBudget`, and `v.PredictedDrift` into the lease alongside the
current generation. It is pure; checking out a *suppressed* verdict is a host
error — there is no call in flight to validate.

Carrying `PredictedDrift` in the lease is a resolved design question
(DESIGN.md): the lease is the OCC read set, and recording the admission-time
prediction in it lets every landing pair predicted vs. actual drift without a
join back to the verdict log. `TestVerdictReproducible` asserts checkout
freezes the verdict's arithmetic; `TestLandingCarriesCalibrationPair` asserts
the pair survives to the landing.

## Connections

Consumes an allowed [[verdict]] from [[route-gate]] or [[route-tiered]];
because the frozen budget is the *effective* one, a debt-admitted call is
validated against the same inflated tolerance it was admitted under
([[suppression-debt]]). Validation happens in `Lease.Land`
([[landing-outcomes]]). Generation bumps are host work
([[host-obligations]]); the predicted/actual pairing is the calibration story
of [[drift-unit]].

## Operational notes

`Gen` is a `uint64` the host owns; landfall never increments it. The lease is
a value type with no clock or state — losing one simply means a landing that
can't be validated, which is itself a host bug worth surfacing.
