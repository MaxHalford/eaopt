package gago

import "math/rand"

// The Initializer is here to create the first generation of individuals in a
// deme. It applies to an individual level and instantiates it's genome gene by
// gene.
type Initializer interface {
	Apply(individual *Individual, generator *rand.Rand)
}

// UniformFloat generates random floating point values between given boundaries.
type UniformFloat struct {
	Lower, Upper float64
}

// Apply the uniform floating point initializer.
func (uf UniformFloat) Apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Decide if positive or negative
		var gene float64
		if generator.Float64() < 0.5 {
			gene = generator.Float64() * uf.Lower
		} else {
			gene = generator.Float64() * uf.Upper
		}
		indi.Genome[i] = gene
	}
}
