package gago

import "math/rand"

// Parenthood crossover mixes two individuals (the parents) into one (the
// offspring). Each parent's contribution to the DNA is determined by the toss
// of a coin. The offspring can inherit from it's mother's genes (coin <= 0.33),
// from it's father's genes (0.33 < coin <= 0.66) or from a random mix of both
// (0.66 < coin <= 1). With this method only the two first selected individuals
// are considered, hence the CSize parameter should be set to 2.
func parenthood(indis Individuals, generator *rand.Rand) Individual {
	mother := indis[0]
	father := indis[1]
	// Create an individual with an empty DNA
	offspring := Individual{make([]float64, len(mother.Dna)), 0.0}
	// For every gene in the parent's DNA
	for j := range mother.Dna {
		// Flip a coin and decide what to do
		coin := rand.Float64()
		switch {
		// The offspring receives the mother's gene
		case coin <= 0.33:
			offspring.Dna[j] = mother.Dna[j]
		// The offspring receives the father's gene
		case coin <= 0.66:
			offspring.Dna[j] = father.Dna[j]
		// The offspring receives a mixture of his parent's genes
		default:
			// Random weight for each individual
			pMother := generator.Float64()
			pFather := 1 - pMother
			offspring.Dna[j] = pMother*mother.Dna[j] + pFather*father.Dna[j]
		}
	}
	return offspring
}
