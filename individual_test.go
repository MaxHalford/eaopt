package gago

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestIndividualFloat(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var nbGenes = 4
	var indis = Individuals{
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
	}
	var ff = FloatFunction{func(X []float64) float64 {
		sum := 0.0
		for _, x := range X {
			sum += math.Abs(x)
		}
		return sum
	}}
	var init = FloatUniform{-5.0, 5.0}
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
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var nbGenes = 4
	var indis = Individuals{
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
	}
	var target = []string{"T", "E", "S", "T"}
	var ff = StringFunction{func(S []string) float64 {
		sum := 0.0
		for i := range S {
			if target[i] != S[i] {
				sum++
			}
		}
		return sum
	}}
	var alphabet = []string{"A", "B", "C", "D"}
	var init = StringUniform{alphabet}
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
