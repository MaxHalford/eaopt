package gago2

import "math/rand"

// Selector chooses a subset of size n from a group of individuals.
type Selector interface {
	Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int)
}

// SelTournament selection chooses an individual through tournament selection.
// The tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness.
type SelTournament struct {
	NbrParticipants int
}

// Apply tournament selection.
func (sel SelTournament) Apply(n int, indis Individuals, rng *rand.Rand) (Individuals, []int) {
	var (
		indexes = make([]int, n)
		winners = make(Individuals, n)
	)
	for i := range winners {
		// Sample the GA
		var roundIndexes, sample = indis.sample(sel.NbrParticipants, rng)
		// The winner is the best individual participating in the tournament
		sample.Sort()
		indexes[i] = roundIndexes[0]
		winners[i] = sample[0]
	}
	return winners, indexes
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
