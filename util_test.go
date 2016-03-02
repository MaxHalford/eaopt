package gago

import (
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetIndex(t *testing.T) {
	var test []interface{}
	test = append(test, 1)
	test = append(test, "イースター")
	// Integer in array
	if getIndex(1, test) != 0 {
		t.Error("Problem with getIndex")
	}
	// String in array
	if getIndex("イースター", test) != 1 {
		t.Error("Problem with getIndex")
	}
	// Element in array
	if getIndex("tamago", test) != -1 {
		t.Error("Problem with getIndex")
	}
}

func TestGenerateWeights(t *testing.T) {
	var sizes = []int{1, 30, 10000}
	var limit = math.Pow(1, -10)
	for _, size := range sizes {
		var weights = generateWeights(size)
		// Test the length of the resulting slice
		if len(weights) != size {
			t.Error("Size problem with generateWeights")
		}
		// Test the elements in the slice sum up to 1
		var sum float64
		for _, weight := range weights {
			sum += weight
		}
		if math.Abs(sum-1.0) > limit {
			t.Error("Sum problem with generateWeights")
		}
	}
}

func TestShuffleStrings(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var strings = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	var shuffled = shuffleStrings(strings, generator)
	// Check the shuffled slice is different from the original one
	if &strings == &shuffled {
		t.Error("Problem with shuffleIndividuals")
	}
}

func TestShuffleIndividuals(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var size = 10
	var indis = make(Individuals, size)
	for i := 0; i < size; i++ {
		indis[i] = Individual{make([]interface{}, 1), float64(i)}
	}
	var shuffled = shuffleIndividuals(indis, generator)
	// Check the shuffled slice is different from the original one
	if &indis == &shuffled {
		t.Error("Problem with shuffleIndividuals")
	}
}

func TestSampleIndividuals(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var nbGenes = 4
	var indis = Individuals{
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
		Individual{make([]interface{}, nbGenes), 0.0},
	}
	var size = 3
	var sample = sampleIndividuals(size, indis, generator)
	// Check the size of the sample
	if len(sample) != size {
		t.Error("Wrong sample size")
	}
	// Check the sampled individuals come from the original population
	for _, a := range sample {
		var exists = false
		for _, b := range indis {
			if reflect.DeepEqual(a, b) {
				exists = true
			}
		}
		if exists == false {
			t.Error("Problem with sampleIndividuals")
		}
	}
	// Check the sampled individuals have new references
	for _, a := range sample {
		var referenced = false
		for _, b := range indis {
			if &a == &b {
				referenced = true
			}
		}
		if referenced == true {
			t.Error("Problem with sampleIndividuals")
		}
	}
}
