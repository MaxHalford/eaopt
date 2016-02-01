package main

import (
	"fmt"
	m "math"
	"runtime"

	"github.com/MaxHalford/gago"
)

func simulate(start, end int) []float64 {
	data := []float64{}
	for x := start; x <= end; x++ {
		value := 1 * m.Pow(float64(x), 2)
		data = append(data, value)
	}
	return data
}

var (
	data = simulate(1, 20)
)

func leastSquares(X []float64) float64 {
	error := 0.0
	for i, target := range data {
		x := i + 1
		difference := target - (X[0] * m.Pow(float64(x), 2))
		error += m.Pow(difference, 2)
	}
	return error
}

func main() {
	runtime.GOMAXPROCS(-1)
	// Instantiate a population
	ga := gago.Default
	// Number of demes
	ga.NbDemes = 4
	// Number of individuals in each deme
	ga.NbIndividuals = 30
	// Initial random boundaries
	ga.Boundary = 10.0
	// Mutation rate
	ga.MRate = 0.2
	// Fitness function
	function := leastSquares
	// Number of variables the function takes as input
	variables := 1
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 1000; i++ {
		ga.Enhance()
	}
	fmt.Println(ga.Best)
}
