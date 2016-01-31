package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Griewank minimum is 0 reached in (0, ..., 0)
// Recommended search domain is [-600, 600]
func Griewank(X []float64) float64 {
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
	// Instantiate a genetic algorithm
	ga := genalg.GA
	// Fitness function
	function := Griewank
	// Number of variables the function takes as input
	variables := 4
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 50; i++ {
		fmt.Println(i, "-", ga.Best.Fitness)
		ga.Enhance()
	}
}
