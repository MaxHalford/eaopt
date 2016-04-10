package gago

import "math/rand"

// The Initializer is here to create the first generation of individuals in a
// population. It applies to an individual level and instantiates it's genome gene by
// gene.
type Initializer interface {
	apply(individual *Individual, generator *rand.Rand)
}

// IFUniform generates random floating points x such that lower < x < upper.
type IFUniform struct {
	Lower, Upper float64
}

// Apply the IFUniform initializer.
func (ifu IFUniform) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		// Decide if positive or negative
		var gene float64
		if generator.Float64() < 0.5 {
			gene = generator.Float64() * ifu.Lower
		} else {
			gene = generator.Float64() * ifu.Upper
		}
		indi.Genome[i] = gene
	}
}

// IFGaussian generates random floating point values sampled from a normal
// distribution.
type IFGaussian struct {
	Mean, Std float64
}

// Apply the IFGaussian initializer.
func (ifg IFGaussian) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = generator.NormFloat64()*ifg.Std + ifg.Mean
	}
}

// ISUniform generates random string slices based on a given corpus.
type ISUniform struct {
	Corpus []string
}

// Apply the ISUniform initializer.
func (isu ISUniform) apply(indi *Individual, generator *rand.Rand) {
	for i := range indi.Genome {
		indi.Genome[i] = isu.Corpus[generator.Intn(len(isu.Corpus))]
	}
}

// ISUnique generates random string slices based on a given corpus, each element
// from the corpus is only represented once in each slice. The method starts by
// shuffling, it then assigns the elements of the corpus in increasing index
// order to an individual. Usually the length of the individual's genome should
// match the length of the corpus.
type ISUnique struct {
	Corpus []string
}

// Apply the ISUnique initializer.
func (isu ISUnique) apply(indi *Individual, generator *rand.Rand) {
	var strings = shuffleStrings(isu.Corpus, generator)
	for i := range indi.Genome {
		indi.Genome[i] = strings[i]
	}
}
