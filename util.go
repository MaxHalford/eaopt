package gago

import (
	"errors"
	"math"
	"math/rand"
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

// Shuffle a slice of strings.
func shuffleStrings(strings []string, rng *rand.Rand) []string {
	var shuffled = make([]string, len(strings))
	for i, j := range rng.Perm(len(strings)) {
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
// specifically Algorithm R. It can be proven by induction that each integer
// has probability of 1/(max-min) to be selected.
func randomInts(k, min, max int, rng *rand.Rand) ([]int, error) {
	if max-min < k {
		return []int{}, errors.New("k has to be superior or equak to max - min")
	}
	var ints = make([]int, k)
	for i := min; i < min+k; i++ {
		ints[i] = i
	}
	for i := min + k; i < max; i++ {
		var j = rng.Intn(i)
		if j < k {
			ints[j] = i
		}
	}
	return ints, nil
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Generate a random string of size n.
func randomString(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rng.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
