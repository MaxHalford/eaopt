package gago

import "testing"

var presets = []GA{
	Generational(MakeVector),
	SimulatedAnnealing(MakeVector),
	HillClimbing(MakeVector),
}

func TestPresetsValid(t *testing.T) {
	for _, preset := range presets {
		if preset.Validate() != nil {
			t.Error("The preset parameters are invalid")
		}
	}
}
