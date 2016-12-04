package gago

import "math/rand"

var (
	ga = GA{
		MakeGenome: MakeVector,
		Topology: Topology{
			NPopulations: 2,
			NIndividuals: 50,
		},
		Model: ModGenerational{
			Selector: SelTournament{
				NParticipants: 3,
			},
			MutRate: 0.5,
		},
		Migrator:     MigRing{10},
		MigFrequency: 3,
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

func (X Vector) Evaluate() float64 {
	var sum float64
	for _, x := range X {
		sum += x
	}
	return sum
}

func (X Vector) Mutate(rng *rand.Rand) {
	MutNormalFloat64(X, rng, 0.5)
}

func (X Vector) Crossover(Y Genome, rng *rand.Rand) (Genome, Genome) {
	var o1, o2 = CrossNPointFloat64(X, Y.(Vector), 2, rng)
	return Vector(o1), Vector(o2)
}

func MakeVector(rng *rand.Rand) Genome {
	return Vector(InitUnifFloat64(4, -10, 10, rng))
}
