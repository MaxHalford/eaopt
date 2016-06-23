package gago

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIndividualsSample(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbIndis   = 5
		nbGenes   = 4
		indis     = makeIndividuals(nbIndis, nbGenes)
		size      = 3
		_, sample = indis.sample(size, generator)
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
		indis   = makeIndividuals(nbIndis, nbGenes)
	)
	// Assign fitness in decreasing order
	for i := range indis {
		indis[i].Fitness = float64(len(indis) - i)
	}
	// Sort
	indis.sort()
	// Check fitnesses are in increasing order
	for i := 1; i < len(indis); i++ {
		if indis[i-1].Fitness > indis[i].Fitness {
			t.Error("Individuals are not sorted")
		}
	}
}
