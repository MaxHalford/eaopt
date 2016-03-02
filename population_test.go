package gago

import "testing"

var (
	floatGa  = Float
	demeSize = 2
)

func init() {
	floatGa.Ff = FloatFunction{func(X []float64) float64 {
		sum := 0.0
		for _, x := range X {
			sum += x
		}
		return sum
	}}
	floatGa.Initialize(demeSize)
	floatGa.Enhance()
}

func TestSizes(t *testing.T) {
	// Number of demes
	if len(floatGa.Demes) != floatGa.NbDemes {
		t.Error("Wrong number of demes")
	}
	// Number of individuals
	for _, deme := range floatGa.Demes {
		if len(deme.Individuals) != floatGa.NbIndividuals {
			t.Error("Wrong number of individuals")
		}
	}
	// Genome size
	for _, deme := range floatGa.Demes {
		for _, individual := range deme.Individuals {
			if len(individual.Genome) != demeSize {
				t.Error("Wrong genome size")
			}
		}
	}
}

func TestGenerators(t *testing.T) {
	for i := range floatGa.Demes {
		for j := i + 1; j < len(floatGa.Demes); j++ {
			if &floatGa.Demes[i].generator == &floatGa.Demes[j].generator {
				t.Error("Deme share generator")
			}
		}
	}
}

func TestSorted(t *testing.T) {
	for _, deme := range floatGa.Demes {
		for i := 0; i < len(deme.Individuals)-1; i++ {
			if deme.Individuals[i].Fitness > deme.Individuals[i+1].Fitness {
				t.Error("Individuals are not sorted in increasing order")
			}
		}
	}
}

func TestBest(t *testing.T) {
	for _, deme := range floatGa.Demes {
		for _, indi := range deme.Individuals {
			if floatGa.Best.Fitness > indi.Fitness {
				t.Error("The current best individual is not the overall best")
			}
		}
	}
}

func TestPointers(t *testing.T) {
	for _, deme := range floatGa.Demes {
		for _, indi := range deme.Individuals {
			if &floatGa.Best == &indi {
				t.Error("The current best individual shares a pointer with another individual")
			}
		}
	}
}
