package gago

import "strconv"

// floatToString converts a float into a string.
func floatToString(float float64) string {
	return strconv.FormatFloat(float, 'f', 6, 64)
}

// floatSliceToString converts a slice of floats into a string by using the
// floatToString method.
func floatSliceToString(slice []float64) string {
	var str string
	for _, float := range slice {
		str += floatToString(float) + ", "
	}
	return str
}

// Display an individual.
func (indi Individual) String() string {
	return "Fitness: " + floatToString(indi.Fitness) +
		" | DNA: [" + floatSliceToString(indi.Dna) + "]"
}
