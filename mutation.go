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

// MutSplice a genome and glue it back in another order.
type MutSplice struct {
	// Mutation rate
	Rate float64
}

// Apply splice mutation.
func (ms MutSplice) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < ms.Rate {
		// Choose where to split the genome
		var split = rand.Intn(len(indi.Genome))
		// Splice and glue
		indi.Genome = append(indi.Genome[split:], indi.Genome[:split]...)
	}
}

// MutPermute permutes two genes.
type MutPermute struct {
	// Mutation rate
	Rate float64
}

// Apply permutation mutation.
func (mp MutPermute) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < mp.Rate {
		// Choose two points on the genome
		var (
			r = generator.Perm(len(indi.Genome))[:2]
			i = r[0]
			j = r[1]
		)
		// Permute the genes
		indi.Genome[i], indi.Genome[j] = indi.Genome[j], indi.Genome[i]
	}
}
