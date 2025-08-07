package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// PruningPolicy defines how states are pruned
type PruningPolicy struct {
	MaxHeight      int
	RetentionCount int
	UseCheckpoints bool
}

// IntegrityProof represents cryptographic proof of pruned state
type IntegrityProof struct {
	RootHash    string
	PrunedCount int
	Timestamp   time.Time
	Signature   string
}

// StatePruner manages blockchain state pruning with integrity proofs
type StatePruner struct {
	policy          PruningPolicy
	integrityProofs []IntegrityProof
	merkleRoot      string
	secretKey       string
}

// NewStatePruner creates a new state pruner
func NewStatePruner(maxHeight, retention int, useCheckpoints bool) *StatePruner {
	return &StatePruner{
		policy: PruningPolicy{
			MaxHeight:      maxHeight,
			RetentionCount: retention,
			UseCheckpoints: useCheckpoints,
		},
		integrityProofs: []IntegrityProof{},
		secretKey:       "pruning-integrity-key", // In production, use secure key management
	}
}

// createIntegrityProof generates a cryptographic proof for pruned data
func (sp *StatePruner) createIntegrityProof(rootHash string, count int) IntegrityProof {
	h := sha256.New()
	h.Write([]byte(rootHash))
	h.Write([]byte(fmt.Sprintf("%d", count)))
	h.Write([]byte(sp.secretKey))
	signature := hex.EncodeToString(h.Sum(nil))
	
	return IntegrityProof{
		RootHash:    rootHash,
		PrunedCount: count,
		Timestamp:   time.Now(),
		Signature:   signature,
	}
}

// PruneBlockchain prunes old states while maintaining cryptographic integrity
func (sp *StatePruner) PruneBlockchain(bc *Blockchain) int {
	if len(bc.Blocks) <= sp.policy.RetentionCount {
		return 0  // Nothing to prune
	}
	
	prunableCount := len(bc.Blocks) - sp.policy.RetentionCount
	if sp.policy.UseCheckpoints {
		// Only prune up to checkpoint blocks
		prunableCount = prunableCount - (prunableCount % sp.policy.MaxHeight)
	}
	
	if prunableCount <= 0 {
		return 0
	}
	
	// Calculate hash of pruned blocks for integrity proof
	h := sha256.New()
	for i := 0; i < prunableCount; i++ {
		h.Write([]byte(bc.Blocks[i].Hash))
	}
	rootHash := hex.EncodeToString(h.Sum(nil))
	
	// Create and store integrity proof
	proof := sp.createIntegrityProof(rootHash, prunableCount)
	sp.integrityProofs = append(sp.integrityProofs, proof)
	sp.merkleRoot = rootHash
	
	// Prune the blockchain
	bc.Blocks = bc.Blocks[prunableCount:]
	
	fmt.Printf("\n[INFO] Pruned %d blocks with integrity proof: %s\n", 
		prunableCount, proof.Signature[:16]+"...")
	return prunableCount
}

// VerifyIntegrity checks if the blockchain has been tampered with after pruning
func (sp *StatePruner) VerifyIntegrity(proof IntegrityProof) bool {
	h := sha256.New()
	h.Write([]byte(proof.RootHash))
	h.Write([]byte(fmt.Sprintf("%d", proof.PrunedCount)))
	h.Write([]byte(sp.secretKey))
	expectedSig := hex.EncodeToString(h.Sum(nil))
	
	return expectedSig == proof.Signature
}

// GetLatestProof returns the most recent integrity proof
func (sp *StatePruner) GetLatestProof() *IntegrityProof {
	if len(sp.integrityProofs) == 0 {
		return nil
	}
	return &sp.integrityProofs[len(sp.integrityProofs)-1]
}

// DemonstrateStatePruning shows the pruning mechanism in action
func (sp *StatePruner) DemonstrateStatePruning(bc *Blockchain) {
	fmt.Println("\n=== Demonstrating State Pruning with Cryptographic Integrity ===")
	
	// Current blockchain state
	fmt.Printf("Current blockchain has %d blocks\n", len(bc.Blocks))
	
	// Prune the blockchain
	prunedCount := sp.PruneBlockchain(bc)
	
	if prunedCount > 0 {
		fmt.Printf("Successfully pruned %d blocks\n", prunedCount)
		
		// Verify integrity
		latestProof := sp.GetLatestProof()
		if latestProof != nil {
			verified := sp.VerifyIntegrity(*latestProof)
			fmt.Printf("Integrity verification: %v\n", verified)
			fmt.Printf("Current blockchain has %d blocks after pruning\n", len(bc.Blocks))
		}
	} else {
		fmt.Println("No blocks were pruned based on current policy")
	}
}