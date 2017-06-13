package gago

// Generational returns a GA instance that uses the generational model.
func Generational(GenomeFactory GenomeFactory) GA {
	return GA{
		GenomeFactory: GenomeFactory,
		NPops:         2,
		PopSize:       50,
		Model: ModGenerational{
			Selector: SelTournament{
				NParticipants: 3,
			},
			MutRate: 0.5,
		},
	}
}

// SimulatedAnnealing returns a GA instance that mimicks a basic simulated
// annealing procedure.
func SimulatedAnnealing(GenomeFactory GenomeFactory) GA {
	return GA{
		GenomeFactory: GenomeFactory,
		NPops:         1,
		PopSize:       1,
		Model: ModSimAnn{
			T:     100,  // Starting temperature
			Tmin:  1,    // Stopping temperature
			Alpha: 0.99, // Decrease rate per iteration
		},
	}
}

// HillClimbing returns a GA instance that mimicks a basic hill-climbing
// procedure.
func HillClimbing(GenomeFactory GenomeFactory) GA {
	return GA{
		GenomeFactory: GenomeFactory,
		NPops:         1,
		PopSize:       1,
		Model: ModMutationOnly{
			NChosen:  1,
			Selector: SelElitism{},
			Strict:   true,
		},
	}
}
