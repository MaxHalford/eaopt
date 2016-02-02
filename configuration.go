package gago

// Default configuration.
var Default = Population{
	NbDemes:       1,
	NbIndividuals: 30,
	Boundary:      100.0,
	Selection:     tournament,
	CrossMethod:   parenthood,
	CrossSize:     2,
	MutMethod:     normal,
	MutRate:       0.1,
	MutIntensity:  1,
}

// Medium configuration.
var Medium = Population{
	NbDemes:       4,
	NbIndividuals: 50,
	Boundary:      100.0,
	Selection:     tournament,
	CrossMethod:   parenthood,
	CrossSize:     2,
	MutMethod:     normal,
	MutRate:       0.2,
	MutIntensity:  1,
}
