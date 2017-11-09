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
func newPopulation(size int, newGenome NewGenome, rng *rand.Rand) Population {
	var (
		popRNG = rand.New(rand.NewSource(rng.Int63()))
		pop    = Population{
			Individuals: newIndividuals(size, newGenome, popRNG),
			ID:          randString(3, popRNG),
			rng:         popRNG,
		}
	)
	return pop
}

// Log a Population's current statistics with a provided log.Logger.
func (pop Population) Log(logger *log.Logger) {
	logger.Printf(
		"pop_id=%s min=%f max=%f avg=%f std=%f",
		pop.ID,
		pop.Individuals.FitMin(),
		pop.Individuals.FitMax(),
		pop.Individuals.FitAvg(),
		pop.Individuals.FitStd(),
	)
}

// Populations type is necessary for migration and speciation purposes.
type Populations []Population
