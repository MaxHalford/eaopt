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

	// Fields that have to be provided by the user
	Ff             FitnessFunction // Fitness function to evaluate individuals (depends on the problem)
	Initializer    Initializer     // Method for initializing a new individual (depends on the problem)
	MigFrequency   int             // Frequency at which migrations occur
	Migrator       Migrator
	Model          Model
	NbrClusters    int // Number of clusters each populations is split into before evolving
	NbrGenes       int // Number of genes in each individual (imposed by the problem)
	NbrIndividuals int // Initial number of individuals in each population
	NbrPopulations int // Number of populations

	// Fields that are generated at runtime
	Best        Individual // Overall best individual (dummy initialization at the beginning)
	Duration    time.Duration
	Generations int
	Populations Populations
}

// Validate the parameters of a GA to ensure it will run correctly. Some
// settings or combination of settings may be incoherent during runtime.
func (ga GA) Validate() error {
	// Check the fitness function presence
	if ga.Ff == nil {
		return errors.New("'Ff' cannot be nil")
	}
	// Check the initialization method presence
	if ga.Initializer == nil {
		return errors.New("'Initializer' cannot be nil")
	}
	// Check the migration frequency in the presence of a migrator
	if ga.Migrator != nil && ga.MigFrequency < 1 {
		return errors.New("'MigFrequency' should be strictly higher than 0")
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
	// Check the number of clusters
	if ga.NbrClusters < 0 {
		return errors.New("'NbrClusters' should be higher or equal to 1 if provided")
	}
	// Check the number of genes
	if ga.NbrGenes < 1 {
		return errors.New("'NbrGenes' should be higher or equal to 1")
	}
	// Check the number of individuals
	if ga.NbrIndividuals < 2 {
		return errors.New("'NbrIndividuals' should be higher or equal to 2")
	}
	// Check the number of populations
	if ga.NbrPopulations < 1 {
		return errors.New("'NbrPopulations' should be higher or equal to 1")
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
	ga.Populations = make([]Population, ga.NbrPopulations)
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = makePopulation(
				ga.NbrIndividuals,
				ga.NbrGenes,
				ga.Ff,
				ga.Initializer,
			)
			// Evaluate it's individuals
			ga.Populations[j].Individuals.Evaluate(ga.Ff)
			// Sort it's individuals
			ga.Populations[j].Individuals.Sort()
		}(i)
	}
	wg.Wait()
	// Best individual (dummy initialization)
	ga.Best = makeIndividual(ga.NbrGenes, rand.New(rand.NewSource(time.Now().UnixNano())))
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
	// Increment the generations counter at the beginning to not migrate at generation 0
	ga.Generations++
	// Migrate the individuals between the populations if there are enough
	// populations, there is a migrator and the migration frequency divides the
	// generation count
	if ga.NbrPopulations > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations)
	}
	// Use a wait group to enhance the populations in parallel
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Apply clustering if a number of clusters has been given
			if ga.NbrClusters > 0 {
				var clusters = ga.Populations[j].cluster(ga.NbrClusters)
				// Apply the evolution model to each cluster
				for k := range clusters {
					ga.Model.Apply(&clusters[k])
				}
				// Merge each cluster back into the original population
				ga.Populations[j].Individuals = clusters.merge()
			} else {
				// Else apply the evolution model to the entire population
				ga.Model.Apply(&ga.Populations[j])
			}
			// Evaluate and sort
			ga.Populations[j].Individuals.Evaluate(ga.Ff)
			ga.Populations[j].Individuals.Sort()
			ga.Populations[j].Duration += time.Since(start)
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Duration += time.Since(start)
}
