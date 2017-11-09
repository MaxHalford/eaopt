package gago

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// A GA contains population which themselves contain individuals.
type GA struct {
	// Fields that are provided by the user
	NewGenome    NewGenome    `json:"-"`
	NPops        int          `json:"-"` // Number of Populations
	PopSize      int          `json:"-"` // Number of Individuls per Population
	Model        Model        `json:"-"`
	Migrator     Migrator     `json:"-"`
	MigFrequency int          `json:"-"` // Frequency at which migrations occur
	Speciator    Speciator    `json:"-"`
	Logger       *log.Logger  `json:"-"`
	Callback     func(ga *GA) `json:"-"`
	RNG          *rand.Rand   `json:"-"`

	// Fields that are generated at runtime
	Populations Populations   `json:"pops"`
	Best        Individual    `json:"overall_best"`
	CurrentBest Individual    `json:"generation_best"`
	Age         time.Duration `json:"duration"`
	Generations int           `json:"generations"`
}

// Validate the parameters of a GA to ensure it will run correctly; some
// settings or combination of settings may be incoherent during runtime.
func (ga GA) Validate() error {
	// Check the NewGenome presence
	if ga.NewGenome == nil {
		return errors.New("NewGenome cannot be nil")
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

// Find the best current Individual in each Population and then compare the best
// overall Individual to the current best Individual.
func (ga *GA) findBest() {
	// Start by finding the current best Individual
	ga.CurrentBest = Individual{Fitness: math.Inf(1)}
	for _, pop := range ga.Populations {
		if !pop.Individuals.IsSortedByFitness() {
			pop.Individuals.SortByFitness()
		}
		if pop.Individuals[0].Fitness < ga.CurrentBest.Fitness {
			ga.CurrentBest = pop.Individuals[0].Clone(pop.rng)
		}
	}
	// Compare the current best Individual to the overall Individual
	if ga.CurrentBest.Fitness < ga.Best.Fitness {
		ga.Best = ga.CurrentBest.Clone(ga.RNG)
	}
}

// Initialized indicates if the GA has been initialized or not.
func (ga GA) Initialized() bool {
	if len(ga.Populations) != ga.NPops {
		return false
	}
	return true
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population. Running Initialize after running Enhance will
// reset the GA entirely.
func (ga *GA) Initialize() {
	// Initialize the random number generator if it hasn't been set
	if ga.RNG == nil {
		ga.RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	// Initialize the populations asynchronously
	ga.Populations = make([]Population, ga.NPops)
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = newPopulation(ga.PopSize, ga.NewGenome, ga.RNG)
			// Evaluate its individuals
			ga.Populations[j].Individuals.Evaluate()
			// Sort its individuals
			ga.Populations[j].Individuals.SortByFitness()
			// Log current statistics if a logger has been provided
			if ga.Logger != nil {
				ga.Populations[j].Log(ga.Logger)
			}
		}(i)
	}
	wg.Wait()
	// Find the initial best Individual
	ga.Best = ga.Populations[0].Individuals[0]
	ga.findBest()
	// Execute the callback if it has been set
	if ga.Callback != nil {
		ga.Callback(ga)
	}
}

// Enhance each population in the GA. The population level operations are done
// in parallel with a wait group. After all the population operations have been
// run, the GA level operations are run.
func (ga *GA) Enhance() error {
	var start = time.Now()
	ga.Generations++
	// Check the GA has been initialized
	if !ga.Initialized() {
		return errors.New("The GA has not been initialized")
	}
	// Migrate the individuals between the populations if there are at least 2
	// Populations and that there is a migrator and that the migration frequency
	// divides the generation count
	if len(ga.Populations) > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations, ga.RNG)
	}
	var g errgroup.Group
	for i := range ga.Populations {
		i := i // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			var err error
			// Apply speciation if a positive number of species has been specified
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
				ga.Populations[i].Log(ga.Logger)
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
	// Execute the callback if it has been set
	if ga.Callback != nil {
		ga.Callback(ga)
	}
	// No error
	return nil
}

func (pop *Population) speciateEvolveMerge(spec Speciator, model Model) error {
	var (
		species, err = spec.Apply(pop.Individuals, pop.rng)
		pops         = make([]Population, len(species))
	)
	if err != nil {
		return err
	}
	// Create a subpopulation from each specie so that the evolution Model can
	// be applied to it.
	for i, specie := range species {
		pops[i] = Population{
			Individuals: specie,
			Age:         pop.Age,
			Generations: pop.Generations,
			ID:          randString(len(pop.ID), pop.rng),
			rng:         pop.rng,
		}
		err = model.Apply(&pops[i])
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
