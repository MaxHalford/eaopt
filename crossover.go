package gago

import "math/rand"

// Crossover mixes two or more individuals into a new individual called the
// offspring.
type Crossover interface {
	Apply(selector Selector, individuals Individuals, generator *rand.Rand) Individual
}

// Parenthood crossover combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the toss of a coin. The offspring can inherit from it's mother's genes
// (coin <= 0.33), from it's father's genes (0.33 < coin <= 0.66) or from a
// random mix of both (0.66 < coin <= 1). A coin is thrown for each gene. With
// this method only the two first selected individuals are considered, hence the
// CrossSize parameter should be set to 2.
type Parenthood struct{}

// Apply parenthood crossover.
func (ph Parenthood) Apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var mother = s.Apply(indis, generator)
	var father = s.Apply(indis, generator)
	// Create an individual with an empty Genome
	var offspring = Individual{make([]float64, len(mother.Genome)), 0.0}
	// For every gene in the parent's Genome
	for i := range offspring.Genome {
		// Flip a coin and decide what to do
		var coin = rand.Float64()
		switch {
		// The offspring receives the mother's gene
		case coin <= 0.33:
			offspring.Genome[i] = mother.Genome[i]
		// The offspring receives the father's gene
		case coin <= 0.66:
			offspring.Genome[i] = father.Genome[i]
		// The offspring receives a mixture of his parent's genes
		default:
			// Random weight for each individual
			var pMother = generator.Float64()
			var pFather = 1 - pMother
			offspring.Genome[i] = pMother*mother.Genome[i] + pFather*father.Genome[i]
		}
	}
	return offspring
}

// FitnessProportionate crossover combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this crossover method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring.
type FitnessProportionate struct {
	// Should be any integer above or equal to two
	NbIndividuals int
}

// Apply fitness proportional crossover.
func (fpc FitnessProportionate) Apply(indis Individuals, generator *rand.Rand) Individual {
	// Create an individual with an empty Genome
	var offspring = Individual{make([]float64, len(indis[0].Genome)), 0.0}
	// For every gene in the parent's Genome
	for i := range offspring.Genome {
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
			gene += indis[j].Genome[i] * weights[j] / total
		}
		// Assign the new gene to the offspring
		offspring.Genome[i] = gene
	}
	return offspring
}
