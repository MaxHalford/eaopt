package gago

import "testing"

func TestCastInts(t *testing.T) {
	var (
		ints     = []int{1, 2, 3}
		uncasted = uncastInts(ints)
		casted   = castInts(uncasted)
	)
	// Check lengths
	if len(uncasted) != len(ints) {
		t.Error("Uncasted slice does not have the right size")
	}
	// Check values
	for i, v := range casted {
		if v != ints[i] {
			t.Error("Casted slice does not contain the right values")
		}
	}
}

func TestCastFloat64s(t *testing.T) {
	var (
		floats   = []float64{1, 2, 3}
		uncasted = uncastFloat64s(floats)
		casted   = castFloat64s(uncasted)
	)
	// Check lengths
	if len(uncasted) != len(floats) {
		t.Error("Uncasted slice does not have the right size")
	}
	// Check values
	for i, v := range casted {
		if v != floats[i] {
			t.Error("Casted slice does not contain the right values")
		}
	}
}

func TestCastStrings(t *testing.T) {
	var (
		strings  = []string{"a", "b", "c"}
		uncasted = uncastStrings(strings)
		casted   = castStrings(uncasted)
	)
	// Check lengths
	if len(uncasted) != len(strings) {
		t.Error("Uncasted slice does not have the right size")
	}
	// Check values
	for i, v := range casted {
		if v != strings[i] {
			t.Error("Casted slice does not contain the right values")
		}
	}
}
