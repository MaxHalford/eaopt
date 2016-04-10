package gago

import "testing"

var (
	ga            GA
	nbPopulations = 4
	nbIndividuals = 30
	nbGenes       = 2
)

func init() {
	ga = GAFloat(nbGenes, func(X []float64) float64 {
		sum := 0.0
		for _, x := range X {
			sum += x
		}
		return sum
	})
	ga.NbPopulations = nbPopulations
	ga.NbIndividuals = nbIndividuals
	ga.Initialize()
	ga.Enhance()
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

func TestPointers(t *testing.T) {
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if &ga.Best == &indi {
				t.Error("The current best individual shares a pointer with another individual")
			}
		}
	}
}
