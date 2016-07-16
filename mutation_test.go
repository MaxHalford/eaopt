package gago

import (
	"math/rand"
	"testing"
	"time"
)

var mutators = []Mutator{
	MutNormalF{
		Rate: 1,
		Std:  1,
	},
	MutSplice{},
	MutPermute{
		Max: 3,
	},
	MutUniformS{
		Corpus: []string{"t", "e", "s"},
	},
}

func TestPotentMutators(t *testing.T) {
	var (
		nbGenes = 4
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		indi    = makeIndividual(nbGenes, rng)
		init    = InitUniformF{-5.0, 5.0}
	)
	init.apply(&indi, rng)
	var genome = make([]interface{}, len(indi.Genome))
	copy(genome, indi.Genome)
	// Probability to mutate is equal to 1
	for _, mutator := range mutators {
		mutator.Apply(&indi, rng)
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
