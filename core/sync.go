package core

import (
	"sync"
)

// SyncManager handles basic cross-shard synchronization
type SyncManager struct {
	mutex sync.Mutex
}

// NewSyncManager creates a new SyncManager
func NewSyncManager() *SyncManager {
	return &SyncManager{}
}

// SyncBlock transfers a block between shards
func (sm *SyncManager) SyncBlock(source, destination *Shard, blockIndex int) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if blockIndex < 0 || blockIndex >= len(source.Blocks) {
		return false
	}

	// Transfer block
	block := source.Blocks[blockIndex]
	destination.AddBlock(block)

	// Remove from source
	source.Blocks = append(source.Blocks[:blockIndex], source.Blocks[blockIndex+1:]...)

	// Rebuild Merkle trees
	var sourceData, destData []string
	for _, b := range source.Blocks {
		sourceData = append(sourceData, b.Data)
	}
	for _, b := range destination.Blocks {
		destData = append(destData, b.Data)
	}
	source.Tree = NewMerkleTree(sourceData)
	destination.Tree = NewMerkleTree(destData)

	return true
}
