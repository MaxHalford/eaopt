package gago

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// A GA contains population which themselves contain individuals.
type GA struct {
	// Fields that are provided by the user
	GenomeFactory GenomeFactory `json:"-"`
	NPops         int           `json:"-"`
	PopSize       int           `json:"-"`
	Model         Model         `json:"-"`
	Migrator      Migrator      `json:"-"`
	MigFrequency  int           `json:"-"` // Frequency at which migrations occur
	Speciator     Speciator     `json:"-"`
	Logger        *log.Logger   `json:"-"`

	// Fields that are generated at runtime
	Populations Populations   `json:"pops"`
	Best        Individual    `json:"best"` // Overall best individual
	Age         time.Duration `json:"duration"`
	Generations int           `json:"generations"`
	rng         *rand.Rand
}

// Validate the parameters of a GA to ensure it will run correctly; some
// settings or combination of settings may be incoherent during runtime.
func (ga GA) Validate() error {
	// Check the GenomeFactory presence
	if ga.GenomeFactory == nil {
		return errors.New("GenomeFactory cannot be nil")
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
	ga.rng = newRandomNumberGenerator()
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = newPopulation(
				ga.PopSize,
				ga.GenomeFactory,
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
	var rng = newRandomNumberGenerator()
	ga.Best = NewIndividual(ga.GenomeFactory(rng), rng)
	ga.findBest()
}

// Enhance each population in the GA. The population level operations are done
// in parallel with a wait group. After all the population operations have been
// run, the GA level operations are run.
func (ga *GA) Enhance() error {
	var start = time.Now()
	ga.Generations++
	// Migrate the individuals between the populations if there are enough
	// populations, there is a migrator and the migration frequency divides the
	// generation count
	if len(ga.Populations) > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations, ga.rng)
	}
	var g errgroup.Group
	for i := range ga.Populations {
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			var err error
			// Apply speciation if a positive number of species has been speficied
			if ga.Speciator != nil {
				err = ga.Populations[i].speciateEvolveMerge(ga.Speciator, ga.Model)
				if err != nil {
					return err
				}
			} else {
				// Else apply the evolution model to the entire population
				err = ga.Model.Apply(&ga.Populations[i])
				if err != nil {
					return err
				}
			}
			// Evaluate and sort
			ga.Populations[i].Individuals.Evaluate()
			ga.Populations[i].Individuals.SortByFitness()
			ga.Populations[i].Age += time.Since(start)
			ga.Populations[i].Generations++
			// Log current statistics if a logger has been provided
			if ga.Logger != nil {
				go ga.Populations[i].Log(ga.Logger)
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Age += time.Since(start)
	return nil
}

func (pop *Population) speciateEvolveMerge(spec Speciator, model Model) error {
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
			ID:          randString(len(pop.ID), pop.rng),
			rng:         pop.rng,
		}
		var err = model.Apply(&pops[i])
		if err != nil {
			return err
		}
	}
	// Merge each species back into the original population
	var i int
	for _, subpop := range pops {
		copy(pop.Individuals[i:i+len(subpop.Individuals)], subpop.Individuals)
		i += len(subpop.Individuals)
	}
	return nil
}
