package main

import (
	m "math"

	"github.com/MaxHalford/gago"
)

// Griewank minimum is 0 reached in (0, ..., 0)
// Recommended search domain is [-600, 600]
func griewank(X []float64) float64 {
	sum := 0.0
	prod := 1.0
	for _, x := range X {
		sum += m.Pow(x, 2) / 4000
	}
	for i, x := range X {
		prod *= m.Cos(x / m.Sqrt(float64(i+1)))
	}
	return sum - prod + 1
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{griewank}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 50; i++ {
		ga.Best.Display()
		ga.Enhance()
	}
}
