package gago

// Float problem configuration.
var Float = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer:   FloatUniform{-10, 10},
	Selector:      Tournament{3},
	Crossover:     FloatParenthood{},
	Mutators:      []Mutator{FloatNormal{0.1, 1}},
	Migrator:      Shuffle{},
}
