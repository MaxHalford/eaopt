package gago

import "math/rand"

// Selector chooses a sample of individual from a larger of individuals.
type Selector interface {
	// Apply select one individual
	apply(size int, individuals Individuals, generator *rand.Rand) Individuals
}

// Tournament selection chooses an individual through tournament selection. The
// tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness.
type Tournament struct {
	NbParticipants int
}

// Apply tournament selection.
func (ts Tournament) apply(size int, indis Individuals, generator *rand.Rand) Individuals {
	var winners = make(Individuals, size)
	for i := range winners {
		// Sample the population
		var sample = sampleIndividuals(ts.NbParticipants, indis, generator)
		// The winner is the best individual participating in the tournament
		sample.Sort()
		winners[i] = sample[0]
	}
	return winners
}

// Elitism selection returns the best individuals in the population.
type Elitism struct{}

// Apply elitism selection.
func (eli Elitism) apply(size int, indis Individuals, generator *rand.Rand) Individuals {
	return indis[:size]
}
