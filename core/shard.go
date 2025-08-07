package core

import (
	"fmt"
	"sync"
)

const MinBlocksPerShard = 2
const MaxBlocksPerShard = 3 // Example threshold for demo/testing

type Shard struct {
	ID     int
	Blocks []Block
	Tree   *MerkleTree
	mutex  sync.Mutex
}

type ShardManager struct {
	Shards *RBTree // Use Red-Black Tree for shard storage
}

// NewShard creates a new shard with a unique ID
func NewShard(id int) *Shard {
	return &Shard{
		ID:     id,
		Blocks: []Block{},
	}
}

// AddBlock adds a block to a shard and rebuilds its Merkle tree
func (s *Shard) AddBlock(block Block) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Blocks = append(s.Blocks, block)

	var data []string
	for _, b := range s.Blocks {
		data = append(data, b.Data)
	}
	s.Tree = NewMerkleTree(data)
}

// GetRoot returns the Merkle root of this shard
func (s *Shard) GetRoot() string {
	if s.Tree != nil {
		return s.Tree.GetRootHash()
	}
	return ""
}

// Initialize shard manager
func NewShardManager() *ShardManager {
	tree := NewRBTree()
	shard := NewShard(0)
	tree.Insert(shard) // Start with one shard
	return &ShardManager{
		Shards: tree,
	}
}

// DistributeBlock handles dynamic allocation
func (sm *ShardManager) DistributeBlock(block Block) {
	// Get the last shard (highest ID)
	shards := sm.Shards.GetAllShards()
	lastShard := shards[len(shards)-1]
	lastShard.AddBlock(block)

	// Trigger rebalance if needed
	sm.RebalanceShards()
}

// RebalanceShards splits or keeps shards based on block count
func (sm *ShardManager) RebalanceShards() {
	currentShards := sm.Shards.GetAllShards()
	newTree := NewRBTree()
	shardIDCounter := len(currentShards)

	for _, shard := range currentShards {
		if len(shard.Blocks) > MaxBlocksPerShard {
			// Split this shard
			mid := len(shard.Blocks) / 2
			leftBlocks := shard.Blocks[:mid]
			rightBlocks := shard.Blocks[mid:]

			// Update original shard with left half
			shard.Blocks = leftBlocks
			shard.Tree = NewMerkleTree(getDataStrings(leftBlocks))
			newTree.Insert(shard)

			// Create new shard with right half
			newShard := NewShard(shardIDCounter)
			for _, b := range rightBlocks {
				newShard.AddBlock(b)
			}
			newTree.Insert(newShard)
			shardIDCounter++
		} else {
			newTree.Insert(shard)
		}
	}

	sm.Shards = newTree
}

func getDataStrings(blocks []Block) []string {
	var data []string
	for _, b := range blocks {
		data = append(data, b.Data)
	}
	return data
}

// PrintShardState displays Merkle roots and shard info
func (sm *ShardManager) PrintShardState() {
	fmt.Println("\n==== Shard Merkle Forest ====")
	for _, shard := range sm.Shards.GetAllShards() {
		fmt.Printf("Shard #%d → Root: %s | Blocks: %d\n", shard.ID, shard.GetRoot(), len(shard.Blocks))
	}
	// Print Red-Black Tree structure
	sm.Shards.PrintTree()
}

// MergeShards merges underutilized shards
func (sm *ShardManager) MergeShards(threshold int) {
	currentShards := sm.Shards.GetAllShards()
	newTree := NewRBTree()
	used := make(map[int]bool)

	for i := 0; i < len(currentShards); i++ {
		if used[i] {
			continue
		}

		current := currentShards[i]

		// Try to find another small shard to merge with
		if len(current.Blocks) < threshold && i+1 < len(currentShards) && !used[i+1] {
			next := currentShards[i+1]

			// Merge blocks
			merged := NewShard(current.ID)
			merged.Blocks = append(current.Blocks, next.Blocks...)

			// Rebuild Merkle tree
			var data []string
			for _, b := range merged.Blocks {
				data = append(data, b.Data)
			}
			merged.Tree = NewMerkleTree(data)

			newTree.Insert(merged)
			used[i] = true
			used[i+1] = true

			fmt.Printf("[MERGE] Shard #%d and Shard #%d merged into Shard #%d\n", current.ID, next.ID, current.ID)
		} else {
			// Keep the shard as-is
			newTree.Insert(current)
			used[i] = true
		}
	}

	sm.Shards = newTree
}

// FindShard retrieves a shard by ID in O(log n) time
func (sm *ShardManager) FindShard(id int) (*Shard, bool) {
	return sm.Shards.FindShard(id)
}

// ReconstructState returns the Merkle root of a shard for state verification
func (sm *ShardManager) ReconstructState(shardID int) (string, bool) {
	shard, exists := sm.FindShard(shardID)
	if !exists {
		return "", false
	}
	return shard.GetRoot(), true
}
