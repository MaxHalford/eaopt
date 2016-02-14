package gago

import "math/rand"

// A Deme contains individuals. Individuals mate within a deme. Individuals can
// migrate from one deme to another.
type Deme struct {
	// Number of individuals in the deme, it is defined for convenience
	Size int
	// Individuals
	Individuals Individuals
	// Each deme has a random number generator to bypass the global rand mutex
	Generator *rand.Rand
}

// Initialize each individual in a deme.
func (deme *Deme) Initialize(NbGenes int, init Initializer) {
	for i := range deme.Individuals {
		var individual = Individual{make([]float64, NbGenes), 0.0}
		init.Apply(&individual, deme.Generator)
		deme.Individuals[i] = individual
	}
}

// Evaluate the fitness of each individual in a deme.
func (deme *Deme) Evaluate(ff func([]float64) float64) {
	for i := range deme.Individuals {
		deme.Individuals[i].Evaluate(ff)
	}
}

// Sort the individuals in a deme. This method is merely a convenience for
// calling the Individuals.Sort method within the deme from the population.
func (deme *Deme) Sort() {
	deme.Individuals.Sort()
}

// Mutate each individual in a deme.
func (deme *Deme) Mutate(m Mutator) {
	for _, individual := range deme.Individuals {
		// Use the pointer to the individual to perform mutation
		m.Apply(&individual, deme.Generator)
	}
}

// Crossover replaces the population with new individuals called offsprings. The
// takes as arguments a selection method, a crossover method and the size of the
// crossover. The size of the crossover is the number of individuals whose genes
// will be mixed to generate an offspring with the crossover function.
func (deme *Deme) Crossover(selector Selector, c Crossover) {
	// Create an empty slice of individuals to store the offsprings
	var offsprings = make(Individuals, deme.Size)
	for i := 0; i < deme.Size; i++ {
		// Generate an offspring from the selected individuals
		offsprings[i] = c.Apply(selector, deme.Individuals, deme.Generator)
	}
	// Replace the population with the offsprings
	deme.Individuals = offsprings
}
