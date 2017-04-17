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
	// Fields that are provided by the user
	MakeGenome   GenomeMaker
	NPops        int
	PopSize      int
	Model        Model
	Migrator     Migrator
	MigFrequency int // Frequency at which migrations occur
	Speciator    Speciator
	Logger       *log.Logger

	// Fields that are generated at runtime
	Populations Populations
	Best        Individual // Overall best individual (dummy initialization at the beginning)
	Age         time.Duration
	Generations int
	rng         *rand.Rand
}

// Validate the parameters of a GA to ensure it will run correctly; some
// settings or combination of settings may be incoherent during runtime.
func (ga GA) Validate() error {
	// Check the GenomeMaker presence
	if ga.MakeGenome == nil {
		return errors.New("GenomeMaker cannot be nil")
	}
	// Check the number of populations is higher than 0
	if ga.NPops < 1 {
		return errors.New("NPops should be higher than 0")
	}
	// Check the number of individuals per population is higher than 0
	if ga.PopSize < 1 {
		return errors.New("PopSize should be higher than 0")
	}
	// Check the model presence
	if ga.Model == nil {
		return errors.New("Model cannot be nil")
	}
	// Check the model is valid
	var modelErr = ga.Model.Validate()
	if modelErr != nil {
		return modelErr
	}
	// Check the migration frequency if a Migrator has been provided
	if ga.Migrator != nil && ga.MigFrequency < 1 {
		return errors.New("MigFrequency should be strictly higher than 0")
	}
	// Check the speciator is valid if it has been provided
	if ga.Speciator != nil {
		if specErr := ga.Speciator.Validate(); specErr != nil {
			return specErr
		}
	}
	// No error
	return nil
}

// Find the best individual in each population and then compare the best overall
// individual to the current best individual. This method supposes that the
// populations have been preemptively ascendingly sorted by fitness so that
// checking the first individual of each population is sufficient.
func (ga *GA) findBest() {
	for _, pop := range ga.Populations {
		var best = pop.Individuals[0]
		if best.Fitness < ga.Best.Fitness {
			ga.Best = best
		}
	}
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population. Running Initialize after running Enhance will
// reset the GA entirely.
func (ga *GA) Initialize() {
	ga.Populations = make([]Population, ga.NPops)
	ga.rng = makeRandomNumberGenerator()
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = makePopulation(
				ga.PopSize,
				ga.MakeGenome,
				randString(3, ga.rng),
			)
			// Evaluate its individuals
			ga.Populations[j].Individuals.Evaluate()
			// Sort its individuals
			ga.Populations[j].Individuals.SortByFitness()
			// Log current statistics if a logger has been provided
			if ga.Logger != nil {
				go ga.Populations[j].Log(ga.Logger)
			}
		}(i)
	}
	wg.Wait()
	// The initial best individual is initialized randomly
	var rng = makeRandomNumberGenerator()
	ga.Best = MakeIndividual(ga.MakeGenome(rng), rng)
	ga.findBest()
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
	if len(ga.Populations) > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations, ga.rng)
	}
	// Use a wait group to enhance the populations in parallel
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Apply speciation if a positive number of species has been speficied
			if ga.Speciator != nil {
				ga.Populations[j].speciateEvolveMerge(ga.Speciator, ga.Model)
			} else {
				// Else apply the evolution model to the entire population
				ga.Model.Apply(&ga.Populations[j])
			}
			// Evaluate and sort
			ga.Populations[j].Individuals.Evaluate()
			ga.Populations[j].Individuals.SortByFitness()
			ga.Populations[j].Age += time.Since(start)
			ga.Populations[j].Generations++
			// Log current statistics if a logger has been provided
			if ga.Logger != nil {
				go ga.Populations[j].Log(ga.Logger)
			}
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Age += time.Since(start)
}

func (pop *Population) speciateEvolveMerge(spec Speciator, model Model) {
	var (
		species = spec.Apply(pop.Individuals, pop.rng)
		pops    = make([]Population, len(species))
	)
	// Create a slice of population from the obtained species and evolve each one separately
	for i, specie := range species {
		pops[i] = Population{
			Individuals: specie,
			Age:         pop.Age,
			Generations: pop.Generations,
			ID:          randString(3, pop.rng),
			rng:         pop.rng,
		}
		model.Apply(&pops[i])
	}
	// Merge each species back into the original population
	var i int
	for _, pop := range pops {
		copy(pop.Individuals[i:i+len(pop.Individuals)], pop.Individuals)
		i += len(pop.Individuals)
	}
}
