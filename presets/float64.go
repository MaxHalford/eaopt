package presets

import "github.com/MaxHalford/gago"

// Float64 returns a configuration for minimizing continuous mathematical
// functions with a given number of variables.
func Float64(n int, function func([]float64) float64) gago.GA {
	return gago.GA{
		Ff: gago.Float64Function{
			Image: function,
		},
		Initializer: gago.InitUniformF{
			Lower: -1,
			Upper: 1,
		},
		Topology: gago.Topology{
			NbrPopulations: 2,
			NbrIndividuals: 30,
			NbrGenes:       n,
		},
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NbrParticipants: 3,
			},
			Crossover: gago.CrossUniformF{},
			Mutator: gago.MutNormalF{
				Rate: 0.5,
				Std:  3,
			},
			MutRate: 0.5,
		},
		Migrator:     gago.MigShuffle{},
		MigFrequency: 10,
	}
}
