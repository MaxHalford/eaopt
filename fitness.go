package gago

// FitnessFunction wraps user defined functions in order to generalize other
// functions.
type FitnessFunction interface {
	apply(genome Genome) float64
}

// FloatFunction is for functions with floating point slices as input.
type FloatFunction struct {
	Image func([]float64) float64
}

// Apply the fitness function in wrapped in FloatFunction.
func (ff FloatFunction) apply(genome Genome) float64 {
	return ff.Image(genome.CastFloat())
}

// StringFunction is for function with string slices as input.
type StringFunction struct {
	Image func([]string) float64
}

// Apply the fitness function in wrapped in StringFunction.
func (sf StringFunction) apply(genome Genome) float64 {
	return sf.Image(genome.CastString())
}
