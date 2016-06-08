package gago

import (
	"math"
	"testing"
	"time"
)

var (
	ga            GA
	nbPopulations = 4
	nbIndividuals = 30
	nbGenes       = 2
	nbParents     = 6
	nbGenerations = 10
)

func init() {
	ga.NbPopulations = nbPopulations
	ga.NbIndividuals = nbIndividuals
	ga.NbGenes = nbGenes
	ga.NbParents = nbParents
	ga.Initializer = InitUniformF{
		Lower: -1,
		Upper: 1,
	}
	ga.Ff = FloatFunction{
		Image: func(X []float64) float64 {
			var sum float64
			for _, x := range X {
				sum += x
			}
			return sum
		},
	}
	ga.Selector = SelTournament{
		NbParticipants: 3,
	}
	ga.Crossover = CrossUniformF{}
	ga.Mutator = MutNormalF{
		Rate: 0.5,
		Std:  3,
	}
	ga.MutRate = 0.5
	ga.Migrator = MigShuffle{}
	ga.MigFrequency = 10
	ga.Initialize()
	for i := 0; i < nbGenerations; i++ {
		ga.Enhance()
	}
}

func TestValidationSuccess(t *testing.T) {
	var err = ga.Validate()
	if err != nil {
		t.Error("GA parameters are invalid")
	}
}

func TestValidationNbPopulations(t *testing.T) {
	// Check invalid number of population
	ga.NbPopulations = 0
	if ga.Validate() == nil {
		t.Error("Invalid number of populations didn't return an error")
	}
}

func TestValidationNbIndividuals(t *testing.T) {
	// Check invalid number of individuals
	ga.NbIndividuals = 1
	if ga.Validate() == nil {
		t.Error("Invalid number of individuals didn't return an error")
	}
}

func TestValidationNbGenes(t *testing.T) {
	// Check invalid number of genes
	ga.NbGenes = 0
	if ga.Validate() == nil {
		t.Error("Invalid number of genes didn't return an error")
	}
}

func TestValidationNbParents(t *testing.T) {
	// Check invalid number of parents
	ga.NbParents = -1
	if ga.Validate() == nil {
		t.Error("Invalid number of parents didn't return an error")
	}
}

func TestValidationNbParentsNbIndividuals(t *testing.T) {
	// Check number of parents higher than number of individuals
	ga.NbIndividuals = 1
	ga.NbParents = 2
	if ga.Validate() == nil {
		t.Error("More parents than individuals didn't return an error")
	}
}

func TestValidationFf(t *testing.T) {
	// Check presence of fitness function
	ga.Ff = nil
	if ga.Validate() == nil {
		t.Error("Nil fitness function didn't return an error")
	}
}

func TestValidationInit(t *testing.T) {
	// Check presence of initializer
	ga.Initializer = nil
	if ga.Validate() == nil {
		t.Error("Nil initializer didn't return an error")
	}
}

func TestValidationSel(t *testing.T) {
	// Check presence of selector
	ga.Selector = nil
	if ga.Validate() == nil {
		t.Error("Nil selector didn't return an error")
	}
}

func TestValidationCross(t *testing.T) {
	// Check presence of crossover
	ga.Crossover = nil
	if ga.Validate() == nil {
		t.Error("Nil crossover didn't return an error")
	}
}

func TestValidationMutRateTooLow(t *testing.T) {
	// Check mutation rate
	ga.MutRate = -1
	if ga.Validate() == nil {
		t.Error("Invalid mutation rate")
	}
}

func TestValidationTooHigh(t *testing.T) {
	// Check mutation rate
	ga.MutRate = 2
	if ga.Validate() == nil {
		t.Error("Invalid mutation rate")
	}
}

func TestValidationMigFrequency(t *testing.T) {
	// Check migration frequency
	ga.MigFrequency = 0
	if ga.Validate() == nil {
		t.Error("Invalid migration frequency")
	}
}

func TestSizes(t *testing.T) {
	// Number of Populations
	if len(ga.Populations) != nbPopulations {
		t.Error("Wrong number of Populations")
	}
	// Number of individuals
	for _, pop := range ga.Populations {
		if len(pop.Individuals) != nbIndividuals {
			t.Error("Wrong number of individuals")
		}
	}
	// Genome size
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if len(indi.Genome) != nbGenes {
				t.Error("Wrong genome size")
			}
		}
	}
}

func TestGenerators(t *testing.T) {
	for i := range ga.Populations {
		for j := i + 1; j < len(ga.Populations); j++ {
			if &ga.Populations[i].generator == &ga.Populations[j].generator {
				t.Error("population share generator")
			}
		}
	}
}

func TestSorted(t *testing.T) {
	for _, pop := range ga.Populations {
		for i := 0; i < len(pop.Individuals)-1; i++ {
			if pop.Individuals[i].Fitness > pop.Individuals[i+1].Fitness {
				t.Error("Individuals are not sorted in increasing order")
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
	ga.Populations[0].Individuals[0].Fitness = math.Inf(-1)
	ga.findBest()
	if ga.Best.Fitness != math.Inf(-1) {
		t.Error("findBest didn't work")
	}
	ga.Best.Fitness = 42
	if ga.Populations[0].Individuals[0].Fitness == 42 {
		t.Error("Best individual is linked to an individual")
	}
}

func TestGenerations(t *testing.T) {
	if ga.Generations != nbGenerations {
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
