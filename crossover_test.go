package gago

import (
	"math"
	"testing"
)

func TestCrossUniform(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = MakeVector(rng).(Vector)
		p2     = MakeVector(rng).(Vector)
		o1, o2 = CrossUniform(p1, p2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossUniform should not produce offsprings with different sizes")
	}
	// Check new values are contained in hyper-rectangle defined by parents
	var (
		bounded = func(x, lower, upper float64) bool { return x > lower && x < upper }
		lower   float64
		upper   float64
	)
	for i := 0; i < len(p1); i++ {
		lower = math.Min(p1[i], p2[i])
		upper = math.Max(p1[i], p2[i])
		if !bounded(o1[i], lower, upper) || !bounded(o2[i], lower, upper) {
			t.Error("New values are not contained in hyper-rectangle")
		}
	}
}
