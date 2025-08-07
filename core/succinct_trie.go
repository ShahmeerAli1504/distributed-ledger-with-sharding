package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// TrieNode represents a node in the succinct trie
type TrieNode struct {
	Children map[byte]*TrieNode // Map of child nodes by byte
	Value    string             // Stored value (block data) for leaf nodes
	Hash     string             // Hash of the node (for Merkle root)
}

// SuccinctTrie represents a compact state trie
type SuccinctTrie struct {
	Root *TrieNode
}

// NewSuccinctTrie creates a new succinct trie
func NewSuccinctTrie() *SuccinctTrie {
	return &SuccinctTrie{
		Root: &TrieNode{
			Children: make(map[byte]*TrieNode),
			Value:    "",
		},
	}
}

// Insert adds a key-value pair (block hash, block data) to the trie
func (st *SuccinctTrie) Insert(key, value string) {
	current := st.Root
	keyBytes := []byte(key)

	// Traverse/create path based on key bytes
	for _, b := range keyBytes {
		if _, exists := current.Children[b]; !exists {
			current.Children[b] = &TrieNode{
				Children: make(map[byte]*TrieNode),
				Value:    "",
			}
		}
		current = current.Children[b]
	}

	// Store value at leaf and compute hash
	current.Value = value
	current.Hash = st.computeNodeHash(current)
	st.updateHashes(st.Root)
}

// Get retrieves the value associated with a key
func (st *SuccinctTrie) Get(key string) (string, bool) {
	current := st.Root
	keyBytes := []byte(key)

	for _, b := range keyBytes {
		if next, exists := current.Children[b]; exists {
			current = next
		} else {
			return "", false
		}
	}

	if current.Value != "" {
		return current.Value, true
	}
	return "", false
}

// GetMerkleRoot returns the Merkle root of the trie
func (st *SuccinctTrie) GetMerkleRoot() string {
	return st.Root.Hash
}

// computeNodeHash calculates the hash of a node
func (st *SuccinctTrie) computeNodeHash(node *TrieNode) string {
	if node == nil {
		return ""
	}

	data := node.Value
	for b, child := range node.Children {
		data += string(b) + child.Hash
	}

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// updateHashes recomputes hashes bottom-up after insertion
func (st *SuccinctTrie) updateHashes(node *TrieNode) {
	if node == nil {
		return
	}

	for _, child := range node.Children {
		st.updateHashes(child)
	}

	node.Hash = st.computeNodeHash(node)
}

// PrintTrie displays the trie structure (for debugging)
func (st *SuccinctTrie) PrintTrie() {
	fmt.Println("\n--- Succinct Trie State ---")
	st.printNode(st.Root, 0)
}

func (st *SuccinctTrie) printNode(node *TrieNode, level int) {
	if node == nil {
		return
	}

	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "  "
	}

	if node.Value != "" {
		fmt.Printf("%sLeaf: Value=%s, Hash=%s\n", prefix, node.Value, node.Hash)
	} else {
		fmt.Printf("%sNode: Hash=%s\n", prefix, node.Hash)
	}

	for b, child := range node.Children {
		fmt.Printf("%sChild [%c]:\n", prefix+"  ", b)
		st.printNode(child, level+1)
	}
}
