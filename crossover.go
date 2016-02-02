package gago

import "math/rand"

// Parenthood crossover combines two individuals (the parents) into one (the
// offspring). Each parent's contribution to the DNA is determined by the toss
// of a coin. The offspring can inherit from it's mother's genes (coin <= 0.33),
// from it's father's genes (0.33 < coin <= 0.66) or from a random mix of both
// (0.66 < coin <= 1). A coin is thrown for each gene. With this method only the
// two first selected individuals are considered, hence the CrossSize parameter
// should be set to 2.
func Parenthood(indis Individuals, generator *rand.Rand) Individual {
	var mother = indis[0]
	var father = indis[1]
	// Create an individual with an empty DNA
	var offspring = Individual{make([]float64, len(mother.Dna)), 0.0}
	// For every gene in the parent's DNA
	for i := range offspring.Dna {
		// Flip a coin and decide what to do
		var coin = rand.Float64()
		switch {
		// The offspring receives the mother's gene
		case coin <= 0.33:
			offspring.Dna[i] = mother.Dna[i]
		// The offspring receives the father's gene
		case coin <= 0.66:
			offspring.Dna[i] = father.Dna[i]
		// The offspring receives a mixture of his parent's genes
		default:
			// Random weight for each individual
			var pMother = generator.Float64()
			var pFather = 1 - pMother
			offspring.Dna[i] = pMother*mother.Dna[i] + pFather*father.Dna[i]
		}
	}
	return offspring
}

// FitnessProportional crossover combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this crossover method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring.
func FitnessProportional(indis Individuals, generator *rand.Rand) Individual {
	// Create an individual with an empty DNA
	var offspring = Individual{make([]float64, len(indis[0].Dna)), 0.0}
	// For every gene in the parent's DNA
	for i := range offspring.Dna {
		// Weight of each individual in the crossover
		var weights = make([]float64, len(indis))
		// Sum of the weights
		var total float64
		// Assign a weight to each individual
		for j := range indis {
			weights[j] = rand.Float64()
			total += weights[j]
		}
		// Create the new gene as the product of the individuals' genes
		var gene float64
		for j := range indis {
			gene += indis[j].Dna[i] * weights[j] / total
		}
		// Assign the new gene to the offspring
		offspring.Dna[i] = gene
	}
	return offspring
}
