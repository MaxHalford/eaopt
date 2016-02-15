// Adapted from Programming Collective Intelligence, O'Reilly, Chapter 5.

// Not finished yet.
package main

import "github.com/MaxHalford/gago"

var (
	// Dorm room names and spaces
	dorms = map[string]int{
		"Zeus":     2,
		"Athena":   2,
		"Hercules": 2,
		"Bacchus":  2,
		"Pluto":    2,
	}
	// Student ordered preferences
	preferences = map[int][]string{
		0: []string{"Bacchus", "Hercules"},
		1: []string{"Zeus", "Pluto"},
		2: []string{"Athena", "Zeus"},
		3: []string{"Zeus", "Pluto"},
		4: []string{"Athena", "Bacchus"},
		5: []string{"Hercules", "Pluto"},
		6: []string{"Pluto", "Athena"},
		7: []string{"Bacchus", "Hercules"},
		8: []string{"Bacchus", "Hercules"},
		9: []string{"Hercules", "Athena"},
	}
)

// Happiness of all the people based on their ordered preferences
func happiness(assignment []string) float64 {
	var score = 0
	for i, dorm := range assignment {
		var position = len(preferences[i])
		// Go through the person's choices
		for j := len(preferences[i]) - 1; j >= 0; j-- {
			// Check if the person has been assigned to more preferable choice
			if dorm == preferences[i][j] {
				position = j
			}
		}
		// Increment the score by the position of the assignement
		score += position
	}
	return float64(score)
}

func main() {
	// Get the name of the dorms
	DormNames := make([]string, len(dorms))
	i := 0
	for dormName := range dorms {
		DormNames[i] = dormName
		i++
	}
	// Instantiate a population
	ga := gago.String
	// Use a custom initializer
	ga.Initializer = gago.UniformString{DormNames}
	// Use a custom mutator
	ga.Mutator = gago.Corpus{0.1, DormNames}
	// Wrap the evaluation function
	ga.Ff = gago.StringFunction{happiness}
	// Initialize with the number of individuals as the number of variables
	ga.Initialize(len(preferences))
	// Enhancement
	for i := 0; i < 1000; i++ {
		ga.Best.Display()
		ga.Enhance()
	}
}
