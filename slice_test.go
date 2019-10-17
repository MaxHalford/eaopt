package eaopt

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	var s = StringSlice{"イ", "ー", "ス", "タ", "ー"}
	if i, err := search("イ", s); i != 0 || err != nil {
		t.Error("Problem with search 1")
	}
	if i, err := search("ー", s); i != 1 || err != nil {
		t.Error("Problem with search 2")
	}
	if i, err := search("ス", s); i != 2 || err != nil {
		t.Error("Problem with search 3")
	}
	if i, err := search("タ", s); i != 3 || err != nil {
		t.Error("Problem with search 4")
	}
	if _, err := search("|", s); err == nil {
		t.Error("Problem with search 5")
	}
}

func TestNewIndexLookup(t *testing.T) {
	var testCases = []struct {
		slice  Slice
		lookup map[interface{}]int
	}{
		{
			slice: IntSlice{1, 2, 3},
			lookup: map[interface{}]int{
				1: 0,
				2: 1,
				3: 2,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var lookup = newIndexLookup(tc.slice)
			for k, v := range lookup {
				if v != tc.lookup[k] {
					t.Error("createLookup didn't work as expected")
				}
			}
		})
	}
}

func TestGetCycles(t *testing.T) {
	var testCases = []struct {
		x      []int
		y      []int
		cycles [][]int
	}{
		{
			x: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			y: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			cycles: [][]int{
				{0, 8, 3, 7},
				{1, 2, 6, 4},
				{5},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var cycles = getCycles(IntSlice(tc.x), IntSlice(tc.y))
			for i, cycle := range cycles {
				for j, c := range cycle {
					if c != tc.cycles[i][j] {
						t.Error("getCycles didn't work as expected")
					}
				}
			}
		})
	}
}

func TestGetNeighbours(t *testing.T) {
	var testCases = []struct {
		x          Slice
		neighbours map[interface{}]set
	}{
		{
			x: IntSlice{1, 2, 3, 4, 5, 6, 7, 8, 9},
			neighbours: map[interface{}]set{
				1: {9: true, 2: true},
				2: {1: true, 3: true},
				3: {2: true, 4: true},
				4: {3: true, 5: true},
				5: {4: true, 6: true},
				6: {5: true, 7: true},
				7: {6: true, 8: true},
				8: {7: true, 9: true},
				9: {8: true, 1: true},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var neighbours = getNeighbours(tc.x)
			for i, set := range neighbours {
				for j := range set {
					if !tc.neighbours[i][j] {
						t.Error("getNeighbours didn't work as expected")
					}
				}
			}
		})
	}
}

func TestIntSliceAt(t *testing.T) {
	var ints = IntSlice{1, 2, 3}
	if ints.At(0) != 1 || ints.At(1) != 2 || ints.At(2) != 3 {
		t.Error("IntSlice At method has unexpected behavior")
	}
}

func TestIntSliceLen(t *testing.T) {
	var ints = IntSlice{1, 2, 3}
	if ints.Len() != 3 {
		t.Error("IntSlice Len method has unexpected behavior")
	}
}

func TestIntSliceSwap(t *testing.T) {
	var ints = IntSlice{1, 2, 3}
	ints.Swap(0, 2)
	if ints.At(0) != 3 || ints.At(2) != 1 {
		t.Error("IntSlice Swap method has unexpected behavior")
	}
}

func TestIntSliceSlice(t *testing.T) {
	var ints = IntSlice{1, 2, 3}.Slice(1, 2)
	if ints.Len() != 1 || ints.At(0) != 2 {
		t.Error("IntSlice Slice method has unexpected behavior")
	}
}

func TestIntSliceSplit(t *testing.T) {
	var a, b = IntSlice{1, 2, 3}.Split(1)
	if a.Len() != 1 || b.Len() != 2 || a.At(0) != 1 || b.At(0) != 2 || b.At(1) != 3 {
		t.Error("IntSlice Split method has unexpected behavior")
	}
}

func TestIntSliceAppend(t *testing.T) {
	var ints = IntSlice{1}.Append(IntSlice{2})
	if ints.Len() != 2 || ints.At(0) != 1 || ints.At(1) != 2 {
		t.Error("IntSlice Append method has unexpected behavior")
	}
}

func TestIntSliceReplace(t *testing.T) {
	var ints = IntSlice{1}
	ints.Replace(IntSlice{2})
	if ints.Len() != 1 || ints.At(0) != 2 {
		t.Error("IntSlice Replace method has unexpected behavior")
	}
}

func TestIntSliceCopy(t *testing.T) {
	var (
		ints  = IntSlice{1}
		clone = ints.Copy()
	)
	clone.Replace(IntSlice{2})
	if ints.At(0) != 1 {
		t.Error("IntSlice Copy method has unexpected behavior")
	}
}

func TestFloat64SliceAt(t *testing.T) {
	var floats = Float64Slice{1, 2, 3}
	if floats.At(0) != 1.0 || floats.At(1) != 2.0 || floats.At(2) != 3.0 {
		t.Error("Float64Slice At method has unexpected behavior")
	}
}

func TestFloat64SliceLen(t *testing.T) {
	var floats = Float64Slice{1, 2, 3}
	if floats.Len() != 3 {
		t.Error("Float64Slice Len method has unexpected behavior")
	}
}

func TestFloat64SliceSwap(t *testing.T) {
	var floats = Float64Slice{1, 2, 3}
	floats.Swap(0, 2)
	if floats.At(0) != 3.0 || floats.At(2) != 1.0 {
		t.Error("Float64Slice Swap method has unexpected behavior")
	}
}

func TestFloat64SliceSlice(t *testing.T) {
	var floats = Float64Slice{1, 2, 3}.Slice(1, 2)
	if floats.Len() != 1 || floats.At(0) != 2.0 {
		t.Error("Float64Slice Slice method has unexpected behavior")
	}
}

func TestFloat64SliceSplit(t *testing.T) {
	var a, b = Float64Slice{1, 2, 3}.Split(1)
	if a.Len() != 1 || b.Len() != 2 || a.At(0) != 1.0 || b.At(0) != 2.0 || b.At(1) != 3.0 {
		t.Error("Float64Slice Split method has unexpected behavior")
	}
}

func TestFloat64SliceAppend(t *testing.T) {
	var floats = Float64Slice{1}.Append(Float64Slice{2})
	if floats.Len() != 2 || floats.At(0) != 1.0 || floats.At(1) != 2.0 {
		t.Error("Float64Slice Append method has unexpected behavior")
	}
}

func TestFloat64SliceReplace(t *testing.T) {
	var floats = Float64Slice{1}
	floats.Replace(Float64Slice{2})
	if floats.Len() != 1 || floats.At(0) != 2.0 {
		t.Error("Float64Slice Replace method has unexpected behavior")
	}
}

func TestFloat64SliceCopy(t *testing.T) {
	var (
		floats = Float64Slice{1}
		clone  = floats.Copy()
	)
	clone.Replace(Float64Slice{2})
	if floats.At(0) != 1.0 {
		t.Error("IntSlice Copy method has unexpected behavior")
	}
}

func TestStringSliceAt(t *testing.T) {
	var strings = StringSlice{"a", "b", "c"}
	if strings.At(0) != "a" || strings.At(1) != "b" || strings.At(2) != "c" {
		t.Error("StringSlice At method has unexpected behavior")
	}
}

func TestStringSliceLen(t *testing.T) {
	var strings = StringSlice{"a", "b", "c"}
	if strings.Len() != 3 {
		t.Error("StringSlice Len method has unexpected behavior")
	}
}

func TestStringSliceSwap(t *testing.T) {
	var strings = StringSlice{"a", "b", "c"}
	strings.Swap(0, 2)
	if strings.At(0) != "c" || strings.At(2) != "a" {
		t.Error("StringSlice Swap method has unexpected behavior")
	}
}

func TestStringSliceSlice(t *testing.T) {
	var strings = StringSlice{"a", "b", "c"}.Slice(1, 2)
	if strings.Len() != 1 || strings.At(0) != "b" {
		t.Error("StringSlice Slice method has unexpected behavior")
	}
}

func TestStringSliceSplit(t *testing.T) {
	var a, b = StringSlice{"a", "b", "c"}.Split(1)
	if a.Len() != 1 || b.Len() != 2 || a.At(0) != "a" || b.At(0) != "b" || b.At(1) != "c" {
		t.Error("StringSlice Split method has unexpected behavior")
	}
}

func TestStringSliceAppend(t *testing.T) {
	var strings = StringSlice{"a"}.Append(StringSlice{"b"})
	if strings.Len() != 2 || strings.At(0) != "a" || strings.At(1) != "b" {
		t.Error("StringSlice Append method has unexpected behavior")
	}
}

func TestStringSliceReplace(t *testing.T) {
	var strings = StringSlice{"a"}
	strings.Replace(StringSlice{"b"})
	if strings.Len() != 1 || strings.At(0) != "b" {
		t.Error("StringSlice Replace method has unexpected behavior")
	}
}

func TestStringSliceCopy(t *testing.T) {
	var (
		strings = StringSlice{"a"}
		clone   = strings.Copy()
	)
	clone.Replace(StringSlice{"b"})
	if strings.At(0) == "b" {
		t.Error("StringSlice Copy method has unexpected behavior")
	}
}
