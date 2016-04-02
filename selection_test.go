package gago

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestTournament(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var size = 3
	var indis = make(Individuals, size)
	for i := 0; i < size; i++ {
		indis[i] = Individual{make([]interface{}, 2), float64(i)}
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var selector = STournament{size}
	var sample = selector.apply(size, indis, generator)
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
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var size = 3
	var indis = make(Individuals, size)
	for i := 0; i < size; i++ {
		indis[i] = Individual{make([]interface{}, 2), float64(i)}
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	var selector = SElitism{}
	var sample = selector.apply(size, indis, generator)
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
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var size = 3
	var indis = make(Individuals, size)
	for i := 0; i < size; i++ {
		indis[i] = Individual{make([]interface{}, 2), float64(i)}
	}
	var original = make([]Individual, len(indis))
	copy(original, indis)
	// All the individuals participate in the tournament
	var elitism = SElitism{}
	var tournament = STournament{size}
	var a = elitism.apply(size, indis, generator)
	var b = tournament.apply(size, indis, generator)
	a.Sort()
	b.Sort()
	// Check the individual is from the initial population
	if a[0].Fitness != b[0].Fitness {
		t.Error("Elitism and full tournament selection differed")
	}
}
