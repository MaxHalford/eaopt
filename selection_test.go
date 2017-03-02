package gago

import (
	"fmt"
	"testing"
)

var (
	validSelectors = []Selector{
		SelElitism{},
		SelTournament{1},
		SelTournament{3},
		SelRoulette{},
	}
	invalidSelectors = []Selector{
		SelTournament{0},
		SelTournament{-1},
	}
)

func TestSelectionSize(t *testing.T) {
	var (
		rng       = makeRandomNumberGenerator()
		indis     = makeIndividuals(30, MakeVector, rng)
		selectors = []Selector{
			SelTournament{
				NParticipants: 3,
			},
			SelElitism{},
		}
	)
	for _, selector := range selectors {
		for _, n := range []int{1, 2, 10, 30} {
			var selected, _ = selector.Apply(n, indis, rng)
			if len(selected) != n {
				t.Error("Selector didn't select the expected number of individuals")
			}
		}
	}
}

func TestSelElitism(t *testing.T) {
	var (
		rng      = makeRandomNumberGenerator()
		indis    = makeIndividuals(30, MakeVector, rng)
		selector = SelElitism{}
	)
	indis.Evaluate()
	for _, n := range []int{1, 2, 10, 30} {
		var _, indexes = selector.Apply(n, indis, rng)
		for i, index := range indexes {
			if index != i {
				t.Error("SelElitism didn't select the expected individuals")
			}
		}
	}
}

func TestSelTournament(t *testing.T) {
	var (
		rng        = makeRandomNumberGenerator()
		indis      = makeIndividuals(30, MakeVector, rng)
		sel        = SelTournament{len(indis)}
		_, indexes = sel.Apply(1, indis, rng)
	)
	indis.Evaluate()
	if indexes[0] != 0 {
		t.Error("Full SelTournament didn't select the best individual")
	}
}

func TestGetWeights(t *testing.T) {
	var testCases = []struct {
		fitnesses []float64
		weights   []float64
	}{
		{[]float64{-10, -8, -5}, []float64{5.0 / 8, 1, 1}},
		{[]float64{-2, 0, 2, 3}, []float64{5.0 / 9, 8.0 / 9, 1, 1}},
	}
	for _, test := range testCases {
		var weights = getWeights(test.fitnesses)
		for i := range weights {
			if weights[i] != test.weights[i] {
				fmt.Println(weights[i], test.weights[i])
				t.Error("getWeights didn't work as expected")
			}
		}
	}
}

func TestSelRoulette(t *testing.T) {
	var (
		rng   = makeRandomNumberGenerator()
		indis = makeIndividuals(30, MakeVector, rng)
		sel   = SelRoulette{}
	)
	indis.Evaluate()
	for _, n := range []int{0, 1, 10, 30} {
		var selected, _ = sel.Apply(n, indis, rng)
		if len(selected) != n {
			t.Error("SelRoulette didn't select the right number of individuals")
		}
	}
}

// TestSelectorsValidate checks that each selector's Validate method doesn't
// return an error in case of a valid model and vice-versa
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
