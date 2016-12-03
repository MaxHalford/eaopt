package gago

import "math/rand"

// InitUnifFloat64 generates random float64s x such that lower < x < upper.
func InitUnifFloat64(n int, lower, upper float64, rng *rand.Rand) (vector []float64) {
	vector = make([]float64, n)
	for i := range vector {
		vector[i] = lower + rng.Float64()*(upper-lower)
	}
	return
}

// InitNormFloat64 generates random float64s sampled from a normal
// distribution.
func InitNormFloat64(n int, mean, std float64, rng *rand.Rand) (vector []float64) {
	vector = make([]float64, n)
	for i := range vector {
		vector[i] = rng.NormFloat64()*std + mean
	}
	return
}

// InitUnifString generates random strings based on a given corpus. The strings
// are not necessarily distinct.
func InitUnifString(n int, corpus []string, rng *rand.Rand) (strings []string) {
	strings = make([]string, n)
	for i := range strings {
		strings[i] = corpus[rng.Intn(len(corpus))]
	}
	return
}

// InitUniqueString generates random string slices based on a given corpus, each
// element from the corpus is only represented once in each slice. The method
// starts by shuffling, it then assigns the elements of the corpus in increasing
// index order to an individual.
func InitUniqueString(n int, corpus []string, rng *rand.Rand) (strings []string) {
	strings = make([]string, n)
	for i, v := range randomInts(n, 0, len(corpus), rng) {
		strings[i] = corpus[v]
	}
	return
}
