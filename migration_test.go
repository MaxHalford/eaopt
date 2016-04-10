package gago

import "testing"

var PopulationSizes = []int{1, 2, 4}

func TestShuffle(t *testing.T) {
	for _, size := range PopulationSizes {
		// Instantiate a population
		var ga = GAFloat(2, nil)
		ga.Migrator = MigShuffle{}
		ga.NbPopulations = size
		// Apply the migration method
		ga.migrate()
		// Check the Population sizes haven't changed
		for _, pop := range ga.Populations {
			if len(pop.Individuals) != ga.NbIndividuals {
				t.Error("Shuffle migration changed the Population sizes")
			}
		}
	}
}
