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

var latinAlphabet = strings.Split("abcdefghijklmnopqrstuvwxyz ", "")

// String problem configuration.
var String = Population{
	NbDemes:       2,
	NbIndividuals: 30,
	Initializer:   UniformString{latinAlphabet},
	Selector:      Tournament{3},
	Breeder:       Crossover{},
	Mutator:       Corpus{0.1, latinAlphabet},
	Migrator:      Shuffle{},
}
