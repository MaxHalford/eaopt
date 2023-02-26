package eaopt

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

// A Population contains individuals. Individuals mate within a population.
// Individuals can migrate from one population to another. Each population has a
// random number generator to bypass the global rand mutex.
type Population struct {
	Individuals     Individuals                  `json:"indis"`
	Age             time.Duration                `json:"age"`
	Generations     uint                         `json:"generations"`
	ID              string                       `json:"id"`
	RNG             *rand.Rand                   `json:"-"`
	JSONUnmarshaler func([]byte) (Genome, error) `json:"-"`
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

func newPopulationsFromBytes(populationCount uint, b []byte, RNG *rand.Rand, unmarshaler func([]byte) (Genome, error)) ([]Population, error) {
	pops := make([]Population, populationCount)
	for i := range pops {
		pops[i].RNG = rand.New(rand.NewSource(RNG.Int63()))
		pops[i].JSONUnmarshaler = unmarshaler
	}
	err := json.Unmarshal(b, &pops)
	return pops, err
}

// Log a Population's current statistics with a provided log.Logger.
func (pop Population) Log(logger *log.Logger) {
	logger.Print(pop.stats())
}

func (pop Population) String() string {
	return pop.stats()
}

func (pop Population) stats() string {
	return fmt.Sprintf("pop_id=%s min=%f max=%f avg=%f std=%f",
		pop.ID,
		pop.Individuals.FitMin(),
		pop.Individuals.FitMax(),
		pop.Individuals.FitAvg(),
		pop.Individuals.FitStd(),
	)
}

// UnmarshalJSON implements a JSON unmarshaler for Populations. This override
// is required because the JSON unmarshaler has no idea how to create your
// implementation of the Genome interface. See setup_test.go:VectorJSONUnmarshaler
// for an example JSON unmarshaler function. See ga_test.go:TestMarshalGA for
// an example of using this custom decoder in a GA instance.
func (pop *Population) UnmarshalJSON(data []byte) error {

	var decoded struct {
		Age         time.Duration
		Generations uint
		ID          string
		Indis       []interface{}
	}
	err := json.Unmarshal(data, &decoded)
	if err != nil {
		return err
	}

	pop.Age = decoded.Age
	pop.Generations = decoded.Generations
	pop.ID = decoded.ID
	if pop.JSONUnmarshaler != nil {
		for _, v := range decoded.Indis {
			val, err := json.Marshal(v.(map[string]interface{})["genome"])
			if err != nil {
				return err
			}
			genome, err := pop.JSONUnmarshaler(val)
			if err != nil {
				return err
			}
			pop.Individuals = append(pop.Individuals, Individual{
				Genome:  genome,
				Fitness: v.(map[string]interface{})["fitness"].(float64),
				ID:      v.(map[string]interface{})["id"].(string),
			})
		}
	}
	return nil
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
