package gago

// A Genome contains the genes of an individual.
type Genome []interface{}

// CastFloat casts each to a float.
func (g Genome) CastFloat() []float64 {
	var casted = make([]float64, len(g))
	for i := range g {
		casted[i] = g[i].(float64)
	}
	return casted
}

// CastString casts each gene to a string.
func (g Genome) CastString() []string {
	var casted = make([]string, len(g))
	for i := range g {
		casted[i] = g[i].(string)
	}
	return casted
}
