package eaopt

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func ExampleSPSO() {
	// Instantiate SPSO
	var spso, err = NewDefaultSPSO()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fix random number generation
	spso.GA.RNG = rand.New(rand.NewSource(42))

	// Define function to minimize
	var styblinskiTang = func(x []float64) (y float64) {
		for _, xi := range x {
			y += math.Pow(xi, 4) - 16*math.Pow(xi, 2) + 5*xi
		}
		return 0.5 * y
	}

	// Run minimization
	x, y, err := spso.Minimize(styblinskiTang, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	fmt.Printf("Found minimum of %.5f in %v\n", y, x)
	// Output:
	// Found minimum of -78.23783 in [-2.8586916496046983 -2.9619895273744623]
}

func TestParticleCrossover(t *testing.T) {
	var spso, err = NewDefaultSPSO()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	spso.NDims = 2
	var (
		rng = newRand()
		p1  = spso.newParticle(rng).(*Particle)
		p2  = spso.newParticle(rng).(*Particle)
		p1c = p1.Clone().(*Particle)
		p2c = p2.Clone().(*Particle)
	)
	if reflect.DeepEqual(p1.CurrentX, p2.CurrentX) {
		t.Errorf("Expected mismatch")
	}
	if !reflect.DeepEqual(p1.CurrentX, p1c.CurrentX) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.CurrentX, p2c.CurrentX) {
		t.Errorf("Expected no mismatch")
	}
	p1.Crossover(p2, rng)
	if !reflect.DeepEqual(p1.CurrentX, p1c.CurrentX) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.CurrentX, p2c.CurrentX) {
		t.Errorf("Expected no mismatch")
	}
}

func TestNewSPSO(t *testing.T) {
	var testCases = []struct {
		f func() error
	}{
		{func() error { _, err := NewSPSO(0, 30, -5, 5, 0.5, false, nil); return err }},
		{func() error { _, err := NewSPSO(40, 0, -5, 5, 0.5, false, nil); return err }},
		{func() error { _, err := NewSPSO(40, 30, 5, -5, 0.5, false, nil); return err }},
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

func TestNewDefaultSPSO(t *testing.T) {
	var spso, err = NewDefaultSPSO()
	spso.GA.ParallelEval = true
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	if _, _, err = spso.Minimize(bowl, 2); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
