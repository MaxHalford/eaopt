package eaopt

import (
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

func NewVector(rng *rand.Rand) Genome {
	return Vector(InitUnifFloat64(4, -10, 10, rng))
}

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

type SpecRuntimeError struct{}

func (spec SpecRuntimeError) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	return []Individuals{indis, indis}, errors.New("")
}
func (spec SpecRuntimeError) Validate() error { return nil }

type RuntimeErrorGenome struct{ Vector }

func (spec RuntimeErrorGenome) Evaluate() (float64, error) {
	return 0, errors.New("error")
}

func NewRuntimeErrorGenome(rng *rand.Rand) Genome {
	return RuntimeErrorGenome{[]float64{42}}
}
