package gago

import "testing"

func TestSelectionSize(t *testing.T) {
	var (
		rng       = makeRandomNumberGenerator()
		indis     = makeIndividuals(30, MakeVector, rng)
		selectors = []Selector{
			SelTournament{
				NbrParticipants: 3,
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
