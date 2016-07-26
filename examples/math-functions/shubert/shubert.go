package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago/presets"
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
	// Instantiate a GA with 2 variables and the fitness function
	var ga = presets.Float64(2, shubert)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 30; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
