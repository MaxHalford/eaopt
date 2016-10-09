package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago/presets"
)

// Sphere function minimum is 0 reached in (0, ..., 0).
// Any search domain is fine.
func sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a GA with 2 variables and the fitness function
	var ga = presets.Float64(2, sphere)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 8; i++ {
		ga.Enhance()
		// Display the current best solution
		fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
	}
}

>>> The best obtained solution is 0.018692
>>> The best obtained solution is 0.003595
>>> The best obtained solution is 0.000473
>>> The best obtained solution is 0.000311
>>> The best obtained solution is 0.000235
>>> The best obtained solution is 0.000064
>>> The best obtained solution is 0.000006
>>> The best obtained solution is 0.000001
