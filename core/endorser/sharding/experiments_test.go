package sharding_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/endorser/sharding"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

// ExperimentConfig holds parameters for running an experiment
type ExperimentConfig struct {
	FaultTolerance int           // f (N = 2f + 1)
	TxCount        int           // Total transactions
	ClientCount    int           // Number of concurrent clients (threads)
	DependencyRate float64       // 0.0 to 1.0 (percentage of txs that conflict)
	Duration       time.Duration // Max duration
	LossProbability float64      // Network packet loss probability
}

// Network simulates the network between Raft nodes
type Network struct {
	nodes   map[uint64]*sharding.ShardLeader
	latency time.Duration
	loss    float64 // Probability of packet loss (0.0 - 1.0)
	mu      sync.RWMutex
}

func NewNetwork(latency time.Duration, loss float64) *Network {
	return &Network{
		nodes:   make(map[uint64]*sharding.ShardLeader),
		latency: latency,
		loss:    loss,
	}
}

func (n *Network) AddNode(id uint64, node *sharding.ShardLeader) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.nodes[id] = node
}

func (n *Network) Run() {
	// In a real implementation, this would start a background router.
	// For this test, nodes push to the network via helper methods.
}

func (n *Network) RouteMessages(fromID uint64, msgs []raftpb.Message) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, msg := range msgs {
		msg := msg // copy
		// Simulate packet loss
		if n.loss > 0 && rand.Float64() < n.loss {
			continue
		}

		target := n.nodes[msg.To]
		if target == nil {
			continue
		}

		// Simulate latency
		go func(t *sharding.ShardLeader, m raftpb.Message) {
			if n.latency > 0 {
				time.Sleep(n.latency)
			}
			t.Step(context.TODO(), m)
		}(target, msg)
	}
}

// RunExperiment executes a single experiment with the given configuration
func RunExperiment(t *testing.T, config ExperimentConfig) {
	f := config.FaultTolerance
	n := 2*f + 1
	nodes := make([]*sharding.ShardLeader, n)
	net := NewNetwork(5*time.Millisecond, config.LossProbability)

	// Create nodes
	replicaNodes := make([]string, n)
	for i := 0; i < n; i++ {
		replicaNodes[i] = fmt.Sprintf("node%d", i+1)
	}

	for i := 0; i < n; i++ {
		shardConfig := sharding.ShardConfig{
			ShardID:      "experiment-shard",
			ReplicaNodes: replicaNodes,
			ReplicaID:    uint64(i + 1),
		}
		
		node, err := sharding.NewShardLeader(shardConfig, 100*time.Millisecond, 50) 
		if err != nil {
			t.Fatalf("Failed to create node %d: %v", i+1, err)
		}
		nodes[i] = node
		net.AddNode(uint64(i+1), node)
		
		// Start message router for this node
		go func(node *sharding.ShardLeader, id uint64) {
			for msgs := range node.MessagesC() {
				net.RouteMessages(id, msgs)
			}
		}(node, uint64(i+1))
	}

	// Wait for leader election
	// We just push some traffic to trigger it, simple wait
	time.Sleep(2 * time.Second)

	// Run workload
	var wg sync.WaitGroup
	startTime := time.Now()
	successCount := 0
	var mu sync.Mutex

	// Send transactions to ANY node
	workloadCh := make(chan int, config.TxCount)
	for i := 0; i < config.TxCount; i++ {
		workloadCh <- i
	}
	close(workloadCh)

	// Dependency management
	// We simulate dependency by having a "hot key" that many transactions try to write to
	hotKey := "hot-key"
	normalKeyPrefix := "key-"

	// Simulate client threads
	for c := 0; c < config.ClientCount; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range workloadCh {
				// Pick random node
				nodeID := rand.Intn(n)
				node := nodes[nodeID]

				// Determine if this tx has a dependency (accesses hot key)
				isDependent := false
				if rand.Float64() < config.DependencyRate {
					isDependent = true
				}

				key := fmt.Sprintf("%s%d", normalKeyPrefix, i)
				if isDependent {
					key = hotKey // Contention on this key
				}

				req := &sharding.PrepareRequest{
					TxID:      fmt.Sprintf("tx-%d-%d", f, i),
					ShardID:   "experiment-shard",
					WriteSet:  map[string][]byte{key: []byte("value")},
					Timestamp: time.Now(),
				}

				// Send to node
				// Note: ProposeC might block if full, so we select
				select {
				case node.ProposeC() <- req:
					// Request sent
				case <-time.After(500 * time.Millisecond):
					// Send timeout
				}
			}
		}()
	}

	// Monitor Commit Channels
	committedTxs := make(map[string]bool)
	allCommits := make(chan *sharding.PrepareProof, 10000)
	for _, node := range nodes {
		go func(c <-chan *sharding.PrepareProof) {
			for proof := range c {
				allCommits <- proof
			}
		}(node.CommitC())
	}

	// Run for duration or until done
	timeout := time.After(config.Duration)
	
Loop:
	for {
		select {
		case proof := <-allCommits:
			mu.Lock()
			if !committedTxs[proof.TxID] {
				committedTxs[proof.TxID] = true
				successCount++
			}
			completion := float64(len(committedTxs)) / float64(config.TxCount)
			mu.Unlock()
			
			if completion >= 1.0 {
				break Loop
			}
			
		case <-timeout:
			break Loop
		}
	}

	elapsed := time.Since(startTime)
	if elapsed > config.Duration {
		elapsed = config.Duration
	}

	throughput := float64(successCount) / elapsed.Seconds()
	
	fmt.Printf("Config: f=%d, Txs=%d, Clients=%d, Dep=%.2f -> Throughput: %.2f tx/s, Success: %.2f%%\n",
		f, config.TxCount, config.ClientCount, config.DependencyRate, throughput, float64(successCount)/float64(config.TxCount)*100)

	// Cleanup
	for _, node := range nodes {
		node.Stop()
	}
}

// Preserve original test function for backward compatibility
func TestExperiments(t *testing.T) {
	// Fault tolerance levels to test: f=0, 1, 2, 3
	fs := []int{0, 1, 2, 3}
	
	fmt.Printf("Starting Sharded Raft Experiments (Original Loop)\n")
	fmt.Printf("=================================\n")

	for _, f := range fs {
		config := ExperimentConfig{
			FaultTolerance: f,
			TxCount:        1000,
			ClientCount:    10,
			DependencyRate: 0.0,
			Duration:       10 * time.Second,
			LossProbability: 0.0,
		}
		RunExperiment(t, config)
	}
}
