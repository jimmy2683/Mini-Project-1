package sharding_test

import (
	"testing"
	"time"
)

// Experiment 2: Throughput and transaction reject rate vs Dependency in a block
func TestExperiment2_Dependency(t *testing.T) {
	dependencyRates := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5}
	
	t.Log("Experiment 2: Varying Dependency Rate")
	t.Log("-------------------------------------")

	for _, depRate := range dependencyRates {
		config := ExperimentConfig{
			FaultTolerance: 1,    // Cluster size 3
			TxCount:        1000, // Fixed transactions
			ClientCount:    32,   // Fixed threads
			DependencyRate: depRate,
			Duration:       30 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
