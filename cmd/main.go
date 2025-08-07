package main

import (
	"blockchain-system/core"
	"fmt"
	"math/big"
	"time"
)

func main() {
	// === 1. Blockchain Initialization ===
	bc := core.NewBlockchain()
	bc.AddBlock("First Block after Genesis")
	bc.AddBlock("Second Block")
	bc.AddBlock("Third Block")
	bc.AddBlock("Fourth Block")
	bc.AddBlock("Fifth Block")
	bc.AddBlock("Sixth Block")
	bc.AddBlock("Seventh Block")
	bc.AddBlock("Eighth Block")

	for _, block := range bc.Blocks {
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Data:", block.Data)
		fmt.Println("Prev Hash:", block.PrevHash)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("====================================")
	}

	// === 2. Merkle Forest (Shard Distribution) ===
	sm := core.NewShardManager()
	for _, block := range bc.Blocks {
		sm.DistributeBlock(block)
	}
	sm.PrintShardState()

	// Demonstrate logarithmic-time shard discovery
	fmt.Println("\n[INFO] Demonstrating logarithmic-time shard discovery")
	shardID := 0 // Example shard ID
	if shard, exists := sm.FindShard(shardID); exists {
		fmt.Printf("Found Shard #%d with %d blocks\n", shard.ID, len(shard.Blocks))
	} else {
		fmt.Printf("Shard #%d not found\n", shardID)
	}

	// Demonstrate state reconstruction
	fmt.Println("\n[INFO] Demonstrating state reconstruction")
	if root, exists := sm.ReconstructState(shardID); exists {
		fmt.Printf("Shard #%d Merkle Root: %s\n", shardID, root)
	} else {
		fmt.Printf("Cannot reconstruct state for Shard #%d\n", shardID)
	}

	// === 3. Atomic Cross-Shard Transfer with Homomorphic Authentication ===
	fmt.Println("\n[INFO] Simulating atomic cross-shard transfer")
	// Use enhanced sync manager with homomorphic authentication
	enhancedSyncManager := core.NewEnhancedSyncManager("secret-key-123")
	shards := sm.Shards.GetAllShards()
	if len(shards) >= 2 {
		// Successful transfer
		fmt.Println("\n[INFO] Attempting successful transfer")
		if enhancedSyncManager.CreateAuthenticatedTransfer(shards[0], shards[1], 0) {
			verified := enhancedSyncManager.VerifyAndApplyTransfer(shards[0], shards[1], 0)
			fmt.Printf("Transfer from Shard #%d to #%d: %v\n", shards[0].ID, shards[1].ID, verified)
		}
		sm.PrintShardState()

		// Simulate failed transfer (e.g., invalid commitment)
		fmt.Println("\n[INFO] Attempting failed transfer to demonstrate rollback")
		// Use a new sync manager with a wrong key to force rollback
		faultySyncManager := core.NewEnhancedSyncManager("wrong-key")
		if faultySyncManager.CreateAuthenticatedTransfer(shards[0], shards[1], 0) {
			verified := faultySyncManager.VerifyAndApplyTransfer(shards[0], shards[1], 0)
			fmt.Printf("Transfer from Shard #%d to #%d: %v (should fail and rollback)\n", shards[0].ID, shards[1].ID, verified)
		}
		sm.PrintShardState()
	} else {
		fmt.Println("Not enough shards for transfer demo")
	}

	// === 4. Shard Merging ===
	fmt.Println("\nChecking for underutilized shards to merge...")
	sm.MergeShards(2)
	sm.PrintShardState()

	// === 5. BFT Consensus Round ===
	bft := core.NewBFTManager(10)
	bft.RunConsensus()

	// === 6. Hybrid Consensus ===
	consensus := &core.ConsensusManager{BFT: bft}
	consensus.RunHybridConsensus()

	// === 7. Zero-Knowledge Proof Demo ===
	zk := &core.ZKProver{}
	zk.TestZKP()

	// === 8. RSA Cryptographic Accumulator Demo ===
	fmt.Println("\n=== RSA Cryptographic Accumulator Demonstration ===")
	acc := core.NewRSAAccumulator()
	for _, block := range bc.Blocks[1:4] { // Use first three blocks after genesis
		acc.AddElement(block.Hash)
		fmt.Printf("Added block #%d hash: %s\n", block.Index, block.Hash)
	}
	fmt.Printf("Accumulator State: %s\n", acc.State.Text(16))
	// Verify membership for a block
	if proof, exists := acc.Proofs[bc.Blocks[1].Hash]; exists {
		valid := acc.VerifyMembership(bc.Blocks[1].Hash, proof)
		fmt.Printf("Membership proof for block #%d: %v\n", bc.Blocks[1].Index, valid)
	}
	// Test non-member
	nonMember := "invalid_hash"
	fakeProof := new(big.Int).Set(acc.G)
	valid := acc.VerifyMembership(nonMember, fakeProof)
	fmt.Printf("Membership proof for invalid hash: %v\n", nonMember, valid)

	// === 9. State Pruning + Compact State Representation ===
	fmt.Println("\n=== State Pruning + Compact State Representation ===")
	smgr := core.NewStateManager(2)
	for _, block := range bc.Blocks {
		smgr.AddBlock(block)
	}
	smgr.PrintState()

	// Demonstrate retrieving data from trie
	fmt.Println("\nRetrieving block data from succinct trie:")
	for _, block := range bc.Blocks[:2] { // First two blocks should be in archive trie
		if data, exists := smgr.ArchiveTrie.Get(block.Hash); exists {
			fmt.Printf("Block #%d (Archived) - Data: %s\n", block.Index, data)
		}
	}
	for _, block := range bc.Blocks[2:] { // Remaining blocks in active trie
		if data, exists := smgr.ActiveTrie.Get(block.Hash); exists {
			fmt.Printf("Block #%d (Active) - Data: %s\n", block.Index, data)
		}
	}

	// === 10. Multi-Party Computation Demonstration ===
	fmt.Println("\n=== Multi-Party Computation Demonstration ===")
	// Create nodes for MPC simulation
	nodes := []*core.Node{
		{ID: 0, Byzantine: false},
		{ID: 1, Byzantine: false},
		{ID: 2, Byzantine: false},
		{ID: 3, Byzantine: true}, // Byzantine node
		{ID: 4, Byzantine: false},
	}

	// Setup MPC protocol with threshold 3 (need at least 3 honest nodes)
	mpcProtocol := core.NewMPCProtocol(nodes, 3)

	// Simulate threshold signing
	signature := mpcProtocol.SimulateMPCSignature("Important blockchain message")
	fmt.Println("MPC Signature verification:", len(signature) > 0)

	// === 11. Probabilistic Verification with Bloom Filters ===
	fmt.Println("\n=== Probabilistic Verification Demonstration ===")

	// Create test data for Merkle tree
	testData := []string{
		"Transaction 1",
		"Transaction 2",
		"Transaction 3",
		"Transaction 4",
		"Transaction 5",
	}

	// Create a proof-compressing Merkle tree
	pcmt := core.NewProofCompressingMerkleTree(testData)
	fmt.Println("Merkle Root Hash:", pcmt.GetRootHash())

	// Verify data probabilistically
	fmt.Println("Probabilistic verification of 'Transaction 2':",
		pcmt.VerifyDataProbabilistic("Transaction 2"))
	fmt.Println("Probabilistic verification of 'Transaction 6':",
		pcmt.VerifyDataProbabilistic("Transaction 6"))

	// Create bloom filter directly for demonstration
	bloomFilter := core.NewBloomFilter(1024, 3)
	for _, tx := range testData {
		bloomFilter.Add(tx)
	}
	fmt.Printf("Bloom filter false positive rate: %.6f%%\n",
		bloomFilter.FalsePositiveRate()*100)

	// === 12. Advanced CAP Optimization Test ===
	fmt.Println("\n=== Advanced CAP Theorem Optimization Test ===")
	adaptiveCAP := core.NewAdaptiveCapacityManager("node1")
	demoAdaptiveCAP(adaptiveCAP)
	fmt.Println("Advanced CAP Optimization Test Complete")

	// === 13. Homomorphic Commitment Demonstration ===
	fmt.Println("\n=== Homomorphic Commitment Demonstration ===")
	// Create homomorphic authenticator
	auth := core.NewHomomorphicAuthenticator("secret-commitment-key")

	// Create individual commitments
	commitment1 := core.HomomorphicCommitment{
		Value:      "Data piece 1",
		Commitment: auth.AuthenticateData("Data piece 1"),
	}

	commitment2 := core.HomomorphicCommitment{
		Value:      "Data piece 2",
		Commitment: auth.AuthenticateData("Data piece 2"),
	}

	// Combine commitments
	combinedCommitment := auth.CombineCommitments([]core.HomomorphicCommitment{
		commitment1, commitment2,
	})

	fmt.Println("Combined value:", combinedCommitment.Value)
	fmt.Println("Combined commitment:", combinedCommitment.Commitment)

	// Verify individual commitments
	fmt.Println("Verification of commitment 1:",
		auth.VerifyAuthentication(commitment1.Value, commitment1.Commitment))
	fmt.Println("Verification of commitment 2:",
		auth.VerifyAuthentication(commitment2.Value, commitment2.Commitment))

	// === 14. State Pruning with Cryptographic Integrity ===
	fmt.Println("\n=== State Pruning with Cryptographic Integrity ===")
	// Create a larger blockchain for demonstration
	prunableBC := core.NewBlockchain()
	for i := 0; i < 20; i++ {
		prunableBC.AddBlock(fmt.Sprintf("Block %d for pruning demo", i))
	}
	fmt.Printf("Created blockchain with %d blocks\n", len(prunableBC.Blocks))
	// Initialize state pruner with policy
	// Keep the last 10 blocks, use height-based checkpoints
	statePruner := core.NewStatePruner(5, 10, true)
	statePruner.DemonstrateStatePruning(prunableBC)
}

// === Helper: Adaptive CAP Simulation ===
func demoAdaptiveCAP(acm *core.AdaptiveCapacityManager) {
	metrics1 := core.NetworkMetrics{
		Latency:    100 * time.Millisecond,
		Throughput: 500.0,
		ErrorRate:  0.01,
		NodeID:     "node1",
		Timestamp:  time.Now(),
	}
	metrics2 := core.NetworkMetrics{
		Latency:    200 * time.Millisecond,
		Throughput: 300.0,
		ErrorRate:  0.05,
		NodeID:     "node2",
		Timestamp:  time.Now(),
	}

	acm.RecordMetrics(metrics1)
	acm.RecordMetrics(metrics2)

	fmt.Println("Node Capacities:")
	fmt.Println("- Node1:", acm.GetNodeCapacity("node1"))
	fmt.Println("- Node2:", acm.GetNodeCapacity("node2"))

	// Simulate degraded performance on node1
	metrics1.Latency = 300 * time.Millisecond
	metrics1.ErrorRate = 0.1
	metrics1.Timestamp = time.Now()
	acm.RecordMetrics(metrics1)

	fmt.Println("\nCapacities after network degradation:")
	fmt.Println("- Node1:", acm.GetNodeCapacity("node1"))
	fmt.Println("- Node2:", acm.GetNodeCapacity("node2"))

	// Global view
	fmt.Println("\nGlobal network view:")
	globalView := acm.GetGlobalView()
	for nodeID, capacity := range globalView {
		fmt.Printf("- %s: %.2f\n", nodeID, capacity)
	}
	// === Consistency Orchestration Simulation ===
	orch := core.NewOrchestrator()

	// Simulate varying network conditions
	fmt.Println("\nEvaluating network conditions for consistency adjustment...")
	orch.EvaluateNetwork(80*time.Millisecond, 0.01) // Strong
	orch.PrintStatus()

	orch.EvaluateNetwork(150*time.Millisecond, 0.04) // Causal
	orch.PrintStatus()

	orch.EvaluateNetwork(300*time.Millisecond, 0.09) // Eventual
	orch.PrintStatus()
}
