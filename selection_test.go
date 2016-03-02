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
	var copy = indis
	var selector = Tournament{3}
	var indi = selector.apply(indis, generator)
	// Check the original population hasn't changed
	if reflect.DeepEqual(indis, copy) == false {
		t.Error("Population has been modified")
	}
	// Check the individual is from the initial population
	if reflect.DeepEqual(indi, indis[0]) == false {
		t.Error("Problem with tournament selection")
	}
}
