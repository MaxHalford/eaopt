package gago

import "math/rand"

// Mutator mutates an individual by modifying part of it's genome.
type Mutator interface {
	// Apply performs the mutation on an individual
	apply(individual *Individual, generator *rand.Rand)
}

// Normal mutation modifies a float gene if a coin toss is under a defined
// mutation rate. It does so for each gene. The new gene value is a random value
// sampled from a normal distribution centered on the gene's current value and
// with the intensity parameter as it's standard deviation.
type Normal struct {
	// Mutation rate
	rate float64
	// Standard deviation
	std float64
}

// Apply normal mutation.
func (norm Normal) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Flip a coin and decide to mutate or not
		if generator.Float64() <= norm.rate {
			// Sample from a normal distribution
			indi.Genome[i] = indi.Genome[i].(float64) * generator.NormFloat64() * norm.std
		}
	}
}

// Corpus modifies a string gene by replacing it with a random character from
// a defined corpus.
type Corpus struct {
	// Mutation rate
	rate float64
	// Slice of strings
	corpus []string
}

func (crp Corpus) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Flip a coin and decide to mutate or not
		if generator.Float64() <= crp.rate {
			// Sample from the corpus
			indi.Genome[i] = crp.corpus[generator.Intn(len(crp.corpus))]
		}
	}
}
