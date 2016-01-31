package genalg

import (
	"math/rand"
	"sort"
)

// An individual represents a potential solution to a problem.
// The individual's DNA is it's genotype, which is a slice
// containing genes. Every gene is a floating point numbers.
// The Fitness is the individual's phenotype and is represented
// by a floating point number.
type Individual struct {
	Dna     []float64
	Fitness float64
}

// Initialize an individual.
func (indi *Individual) initialize(boundary float64) {
	for i := range indi.Dna {
		// Decide if positive or negative
		var sign float64
		if rand.Float64() > 0.5 {
			sign = 1.0
		} else {
			sign = -1.0
		}
		gene := rand.Float64() * sign * boundary
		indi.Dna[i] = gene
	}
}

// Mutate an individual.
func (indi *Individual) mutate(rate float64, std float64) {
	for i := range indi.Dna {
		// Flip a coin and decide to mutate or not
		coin := rand.Float64()
		if coin <= rate {
			indi.Dna[i] *= rand.NormFloat64() * std
		}
	}
}

// Crossover two individuals. Each individual can
// contribute to the DNA of the offsprings based
// on the toss of a coin. Both individuals can
// also mix their genes to create a new one.
func crossover(mother *Individual, father *Individual) Individual {

	//sum := mother.Fitness + father.Fitness
	//pMother := mother.Fitness / sum
	//pFather := father.Fitness / sum

	pMother := rand.Float64()
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

// Evaluate the Fitness of an individual.
func (indi *Individual) evaluate(FitnessFunction func([]float64) float64) {
	indi.Fitness = FitnessFunction(indi.Dna)
}

// The slice of individuals type is necessary
// for sorting and selection purposes.
type Individuals []Individual

// Sort the individuals of a deme in ascending order based
// on their Fitness. The convention is that we always want
// to minimize a function. If the function has to be maximized
// then we can minimize 1/func(X) or -func(X).
func (individuals Individuals) sort() {
	sort.Sort(individuals)
}

func (individuals Individuals) Len() int {
	return len(individuals)
}

func (individuals Individuals) Less(i, j int) bool {
	return individuals[i].Fitness < individuals[j].Fitness
}

func (individuals Individuals) Swap(i, j int) {
	individuals[i], individuals[j] = individuals[j], individuals[i]
}
