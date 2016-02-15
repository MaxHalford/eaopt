package gago

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// A Population contains deme which themselves contain individuals.
type Population struct {
	// Number of demes
	NbDemes int
	// Number of individuals in each deme
	NbIndividuals int
	// Number of genes in each individual (defined by the user)
	NbGenes int
	// Demes
	Demes []Deme
	// Overall best individual (dummy initialization at the begining)
	Best Individual
	// Fitness function to evaluate individuals (defined by the user)
	Ff FitnessFunction
	// Initial random boundaries
	Initializer Initializer
	// Selection method
	Selector Selector
	// Breeding method
	Breeder Breeder
	// Mutation method
	Mutator Mutator
	// Migration method
	Migrator Migrator
}

// Initialize each deme in the population and assign an initial fitness to each
// individual in each deme.
func (pop *Population) Initialize(variables int) {
	// Number of genes in each individual
	pop.NbGenes = variables
	// Create the demes
	pop.Demes = make([]Deme, pop.NbDemes)
	// Best individual (dummy instantiation)
	pop.Best = Individual{make([]interface{}, pop.NbGenes), math.Inf(1)}
	var wg sync.WaitGroup
	for i := range pop.Demes {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Create a new random number generator
			var source = rand.NewSource(time.Now().UnixNano())
			var generator = rand.New(source)
			// Create the deme
			var deme = Deme{make([]Individual, pop.NbIndividuals), generator}
			// Initialize the deme
			deme.initialize(pop.NbGenes, pop.Initializer)
			// Add it to the population
			pop.Demes[j] = deme
			// Initial evaluation
			pop.Demes[j].evaluate(pop.Ff)
			// Sort the deme
			pop.Demes[j].sort()
		}(i)
	}
	wg.Wait()
	// Find the best individual
	pop.findBest()
}

// Find the best individual in a deme and check if it's better than the current
// best individual. The deme's best individual is the first one if the deme has
// been sorted.
func (pop *Population) findBest() {
	for _, deme := range pop.Demes {
		if deme.Individuals[0].Fitness < pop.Best.Fitness {
			pop.Best = deme.Individuals[0]
		}
	}
}

// Migrate allows demes to exchange individuals through the migration protocol
// defined in pop.MigMethod. This is a convenience method for calling purposes.
func (pop *Population) migrate() {
	// Use the pointer to the demes to perform migration
	pop.Demes = pop.Migrator.apply(pop.Demes)
}

// Enhance each deme in the population. The deme level operations are done in
// parallel with a wait group. After all the deme operations have been run, the
// population level operations are run.
func (pop *Population) Enhance() {
	// Migrate the individuals between the demes
	if pop.Migrator != nil {
		pop.migrate()
	}
	// Use a wait group to run the genetic operators in each deme in parallel
	var wg sync.WaitGroup
	for i := range pop.Demes {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// 1. Breed
			if pop.Breeder != nil {
				pop.Demes[j].breed(pop.Selector, pop.Breeder)
			}
			// 2. Mutate
			if pop.Mutator != nil {
				pop.Demes[j].mutate(pop.Mutator)
			}
			// 3. Evaluate
			pop.Demes[j].evaluate(pop.Ff)
			// 4. Sort
			pop.Demes[j].sort()
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	pop.findBest()
}
