package gago2

import "math/rand"

type Vector []float64

func (v Vector) Evaluate() float64 {
	var sum float64
	for _, x := range v {
		sum += x
	}
	return sum
}

func (v Vector) Mutate(rng *rand.Rand) {}

func (v Vector) Crossover(v2 interface{}, rng *rand.Rand) (Genome, Genome) { return v, v2.(*Vector) }

func MakeVector(rng *rand.Rand) Vector { return []float64{1, 2, 3} }
