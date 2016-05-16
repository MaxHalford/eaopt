package gago

import (
	"math"
	"testing"
)

var (
	nbPopulations = 4
	nbIndividuals = 30
	nbGenes       = 2
	nbParents     = 6
	nbGenerations = 10
	ga            = GA{
		NbPopulations: nbPopulations,
		NbIndividuals: nbIndividuals,
		NbGenes:       nbGenes,
		NbParents:     nbParents,
		Initializer: IFUniform{
			Lower: -1,
			Upper: 1,
		},
		Ff: FloatFunction{
			Image: func(X []float64) float64 {
				var sum float64
				for _, x := range X {
					sum += x
				}
				return sum
			},
		},
		Selector: STournament{
			NbParticipants: 3,
		},
		Crossover: CFUniform{},
		Mutator: MutFNormal{
			Rate: 0.5,
			Std:  3,
		},
		MutRate:      0.5,
		Migrator:     MigShuffle{},
		MigFrequency: 10,
	}
)

func init() {
	ga.Initialize()
	for i := 0; i < nbGenerations; i++ {
		ga.Enhance()
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
