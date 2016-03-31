package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Shubert minimum is -186.7309 reached 18 times
// Recommended search domain is [-10, 10]
func shubert(X []float64) float64 {
	a := 0.0
	b := 0.0
	for i := 1; i <= 5; i++ {
		a += float64(i) * m.Cos(float64((i+1))*X[0]+float64(i))
		b += float64(i) * m.Cos(float64((i+1))*X[1]+float64(i))
	}
	return a * b
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{shubert}
	// More demes
	ga.NbDemes = 4
	// More individual
	ga.NbIndividuals = 100
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhancement
	for i := 0; i < 100; i++ {
		fmt.Println(ga.Best.Fitness)
		ga.Enhance()
	}
}
