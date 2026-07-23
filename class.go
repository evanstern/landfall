// Package landfall is staleness-budgeted admission control for slow oracles.
//
// A fast loop advances the world; a slow oracle (LLM, heavy planner, remote
// model, human) answers with latency. Every oracle query is a bet that the
// world at landing time still resembles the world at ask time. landfall
// prices that bet before the call is spent — deterministically, with no model
// in the enforcement path — and re-validates it when the answer lands.
//
//	predicted drift = world velocity × predicted time-to-land
//	allow iff predicted drift ≤ the decision class's budget
//
// Generalized from promptworld's cognition horizon (internal/cognition).
package landfall

import (
	"fmt"
	"math"
)

// Degrade names what the caller does when a class is suppressed. Installing
// the gate forces the degraded mode to exist: a Class without a legal Degrade
// fails Validate, so every slow-oracle dependency has a declared fallback.
type Degrade string

const (
	// DegradeSkip: nothing replaces the answer. Recorded, never silent.
	DegradeSkip Degrade = "skip"
	// DegradeReflex: a deterministic fallback policy covers.
	DegradeReflex Degrade = "reflex"
	// DegradeTemplate: pre-authored output stands in.
	DegradeTemplate Degrade = "template"
	// DegradeFasterTier: retry the route against a faster tier (RouteTiered
	// does this implicitly; a plain Route caller handles it explicitly).
	DegradeFasterTier Degrade = "faster-tier"
)

var legalDegrade = map[Degrade]bool{
	DegradeSkip: true, DegradeReflex: true, DegradeTemplate: true, DegradeFasterTier: true,
}

// Class is one registered category of oracle-reaching decision.
//
// The latency model is factored shape × host: Points is the cost of the
// prompt shape (host-independent; Fibonacci by convention), and the host's
// measured seconds-per-point is supplied at route time. A new class needs no
// latency observations — assign points and it inherits the tier's calibration.
//
// Budget is the answer's half-life in drift-units — how much decision-relevant
// change the answer survives. With velocity measured in salient events per
// second, a Budget of 5 reads "this decision survives 5 salient events."
// Budgets are doctrine: reviewed constants, never runtime-tuned.
type Class struct {
	Name    string
	Points  int
	Budget  float64
	Degrade Degrade
	// DebtFactor governs starvation relief: each suppressed call accrues one
	// unit of debt, and the effective budget is Budget × (1 + debt×DebtFactor),
	// so persistent suppression eventually admits a call. Zero means starvation
	// is acceptable for this class (e.g. purely cosmetic decisions).
	DebtFactor float64
}

// Validate checks class invariants. Hosts should call this for every class at
// startup and refuse to run otherwise — no unclassified oracle calls.
func (c Class) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("landfall: class with empty name")
	}
	if c.Points <= 0 {
		return fmt.Errorf("landfall: class %q points %d not positive", c.Name, c.Points)
	}
	if c.Budget <= 0 || math.IsInf(c.Budget, 0) || math.IsNaN(c.Budget) {
		return fmt.Errorf("landfall: class %q budget %v not a positive finite number", c.Name, c.Budget)
	}
	if !legalDegrade[c.Degrade] {
		return fmt.Errorf("landfall: class %q degrade %q not one of skip|reflex|template|faster-tier", c.Name, c.Degrade)
	}
	if c.DebtFactor < 0 {
		return fmt.Errorf("landfall: class %q debt factor %v negative", c.Name, c.DebtFactor)
	}
	return nil
}
