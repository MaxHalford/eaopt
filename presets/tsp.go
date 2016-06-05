package presets

import (
	"math/rand"

	"github.com/MaxHalford/gago"
)

type tspMutator struct{}

var permute = gago.MutPermute{Max: 3}
var splice = gago.MutSplice{}

func (tspmut tspMutator) Apply(indi *gago.Individual, generator *rand.Rand) {
	if generator.Float64() < 0.6 {
		permute.Apply(indi, generator)
	}
	if generator.Float64() < 0.3 {
		splice.Apply(indi, generator)
	}
}

// TSP returns a configuration for solving Travelling Salesman Problems given
// a corpus of positions that are associated to coordinates in the fitness
// function.
func TSP(places []string, distance func([]string) float64) gago.GA {
	return gago.GA{
		NbPopulations: 1,
		NbIndividuals: 30,
		NbGenes:       len(places),
		NbParents:     3,
		Ff: gago.StringFunction{
			Image: distance,
		},
		Initializer: gago.InitUniqueS{
			Corpus: places,
		},
		Selector:     gago.SelElitism{},
		Crossover:    gago.CrossPMX{},
		Mutator:      tspMutator{},
		MutRate:      1,
		Migrator:     gago.MigShuffle{},
		MigFrequency: 10,
	}
}
