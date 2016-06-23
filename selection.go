package gago

import "math/rand"

// Selector chooses a subset of size n from a group of individuals.
type Selector interface {
	Apply(n int, indis Individuals, generator *rand.Rand) ([]int, Individuals)
}

// SelTournament selection chooses an individual through tournament selection.
// The tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness.
type SelTournament struct {
	NbParticipants int
}

// Apply tournament selection.
func (sel SelTournament) Apply(n int, indis Individuals, generator *rand.Rand) ([]int, Individuals) {
	var (
		indexes = make([]int, n)
		winners = make(Individuals, n)
	)
	for i := range winners {
		// Sample the GA
		var roundIndexes, sample = indis.sample(sel.NbParticipants, generator)
		// The winner is the best individual participating in the tournament
		sample.sort()
		indexes[i] = roundIndexes[0]
		winners[i] = sample[0]
	}
	return indexes, winners
}

// SelElitism selection returns the best individuals in the GA.
type SelElitism struct{}

// Apply elitism selection.
func (sel SelElitism) Apply(n int, indis Individuals, generator *rand.Rand) ([]int, Individuals) {
	var indexes = make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	return indexes, indis[:n]
}
