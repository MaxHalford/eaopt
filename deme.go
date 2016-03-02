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

// Sort the individuals in a deme. This method is merely and ergonomic layer.
func (deme *Deme) sort() {
	deme.Individuals.Sort()
}

// Mutate each individual in a deme.
func (deme *Deme) mutate(m Mutator) {
	for _, individual := range deme.Individuals {
		m.apply(&individual, deme.generator)
	}
}

// Crossover replaces the population with new individuals called offsprings. The
// method takes as arguments a selection method, a crossover method and the size
// of the crossover. The size of the crossover is the number of individuals
// whose genes will be mixed to generate an offspring with the crossover
// function.
func (deme *Deme) crossover(selector Selector, c Crossover) {
	// Create an empty slice of individuals to store the offsprings
	var offsprings = make(Individuals, len(deme.Individuals))
	for i := range offsprings {
		// Generate an offspring from the selected individuals
		offsprings[i] = c.apply(selector, deme.Individuals, deme.generator)
	}
	// Replace the population with the offsprings
	deme.Individuals = offsprings
}
