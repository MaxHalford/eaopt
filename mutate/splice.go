package mutate

import "math/rand"

// Splice splits a genome in 3 and glues the parts back together in another
// order.
func splice(genome Castable, rng *rand.Rand) []interface{} {
	var slice = genome.Cast()
	// Choose where to start and end the splice
	var (
		end   = rng.Intn(len(slice)-1) + 1
		start = rng.Intn(end)
	)
	// Split the genome into two
	var inner = make([]interface{}, end-start)
	copy(inner, slice[start:end])
	var outer = append(slice[:start], slice[end:]...)
	// Choose where to insert the splice
	var insert = rng.Intn(len(outer))
	// Splice and insert
	slice = append(
		outer[:insert],
		append(
			inner,
			outer[insert:]...,
		)...,
	)
	return slice
}

func SpliceFloat64s(floats []float64, rng *rand.Rand) []interface{} {
	return splice(Float64Slice(floats), rng)
}
