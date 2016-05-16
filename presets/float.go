package presets

import "github.com/MaxHalford/gago"

// Float returns a configuration for minimizing continuous mathematical
// functions with a given number of variables.
func Float(n int, function func([]float64) float64) gago.GA {
	return gago.GA{
		NbPopulations: 2,
		NbIndividuals: 30,
		NbGenes:       n,
		NbParents:     6,
		Ff: gago.FloatFunction{
			Image: function,
		},
		Initializer: gago.IFUniform{
			Lower: -1,
			Upper: 1,
		},
		Selector: gago.STournament{
			NbParticipants: 3,
		},
		Crossover: gago.CFUniform{},
		Mutator: gago.MutFNormal{
			Rate: 0.5,
			Std:  3,
		},
		MutRate:      0.5,
		Migrator:     gago.MigShuffle{},
		MigFrequency: 10,
	}
}
