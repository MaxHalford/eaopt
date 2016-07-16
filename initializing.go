package gago

import "math/rand"

// The Initializer is here to create the first generation of individuals in a
// population. It applies to an individual level and instantiates it's genome gene by
// gene.
type Initializer interface {
	apply(indi *Individual, rng *rand.Rand)
}

// InitUniformF generates random floating points x such that lower < x < upper.
type InitUniformF struct {
	Lower, Upper float64
}

// Apply the InitUniformF initializer.
func (init InitUniformF) apply(indi *Individual, rng *rand.Rand) {
	for i := range indi.Genome {
		var gene float64
		// Decide if positive or negative
		if rng.Float64() < 0.5 {
			gene = rng.Float64() * init.Lower
		} else {
			gene = rng.Float64() * init.Upper
		}
		indi.Genome[i] = gene
	}
}

// InitGaussianF generates random floating point values sampled from a normal
// distribution.
type InitGaussianF struct {
	Mean, Std float64
}

// Apply the InitGaussianF initializer.
func (init InitGaussianF) apply(indi *Individual, rng *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = rng.NormFloat64()*init.Std + init.Mean
	}
}

// InitUniformS generates random string slices based on a given corpus.
type InitUniformS struct {
	Corpus []string
}

// Apply the InitUniformS initializer.
func (init InitUniformS) apply(indi *Individual, rng *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = init.Corpus[rng.Intn(len(init.Corpus))]
	}
}

// InitUniqueS generates random string slices based on a given corpus, each
// element from the corpus is only represented once in each slice. The method
// starts by shuffling, it then assigns the elements of the corpus in increasing
// index order to an individual. Usually the length of the individual's genome
// should match the length of the corpus.
type InitUniqueS struct {
	Corpus []string
}

// Apply the InitUniqueS initializer.
func (init InitUniqueS) apply(indi *Individual, rng *rand.Rand) {
	var strings = shuffleStrings(init.Corpus, rng)
	for i := range indi.Genome {
		indi.Genome[i] = strings[i]
	}
}
