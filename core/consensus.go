package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type ConsensusManager struct {
	BFT *BFTManager
}

// simulateProofOfWork injects random "mining" delay
func (cm *ConsensusManager) simulateProofOfWork() string {
	fmt.Println("\n Simulating Proof of Work...")

	// Simulate randomness
	buf := make([]byte, 32)
	rand.Read(buf)
	hash := sha256.Sum256(buf)

	nonce := hex.EncodeToString(hash[:])
	fmt.Println("PoW Nonce:", nonce)

	return nonce
}

// simulateVRFLeaderElection picks one node using deterministic hash
func (cm *ConsensusManager) simulateVRFLeaderElection(nonce string) *Node {
	fmt.Println(" Performing VRF-based Leader Election...")

	var leader *Node
	highestScore := ""

	for _, node := range cm.BFT.Nodes {
		if node.Byzantine {
			continue
		}

		input := fmt.Sprintf("%s-%d", nonce, node.ID)
		hash := sha256.Sum256([]byte(input))
		score := hex.EncodeToString(hash[:])

		if score > highestScore {
			highestScore = score
			leader = node
		}
	}

	if leader != nil {
		fmt.Printf(" Leader Elected: Node #%d\n", leader.ID)
	} else {
		fmt.Println(" No eligible leader found.")
	}

	return leader
}

// RunHybridConsensus executes PoW + BFT + VRF
func (cm *ConsensusManager) RunHybridConsensus() {
	fmt.Println("\n Running Hybrid Consensus Protocol")

	nonce := cm.simulateProofOfWork()
	time.Sleep(1 * time.Second)

	leader := cm.simulateVRFLeaderElection(nonce)
	if leader == nil {
		fmt.Println(" Consensus aborted: No leader.")
		return
	}

	cm.BFT.RunConsensus()
}
