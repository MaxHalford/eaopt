package gago

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

var newGA = func() GA {
	return GA{
		NewGenome: NewVector,
		NPops:     2,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		Migrator:     MigRing{10},
		MigFrequency: 3,
		RNG:          rand.New(rand.NewSource(42)),
		Logger:       log.New(os.Stdin, "", log.Ldate|log.Ltime),
	}
}

// newRand returns a new random number generator with a random seed.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Vector []float64

// Implement the Genome interface

func (X Vector) Evaluate() float64 {
	var sum float64
	for _, x := range X {
		sum += x
	}
	return sum
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

// Minkowski distance with p = 1
func l1Distance(x1, x2 Individual) (dist float64) {
	var (
		g1 = x1.Genome.(Vector)
		g2 = x2.Genome.(Vector)
	)
	for i := range g1 {
		dist += math.Abs(g1[i] - g2[i])
	}
	return
}

// Identity model

type ModIdentity struct{}

func (mod ModIdentity) Apply(pop *Population) error { return nil }
func (mod ModIdentity) Validate() error             { return nil }

// Runtime error model

type ModRuntimeError struct{}

func (mod ModRuntimeError) Apply(pop *Population) error { return errors.New("") }
func (mod ModRuntimeError) Validate() error             { return nil }

// Runtime error speciator

type SpecRuntimeError struct{}

func (spec SpecRuntimeError) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	return []Individuals{indis, indis}, errors.New("")
}
func (spec SpecRuntimeError) Validate() error { return nil }
