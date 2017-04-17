package main

import (
	"fmt"
	m "math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector with the Drop-Wave function which takes two variables as
// input and reaches a minimum of -1 in (0, 0).
func (X Vector) Evaluate() float64 {
	var (
		numerator   = 1 + m.Cos(12*m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2)))
		denominator = 0.5*(m.Pow(X[0], 2)+m.Pow(X[1], 2)) + 2
	)
	return -numerator / denominator
}

// Mutate a Vector by applying by resampling each element from a normal
// distribution with probability 0.8.
func (X Vector) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying 2-point crossover.
func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNXFloat64(X, Y.(Vector), 2, rng) // Returns two float64 slices
	return Vector(o1), Vector(o2)
}

// MakeVector returns a random vector by generating 2 values uniformally
// distributed between -10 and 10.
func MakeVector(rng *rand.Rand) gago.Genome {
	return Vector(gago.InitUnifFloat64(2, -10, 10, rng))
}

func main() {
	var ga = gago.Generational(MakeVector)
	ga.Initialize()

	fmt.Printf("Best fitness at generation 0: %f\n", ga.Best.Fitness)
	for i := 1; i < 10; i++ {
		ga.Enhance()
		fmt.Printf("Best fitness at generation %d: %f\n", i, ga.Best.Fitness)
	}
}
