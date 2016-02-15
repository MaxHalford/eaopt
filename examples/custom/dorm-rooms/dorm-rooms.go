// Adapted from Programming Collective Intelligence, O'Reilly, Chapter 5.
package main

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
	preferences = map[string][]string{
		"Toby":   []string{"Bacchus", "Hercules"},
		"Steve":  []string{"Zeus", "Pluto"},
		"Andrea": []string{"Athena", "Zeus"},
		"Sarah":  []string{"Zeus", "Pluto"},
		"Dave":   []string{"Athena", "Bacchus"},
		"Jeff":   []string{"Hercules", "Pluto"},
		"Fred":   []string{"Pluto", "Athena"},
		"Suzie":  []string{"Bacchus", "Hercules"},
		"Laura":  []string{"Bacchus", "Hercules"},
		"Neil":   []string{"Hercules", "Athena"},
	}
)
