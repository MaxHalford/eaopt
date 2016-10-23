package preset

import "github.com/MaxHalford/gago"

// SimAnn is a GA instance that mimicks a basic simulated annealing procedure.
func SimAnn(GenomeMaker gago.GenomeMaker) gago.GA {
	var ga = gago.GA{
		MakeGenome: GenomeMaker,
		Topology: gago.Topology{
			NbrPopulations: 1,
			NbrIndividuals: 1,
		},
		Model: gago.ModSimAnn{
			T:     100,  // Starting temperature
			Tmin:  1,    // Stopping temperature
			Alpha: 0.99, // Decrease rate per iteration
		},
	}
	ga.Initialize()
	return ga
}
