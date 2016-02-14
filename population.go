package gago

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// A Population contains Demes which contains Individuals.
type Population struct {
	// Number of demes
	NbDemes int
	// Number of individuals in each deme
	NbIndividuals int
	// Number of genes in each individual (defined by the user)
	NbGenes int
	// Fitness function to evaluate individuals (defined by the user)
	Ff func([]float64) float64
	// Demes
	Demes []Deme
	// Overall best individual (dummy initialization at the begining)
	Best Individual
	// Initial random boundaries
	Initializer Initializer
	// Selection method
	Selector Selector
	// Crossover method
	Crossover Crossover
	// Mutation method
	Mutator Mutator
	// Migration method
	Migrator Migrator
}

// Initialize each deme in the population and assign an initial fitness to each
// individual in each deme.
func (pop *Population) Initialize(ff func([]float64) float64, variables int) {
	// Fitness function
	pop.Ff = ff
	// Number of genes in each individual
	pop.NbGenes = variables
	// Create the demes
	pop.Demes = make([]Deme, pop.NbDemes)
	// Best individual (dummy instantiation)
	pop.Best = Individual{make([]float64, pop.NbGenes), math.Inf(1)}
	for i := range pop.Demes {
		// Create a new random number generator
		var source = rand.NewSource(time.Now().UnixNano())
		var generator = rand.New(source)
		// Create the deme
		var deme = Deme{pop.NbIndividuals, make([]Individual, pop.NbIndividuals), generator}
		// Initialize the deme
		deme.Initialize(pop.NbGenes, pop.Initializer)
		// Add it to the population
		pop.Demes[i] = deme
		// Initial evaluation
		pop.Demes[i].Evaluate(pop.Ff)
	}
	pop.FindBest()
}

// Migrate allows demes to exchange individuals through the migration protocol
// defined in pop.MigMethod. This is a convenience method for calling purposes.
func (pop *Population) Migrate() {
	// Use the pointer to the demes to perform migration
	pop.Demes = pop.Migrator.Apply(pop.Demes)
}

// FindBest stores the best individual over all demes.
func (pop *Population) FindBest() {
	for _, deme := range pop.Demes {
		if deme.Individuals[0].Fitness < pop.Best.Fitness {
			pop.Best = deme.Individuals[0]
		}
	}
}

// Enhance each deme in the population. The deme level operations are done in
// parallel with a wait group. After all the deme operations have been run, the
// population level operations are run.
func (pop *Population) Enhance() {
	var wg sync.WaitGroup
	for i := range pop.Demes {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			pop.Demes[j].Crossover(pop.Selector, pop.Crossover)
			pop.Demes[j].Mutate(pop.Mutator)
			pop.Demes[j].Evaluate(pop.Ff)
			pop.Demes[j].Sort()
		}(i)
	}
	wg.Wait()
	// Check if there is a new best individual
	pop.FindBest()
	// Migrate the individuals between the demes
	pop.Migrate()
}
