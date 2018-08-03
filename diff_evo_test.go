package eaopt

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func ExampleDiffEvo() {
	// Instantiate DiffEvo
	var de, err = NewDefaultDiffEvo()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fix random number generation
	de.GA.RNG = rand.New(rand.NewSource(42))

	// Run minimization
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	X, y, err := de.Minimize(bowl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	fmt.Printf("Found minimum of %.5f in %v\n", y, X)
	// Output:
	// Found minimum of 0.00000 in [6.0034503946953274e-05 0.00013615643701126828]
}

func TestAgentCrossover(t *testing.T) {
	var de, err = NewDefaultDiffEvo()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	de.NDims = 2
	var (
		rng = newRand()
		p1  = de.newAgent(rng).(*Agent)
		p2  = de.newAgent(rng).(*Agent)
		p1c = p1.Clone().(*Agent)
		p2c = p2.Clone().(*Agent)
	)
	if reflect.DeepEqual(p1.X, p2.X) {
		t.Errorf("Expected mismatch")
	}
	if !reflect.DeepEqual(p1.X, p1c.X) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.X, p2c.X) {
		t.Errorf("Expected no mismatch")
	}
	p1.Crossover(p2, rng)
	if !reflect.DeepEqual(p1.X, p1c.X) {
		t.Errorf("Expected no mismatch")
	}
	if !reflect.DeepEqual(p2.X, p2c.X) {
		t.Errorf("Expected no mismatch")
	}
}

func TestNewDiffEvo(t *testing.T) {
	var testCases = []struct {
		f func() error
	}{
		{func() error { _, err := NewDiffEvo(0, 30, -5, 5, 0.5, 0.2, false, nil); return err }},
		{func() error { _, err := NewDiffEvo(1, 30, -5, 5, 0.5, 0.2, false, nil); return err }},
		{func() error { _, err := NewDiffEvo(2, 30, -5, 5, 0.5, 0.2, false, nil); return err }},
		{func() error { _, err := NewDiffEvo(3, 30, -5, 5, 0.5, 0.2, false, nil); return err }},
		{func() error { _, err := NewDiffEvo(40, 0, -5, 5, 0.5, 0.2, false, nil); return err }},
		{func() error { _, err := NewDiffEvo(40, 30, 5, -5, 0.5, 0.2, false, nil); return err }},
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

func TestNewDefaultDiffEvo(t *testing.T) {
	var de, err = NewDefaultDiffEvo()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	if _, _, err = de.Minimize(bowl, 2); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
