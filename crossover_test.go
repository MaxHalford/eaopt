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
	{CrossPoint{NbPoints: 2}, InitUniformF{-5.0, 5.0}},
	{CrossUniformF{}, InitUniformF{-5.0, 5.0}},
	{CrossPMX{}, InitUniqueS{[]string{"A", "B", "C", "D"}}},
}

func TestCrossovers(t *testing.T) {
	var (
		src      = rand.NewSource(time.Now().UnixNano())
		rng      = rand.New(src)
		nbIndis  = 5
		nbGenes  = 4
		selector = SelTournament{2}
	)
	for _, c := range crossovers {
		var indis = makeIndividuals(nbIndis, nbGenes, rng)
		// Assign genomes
		for _, indi := range indis {
			c.init.apply(&indi, rng)
		}
		var (
			parents, _             = selector.Apply(2, indis, rng)
			offspring1, offspring2 = c.crossover.Apply(parents[0], parents[1], rng)
			offsprings             = Individuals{offspring1, offspring2}
		)
		// Check the offspring has a valid genome
		for _, offspring := range offsprings {
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
}
