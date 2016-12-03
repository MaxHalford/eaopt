package gago

// Generational returns a GA instance that uses the generational model.
func Generational(MakeGenome GenomeMaker) GA {
	var ga = GA{
		MakeGenome: MakeGenome,
		Topology: Topology{
			NPopulations: 2,
			NIndividuals: 50,
		},
		Model: ModGenerational{
			Selector: SelTournament{
				NParticipants: 3,
			},
			MutRate: 0.5,
		},
	}
	ga.Initialize()
	return ga
}

// SimulatedAnnealing returns a GA instance that mimicks a basic simulated
// annealing procedure.
func SimulatedAnnealing(MakeGenome GenomeMaker) GA {
	var ga = GA{
		MakeGenome: MakeGenome,
		Topology: Topology{
			NPopulations: 1,
			NIndividuals: 1,
		},
		Model: ModSimAnn{
			T:     100,  // Starting temperature
			Tmin:  1,    // Stopping temperature
			Alpha: 0.99, // Decrease rate per iteration
		},
	}
	ga.Initialize()
	return ga
}

// HillClimbing returns a GA instance that mimicks a basic hill-climbing
// procedure.
func HillClimbing(MakeGenome GenomeMaker) GA {
	var ga = GA{
		MakeGenome: MakeGenome,
		Topology: Topology{
			NPopulations: 1,
			NIndividuals: 1,
		},
		Model: ModMutationOnly{
			NChosen:  1,
			Selector: SelElitism{},
			Strict:   true,
		},
	}
	ga.Initialize()
	return ga
}
