package gago

import (
	"math/rand"
	"testing"
	"time"
)

var mutators = []Mutator{
	MutFNormal{
		Rate: 1,
		Std:  1,
	},
	MutSplice{},
	MutPermute{
		Max: 3,
	},
	MutSUniform{
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
	for _, mutator := range mutators {
		mutator.Apply(&indi, generator)
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
