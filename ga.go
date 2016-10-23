package gago

import (
	"errors"
	"sync"
	"time"
)

// A Topology holds all the information relative to the size of a GA.
type Topology struct {
	NbrPopulations int // Number of populations
	NbrClusters    int // Number of clusters each populations is split into before evolving
	NbrIndividuals int // Initial number of individuals in each population
}

// Validate the properties of a Topology.
func (topo Topology) Validate() error {
	// Check the number of clusters
	if topo.NbrClusters < 0 {
		return errors.New("'NbrClusters' should be higher or equal to 1 if provided")
	}
	// Check the number of individuals
	if topo.NbrIndividuals < 2 {
		return errors.New("'NbrIndividuals' should be higher or equal to 2")
	}
	// Check the number of populations
	if topo.NbrPopulations < 1 {
		return errors.New("'NbrPopulations' should be higher or equal to 1")
	}
	// No error
	return nil
}

// A GA contains population which themselves contain individuals.
type GA struct {
	// Fields that are provided by the user
	MakeGenome   GenomeMaker
	Topology     Topology
	Model        Model
	MigFrequency int // Frequency at which migrations occur
	Migrator     Migrator

	// Fields that are generated at runtime
	Best        Individual // Overall best individual (dummy initialization at the beginning)
	Duration    time.Duration
	Generations int
	Populations Populations
}

// Validate the parameters of a GA to ensure it will run correctly; some
// settings or combination of settings may be incoherent during runtime.
func (ga GA) Validate() error {
	// Check the GenomeMaker presence
	if ga.MakeGenome == nil {
		return errors.New("'GenomeMaker' cannot be nil")
	}
	// Check the topology is valid
	var topoErr = ga.Topology.Validate()
	if topoErr != nil {
		return topoErr
	}
	// Check the model presence
	if ga.Model == nil {
		return errors.New("'Model' cannot be nil")
	}
	// Check the model is valid
	var modelErr = ga.Model.Validate()
	if modelErr != nil {
		return modelErr
	}
	// Check the migration frequency in the presence of a migrator
	if ga.Migrator != nil && ga.MigFrequency < 1 {
		return errors.New("'MigFrequency' should be strictly higher than 0")
	}
	// No error
	return nil
}

// Initialize each population in the GA and assign an initial fitness to each
// individual in each population. Running Initialize after running Enhance will
// reset the GA entirely.
func (ga *GA) Initialize() {
	// Reset the number of generations and the elapsed duration
	ga.Generations = 0
	ga.Duration = 0
	// Create the populations
	ga.Populations = make([]Population, ga.Topology.NbrPopulations)
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Generate a population
			ga.Populations[j] = makePopulation(ga.Topology.NbrIndividuals, ga.MakeGenome)
			// Evaluate it's individuals
			ga.Populations[j].Individuals.Evaluate()
			// Sort it's individuals
			ga.Populations[j].Individuals.Sort()
		}(i)
	}
	wg.Wait()
	// The initial best individual is initialized randomly
	ga.Best = MakeIndividual(ga.MakeGenome(makeRandomNumberGenerator()))
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
	// Increment the generations counter at the beginning to not migrate at generation 0
	ga.Generations++
	// Migrate the individuals between the populations if there are enough
	// populations, there is a migrator and the migration frequency divides the
	// generation count
	if ga.Topology.NbrPopulations > 1 && ga.Migrator != nil && ga.Generations%ga.MigFrequency == 0 {
		ga.Migrator.Apply(ga.Populations)
	}
	// Use a wait group to enhance the populations in parallel
	var wg sync.WaitGroup
	for i := range ga.Populations {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			// Apply clustering if a number of clusters has been given
			if ga.Topology.NbrClusters > 0 {
				var clusters = ga.Populations[j].cluster(ga.Topology.NbrClusters)
				// Apply the evolution model to each cluster
				for k := range clusters {
					ga.Model.Apply(&clusters[k])
				}
				// Merge each cluster back into the original population
				ga.Populations[j].Individuals = clusters.merge()
			} else {
				// Else apply the evolution model to the entire population
				ga.Model.Apply(&ga.Populations[j])
			}
			// Evaluate and sort
			ga.Populations[j].Individuals.Evaluate()
			ga.Populations[j].Individuals.Sort()
			ga.Populations[j].Duration += time.Since(start)
		}(i)
	}
	wg.Wait()
	// Check if there is an individual that is better than the current one
	ga.findBest()
	ga.Duration += time.Since(start)
}
