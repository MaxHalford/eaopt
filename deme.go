package gago

import "math/rand"

// A Deme contains individuals. Individuals mate within a deme. Individuals can
// migrate from one deme to another.
type Deme struct {
	// Number of individuals in the deme
	Size int
	// Individuals
	Individuals Individuals
	// Each deme has a random number generator to bypass the global rand mutex
	Generator *rand.Rand
}

// Initialize each individual in a deme.
func (deme *Deme) initialize(indiSize int, boundary float64) {
	for i := range deme.Individuals {
		individual := Individual{make([]float64, indiSize), 0.0}
		individual.initialize(boundary, deme.Generator)
		deme.Individuals[i] = individual
	}
}

// Evaluate the fitness of each individual in a deme.
func (deme *Deme) evaluate(fitnessFunction func([]float64) float64) {
	for i := range deme.Individuals {
		deme.Individuals[i].evaluate(fitnessFunction)
	}
}

// Sort the individuals in a deme.
func (deme *Deme) sort() {
	deme.Individuals.sort()
}

// Mutate each individual in a deme.
func (deme *Deme) mutate(mutMethod func(indi *Individual, mutRate, mutIntensity float64, generator *rand.Rand),
	mutRate, mutIntensity float64) {
	for _, individual := range deme.Individuals {
		// Use the pointer to the individual to mutate
		mutMethod(&individual, mutRate, mutIntensity, deme.Generator)
	}
}

// Crossover replaces the population with new individuals called offsprings. The
// takes as arguments a selection method, a crossover method and the size of the
// crossover. The size of the crossover is the number of individuals whose genes
// will be mixed to generate an offspring with the crossover function.
func (deme *Deme) crossover(selection func(Individuals, *rand.Rand) Individual,
	crossMethod func(Individuals, *rand.Rand) Individual, crossSize int) {
	// Create an empty slice of individuals to store the offsprings
	offsprings := make(Individuals, deme.Size)
	for i := 0; i < deme.Size; i++ {
		// Select individuals to perform crossover
		selected := make(Individuals, crossSize)
		for j := 0; j < crossSize; j++ {
			selected[j] = selection(deme.Individuals, deme.Generator)
		}
		// Generate an offspring from the selected individuals
		offsprings[i] = crossMethod(selected, deme.Generator)
	}
	// Replace the population with the offsprings
	deme.Individuals = offsprings
}
