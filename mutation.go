package gago

import "math/rand"

// Mutator mutates an individual by modifying part of it's genome.
type Mutator interface {
	// Apply performs the mutation on an individual
	apply(individual *Individual, generator *rand.Rand)
}

// FloatNormal mutation modifies a float gene if a coin toss is under a defined
// mutation rate. It does so for each gene. The new gene value is a random value
// sampled from a normal distribution centered on the gene's current value and
// with the intensity parameter as it's standard deviation. Only works for
// floating point values.
type FloatNormal struct {
	// Mutation rate
	Rate float64
	// Standard deviation
	Std float64
}

// Apply normal mutation.
func (fnorm FloatNormal) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Flip a coin and decide to mutate or not
		if generator.Float64() < fnorm.Rate {
			// Sample from a normal distribution
			indi.Genome[i] = indi.Genome[i].(float64) * generator.NormFloat64() * fnorm.Std
		}
	}
}

// Splice a genome and glue it back in another order.
type Splice struct {
	// Mutation rate
	Rate float64
}

// Apply splice mutation.
func (spl Splice) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < spl.Rate {
		// Choose where to split the genome
		var split = rand.Intn(len(indi.Genome))
		// Splice and glue
		indi.Genome = append(indi.Genome[split:], indi.Genome[:split]...)
	}
}

// Permute two genes.
type Permute struct {
	// Mutation rate
	Rate float64
}

// Apply permute mutation.
func (perm Permute) apply(indi *Individual, generator *rand.Rand) {
	if generator.Float64() < perm.Rate {
		// Choose two points on the genome
		var i = generator.Intn(len(indi.Genome))
		var j = i
		// Make sure both points are different
		for i == j {
			j = generator.Intn(len(indi.Genome))
		}
		// Permute the genes
		indi.Genome[i], indi.Genome[j] = indi.Genome[j], indi.Genome[i]
	}
}
