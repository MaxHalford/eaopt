package gago

import (
	"log"
	"math"
	"math/rand"
	"os"
)

var (
	ga = GA{
		MakeGenome: MakeVector,
		NPops:      2,
		PopSize:    50,
		Model: ModGenerational{
			Selector: SelTournament{
				NParticipants: 3,
			},
			MutRate: 0.5,
		},
		Migrator:     MigRing{10},
		MigFrequency: 3,
		Logger:       log.New(os.Stdin, "", log.Ldate|log.Ltime),
	}
	nbrGenerations = 5 // Initial number of generations to enhance
)

func init() {
	ga.Initialize()
	for i := 0; i < nbrGenerations; i++ {
		ga.Enhance()
	}
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

func (X Vector) Crossover(Y Genome, rng *rand.Rand) (Genome, Genome) {
	var o1, o2 = CrossGNXFloat64(X, Y.(Vector), 2, rng)
	return Vector(o1), Vector(o2)
}

// Implement the Slice interface

func (X Vector) Swap(i, j int) {
	X[i], X[j] = X[j], X[i]
}

func (X Vector) Len() int {
	return len(X)
}

func MakeVector(rng *rand.Rand) Genome {
	return Vector(InitUnifFloat64(4, -10, 10, rng))
}

// Minkowski distance with p = 1
func l1Distance(x1, x2 Individual) (dist float64) {
	var g1 = x1.Genome.(Vector)
	var g2 = x2.Genome.(Vector)
	for i := range g1 {
		dist += math.Abs(g1[i] - g2[i])
	}
	return
}

// Identity model

type ModIdentity struct{}

func (mod ModIdentity) Apply(pop *Population) {}
func (mod ModIdentity) Validate() error       { return nil }
