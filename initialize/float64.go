package initialize

import "math/rand"

// UniformFloat64 generates random floating points x such that lower < x < upper.
func UniformFloat64(n int, rng *rand.Rand, lower, upper float64) []float64 {
	var vector = make([]float64, n)
	for i := range vector {
		var v float64
		if rng.Float64() < 0.5 {
			v = rng.Float64() * lower
		} else {
			v = rng.Float64() * upper
		}
		vector[i] = v
	}
	return vector
}

// NormalFloat64 generates random floating point values sampled from a normal
// distribution.
func NormalFloat64(n int, rng *rand.Rand, mean, std float64) []float64 {
	var vector = make([]float64, n)
	for i := range vector {
		vector[i] = rng.NormFloat64()*std + mean
	}
	return vector
}
