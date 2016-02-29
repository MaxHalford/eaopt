package gago

// Float problem configuration.
var Float = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer:   UniformFloat{-10, 10},
	Selector:      Tournament{3},
	Crossover:     Parenthood{},
	Mutator:       Normal{0.1, 1},
	Migrator:      Shuffle{},
}
