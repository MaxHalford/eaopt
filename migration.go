package gago

import (
	"errors"
	"math/rand"
)

// Migrator applies crossover to the GA level, as such it doesn't
// require an independent random number generator and can use the global one.
type Migrator interface {
	Apply(pops Populations, rng *rand.Rand)
	Validate() error
}

// MigRing migration exchanges individuals between consecutive Populations in a
// random fashion. One by one, each population exchanges NMigrants individuals
// at random with the next population. NMigrants should be higher than the
// number of individuals in each population, else all the individuals will
// migrate and it will be as if nothing happened.
type MigRing struct {
	NMigrants int // Number of migrants per exchange between Populations
}

// Apply MigRing.
func (mig MigRing) Apply(pops Populations, rng *rand.Rand) {
	for i := 0; i < len(pops)-1; i++ {
		for _, k := range randomInts(mig.NMigrants, 0, len(pops[i].Individuals), rng) {
			pops[i].Individuals[k], pops[i+1].Individuals[k] = pops[i+1].Individuals[k], pops[i].Individuals[k]
		}
	}
}

// Validate MigRing fields.
func (mig MigRing) Validate() error {
	if mig.NMigrants < 1 {
		return errors.New("NMigrants should be higher than 0")
	}
	return nil
}
