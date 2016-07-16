package gago

import "testing"

var (
	migrators = []Migrator{
		MigShuffle{},
	}
)

func TestMigSizes(t *testing.T) {
	var (
		initializer = InitUniformF{
			Lower: -1,
			Upper: 1,
		}
		populationSizes = []int{1, 2, 4}
		nbIndis         = []int{1, 2, 10}
		ff              = Float64Function{func(X []float64) float64 {
			sum := 0.0
			for _, x := range X {
				sum += x
			}
			return sum
		}}
	)
	for _, migrator := range migrators {
		for _, size := range populationSizes {
			for _, n := range nbIndis {
				// Instantiate populations
				var pops = make([]Population, size)
				for i := 0; i < size; i++ {
					pops[i] = makePopulation(n, 2, ff, initializer)
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
