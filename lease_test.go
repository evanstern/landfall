package landfall

import "testing"

func TestLeaseLand(t *testing.T) {
	l := Lease{Gen: 7, Class: "planner", Budget: 10}
	if got := l.Land(7, 5); got != Landed {
		t.Errorf("in-budget same-gen: %v, want landed", got)
	}
	if got := l.Land(8, 0); got != Superseded {
		t.Errorf("gen bump: %v, want superseded", got)
	}
	if got := l.Land(7, 11); got != Stale {
		t.Errorf("over-budget: %v, want stale", got)
	}
	// Supersession trumps staleness: a salient event invalidates regardless.
	if got := l.Land(8, 11); got != Superseded {
		t.Errorf("gen bump + over-budget: %v, want superseded", got)
	}
	if got := l.Land(7, 10); got != Landed {
		t.Errorf("drift exactly at budget: %v, want landed", got)
	}
}
