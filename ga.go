package gago

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

// A GA contains population which themselves contain individuals.
type GA struct {
	NbPopulations int // Number of populations
	NbIndividuals int // Initial number of individuals in each population
	NbGenes       int // Number of genes in each individual (imposed by the problem)
	Populations   []Population
	Best          Individual      // Overall best individual (dummy initialization at the beginning)
	Ff            FitnessFunction // Fitness function to evaluate individuals (imposed by the problem)
	Initializer   Initializer
	Model         Model
	Migrator      Migrator
	MigFrequency  int // Migration frequency
	Generations   int
	Duration      time.Duration
}

// Validate the parameters of a GA to ensure it will run correctly. Some
// settings or combination of settings may be incoherent during runtime.
func (ga *GA) Validate() error {
	// Check the number of populations
	if ga.NbPopulations < 1 {
		return errors.New("'NbPopulations' should be higher or equal to 1")
	}
	// Check the number of individuals
	if ga.NbIndividuals < 2 {
		return errors.New("'NbIndividuals' should be higher or equal to 2")
	}
	// Check the number of genes
	if ga.NbGenes < 1 {
		return errors.New("'NbGenes' should be higher or equal to 1")
	}
	// Check the fitness function presence
	if ga.Ff == nil {
		return errors.New("'Ff' cannot be nil")
	}
	// Check the initialization method presence
	if ga.Initializer == nil {
		return errors.New("'Initializer' cannot be nil")
	}
	// Check the model presence
	if ga.Model == nil {
		return errors.New("'Model' cannot be nil")
	}
	// Check the model is valid
	var modelErr = ga.Model.Validate()
	if modelErr != nil {
		return modelErr
	}
	// Check the migration frequency in the presence of a migrator
	if ga.Migrator != nil && ga.MigFrequency < 1 {
		return errors.New("'MigFrequency' should be strictly higher than 0")
	}
	// No error
	return nil
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population. Running Initialize after running Enhance will
// reset the GA entirely.
func (ga *GA) Initialize() {
	// Validate the parameters of the GA
	var err = ga.Validate()
	if err != nil {
		log.Fatal(err)
	}
	// Reset the number of generations and the elapsed duration
	ga.Generations = 0
	ga.Duration = 0
	// Create the populations
	ga.Populations = make([]Population, ga.NbPopulations)
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = makePopulation(ga.NbIndividuals, ga.NbGenes, ga.Initializer)
			// Evaluate it's individuals
			ga.Populations[j].Individuals.evaluate(ga.Ff)
			// Sort it's individuals
			ga.Populations[j].Individuals.sort()
		}(i)
	}
	wg.Wait()
	// Best individual (dummy initialization)
	ga.Best = makeIndividual(ga.NbGenes, rand.New(rand.NewSource(time.Now().UnixNano())))
	// Find the best individual
	ga.findBest()
}

// Find the best individual in each population and then compare the best overall
// individual to the current best individual.
func (ga *GA) findBest() {
	for _, pop := range ga.Populations {
		var best = pop.Individuals[0]
		if best.Fitness < ga.Best.Fitness {
			ga.Best = best
		}
	}
}

// Enhance each population in the GA. The population level operations are done
// in parallel with a wait group. After all the population operations have been
// run, the GA level operations are run.
func (ga *GA) Enhance() {
	var start = time.Now()
	ga.Generations++
	// Migrate the individuals between the populations if there are enough
	// populations, there is a migrator and the migration frequency divides the
	// generation count
	if ga.NbPopulations > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations)
	}
	// Use a wait group to enhance the populations in parallel
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Apply the evolution model
			ga.Model.Apply(&ga.Populations[j])
			// Evaluate and sort
			ga.Populations[j].Individuals.evaluate(ga.Ff)
			ga.Populations[j].Individuals.sort()
			ga.Populations[j].Duration += time.Since(start)
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Duration += time.Since(start)
}
