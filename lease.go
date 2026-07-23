package landfall

// The gate is ex-ante; the lease is its ex-post complement. Even an allowed
// call can be invalidated mid-flight: a salient event may bump the world
// generation, or actual drift may exceed what the gate predicted. The two
// checks share the drift unit, so a log of (PredictedDrift, actual drift)
// pairs is the gate's own calibration report.

// Lease is checked out when an allowed call launches and validated when the
// answer lands. Gen is the world generation at launch; the host bumps its
// generation on salient events. Budget is the effective budget the gate
// allowed against, frozen at launch.
type Lease struct {
	Gen    uint64
	Class  string
	Budget float64
}

// Outcome is the landing verdict. Every outcome — not just Landed — must be
// recorded by the host: a dropped answer without a trace is how "why is the
// agent so passive" mystery bugs are made.
type Outcome int

const (
	// Landed: the world still matches; apply the answer.
	Landed Outcome = iota
	// Superseded: a salient event bumped the generation mid-flight. The
	// answer must not be applied as-is; hosts may still warm a cache with it.
	Superseded
	// Stale: no supersession, but actual drift exceeded the budget the gate
	// allowed against — the gate's prediction was wrong. Record loudly; a
	// pattern of Stale outcomes means the velocity estimate or the tier
	// calibration is lying.
	Stale
)

func (o Outcome) String() string {
	switch o {
	case Landed:
		return "landed"
	case Superseded:
		return "superseded"
	case Stale:
		return "stale"
	}
	return "unknown"
}

// Land validates the lease against the world at landing time: the current
// generation and the actual drift accumulated during flight. Supersession
// trumps staleness — a salient event invalidates the answer regardless of
// how little drift accrued. Pure, like Route.
func (l Lease) Land(currentGen uint64, actualDrift float64) Outcome {
	if currentGen != l.Gen {
		return Superseded
	}
	if actualDrift > l.Budget {
		return Stale
	}
	return Landed
}
