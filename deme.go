package gago

import "math/rand"

// A Deme contains individuals. Individuals mate within a deme. Individuals can
// migrate from one deme to another.
type Deme struct {
	// Individuals
	Individuals Individuals
	// Each deme has a random number generator to bypass the global rand mutex
	generator *rand.Rand
}

// Initialize each individual in a deme.
func (deme *Deme) initialize(NbGenes int, init Initializer) {
	for i := range deme.Individuals {
		var individual = Individual{make([]interface{}, NbGenes), 0.0}
		init.apply(&individual, deme.generator)
		deme.Individuals[i] = individual
	}
}

// Evaluate the fitness of each individual in a deme.
func (deme *Deme) evaluate(ff FitnessFunction) {
	for i := range deme.Individuals {
		deme.Individuals[i].Evaluate(ff)
	}
}

// Sort the individuals in a deme. This method is merely a convenience for
// calling the Individuals.Sort method within the deme from the population.
func (deme *Deme) sort() {
	deme.Individuals.Sort()
}

// Mutate each individual in a deme.
func (deme *Deme) mutate(m Mutator) {
	for _, individual := range deme.Individuals {
		// Use the pointer to the individual to perform mutation
		m.apply(&individual, deme.generator)
	}
}

// Breed replaces the population with new individuals called offsprings. The
// method takes as arguments a selection method, a breeding method and the size
// of the breeding. The size of the breeding is the number of individuals whose
// genes will be mixed to generate an offspring with the breeding function.
func (deme *Deme) breed(selector Selector, c Breeder) {
	// Create an empty slice of individuals to store the offsprings
	var offsprings = make(Individuals, len(deme.Individuals))
	for i := range offsprings {
		// Generate an offspring from the selected individuals
		offsprings[i] = c.apply(selector, deme.Individuals, deme.generator)
	}
	// Replace the population with the offsprings
	deme.Individuals = offsprings
}
