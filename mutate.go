package gago

import "math/rand"

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
func MutPermute(genome []interface{}, n int, rng *rand.Rand) {
	// Nothing to permute
	if len(genome) <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		// Choose two points on the genome
		var (
			points = randomInts(2, 0, len(genome), rng)
			i      = points[0]
			j      = points[1]
		)
		// Permute the genes
		genome[i], genome[j] = genome[j], genome[i]
	}
}

// MutPermuteFloat64 is a convenience function for calling MutPermute on a
// float64 slice.
func MutPermuteFloat64(values []float64, n int, rng *rand.Rand) {
	var genome = uncastFloat64s(values)
	MutPermute(genome, n, rng)
	copy(values, castFloat64s(genome))
}

// MutPermuteInt is a convenience function for calling MutPermute on an int
// slice.
func MutPermuteInt(values []int, n int, rng *rand.Rand) {
	var genome = uncastInts(values)
	MutPermute(genome, n, rng)
	copy(values, castInts(genome))
}

// MutPermuteString is a convenience function for calling MutPermute on a string
// slice.
func MutPermuteString(values []string, n int, rng *rand.Rand) {
	var genome = uncastStrings(values)
	MutPermute(genome, n, rng)
	copy(values, castStrings(genome))
}

// MutSplice splits a genome in 2 and glues the pieces back together in a
// different order.
func MutSplice(genome []interface{}, rng *rand.Rand) {
	var split = rng.Intn(len(genome)-1) + 1
	copy(genome, append(genome[split:], genome[:split]...))
}

// MutSpliceFloat64 is a convenience function for calling MutSplice on a float64
// slice.
func MutSpliceFloat64(values []float64, rng *rand.Rand) {
	var genome = uncastFloat64s(values)
	MutSplice(genome, rng)
	copy(values, castFloat64s(genome))
}

// MutSpliceInt is a convenience function for calling MutSplice on an int
// slice.
func MutSpliceInt(values []int, rng *rand.Rand) {
	var genome = uncastInts(values)
	MutSplice(genome, rng)
	copy(values, castInts(genome))
}

// MutSpliceString is a convenience function for calling MutSplice on a string
// slice.
func MutSpliceString(values []string, rng *rand.Rand) {
	var genome = uncastStrings(values)
	MutSplice(genome, rng)
	copy(values, castStrings(genome))
}
