package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Sphere function minimum is 0 reached in (0, ..., 0).
// Any search domain is fine.
func Sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := gago.Default
	// Fitness function
	function := Sphere
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 1000; i++ {
		fmt.Println(ga.Best)
		ga.Enhance()
	}
}
