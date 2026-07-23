package landfall

// The gate is ex-ante; the lease is its ex-post complement. Even an allowed
// call can be invalidated mid-flight: a salient event may bump the world
// generation, or actual drift may exceed what the gate predicted. The two
// checks share the drift unit, and the lease carries the gate's prediction
// with it, so every landing pairs (PredictedDrift, ActualDrift) in one
// record — the gate's calibration report, no join back to the verdict log.

// Lease is checked out when an allowed call launches and validated when the
// answer lands. Gen is the world generation at launch; the host bumps its
// generation on salient events. Budget is the effective budget the gate
// allowed against and PredictedDrift is the drift the gate bet on — both
// frozen at launch, the read set the landing validates against.
type Lease struct {
	Gen            uint64
	Class          string
	Budget         float64
	PredictedDrift float64
}

// Checkout freezes an allowed verdict into a Lease at launch: the class, the
// effective budget the gate allowed against, and the predicted drift it bet
// on. Pure; checking out a suppressed verdict is a host error — there is no
// call in flight to validate.
func Checkout(v Verdict, gen uint64) Lease {
	return Lease{Gen: gen, Class: v.Class, Budget: v.EffectiveBudget, PredictedDrift: v.PredictedDrift}
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

// Landing is the landing outcome plus the numbers that produced it — the
// ex-post sibling of Verdict. PredictedDrift (the gate's bet, frozen at
// checkout) and ActualDrift share a unit, so the landing log alone is the
// gate's calibration report. Like a Verdict, any logged Landing is
// reproducible from its own recorded fields.
type Landing struct {
	Outcome        Outcome
	Class          string
	Gen            uint64  // world generation at checkout
	CurrentGen     uint64  // world generation at landing
	Budget         float64 // effective budget frozen at checkout
	PredictedDrift float64 // the gate's ex-ante bet, frozen at checkout
	ActualDrift    float64 // drift accumulated during flight
}

// Land validates the lease against the world at landing time: the current
// generation and the actual drift accumulated during flight. Supersession
// trumps staleness — a salient event invalidates the answer regardless of
// how little drift accrued. Pure, like Route.
func (l Lease) Land(currentGen uint64, actualDrift float64) Landing {
	ld := Landing{
		Class: l.Class, Gen: l.Gen, CurrentGen: currentGen,
		Budget: l.Budget, PredictedDrift: l.PredictedDrift, ActualDrift: actualDrift,
	}
	switch {
	case currentGen != l.Gen:
		ld.Outcome = Superseded
	case actualDrift > l.Budget:
		ld.Outcome = Stale
	default:
		ld.Outcome = Landed
	}
	return ld
}
