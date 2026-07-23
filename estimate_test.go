package landfall

import "testing"

func TestEstimatorFollowsDriftRejectsSpikes(t *testing.T) {
	e := NewEstimator(10, DefaultEstimatorConfig())
	// A one-shot spike (>3× estimate) is excluded from the EWMA.
	e.Sample(100)
	if got := e.Estimate(); got != 10 {
		t.Errorf("spike moved the estimate: %g", got)
	}
	// Systemic drift is followed.
	for i := 0; i < 50; i++ {
		e.Sample(20)
	}
	if got := e.Estimate(); got < 19 || got > 20 {
		t.Errorf("estimate %g did not converge toward 20", got)
	}
}

func TestEstimatorBreachSignalFiresOnceAndRearms(t *testing.T) {
	cfg := DefaultEstimatorConfig()
	e := NewEstimator(10, cfg)
	fill := func(spike bool, n int) (fired int) {
		v := 10.0
		if spike {
			v = 1000
		}
		for i := 0; i < n; i++ {
			if e.Sample(v) {
				fired++
			}
		}
		return fired
	}
	// No breach verdicts until the window is full; the first sample after a
	// fully-spiking window fires, and exactly once.
	if fired := fill(true, cfg.WindowSize); fired != 0 {
		t.Errorf("breach fired %d times while filling the window, want 0", fired)
	}
	if fired := fill(true, 5); fired != 1 {
		t.Errorf("breach fired %d times past a spiking window, want exactly 1", fired)
	}
	// Recovery re-arms; a fresh breach fires again.
	fill(false, cfg.WindowSize)
	if fired := fill(true, cfg.WindowSize+1); fired != 1 {
		t.Errorf("re-armed breach fired %d times, want exactly 1", fired)
	}
}
