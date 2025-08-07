package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ZKProver simulates creating a zero-knowledge proof
type ZKProver struct{}

// ProveKnowledge simulates proving knowledge of a value without revealing it
func (zk *ZKProver) ProveKnowledge(secret string) string {
	// In a real ZKP, you'd generate a proof without exposing the secret
	// Here we hash the secret to simulate a commitment
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

// VerifyProof checks if a given proof matches the actual secret
func (zk *ZKProver) VerifyProof(proof string, claimed string) bool {
	expected := zk.ProveKnowledge(claimed)
	return expected == proof
}

// TestZKP simulates proving and verifying knowledge
func (zk *ZKProver) TestZKP() {
	fmt.Println("\nSimulating Zero-Knowledge Proof...")

	secret := "SuperSecretTransaction"
	proof := zk.ProveKnowledge(secret)

	// Simulate attacker not knowing secret
	wrong := "FakeData"

	fmt.Println("Original Proof:", proof)

	valid := zk.VerifyProof(proof, secret)
	invalid := zk.VerifyProof(proof, wrong)

	fmt.Println("Proof Valid?", valid)
	fmt.Println("Proof Valid with Wrong Data?", invalid)
}
