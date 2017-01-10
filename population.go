package gago

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
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

// Log a Population's current statistics with a provided log.Logger.
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

// Post a Population's current statistics with a provided URL.
func (pop Population) Post(url string) {
	var payload = new(bytes.Buffer)
	json.NewEncoder(payload).Encode(
		struct {
			ID  int     `json:"id"`
			Min float64 `json:"min_fitness"`
			Max float64 `json:"max_fitness"`
			Avg float64 `json:"avg_fitness"`
			Std float64 `json:"std_fitness"`
		}{
			ID:  pop.ID,
			Min: pop.Individuals.FitMin(),
			Max: pop.Individuals.FitMax(),
			Avg: pop.Individuals.FitAvg(),
			Std: pop.Individuals.FitStd(),
		},
	)
	var resp, _ = http.Post(url, "application/json; charset=utf-8", payload)
	defer resp.Body.Close()
}

// Populations type is necessary for migration and speciation purposes.
type Populations []Population
