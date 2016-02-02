package gago

import (
	"strconv"
	"strings"
)

// floatToString converts a float into a string.
func floatToString(float float64) string {
	return strconv.FormatFloat(float, 'f', 6, 64)
}

// floatSliceToString converts a slice of floats into a string by using the
// previously defined floatToString method.
func floatSliceToString(slice []float64) string {
	var str = make([]string, len(slice))
	for i, float := range slice {
		str[i] = floatToString(float)
	}
	return strings.Join(str, ", ")
}

// Display an individual.
func (indi Individual) String() string {
	return "Fitness: " + floatToString(indi.Fitness) +
		" | DNA: [" + floatSliceToString(indi.Dna) + "]"
}
