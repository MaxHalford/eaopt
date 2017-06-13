package gago

import "testing"

var presets = []GA{
	Generational(NewVector),
	SimulatedAnnealing(NewVector),
	HillClimbing(NewVector),
}

func TestPresetsValid(t *testing.T) {
	for _, preset := range presets {
		if preset.Validate() != nil {
			t.Error("The preset parameters are invalid")
		}
	}
}
