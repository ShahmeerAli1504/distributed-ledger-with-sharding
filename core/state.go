package core

import (
	"fmt"
)

type ArchivedBlock struct {
	Index int
	Data  string
	Hash  string
}

type StateManager struct {
	PrunedBlocks   []ArchivedBlock
	ActiveBlocks   []Block
	ActiveTrie     *SuccinctTrie // Trie for active blocks
	ArchiveTrie    *SuccinctTrie // Trie for archived blocks
	MaxActiveCount int
}

func NewStateManager(maxActive int) *StateManager {
	return &StateManager{
		PrunedBlocks:   []ArchivedBlock{},
		ActiveBlocks:   []Block{},
		ActiveTrie:     NewSuccinctTrie(),
		ArchiveTrie:    NewSuccinctTrie(),
		MaxActiveCount: maxActive,
	}
}

// AddBlock adds a new block and prunes if limit exceeded
func (sm *StateManager) AddBlock(block Block) {
	sm.ActiveBlocks = append(sm.ActiveBlocks, block)
	// Insert block into active trie (key: block hash, value: block data)
	sm.ActiveTrie.Insert(block.Hash, block.Data)

	if len(sm.ActiveBlocks) > sm.MaxActiveCount {
		archived := sm.ActiveBlocks[0]
		sm.PrunedBlocks = append(sm.PrunedBlocks, ArchivedBlock{
			Index: archived.Index,
			Data:  archived.Data,
			Hash:  archived.Hash,
		})
		// Move to archive trie
		sm.ArchiveTrie.Insert(archived.Hash, archived.Data)
		sm.ActiveBlocks = sm.ActiveBlocks[1:]
		// Remove from active trie (optional, or keep for history)
		// For simplicity, we keep it in active trie but rely on ActiveBlocks for state
	}
}

// UpdateArchiveRoot sets the archive root to the trie’s Merkle root
func (sm *StateManager) UpdateArchiveRoot() {
	sm.ArchiveTrie.updateHashes(sm.ArchiveTrie.Root)
}

// GetActiveRoot returns the Merkle root of active blocks
func (sm *StateManager) GetActiveRoot() string {
	return sm.ActiveTrie.GetMerkleRoot()
}

// GetArchiveRoot returns the Merkle root of archived blocks
func (sm *StateManager) GetArchiveRoot() string {
	return sm.ArchiveTrie.GetMerkleRoot()
}

// PrintState displays current active state and archive status
func (sm *StateManager) PrintState() {
	fmt.Println("\n--- State Manager ---")
	fmt.Printf("Active Blocks: %d\n", len(sm.ActiveBlocks))
	for _, b := range sm.ActiveBlocks {
		fmt.Printf("Block #%d - Hash: %s\n", b.Index, b.Hash)
	}
	fmt.Printf("Active Trie Merkle Root: %s\n", sm.GetActiveRoot())
	fmt.Printf("Archived Blocks: %d\n", len(sm.PrunedBlocks))
	fmt.Printf("Archive Trie Merkle Root: %s\n", sm.GetArchiveRoot())
	// Print trie structures for debugging
	fmt.Println("\nActive Trie Structure:")
	sm.ActiveTrie.PrintTrie()
	fmt.Println("\nArchive Trie Structure:")
	sm.ArchiveTrie.PrintTrie()
}
