package gago

import (
	"math"
	"testing"
)

func L1Distance(x1, x2 Genome) (dist float64) {
	var g1 = x1.(Vector)
	var g2 = x2.(Vector)
	for i := range g1 {
		dist += math.Abs(g1[i] - g2[i])
	}
	return
}

func TestDistanceMemoizer(t *testing.T) {
	var dm = makeDistanceMemoizer(L1Distance)
}
