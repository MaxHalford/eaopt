package gago

// GAFloat returns a configuration for minimizing continuous mathematical
// functions with a given number of variables.
func GAFloat(n int, function func([]float64) float64) GA {
	return GA{
		NbPopulations: 4,
		NbIndividuals: 30,
		NbGenes:       n,
		Ff: FloatFunction{
			Image: function,
		},
		Initializer: IFUniform{
			Lower: -1,
			Upper: 1,
		},
		Selector: STournament{
			NbParticipants: 3,
		},
		Crossover: CFUniform{},
		Mutators: []Mutator{
			MutFNormal{
				Rate: 0.5,
				Std:  3,
			},
		},
		Migrator: MigShuffle{},
	}
}

// GATSP returns a configuration for solving Travelling Salesman Problems given
// a corpus of positions that are associated to coordinates in the fitness
// function.
func GATSP(corpus []string, distance func([]string) float64) GA {
	return GA{
		NbPopulations: 4,
		NbIndividuals: 30,
		NbGenes:       len(corpus),
		Ff: StringFunction{
			Image: distance,
		},
		Initializer: ISUnique{
			Corpus: corpus,
		},
		Selector:  SElitism{},
		Crossover: CPMX{},
		Mutators: []Mutator{
			MutPermute{
				Rate: 0.5,
				Max:  3,
			},
			MutSplice{
				Rate: 0.5,
			},
		},
		Migrator: MigShuffle{},
	}
}
