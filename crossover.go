package gago

import (
	"math/rand"
	"sort"
)

// Crossover generates new individuals called "offsprings" are by mixing the
// genomes of a sample group of individuals. Instead of using the same sample
// for generating each offspring, each crossover resamples the population in
// order to preserve diversity.
type Crossover interface {
	Apply(indis Individuals, generator *rand.Rand) Individuals
}

// CrossPoint selects identical random points on each parent's genome and
// exchanges mirroring segments. It generalizes one-point crossover and
// two-point crossover to n-point crossover.
type CrossPoint struct {
	NbPoints int
}

// Apply n-point crossover.
func (cross CrossPoint) Apply(indis Individuals, generator *rand.Rand) Individuals {
	var (
		// Choose two individuals at random
		parents = indis.sample(2, generator)
		// Choose n random points along the genome
		points = generator.Perm(len(parents[0].Genome))[:cross.NbPoints]
	)
	// Sort the points
	sort.Ints(points)
	// Add the start and end of the genome points
	points = append([]int{0}, points...)
	points = append(points, len(parents[0].Genome))
	// Create offsprings
	var (
		nbGenes    = len(parents[0].Genome)
		offsprings = makeIndividuals(len(parents), nbGenes)
		// Use switching indexes to know which parent's genome to copy
		a = 0
		b = 1
	)
	for i := 0; i < len(points)-1; i++ {
		// Copy the first parent's segment onto the first offspring's segment
		copy(
			offsprings[0].Genome[points[i]:points[i+1]],
			parents[a].Genome[points[i]:points[i+1]],
		)
		// Copy the second parent's segment onto the second offspring's segment
		copy(
			offsprings[1].Genome[points[i]:points[i+1]],
			parents[b].Genome[points[i]:points[i+1]],
		)
		// Alternate for the new copying
		a, b = b, a
	}
	return offsprings
}

// CrossUniformF crossover combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the value of a probability p. Each offspring receives a proportion of both of
// it's parents genomes. The new values are located in the hyper-rectangle
// defined between both parent's position in Cartesian space.
type CrossUniformF struct{}

// Apply uniform float crossover.
func (cross CrossUniformF) Apply(indis Individuals, generator *rand.Rand) Individuals {
	var (
		parents    = indis.sample(2, generator)
		mother     = parents[0]
		father     = parents[1]
		nbGenes    = len(mother.Genome)
		offsprings = makeIndividuals(len(parents), nbGenes)
	)
	// For every gene
	for i := 0; i < nbGenes; i++ {
		// Pick a random number between 0 and 1
		var p = generator.Float64()
		offsprings[0].Genome[i] = p*mother.Genome[i].(float64) + (1-p)*father.Genome[i].(float64)
		offsprings[1].Genome[i] = (1-p)*mother.Genome[i].(float64) + p*father.Genome[i].(float64)
	}
	return offsprings
}

// CrossProportionateF crossover combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this crossover method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring. Only works for floating point values.
type CrossProportionateF struct {
	// Should be any integer above or equal to two
	NbParents int
}

// Apply proportionate float crossover.
func (cross CrossProportionateF) Apply(indis Individuals, generator *rand.Rand) Individuals {
	var (
		parents   = indis.sample(cross.NbParents, generator)
		nbGenes   = len(parents[0].Genome)
		offspring = makeIndividual(nbGenes)
	)
	// For every gene in the parent's genome
	for i := range offspring.Genome {
		var (
			// Weight of each individual in the crossover
			weights = generateWeights(len(indis))
			// Create the new gene as the product of the individuals' genes
			gene float64
		)
		for j := range parents {
			gene += parents[j].Genome[i].(float64) * weights[j]
		}
		// Assign the new gene to the offspring
		offspring.Genome[i] = gene
	}
	return Individuals{offspring}
}

// CrossPMX (Partially Mapped Crossover) randomly picks a crossover point. The
// offsprings are generated by copying one of the parents and then copying the
// other parent's values up to the crossover point. Each gene that is replaced
// is permuted with the gene that is copied in the first parent's genome. Two
// offsprings are generated in such a way (because there are two parents). This
// crossover method ensures the offspring's genomes are composed of unique
// genes, which is particularly useful for permutation problems such as the
// Traveling Salesman Problem (TSP).
type CrossPMX struct{}

// Apply partially mixed crossover.
func (c CrossPMX) Apply(indis Individuals, generator *rand.Rand) Individuals {
	var (
		parents    = indis.sample(2, generator)
		nbGenes    = len(parents[0].Genome)
		offsprings = makeIndividuals(len(parents), nbGenes)
	)
	copy(offsprings[0].Genome, parents[0].Genome)
	copy(offsprings[1].Genome, parents[1].Genome)
	// Choose a random crossover point p such that 0 < p < (nbGenes - 1)
	var (
		p = generator.Intn(nbGenes-2) + 1
		a int
		b int
	)
	// Paste the father's genome up to the crossover point
	for i := 0; i < p; i++ {
		// Find where the second parent's gene is in the first offspring's genome
		a = getIndex(parents[1].Genome[i], offsprings[0].Genome)
		// Swap the genes
		offsprings[0].Genome[a], offsprings[0].Genome[i] = offsprings[0].Genome[i], parents[1].Genome[i]
		// Find where the first parent's gene is in the second offspring's genome
		b = getIndex(parents[0].Genome[i], offsprings[1].Genome)
		// Swap the genes
		offsprings[1].Genome[b], offsprings[1].Genome[i] = offsprings[1].Genome[i], parents[0].Genome[i]
	}
	return offsprings
}
