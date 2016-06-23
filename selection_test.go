package gago

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestTournament(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		size      = 3
		nbGenes   = 2
		indis     = make(Individuals, size)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var selector = SelTournament{size}
	var _, sample = selector.Apply(size, indis, generator)
	// Check the size of the sample
	if len(sample) != size {
		t.Error("Wrong sample size")
	}
	// Check the original population hasn't changed
	for i := range indis {
		if reflect.DeepEqual(indis[i], original[i]) == false {
			t.Error("Population has been modified")
		}
	}
	// Check the individual is from the initial population
	if reflect.DeepEqual(indis[0], sample[0]) == false {
		t.Error("Problem with tournament selection")
	}
}

func TestElitism(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbGenes   = 2
		size      = 3
		indis     = make(Individuals, size)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	var selector = SelElitism{}
	var _, sample = selector.Apply(size, indis, generator)
	// Check the size of the sample
	if len(sample) != size {
		t.Error("Wrong sample size")
	}
	// Check the original population hasn't changed
	for i := range indis {
		if reflect.DeepEqual(indis[i], original[i]) == false {
			t.Error("Population has been modified")
		}
	}
	// Check the individual is from the initial population
	if reflect.DeepEqual(indis[0], sample[0]) == false {
		t.Error("Problem with elitism selection")
	}
}

func TestTournamentAndElitism(t *testing.T) {
	var (
		source    = rand.NewSource(time.Now().UnixNano())
		generator = rand.New(source)
		nbGenes   = 2
		size      = 3
		indis     = make(Individuals, size)
	)
	for i := 0; i < size; i++ {
		indis[i] = makeIndividual(nbGenes)
		indis[i].Fitness = float64(i)
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var elitism = SelElitism{}
	var tournament = SelTournament{size}
	var _, a = elitism.Apply(size, indis, generator)
	var _, b = tournament.Apply(size, indis, generator)
	a.sort()
	b.sort()
	// Check the individual is from the initial population
	if a[0].Fitness != b[0].Fitness {
		t.Error("Elitism and full tournament selection differed")
	}
}
