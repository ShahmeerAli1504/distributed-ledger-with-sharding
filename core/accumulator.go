package core

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

// RSAAccumulator represents an RSA-based cryptographic accumulator
type RSAAccumulator struct {
	N        *big.Int            // RSA modulus (product of two safe primes)
	G        *big.Int            // Generator
	State    *big.Int            // Current accumulator state
	Elements []string            // For demo: store elements (in production, store only state)
	Proofs   map[string]*big.Int // Membership proofs for elements
}

// NewRSAAccumulator initializes a new RSA accumulator
func NewRSAAccumulator() *RSAAccumulator {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate two safe primes (simplified for demo; in production, use crypto/rsa)
	p := big.NewInt(251) // Example safe prime
	q := big.NewInt(239) // Example safe prime
	n := new(big.Int).Mul(p, q)
	g := big.NewInt(3) // Example generator

	// Initialize accumulator state as g
	state := new(big.Int).Set(g)

	return &RSAAccumulator{
		N:        n,
		G:        g,
		State:    state,
		Elements: []string{},
		Proofs:   make(map[string]*big.Int),
	}
}

// hashToPrime converts a string to a prime number (simplified for demo)
func (acc *RSAAccumulator) hashToPrime(data string) *big.Int {
	hash := sha256.Sum256([]byte(data))
	hashInt := new(big.Int).SetBytes(hash[:])
	// Ensure it's odd (for simplicity, add 1 if even)
	if hashInt.Bit(0) == 0 {
		hashInt.Add(hashInt, big.NewInt(1))
	}
	// In production, use a primality test (e.g., Miller-Rabin)
	return hashInt
}

// AddElement adds an element to the accumulator
func (acc *RSAAccumulator) AddElement(element string) {
	prime := acc.hashToPrime(element)
	// Update state: state = state^prime mod N
	acc.State.Exp(acc.State, prime, acc.N)
	acc.Elements = append(acc.Elements, element)

	// Generate membership proof: proof = G^(product of all other primes) mod N
	proof := new(big.Int).Set(acc.G)
	for _, e := range acc.Elements {
		if e != element {
			proof.Exp(proof, acc.hashToPrime(e), acc.N)
		}
	}
	acc.Proofs[element] = proof
}

// VerifyMembership checks if an element is in the accumulator
func (acc *RSAAccumulator) VerifyMembership(element string, proof *big.Int) bool {
	prime := acc.hashToPrime(element)
	// Verify: proof^prime mod N == state
	result := new(big.Int).Exp(proof, prime, acc.N)
	return result.Cmp(acc.State) == 0
}

// TestAccumulator demonstrates the RSA accumulator
func (acc *RSAAccumulator) TestAccumulator() {
	fmt.Println("\n=== Testing RSA Cryptographic Accumulator ===")

	// Add some block hashes (simplified as strings)
	elements := []string{
		"block_hash_1",
		"block_hash_2",
		"block_hash_3",
	}

	for _, e := range elements {
		acc.AddElement(e)
		fmt.Printf("Added element: %s\n", e)
	}

	fmt.Printf("Accumulator State: %s\n", acc.State.Text(16))

	// Test membership
	for _, e := range elements {
		proof, exists := acc.Proofs[e]
		if exists {
			valid := acc.VerifyMembership(e, proof)
			fmt.Printf("Membership proof for %s: %v\n", e, valid)
		}
	}

	// Test non-member
	nonMember := "block_hash_4"
	fakeProof := new(big.Int).Set(acc.G) // Invalid proof
	valid := acc.VerifyMembership(nonMember, fakeProof)
	fmt.Printf("Membership proof for %s (non-member): %v\n", nonMember, valid)
}
