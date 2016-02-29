package gago

import (
	"math"
	"testing"
)

func TestFindPosition(t *testing.T) {
	var test []interface{}
	test = append(test, 1)
	test = append(test, "イースター")
	// Integer in array
	if findPosition(test, 1) != 0 {
		t.Error("Problem with findPosition")
	}
	// String in array
	if findPosition(test, "イースター") != 1 {
		t.Error("Problem with findPosition")
	}
	// Element in array
	if findPosition(test, "tamago") != -1 {
		t.Error("Problem with findPosition")
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
