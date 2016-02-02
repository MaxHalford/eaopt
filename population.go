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
	Boundary float64
	// Selection method
	SelMethod func(Individuals, *rand.Rand) Individual
	// Crossover method
	CrossMethod func(Individuals, *rand.Rand) Individual
	// Crossover size
	CrossSize int
	// Mutation method
	MutMethod func(indi *Individual, rate float64, intensity float64, generator *rand.Rand)
	// Mutation rate
	MutRate float64
	// Mutation intensity
	MutIntensity float64
	// Migration method
	MigMethod func(demes []Deme) []Deme
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
		deme.Initialize(pop.NbGenes, pop.Boundary)
		// Add it to the population
		pop.Demes[i] = deme
		// Initial evaluation
		pop.Demes[i].Evaluate(pop.Ff)
	}
}

// Migrate allows demes to exchange individuals through the migration protocol
// defined in pop.MigMethod. This is a convenience method for calling purposes.
func (pop *Population) Migrate() {
	// Use the pointer to the demes to perform migration
	pop.Demes = pop.MigMethod(pop.Demes)
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
// parallel with a wait group.
func (pop *Population) Enhance() {
	var wg sync.WaitGroup
	for i := range pop.Demes {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			pop.Demes[j].Crossover(pop.SelMethod, pop.CrossMethod, pop.CrossSize)
			pop.Demes[j].Mutate(pop.MutMethod, pop.MutRate, pop.MutIntensity)
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
