package gago

import (
	"errors"
	"math/rand"
	"sort"
)

// Selector chooses a subset of size n from a group of individuals. The group of
// individuals a Selector is applied to is expected to be sorted.
type Selector interface {
	Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int)
	Validate() error
}

// SelElitism selection returns the n best individuals of a group.
type SelElitism struct{}

// Apply SelElitism.
func (sel SelElitism) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int) {
	var indexes = make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	return indis[:n], indexes
}

// Validate SelElitism fields.
func (sel SelElitism) Validate() error {
	return nil
}

// SelTournament samples individuals through tournament selection. The
// tournament is composed of randomly chosen individuals. The winner of the
// tournament is the chosen individual with the lowest fitness. The obtained
// individuals are not necessarily unique.
type SelTournament struct {
	NParticipants int
}

// Apply SelTournament.
func (sel SelTournament) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int) {
	var (
		selected = make(Individuals, n)
		indexes  = make([]int, n)
	)
	for i := range selected {
		var sample, roundIndexes = indis.sample(sel.NParticipants, rng)
		sample.SortByFitness()
		indexes[i] = roundIndexes[0]
		selected[i] = sample[0]
	}
	return selected, indexes
}

// Validate SelTournament fields.
func (sel SelTournament) Validate() error {
	if sel.NParticipants < 1 {
		return errors.New("NParticipants should be higher than 0")
	}
	return nil
}

// SelRoulette samples individuals through roulette wheel selection (also known
// as fitness proportionate selection).
type SelRoulette struct{}

func getWeights(fitnesses []float64) []float64 {
	var (
		n       = len(fitnesses)
		weights = make([]float64, n)
	)
	for i, v := range fitnesses {
		weights[i] = fitnesses[n-1] - v
	}
	return cumsum(divide(weights, sumFloat64s(weights)))
}

// Apply SelRoulette.
func (sel SelRoulette) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int) {
	var (
		selected = make(Individuals, n)
		indexes  = make([]int, n)
		weights  = getWeights(indis.getFitnesses())
	)
	for i := range selected {
		var (
			index  = sort.SearchFloat64s(weights, rand.Float64())
			winner = indis[index]
		)
		indexes[i] = index
		selected[i] = winner
	}
	return selected, indexes
}

// Validate SelRoulette fields.
func (sel SelRoulette) Validate() error {
	return nil
}
