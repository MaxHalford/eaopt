package eaopt

import (
	"math"
	"math/rand"
	"runtime"
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
func newIndividuals(n uint, parallel bool, newGenome func(rng *rand.Rand) Genome, rng *rand.Rand) Individuals {
	var indis = make(Individuals, n)
	if !parallel {
		for i := range indis {
			indis[i] = NewIndividual(newGenome(rng), rng)
		}
		return indis
	}
	
	var (
		nWorkers  = uint(runtime.GOMAXPROCS(-1))
		chunkSize = (n + nWorkers - 1) / nWorkers
		g         errgroup.Group
	)

	for a := uint(0); a < n; a += chunkSize {
		a := a // https://golang.org/doc/faq#closures_and_goroutines
		var b = minUint(a+chunkSize, n)
		seed := rng.Int63()
		g.Go(func() error {
			indRNG := rand.New(rand.NewSource(seed))
			for i := a; i < b; i++ {
				indis[i] = NewIndividual(newGenome(indRNG), indRNG)
			}
			return nil
		})
	}
	g.Wait()
	return indis
}

// Evaluate each Individual in a slice.
func (indis Individuals) Evaluate(parallel bool) error {
	if !parallel {
		var err error
		for i := range indis {
			err = indis[i].Evaluate()
			if err != nil {
				return err
			}
		}
		return nil
	}

	var (
		nWorkers  = runtime.GOMAXPROCS(-1)
		n         = len(indis)
		chunkSize = (n + nWorkers - 1) / nWorkers
		g         errgroup.Group
	)

	for a := 0; a < n; a += chunkSize {
		a := a // https://golang.org/doc/faq#closures_and_goroutines
		var b = minInt(a+chunkSize, n)
		g.Go(func() error {
			return indis[a:b].Evaluate(false)
		})
	}

	return g.Wait()
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

// FitMax returns the worst fitness of a slice of individuals.
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
