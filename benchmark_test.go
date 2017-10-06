package gago

import (
	"runtime"
	"testing"
)

func BenchmarkEnhance1Pop(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     1,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}
	ga.Initialize()
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}

func BenchmarkEnhance2Pops(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ga = GA{
		NewGenome: NewVector,
		NPops:     2,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}
	ga.Initialize()
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}

func BenchmarkEnhance3Pops(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     3,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}
	ga.Initialize()
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}
func BenchmarkEnhance4Pops(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     4,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
	}
	ga.Initialize()
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}
