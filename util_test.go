package gago

import (
	"math"
	"math/rand"
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
	var (
		sizes = []int{1, 30, 10000}
		limit = math.Pow(1, -10)
	)
	for _, size := range sizes {
		var weights = randomWeights(size)
		// Test the length of the resulting slice
		if len(weights) != size {
			t.Error("Size problem with randomWeights")
		}
		// Test the elements in the slice sum up to 1
		var sum float64
		for _, weight := range weights {
			sum += weight
		}
		if math.Abs(sum-1.0) > limit {
			t.Error("Sum problem with randomWeights")
		}
	}
}

func TestMin(t *testing.T) {
	var testCases = [][]int{
		[]int{1, 2},
		[]int{2, 2},
		[]int{2, 3},
	}
	for _, pair := range testCases {
		if min(pair[0], pair[1]) != pair[0] {
			t.Error("min didn't find the smallest integer")
		}
	}
}

func TestMean(t *testing.T) {
	var testCases = []struct {
		values []float64
		mean   float64
	}{
		{[]float64{1.0}, 1.0},
		{[]float64{1.0, 2.0}, 1.5},
		{[]float64{-1.0, 1.0}, 0.0},
	}
	for _, testCase := range testCases {
		if mean(testCase.values) != testCase.mean {
			t.Error("mean didn't work as expected")
		}
	}
}

func TestVariance(t *testing.T) {
	var testCases = []struct {
		values   []float64
		variance float64
	}{
		{[]float64{1.0}, 0.0},
		{[]float64{-1.0, 1.0}, 1.0},
		{[]float64{-2.0, 2.0}, 4.0},
	}
	for _, testCase := range testCases {
		if variance(testCase.values) != testCase.variance {
			t.Error("variance didn't work as expected")
		}
	}
}

func TestRandomInts(t *testing.T) {
	var (
		src       = rand.NewSource(time.Now().UnixNano())
		rng       = rand.New(src)
		testCases = []struct {
			k, min, max int
		}{
			{1, 0, 1},
			{1, 0, 2},
			{2, 0, 2},
		}
	)
	for _, tc := range testCases {
		var ints = randomInts(tc.k, tc.min, tc.max, rng)
		// Check the number of generated integers
		if len(ints) != tc.k {
			t.Error("randomInts didn't generate the right number of integers")
		}
		// Check the bounds of each generated integer
		for _, integer := range ints {
			if integer < tc.min || integer >= tc.max {
				t.Error("randomInts didn't generate integers in the desired range")
			}
		}
		// Check the generated integers are unique
		for i, a := range ints {
			for j, b := range ints {
				if i != j && a == b {
					t.Error("randomInts didn't generate unique integers")
				}
			}
		}
	}
}
