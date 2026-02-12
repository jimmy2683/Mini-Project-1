package sharding_test

import (
	"testing"
	"time"
)

// Experiment 1: Throughput and transaction reject rate vs Number of transactions in a block
// Note: In our simulation, "Transactions in a block" is proxied by total load since block size is internal to Raft.
// We vary the total load to see how the system handles increasing pressure.
func TestExperiment1_Transactions(t *testing.T) {
	txCounts := []int{1000, 2000, 3000, 4000, 5000}
	
	t.Log("Experiment 1: Varying Transaction Count")
	t.Log("---------------------------------------")

	for _, txCount := range txCounts {
		config := ExperimentConfig{
			FaultTolerance: 1,    // Cluster size 3 (N=2f+1)
			TxCount:        txCount,
			ClientCount:    32,   // Fixed threads
			DependencyRate: 0.40, // Fixed 40% dependency
			Duration:       30 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
