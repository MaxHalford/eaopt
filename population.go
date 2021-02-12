package eaopt

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

// A Population contains individuals. Individuals mate within a population.
// Individuals can migrate from one population to another. Each population has a
// random number generator to bypass the global rand mutex.
type Population struct {
	Individuals Individuals   `json:"indis"`
	Age         time.Duration `json:"age"`
	Generations uint          `json:"generations"`
	ID          string        `json:"id"`
	RNG         *rand.Rand
}

// Generate a new population.
func newPopulation(size uint, parallel bool, newGenome func(rng *rand.Rand) Genome, rng *rand.Rand) Population {
	var (
		popRNG = rand.New(rand.NewSource(rng.Int63()))
		pop    = Population{
			Individuals: newIndividuals(size, parallel, newGenome, popRNG),
			ID:          randString(3, popRNG),
			RNG:         popRNG,
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

// Apply a function to a slice of Populations.
func (pops Populations) Apply(f func(pop *Population) error) error {
	var g errgroup.Group
	for i := range pops {
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return f(&pops[i])
		})
	}
	return g.Wait()
}
