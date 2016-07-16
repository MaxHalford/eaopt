package gago

import (
	"math/rand"
	"testing"
	"time"
)

func TestTournament(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		size    = 3
		nbGenes = 2
		indis   = make(Individuals, size)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes, rng)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var (
		selector  = SelTournament{size}
		sample, _ = selector.Apply(size, indis, rng)
	)
	// Check the size of the sample
	if len(sample) != size {
		t.Error("Wrong sample size")
	}
	// Check the original population hasn't changed
	for i := range indis {
		if indis[i].Name != original[i].Name {
			t.Error("Population has been modified")
		}
	}
	// Check the individual is from the initial population
	if indis[0].Name != sample[0].Name {
		t.Error("Problem with tournament selection")
	}
}

func TestElitism(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		nbGenes = 2
		size    = 3
		indis   = makeIndividuals(size, nbGenes, rng)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes, rng)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	var (
		selector  = SelElitism{}
		sample, _ = selector.Apply(size, indis, rng)
	)
	// Check the size of the sample
	if len(sample) != size {
		t.Error("Wrong sample size")
	}
	// Check the original population hasn't changed
	for i := range indis {
		if indis[i].Name != original[i].Name {
			t.Error("Population has been modified")
		}
	}
	// Check the individual is from the initial population
	if indis[0].Name != sample[0].Name {
		t.Error("Problem with elitism selection")
	}
}

func TestTournamentAndElitism(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		nbGenes = 2
		size    = 3
		indis   = makeIndividuals(size, nbGenes, rng)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes, rng)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var (
		elitism    = SelElitism{}
		tournament = SelTournament{size}
		elite, _   = elitism.Apply(size, indis, rng)
		tourney, _ = tournament.Apply(size, indis, rng)
	)
	elite.sort()
	tourney.sort()
	// Check the individual is from the initial population
	if elite[0].Fitness != tourney[0].Fitness {
		t.Error("Elitism and full tournament selection differed")
	}
}
