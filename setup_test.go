package eaopt

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
)

type Vector []float64

func (X Vector) Evaluate() (float64, error) {
	var sum float64
	for _, x := range X {
		sum += x
	}
	return sum, nil
}

func (X Vector) Mutate(rng *rand.Rand) {
	MutNormalFloat64(X, 0.5, rng)
}

func (X Vector) Crossover(Y Genome, rng *rand.Rand) {
	CrossUniformFloat64(X, Y.(Vector), rng)
}

func (X Vector) Clone() Genome {
	var XX = make(Vector, len(X))
	copy(XX, X)
	return XX
}

// VectorJSONUnmarshaler handles unmarshaling of JSON encoded Vectors.
func VectorJSONUnmarshaler(data []byte) (Genome, error) {
	var values []float64
	err := json.Unmarshal(data, &values)
	return Vector(values), err
}

func NewVector(rng *rand.Rand) Genome {
	return Vector(InitUnifFloat64(4, -10, 10, rng))
}

type ErrorGenome struct{}

func (eg ErrorGenome) Evaluate() (float64, error)         { return 0, errors.New("") }
func (eg ErrorGenome) Mutate(rng *rand.Rand)              {}
func (eg ErrorGenome) Crossover(Y Genome, rng *rand.Rand) {}
func (eg ErrorGenome) Clone() Genome                      { return ErrorGenome{} }

func NewErrorGenome(rng *rand.Rand) Genome { return ErrorGenome{} }

func l1Distance(x1, x2 Individual) (d float64) {
	var (
		g1 = x1.Genome.(Vector)
		g2 = x2.Genome.(Vector)
	)
	for i := range g1 {
		d += math.Abs(g1[i] - g2[i])
	}
	return
}

type ModIdentity struct{}

func (mod ModIdentity) Apply(pop *Population) error { return nil }
func (mod ModIdentity) Validate() error             { return nil }

type ModRuntimeError struct{}

func (mod ModRuntimeError) Apply(pop *Population) error { return errors.New("") }
func (mod ModRuntimeError) Validate() error             { return nil }

type ModValidateError struct{}

func (mod ModValidateError) Apply(pop *Population) error { return errors.New("") }
func (mod ModValidateError) Validate() error             { return errors.New("") }

type SpecRuntimeError struct{}

func (spec SpecRuntimeError) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	return []Individuals{indis}, errors.New("")
}
func (spec SpecRuntimeError) Validate() error { return nil }

type SpecValidateError struct{}

func (spec SpecValidateError) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	return []Individuals{indis}, errors.New("")
}
func (spec SpecValidateError) Validate() error { return errors.New("") }

type RuntimeErrorGenome struct{ Vector }

func (spec RuntimeErrorGenome) Evaluate() (float64, error) {
	return 0, errors.New("error")
}

func NewRuntimeErrorGenome(rng *rand.Rand) Genome {
	return RuntimeErrorGenome{[]float64{42}}
}
