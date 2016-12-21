package gago

import (
	"math/rand"
	"time"
)

// A Population contains individuals. Individuals mate within a population.
// Individuals can migrate from one population to another. Each population has a
// random number generator to bypass the global rand mutex.
type Population struct {
	Individuals Individuals
	Duration    time.Duration
	rng         *rand.Rand
}

// Generate a new population.
func makePopulation(nbrIndis int, gm GenomeMaker) Population {
	var (
		rng = makeRandomNumberGenerator()
		pop = Population{
			Individuals: makeIndividuals(nbrIndis, gm, rng),
			rng:         rng,
		}
	)
	return pop
}

// Populations type is necessary for migration and speciation purposes.
type Populations []Population
