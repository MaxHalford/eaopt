package gago

// Default configuration
var Default = Population{
	NbDemes:       1,
	NbIndividuals: 30,
	Boundary:      100.0,
	Selection:     tournament,
	Crossover:     crossover,
    CSize: 2,
	Mutate:        mutate,
	MRate:          0.1,
}
