package main

import (
	"fmt"
	m "math"
)

// StyblinskiTang minimum is -39.16599d reached in (-2.903534, ..., -2.903534)
// Recommended search domain is [-5, 5]
func styblinskiTang(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 4) - 16*m.Pow(x, 2) + 5*x
	}
	return sum / 2
}

func main() {
	// Instantiate a GA with 2 variables and the fitness function
	var ga = preset.Float64(2, styblinskiTang)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 50; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
