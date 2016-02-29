package gago

import "math/rand"

// Find where an element is in a slice.
func findPosition(array []interface{}, element interface{}) int {
	var index = 0
	for array[index] != element {
		if index >= len(array)-1 {
			return -1
		}
		index++
	}
	return index
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
