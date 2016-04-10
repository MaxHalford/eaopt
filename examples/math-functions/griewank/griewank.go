package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Griewank minimum is 0 reached in (0, ..., 0)
// Recommended search domain is [-600, 600]
func griewank(X []float64) float64 {
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
	// Instantiate a GA with 2 variables and the fitness function
	var ga = gago.GAFloat(2, griewank)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 50; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
