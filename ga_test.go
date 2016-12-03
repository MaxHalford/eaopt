package gago

import (
	"math"
	"testing"
	"time"
)

func TestValidationSuccess(t *testing.T) {
	var err = ga.Validate()
	if err != nil {
		t.Error("GA parameters are invalid")
	}
}

func TestValidationGenomeMaker(t *testing.T) {
	var genomeMaker = ga.MakeGenome
	ga.MakeGenome = nil
	if ga.Validate() == nil {
		t.Error("Nil GenomeMaker should return an error")
	}
	ga.MakeGenome = genomeMaker
}

func TestValidationNPopulations(t *testing.T) {
	var nbrPopulations = ga.Topology.NPopulations
	ga.Topology.NPopulations = -1
	if ga.Topology.Validate() == nil {
		t.Error("Invalid number of Populations should return an error")
	}
	ga.Topology.NPopulations = nbrPopulations
}

func TestValidationNClusters(t *testing.T) {
	var nbrClusters = ga.Topology.NClusters
	ga.Topology.NClusters = -1
	if ga.Topology.Validate() == nil {
		t.Error("Invalid number of Clusters should return an error")
	}
	ga.Topology.NClusters = nbrClusters
}

func TestValidationNIndividuals(t *testing.T) {
	var nbrIndividuals = ga.Topology.NIndividuals
	ga.Topology.NIndividuals = -1
	if ga.Topology.Validate() == nil {
		t.Error("Invalid number of Individuals should return an error")
	}
	ga.Topology.NIndividuals = nbrIndividuals
}

func TestValidationModel(t *testing.T) {
	var model = ga.Model
	ga.Model = nil
	if ga.Validate() == nil {
		t.Error("Nil Model should return an error")
	}
	ga.Model = model
}

func TestValidationMigFrequency(t *testing.T) {
	var migFrequency = ga.MigFrequency
	ga.Migrator = MigRing{}
	ga.MigFrequency = 0
	if ga.Validate() == nil {
		t.Error("Invalid MigFrequency should return an error")
	}
	ga.Migrator = nil
	ga.MigFrequency = migFrequency
}

func TestRandomNumberGenerators(t *testing.T) {
	for i, pop1 := range ga.Populations {
		for j, pop2 := range ga.Populations {
			if i != j && &pop1.rng == &pop2.rng {
				t.Error("Population should not share random number generators")
			}
		}
	}
}

func TestBest(t *testing.T) {
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if ga.Best.Fitness > indi.Fitness {
				t.Error("The current best individual is not the overall best")
			}
		}
	}
}

func TestFindBest(t *testing.T) {
	// Check sure the findBest method works as expected
	var fitness = ga.Populations[0].Individuals[0].Fitness
	ga.Populations[0].Individuals[0].Fitness = math.Inf(-1)
	ga.findBest()
	if ga.Best.Fitness != math.Inf(-1) {
		t.Error("findBest didn't work")
	}
	ga.Populations[0].Individuals[0].Fitness = fitness
	// Check the best individual doesn't a share a pointer with anyone
	fitness = ga.Best.Fitness
	ga.Best.Fitness = 42
	if ga.Populations[0].Individuals[0].Fitness == 42 {
		t.Error("Best individual shares a pointer with an individual in the populations")
	}
	ga.Best.Fitness = fitness
}

func TestGenerations(t *testing.T) {
	if ga.Generations != nbrGenerations {
		t.Error("Generations counter wasn't incremented")
	}
}

// TestDuration verifies the sum of the duration of each population is higher
// the actual duration. This is due to the fact that each population runs on a
// separate core.
func TestDuration(t *testing.T) {
	var totalDuration time.Duration
	for _, pop := range ga.Populations {
		totalDuration += pop.Duration
	}
	if totalDuration < ga.Duration {
		t.Error("Inefficient parallelism")
	}
}

func BenchmarkEnhance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}
