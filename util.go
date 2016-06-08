package gago

import "math/rand"

// Find where an element is in a slice.
func getIndex(element interface{}, array []interface{}) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	// Element not in array
	return -1
}

// Generate random weights that sum up to 1?
func generateWeights(size int) []float64 {
	var weights = make([]float64, size)
	// Sum of the weights
	var total float64
	// Assign a weight to each individual
	for i := range weights {
		weights[i] = rand.Float64()
		total += weights[i]
	}
	// Normalize the weights
	for i := range weights {
		weights[i] /= total
	}
	return weights
}

// Shuffle a slice of strings.
func shuffleStrings(strings []string, generator *rand.Rand) []string {
	var shuffled = make([]string, len(strings))
	for i, j := range generator.Perm(len(strings)) {
		shuffled[j] = strings[i]
	}
	return shuffled
}

// Find the strict minimum between two integers
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
