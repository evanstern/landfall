package landfall

import "sync"

// EstimatorConfig tunes a tier's live seconds-per-point estimator. The zero
// value is not usable; start from DefaultEstimatorConfig.
type EstimatorConfig struct {
	// Alpha is the EWMA weight for a fresh sample.
	Alpha float64
	// SpikeFactor: a sample beyond SpikeFactor × the current estimate is a
	// spike — excluded from the EWMA but counted, so one-shot lag never
	// poisons the estimate while systemic drift is still followed.
	SpikeFactor float64
	// WindowSize is the ring of recent samples the breach signal watches.
	WindowSize int
	// BreachRate: when the spike rate over a full window exceeds this, the
	// estimator recommends recalibration (once, re-armed on recovery).
	BreachRate float64
}

// DefaultEstimatorConfig carries promptworld's tuned values.
func DefaultEstimatorConfig() EstimatorConfig {
	return EstimatorConfig{Alpha: 0.2, SpikeFactor: 3.0, WindowSize: 20, BreachRate: 0.3}
}

// Estimator is the live seconds-per-point estimate for one tier. Seed it
// pessimistically when uncalibrated — the doctrine is fail toward the degrade
// floor, never toward stale action. Feed it successes only; failures belong
// to a circuit breaker, not the latency model (a failed call's duration says
// nothing about how long a landed answer takes).
//
// Process-lifetime only: persist calibration under a human's hand and re-seed
// on restart, so the recorded baseline never drifts silently.
type Estimator struct {
	cfg      EstimatorConfig
	mu       sync.Mutex
	estimate float64
	window   []bool // true = spike
	wi, wn   int
	samples  int
	spikes   int
	breached bool
}

func NewEstimator(seed float64, cfg EstimatorConfig) *Estimator {
	return &Estimator{cfg: cfg, estimate: seed, window: make([]bool, cfg.WindowSize)}
}

// Estimate returns the current seconds-per-point.
func (e *Estimator) Estimate() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.estimate
}

// Stats returns the estimate, rolling spike rate, and lifetime counts.
func (e *Estimator) Stats() (estimate, spikeRate float64, samples, spikes int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.estimate, e.rateLocked(), e.samples, e.spikes
}

func (e *Estimator) rateLocked() float64 {
	if e.wn == 0 {
		return 0
	}
	n := 0
	for i := 0; i < e.wn; i++ {
		if e.window[i] {
			n++
		}
	}
	return float64(n) / float64(e.wn)
}

// Sample feeds one successful call's observed seconds-per-point. Returns true
// exactly when the spike rate over a full window first breaches BreachRate —
// the recalibration signal, re-armed once the rate falls back under.
func (e *Estimator) Sample(secPerPoint float64) (recalibrate bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.samples++
	spike := secPerPoint > e.cfg.SpikeFactor*e.estimate
	if spike {
		e.spikes++
	} else {
		e.estimate = (1-e.cfg.Alpha)*e.estimate + e.cfg.Alpha*secPerPoint
	}
	e.window[e.wi] = spike
	e.wi = (e.wi + 1) % e.cfg.WindowSize
	if e.wn < e.cfg.WindowSize {
		e.wn++
		return false
	}
	if e.rateLocked() > e.cfg.BreachRate {
		if !e.breached {
			e.breached = true
			return true
		}
		return false
	}
	e.breached = false
	return false
}
