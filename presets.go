package gago

// Float problem configuration.
var Float = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer: IFUniform{
		Lower: -10,
		Upper: 10,
	},
	Selector: STournament{
		NbParticipants: 3,
	},
	Crossover: CFUniform{},
	Mutators: []Mutator{
		MutFNormal{
			Rate: 0.1,
			Std:  1,
		},
	},
	Migrator: MigShuffle{},
}
