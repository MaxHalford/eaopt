package gago

import (
	"fmt"
	"math"
	"testing"
)

func TestIndividualsString(t *testing.T) {
	var testCases = []struct {
		indis Individuals
		str   string
	}{
		{
			indis: Individuals{
				Individual{
					Genome:    Vector{0, 1, 2},
					Fitness:   42,
					Evaluated: true,
					ID:        "bob",
				},
				Individual{
					Genome:    Vector{0, 1, 2},
					Fitness:   42,
					Evaluated: false,
					ID:        "ALICE",
				},
			},
			str: "bob - 42.000 - [0 1 2]\nALICE - ??? - [0 1 2]",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			if tc.indis.String() != tc.str {
				t.Errorf("Expected %s, got %s", tc.str, tc.indis.String())
			}
		})
	}
}

func TestNewIndividuals(t *testing.T) {
	var rng = newRand()
	for _, n := range []int{1, 2, 42} {
		var indis = newIndividuals(n, NewVector, rng)
		if len(indis) != n {
			t.Error("newIndividuals didn't generate the right number of individuals")
		}
	}
}

func TestCloneIndividuals(t *testing.T) {
	var (
		rng    = newRand()
		indis  = newIndividuals(20, NewVector, rng)
		clones = indis.Clone(rng)
	)
	for _, indi := range indis {
		for _, clone := range clones {
			if &indi == &clone || indi.ID == clone.ID {
				t.Error("Cloning did not work as expected")
			}
		}
	}
}

func TestEvaluateIndividuals(t *testing.T) {
	var indis = newIndividuals(10, NewVector, newRand())
	for _, indi := range indis {
		if indi.Evaluated {
			t.Error("Individual shouldn't have Evaluated set to True")
		}
	}
	indis.Evaluate(false)
	for _, indi := range indis {
		if !indi.Evaluated {
			t.Error("Individual should have Evaluated set to True")
		}
	}
}

func TestEvaluateIndividualsParallel(t *testing.T) {
	var indis = newIndividuals(10, NewVector, newRand())
	for _, indi := range indis {
		if indi.Evaluated {
			t.Error("Individual shouldn't have Evaluated set to True")
		}
	}
	indis.Evaluate(true)
	for _, indi := range indis {
		if !indi.Evaluated {
			t.Error("Individual should have Evaluated set to True")
		}
	}
}

func TestMutateIndividuals(t *testing.T) {
	var (
		rng   = newRand()
		indis = newIndividuals(10, NewVector, rng)
	)
	indis.Evaluate(false)
	indis.Mutate(1, rng)
	for _, indi := range indis {
		if indi.Evaluated {
			t.Error("Individual shouldn't have Evaluated set to True")
		}
	}
}

func TestIndividualsSortByFitness(t *testing.T) {
	var indis = newIndividuals(10, NewVector, newRand())
	// Assign a fitness to each individual in decreasing order
	for i := range indis {
		indis[i].Fitness = float64(len(indis) - i)
	}
	indis.SortByFitness()
	// Check fitnesses are in increasing order
	for i := 1; i < len(indis); i++ {
		if indis[i-1].Fitness > indis[i].Fitness {
			t.Error("Individuals are not sorted")
		}
	}
}

func TestGetFitnesses(t *testing.T) {
	var (
		indis = Individuals{
			Individual{Fitness: 0.0},
			Individual{Fitness: 1.0},
			Individual{Fitness: 2.0},
		}
		target    = []float64{0.0, 1.0, 2.0}
		fitnesses = indis.getFitnesses()
	)
	for i, fitness := range fitnesses {
		if fitness != target[i] {
			t.Error("getFitnesses didn't work as expected")
		}
	}
}

func TestFitMin(t *testing.T) {
	var testCases = []struct {
		indis Individuals
		min   float64
	}{
		{Individuals{
			Individual{Fitness: 1.0},
		}, 1.0},
		{Individuals{
			Individual{Fitness: 2.0},
			Individual{Fitness: 1.0},
		}, 1.0},
		{Individuals{
			Individual{Fitness: 1.0},
			Individual{Fitness: -1.0},
		}, -1.0},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			if tc.indis.FitMin() != tc.min {
				t.Error("FitMin didn't work as expected")
			}
		})
	}
}

func TestFitMax(t *testing.T) {
	var testCases = []struct {
		indis Individuals
		max   float64
	}{
		{Individuals{
			Individual{Fitness: 1.0},
		}, 1.0},
		{Individuals{
			Individual{Fitness: 2.0},
			Individual{Fitness: 1.0},
		}, 2.0},
		{Individuals{
			Individual{Fitness: 1.0},
			Individual{Fitness: -1.0},
		}, 1.0},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			if tc.indis.FitMax() != tc.max {
				t.Error("FitMax didn't work as expected")
			}
		})
	}
}

func TestFitAvg(t *testing.T) {
	var testCases = []struct {
		indis Individuals
		mean  float64
	}{
		{Individuals{
			Individual{Fitness: 1.0},
		}, 1.0},
		{Individuals{
			Individual{Fitness: 1.0},
			Individual{Fitness: 2.0},
		}, 1.5},
		{Individuals{
			Individual{Fitness: -1.0},
			Individual{Fitness: 1.0},
		}, 0.0},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			if tc.indis.FitAvg() != tc.mean {
				t.Error("FitAvg didn't work as expected")
			}
		})
	}
}

func TestFitStd(t *testing.T) {
	var testCases = []struct {
		indis    Individuals
		variance float64
	}{
		{Individuals{
			Individual{Fitness: 1.0},
		}, 0.0},
		{Individuals{
			Individual{Fitness: -1.0},
			Individual{Fitness: 1.0},
		}, 1.0},
		{Individuals{
			Individual{Fitness: -2.0},
			Individual{Fitness: 2.0},
		}, 4.0},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			if tc.indis.FitStd() != math.Sqrt(tc.variance) {
				t.Error("FitStd didn't work as expected")
			}
		})
	}
}
