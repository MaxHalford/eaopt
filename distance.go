package gago

// A Metric returns the distance between two genomes.
type Metric func(a, b Individual) float64

// A DistanceMemoizer computes and stores Metric calculations.
type DistanceMemoizer struct {
	Metric        Metric
	Distances     map[string]map[string]float64
	nCalculations int // Total number of calls to Metric for testing purposes
}

// newDistanceMemoizer initializes a DistanceMemoizer.
func newDistanceMemoizer(metric Metric) DistanceMemoizer {
	return DistanceMemoizer{
		Metric:    metric,
		Distances: make(map[string]map[string]float64),
	}
}

// GetDistance returns the distance between two Individuals based on the
// DistanceMemoizer's Metric field. If the two individuals share the same ID
// then GetDistance returns 0. DistanceMemoizer stores the calculated distances
// so that if GetDistance is called twice with the two same Individuals then
// the second call will return the stored distance instead of recomputing it.
func (dm *DistanceMemoizer) GetDistance(a, b Individual) float64 {
	// Check if the two individuals are the same before proceding
	if a.ID == b.ID {
		return 0
	}
	// Create maps if the genomes have never been encountered
	if _, ok := dm.Distances[a.ID]; !ok {
		dm.Distances[a.ID] = make(map[string]float64)
	} else {
		// Check if the distance between the two genomes has been calculated
		if dist, ok := dm.Distances[a.ID][b.ID]; ok {
			return dist
		}
	}
	if _, ok := dm.Distances[b.ID]; !ok {
		dm.Distances[b.ID] = make(map[string]float64)
	}
	// Calculate the distance between the genomes and memoize it
	var dist = dm.Metric(a, b)
	dm.Distances[a.ID][b.ID] = dist
	dm.Distances[b.ID][a.ID] = dist
	dm.nCalculations++
	return dist
}

// Return the average distance between a Individual and a slice of Individuals.
func calcAvgDistances(indis Individuals, dm DistanceMemoizer) map[string]float64 {
	var avgDistances = make(map[string]float64)
	for _, a := range indis {
		for _, b := range indis {
			avgDistances[a.ID] += dm.GetDistance(a, b)
		}
		avgDistances[a.ID] /= float64(len(indis) - 1)
	}
	return avgDistances
}
