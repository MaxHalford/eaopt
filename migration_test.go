package gago

import "testing"

var demeSizes = []int{1, 2, 4}

func TestShuffle(t *testing.T) {
	for _, size := range demeSizes {
		// Instantiate a population
		var pop = Float
		pop.Migrator = MigShuffle{}
		pop.NbDemes = size
		// Apply the migration method
		pop.migrate()
		// Check the deme sizes haven't changed
		for _, deme := range pop.Demes {
			if len(deme.Individuals) != pop.NbIndividuals {
				t.Error("Shuffle migration changed the deme sizes")
			}
		}
	}
}
