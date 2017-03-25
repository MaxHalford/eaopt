package gago

import "errors"

// A Slice is a genome with a list-like structure.
type Slice interface {
	At(i int) interface{}
	Set(i int, v interface{})
	Len() int
	Swap(i, j int)
	Slice(a, b int) Slice
	Split(k int) (Slice, Slice)
	Append(Slice) Slice
	Replace(Slice)
	Clone() Slice
}

// Search for the first index of an element in a Slice.
func search(v interface{}, s Slice) (int, error) {
	for i := 0; i < s.Len(); i++ {
		if s.At(i) == v {
			return i, nil
		}
	}
	// Element not in slice
	return 0, errors.New("Value not contained in slice")
}

// Make a lookup table from a slice, mapping values to indexes.
func makeIndexLookup(s Slice) map[interface{}]int {
	var lookup = make(map[interface{}]int)
	for i := 0; i < s.Len(); i++ {
		lookup[s.At(i)] = i
	}
	return lookup
}

// getCycles determines the cycles that exist between two slices. A cycle is a
// list of indexes indicating mirroring values between each slice.
func getCycles(s1, s2 Slice) (cycles [][]int) {
	var (
		s1Lookup = makeIndexLookup(s1) // Matches values to indexes for quick lookup
		visited  = make(map[int]bool)  // Indicates if an index is already in a cycle or not
	)
	for i := 0; i < s1.Len(); i++ {
		if !visited[i] {
			visited[i] = true
			var (
				cycle = []int{i}
				j     = s1Lookup[s2.At(i)]
			)
			// Continue building the cycle until it closes in on itself
			for j != cycle[0] {
				cycle = append(cycle, j)
				visited[j] = true
				j = s1Lookup[s2.At(j)]
			}
			cycles = append(cycles, cycle)
		}
	}
	return
}

// getNeighbours converts a slice into an adjacency map mapping values to left
// and right neighbours. The values of the map are sets.
func getNeighbours(s Slice) map[interface{}]set {
	var (
		neighbours = make(map[interface{}]set)
		n          = s.Len()
	)
	neighbours[s.At(0)] = set{s.At(n - 1): true, s.At(1): true}
	for i := 1; i < n-1; i++ {
		neighbours[s.At(i)] = set{s.At(i - 1): true, s.At(i + 1): true}
	}
	neighbours[s.At(n-1)] = set{s.At(n - 2): true, s.At(0): true}
	return neighbours
}

// IntSlice attaches the methods of Slice to []float64
type IntSlice []int

// At method from Slice
func (s IntSlice) At(i int) interface{} {
	return s[i]
}

// Set method from Slice
func (s IntSlice) Set(i int, v interface{}) {
	s[i] = v.(int)
}

// Len method from Slice
func (s IntSlice) Len() int {
	return len(s)
}

// Swap method from Slice
func (s IntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Slice method from Slice
func (s IntSlice) Slice(a, b int) Slice {
	return s[a:b]
}

// Split method from Slice
func (s IntSlice) Split(k int) (Slice, Slice) {
	return s[:k], s[k:]
}

// Append method from Slice
func (s IntSlice) Append(t Slice) Slice {
	return append(s, t.(IntSlice)...)
}

// Replace method from Slice
func (s IntSlice) Replace(t Slice) {
	copy(s, t.(IntSlice))
}

// Clone method from Slice
func (s IntSlice) Clone() Slice {
	var clone = make(IntSlice, len(s))
	copy(clone, s)
	return clone
}

// Float64Slice attaches the methods of Slice to []float64
type Float64Slice []float64

// At method from Slice
func (s Float64Slice) At(i int) interface{} {
	return s[i]
}

// Set method from Slice
func (s Float64Slice) Set(i int, v interface{}) {
	s[i] = v.(float64)
}

// Len method from Slice
func (s Float64Slice) Len() int {
	return len(s)
}

// Swap method from Slice
func (s Float64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Slice method from Slice
func (s Float64Slice) Slice(a, b int) Slice {
	return s[a:b]
}

// Split method from Slice
func (s Float64Slice) Split(k int) (Slice, Slice) {
	return s[:k], s[k:]
}

// Append method from Slice
func (s Float64Slice) Append(t Slice) Slice {
	return append(s, t.(Float64Slice)...)
}

// Replace method from Slice
func (s Float64Slice) Replace(t Slice) {
	copy(s, t.(Float64Slice))
}

// Clone method from Slice
func (s Float64Slice) Clone() Slice {
	var clone = make(Float64Slice, len(s))
	copy(clone, s)
	return clone
}

// StringSlice attaches the methods of Slice to []float64
type StringSlice []string

// At method from Slice
func (s StringSlice) At(i int) interface{} {
	return s[i]
}

// Set method from Slice
func (s StringSlice) Set(i int, v interface{}) {
	s[i] = v.(string)
}

// Len method from Slice
func (s StringSlice) Len() int {
	return len(s)
}

// Swap method from Slice
func (s StringSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Slice method from Slice
func (s StringSlice) Slice(a, b int) Slice {
	return s[a:b]
}

// Split method from Slice
func (s StringSlice) Split(k int) (Slice, Slice) {
	return s[:k], s[k:]
}

// Append method from Slice
func (s StringSlice) Append(t Slice) Slice {
	return append(s, t.(StringSlice)...)
}

// Replace method from Slice
func (s StringSlice) Replace(t Slice) {
	copy(s, t.(StringSlice))
}

// Clone method from Slice
func (s StringSlice) Clone() Slice {
	var clone = make(StringSlice, len(s))
	copy(clone, s)
	return clone
}
