package gago

import "math/rand"

// Migrator applies crossover to the GA level, as such it doesn't
// require an independent random number generator and can use the global one.
type Migrator interface {
	Apply([]Population)
}

// MigShuffle migration exchanges individuals between Populations in a random fashion.
type MigShuffle struct{}

// Apply shuffle migration.
func (mig MigShuffle) Apply(pops []Population) {
	for i := 0; i < len(pops); i++ {
		for j := i + 1; j < len(pops); j++ {
			// Choose where to split the individuals
			var split = rand.Int() % len(pops[i].Individuals)
			// Create a temporary slice of individuals in order to switch
			var tmp = make([]Individual, len(pops[i].Individuals))
			copy(tmp, pops[i].Individuals)
			// Perform the switch
			pops[i].Individuals = append(tmp[:split], pops[j].Individuals[split:]...)
			pops[j].Individuals = append(pops[j].Individuals[:split], tmp[split:]...)
		}
	}
}
