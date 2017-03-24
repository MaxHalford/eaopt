package gago

type DistanceMemoizer struct {
	Metric    func(a, b Genome) float64
	Distances map[Genome]map[Genome]float64
}

func makeDistanceMemoizer(metric func(a, b Genome) float64) DistanceMemoizer {
	return DistanceMemoizer{
		Metric:    metric,
		Distances: make(map[Genome]map[Genome]float64),
	}
}

func (dm *DistanceMemoizer) getDistance(a, b Genome) float64 {
	if dist, ok := dm.Distances[a][b]; ok {
		return dist
	}
	if dist, ok := dm.Distances[b][a]; ok {
		return dist
	}
	var dist = dm.Metric(a, b)
	dm.Distances[a][b] = dist
	dm.Distances[b][a] = dist
	return dist
}
