package gago

import "math/rand"

// Mutator mutates an individual by modifying part of it's genome.
type Mutator interface {
	// Apply performs the mutation on an individual
	apply(individual *Individual, generator *rand.Rand)
}

// MutFNormal mutation modifies a float gene if a coin toss is under a defined
// mutation rate. It does so for each gene. The new gene value is a random value
// sampled from a normal distribution centered on the gene's current value and
// with the intensity parameter as it's standard deviation. Only works for
// floating point values.
type MutFNormal struct {
	// Mutation rate
	Rate float64
	// Standard deviation
	Std float64
}

// Apply normal mutation.
func (mfn MutFNormal) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Flip a coin and decide to mutate or not
		if generator.Float64() < mfn.Rate {
			// Sample from a normal distribution
			indi.Genome[i] = indi.Genome[i].(float64) * generator.NormFloat64() * mfn.Std
		}
	}
}

// MutSplice splices a genome in 3 and glues the parts back together in another
// order.
type MutSplice struct {
	// Mutation rate
	Rate float64
}

// Apply splice mutation.
func (ms MutSplice) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < ms.Rate {
		// Choose where to start and end the splice
		var end = rand.Intn(len(indi.Genome))
		var start = rand.Intn(len(indi.Genome))
		if end < start {
			start, end = end, start
		}
		// Split the genome into two
		var inner = make(Genome, end-start)
		copy(inner, indi.Genome[start:end])
		var outer = append(indi.Genome[:start], indi.Genome[end:]...)
		// Choose where to insert the splice
		var insert = rand.Intn(len(outer))
		// Splice and insert
		indi.Genome = append(
			outer[:insert],
			append(inner, outer[insert:]...)...,
		)
	}
}

// MutPermute permutes two genes.
type MutPermute struct {
	// Mutation rate
	Rate float64
	// Maximum number of permutation
	Max int
}

// Apply permutation mutation.
func (mp MutPermute) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < mp.Rate {
		for i := 0; i < generator.Intn(mp.Max); i++ {
			// Choose two points on the genome
			var (
				i = generator.Intn(len(indi.Genome))
				j = i
			)
			// Make sure both points are different
			for i == j {
				j = generator.Intn(len(indi.Genome))
			}
			// Permute the genes
			indi.Genome[i], indi.Genome[j] = indi.Genome[j], indi.Genome[i]
		}
	}
}
