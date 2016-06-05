package gago

import "math/rand"

// The Initializer is here to create the first generation of individuals in a
// population. It applies to an individual level and instantiates it's genome gene by
// gene.
type Initializer interface {
	apply(individual *Individual, generator *rand.Rand)
}

// InitUniformF generates random floating points x such that lower < x < upper.
type InitUniformF struct {
	Lower, Upper float64
}

// Apply the InitUniformF initializer.
func (iuf InitUniformF) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Decide if positive or negative
		var gene float64
		if generator.Float64() < 0.5 {
			gene = generator.Float64() * iuf.Lower
		} else {
			gene = generator.Float64() * iuf.Upper
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
func (igf InitGaussianF) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = generator.NormFloat64()*igf.Std + igf.Mean
	}
}

// InitUniformS generates random string slices based on a given corpus.
type InitUniformS struct {
	Corpus []string
}

// Apply the InitUniformS initializer.
func (ius InitUniformS) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = ius.Corpus[generator.Intn(len(ius.Corpus))]
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
func (ius InitUniqueS) apply(indi *Individual, generator *rand.Rand) {
	var strings = shuffleStrings(ius.Corpus, generator)
	for i := range indi.Genome {
		indi.Genome[i] = strings[i]
	}
}
