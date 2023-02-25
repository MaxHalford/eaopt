package eaopt

import "testing"

func TestMigSizes(t *testing.T) {
	var (
		rng       = newRand()
		migrators = []Migrator{
			MigRing{
				NMigrants: 5,
			},
		}
	)
	for _, migrator := range migrators {
		for _, nPops := range []uint{2, 3, 10} {
			var fitnessMeans = make([]float64, nPops)
			for _, popSize := range []uint{6, 10, 30} {
				// Instantiate populations
				var pops = make([]Population, nPops)
				for i := range pops {
					pops[i] = newPopulation(popSize, false, NewVector, rng)
					pops[i].Individuals.Evaluate(false)
					fitnessMeans[i] = pops[i].Individuals.FitAvg()
				}
				migrator.Apply(pops, rng)
				// Check the Population sizes haven't changed
				for _, pop := range pops {
					if len(pop.Individuals) != int(popSize) {
						t.Error("Migration changed the Population sizes")
					}
				}
				// Check the average fitnesses have changed
				for i, pop := range pops {
					if pop.Individuals.FitAvg() == fitnessMeans[i] {
						t.Error("Average fitnesses didn't change")
					}
				}
			}
		}
	}
}

func TestMigRingValidate(t *testing.T) {
	var mig = MigRing{1}
	if err := mig.Validate(); err != nil {
		t.Error("Validation should not have raised error")
	}
	// Set NMigrants lower than 1
	mig.NMigrants = 0
	if err := mig.Validate(); err == nil {
		t.Error("Validation should raised error")
	}
}
