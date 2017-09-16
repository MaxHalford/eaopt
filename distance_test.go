package gago

import (
	"errors"
	"fmt"
	"testing"
)

func TestDistanceMemoizer(t *testing.T) {
	var (
		dm = newDistanceMemoizer(l1Distance)
		a  = Individual{Genome: Vector{1, 1, 1}, ID: "1"}
		b  = Individual{Genome: Vector{3, 3, 3}, ID: "2"}
		c  = Individual{Genome: Vector{6, 6, 6}, ID: "3"}
	)
	// Check the number of calculations is initially 0
	if dm.nCalculations != 0 {
		t.Error("nCalculations is not initialized to 0")
	}
	// Check the distance between the 1st and itself
	if dm.GetDistance(a, a) != 0 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations is initially 0
	if dm.nCalculations != 0 {
		t.Error("nCalculations should not have increased")
	}
	// Check the distance between the 1st and the 2nd individual
	if dm.GetDistance(a, b) != 6 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations has increased by 1
	if dm.nCalculations != 1 {
		t.Error("nCalculations has not increased")
	}
	// Check the distance between the 2nd and the 1st individual
	if dm.GetDistance(b, a) != 6 {
		t.Error("Wrong calculated distance")
	}
	// Check the number of calculations has not increased
	if dm.nCalculations != 1 {
		t.Error("nCalculations has increased")
	}
	// Check the distance between the 1st and the 3rd individual
	if dm.GetDistance(a, c) != 15 {
		t.Error("Wrong calculated distance")
	}
	// Check the distance between the 1st and the 3rd individual
	if dm.GetDistance(b, c) != 9 {
		t.Error("Wrong calculated distance")
	}
}

func TestSortByDistanceToMedoid(t *testing.T) {
	var (
		indis = Individuals{
			Individual{Genome: Vector{3, 3, 3}, Fitness: 0},
			Individual{Genome: Vector{2, 2, 2}, Fitness: 1},
			Individual{Genome: Vector{5, 5, 5}, Fitness: 2},
		}
		dm = newDistanceMemoizer(l1Distance)
	)
	indis.SortByDistanceToMedoid(dm)
	for i := range indis {
		if indis[i].Fitness != float64(i) {
			t.Error("Individuals were not sorted according to their distance to the medoid")
		}
	}
}

func TestRebalanceClusters(t *testing.T) {
	var testCases = []struct {
		clusters        []Individuals
		dm              DistanceMemoizer
		minClusterSize  int
		newClusterSizes []int
		err             error
	}{
		{
			clusters: []Individuals{
				Individuals{
					Individual{Genome: Vector{1, 1, 1}, ID: "1"},
					Individual{Genome: Vector{1, 1, 1}, ID: "2"},
					Individual{Genome: Vector{1, 1, 1}, ID: "3"},
					Individual{Genome: Vector{2, 2, 2}, ID: "4"}, // Second furthest away from the cluster
					Individual{Genome: Vector{3, 3, 3}, ID: "5"}, // Furthest away from the cluster
				},
				Individuals{
					Individual{Genome: Vector{2, 2, 2}, ID: "6"},
				},
				Individuals{
					Individual{Genome: Vector{3, 3, 3}, ID: "7"},
				},
			},
			dm:              newDistanceMemoizer(l1Distance),
			minClusterSize:  2,
			newClusterSizes: []int{3, 2, 2},
			err:             nil,
		},
		{
			clusters: []Individuals{
				Individuals{
					Individual{Genome: Vector{1, 1, 1}, ID: "1"},
					Individual{Genome: Vector{1, 1, 1}, ID: "2"},
				},
				Individuals{},
			},
			dm:              newDistanceMemoizer(l1Distance),
			minClusterSize:  1,
			newClusterSizes: []int{2, 0},
			err:             errors.New("Cluster number 2 is empty"),
		},
		{
			clusters: []Individuals{
				Individuals{
					Individual{Genome: Vector{1, 1, 1}, ID: "1"},
					Individual{Genome: Vector{1, 1, 1}, ID: "2"},
				},
				Individuals{
					Individual{Genome: Vector{1, 1, 1}, ID: "3"},
				},
			},
			dm:              newDistanceMemoizer(l1Distance),
			minClusterSize:  2,
			newClusterSizes: []int{2, 0},
			err:             errors.New("Not enough individuals to rebalance"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var err = rebalanceClusters(tc.clusters, tc.dm, tc.minClusterSize)
			// Check if the error is nil or not
			if (err == nil) != (tc.err == nil) {
				t.Errorf("Wrong error in test case number %d", i)
			}
			// Check new cluster sizes
			if err == nil {
				for j, cluster := range tc.clusters {
					if len(cluster) != tc.newClusterSizes[j] {
						t.Errorf("Wrong new cluster size in test case number %d", i)
					}
				}
			}
		})
	}
}

// If a cluster is empty then rebalancing is impossible
func TestRebalanceClustersEmptyCluster(t *testing.T) {
	var (
		clusters = []Individuals{
			Individuals{
				Individual{Genome: Vector{1, 1, 1}, ID: "1"},
				Individual{Genome: Vector{1, 1, 1}, ID: "2"},
				Individual{Genome: Vector{1, 1, 1}, ID: "3"},
			},
			Individuals{},
		}
		dm = newDistanceMemoizer(l1Distance)
	)
	var err = rebalanceClusters(clusters, dm, 2)
	if err == nil {
		t.Error("rebalanceClusters should have returned an error")
	}
}

// It's impossible to put 2 Individuals inside each cluster if there are 3
// clusters and 5 individuals in total
func TestRebalanceClustersTooManyMissing(t *testing.T) {
	var (
		clusters = []Individuals{
			Individuals{
				Individual{Genome: Vector{1, 1, 1}, ID: "1"},
				Individual{Genome: Vector{1, 1, 1}, ID: "2"},
				Individual{Genome: Vector{1, 1, 1}, ID: "3"},
			},
			Individuals{
				Individual{Genome: Vector{2, 2, 2}, ID: "6"},
			},
			Individuals{
				Individual{Genome: Vector{3, 3, 3}, ID: "7"},
			},
		}
		dm = newDistanceMemoizer(l1Distance)
	)
	var err = rebalanceClusters(clusters, dm, 2)
	if err == nil {
		t.Error("rebalanceClusters should have returned an error")
	}
}
