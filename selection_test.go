package gago

import "testing"

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
	if indexes[0] != 0 {
		t.Error("Full SelTournament didn't select the best individual")
	}
}

func TestGetWeights(t *testing.T) {
	var testCases = []struct {
		fitnesses []float64
		weights   []float64
	}{
		{[]float64{-10, -8, -5}, []float64{5 / 8, 1, 1}},
		{[]float64{-2, 0, 2, 3}, []float64{5 / 9, 8 / 9, 1, 1}},
		{[]float64{1, 2, 3, 4}, []float64{0.1, 0.3, 0.6, 1}},
	}
	for _, test := range testCases {
		var weights = getWeights(test.fitnesses)
		for i := range weights {
			if weights[i] != test.weights[i] {
				t.Error("getWeights didn't work as expected")
			}
		}
	}
}

func TestSpin(t *testing.T) {
	var testCases = []struct {
		value float64
		wheel []float64
		index int
	}{
		{0.1, []float64{0.3, 0.7, 1}, 0},
		{0.3, []float64{0.3, 0.7, 1}, 1},
		{0.8, []float64{0.3, 0.7, 1}, 2},
	}
	for _, test := range testCases {
		var index = spin(test.value, test.wheel)
		if index != test.index {
			t.Error("spin didn't work as expected")
		}
	}
}

func TestSelRoulette(t *testing.T) {
	var (
		rng   = makeRandomNumberGenerator()
		indis = makeIndividuals(30, MakeVector, rng)
		sel   = SelTournament{len(indis)}
	)
	for _, n := range []int{0, 1, 10, 30} {
		var selected, _ = sel.Apply(n, indis, rng)
		if len(selected) != n {
			t.Error("SelRoulette didn't select the right number of individuals")
		}
	}
}
