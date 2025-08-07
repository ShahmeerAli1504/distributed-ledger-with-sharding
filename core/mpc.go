package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// MPCProtocol represents a basic multi-party computation protocol
type MPCProtocol struct {
	Participants []*Node
	Threshold    int
	SecretShares map[int]*big.Int // Node ID -> Secret share
	Prime        *big.Int         // Prime modulus for Shamir's secret sharing
}

// NewMPCProtocol creates a new MPC protocol instance
func NewMPCProtocol(participants []*Node, threshold int) *MPCProtocol {
	// Generate a large prime for the field
	prime, _ := rand.Prime(rand.Reader, 256)
	
	return &MPCProtocol{
		Participants: participants,
		Threshold:    threshold,
		SecretShares: make(map[int]*big.Int),
		Prime:        prime,
	}
}

// generatePolynomial creates a random polynomial f(x) of degree threshold-1
// such that f(0) = secret
func (mpc *MPCProtocol) generatePolynomial(secret *big.Int, degree int) []*big.Int {
	coefficients := make([]*big.Int, degree+1)
	coefficients[0] = new(big.Int).Set(secret)
	
	for i := 1; i <= degree; i++ {
		// Random coefficient
		coeff, _ := rand.Int(rand.Reader, mpc.Prime)
		coefficients[i] = coeff
	}
	
	return coefficients
}

// evaluatePolynomial calculates f(x) for a polynomial
func (mpc *MPCProtocol) evaluatePolynomial(coefficients []*big.Int, x int) *big.Int {
	result := big.NewInt(0)
	xBig := big.NewInt(int64(x))
	
	// Compute: a_0 + a_1 * x + a_2 * x^2 + ... + a_n * x^n
	for i := 0; i < len(coefficients); i++ {
		term := new(big.Int)
		term.Exp(xBig, big.NewInt(int64(i)), mpc.Prime)
		term.Mul(term, coefficients[i])
		term.Mod(term, mpc.Prime)
		
		result.Add(result, term)
		result.Mod(result, mpc.Prime)
	}
	
	return result
}

// ShareSecret splits a secret using Shamir's Secret Sharing
func (mpc *MPCProtocol) ShareSecret(secret *big.Int) {
	// Clear previous shares
	mpc.SecretShares = make(map[int]*big.Int)
	
	// Create a polynomial with our secret as the constant term
	coeffs := mpc.generatePolynomial(secret, mpc.Threshold-1)
	
	// Generate a share for each participant
	for _, participant := range mpc.Participants {
		// Evaluate f(participant.ID)
		share := mpc.evaluatePolynomial(coeffs, participant.ID+1)
		mpc.SecretShares[participant.ID] = share
	}
}

// ReconstructSecret uses Lagrange interpolation to reconstruct the secret
func (mpc *MPCProtocol) ReconstructSecret(shares map[int]*big.Int) (*big.Int, error) {
	if len(shares) < mpc.Threshold {
		return nil, fmt.Errorf("not enough shares: need %d, have %d", mpc.Threshold, len(shares))
	}
	
	// Use first 'threshold' shares
	points := make([]struct{ x, y *big.Int }, 0, mpc.Threshold)
	for id, share := range shares {
		if len(points) >= mpc.Threshold {
			break
		}
		points = append(points, struct{ x, y *big.Int }{
			x: big.NewInt(int64(id + 1)),
			y: new(big.Int).Set(share),
		})
	}
	
	// Compute the secret using Lagrange interpolation at x=0
	secret := big.NewInt(0)
	
	for i, p := range points {
		// Compute Lagrange basis polynomial l_i(0)
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		
		for j, p2 := range points {
			if i == j {
				continue
			}
			
			// 0 - x_j
			term := new(big.Int).Neg(p2.x)
			term.Mod(term, mpc.Prime)
			numerator.Mul(numerator, term)
			numerator.Mod(numerator, mpc.Prime)
			
			// x_i - x_j
			term = new(big.Int).Sub(p.x, p2.x)
			term.Mod(term, mpc.Prime)
			denominator.Mul(denominator, term)
			denominator.Mod(denominator, mpc.Prime)
		}
		
		// Calculate the inverse of denominator
		inverseDenom := new(big.Int).ModInverse(denominator, mpc.Prime)
		
		// l_i(0) = numerator / denominator
		basis := new(big.Int).Mul(numerator, inverseDenom)
		basis.Mod(basis, mpc.Prime)
		
		// Add l_i(0) * y_i to result
		term := new(big.Int).Mul(basis, p.y)
		term.Mod(term, mpc.Prime)
		secret.Add(secret, term)
		secret.Mod(secret, mpc.Prime)
	}
	
	return secret, nil
}

// SimulateMPCSignature demonstrates threshold signing
func (mpc *MPCProtocol) SimulateMPCSignature(message string) string {
	fmt.Println("\nSimulating Multi-Party Computation Threshold Signing...")
	
	// Generate a random secret key
	secretKey, _ := rand.Int(rand.Reader, big.NewInt(1).Exp(big.NewInt(2), big.NewInt(128), nil))
	fmt.Println("Original Secret Key:", secretKey.String())
	
	// Share the secret among participants
	mpc.ShareSecret(secretKey)
	fmt.Printf("Secret shared among %d participants with threshold %d\n", 
		len(mpc.Participants), mpc.Threshold)
	
	// Simulate some participants coming together to sign
	collectedShares := make(map[int]*big.Int)
	participantCount := 0
	
	for _, node := range mpc.Participants {
		if !node.Byzantine && participantCount < mpc.Threshold+1 {
			share := mpc.SecretShares[node.ID]
			collectedShares[node.ID] = share
			participantCount++
		}
	}
	
	fmt.Printf("Collected %d shares from honest participants\n", len(collectedShares))
	
	// Reconstruct the secret
	reconstructed, err := mpc.ReconstructSecret(collectedShares)
	if err != nil {
		fmt.Println("Error reconstructing secret:", err)
		return ""
	}
	
	fmt.Println("Reconstructed Secret Key:", reconstructed.String())
	
	// Use the reconstructed key to "sign" the message
	h := sha256.New()
	h.Write([]byte(message))
	h.Write(reconstructed.Bytes())
	signature := hex.EncodeToString(h.Sum(nil))
	
	fmt.Println("Threshold Signature:", signature)
	return signature
}