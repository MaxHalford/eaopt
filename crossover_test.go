package gago

import (
	"math"
	"testing"
)

func TestCrossUniformFloat64(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = MakeVector(rng).(Vector)
		p2     = MakeVector(rng).(Vector)
		o1, o2 = CrossUniformFloat64(p1, p2, rng)
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

func TestNPointCrossover(t *testing.T) {
	var testCases = []struct {
		p1      []int
		p2      []int
		indexes []int
		o1      []int
		o2      []int
	}{
		{
			p1:      []int{1, 2, 3},
			p2:      []int{3, 2, 1},
			indexes: []int{1},
			o1:      []int{1, 2, 1},
			o2:      []int{3, 2, 3},
		},
		{
			p1:      []int{1, 2, 3, 4, 5, 6, 7},
			p2:      []int{7, 6, 5, 4, 3, 2, 1},
			indexes: []int{2, 5},
			o1:      []int{1, 2, 5, 4, 3, 6, 7},
			o2:      []int{7, 6, 3, 4, 5, 2, 1},
		},
	}
	for _, test := range testCases {
		var (
			n      = len(test.p1)
			o1, o2 = nPointCrossover(uncastInts(test.p1), uncastInts(test.p2), test.indexes)
		)
		for i := 0; i < n; i++ {
			if o1[i] != test.o1[i] || o2[i] != test.o2[i] {
				t.Error("Something went wrong during n-point crossover")
			}
		}
	}
}

func TestCrossNPointFloat64(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []float64{1, 2, 3}
		p2     = []float64{3, 2, 1}
		o1, o2 = CrossNPointFloat64(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossNPointFloat64 should not produce offsprings with different sizes")
	}
}

func TestCrossNPointInt(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []int{1, 2, 3}
		p2     = []int{3, 2, 1}
		o1, o2 = CrossNPointInt(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossNPointInt should not produce offsprings with different sizes")
	}
}

func TestCrossNPointString(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []string{"a", "b", "c"}
		p2     = []string{"c", "b", "a"}
		o1, o2 = CrossNPointString(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossNPointString should not produce offsprings with different sizes")
	}
}

func TestPMXCrossover(t *testing.T) {
	var testCases = []struct {
		p1 []int
		p2 []int
		a  int
		b  int
		o1 []int
		o2 []int
	}{
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  3,
			b:  7,
			o1: []int{9, 3, 2, 4, 5, 6, 7, 1, 8},
			o2: []int{1, 7, 3, 8, 2, 6, 5, 4, 9},
		},
	}
	for _, test := range testCases {
		var (
			n      = len(test.p1)
			o1, o2 = pmxCrossover(uncastInts(test.p1), uncastInts(test.p2), test.a, test.b)
		)
		for i := 0; i < n; i++ {
			if o1[i] != test.o1[i] || o2[i] != test.o2[i] {
				t.Error("Something went wrong during PMX crossover")
			}
		}
	}
}

func TestCrossPMXFloat64(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []float64{1, 2, 3}
		p2     = []float64{3, 2, 1}
		o1, o2 = CrossPMXFloat64(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossPMXFloat64 should not produce offsprings with different sizes")
	}
}

func TestCrossPMXInt(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []int{1, 2, 3}
		p2     = []int{3, 2, 1}
		o1, o2 = CrossPMXInt(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossPMXInt should not produce offsprings with different sizes")
	}
}

func TestCrossPMXString(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		p1     = []string{"a", "b", "c"}
		p2     = []string{"c", "b", "a"}
		o1, o2 = CrossPMXString(p1, p2, 2, rng)
	)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossPMXString should not produce offsprings with different sizes")
	}
}
