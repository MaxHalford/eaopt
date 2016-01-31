package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Rastrigin minimum is 0 reached in (0, ..., 0)
// Recommended search domain is [-5.12, 5.12]
func Rastrigin(X []float64) float64 {
	sum := 10.0 * float64(len(X))
	for _, x := range X {
		sum += m.Pow(x, 2) - 10*m.Cos(2*m.Pi*x)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := genalg.GA
	// Fitness function
	function := Rastrigin
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 20; i++ {
		fmt.Println(ga.Best)
		ga.Enhance()
	}
}
