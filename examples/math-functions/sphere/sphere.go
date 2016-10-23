package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
	"github.com/MaxHalford/gago/initialize"
	"github.com/MaxHalford/gago/preset"
)

// A Vector contains float64s.
type Vector struct {
	Values []float64
}

// Evaluate a Vector with the Sphere function (min of 0 in (0, ..., 0)).
func (v Vector) Evaluate() float64 {
	var sum float64
	for _, x := range v.Values {
		sum += math.Pow(x, 2)
	}
	return sum
}

// Mutate a Vector.
func (v Vector) Mutate(rng *rand.Rand) {
	v.Values.Splice(rng)
}

// Crossover a Vector with another Vector.
func (v Vector) Crossover(v2 interface{}, rng *rand.Rand) (gago.Genome, gago.Genome) {
	return v, v2.(Vector)
}

// MakeVector returns a random vector
func MakeVector(rng *rand.Rand) gago.Genome {
	return Vector{
		Values: initialize.UniformFloat64(5, rng, -10, 10),
	}
}

func main() {
	var ga = preset.SimAnn(MakeVector)
	ga.Enhance()
	fmt.Printf("Best -> %f\n", ga.Best.Fitness)
}
