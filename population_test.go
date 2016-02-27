package gago

import "testing"

var (
	ga Population
)

func init() {
	// Instantiate a population
	ga = Population{
		NbDemes:       4,
		NbIndividuals: 30,
		Initializer:   UniformFloat{-10, 10},
		Selector:      Tournament{3},
		Breeder:       Parenthood{},
		Mutator:       Normal{0.1, 1},
		Migrator:      Shuffle{},
	}
	ga.Ff = FloatFunction{func(X []float64) float64 {
		sum := 0.0
		for _, x := range X {
			sum += x
		}
		return sum
	}}
	ga.Initialize(2)
	ga.Enhance()
}

func TestSizes(t *testing.T) {
	if len(ga.Demes) != ga.NbDemes {
		t.Error("Wrong number of demes")
	}
	for _, deme := range ga.Demes {
		if len(deme.Individuals) != ga.NbIndividuals {
			t.Error("Wrong number of individuals")
		}
	}
}

func TestGenerators(t *testing.T) {
	for i := range ga.Demes {
		for j := i + 1; j < len(ga.Demes); j++ {
			if &ga.Demes[i].generator == &ga.Demes[j].generator {
				t.Error("Deme share generator pointers")
			}
		}
	}
}

func TestSorted(t *testing.T) {
	for _, deme := range ga.Demes {
		for i := 0; i < len(deme.Individuals)-1; i++ {
			if deme.Individuals[i].Fitness > deme.Individuals[i+1].Fitness {
				t.Error("Individuals are not sorted in increasing order")
			}
		}
	}
}

func TestBest(t *testing.T) {
	for _, deme := range ga.Demes {
		for _, indi := range deme.Individuals {
			if ga.Best.Fitness > indi.Fitness {
				t.Error("The current best individual is not the overall best")
			}
		}
	}
}

func TestPointers(t *testing.T) {
	for _, deme := range ga.Demes {
		for _, indi := range deme.Individuals {
			if &ga.Best == &indi {
				t.Error("The current best individual shares a pointer with another individual")
			}
		}
	}
}
