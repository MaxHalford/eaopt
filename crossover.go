package gago

import (
	"math/rand"
	"sort"
)

// Crossover mixes two or more individuals into a new individual called the
// offspring.
type Crossover interface {
	apply(selector Selector, individuals Individuals, generator *rand.Rand) Individual
}

// CPoint selects identical random points on each parent's genome and exchanges
// mirroring segments. It generalizes one-point crossover and two-point
// crossover to n-point crossover. One of the generated offsprings is returned
// at random.
type CPoint struct {
	NbPoints int
}

// Apply n-point crossover.
func (cp CPoint) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var parents = s.apply(2, indis, generator)
	// Choose n random points along the genome
	var points = generator.Perm(len(parents[0].Genome))[:cp.NbPoints]
	// Sort the points
	sort.Ints(points)
	// Add the start and end of the genome points
	points = append([]int{0}, points...)
	points = append(points, len(parents[0].Genome))
	// Create an offspring
	var offspring = Individual{make([]interface{}, len(parents[0].Genome)), 0.0}
	// Counter for choosing on which parent to choose the segment from (0 or 1)
	var p = 0
	for i := 0; i < len(points)-1; i++ {
		// Copy the parent's segment onto the offspring's segment
		var segment = parents[p].Genome[points[i]:points[i+1]]
		copy(offspring.Genome[points[i]:points[i+1]], segment)
		// p will alternate between 0 and 1
		p = (p - 1) * (p - 1)
	}
	return offspring
}

// CFUniform crossover combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the value of a probability p. The offspring can inherit from it's mother's
// genes (p <= 0.33), from it's father's genes (0.33 < p <= 0.66) or from a
// random mix of both (0.66 < p <= 1). A coin is thrown for each gene. Only
// works for floating point values.
type CFUniform struct{}

// Apply uniform float crossover.
func (cfu CFUniform) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var parents = s.apply(2, indis, generator)
	var mother = parents[0]
	var father = parents[1]
	// Create an offspring
	var offspring = Individual{make([]interface{}, len(mother.Genome)), 0.0}
	// For every gene in the parent's genome
	for i := range offspring.Genome {
		// Pick a random number and decide what to do
		var p = rand.Float64()
		switch {
		// The offspring receives the mother's gene
		case p <= 0.33:
			offspring.Genome[i] = mother.Genome[i]
			// The offspring receives the father's gene
		case p <= 0.66:
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

// CFProportionate crossover combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this crossover method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring. Only works for floating point values.
type CFProportionate struct {
	// Should be any integer above or equal to two
	NbIndividuals int
}

// Apply proportionate float crossover.
func (cfp CFProportionate) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose individuals at random
	var parents = s.apply(cfp.NbIndividuals, indis, generator)
	// Create an offspring
	var offspring = Individual{make([]interface{}, len(indis[0].Genome)), 0.0}
	// For every gene in the parent's genome
	for i := range offspring.Genome {
		// Weight of each individual in the crossover
		var weights = generateWeights(len(indis))
		// Create the new gene as the product of the individuals' genes
		var gene float64
		for j := range parents {
			gene += parents[j].Genome[i].(float64) * weights[j]
		}
		// Assign the new gene to the offspring
		offspring.Genome[i] = gene
	}
	return offspring
}

// CPMX (Partially Mapped Crossover) randomly picks a crossover point. The
// offspring is built by copying one of the parents and then copying the other
// parent's values up to the crossover point. Each gene that is replaced is
// permuted with the gene that is copied in the first parent's genome. Two
// offsprings are generated in such a way (because there are two parents).
// This crossover method ensures the offspring's genomes are composed of unique
// genes, which is particularly useful for permutation problems such as the
// Traveling Salesman Problem (TSP).
type CPMX struct{}

// Apply partially mixed crossover.
func (pmx CPMX) apply(s Selector, indis Individuals, generator *rand.Rand) Individual {
	// Choose two individuals at random
	var parents = s.apply(2, indis, generator)
	var mother = parents[0]
	var father = parents[1]
	// Create an offspring with the mother's genome
	var offspring = Individual{make([]interface{}, len(mother.Genome)), 0.0}
	copy(offspring.Genome, mother.Genome)
	// Choose a random crossover point
	var point = generator.Intn(len(mother.Genome))
	// Paste the father's genome up to the crossover point
	for i := 0; i < point; i++ {
		// Find where the element is in the offspring's genome
		var position = getIndex(father.Genome[i], offspring.Genome)
		// Swap the genes
		offspring.Genome[position], offspring.Genome[i] = offspring.Genome[i], father.Genome[i]
	}
	return offspring
}
