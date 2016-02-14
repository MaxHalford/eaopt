package gago

import (
	"sort"
	"strconv"
	"strings"
)

// An Individual represents a potential solution to a problem. The individual's
// is defined by it's genome, which is a slice containing genes. Every gene is a
// floating point numbers. The fitness is the individual's phenotype and is
// represented by a floating point number.
type Individual struct {
	Genome  []float64
	Fitness float64
}

// Evaluate the fitness of an individual.
func (indi *Individual) Evaluate(FitnessFunction func([]float64) float64) {
	indi.Fitness = FitnessFunction(indi.Genome)
}

// floatToString converts a float into a string.
func floatToString(float float64) string {
	return strconv.FormatFloat(float, 'f', 6, 64)
}

// floatSliceToString converts a slice of floats into a string by using the
// previously defined floatToString method.
func floatSliceToString(slice []float64) string {
	var str = make([]string, len(slice))
	for i, float := range slice {
		str[i] = floatToString(float)
	}
	return strings.Join(str, ", ")
}

// Display an individual.
func (indi Individual) String() string {
	return "Fitness: " + floatToString(indi.Fitness) +
		" | Genome: [" + floatSliceToString(indi.Genome) + "]"
}

// Individuals type is necessary for sorting and selection purposes.
type Individuals []Individual

// Sort the individuals of a deme in ascending order based on their fitness. The
// convention is that we always want to minimize a function. A function f(x) can
// be function maximized by minimizing -f(x) or 1/f(x).
func (indis Individuals) Sort() {
	sort.Sort(indis)
}

func (indis Individuals) Len() int {
	return len(indis)
}

func (indis Individuals) Less(i, j int) bool {
	return indis[i].Fitness < indis[j].Fitness
}

func (indis Individuals) Swap(i, j int) {
	indis[i], indis[j] = indis[j], indis[i]
}
