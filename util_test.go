package gago

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestElementInSlice(t *testing.T) {
	var testCases = []struct {
		element interface{}
		slice   []interface{}
		in      bool
	}{
		{
			element: 1,
			slice:   uncastInts([]int{1, 2, 3}),
			in:      true,
		},
		{
			element: 4,
			slice:   uncastInts([]int{1, 2, 3}),
			in:      false,
		},
	}
	for _, test := range testCases {
		if elementInSlice(test.element, test.slice) != test.in {
			t.Error("elementInSlice is not behaving as expected")
		}
	}
}

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

func TestMakeLookup(t *testing.T) {
	var testCases = []struct {
		slice  []interface{}
		lookup map[interface{}]int
	}{
		{
			slice: uncastInts([]int{1, 2, 3}),
			lookup: map[interface{}]int{
				1: 0,
				2: 1,
				3: 2,
			},
		},
	}
	for _, test := range testCases {
		var lookup = makeLookup(test.slice)
		for k, v := range lookup {
			if v != test.lookup[k] {
				t.Error("createLookup didn't work as expected")
			}
		}
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
	for _, test := range testCases {
		if min(test[0], test[1]) != test[0] {
			t.Error("min didn't find the smallest integer")
		}
	}
}

func TestSum(t *testing.T) {
	var testCases = []struct {
		floats []float64
		total  float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{-1, 1}, 0},
		{[]float64{1.42, 42.1}, 43.52},
	}
	for _, test := range testCases {
		if sum(test.floats) != test.total {
			t.Error("sum didn't work as expected")
		}
	}
}

func TestMean(t *testing.T) {
	var testCases = []struct {
		floats []float64
		mean   float64
	}{
		{[]float64{1.0}, 1.0},
		{[]float64{1.0, 2.0}, 1.5},
		{[]float64{-1.0, 1.0}, 0.0},
	}
	for _, test := range testCases {
		if mean(test.floats) != test.mean {
			t.Error("mean didn't work as expected")
		}
	}
}

func TestVariance(t *testing.T) {
	var testCases = []struct {
		floats   []float64
		variance float64
	}{
		{[]float64{1.0}, 0.0},
		{[]float64{-1.0, 1.0}, 1.0},
		{[]float64{-2.0, 2.0}, 4.0},
	}
	for _, test := range testCases {
		if variance(test.floats) != test.variance {
			t.Error("variance didn't work as expected")
		}
	}
}

func TestDivide(t *testing.T) {
	var testCases = []struct {
		floats  []float64
		value   float64
		divided []float64
	}{
		{[]float64{1, 1}, 1, []float64{1, 1}},
		{[]float64{1, 1}, 2, []float64{0.5, 0.5}},
		{[]float64{42, -42}, 21, []float64{2, -2}},
	}
	for _, test := range testCases {
		var divided = divide(test.floats, test.value)
		for i := range divided {
			if divided[i] != test.divided[i] {
				t.Error("divided didn't work as expected")
			}
		}
	}
}

func TestCumsum(t *testing.T) {
	var testCases = []struct {
		floats []float64
		summed []float64
	}{
		{[]float64{1, 2, 3, 4}, []float64{1, 3, 6, 10}},
		{[]float64{-1, 0, 1}, []float64{-1, -1, 0}},
	}
	for _, test := range testCases {
		var summed = cumsum(test.floats)
		for i := range summed {
			if summed[i] != test.summed[i] {
				t.Error("cumsum didn't work as expected")
			}
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
	for _, test := range testCases {
		var ints = randomInts(test.k, test.min, test.max, rng)
		// Check the number of generated integers
		if len(ints) != test.k {
			t.Error("randomInts didn't generate the right number of integers")
		}
		// Check the bounds of each generated integer
		for _, integer := range ints {
			if integer < test.min || integer >= test.max {
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
