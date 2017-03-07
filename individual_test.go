package gago

import (
	"math"
	"reflect"
	"testing"
)

func TestDeepCopyIndividual(t *testing.T) {
	var (
		genome = MakeVector(makeRandomNumberGenerator())
		indi1  = MakeIndividual(genome)
		indi2  = indi1.DeepCopy()
	)
	if &indi1 == &indi2 || &indi1.Genome == &indi2.Genome {
		t.Error("Individual was not deep copied")
	}
}

func TestEvaluateIndividual(t *testing.T) {
	var (
		genome = MakeVector(makeRandomNumberGenerator())
		indi   = MakeIndividual(genome)
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
		rng    = makeRandomNumberGenerator()
		genome = MakeVector(rng)
		indi   = MakeIndividual(genome)
	)
	indi.Evaluate()
	indi.Mutate(rng)
	if indi.Evaluated {
		t.Error("Individual shouldn't have Evaluated set to True")
	}
}

func TestCrossoverIndividual(t *testing.T) {
	var (
		rng                    = makeRandomNumberGenerator()
		indi1                  = MakeIndividual(MakeVector(rng))
		indi2                  = MakeIndividual(MakeVector(rng))
		offspring1, offspring2 = indi1.Crossover(indi2, rng)
	)
	if offspring1.Evaluated || offspring2.Evaluated {
		t.Error("Offsprings shouldn't have Evaluated set to True")
	}
	if &offspring1 == &indi1 || &offspring1 == &indi2 || &offspring2 == &indi1 || &offspring2 == &indi2 {
		t.Error("Offsprings shouldn't share pointers with parents")
	}
}

func TestMakeIndividuals(t *testing.T) {
	var rng = makeRandomNumberGenerator()
	for _, n := range []int{1, 2, 42} {
		var indis = makeIndividuals(n, MakeVector, rng)
		if len(indis) != n {
			t.Error("makeIndividuals didn't generate the right number of individuals")
		}
	}
}

func TestEvaluateIndividuals(t *testing.T) {
	var indis = makeIndividuals(10, MakeVector, makeRandomNumberGenerator())
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
		rng   = makeRandomNumberGenerator()
		indis = makeIndividuals(10, MakeVector, rng)
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
	var indis = makeIndividuals(10, MakeVector, makeRandomNumberGenerator())
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

func TestIndividualsSample(t *testing.T) {
	var (
		rng        = makeRandomNumberGenerator()
		indis      = makeIndividuals(10, MakeVector, rng)
		sampleSize = 3
		sample, _  = indis.sample(sampleSize, rng)
	)
	if len(sample) != sampleSize {
		t.Error("Wrong sample size")
	}
	// Check the sampled individuals come from the original population
	for _, a := range sample {
		var exists = false
		for _, b := range indis {
			if reflect.DeepEqual(a, b) {
				exists = true
			}
		}
		if exists == false {
			t.Error("Sampled individuals should come from the original population")
		}
	}
	// Check the sampled individuals have new references
	for _, a := range sample {
		var referenced = false
		for _, b := range indis {
			if &a == &b {
				referenced = true
			}
		}
		if referenced == true {
			t.Error("Sampled individuals shouln't share pointers with original population")
		}
	}
}

func TestGetFitnesses(t *testing.T) {
	var (
		indis = Individuals{
			Individual{nil, 0.0, false},
			Individual{nil, 1.0, false},
			Individual{nil, 2.0, false},
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
			Individual{nil, 1.0, false},
		}, 1.0},
		{Individuals{
			Individual{nil, 2.0, false},
			Individual{nil, 1.0, false},
		}, 1.0},
		{Individuals{
			Individual{nil, 1.0, false},
			Individual{nil, -1.0, false},
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
			Individual{nil, 1.0, false},
		}, 1.0},
		{Individuals{
			Individual{nil, 2.0, false},
			Individual{nil, 1.0, false},
		}, 2.0},
		{Individuals{
			Individual{nil, 1.0, false},
			Individual{nil, -1.0, false},
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
			Individual{nil, 1.0, false},
		}, 1.0},
		{Individuals{
			Individual{nil, 1.0, false},
			Individual{nil, 2.0, false},
		}, 1.5},
		{Individuals{
			Individual{nil, -1.0, false},
			Individual{nil, 1.0, false},
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
			Individual{nil, 1.0, false},
		}, 0.0},
		{Individuals{
			Individual{nil, -1.0, false},
			Individual{nil, 1.0, false},
		}, 1.0},
		{Individuals{
			Individual{nil, -2.0, false},
			Individual{nil, 2.0, false},
		}, 4.0},
	}
	for _, test := range testCases {
		if test.indis.FitStd() != math.Sqrt(test.variance) {
			t.Error("FitStd didn't work as expected")
		}
	}
}
