package gago

import "math/rand"

// Breeder mixes two or more individuals into a new individual called the
// offspring.
type Breeder interface {
	apply(selector Selector, individuals Individuals, generator *rand.Rand) Individual
}

// Crossover breeding creates a new genome by splicing the genomes of two
// selected individuals (the mother and the father) and gluing them back
// together in a random way.
type Crossover struct{}

// Apply crossover breeding.
func (cv Crossover) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var mother = s.apply(indis, generator)
	var father = s.apply(indis, generator)
	// Create an individual with an empty genome
	var offspring = Individual{make([]interface{}, len(mother.Genome)), 0.0}
	// Choose where to split the genomes
	var split = rand.Intn(len(mother.Genome))
	// Split and glue
	if generator.Float64() < 0.5 {
		offspring.Genome = append(mother.Genome[:split], father.Genome[split:]...)
	} else {
		offspring.Genome = append(father.Genome[:split], mother.Genome[split:]...)
	}
	return offspring
}

// Parenthood breeding combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the toss of a coin. The offspring can inherit from it's mother's genes
// (coin <= 0.33), from it's father's genes (0.33 < coin <= 0.66) or from a
// random mix of both (0.66 < coin <= 1). A coin is thrown for each gene. With
// this method only the two first selected individuals are considered, hence the
// CrossSize parameter should be set to 2.
type Parenthood struct{}

// Apply parenthood breeding.
func (ph Parenthood) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var mother = s.apply(indis, generator)
	var father = s.apply(indis, generator)
	// Create an individual with an empty genome
	var offspring = Individual{make([]interface{}, len(mother.Genome)), 0.0}
	// For every gene in the parent's genome
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
			offspring.Genome[i] = pMother*mother.Genome[i].(float64) + pFather*father.Genome[i].(float64)
		}
	}
	return offspring
}

// FitnessProportionate breeding combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this breeding method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring.
type FitnessProportionate struct {
	// Should be any integer above or equal to two
	NbIndividuals int
}

// Apply fitness proportional breeding.
func (fpc FitnessProportionate) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose individuals at random
	var selected = make(Individuals, fpc.NbIndividuals)
	for i := range selected {
		selected[i] = s.apply(indis, generator)
	}
	// Create an individual with an empty genome
	var offspring = Individual{make([]interface{}, len(indis[0].Genome)), 0.0}
	// For every gene in the parent's genome
	for i := range offspring.Genome {
		// Weight of each individual in the breeding
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
			gene += indis[j].Genome[i].(float64) * weights[j] / total
		}
		// Assign the new gene to the offspring
		offspring.Genome[i] = gene
	}
	return offspring
}
