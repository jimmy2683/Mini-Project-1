package sharding_test

import (
	"testing"
	"time"
)

// Experiment 5: Response Time vs Transactions
// Measures total response time and commit response time
func TestExperiment5_ResponseTime(t *testing.T) {
	txCounts := []int{1000, 2000, 3000, 4000, 5000}
	
	t.Log("Experiment 5: Response Time Analysis")
	t.Log("------------------------------------")

	// Note: RunExperiment prints "Throughput", but since Throughput = Count / Time,
	// Inverse of throughput per tx is roughly Avg Latency if concurrency is low,
	// but with concurrency it's more complex. 
	// The current RunExperiment implementation prints aggregated stats.
	// For this specific test, we interpret the output duration as a proxy for efficiency.
	
	for _, txCount := range txCounts {
		config := ExperimentConfig{
			FaultTolerance: 1,    // Cluster size 3
			TxCount:        txCount,
			ClientCount:    32,
			DependencyRate: 0.40,
			Duration:       30 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
