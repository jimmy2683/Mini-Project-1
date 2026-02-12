package sharding_test

import (
	"testing"
	"time"
)

// Experiment 6: High Scale / Resilience Test
// Replicates the logic of the original experiments_test.go but with higher fault tolerance values
func TestExperiment6_Scale(t *testing.T) {
	// Test larger clusters
	faultTolerances := []int{4, 5} // Cluster sizes 9 and 11
	
	t.Log("Experiment 6: High Scale Cluster Test")
	t.Log("-------------------------------------")

	for _, f := range faultTolerances {
		config := ExperimentConfig{
			FaultTolerance: f,
			TxCount:        1000,
			ClientCount:    32,
			DependencyRate: 0.40,
			Duration:       60 * time.Second, // Allow more time for leader election in large clusters
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
