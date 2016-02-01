package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// StyblinskiTang minimum is -39.16599d reached in (-2.903534, ..., -2.903534)
// Recommended search domain is [-5, 5]
func StyblinskiTang(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 4) - 16*m.Pow(x, 2) + 5*x
	}
	return sum / 2
}

func main() {
	// Instantiate a population
	ga := gago.Default
	// Fitness function
	function := StyblinskiTang
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 50; i++ {
		fmt.Println(ga.Best)
		ga.Enhance()
	}
}
