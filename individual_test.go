package gago

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIndividualFloat(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbGenes   = 4
		indis     = Individuals{
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
		}
		ff = FloatFunction{func(X []float64) float64 {
			sum := 0.0
			for _, x := range X {
				sum += math.Abs(x)
			}
			return sum
		}}
		init = IFUniform{-5.0, 5.0}
	)
	// Assign genomes and fitnesses
	for i, indi := range indis {
		init.apply(&indi, generator)
		indis[i].Evaluate(ff)
	}
	// Check if fitnesses have been assigned
	for _, indi := range indis {
		if indi.Fitness == 0.0 {
			t.Error("No fitness was assigned")
		}
	}
	// Check if the individuals are sorted
	indis.Sort()
	for i := 0; i < len(indis); i++ {
		for j := i + 1; j < len(indis); j++ {
			if indis[i].Fitness > indis[j].Fitness {
				t.Error("Individuals are not sorted")
			}
		}
	}
}

func TestIndividualString(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbGenes   = 4
		indis     = Individuals{
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
		}
		target   = []string{"T", "E", "S", "T"}
		sum      = 0.0
		alphabet = []string{"A", "B", "C", "D"}
		init     = ISUniform{alphabet}
		ff       = StringFunction{func(S []string) float64 {
			for i := range S {
				if target[i] != S[i] {
					sum++
				}
			}
			return sum
		}}
	)

	// Assign genomes and fitnesses
	for i, indi := range indis {
		init.apply(&indi, generator)
		indis[i].Evaluate(ff)
	}
	// Check if fitnesses have been assigned
	for _, indi := range indis {
		if indi.Fitness == 0.0 {
			t.Error("No fitness was assigned")
		}
	}
	// Check if the individuals are sorted
	indis.Sort()
	for i := 0; i < len(indis); i++ {
		for j := i + 1; j < len(indis); j++ {
			if indis[i].Fitness > indis[j].Fitness {
				t.Error("Individuals are not sorted")
			}
		}
	}
}

func TestShuffleIndividuals(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		size      = 10
		indis     = make(Individuals, size)
	)
	for i := 0; i < size; i++ {
		indis[i] = Individual{make([]interface{}, 1), float64(i)}
	}
	var shuffled = indis.shuffle(generator)
	// Check the shuffled slice is different from the original one
	if &shuffled == &indis {
		t.Error("Problem with shuffleIndividuals")
	}
}

func TestSampleIndividuals(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbGenes   = 4
		indis     = Individuals{
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
		}
		size   = 3
		sample = indis.sample(size, generator)
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
