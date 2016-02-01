package gago

import "math/rand"

// Crossover two individuals. Each individual can contribute to the DNA of the
// offsprings based on the toss of a coin. Both individuals can also mix their
// genes to create a new one.
func crossover(individuals Individuals, generator *rand.Rand) Individual {
	mother := individuals[0]
	father := individuals[1]
	// Random weight for each individual
	pMother := generator.Float64()
	pFather := 1 - pMother
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
			offspring.Dna[j] = pMother*mother.Dna[j] + pFather*father.Dna[j]
		}
	}
	return offspring
}
