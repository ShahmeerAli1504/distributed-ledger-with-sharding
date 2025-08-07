package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// HomomorphicCommitment represents a commitment to some data
type HomomorphicCommitment struct {
	Value      string
	Commitment string
}

// HomomorphicAuthenticator handles homomorphic authentication
type HomomorphicAuthenticator struct {
	key []byte
}

// NewHomomorphicAuthenticator creates a new authenticator
func NewHomomorphicAuthenticator(key string) *HomomorphicAuthenticator {
	return &HomomorphicAuthenticator{
		key: []byte(key),
	}
}

// AuthenticateData creates a commitment for some data
func (ha *HomomorphicAuthenticator) AuthenticateData(data string) string {
	mac := hmac.New(sha256.New, ha.key)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyAuthentication checks if a commitment is valid
func (ha *HomomorphicAuthenticator) VerifyAuthentication(data, commitment string) bool {
	expected := ha.AuthenticateData(data)
	return hmac.Equal([]byte(commitment), []byte(expected))
}

// CombineCommitments homomorphically combines multiple commitments
func (ha *HomomorphicAuthenticator) CombineCommitments(commitments []HomomorphicCommitment) HomomorphicCommitment {
	combinedValue := ""
	combinedData := ""
	for _, c := range commitments {
		combinedValue += c.Value
		combinedData += c.Commitment
	}

	// Create a new commitment for the combined data
	combinedCommitment := ha.AuthenticateData(combinedData)

	return HomomorphicCommitment{
		Value:      combinedValue,
		Commitment: combinedCommitment,
	}
}

// EnhancedSyncManager extends SyncManager with homomorphic authentication and atomic transfers
type EnhancedSyncManager struct {
	syncManager      *SyncManager
	authenticator    *HomomorphicAuthenticator
	pendingTransfers map[string]*TransferState // Track pending transfers for 2PC
	mutex            sync.Mutex
}

// TransferState represents the state of a pending transfer
type TransferState struct {
	SourceShard    *Shard
	DestShard      *Shard
	BlockIndex     int
	Commitment     string
	Prepared       bool
	SourceSnapshot []Block // Snapshot for rollback
	DestSnapshot   []Block // Snapshot for rollback
}

// NewEnhancedSyncManager creates a new EnhancedSyncManager
func NewEnhancedSyncManager(key string) *EnhancedSyncManager {
	return &EnhancedSyncManager{
		syncManager:      NewSyncManager(),
		authenticator:    NewHomomorphicAuthenticator(key),
		pendingTransfers: make(map[string]*TransferState),
	}
}

// CreateAuthenticatedTransfer initiates a two-phase commit transfer
func (esm *EnhancedSyncManager) CreateAuthenticatedTransfer(source, destination *Shard, blockIndex int) bool {
	esm.mutex.Lock()
	defer esm.mutex.Unlock()

	// Validate block index
	if blockIndex < 0 || blockIndex >= len(source.Blocks) {
		fmt.Printf("Invalid block index %d for Shard #%d (block count: %d)\n", blockIndex, source.ID, len(source.Blocks))
		return false
	}

	// Create partial state (block hash and data)
	block := source.Blocks[blockIndex]
	partialState := fmt.Sprintf("%s:%s", block.Hash, block.Data)
	commitment := esm.authenticator.AuthenticateData(partialState)

	// Store pending transfer state
	transferID := fmt.Sprintf("%d-%d-%d", source.ID, destination.ID, blockIndex)
	transferState := &TransferState{
		SourceShard:    source,
		DestShard:      destination,
		BlockIndex:     blockIndex,
		Commitment:     commitment,
		Prepared:       false,
		SourceSnapshot: make([]Block, len(source.Blocks)),
		DestSnapshot:   make([]Block, len(destination.Blocks)),
	}
	copy(transferState.SourceSnapshot, source.Blocks)
	copy(transferState.DestSnapshot, destination.Blocks)
	esm.pendingTransfers[transferID] = transferState

	// Phase 1: Prepare (lock resources and validate)
	if !esm.prepareTransfer(transferState) {
		delete(esm.pendingTransfers, transferID)
		return false
	}

	return true
}

// prepareTransfer locks resources and validates the transfer
func (esm *EnhancedSyncManager) prepareTransfer(state *TransferState) bool {
	// Validate commitment
	partialState := fmt.Sprintf("%s:%s", state.SourceShard.Blocks[state.BlockIndex].Hash, state.SourceShard.Blocks[state.BlockIndex].Data)
	if !esm.authenticator.VerifyAuthentication(partialState, state.Commitment) {
		fmt.Printf("Prepare failed: Invalid commitment for transfer from Shard #%d to #%d\n", state.SourceShard.ID, state.DestShard.ID)
		return false
	}

	// Simulate resource locking (e.g., shard mutexes)
	state.SourceShard.mutex.Lock()
	defer state.SourceShard.mutex.Unlock()
	state.DestShard.mutex.Lock()
	defer state.DestShard.mutex.Unlock()

	// Mark as prepared
	state.Prepared = true
	return true
}

// VerifyAndApplyTransfer completes or rolls back the transfer
func (esm *EnhancedSyncManager) VerifyAndApplyTransfer(source, destination *Shard, blockIndex int) bool {
	esm.mutex.Lock()
	defer esm.mutex.Unlock()

	transferID := fmt.Sprintf("%d-%d-%d", source.ID, destination.ID, blockIndex)
	transferState, exists := esm.pendingTransfers[transferID]
	if !exists || !transferState.Prepared {
		fmt.Printf("Transfer %s not found or not prepared\n", transferID)
		return false
	}

	// Validate block index
	if blockIndex < 0 || blockIndex >= len(source.Blocks) {
		fmt.Printf("Invalid block index %d for Shard #%d (block count: %d) during commit\n", blockIndex, source.ID, len(source.Blocks))
		return false
	}

	// Phase 2: Commit or Rollback
	if esm.authenticator.VerifyAuthentication(
		fmt.Sprintf("%s:%s", source.Blocks[blockIndex].Hash, source.Blocks[blockIndex].Data),
		transferState.Commitment) {
		// Commit: Apply transfer
		if esm.syncManager.SyncBlock(source, destination, blockIndex) {
			fmt.Printf("Committed transfer from Shard #%d to #%d\n", source.ID, destination.ID)
			delete(esm.pendingTransfers, transferID)
			return true
		}
	}

	// Rollback: Restore snapshots
	source.Blocks = make([]Block, len(transferState.SourceSnapshot))
	destination.Blocks = make([]Block, len(transferState.DestSnapshot))
	copy(source.Blocks, transferState.SourceSnapshot)
	copy(destination.Blocks, transferState.DestSnapshot)

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

	fmt.Printf("Rolled back transfer from Shard #%d to #%d\n", source.ID, destination.ID)
	delete(esm.pendingTransfers, transferID)
	return false
}
