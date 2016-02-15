package gago

import "math/rand"

// Selector chooses an individual from a group of individuals.
type Selector interface {
	// Apply select one individual
	apply(individuals Individuals, generator *rand.Rand) Individual
}

// Tournament selection chooses an individual through tournament selection. The
// tournament is composed of randomly chosen individuals. The winner of the
// tournament is the individual with the lowest fitness.
type Tournament struct {
	NbParticipants int
}

// Apply tournament selection.
func (ts Tournament) apply(indis Individuals, generator *rand.Rand) Individual {
	// Randomly sample the population
	sample := make(Individuals, ts.NbParticipants)
	for j := 0; j < ts.NbParticipants; j++ {
		index := generator.Intn(len(indis))
		sample[j] = indis[index]
	}
	// The winner is the best individual participating in the tournament
	sample.Sort()
	winner := sample[0]
	return winner
}
