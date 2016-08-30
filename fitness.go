package gago

// FitnessFunction wraps user defined functions in order to generalize other functions.
type FitnessFunction interface {
	Apply(genome Genome) float64
}

// Float64Function is for functions with floating point slices as input.
type Float64Function struct {
	Image func([]float64) float64
}

// Apply the fitness function wrapped in FloatFunction.
func (ff Float64Function) Apply(genome Genome) float64 {
	var casted = make([]float64, len(genome))
	for i := range genome {
		casted[i] = genome[i].(float64)
	}
	return ff.Image(casted)
}

// StringFunction is for functions with string slices as input.
type StringFunction struct {
	Image func([]string) float64
}

// Apply the fitness function wrapped in StringFunction.
func (ff StringFunction) Apply(genome Genome) float64 {
	var casted = make([]string, len(genome))
	for i := range genome {
		casted[i] = genome[i].(string)
	}
	return ff.Image(casted)
}
