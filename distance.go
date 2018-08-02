package eaopt

import (
	"fmt"
	"math"
)

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

// calcAvgDistances returns a map that associates the ID of each provided
// Individual with the average distance the Individual has with the rest of the
// Individuals.
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

func rebalanceClusters(clusters []Individuals, dm DistanceMemoizer, minPerCluster uint) error {
	// Calculate the number of missing Individuals per cluster for each cluster
	// to reach at least minPerCluster Individuals.
	var missing = make([]int, len(clusters))
	for i, cluster := range clusters {
		// Check that the cluster has at least one Individual
		if len(cluster) == 0 {
			return fmt.Errorf("Cluster %d has 0 individuals", i)
		}
		// Calculate the number of missing Individual in the cluster to reach minPerCluster
		missing[i] = int(minPerCluster) - len(cluster)
	}
	// Check if there are enough Individuals to rebalance the clusters.
	if sumInts(missing) >= 0 {
		return fmt.Errorf("Missing %d individuals to be able to rebalance the clusters",
			sumInts(missing))
	}
	// Loop through the clusters that are missing Individuals
	for i, cluster := range clusters {
		// Check if the cluster is missing Individuals
		if missing[i] <= 0 {
			continue
		}
		// Assign new Individuals to the cluster while it is missing some
		for missing[i] > 0 {
			// Determine the medoid
			cluster.SortByDistanceToMedoid(dm)
			var medoid = cluster[0]
			// Go through the Individuals of the other clusters and find the one
			// closest to the computed medoid
			var (
				cci     int // Closest cluster index
				cii     int // Closest Individual index
				minDist = math.Inf(1)
			)
			for j := range clusters {
				// Check that the cluster has Individuals to spare
				if i == j || missing[j] >= 0 {
					continue
				}
				// Find the closest Individual to the medoid inside the cluster
				for k, indi := range clusters[j] {
					var dist = dm.GetDistance(indi, medoid)
					if dist < minDist {
						cci = j
						cii = k
						minDist = dist
					}
				}
			}
			// Add the closest Individual to the cluster
			clusters[i] = append(clusters[i], clusters[cci][cii])
			// Remove the closest Individual from the cluster it belonged to
			clusters[cci] = append(clusters[cci][:cii], clusters[cci][cii+1:]...)
			missing[i]--
		}
	}
	return nil
}
