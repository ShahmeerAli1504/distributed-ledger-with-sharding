package core

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	ID           int
	Reputation   float64
	Byzantine    bool
	LastResponse time.Time
}

type BFTManager struct {
	Nodes []*Node
}

// NewBFTManager initializes N nodes
func NewBFTManager(total int) *BFTManager {
	var nodes []*Node
	for i := 0; i < total; i++ {
		nodes = append(nodes, &Node{
			ID:         i,
			Reputation: rand.Float64(),    // simulate history
			Byzantine:  rand.Intn(10) < 2, // ~20% faulty
		})
	}
	return &BFTManager{Nodes: nodes}
}

// SelectConsensusParticipants picks top honest nodes
func (bft *BFTManager) SelectConsensusParticipants() []*Node {
	var selected []*Node
	for _, node := range bft.Nodes {
		if !node.Byzantine && node.Reputation > 0.5 {
			selected = append(selected, node)
		}
	}
	return selected
}

// RunConsensus simulates a voting round
func (bft *BFTManager) RunConsensus() {
	fmt.Println("\nRunning BFT Consensus...")
	participants := bft.SelectConsensusParticipants()

	if len(participants) >= (len(bft.Nodes)*2)/3 {
		fmt.Printf("Consensus Reached with %d honest nodes\n", len(participants))
	} else {
		fmt.Printf("Consensus Failed (only %d honest nodes)\n", len(participants))
	}

	for _, node := range participants {
		fmt.Printf("Node #%d | Reputation: %.2f\n", node.ID, node.Reputation)
	}
}
