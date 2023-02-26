package eaopt

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func _ExampleOES() {
	// Instantiate DiffEvo
	var oes, err = NewDefaultOES()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fix random number generation
	oes.GA.RNG = rand.New(rand.NewSource(42))

	// Define function to minimize
	var rastrigin = func(x []float64) (y float64) {
		y = 10 * float64(len(x))
		for _, xi := range x {
			y += math.Pow(xi, 2) - 10*math.Cos(2*math.Pi*xi)
		}
		return y
	}

	// Run minimization
	X, y, err := oes.Minimize(rastrigin, []float64{0, 0})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	fmt.Printf("Found minimum of %.5f in %v\n", y, X)
	// Output:
	// Found minimum of 0.02270 in [0.006807861794722094 -0.008251984117745246]
}

func TestPointCrossover(t *testing.T) {
	var oes, err = NewDefaultOES()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	oes.Mu = []float64{1, 1}
	var (
		rng = newRand()
		p1  = oes.newPoint(rng).(*oesPoint)
		p2  = oes.newPoint(rng).(*oesPoint)
		p1c = p1.Clone().(*oesPoint)
		p2c = p2.Clone().(*oesPoint)
	)
	if reflect.DeepEqual(p1.x, p2.x) {
		t.Errorf("Expected mismatch")
	}
	if !reflect.DeepEqual(p1.x, p1c.x) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.x, p2c.x) {
		t.Errorf("Expected no mismatch")
	}
	p1.Crossover(p2, rng)
	if !reflect.DeepEqual(p1.x, p1c.x) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.x, p2c.x) {
		t.Errorf("Expected no mismatch")
	}
}

func TestNewOES(t *testing.T) {
	var testCases = []struct {
		f func() error
	}{
		{func() error { _, err := NewOES(2, 30, 1, 0.1, false, nil); return err }},
		{func() error { _, err := NewOES(100, 0, 1, 0.1, false, nil); return err }},
		{func() error { _, err := NewOES(100, 30, 0, 0.1, false, nil); return err }},
		{func() error { _, err := NewOES(100, 30, 1, 0, false, nil); return err }},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var err = tc.f()
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
}

func TestNewDefaultOES(t *testing.T) {
	var oes, err = NewDefaultOES()
	oes.GA.ParallelEval = true
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	if _, _, err = oes.Minimize(bowl, []float64{5, 5}); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
