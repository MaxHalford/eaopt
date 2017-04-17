package gago

import "testing"

func TestMigSizes(t *testing.T) {
	var (
		rng       = makeRandomNumberGenerator()
		migrators = []Migrator{
			MigRing{
				NMigrants: 5,
			},
		}
	)
	for _, migrator := range migrators {
		for _, nbrPops := range []int{2, 3, 10} {
			var fitnessMeans = make([]float64, nbrPops)
			for _, nbrIndis := range []int{6, 10, 30} {
				// Instantiate populations
				var pops = make([]Population, nbrPops)
				for i := range pops {
					pops[i] = makePopulation(nbrIndis, MakeVector, randString(3, rng))
					pops[i].Individuals.Evaluate()
					fitnessMeans[i] = pops[i].Individuals.FitAvg()
				}
				migrator.Apply(pops, rng)
				// Check the Population sizes haven't changed
				for _, pop := range pops {
					if len(pop.Individuals) != nbrIndis {
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
