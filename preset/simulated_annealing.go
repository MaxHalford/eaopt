package preset

import "github.com/MaxHalford/gago2"

// SimAnn is a GA instance that mimicks a basic simulated annealing procedure.
func SimAnn(GenomeMaker gago2.GenomeMaker) gago2.GA {
	var ga = gago2.GA{
		MakeGenome: GenomeMaker,
		Topology: gago2.Topology{
			NbrPopulations: 1,
			NbrIndividuals: 1,
		},
		Model: gago2.ModSimAnn{
			T:     100,  // Starting temperature
			Tmin:  1,    // Stopping temperature
			Alpha: 0.99, // Decrease rate per iteration
		},
	}
	ga.Initialize()
	return ga
}
