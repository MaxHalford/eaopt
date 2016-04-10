package gago

import (
	"math/rand"
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
	// Mutation methods
	Mutators []Mutator
	// Migration method
	Migrator Migrator
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population.
func (ga *GA) Initialize() {
	// Create the Populations
	ga.Populations = make([]Population, ga.NbPopulations)
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Create a new random number generator
			var (
				source    = rand.NewSource(time.Now().UnixNano())
				generator = rand.New(source)
			)
			// Generate a population
			ga.Populations[j] = makePopulation(ga.NbIndividuals, ga.NbGenes,
				ga.Initializer, generator)
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
	// Get each population's best individual
	var best = make(Individuals, ga.NbPopulations)
	for i, population := range ga.Populations {
		best[i] = population.Individuals[0]
	}
	// Sort the best individuals
	best.sort()
	// Get the overall best individual
	var overallBest = best[0]
	// Compare it to the current best individual
	if overallBest.Fitness < ga.Best.Fitness {
		ga.Best = overallBest
	}
}

// Migrate allows Populations to exchange individuals through the migration protocol
// defined in ga.MigMethod. This is a convenience method for calling purposes.
func (ga *GA) migrate() {
	// Use the pointer to the Populations to perform migration
	ga.Migrator.apply(ga.Populations)
}

// Enhance each population in the GA. The population level operations are done in
// parallel with a wait group. After all the population operations have been run, the
// GA level operations are run.
func (ga *GA) Enhance() {
	// Migrate the individuals between the Populations
	if ga.Migrator != nil {
		ga.migrate()
	}
	// Use a wait group to run the genetic operators in each population in parallel
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// 1. Crossover
			if ga.Crossover != nil {
				ga.Populations[j].crossover(ga.Selector, ga.Crossover)
			}
			// 2. Mutate
			if ga.Mutators != nil {
				ga.Populations[j].mutate(ga.Mutators)
			}
			// 3. Evaluate
			ga.Populations[j].Individuals.evaluate(ga.Ff)
			// 4. Sort
			ga.Populations[j].Individuals.sort()
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
}
