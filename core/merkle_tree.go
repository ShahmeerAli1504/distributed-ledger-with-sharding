package core

import (
	"crypto/sha256"
	"encoding/hex"
)

// MerkleTree represents a Merkle Tree for efficient data verification
type MerkleTree struct {
	Root   string
	Leaves []string
}

// NewMerkleTree creates a new Merkle Tree from a slice of data
func NewMerkleTree(data []string) *MerkleTree {
	if len(data) == 0 {
		// Return a tree with a default root (hash of empty string)
		hash := sha256.Sum256([]byte(""))
		return &MerkleTree{
			Root:   hex.EncodeToString(hash[:]),
			Leaves: []string{},
		}
	}

	// Create leaf nodes
	leaves := make([]string, len(data))
	for i, d := range data {
		hash := sha256.Sum256([]byte(d))
		leaves[i] = hex.EncodeToString(hash[:])
	}

	// Build the tree
	return &MerkleTree{
		Root:   buildMerkleTree(leaves),
		Leaves: leaves,
	}
}

// buildMerkleTree constructs the Merkle Tree and returns the root hash
func buildMerkleTree(leaves []string) string {
	if len(leaves) == 0 {
		hash := sha256.Sum256([]byte(""))
		return hex.EncodeToString(hash[:])
	}

	if len(leaves) == 1 {
		return leaves[0]
	}

	// Build parent nodes
	var parents []string
	for i := 0; i < len(leaves); i += 2 {
		var combined []byte
		if i+1 < len(leaves) {
			combined = append([]byte(leaves[i]), []byte(leaves[i+1])...)
		} else {
			combined = []byte(leaves[i])
		}
		hash := sha256.Sum256(combined)
		parents = append(parents, hex.EncodeToString(hash[:]))
	}

	// Recursively build the tree
	return buildMerkleTree(parents)
}

// GetRootHash returns the Merkle Tree root
func (mt *MerkleTree) GetRootHash() string {
	return mt.Root
}
