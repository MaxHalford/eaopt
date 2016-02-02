package gago

import "math/rand"

// Shuffle exchanges individuals between demes in a random fashion.
func Shuffle(demes []Deme) []Deme {
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
