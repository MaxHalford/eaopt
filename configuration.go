package gago

// Default configuration.
var Default = Population{
	NbDemes:       1,
	NbIndividuals: 30,
	Boundary:      100.0,
	SelMethod:     Tournament,
	CrossMethod:   Parenthood,
	CrossSize:     2,
	MutMethod:     Normal,
	MutRate:       0.1,
	MutIntensity:  1,
	MigMethod:     Shuffle,
}

// Medium configuration.
var Medium = Population{
	NbDemes:       4,
	NbIndividuals: 50,
	Boundary:      100.0,
	SelMethod:     Tournament,
	CrossMethod:   Parenthood,
	CrossSize:     2,
	MutMethod:     Normal,
	MutRate:       0.2,
	MutIntensity:  1,
	MigMethod:     Shuffle,
}
