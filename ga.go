package eaopt

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// A GA contains populations which themselves contain individuals.
type GA struct {
	GAConfig `json:"-"`

	// Fields generated at runtime
	Populations Populations   `json:"populations"`
	HallOfFame  Individuals   `json:"hall_of_fame"` // Sorted best Individuals ever encountered
	Age         time.Duration `json:"duration"`     // Duration during which the GA has been evolved
	Generations uint          `json:"generations"`  // Number of generations the GA has been evolved
}

// Find the best current Individual in each population and then compare the best
// overall Individual to the current best Individual. The Individuals in each
// population are expected to be sorted.
func updateHallOfFame(hof Individuals, indis Individuals, rng *rand.Rand) {
	var k = len(hof)
	// Start by finding the current best Individual
	for _, indi := range indis[:minInt(k, len(indis))] {
		// Find if and where the Individual should fit in the hall of fame
		var (
			f = func(i int) bool { return indi.Fitness < hof[i].Fitness }
			i = sort.Search(k, f)
		)
		if i < k {
			// Shift the hall of fame to the right
			copy(hof[i+1:], hof[i:])
			// Insert the new Individual
			hof[i] = indi.Clone(rng)
		}
	}
}

func (ga *GA) init(newGenome func(rng *rand.Rand) Genome) error {
	// Reset counters
	ga.Generations = 0
	ga.Age = 0

	// Create the initial Populations
	ga.Populations = make(Populations, ga.NPops)
	for i := range ga.Populations {
		ga.Populations[i] = newPopulation(ga.PopSize, newGenome, ga.RNG)
		// Evaluate and sort
		err := ga.Populations[i].Individuals.Evaluate(ga.ParallelEval)
		if err != nil {
			return err
		}
		ga.Populations[i].Individuals.SortByFitness()
		// Log current statistics if a logger has been provided
		if ga.Logger != nil {
			ga.Populations[i].Log(ga.Logger)
		}
	}

	// Initialize the hall of fame
	ga.HallOfFame = make(Individuals, ga.HofSize)
	for i := range ga.HallOfFame {
		ga.HallOfFame[i] = Individual{Fitness: math.Inf(1)}
	}
	for _, pop := range ga.Populations {
		updateHallOfFame(ga.HallOfFame, pop.Individuals, pop.RNG)
	}

	// Execute the callback if it has been set
	if ga.Callback != nil {
		ga.Callback(ga)
	}

	return nil
}

// Evolve a GA's Populations in parallel.
func (ga *GA) evolve() error {
	var start = time.Now()
	ga.Generations++

	// Migrate the individuals between the populations if there are at least 2
	// Populations and that there is a migrator and that the migration frequency
	// divides the generation count
	if len(ga.Populations) > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations, ga.RNG)
	}

	var f = func(pop *Population) error {
		var err error
		// Apply speciation if a positive number of species has been specified
		if ga.Speciator != nil {
			err = pop.speciateEvolveMerge(ga.Speciator, ga.Model)
			if err != nil {
				return err
			}
		} else {
			// Else apply the evolution model to the entire population
			err = ga.Model.Apply(pop)
			if err != nil {
				return err
			}
		}
		// Evaluate and sort
		err = pop.Individuals.Evaluate(ga.ParallelEval)
		if err != nil {
			return err
		}
		pop.Individuals.SortByFitness()
		// Record time spent evolving
		pop.Age += time.Since(start)
		pop.Generations++
		// Log current statistics if a logger has been provided
		if ga.Logger != nil {
			pop.Log(ga.Logger)
		}
		return err
	}

	var err = ga.Populations.Apply(f)
	if err != nil {
		return err
	}
	// Update HallOfFame
	for _, pop := range ga.Populations {
		updateHallOfFame(ga.HallOfFame, pop.Individuals, pop.RNG)
	}

	ga.Age += time.Since(start)

	// Execute the callback if it has been set
	if ga.Callback != nil {
		ga.Callback(ga)
	}

	return nil
}

// Minimize evolves the GA's Populations following the given evolutionary
// method. The GA's hall of fame is updated after each generation.
func (ga *GA) Minimize(newGenome func(rng *rand.Rand) Genome) error {
	// Initialize the GA
	var err = ga.init(newGenome)
	if err != nil {
		return err
	}

	// Go through the generations
	for i := uint(0); i < ga.NGenerations; i++ {
		// Check for early stopping
		if ga.EarlyStop != nil && ga.EarlyStop(ga) {
			return nil
		}
		if err := ga.evolve(); err != nil {
			return err
		}
	}
	return nil
}

func (pop *Population) speciateEvolveMerge(spec Speciator, model Model) error {
	var (
		species, err = spec.Apply(pop.Individuals, pop.RNG)
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
			ID:          randString(len(pop.ID), pop.RNG),
			RNG:         pop.RNG,
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
