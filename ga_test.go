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

func TestValidationGenomeFactory(t *testing.T) {
	var genomeFactory = ga.GenomeFactory
	ga.GenomeFactory = nil
	if ga.Validate() == nil {
		t.Error("Nil GenomeFactory should return an error")
	}
	ga.GenomeFactory = genomeFactory
}

func TestValidationNPopulations(t *testing.T) {
	var nPops = ga.NPops
	ga.NPops = -1
	if ga.Validate() == nil {
		t.Error("Invalid number of Populations should return an error")
	}
	ga.NPops = nPops
}

func TestValidationNIndividuals(t *testing.T) {
	var popSize = ga.PopSize
	ga.PopSize = -1
	if ga.Validate() == nil {
		t.Error("Invalid number of Individuals should return an error")
	}
	ga.PopSize = popSize
}

func TestValidationModel(t *testing.T) {
	var model = ga.Model
	// Check nil model raises error
	ga.Model = nil
	if ga.Validate() == nil {
		t.Error("Nil Model should return an error")
	}
	// Check invalid model raises error
	ga.Model = ModGenerational{
		Selector: SelTournament{
			NContestants: 3,
		},
		MutRate: -1,
	}
	if ga.Validate() == nil {
		t.Error("Invalid Model should return an error")
	}
	ga.Model = model
}

func TestValidationMigFrequency(t *testing.T) {
	var (
		migrator     = ga.Migrator
		migFrequency = ga.MigFrequency
	)
	ga.Migrator = MigRing{}
	ga.MigFrequency = 0
	if ga.Validate() == nil {
		t.Error("Invalid MigFrequency should return an error")
	}
	ga.Migrator = migrator
	ga.MigFrequency = migFrequency
}

func TestValidationSpeciator(t *testing.T) {
	var speciator = ga.Speciator
	ga.Speciator = SpecFitnessInterval{0}
	if ga.Validate() == nil {
		t.Error("Invalid Speciator should return an error")
	}
	ga.Speciator = speciator
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
		totalDuration += pop.Age
	}
	if totalDuration < ga.Age {
		t.Error("Inefficient parallelism")
	}
}

func TestSpeciateEvolveMerge(t *testing.T) {
	var (
		rng = newRandomNumberGenerator()
		pop = Population{ID: "42", rng: rng, Individuals: Individuals{}}
	)
	for i := 0; i < 7; i++ {
		pop.Individuals = append(pop.Individuals, Individual{Fitness: float64(i)})
	}
	pop.speciateEvolveMerge(SpecFitnessInterval{3}, ModIdentity{})
	for i := 0; i < 7; i++ {
		if pop.Individuals[i].Fitness != float64(i) {
			t.Error("speciateEvolveMerge did not work as expected")
		}
	}
}
