package eaopt

import (
	"fmt"
	"math/rand"
)

func ExampleDiffEvo() {
	// Instantiate DiffEvo
	var de, err = NewDefaultDiffEvo()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fix random number generation
	de.GA.RNG = rand.New(rand.NewSource(42))

	// Run minimization
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	X, y, err := de.Minimize(bowl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	fmt.Printf("Found minimum of %.5f in %v\n", y, X)
	// Output:
	// Found minimum of 0.00000 in [6.0034503946953274e-05 0.00013615643701126828]
}
