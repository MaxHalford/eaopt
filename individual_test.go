package gago

import (
	"math"
	"testing"
)

func TestCloneIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi1  = NewIndividual(genome, rng)
		indi2  = indi1.Clone(rng)
	)
	if &indi1 == &indi2 || &indi1.Genome == &indi2.Genome {
		t.Error("Individual was not deep copied")
	}
}

func TestEvaluateIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi   = NewIndividual(genome, rng)
	)
	if indi.Evaluated {
		t.Error("Individual shouldn't have Evaluated set to True")
	}
	indi.Evaluate()
	if !indi.Evaluated {
		t.Error("Individual should have Evaluated set to True")
	}
}

func TestMutateIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi   = NewIndividual(genome, rng)
	)
	indi.Evaluate()
	indi.Mutate(rng)
	if indi.Evaluated {
		t.Error("Individual shouldn't have Evaluated set to True")
	}
}

func TestCrossoverIndividual(t *testing.T) {
	var (
		rng                    = newRandomNumberGenerator()
		indi1                  = NewIndividual(NewVector(rng), rng)
		indi2                  = NewIndividual(NewVector(rng), rng)
		offspring1, offspring2 = indi1.Crossover(indi2, rng)
	)
	if offspring1.Evaluated || offspring2.Evaluated {
		t.Error("Offsprings shouldn't have Evaluated set to True")
	}
	if &offspring1 == &indi1 || &offspring1 == &indi2 || &offspring2 == &indi1 || &offspring2 == &indi2 {
		t.Error("Offsprings shouldn't share pointers with parents")
	}
}

func TestNewIndividuals(t *testing.T) {
	var rng = newRandomNumberGenerator()
	for _, n := range []int{1, 2, 42} {
		var indis = newIndividuals(n, NewVector, rng)
		if len(indis) != n {
			t.Error("newIndividuals didn't generate the right number of individuals")
		}
	}
}

func TestCloneIndividuals(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
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
	var indis = newIndividuals(10, NewVector, newRandomNumberGenerator())
	for _, indi := range indis {
		if indi.Evaluated {
			t.Error("Individual shouldn't have Evaluated set to True")
		}
	}
	indis.Evaluate()
	for _, indi := range indis {
		if !indi.Evaluated {
			t.Error("Individual should have Evaluated set to True")
		}
	}
}

func TestMutateIndividuals(t *testing.T) {
	var (
		rng   = newRandomNumberGenerator()
		indis = newIndividuals(10, NewVector, rng)
	)
	indis.Evaluate()
	indis.Mutate(1, rng)
	for _, indi := range indis {
		if indi.Evaluated {
			t.Error("Individual shouldn't have Evaluated set to True")
		}
	}
}

func TestIndividualsSortByFitness(t *testing.T) {
	var indis = newIndividuals(10, NewVector, newRandomNumberGenerator())
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
	for _, test := range testCases {
		if test.indis.FitMin() != test.min {
			t.Error("FitMin didn't work as expected")
		}
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
	for _, test := range testCases {
		if test.indis.FitMax() != test.max {
			t.Error("FitMax didn't work as expected")
		}
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
	for _, test := range testCases {
		if test.indis.FitAvg() != test.mean {
			t.Error("FitAvg didn't work as expected")
		}
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
	for _, test := range testCases {
		if test.indis.FitStd() != math.Sqrt(test.variance) {
			t.Error("FitStd didn't work as expected")
		}
	}
}
