package mutate

import "math/rand"

// Splice splits a genome in 3 and glues the parts back together in another
// order.
func Splice(genome []interface{}, rng *rand.Rand) {
	// Choose where to start and end the splice
	var (
		end   = rng.Intn(len(genome)-1) + 1
		start = rng.Intn(end)
	)
	// Split the genome into two
	var inner = make([]interface{}, end-start)
	copy(inner, genome[start:end])
	var outer = append(genome[:start], genome[end:]...)
	// Choose where to insert the splice
	var insert = rng.Intn(len(outer))
	// Splice and insert
	genome = append(
		outer[:insert],
		append(
			inner,
			outer[insert:]...,
		)...,
	)
}

// Convenience types for common cases

// Splice a []int.
func (ints IntSlice) Splice(rng *rand.Rand) {
	var genome = make([]interface{}, len(ints))
	for i, v := range ints {
		genome[i] = v
	}
	Splice(genome, rng)
}

// Splice a []float64.
func (floats Float64Slice) Splice(rng *rand.Rand) {
	var genome = make([]interface{}, len(floats))
	for i, v := range floats {
		genome[i] = v
	}
	Splice(genome, rng)
}

// Splice a []string.
func (strings StringSlice) Splice(rng *rand.Rand) {
	var genome = make([]interface{}, len(strings))
	for i, v := range strings {
		genome[i] = v
	}
	Splice(genome, rng)
}
