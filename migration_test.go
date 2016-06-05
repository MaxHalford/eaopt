package gago

import "testing"

var (
	migrators = []Migrator{
		MigShuffle{},
	}
	initializer = InitUniformF{
		Lower: -1,
		Upper: 1,
	}
)

func TestMigSizes(t *testing.T) {
	var (
		populationSizes = []int{1, 2, 4}
		nbIndis         = []int{1, 2, 10}
	)
	for _, migrator := range migrators {
		for _, size := range populationSizes {
			for _, n := range nbIndis {
				// Instantiate populations
				var pops = make([]Population, size)
				for i := 0; i < size; i++ {
					pops[i] = makePopulation(n, 2, initializer)
				}
				// Apply the migration method
				migrator.Apply(pops)
				// Check the Population sizes haven't changed
				for _, pop := range pops {
					if len(pop.Individuals) != n {
						t.Error("Shuffle migration changed the Population sizes")
					}
				}
			}
		}
	}
}
