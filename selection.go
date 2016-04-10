package gago

import "math/rand"

// Selector chooses a sample of individual from a larger of individuals.
type Selector interface {
	// Apply select one individual
	apply(size int, individuals Individuals, generator *rand.Rand) Individuals
}

// STournament selection chooses an individual through tournament selection. The
// tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness.
type STournament struct {
	NbParticipants int
}

// Apply tournament selection.
func (ts STournament) apply(size int, indis Individuals, generator *rand.Rand) Individuals {
	var winners = make(Individuals, size)
	for i := range winners {
		// Sample the GA
		var sample = indis.sample(ts.NbParticipants, generator)
		// The winner is the best individual participating in the tournament
		sample.sort()
		winners[i] = sample[0]
	}
	return winners
}

// SElitism selection returns the best individuals in the GA.
type SElitism struct{}

// Apply elitism selection.
func (eli SElitism) apply(size int, indis Individuals, generator *rand.Rand) Individuals {
	return indis[:size]
}
