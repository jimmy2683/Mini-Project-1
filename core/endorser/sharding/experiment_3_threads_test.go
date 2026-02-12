package sharding_test

import (
	"testing"
	"time"
)

// Experiment 3: Throughput and transaction reject rate vs Number of threads
func TestExperiment3_Threads(t *testing.T) {
	// Powers of 2: 1, 2, 4, 8, 16, 32
	threadCounts := []int{1, 2, 4, 8, 16, 32}
	
	t.Log("Experiment 3: Varying Thread Count")
	t.Log("----------------------------------")

	for _, threads := range threadCounts {
		config := ExperimentConfig{
			FaultTolerance: 1,    // Cluster size 3
			TxCount:        1000, // Fixed transactions
			ClientCount:    threads,
			DependencyRate: 0.40, // Fixed 40% dependency
			Duration:       30 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
