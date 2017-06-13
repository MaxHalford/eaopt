package gago

import (
	"testing"
)

func TestDistanceMemoizer(t *testing.T) {
	var (
		dm = newDistanceMemoizer(l1Distance)
		a  = Individual{Genome: Vector{1, 1, 1}, ID: "1"}
		b  = Individual{Genome: Vector{3, 3, 3}, ID: "2"}
		c  = Individual{Genome: Vector{6, 6, 6}, ID: "3"}
	)
	// Check the number of calculations is initially 0
	if dm.nCalculations != 0 {
		t.Error("nCalculations is not initialized to 0")
	}
	// Check the distance between the 1st and itself
	if dm.GetDistance(a, a) != 0 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations is initially 0
	if dm.nCalculations != 0 {
		t.Error("nCalculations should not have increased")
	}
	// Check the distance between the 1st and the 2nd individual
	if dm.GetDistance(a, b) != 6 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations has increased by 1
	if dm.nCalculations != 1 {
		t.Error("nCalculations has not increased")
	}
	// Check the distance between the 2nd and the 1st individual
	if dm.GetDistance(b, a) != 6 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations has not increased
	if dm.nCalculations != 1 {
		t.Error("nCalculations has increased")
	}
	// Check the distance between the 1st and the 3rd individual
	if dm.GetDistance(a, c) != 15 {
		t.Error("Wrong calculated distance")
	}
	// Check the distance between the 1st and the 3rd individual
	if dm.GetDistance(b, c) != 9 {
		t.Error("Wrong calculated distance")
	}
}
