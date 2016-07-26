package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago/presets"
)

// DropWave minimum is -1 reached in (0, 0)
// Recommended search domain is [-5.12, 5.12]
func dropWave(X []float64) float64 {
	numerator := 1 + m.Cos(12*m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2)))
	denominator := 0.5*(m.Pow(X[0], 2)+m.Pow(X[1], 2)) + 2
	return -numerator / denominator
}

func main() {
	// Instantiate a GA with 2 variables and the fitness function
	var ga = presets.Float64(2, dropWave)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 30; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
