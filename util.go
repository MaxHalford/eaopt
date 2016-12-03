package gago

import (
	"math"
	"math/rand"
	"time"
)

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

// Generate random weights that sum up to 1.
func randomWeights(size int) []float64 {
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

// Find the strict minimum between two ints.
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Compute the mean of a slice of a float64 slice.
func mean(slice []float64) float64 {
	var sum float64
	for _, v := range slice {
		sum += v
	}
	return sum / float64(len(slice))
}

// Compute the variance of a float64 slice.
func variance(slice []float64) float64 {
	// Compute the squares
	var squares = make([]float64, len(slice))
	for i, v := range slice {
		squares[i] = math.Pow(v, 2)
	}
	return mean(squares) - math.Pow(mean(slice), 2)
}

// Sample k unique integers in range [min, max) using reservoir sampling,
// specifically Vitter's Algorithm R.
func randomInts(k, min, max int, rng *rand.Rand) (ints []int) {
	ints = make([]int, k)
	for i := 0; i < k; i++ {
		ints[i] = i + min
	}
	for i := k; i < max-min; i++ {
		var j = rng.Intn(i + 1)
		if j < k {
			ints[j] = i + min
		}
	}
	return
}

func makeRandomNumberGenerator() (rng *rand.Rand) {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	return
}
