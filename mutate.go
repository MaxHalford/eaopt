package gago

import (
	"math/rand"
)

// Type specific mutations for slices

// MutNormalFloat64 modifies a float64 gene if a coin toss is under a defined
// mutation rate. The new gene value is a random value sampled from a normal
// distribution centered on the gene's current value and with a standard
// deviation proportional to the current value. It does so for each gene.
func MutNormalFloat64(genome []float64, rate float64, rng *rand.Rand) {
	for i := range genome {
		// Flip a coin and decide to mutate or not
		if rng.Float64() < rate {
			genome[i] += rng.NormFloat64() * genome[i]
		}
	}
}

// MutUniformString picks a gene at random and replaces it with a random from a
// provided corpus. It repeats this n times.
func MutUniformString(genome []string, corpus []string, n int, rng *rand.Rand) {
	for i := 0; i < n; i++ {
		var (
			element = corpus[rng.Intn(len(corpus))]
			pos     = rng.Intn(len(genome))
		)
		genome[pos] = element
	}
}

// Generic mutations for slices

// MutPermute permutes two genes at random n times.
func MutPermute(genome Slice, n int, rng *rand.Rand) {
	// Nothing to permute
	if genome.Len() <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		// Choose two points on the genome
		var points = randomInts(2, 0, genome.Len(), rng)
		genome.Swap(points[0], points[1])
	}
}

// MutPermuteInt calls MutPermute on an int slice.
func MutPermuteInt(s []int, n int, rng *rand.Rand) {
	MutPermute(IntSlice(s), n, rng)
}

// MutPermuteFloat64 calls MutPermute on a float64 slice.
func MutPermuteFloat64(s []float64, n int, rng *rand.Rand) {
	MutPermute(Float64Slice(s), n, rng)
}

// MutPermuteString callsMutPermute on a string slice.
func MutPermuteString(s []string, n int, rng *rand.Rand) {
	MutPermute(StringSlice(s), n, rng)
}

// MutSplice splits a genome in 2 and glues the pieces back together in reverse
// order.
func MutSplice(genome Slice, rng *rand.Rand) {
	var (
		k    = rng.Intn(genome.Len()-1) + 1
		a, b = genome.Split(k)
	)
	genome.Replace(b.Append(a))
}

// MutSpliceInt calls MutSplice on an int slice.
func MutSpliceInt(s []int, rng *rand.Rand) {
	MutSplice(IntSlice(s), rng)
}

// MutSpliceFloat64 calls MutSplice on a float64 slice.
func MutSpliceFloat64(s []float64, rng *rand.Rand) {
	MutSplice(Float64Slice(s), rng)
}

// MutSpliceString calls MutSplice on a string slice.
func MutSpliceString(s []string, rng *rand.Rand) {
	MutSplice(StringSlice(s), rng)
}
