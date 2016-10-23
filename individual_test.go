package gago

import "testing"

func TestCopy(t *testing.T) {
	var (
		genome = MakeVector(makeRandomNumberGenerator())
		indi1  = MakeIndividual(genome)
		indi2  = indi1.DeepCopy()
	)
	if &indi1 == &indi2 || &indi1.Genome == &indi2.Genome {
		t.Error("Individual was not deep copied")
	}
}
