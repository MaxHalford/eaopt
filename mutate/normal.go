package mutate

import "math/rand"

// MutNormal modifies a float gene if a coin toss is under a defined mutation
// rate. It does so for each gene. The new gene value is a random value sampled
// from a normal distribution centered on the gene's current value and with a
// standard deviation provided by the user. Only works for floating point
// values.
func MutNormal(genome []float64, rng *rand.Rand, rate, std float64) {
	for i := range genome {
		// Flip a coin and decide to mutate or not
		if rng.Float64() < rate {
			// Sample from a normal distribution
			genome[i] += rng.NormFloat64() * std
		}
	}
}
