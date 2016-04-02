package gago

import (
	"math"
	"math/rand"
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
