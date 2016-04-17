package gago

import (
	"math/rand"
	"testing"
	"time"
)

var potentMutators = []Mutator{
	MutFNormal{
		Rate: 1,
		Std:  1,
	},
	MutSplice{
		Rate: 1,
	},
	MutPermute{
		Rate: 1,
		Max:  3,
	},
	MutSUniform{
		Rate:   1,
		Corpus: []string{"t", "e", "s"},
	},
}

var impotentMutators = []Mutator{
	MutFNormal{
		Rate: 0,
		Std:  1,
	},
	MutSplice{
		Rate: 0,
	},
	MutPermute{
		Rate: 0,
		Max:  1,
	},
	MutSUniform{
		Rate:   0,
		Corpus: []string{"t", "e", "s"},
	},
}

func TestPotentMutators(t *testing.T) {
	var (
		nbGenes   = 4
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		indi      = makeIndividual(nbGenes)
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
		// Check the number of differences
		var differences int
		for i := range genome {
			if genome[i] != indi.Genome[i] {
				differences++
			}
		}
		if differences == 0 {
			t.Error("Mutator should have worked and didn't")
		}
	}
}

func TestImpotentMutators(t *testing.T) {
	var (
		nbGenes   = 4
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		indi      = makeIndividual(nbGenes)
		init      = IFUniform{-5.0, 5.0}
	)
	init.apply(&indi, generator)
	var genome = make([]interface{}, len(indi.Genome))
	copy(genome, indi.Genome)
	// Probability to mutate is equal to 1
	for _, mutator := range impotentMutators {
		mutator.apply(&indi, generator)
		// Check the genome is still the same size
		if len(indi.Genome) != nbGenes {
			t.Error("Genome size was changed after mutation")
		}
		// Check the number of differences
		var differences int
		for i := range genome {
			if genome[i] != indi.Genome[i] {
				differences++
			}
		}
		if differences != 0 {
			t.Error("Mutator did work and shouldn't have")
		}
	}
}
