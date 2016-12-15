package gago

import (
	"math/rand"
	"sort"
)

// Type specific mutations for slices

// CrossUniformFloat64 crossover combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the value of a probability p. Each offspring receives a proportion of both of
// it's parents genomes. The new values are located in the hyper-rectangle
// defined between both parent's position in Cartesian space.
func CrossUniformFloat64(p1 []float64, p2 []float64, rng *rand.Rand) (o1 []float64, o2 []float64) {
	var gSize = len(p1)
	o1 = make([]float64, gSize)
	o2 = make([]float64, gSize)
	// For every gene pick a random number between 0 and 1
	for i := 0; i < gSize; i++ {
		var p = rng.Float64()
		o1[i] = p*p1[i] + (1-p)*p2[i]
		o2[i] = (1-p)*p1[i] + p*p2[i]
	}
	return o1, o2
}

// Generic mutations for slices

// Contains the deterministic part of the GNX method for testing purposes.
func gnx(p1, p2 []interface{}, indexes []int) ([]interface{}, []interface{}) {
	var (
		n  = len(p1)
		o1 = make([]interface{}, n)
		o2 = make([]interface{}, n)
		s  = true // Switch
	)
	// Add the first and last indexes
	indexes = append([]int{0}, indexes...)
	indexes = append(indexes, n)
	for i := 0; i < len(indexes)-1; i++ {
		if s {
			copy(o1[indexes[i]:indexes[i+1]], p1[indexes[i]:indexes[i+1]])
			copy(o2[indexes[i]:indexes[i+1]], p2[indexes[i]:indexes[i+1]])
		} else {
			copy(o1[indexes[i]:indexes[i+1]], p2[indexes[i]:indexes[i+1]])
			copy(o2[indexes[i]:indexes[i+1]], p1[indexes[i]:indexes[i+1]])
		}
		s = !s // Alternate for the new copying
	}
	return o1, o2
}

// CrossGNX (Generalized N-point Crossover). An identical point is chosen on
// each parent's genome and the mirroring segments are switched. n determines
// the number of crossovers (aka mirroring segments) to perform. n has to be
// equal or lower than the number of genes in each parent.
func CrossGNX(p1 []interface{}, p2 []interface{}, n int, rng *rand.Rand) (o1 []interface{}, o2 []interface{}) {
	var indexes = randomInts(n, 1, len(p1), rng)
	sort.Ints(indexes)
	return gnx(p1, p2, indexes)
}

// CrossGNXFloat64 is a convenience function for calling CrossGNX on a
// float64 slice.
func CrossGNXFloat64(v1 []float64, v2 []float64, n int, rng *rand.Rand) ([]float64, []float64) {
	var (
		p1, p2 = uncastFloat64s(v1), uncastFloat64s(v2)
		o1, o2 = CrossGNX(p1, p2, n, rng)
	)
	return castFloat64s(o1), castFloat64s(o2)
}

// CrossGNXInt is a convenience function for calling CrossGNX on an int
// slice.
func CrossGNXInt(v1 []int, v2 []int, n int, rng *rand.Rand) ([]int, []int) {
	var (
		p1, p2 = uncastInts(v1), uncastInts(v2)
		o1, o2 = CrossGNX(p1, p2, n, rng)
	)
	return castInts(o1), castInts(o2)
}

// CrossGNXString is a convenience function for calling CrossGNX on a
// float64 slice.
func CrossGNXString(v1 []string, v2 []string, n int, rng *rand.Rand) ([]string, []string) {
	var (
		p1, p2 = uncastStrings(v1), uncastStrings(v2)
		o1, o2 = CrossGNX(p1, p2, n, rng)
	)
	return castStrings(o1), castStrings(o2)
}

// Contains the deterministic part of the PMX method for testing purposes.
func pmx(p1, p2 []interface{}, a, b int) ([]interface{}, []interface{}) {
	var (
		n  = len(p1)
		o1 = make([]interface{}, n)
		o2 = make([]interface{}, n)
	)
	// Copy part of the first parent's genome onto the first offspring
	copy(o1[a:b], p1[a:b])
	copy(o2[a:b], p2[a:b])
	for i := a; i < b; i++ {
		// Find the element in the second parent that has not been copied in the first offspring
		if !elementInSlice(p2[i], o1[a:b]) {
			var j = i
			for o1[j] != nil {
				j = getIndex(o1[j], p2)
			}
			o1[j] = p2[i]
		}
		// Find the element in the first parent that has not been copied in the second offspring
		if !elementInSlice(p1[i], o2[a:b]) {
			var j = i
			for o2[j] != nil {
				j = getIndex(o2[j], p1)
			}
			o2[j] = p1[i]
		}
	}
	// Fill in the offspring's missing values with the opposite parent's values
	for i := 0; i < n; i++ {
		if o1[i] == nil {
			o1[i] = p2[i]
		}
		if o2[i] == nil {
			o2[i] = p1[i]
		}
	}
	return o1, o2
}

// CrossPMX (Partially Mapped Crossover). The offsprings are generated by
// copying one of the parents and then copying the other parent's values up to a
// randomly chosen crossover point. Each gene that is replaced is permuted with
// the gene that is copied in the first parent's genome. Two offsprings are
// generated in such a way (because there are two parents). The PMX method
// preserves gene uniqueness.
func CrossPMX(p1 []interface{}, p2 []interface{}, rng *rand.Rand) (o1 []interface{}, o2 []interface{}) {
	var indexes = randomInts(2, 0, len(p1), rng)
	sort.Ints(indexes)
	o1, o2 = pmx(p1, p2, indexes[0], indexes[1])
	return o1, o2
}

// CrossPMXFloat64 is a convenience function for calling CrossPMX on a
// float64 slice.
func CrossPMXFloat64(v1 []float64, v2 []float64, n int, rng *rand.Rand) ([]float64, []float64) {
	var (
		p1, p2 = uncastFloat64s(v1), uncastFloat64s(v2)
		o1, o2 = CrossPMX(p1, p2, rng)
	)
	return castFloat64s(o1), castFloat64s(o2)
}

// CrossPMXInt is a convenience function for calling CrossPMX on an int
// slice.
func CrossPMXInt(v1 []int, v2 []int, n int, rng *rand.Rand) ([]int, []int) {
	var (
		p1, p2 = uncastInts(v1), uncastInts(v2)
		o1, o2 = CrossPMX(p1, p2, rng)
	)
	return castInts(o1), castInts(o2)
}

// CrossPMXString is a convenience function for calling CrossPMX on a
// float64 slice.
func CrossPMXString(v1 []string, v2 []string, n int, rng *rand.Rand) ([]string, []string) {
	var (
		p1, p2 = uncastStrings(v1), uncastStrings(v2)
		o1, o2 = CrossPMX(p1, p2, rng)
	)
	return castStrings(o1), castStrings(o2)
}

// Contains the deterministic part of the OX method for testing purposes.
func ox(p1, p2 []interface{}, a, b int) ([]interface{}, []interface{}) {
	var (
		n  = len(p1)
		o1 = make([]interface{}, n)
		o2 = make([]interface{}, n)
	)
	// Copy part of the first parent's genome onto the first offspring
	copy(o1[a:b], p1[a:b])
	copy(o2[a:b], p2[a:b])
	// Keep two indicators to know where to fill the offsprings
	var j1, j2 = b, b
	for i := b; i < b+n; i++ {
		var k = i % n
		if !elementInSlice(p2[k], o1[a:b]) {
			o1[j1%n] = p2[k]
			j1++
		}
		if !elementInSlice(p1[k], o2[a:b]) {
			o2[j2%n] = p1[k]
			j2++
		}
	}
	return o1, o2
}

// CrossOX (Ordered Crossover). Part of the first parent's genome is copied onto
// the first offspring's genome. Then the second parent's genome is iterated
// over, starting on the right of the part that was copied. Each gene of the
// second parent's genome is copied onto the next blank gene of the first
// offspring's genome if it wasn't already copied from the first parent. The OX
// method preserves gene uniqueness.
func CrossOX(p1 []interface{}, p2 []interface{}, rng *rand.Rand) (o1 []interface{}, o2 []interface{}) {
	var indexes = randomInts(2, 0, len(p1), rng)
	sort.Ints(indexes)
	o1, o2 = ox(p1, p2, indexes[0], indexes[1])
	return o1, o2
}

// CrossOXFloat64 is a convenience function for calling CrossOX on a float64
// slice.
func CrossOXFloat64(v1 []float64, v2 []float64, rng *rand.Rand) ([]float64, []float64) {
	var (
		p1, p2 = uncastFloat64s(v1), uncastFloat64s(v2)
		o1, o2 = CrossOX(p1, p2, rng)
	)
	return castFloat64s(o1), castFloat64s(o2)
}

// CrossOXInt is a convenience function for calling CrossOX on an int slice.
func CrossOXInt(v1 []int, v2 []int, rng *rand.Rand) ([]int, []int) {
	var (
		p1, p2 = uncastInts(v1), uncastInts(v2)
		o1, o2 = CrossOX(p1, p2, rng)
	)
	return castInts(o1), castInts(o2)
}

// CrossOXString is a convenience function for calling CrossOX on a float64
// slice.
func CrossOXString(v1 []string, v2 []string, rng *rand.Rand) ([]string, []string) {
	var (
		p1, p2 = uncastStrings(v1), uncastStrings(v2)
		o1, o2 = CrossOX(p1, p2, rng)
	)
	return castStrings(o1), castStrings(o2)
}

// getCycles determines the cycles that exist between two slices. A cycle is a
// list of indexes indicating mirroring values between each slice.
func getCycles(x, y []interface{}) (cycles [][]int) {
	var (
		xLookup = makeLookup(x)      // Matches values to indexes for quick lookup
		visited = make(map[int]bool) // Indicates if an index is already in a cycle or not
	)
	for i := 0; i < len(x); i++ {
		if !visited[i] {
			visited[i] = true
			var (
				cycle = []int{i}
				j     = xLookup[y[i]]
			)
			// Continue building the cycle until it closes in on itself
			for j != cycle[0] {
				cycle = append(cycle, j)
				visited[j] = true
				j = xLookup[y[j]]
			}
			cycles = append(cycles, cycle)
		}
	}
	return
}

// CrossCX (Cycle Crossover). Cycles between the parents are indentified, they
// are then copied alternatively onto the offsprings. The CX method is
// deterministic and preserves gene uniqueness.
func CrossCX(p1, p2 []interface{}) ([]interface{}, []interface{}) {
	var (
		n      = len(p1)
		o1     = make([]interface{}, n)
		o2     = make([]interface{}, n)
		cycles = getCycles(p1, p2)
		s      = true // Switch
	)
	for i := 0; i < len(cycles); i++ {
		for _, j := range cycles[i] {
			if s {
				o1[j], o2[j] = p1[j], p2[j]
			} else {
				o2[j], o1[j] = p1[j], p2[j]
			}
		}
		s = !s
	}
	return o1, o2
}

// CrossCXFloat64 is a convenience function for calling CrossCX on a float64
// slice.
func CrossCXFloat64(v1 []float64, v2 []float64) ([]float64, []float64) {
	var (
		p1, p2 = uncastFloat64s(v1), uncastFloat64s(v2)
		o1, o2 = CrossCX(p1, p2)
	)
	return castFloat64s(o1), castFloat64s(o2)
}

// CrossCXInt is a convenience function for calling CrossCX on an int slice.
func CrossCXInt(v1 []int, v2 []int) ([]int, []int) {
	var (
		p1, p2 = uncastInts(v1), uncastInts(v2)
		o1, o2 = CrossCX(p1, p2)
	)
	return castInts(o1), castInts(o2)
}

// CrossCXString is a convenience function for calling CrossCX on a float64
// slice.
func CrossCXString(v1 []string, v2 []string) ([]string, []string) {
	var (
		p1, p2 = uncastStrings(v1), uncastStrings(v2)
		o1, o2 = CrossCX(p1, p2)
	)
	return castStrings(o1), castStrings(o2)
}
