package gago

import "math/rand"

// The Initializer is here to create the first generation of individuals in a
// deme. It applies to an individual level and instantiates it's genome gene by
// gene.
type Initializer interface {
	apply(individual *Individual, generator *rand.Rand)
}

// UniformFloat generates random floating point values between given boundaries.
type UniformFloat struct {
	lower, upper float64
}

// Apply the uniform floating point initializer.
func (uf UniformFloat) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Decide if positive or negative
		var gene float64
		if generator.Float64() < 0.5 {
			gene = generator.Float64() * uf.lower
		} else {
			gene = generator.Float64() * uf.upper
		}
		indi.Genome[i] = gene
	}
}

// UniformString generates random string values from a given corpus.
type UniformString struct {
	Corpus []string
}

// Apply the uniform floating point initializer.
func (us UniformString) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = us.Corpus[generator.Intn(len(us.Corpus))]
	}
}
