package gago

import (
	"log"
	"math/rand"
	"time"
)

// A Population contains individuals. Individuals mate within a population.
// Individuals can migrate from one population to another. Each population has a
// random number generator to bypass the global rand mutex.
type Population struct {
	Individuals Individuals
	Age         time.Duration
	Generations int
	ID          int
	rng         *rand.Rand
}

// Generate a new population.
func makePopulation(nbrIndis int, gm GenomeMaker, id int) Population {
	var (
		rng = makeRandomNumberGenerator()
		pop = Population{
			Individuals: makeIndividuals(nbrIndis, gm, rng),
			ID:          id,
			rng:         rng,
		}
	)
	return pop
}

// Log a Population's current statistics.
func (pop Population) Log(logger *log.Logger) {
	logger.Printf(
		"id=%d min=%f max=%f avg=%f std=%f",
		pop.ID,
		pop.Individuals.FitMin(),
		pop.Individuals.FitMax(),
		pop.Individuals.FitAvg(),
		pop.Individuals.FitStd(),
	)
}

// Populations type is necessary for migration and speciation purposes.
type Populations []Population
