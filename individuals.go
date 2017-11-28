package gago

import (
	"math"
	"math/rand"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
)

// Individuals is a convenience type, methods that belong to an Individual can
// be called declaratively.
type Individuals []Individual

// String representation of a slice of Individuals.
func (indis Individuals) String() string {
	var str string
	for _, indi := range indis {
		str += indi.String() + "\n"
	}
	return strings.TrimSuffix(str, "\n")
}

// Clone returns the same exact same slice of individuals but with different
// pointers and ID fields.
func (indis Individuals) Clone(rng *rand.Rand) Individuals {
	var clones = make(Individuals, len(indis))
	for i, indi := range indis {
		clones[i] = indi.Clone(rng)
	}
	return clones
}

// Generate a slice of n new individuals.
func newIndividuals(n int, newGenome NewGenome, rng *rand.Rand) Individuals {
	var indis = make(Individuals, n)
	for i := range indis {
		indis[i] = NewIndividual(newGenome(rng), rng)
	}
	return indis
}

// Apply a function to a slice of Individuals.
func (indis Individuals) Apply(f func(indi *Individual) error, parallel bool) error {
	if parallel {
		var g errgroup.Group
		for i := range indis {
			i := i // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				return f(&indis[i])
			})
		}
		return g.Wait()
	}
	var err error
	for i := range indis {
		err = f(&indis[i])
		if err != nil {
			return err
		}
	}
	return err
}

// Evaluate each Individual. If parallel is true then each Individual will be
// evaluated in parallel thanks to the golang.org/x/sync/errgroup package. If
// not then a simple sequential loop will be used. Evaluating in parallel is
// only recommended for cases where evaluating an Individual takes a "long"
// time. Indeed there won't necessarily be a speed-up when evaluating in
// parallel. In fact performance can be degraded if evaluating an Individual is
// too cheap.
func (indis Individuals) Evaluate(parallel bool) {
	indis.Apply(
		func(indi *Individual) error { indi.Evaluate(); return nil },
		parallel,
	)
}

// Mutate each individual.
func (indis Individuals) Mutate(mutRate float64, rng *rand.Rand) {
	for i := range indis {
		if rng.Float64() < mutRate {
			indis[i].Mutate(rng)
		}
	}
}

// SortByFitness ascendingly sorts individuals by fitness.
func (indis Individuals) SortByFitness() {
	var less = func(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
	sort.Slice(indis, less)
}

// IsSortedByFitness checks if individuals are ascendingly sorted by fitness.
func (indis Individuals) IsSortedByFitness() bool {
	var less = func(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
	return sort.SliceIsSorted(indis, less)
}

// SortByDistanceToMedoid sorts Individuals according to their distance to the
// medoid. The medoid is the Individual that has the lowest average distance to
// the rest of the Individuals.
func (indis Individuals) SortByDistanceToMedoid(dm DistanceMemoizer) {
	var (
		avgDists = calcAvgDistances(indis, dm)
		less     = func(i, j int) bool {
			return avgDists[indis[i].ID] < avgDists[indis[j].ID]
		}
	)
	sort.Slice(indis, less)
}

// Extract the fitness of a slice of individuals into a float64 slice.
func (indis Individuals) getFitnesses() []float64 {
	var fitnesses = make([]float64, len(indis))
	for i, indi := range indis {
		fitnesses[i] = indi.Fitness
	}
	return fitnesses
}

// FitMin returns the best fitness of a slice of individuals.
func (indis Individuals) FitMin() float64 {
	if indis.IsSortedByFitness() {
		return indis[0].Fitness
	}
	return minFloat64s(indis.getFitnesses())
}

// FitMax returns the best fitness of a slice of individuals.
func (indis Individuals) FitMax() float64 {
	if indis.IsSortedByFitness() {
		return indis[len(indis)-1].Fitness
	}
	return maxFloat64s(indis.getFitnesses())
}

// FitAvg returns the average fitness of a slice of individuals.
func (indis Individuals) FitAvg() float64 {
	return meanFloat64s(indis.getFitnesses())
}

// FitStd returns the standard deviation of the fitness of a slice of
// individuals.
func (indis Individuals) FitStd() float64 {
	return math.Sqrt(varianceFloat64s(indis.getFitnesses()))
}
