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
	Individuals Individuals   `json:"indis"`
	Age         time.Duration `json:"age"`
	Generations int           `json:"generations"`
	ID          string        `json:"id"`
	rng         *rand.Rand
}

// Generate a new population.
func makePopulation(nIndis int, gm GenomeMaker, id string) Population {
	var (
		rng = makeRandomNumberGenerator()
		pop = Population{
			Individuals: makeIndividuals(nIndis, gm, rng),
			ID:          id,
			rng:         rng,
		}
	)
	return pop
}

// Log a Population's current statistics with a provided log.Logger.
func (pop Population) Log(logger *log.Logger) {
	logger.Printf(
		"id=%s min=%f max=%f avg=%f std=%f",
		pop.ID,
		pop.Individuals.FitMin(),
		pop.Individuals.FitMax(),
		pop.Individuals.FitAvg(),
		pop.Individuals.FitStd(),
	)
}

// Populations type is necessary for migration and speciation purposes.
type Populations []Population
