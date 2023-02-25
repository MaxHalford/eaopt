package eaopt

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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

	// Define function to minimize
	var ackley = func(x []float64) float64 {
		var (
			a, b, c = 20.0, 0.2, 2 * math.Pi
			s1, s2  float64
			d       = float64(len(x))
		)
		for _, xi := range x {
			s1 += xi * xi
			s2 += math.Cos(c * xi)
		}
		return -a*math.Exp(-b*math.Sqrt(s1/d)) - math.Exp(s2/d) + a + math.Exp(1)
	}

	// Run minimization
	x, y, err := de.Minimize(ackley, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	os.Stderr.WriteString(fmt.Sprintf("Found minimum of %.5f in %v\n", y, x))
	//fmt.Printf("Found MIN of %.5f in %v\n", y, x)
	// Output:
	// Found minimum of 0.00137 in [0.0004420129693826938 0.000195924625132926]
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
