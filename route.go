package landfall

import (
	"fmt"
	"math"
)

// Verdict is the gate's decision plus the arithmetic that produced it. Every
// verdict — especially a suppression — should be recorded by the host; the
// Arithmetic string and the numeric fields make each gate auditable and each
// decision reproducible from its own record.
type Verdict struct {
	Allow bool
	Class string
	Tier  string // chosen tier; empty for plain Route or when no tier fits

	Points           int
	SecondsPerPoint  float64
	Velocity         float64 // drift-units per wall-second, as fed by the host
	PredictedWallSec float64 // Points × SecondsPerPoint
	PredictedDrift   float64 // PredictedWallSec × Velocity
	Budget           float64 // the class's declared budget
	EffectiveBudget  float64 // debt-inflated: Budget × (1 + Debt×DebtFactor)
	Debt             float64

	Degrade    Degrade // what the caller does when !Allow
	Arithmetic string  // the equation, human-readable
}

// Route decides whether a decision may go to the oracle: pure arithmetic over
// the class's registered values, the host's seconds-per-point for the target
// tier, the current world velocity, and the class's accrued suppression debt.
// No model, no randomness, no wall-clock reads — the same inputs always
// produce the same verdict.
//
// Velocity semantics: 0 is a paused world — zero staleness, always allowed.
// (promptworld's gate used ≤0 for "uncapped speed"; landfall inverts this:
// velocity is the measured drift rate, and an unmeasurable rate is expressed
// as NaN/Inf/negative, which suppresses — an uncalibrated world fails toward
// the degrade floor, never toward stale action.)
func Route(c Class, velocity, secondsPerPoint, debt float64) Verdict {
	v := Verdict{
		Class: c.Name, Points: c.Points, SecondsPerPoint: secondsPerPoint,
		Velocity: velocity, Budget: c.Budget, Debt: debt, Degrade: c.Degrade,
	}
	v.PredictedWallSec = float64(c.Points) * secondsPerPoint
	v.EffectiveBudget = c.Budget * (1 + debt*c.DebtFactor)
	if velocity < 0 || math.IsNaN(velocity) || math.IsInf(velocity, 0) {
		v.Arithmetic = fmt.Sprintf("%dpt x %.1fs/pt at unmeasurable velocity (%v) - suppressed", c.Points, secondsPerPoint, velocity)
		return v
	}
	v.PredictedDrift = v.PredictedWallSec * velocity
	v.Allow = v.PredictedDrift <= v.EffectiveBudget
	rel := "<="
	if !v.Allow {
		rel = ">"
	}
	v.Arithmetic = fmt.Sprintf("%dpt x %.1fs/pt x %g/s = %.1f drift %s budget %.1f (base %.1f, debt %.1f)",
		c.Points, secondsPerPoint, velocity, v.PredictedDrift, rel, v.EffectiveBudget, c.Budget, debt)
	return v
}

// Tier is one oracle the class could route to, with the host's current
// seconds-per-point estimate for it. Quality orders tiers best-first in the
// caller's judgment; RouteTiered prefers higher quality.
//
// Debt is the suppression-debt queue owed to this tier's constraint. Debt
// queues are per-constraint, and under RouteTiered each tier is its own
// constraint: "consulted at quality ≥ this tier's, often enough." The host
// accrues a tier's queue on every verdict that lands below its quality (or
// not at all), and a landing at tier T resets the queues of T and every
// lower-quality tier — a landed local call does NOT repay the cloud tier's
// debt. The cloud queue keeps growing and inflating the cloud tier's own
// effective budget until the cloud tier itself is served.
type Tier struct {
	Name            string
	SecondsPerPoint float64
	Quality         float64
	Debt            float64
}

// RouteTiered is the comparative form of the gate: route to the best tier
// whose predicted drift fits that tier's debt-inflated budget — the big slow
// oracle when the world is calm, a faster one when it is hot, the degrade
// floor when nothing lands. Pure like Route: each tier's debt is passed in
// on the Tier and the chosen tier's debt is recorded in the verdict. One
// tier's debt never inflates another tier's budget. With no tiers, or none
// that fit, the returned verdict is a suppression carrying the best-quality
// tier's arithmetic.
func RouteTiered(c Class, velocity float64, tiers []Tier) Verdict {
	if len(tiers) == 0 {
		v := Route(c, velocity, 0, 0)
		v.Allow = false
		v.Arithmetic = "no tiers offered - suppressed"
		return v
	}
	best := 0
	for i, t := range tiers {
		if t.Quality > tiers[best].Quality {
			best = i
		}
	}
	chosen, ok := best, false
	for i, t := range tiers {
		v := Route(c, velocity, t.SecondsPerPoint, t.Debt)
		if !v.Allow {
			continue
		}
		if !ok || t.Quality > tiers[chosen].Quality {
			chosen, ok = i, true
		}
	}
	v := Route(c, velocity, tiers[chosen].SecondsPerPoint, tiers[chosen].Debt)
	v.Tier = tiers[chosen].Name
	if !ok {
		v.Allow = false
		v.Tier = ""
		v.Arithmetic = fmt.Sprintf("best of %d tiers (%s): %s", len(tiers), tiers[chosen].Name, v.Arithmetic)
	}
	return v
}
