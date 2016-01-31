package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// DropWave minimum is -1 reached in (0, 0)
// Recommended search domain is [-5.12, 5.12]
func DropWave(X []float64) float64 {
	numerator := 1 + m.Cos(12*m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2)))
	denominator := 0.5*(m.Pow(X[0], 2)+m.Pow(X[1], 2)) + 2
	return -numerator / denominator
}

func main() {
	// Instantiate a population
	ga := genalg.GA
	// Fitness function
	function := DropWave
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 30; i++ {
		fmt.Println(ga.Best)
		ga.Enhance()
	}
}
