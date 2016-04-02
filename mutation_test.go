package gago

import (
	"math/rand"
	"testing"
	"time"
)

var potentMutators = []Mutator{
	MutFNormal{Rate: 1.0, Std: 1},
	MutSplice{Rate: 1.0},
	MutPermute{Rate: 1.0},
}

var impotentMutators = []Mutator{
	MutFNormal{Rate: 0.0, Std: 1},
	MutSplice{Rate: 0.0},
	MutPermute{Rate: 0.0},
}

func TestPotentMutators(t *testing.T) {
	var (
		nbGenes   = 4
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		indi      = Individual{make([]interface{}, nbGenes), 0.0}
		init      = IFUniform{-5.0, 5.0}
	)
	init.apply(&indi, generator)
	var genome = make([]interface{}, len(indi.Genome))
	copy(genome, indi.Genome)
	// Probability to mutate is equal to 1
	for _, mutator := range potentMutators {
		mutator.apply(&indi, generator)
		// Check the genome is still the same size
		if len(indi.Genome) != nbGenes {
			t.Error("Genome size was changed after mutation")
		}
		// Check the genome has changed
		for i := range indi.Genome {
			if indi.Genome[i] == genome[i] {
				t.Error("Genome was not modified after potent mutation")
			}
		}
	}
}

func TestImpotentMutators(t *testing.T) {
	var (
		nbGenes   = 4
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		indi      = Individual{make([]interface{}, nbGenes), 0.0}
		init      = IFUniform{-5.0, 5.0}
	)
	init.apply(&indi, generator)
	var genome = make([]interface{}, len(indi.Genome))
	copy(genome, indi.Genome)
	// Probability to mutate is equal to 0
	for _, mutator := range impotentMutators {
		mutator.apply(&indi, generator)
		// Check the genome is still the same size
		if len(indi.Genome) != nbGenes {
			t.Error("Genome size was changed after mutation")
		}
		// Check the genome has changed
		for i := range indi.Genome {
			if indi.Genome[i] != genome[i] {
				t.Error("Genome was modified after impotent mutation")
			}
		}
	}
}
