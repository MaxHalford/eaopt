// See http://tracer.lcc.uma.es/problems/onemax/onemax.html for a description
// of the problem.
package main

import (
	"fmt"
	"math/rand"

	"github.com/MaxHalford/gago"
)

// N is the size of each genome.
const N = 20

// Digits is a slice of ints.
type Digits []int

// Evaluate a slice of digits by summing the number of 1s.
func (X Digits) Evaluate() float64 {
	var sum int
	for _, d := range X {
		sum += d
	}
	return N - float64(sum) // We want to minimize the fitness, hence the reversing
}

// Mutate a slice of digits by permuting it's values.
func (X Digits) Mutate(rng *rand.Rand) {
	gago.MutPermuteInt(X, 3, rng)
}

// Crossover a slice of digits with another by applying 2-point crossover.
func (X Digits) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNXInt(X, Y.(Digits), 2, rng)
	return Digits(o1), Digits(o2)
}

// MakeDigits creates a random slice of digits by randomly picking 1s and 0s.
func MakeDigits(rng *rand.Rand) gago.Genome {
	var digits = make(Digits, N)
	for i := range digits {
		if rng.Float64() < 0.5 {
			digits[i] = 1
		}
	}
	return gago.Genome(digits)
}

func main() {
	var ga = gago.Generational(MakeDigits)
	ga.Initialize()

	for i := 1; i < 10; i++ {
		ga.Enhance()
	}
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
}
