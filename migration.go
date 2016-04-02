package gago

import "math/rand"

// Migrator applies crossover to the population level, as such it doesn't
// require an independent random number generator and can use the global one.
type Migrator interface {
	apply([]Deme) []Deme
}

// MigShuffle migration exchanges individuals between demes in a random fashion.
type MigShuffle struct{}

// Apply shuffle migration.
func (shuffle MigShuffle) apply(demes []Deme) []Deme {
	for i := 0; i < len(demes); i++ {
		for j := i + 1; j < len(demes); j++ {
			// Choose where to split the individuals
			var split = rand.Intn(len(demes[i].Individuals))
			// Create a temporary slice of individuals in order to switch
			var tmp = make([]Individual, len(demes[i].Individuals))
			copy(tmp, demes[i].Individuals)
			// Perform the switch
			demes[i].Individuals = append(tmp[:split], demes[j].Individuals[split:]...)
			demes[j].Individuals = append(demes[j].Individuals[:split], tmp[split:]...)
		}
	}
	return demes
}
