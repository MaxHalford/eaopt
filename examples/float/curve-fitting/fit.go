package main

import (
	m "math"

	"github.com/MaxHalford/gago"
)

// Target function the GA has to approach
func f(x int, B []float64) float64 {
	return B[0] * m.Pow(float64(x), B[1])
}

// Simulate random data based on the target function
func simulate(start, end int, B []float64) []float64 {
	data := []float64{}
	for x := start; x <= end; x++ {
		data = append(data, f(x, B))
	}
	return data
}

var data = simulate(1, 20, []float64{1.0, 2.0})

// Least squares function to evaluate the difference between the GA individuals
// and the originally simulated data.
func leastSquares(B []float64) float64 {
	error := 0.0
	for i, y := range data {
		x := i + 1
		error += m.Pow(y-f(x, B), 2)
	}
	return error
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{leastSquares}
	// Number of demes
	ga.NbDemes = 4
	// Number of individuals in each deme
	ga.NbIndividuals = 100
	// Initialize the genetic algorithm
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 10000; i++ {
		ga.Enhance()
	}
	ga.Best.Display()
}
