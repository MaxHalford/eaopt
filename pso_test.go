package eaopt

import (
	"fmt"
	"math/rand"
)

func ExampleSPSO() {
	// Instantiate SPSO
	var spso, err = NewDefaultSPSO()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fix random number generation
	spso.GA.RNG = rand.New(rand.NewSource(42))

	// Run minimization
	var bowl = func(X []float64) (y float64) {
		for _, x := range X {
			y += x * x
		}
		return
	}
	X, y, err := spso.Minimize(bowl, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output best encountered solution
	fmt.Printf("Found minimum of %.5f in %v\n", y, X)
	// Output:
	// Found minimum of 0.00256 in [0.032498335711881876 -0.03878537756373712]
}
