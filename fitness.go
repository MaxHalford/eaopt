package gago

// FitnessFunction wraps user defined functions in order to generalize other
// functions.
type FitnessFunction interface {
	apply(genome []interface{}) float64
}

// FloatFunction is for functions with floating point slices as input.
type FloatFunction struct {
	Image func([]float64) float64
}

// Apply the fitness function in wrapped in FloatFunction.
func (ff FloatFunction) apply(genome []interface{}) float64 {
	// Apply type insertion over the genome slice
	floats := make([]float64, len(genome))
	for i, gene := range genome {
		floats[i] = gene.(float64)
	}
	// Compute the fitness
	return ff.Image(floats)
}

// StringFunction is for function with string slices as input.
type StringFunction struct {
	Image func([]string) float64
}

// Apply the fitness function in wrapped in FloatFunction.
func (sf StringFunction) apply(genome []interface{}) float64 {
	// Apply type insertion over the genome slice
	strings := make([]string, len(genome))
	for i, gene := range genome {
		strings[i] = gene.(string)
	}
	// Compute the fitness
	return sf.Image(strings)
}
