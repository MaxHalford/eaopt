package gago

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// Selector chooses a subset of size n from a group of individuals. The group of
// individuals a Selector is applied to is expected to be sorted.
type Selector interface {
	Apply(n int, indis Individuals, rng *rand.Rand) (selected Individuals, indexes []int, err error)
	Validate() error
}

// SelElitism selection returns the n best individuals of a group.
type SelElitism struct{}

// Apply SelElitism.
func (sel SelElitism) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int, error) {
	return indis[:n].Clone(rng), newInts(n), nil
}

// Validate SelElitism fields.
func (sel SelElitism) Validate() error {
	return nil
}

// SelTournament samples individuals through tournament selection. The
// tournament is composed of randomly chosen individuals. The winner of the
// tournament is the chosen individual with the lowest fitness. The obtained
// individuals are all distinct.
type SelTournament struct {
	NContestants int
}

// Apply SelTournament.
func (sel SelTournament) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int, error) {
	// Check that the number of individuals is large enough
	if len(indis)-n < sel.NContestants-1 {
		return nil, nil, fmt.Errorf("Not enough individuals to select %d "+
			"with NContestants = %d, have %d individuals and need at least %d",
			n, sel.NContestants, len(indis), sel.NContestants+n-1)
	}
	var (
		winners         = make(Individuals, n)
		indexes         = make([]int, n)
		notSelectedIdxs = newInts(len(indis))
	)
	for i := range winners {
		// Sample contestants
		var (
			contestants, idxs, _ = sampleInts(notSelectedIdxs, sel.NContestants, rng)
			winnerIdx            int
		)
		// Find the best contestant
		winners[i].Fitness = math.Inf(1)
		for j, k := range contestants {
			if indis[k].GetFitness() < winners[i].Fitness {
				winners[i] = indis[k]
				indexes[i] = k
				winnerIdx = idxs[j]
			}
		}
		// Ban the winner from re-participating
		notSelectedIdxs = append(notSelectedIdxs[:winnerIdx], notSelectedIdxs[winnerIdx+1:]...)
	}
	return winners, indexes, nil
}

// Validate SelTournament fields.
func (sel SelTournament) Validate() error {
	if sel.NContestants < 1 {
		return errors.New("NContestants should be higher than 0")
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
func (sel SelRoulette) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int, error) {
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
	return selected, indexes, nil
}

// Validate SelRoulette fields.
func (sel SelRoulette) Validate() error {
	return nil
}
