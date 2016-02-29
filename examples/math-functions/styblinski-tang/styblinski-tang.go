package main

import (
	m "math"

	"github.com/MaxHalford/gago"
)

// StyblinskiTang minimum is -39.16599d reached in (-2.903534, ..., -2.903534)
// Recommended search domain is [-5, 5]
func styblinskiTang(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 4) - 16*m.Pow(x, 2) + 5*x
	}
	return sum / 2
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{styblinskiTang}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 50; i++ {
		ga.Best.Display()
		ga.Enhance()
	}
}
