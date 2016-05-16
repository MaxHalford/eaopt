package gago

import (
	"math/rand"
	"time"
)

// A Population contains individuals. Individuals mate within a population. Individuals can
// migrate from one population to another.
type Population struct {
	// Individuals
	Individuals Individuals
	// Elapsed time
	Duration time.Duration
	// Each population has a random number generator to bypass the global rand mutex
	generator *rand.Rand
}

// Generate a new population.
func makePopulation(nbIndis, nbGenes int, init Initializer, generator *rand.Rand) Population {
	var pop = Population{
		Individuals: makeIndividuals(nbIndis, nbGenes),
		Duration:    0,
		generator:   generator,
	}
	// Randomly initialize each individual's genome
	for i := range pop.Individuals {
		init.apply(&pop.Individuals[i], pop.generator)
	}
	return pop
}

// Crossover replaces the GA with new individuals called offsprings. The
// method takes as arguments a selection method, a crossover method and the size
// of the crossover. The size of the crossover is the number of individuals
// whose genes will be mixed to generate an offspring with the crossover
// function.
func (pop *Population) crossover(parents Individuals, c Crossover) {
	// Create an empty slice of individuals to store the offsprings
	var offsprings = make(Individuals, len(pop.Individuals))
	// Generate offsprings through crossover until there are enough offsprings
	var counter = 0
	for counter < len(offsprings) {
		var children = c.Apply(parents, pop.generator)
		for _, child := range children {
			if counter < len(offsprings) {
				offsprings[counter] = child
			}
			counter++
		}
	}
	// Replace the old population with the new one
	copy(pop.Individuals, offsprings)
}

// Mutate each individual in a population.
func (pop *Population) mutate(m Mutator, mutRate float64) {
	for _, individual := range pop.Individuals {
		if pop.generator.Float64() < mutRate {
			m.Apply(&individual, pop.generator)
		}
	}
}
