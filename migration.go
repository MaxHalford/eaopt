package gago

import "math/rand"

// Migrator applies crossover to the population level, as such it doesn't
// require an independent random number generator and can use the global one.
type Migrator interface {
	Apply([]Deme) []Deme
}

// Shuffle migration exchanges individuals between demes in a random fashion.
type Shuffle struct{}

// Apply shuffle migration
func (shuffle Shuffle) Apply(demes []Deme) []Deme {
	for i := 0; i < len(demes); i++ {
		for j := i + 1; j < len(demes); j++ {
			// Choose where to split the individuals
			var split = rand.Intn(demes[i].Size)
			// Create a temporary slice of individuals in order to switch
			var tmp = make([]Individual, demes[i].Size)
			copy(tmp, demes[i].Individuals)
			// Perform the switch
			demes[i].Individuals = append(demes[i].Individuals[:split], demes[j].Individuals[split:]...)
			demes[j].Individuals = append(demes[j].Individuals[:split], tmp[split:]...)
		}
	}
	return demes
}

// String description of shuffle migration.
func (shuffle Shuffle) String() string {
	return "Shuffle migration."
}
