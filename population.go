package gago

import (
	"math/rand"
	"time"
)

// A Population contains individuals. Individuals mate within a population. Individuals can
// migrate from one population to another.
type Population struct {
	Individuals Individuals
	Duration    time.Duration
	generator   *rand.Rand // Each population has a random number generator to bypass the global rand mutex
}

// Generate a new population.
func makePopulation(nbIndis, nbGenes int, init Initializer) Population {
	var (
		source = rand.NewSource(time.Now().UnixNano())
		pop    = Population{
			Individuals: makeIndividuals(nbIndis, nbGenes),
			Duration:    0,
			generator:   rand.New(source),
		}
	)
	// Randomly initialize each individual's genome
	for i := range pop.Individuals {
		init.apply(&pop.Individuals[i], pop.generator)
	}
	return pop
}
