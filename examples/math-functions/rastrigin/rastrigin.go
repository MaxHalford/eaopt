package main

import (
	"fmt"
	m "math"
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
	// Instantiate a GA with 2 variables and the fitness function
	var ga = preset.Float64(2, rastrigin)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 40; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
