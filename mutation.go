package gago

import "math/rand"

// Mutate an individual.
func mutate(indi *Individual, rate float64, generator *rand.Rand) {
	for i := range indi.Dna {
		// Flip a coin and decide to mutate or not
		coin := generator.Float64()
		if coin <= rate {
			indi.Dna[i] *= generator.NormFloat64() * 1
		}
	}
}
