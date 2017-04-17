package gago

import (
	"math"
	"math/rand"
	"time"
)

// Divide each element in a float64 slice by a given value.
func divide(floats []float64, value float64) []float64 {
	var divided = make([]float64, len(floats))
	for i, v := range floats {
		divided[i] = v / value
	}
	return divided
}

// Compute the cumulative sum of a float64 slice.
func cumsum(floats []float64) []float64 {
	var summed = make([]float64, len(floats))
	copy(summed, floats)
	for i := 1; i < len(summed); i++ {
		summed[i] += summed[i-1]
	}
	return summed
}

// Generate random weights that sum up to 1.
func randomWeights(size int) []float64 {
	var (
		weights = make([]float64, size)
		total   float64
	)
	for i := range weights {
		weights[i] = rand.Float64()
		total += weights[i]
	}
	var normalized = divide(weights, total)
	return normalized
}

// Find the minimum between two ints.
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Compute the sum of a float64 slice.
func sumFloat64s(floats []float64) (sum float64) {
	for _, v := range floats {
		sum += v
	}
	return
}

// Compute the minimum value of a float64 slice.
func minFloat64s(floats []float64) (min float64) {
	min = math.Inf(1)
	for _, f := range floats {
		if f < min {
			min = f
		}
	}
	return
}

// Compute the maximum value of a float64 slice.
func maxFloat64s(floats []float64) (max float64) {
	max = math.Inf(-1)
	for _, f := range floats {
		if f > max {
			max = f
		}
	}
	return
}

// Compute the mean of a float64 slice.
func meanFloat64s(floats []float64) float64 {
	return sumFloat64s(floats) / float64(len(floats))
}

// Compute the variance of a float64 slice.
func varianceFloat64s(floats []float64) float64 {
	var squares = make([]float64, len(floats))
	for i, f := range floats {
		squares[i] = math.Pow(f, 2)
	}
	return meanFloat64s(squares) - math.Pow(meanFloat64s(floats), 2)
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

// makeRandomNumberGenerator returns a new random number generator with a random
// seed.
func makeRandomNumberGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type set map[interface{}]bool

// union merges two slices and ignores duplicates.
func union(x, y set) set {
	var (
		u         = make(set)
		blackList = make(map[interface{}]bool)
	)
	for i := range x {
		u[i] = true
		blackList[i] = true
	}
	for i := range y {
		if !blackList[i] {
			u[i] = true
			blackList[i] = true
		}
	}
	return u
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int, rng *rand.Rand) string {
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
