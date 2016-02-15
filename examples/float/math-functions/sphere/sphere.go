package main

import (
	m "math"

	"github.com/MaxHalford/gago"
)

// Sphere function minimum is 0 reached in (0, ..., 0).
// Any search domain is fine.
func sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{sphere}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 100; i++ {
		ga.Best.Display()
		ga.Enhance()
	}
}
