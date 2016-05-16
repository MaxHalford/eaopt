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
	// Number of generations
	Generations int
	// Elapsed time
	Duration time.Duration
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
	for _, pop := range ga.Populations {
		var best = pop.Individuals[0]
		if best.Fitness < ga.Best.Fitness {
			ga.Best = best
		}
	}
}

// Enhance each population in the GA. The population level operations are done in
// parallel with a wait group. After all the population operations have been run, the
// GA level operations are run.
func (ga *GA) Enhance() {
	var start = time.Now()
	ga.Generations++
	// Migrate the individuals between the Populations
	if ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations)
	}
	// Use a wait group to run the genetic operators in each population in parallel
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
