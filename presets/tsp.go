package presets

import (
	"math/rand"

	"github.com/MaxHalford/gago"
)

type tspMutator struct{}

var (
	permute = gago.MutPermute{Max: 3}
	splice  = gago.MutSplice{}
)

func (tspmut tspMutator) Apply(indi *gago.Individual, rng *rand.Rand) {
	if rng.Float64() < 0.35 {
		permute.Apply(indi, rng)
	}
	if rng.Float64() < 0.45 {
		splice.Apply(indi, rng)
	}
}

// TSP returns a configuration for solving Travelling Salesman Problems given
// a corpus of positions that are associated to coordinates in the fitness
// function.
func TSP(places []string, distance func([]string) float64) gago.GA {
	return gago.GA{
		NbrPopulations: 1,
		NbrIndividuals: 100,
		NbrGenes:       len(places),
		Ff: gago.StringFunction{
			Image: distance,
		},
		Initializer: gago.InitUniqueS{
			Corpus: places,
		},
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NbParticipants: 17,
			},
			Crossover: gago.CrossPMX{},
			Mutator:   tspMutator{},
			MutRate:   1,
		},
		Migrator:     gago.MigShuffle{},
		MigFrequency: 10,
	}
}
