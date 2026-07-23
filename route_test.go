package landfall

import (
	"math"
	"testing"
)

var planner = Class{Name: "planner", Points: 3, Budget: 10, Degrade: DegradeReflex, DebtFactor: 0.5}

func TestRouteBoundary(t *testing.T) {
	// 3pt × 2s/pt = 6s wall. Budget 10 drift-units.
	cases := []struct {
		velocity float64
		allow    bool
	}{
		{0, true},      // paused world: zero staleness, always allowed
		{1.0, true},    // drift 6 <= 10
		{10.0 / 6, true}, // drift exactly at budget
		{2.0, false},   // drift 12 > 10
	}
	for _, c := range cases {
		v := Route(planner, c.velocity, 2.0, 0)
		if v.Allow != c.allow {
			t.Errorf("velocity %g: allow=%v want %v (%s)", c.velocity, v.Allow, c.allow, v.Arithmetic)
		}
	}
}

func TestRouteUnmeasurableVelocitySuppresses(t *testing.T) {
	for _, vel := range []float64{-1, math.NaN(), math.Inf(1)} {
		v := Route(planner, vel, 2.0, 0)
		if v.Allow {
			t.Errorf("velocity %v: must suppress (fail toward the floor)", vel)
		}
		if v.Degrade != DegradeReflex {
			t.Errorf("velocity %v: verdict must carry the degrade mode", vel)
		}
	}
}

func TestRouteDebtRelievesStarvation(t *testing.T) {
	// Velocity 2: drift 12 > budget 10 — suppressed with no debt.
	if v := Route(planner, 2.0, 2.0, 0); v.Allow {
		t.Fatalf("expected suppression at zero debt: %s", v.Arithmetic)
	}
	// Debt 1 × factor 0.5 → effective budget 15 ≥ drift 12 — admitted.
	v := Route(planner, 2.0, 2.0, 1)
	if !v.Allow {
		t.Fatalf("expected debt to admit the call: %s", v.Arithmetic)
	}
	if v.EffectiveBudget != 15 {
		t.Errorf("effective budget %g, want 15", v.EffectiveBudget)
	}
	// A class with DebtFactor 0 never inflates — starvation is declared acceptable.
	cosmetic := planner
	cosmetic.DebtFactor = 0
	if v := Route(cosmetic, 2.0, 2.0, 100); v.Allow {
		t.Error("DebtFactor 0 must never inflate the budget")
	}
}

// TestVerdictReproducible is the pattern's defining property: any logged
// verdict must be reproducible from its own recorded fields.
func TestVerdictReproducible(t *testing.T) {
	for _, vel := range []float64{0, 0.5, 1, 5, 100} {
		for _, debt := range []float64{0, 1, 3} {
			v := Route(planner, vel, 2.5, debt)
			c := Class{Name: v.Class, Points: v.Points, Budget: v.Budget, Degrade: v.Degrade, DebtFactor: planner.DebtFactor}
			r := Route(c, v.Velocity, v.SecondsPerPoint, v.Debt)
			if r != v {
				t.Errorf("verdict not reproducible from its own record:\n got %+v\nwant %+v", r, v)
			}
		}
	}
}

func TestRouteTiered(t *testing.T) {
	tiers := []Tier{
		{Name: "cloud-big", SecondsPerPoint: 10, Quality: 3},
		{Name: "cloud-fast", SecondsPerPoint: 3, Quality: 2},
		{Name: "local", SecondsPerPoint: 1, Quality: 1},
	}
	// Calm world: the best tier fits (3pt × 10s/pt × 0.1/s = 3 <= 10).
	if v := RouteTiered(planner, 0.1, tiers, 0); !v.Allow || v.Tier != "cloud-big" {
		t.Errorf("calm world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	// Hot world: only the local tier lands (3pt × 1s/pt × 3/s = 9 <= 10).
	if v := RouteTiered(planner, 3.0, tiers, 0); !v.Allow || v.Tier != "local" {
		t.Errorf("hot world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	// Frantic world: nothing fits — suppression carries the best tier's math.
	if v := RouteTiered(planner, 50.0, tiers, 0); v.Allow || v.Tier != "" {
		t.Errorf("frantic world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	if v := RouteTiered(planner, 0.1, nil, 0); v.Allow {
		t.Error("no tiers must suppress")
	}
}

func TestClassValidate(t *testing.T) {
	if err := planner.Validate(); err != nil {
		t.Fatalf("valid class rejected: %v", err)
	}
	bad := []Class{
		{Points: 3, Budget: 10, Degrade: DegradeSkip},                           // no name
		{Name: "x", Points: 0, Budget: 10, Degrade: DegradeSkip},                // no points
		{Name: "x", Points: 3, Budget: 0, Degrade: DegradeSkip},                 // no budget
		{Name: "x", Points: 3, Budget: math.Inf(1), Degrade: DegradeSkip},       // infinite budget
		{Name: "x", Points: 3, Budget: 10},                                      // no degrade mode
		{Name: "x", Points: 3, Budget: 10, Degrade: DegradeSkip, DebtFactor: -1}, // negative debt factor
	}
	for i, c := range bad {
		if err := c.Validate(); err == nil {
			t.Errorf("bad class %d passed validation: %+v", i, c)
		}
	}
}
