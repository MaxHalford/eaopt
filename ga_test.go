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
	initializer   = InitUniformF{
		Lower: -1,
		Upper: 1,
	}
	ff = FloatFunction{
		Image: func(X []float64) float64 {
			var sum float64
			for _, x := range X {
				sum += x
			}
			return sum
		},
	}
	selector = SelTournament{
		NbParticipants: 3,
	}
	crossover = CrossUniformF{}
	mutator   = MutNormalF{
		Rate: 0.5,
		Std:  3,
	}
	mutRate      = 0.5
	migrator     = MigShuffle{}
	migFrequency = 10
)

func init() {
	ga.NbPopulations = nbPopulations
	ga.NbIndividuals = nbIndividuals
	ga.NbGenes = nbGenes
	ga.NbParents = nbParents
	ga.Initializer = initializer
	ga.Ff = ff
	ga.Selector = selector
	ga.Crossover = crossover
	ga.Mutator = mutator
	ga.MutRate = mutRate
	ga.Migrator = migrator
	ga.MigFrequency = migFrequency
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
	ga.NbPopulations = nbIndividuals
}

func TestValidationNbIndividuals(t *testing.T) {
	// Check invalid number of individuals
	ga.NbIndividuals = 1
	if ga.Validate() == nil {
		t.Error("Invalid number of individuals didn't return an error")
	}
	ga.NbIndividuals = nbIndividuals
}

func TestValidationNbGenes(t *testing.T) {
	// Check invalid number of genes
	ga.NbGenes = 0
	if ga.Validate() == nil {
		t.Error("Invalid number of genes didn't return an error")
	}
	ga.NbGenes = nbGenes
}

func TestValidationNbParents(t *testing.T) {
	// Check invalid number of parents
	ga.NbParents = -1
	if ga.Validate() == nil {
		t.Error("Invalid number of parents didn't return an error")
	}
	ga.NbParents = nbParents
}

func TestValidationNbParentsNbIndividuals(t *testing.T) {
	// Check number of parents higher than number of individuals
	ga.NbIndividuals = 1
	ga.NbParents = 2
	if ga.Validate() == nil {
		t.Error("More parents than individuals didn't return an error")
	}
	ga.NbIndividuals = nbIndividuals
	ga.NbParents = nbParents
}

func TestValidationFf(t *testing.T) {
	// Check presence of fitness function
	ga.Ff = nil
	if ga.Validate() == nil {
		t.Error("Nil fitness function didn't return an error")
	}
	ga.Ff = ff
}

func TestValidationInit(t *testing.T) {
	// Check presence of initializer
	ga.Initializer = nil
	if ga.Validate() == nil {
		t.Error("Nil initializer didn't return an error")
	}
	ga.Initializer = initializer
}

func TestValidationSel(t *testing.T) {
	// Check presence of selector
	ga.Selector = nil
	if ga.Validate() == nil {
		t.Error("Nil selector didn't return an error")
	}
	ga.Selector = selector
}

func TestValidationCross(t *testing.T) {
	// Check presence of crossover
	ga.Crossover = nil
	if ga.Validate() == nil {
		t.Error("Nil crossover didn't return an error")
	}
	ga.Crossover = crossover
}

func TestValidationMutRateTooLow(t *testing.T) {
	// Check mutation rate
	ga.MutRate = -1
	if ga.Validate() == nil {
		t.Error("Invalid mutation rate")
	}
	ga.MutRate = mutRate
}

func TestValidationTooHigh(t *testing.T) {
	// Check mutation rate
	ga.MutRate = 2
	if ga.Validate() == nil {
		t.Error("Invalid mutation rate")
	}
	ga.MutRate = mutRate
}

func TestValidationMigFrequency(t *testing.T) {
	// Check migration frequency
	ga.MigFrequency = 0
	if ga.Validate() == nil {
		t.Error("Invalid migration frequency")
	}
	ga.MigFrequency = migFrequency
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
