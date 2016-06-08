package gago

import (
	"errors"
	"log"
	"sync"
	"time"
)

// A GA contains population which themselves contain individuals.
type GA struct {
	// Number of populations
	NbPopulations int
	// Number of individuals in each population
	NbIndividuals int
	// Number of genes in each individual (imposed by the problem)
	NbGenes int
	// Number of parents selected for reproduction
	NbParents int
	// Populations
	Populations []Population
	// Overall best individual (dummy initialization at the beginning)
	Best Individual
	// Fitness function to evaluate individuals (imposed by the problem)
	Ff FitnessFunction
	// Initial random boundaries
	Initializer Initializer
	// Selection method
	Selector Selector
	// Crossover method
	Crossover Crossover
	// Mutation method
	Mutator Mutator
	// Mutation rate
	MutRate float64
	// Migration method
	Migrator Migrator
	// Migration frequency
	MigFrequency int
	// Clustering method
	Clusterer Clusterer
	// Number of clusters
	NbClusters int
	// Number of generations
	Generations int
	// Elapsed time
	Duration time.Duration
}

// Validate the parameters of a GA to ensure it will run correctly. Some
// settings or combination of settings may be incoherent during runtime.
func (ga *GA) Validate() error {
	var err error
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
	// Check the number of parents
	if ga.NbParents < 0 || ga.NbParents > ga.NbIndividuals {
		return errors.New(`'NbParents' should be higher or equal to 0 and should
			be lower or equal to 'NbIndividuals'`)
	}
	// Check the fitness function
	if ga.Ff == nil {
		return errors.New("'Ff' cannot be nil")
	}
	// Check the initialization method
	if ga.Initializer == nil {
		return errors.New("'Initializer' cannot be nil")
	}
	// Check the selection method
	if ga.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the crossover method
	if ga.Crossover == nil {
		return errors.New("'Crossover' cannot be nil")
	}
	// Check the mutation rate
	if ga.MutRate < 0 || ga.MutRate > 1 {
		return errors.New("'MutRate' should be comprised between 0 and 1 (included)")
	}
	// Check the migration frequency in the presence of a migrator
	if ga.Migrator != nil && ga.MigFrequency < 1 {
		return errors.New("'MigFrequency' should be strictly higher than 0")
	}
	// No error
	return err
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population. Running Initialize after running Enhance will
// reset the GA entirely.
func (ga *GA) Initialize() {
	// Begin by validating the parameters of the GA
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
	ga.Best = makeIndividual(ga.NbGenes)
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
			// 1. Select
			var parents = ga.Selector.Apply(ga.NbParents, ga.Populations[j].Individuals,
				ga.Populations[j].generator)
			// 2. Crossover
			ga.Populations[j].crossover(parents, ga.Crossover)
			// 3. Mutate
			if ga.Mutator != nil {
				ga.Populations[j].mutate(ga.Mutator, ga.MutRate)
			}
			// 4. Evaluate
			ga.Populations[j].Individuals.evaluate(ga.Ff)
			// 5. Sort
			ga.Populations[j].Individuals.sort()
			ga.Populations[j].Duration += time.Since(start)
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Duration += time.Since(start)
}
