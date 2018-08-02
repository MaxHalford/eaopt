package eaopt

import (
	"fmt"
	"testing"
)

var (
	validSelectors = []Selector{
		SelElitism{},
		SelTournament{3},
		SelRoulette{},
	}
	invalidSelectors = []Selector{
		SelTournament{0},
	}
)

func TestSelectionSize(t *testing.T) {
	var (
		rng       = newRand()
		indis     = newIndividuals(30, NewVector, rng)
		selectors = []Selector{
			SelTournament{
				NContestants: 3,
			},
			SelElitism{},
		}
	)
	for _, selector := range selectors {
		for _, n := range []uint{3, 10, 20} {
			var selected, _, _ = selector.Apply(n, indis, rng)
			if len(selected) != int(n) {
				t.Error("Selector didn't select the expected number of individuals")
			}
		}
	}
}

func TestSelElitism(t *testing.T) {
	var (
		rng      = newRand()
		indis    = newIndividuals(30, NewVector, rng)
		selector = SelElitism{}
	)
	indis.Evaluate(false)
	for _, n := range []uint{1, 2, 10, 30} {
		var _, indexes, _ = selector.Apply(n, indis, rng)
		for i, index := range indexes {
			if index != i {
				t.Error("SelElitism didn't select the expected individuals")
			}
		}
	}
}

func TestSelTournament(t *testing.T) {
	var (
		rng   = newRand()
		indis = newIndividuals(30, NewVector, rng)
	)
	indis.Evaluate(false)
	var selected, _, _ = SelTournament{uint(len(indis))}.Apply(1, indis, rng)
	if selected[0].Fitness != indis.FitMin() {
		t.Error("Full SelTournament didn't select the best individual")
	}
}

func TestBuildWheel(t *testing.T) {
	var testCases = []struct {
		fitnesses []float64
		weights   []float64
	}{
		{[]float64{-10, -8, -5}, []float64{6.0 / 11, 10.0 / 11, 1}},
		{[]float64{-2, 0, 2, 3}, []float64{6.0 / 13, 10.0 / 13, 12.0 / 13, 1}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var weights = buildWheel(tc.fitnesses)
			for i := range weights {
				if weights[i] != tc.weights[i] {
					t.Error("buildWheel didn't work as expected")
				}
			}
		})
	}
}

func TestSelRoulette(t *testing.T) {
	var (
		rng   = newRand()
		indis = newIndividuals(30, NewVector, rng)
		sel   = SelRoulette{}
	)
	indis.Evaluate(false)
	for _, n := range []uint{0, 1, 10, 30} {
		var selected, _, _ = sel.Apply(n, indis, rng)
		if len(selected) != int(n) {
			t.Error("SelRoulette didn't select the right number of individuals")
		}
	}
}

// TestSelectorsValidate checks that each selector's Validate method doesn't
// return an error in case of a valid model and that it does for invalid models.
func TestSelectorsValidate(t *testing.T) {
	// Check valid selectors do not raise an error
	for _, sel := range validSelectors {
		if sel.Validate() != nil {
			t.Error("The selector validation should not have raised an error")
		}
	}
	// Check invalid selectors raise an error
	for _, sel := range invalidSelectors {
		if sel.Validate() == nil {
			t.Error("The selector validation should have raised error")
		}
	}
}
