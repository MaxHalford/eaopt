package gago

import "math/rand"

// Normal mutation individual modifies an individual's gene if a coin toss is
// under a defined mutation rate. It does so for each gene. The new gene value
// is a random value sampled from a normal distribution centered on the gene's
// current value and with the intensity parameter as it's standard deviation.
func Normal(indi *Individual, mutRate, mutIntensity float64, generator *rand.Rand) {
	for i := range indi.Dna {
		// Flip a coin and decide to mutate or not
		if generator.Float64() <= mutRate {
			indi.Dna[i] *= generator.NormFloat64() * mutIntensity
		}
	}
}
