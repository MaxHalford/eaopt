package gago

import (
	"fmt"
	"math/rand"
)

// Find where an element is in a slice.
func getIndex(element interface{}, array []interface{}) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	fmt.Println(element, array)
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
	for i, v := range generator.Perm(len(strings)) {
		shuffled[v] = strings[i]
	}
	return shuffled
}

// Shuffle a slice of individuals.
func shuffleIndividuals(indis Individuals, generator *rand.Rand) Individuals {
	var shuffled = make(Individuals, len(indis))
	for i, v := range generator.Perm(len(indis)) {
		shuffled[v] = indis[i]
	}
	return shuffled
}

// Sample n unique individuals from a slice of individuals
func sampleIndividuals(n int, indis Individuals, generator *rand.Rand) Individuals {
	var shuffled = shuffleIndividuals(indis, generator)
	var sample = shuffled[:n]
	return sample
}
