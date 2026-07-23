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
// verdict must be reproducible from its own recorded fields — and the record
// chain survives checkout: the lease freezes the verdict's arithmetic, so the
// landing log inherits the same reproducibility (TestLandingReproducible).
func TestVerdictReproducible(t *testing.T) {
	for _, vel := range []float64{0, 0.5, 1, 5, 100} {
		for _, debt := range []float64{0, 1, 3} {
			v := Route(planner, vel, 2.5, debt)
			c := Class{Name: v.Class, Points: v.Points, Budget: v.Budget, Degrade: v.Degrade, DebtFactor: planner.DebtFactor}
			r := Route(c, v.Velocity, v.SecondsPerPoint, v.Debt)
			if r != v {
				t.Errorf("verdict not reproducible from its own record:\n got %+v\nwant %+v", r, v)
			}
			l := Checkout(v, 1)
			if l.Class != v.Class || l.Budget != v.EffectiveBudget || l.PredictedDrift != v.PredictedDrift {
				t.Errorf("checkout must freeze the verdict's arithmetic into the lease:\n lease %+v\n verdict %+v", l, v)
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
	if v := RouteTiered(planner, 0.1, tiers); !v.Allow || v.Tier != "cloud-big" {
		t.Errorf("calm world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	// Hot world: only the local tier lands (3pt × 1s/pt × 3/s = 9 <= 10).
	if v := RouteTiered(planner, 3.0, tiers); !v.Allow || v.Tier != "local" {
		t.Errorf("hot world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	// Frantic world: nothing fits — suppression carries the best tier's math.
	if v := RouteTiered(planner, 50.0, tiers); v.Allow || v.Tier != "" {
		t.Errorf("frantic world: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	if v := RouteTiered(planner, 0.1, nil); v.Allow {
		t.Error("no tiers must suppress")
	}
}

// Debt queues are per-constraint: each tier carries its own suppression debt,
// a landed lower tier never repays a higher tier's queue, and one tier's debt
// never inflates another tier's budget.
func TestRouteTieredPerTierDebt(t *testing.T) {
	tiers := func(cloudDebt, localDebt float64) []Tier {
		return []Tier{
			{Name: "cloud-big", SecondsPerPoint: 10, Quality: 3, Debt: cloudDebt},
			{Name: "local", SecondsPerPoint: 1, Quality: 1, Debt: localDebt},
		}
	}
	// Hot world (velocity 3): cloud drift 90 > 10, local drift 9 <= 10.
	// With no debt the class is served locally — quality-starved of cloud.
	if v := RouteTiered(planner, 3.0, tiers(0, 0)); v.Tier != "local" {
		t.Fatalf("hot world, no debt: got tier %q (%s)", v.Tier, v.Arithmetic)
	}
	// The cloud queue keeps accruing while local landings don't repay it, and
	// inflates only the cloud tier's budget until cloud itself fits again:
	// debt 16 × factor 0.5 → effective budget 90 ≥ drift 90 — cloud admitted
	// and preferred over the also-fitting local tier.
	v := RouteTiered(planner, 3.0, tiers(16, 0))
	if !v.Allow || v.Tier != "cloud-big" {
		t.Fatalf("cloud debt must pull routing back up-tier: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	if v.Debt != 16 || v.EffectiveBudget != 90 {
		t.Errorf("verdict must record the chosen tier's debt: debt %g effective %g, want 16 and 90", v.Debt, v.EffectiveBudget)
	}
	// Isolation the other way: however large the local tier's queue grows, it
	// never inflates the cloud tier's budget — cloud stays suppressed.
	if v := RouteTiered(planner, 3.0, tiers(0, 1000)); v.Tier != "local" {
		t.Errorf("local debt must not inflate the cloud budget: got tier %q (%s)", v.Tier, v.Arithmetic)
	}
	// Frantic world (velocity 50): nothing fits debtless; only the tier that
	// owns the debt is admitted by it. local drift 150, debt 28 × 0.5 → 150.
	if v := RouteTiered(planner, 50.0, tiers(0, 28)); !v.Allow || v.Tier != "local" {
		t.Errorf("a tier's own debt must admit it: got tier %q allow=%v (%s)", v.Tier, v.Allow, v.Arithmetic)
	}
	// DebtFactor 0 declares starvation acceptable: no queue inflates anything.
	cosmetic := planner
	cosmetic.DebtFactor = 0
	if v := RouteTiered(cosmetic, 50.0, tiers(1000, 1000)); v.Allow {
		t.Errorf("DebtFactor 0 must never inflate any tier's budget (%s)", v.Arithmetic)
	}
}

// A RouteTiered verdict is reproducible from its own recorded fields, like a
// plain Route verdict: the recorded seconds-per-point and debt are the chosen
// tier's, so Route over them regenerates everything but the tier name.
func TestRouteTieredVerdictReproducible(t *testing.T) {
	tiers := []Tier{
		{Name: "cloud-big", SecondsPerPoint: 10, Quality: 3, Debt: 4},
		{Name: "local", SecondsPerPoint: 1, Quality: 1, Debt: 9},
	}
	for _, vel := range []float64{0, 0.5, 3, 50} {
		v := RouteTiered(planner, vel, tiers)
		c := Class{Name: v.Class, Points: v.Points, Budget: v.Budget, Degrade: v.Degrade, DebtFactor: planner.DebtFactor}
		r := Route(c, v.Velocity, v.SecondsPerPoint, v.Debt)
		r.Tier = v.Tier
		if v.Allow && r != v {
			t.Errorf("tiered verdict not reproducible from its own record:\n got %+v\nwant %+v", r, v)
		}
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
