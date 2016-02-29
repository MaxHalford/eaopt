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
// with the intensity parameter as it's standard deviation. Only works for
// numeric genes.
type Normal struct {
	// Mutation rate
	Rate float64
	// Standard deviation
	Std float64
}

// Apply normal mutation.
func (norm Normal) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Flip a coin and decide to mutate or not
		if generator.Float64() <= norm.Rate {
			// Sample from a normal distribution
			indi.Genome[i] = indi.Genome[i].(float64) * generator.NormFloat64() * norm.Std
		}
	}
}

// Swap two elements of a genome. Preferably used for long genomes of any kind.
type Swap struct{}

// Apply swap mutation.
func (swap Swap) apply(indi *Individual, generator *rand.Rand) {
	// Select two random positions in the genome
	var posA = generator.Intn(len(indi.Genome))
	var posB = generator.Intn(len(indi.Genome))
	// Make sure the two positions are different
	for posA == posB {
		posB = generator.Intn(len(indi.Genome))
	}
	// Perform the swap
	indi.Genome[posA], indi.Genome[posB] = indi.Genome[posB], indi.Genome[posA]
}
