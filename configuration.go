package gago

import "strings"

// Float problem configuration.
var Float = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer:   UniformFloat{-10, 10},
	Selector:      Tournament{3},
	Breeder:       Parenthood{},
	Mutator:       Normal{0.1, 1},
	Migrator:      Shuffle{},
}

var latin = strings.Split("abcdefghijklmnopqrstuvwxyz", "")

// String problem configuration.
var String = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer:   UniformString{latin},
	Selector:      Tournament{3},
	Breeder:       Crossover{},
	Mutator:       Corpus{0.1, latin},
	Migrator:      Shuffle{},
}
