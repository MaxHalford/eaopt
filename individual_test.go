package gago

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIndividualsSample(t *testing.T) {
	var (
		src       = rand.NewSource(time.Now().UnixNano())
		rng       = rand.New(src)
		nbIndis   = 5
		nbGenes   = 4
		indis     = makeIndividuals(nbIndis, nbGenes, rng)
		size      = 3
		_, sample = indis.sample(size, rng)
	)
	// Check the size of the sample
	if len(sample) != size {
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
			t.Error("Problem with sampleIndividuals")
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
			t.Error("Problem with sampleIndividuals")
		}
	}
}

func TestIndividualsSort(t *testing.T) {
	var (
		nbIndis = 5
		nbGenes = 4
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		indis   = makeIndividuals(nbIndis, nbGenes, rng)
	)
	// Assign fitness in decreasing order
	for i := range indis {
		indis[i].Fitness = float64(len(indis) - i)
	}
	// Sort
	indis.Sort()
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
			Individual{nil, 0.0, false, 0, "a"},
			Individual{nil, 1.0, false, 0, "b"},
			Individual{nil, 2.0, false, 0, "c"},
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

func TestFitnessMean(t *testing.T) {
	var testCases = []struct {
		indis Individuals
		mean  float64
	}{
		{Individuals{
			Individual{nil, 1.0, false, 0, "a"},
		}, 1.0},
		{Individuals{
			Individual{nil, 1.0, false, 0, "a"},
			Individual{nil, 2.0, false, 0, "b"},
		}, 1.5},
		{Individuals{
			Individual{nil, -1.0, false, 0, "a"},
			Individual{nil, 1.0, false, 0, "b"},
		}, 0.0},
	}
	for _, testCase := range testCases {
		if testCase.indis.FitnessMean() != testCase.mean {
			t.Error("FitnessMean didn't work as expected")
		}
	}
}

func TestFitnessVariance(t *testing.T) {
	var testCases = []struct {
		indis    Individuals
		variance float64
	}{
		{Individuals{
			Individual{nil, 1.0, false, 0, "a"},
		}, 0.0},
		{Individuals{
			Individual{nil, -1.0, false, 0, "a"},
			Individual{nil, 1.0, false, 0, "b"},
		}, 1.0},
		{Individuals{
			Individual{nil, -2.0, false, 0, "a"},
			Individual{nil, 2.0, false, 0, "b"},
		}, 4.0},
	}
	for _, testCase := range testCases {
		if testCase.indis.FitnessVar() != testCase.variance {
			t.Error("FitnessVar didn't work as expected")
		}
	}
}
