package landfall

import "testing"

func TestLeaseLand(t *testing.T) {
	l := Lease{Gen: 7, Class: "planner", Budget: 10, PredictedDrift: 6}
	if got := l.Land(7, 5).Outcome; got != Landed {
		t.Errorf("in-budget same-gen: %v, want landed", got)
	}
	if got := l.Land(8, 0).Outcome; got != Superseded {
		t.Errorf("gen bump: %v, want superseded", got)
	}
	if got := l.Land(7, 11).Outcome; got != Stale {
		t.Errorf("over-budget: %v, want stale", got)
	}
	// Supersession trumps staleness: a salient event invalidates regardless.
	if got := l.Land(8, 11).Outcome; got != Superseded {
		t.Errorf("gen bump + over-budget: %v, want superseded", got)
	}
	if got := l.Land(7, 10).Outcome; got != Landed {
		t.Errorf("drift exactly at budget: %v, want landed", got)
	}
}

// A Stale landing carries the calibration pair (predicted vs. actual drift)
// in its own record — no join back to the verdict log.
func TestLandingCarriesCalibrationPair(t *testing.T) {
	v := Route(planner, 2.0, 2.0, 1) // drift 12, effective budget 15 — allowed
	if !v.Allow {
		t.Fatalf("expected an allowed verdict: %s", v.Arithmetic)
	}
	l := Checkout(v, 42)
	if l.PredictedDrift != v.PredictedDrift || l.Budget != v.EffectiveBudget {
		t.Fatalf("checkout must freeze the verdict's arithmetic: lease %+v, verdict %+v", l, v)
	}
	ld := l.Land(42, 20)
	if ld.Outcome != Stale {
		t.Fatalf("drift 20 > budget 15: %v, want stale", ld.Outcome)
	}
	if ld.PredictedDrift != 12 || ld.ActualDrift != 20 {
		t.Errorf("landing must pair predicted vs actual: got (%g, %g), want (12, 20)", ld.PredictedDrift, ld.ActualDrift)
	}
}

// TestLandingReproducible is the reproducibility property extended to the
// ex-post record: any logged Landing must be reproducible from its own
// recorded fields, exactly like a Verdict.
func TestLandingReproducible(t *testing.T) {
	for _, vel := range []float64{0, 0.5, 1, 2} {
		for _, gen := range []uint64{3, 4} {
			for _, actual := range []float64{0, 9, 40} {
				v := Route(planner, vel, 2.5, 1)
				ld := Checkout(v, 3).Land(gen, actual)
				l := Lease{Gen: ld.Gen, Class: ld.Class, Budget: ld.Budget, PredictedDrift: ld.PredictedDrift}
				if r := l.Land(ld.CurrentGen, ld.ActualDrift); r != ld {
					t.Errorf("landing not reproducible from its own record:\n got %+v\nwant %+v", r, ld)
				}
			}
		}
	}
}
