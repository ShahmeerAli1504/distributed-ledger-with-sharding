package core

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
)

// ProofCompressingMerkleTree is a Merkle Tree with probabilistic verification
type ProofCompressingMerkleTree struct {
	tree *MerkleTree
}

// NewProofCompressingMerkleTree creates a new proof-compressing Merkle Tree
func NewProofCompressingMerkleTree(data []string) *ProofCompressingMerkleTree {
	return &ProofCompressingMerkleTree{
		tree: NewMerkleTree(data),
	}
}

// GetRootHash returns the root hash of the Merkle Tree
func (pcmt *ProofCompressingMerkleTree) GetRootHash() string {
	return pcmt.tree.GetRootHash()
}

// VerifyDataProbabilistic verifies data membership probabilistically
func (pcmt *ProofCompressingMerkleTree) VerifyDataProbabilistic(data string) bool {
	// Hash the data to check against leaves
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	// Check if the hash is in the leaves
	for _, leaf := range pcmt.tree.Leaves {
		if leaf == hashStr {
			return true
		}
	}
	return false
}

// BloomFilter is a probabilistic data structure for membership testing
type BloomFilter struct {
	bits      []bool
	size      uint
	hashFuncs uint
}

// NewBloomFilter creates a new Bloom Filter
func NewBloomFilter(size, hashFuncs uint) *BloomFilter {
	return &BloomFilter{
		bits:      make([]bool, size),
		size:      size,
		hashFuncs: hashFuncs,
	}
}

// Add adds an item to the Bloom Filter
func (bf *BloomFilter) Add(item string) {
	for i := uint(0); i < bf.hashFuncs; i++ {
		hash := simpleHash(item, i)
		index := hash % uint64(bf.size)
		bf.bits[index] = true
	}
}

// Test checks if an item is possibly in the Bloom Filter
func (bf *BloomFilter) Test(item string) bool {
	for i := uint(0); i < bf.hashFuncs; i++ {
		hash := simpleHash(item, i)
		index := hash % uint64(bf.size)
		if !bf.bits[index] {
			return false
		}
	}
	return true
}

// FalsePositiveRate estimates the false positive rate
func (bf *BloomFilter) FalsePositiveRate() float64 {
	// Approximate false positive rate: (1 - e^(-k*n/m))^k
	// where k = hashFuncs, n = number of items (approximated), m = size
	n := float64(bf.size) / 10 // Assume ~10% of bits are set for estimation
	k := float64(bf.hashFuncs)
	m := float64(bf.size)
	return math.Pow(1.0-math.Exp(-k*n/m), k)
}

// simpleHash generates a hash for Bloom Filter
func simpleHash(data string, seed uint) uint64 {
	hash := sha256.Sum256([]byte(data + string(seed)))
	var result uint64
	for i := 0; i < 8; i++ {
		result = (result << 8) | uint64(hash[i])
	}
	return result
}
