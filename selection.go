package gago

import "math/rand"

// Selector chooses a subset of size n from a group of individuals. The group of
// individuals a Selector is applied to is expected to be sorted.
type Selector interface {
	Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int)
}

// SelElitism selection returns the n best individuals of a group.
type SelElitism struct{}

// Apply elitism selection.
func (sel SelElitism) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int) {
	var indexes = make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	return indis[:n], indexes
}

// SelTournament selection chooses an individual through tournament selection.
// The tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness. The obtained
// individuals are not necessarily unique.
type SelTournament struct {
	NParticipants int
}

// Apply tournament selection.
func (sel SelTournament) Apply(n int, indis Individuals, rng *rand.Rand) (winners Individuals, indexes []int) {
	winners = make(Individuals, n)
	indexes = make([]int, n)
	for i := range winners {
		var sample, roundIndexes = indis.sample(sel.NParticipants, rng)
		sample.Sort()
		indexes[i] = roundIndexes[0]
		winners[i] = sample[0]
	}
	return
}
