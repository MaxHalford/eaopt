package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Rastrigin minimum is 0 reached in (0, ..., 0)
// Recommended search domain is [-5.12, 5.12]
func rastrigin(X []float64) float64 {
	sum := 10.0 * float64(len(X))
	for _, x := range X {
		sum += m.Pow(x, 2) - 10*m.Cos(2*m.Pi*x)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{rastrigin}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 40; i++ {
		fmt.Println(ga.Best.Fitness)
		ga.Enhance()
	}
}
