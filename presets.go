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
func GATSP(places []string, distance func([]string) float64) GA {
	return GA{
		NbPopulations: 4,
		NbIndividuals: 30,
		NbGenes:       len(places),
		Ff: StringFunction{
			Image: distance,
		},
		Initializer: ISUnique{
			Corpus: places,
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

// GAAlignment returns a configuration for solving string alignment problems.
// function. The output will be a genome of a certain length with genes
// belonging to a corpus of elements.
func GAAlignment(length int, corpus []string, distance func([]string) float64) GA {
	return GA{
		NbPopulations: 4,
		NbIndividuals: 30,
		NbGenes:       length,
		Ff: StringFunction{
			Image: distance,
		},
		Initializer: ISUniform{
			Corpus: corpus,
		},
		Selector: STournament{
			NbParticipants: 3,
		},
		Crossover: CPoint{
			NbPoints: 2,
		},
		Mutators: []Mutator{
			MutPermute{
				Rate: 0.5,
				Max:  3,
			},
			MutSplice{
				Rate: 0.5,
			},
			MutSUniform{
				Rate:   0.5,
				Corpus: corpus,
			},
		},
		Migrator: MigShuffle{},
	}
}
