package eaopt

import (
	"testing"
)

func TestMin(t *testing.T) {
	var testCases = [][]int{
		{1, 2},
		{2, 2},
		{2, 3},
	}
	for _, test := range testCases {
		if minInt(test[0], test[1]) != test[0] {
			t.Error("min didn't find the smallest integer")
		}
	}
}

func TestSumFloat64s(t *testing.T) {
	var testCases = []struct {
		floats []float64
		total  float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{-1, 1}, 0},
		{[]float64{1.42, 42.1}, 43.52},
	}
	for _, test := range testCases {
		if sumFloat64s(test.floats) != test.total {
			t.Error("sumFloat64s didn't work as expected")
		}
	}
}

func TestMinFloat64s(t *testing.T) {
	var testCases = []struct {
		floats []float64
		min    float64
	}{
		{[]float64{1.0}, 1.0},
		{[]float64{1.0, 2.0}, 1.0},
		{[]float64{-1.0, 1.0}, -1.0},
	}
	for _, test := range testCases {
		if minFloat64s(test.floats) != test.min {
			t.Error("meanFloat64s didn't work as expected")
		}
	}
}

func TestMaxFloat64s(t *testing.T) {
	var testCases = []struct {
		floats []float64
		max    float64
	}{
		{[]float64{1.0}, 1.0},
		{[]float64{1.0, 2.0}, 2.0},
		{[]float64{-1.0, 1.0}, 1.0},
	}
	for _, test := range testCases {
		if maxFloat64s(test.floats) != test.max {
			t.Error("maxFloat64s didn't work as expected")
		}
	}
}

func TestMeanFloat64s(t *testing.T) {
	var testCases = []struct {
		floats []float64
		mean   float64
	}{
		{[]float64{1.0}, 1.0},
		{[]float64{1.0, 2.0}, 1.5},
		{[]float64{-1.0, 1.0}, 0.0},
	}
	for _, test := range testCases {
		if meanFloat64s(test.floats) != test.mean {
			t.Error("meanFloat64s didn't work as expected")
		}
	}
}

func TestVarianceFloat64s(t *testing.T) {
	var testCases = []struct {
		floats   []float64
		variance float64
	}{
		{[]float64{1.0}, 0.0},
		{[]float64{-1.0, 1.0}, 1.0},
		{[]float64{-2.0, 2.0}, 4.0},
	}
	for _, test := range testCases {
		if varianceFloat64s(test.floats) != test.variance {
			t.Error("varianceFloat64s didn't work as expected")
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

func TestUnion(t *testing.T) {
	var testCases = []struct {
		x set
		y set
		u set
	}{
		{
			x: set{1: true, 2: true, 3: true},
			y: set{4: true, 5: true, 6: true},
			u: set{1: true, 2: true, 3: true, 4: true, 5: true, 6: true},
		},
		{
			x: set{1: true, 2: true, 3: true},
			y: set{2: true, 3: true, 4: true},
			u: set{1: true, 2: true, 3: true, 4: true},
		},
	}
	for _, test := range testCases {
		var u = union(test.x, test.y)
		for i := range u {
			if !test.u[i] {
				t.Error("union didn't work as expected")
			}
		}
	}
}
