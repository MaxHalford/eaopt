package gago

import (
	"math/rand"
	"testing"
	"time"
)

var crossovers = []struct {
	crossover Crossover
	init      Initializer
}{
	{FloatParenthood{}, FloatUniform{-5.0, 5.0}},
	{FloatFitnessProportionate{2}, FloatUniform{-5.0, 5.0}},
	{PartiallyMappedCrossover{}, StringUnique{[]string{"A", "B", "C", "D"}}},
}

func TestCrossovers(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var nbGenes = 4
	var selector = Tournament{2}
	for _, c := range crossovers {
		var indis = Individuals{
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
			Individual{make([]interface{}, nbGenes), 0.0},
		}
		// Assign genomes
		for _, indi := range indis {
			c.init.apply(&indi, generator)
		}
		var offspring = c.crossover.apply(selector, indis, generator)
		// Check the offspring has a valid genome
		if len(offspring.Genome) != nbGenes {
			t.Error("Crossover operator generated genome of invalid length")
		}
		// Check the offspring doesn't share a pointer with the original population
		for _, indi := range indis {
			if &offspring == &indi {
				t.Error("The offspring shares a pointer with another individual")
			}
		}
	}
}
