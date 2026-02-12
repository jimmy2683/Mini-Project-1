package sharding_test

import (
	"testing"
	"time"
)

// Experiment 4: Throughput and transaction reject rate vs Size of the cluster
func TestExperiment4_ClusterSize(t *testing.T) {
	// Cluster sizes: 1 (f=0), 3 (f=1), 5 (f=2), 7 (f=3)
	// Relationship: Size = 2f + 1 => f = (Size - 1) / 2
	faultTolerances := []int{0, 1, 2, 3}
	
	t.Log("Experiment 4: Varying Cluster Size (Fault Tolerance)")
	t.Log("--------------------------------------------------")

	for _, f := range faultTolerances {
		config := ExperimentConfig{
			FaultTolerance: f,
			TxCount:        1000,
			ClientCount:    32,
			DependencyRate: 0.40,
			Duration:       30 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
