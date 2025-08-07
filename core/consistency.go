package core

import (
	"fmt"
	"time"
)

type ConsistencyLevel string

const (
	Strong   ConsistencyLevel = "Strong"
	Causal   ConsistencyLevel = "Causal"
	Eventual ConsistencyLevel = "Eventual"
)

type ConsistencyOrchestrator struct {
	CurrentLevel ConsistencyLevel
	LastLatency  time.Duration
	ErrorRate    float64
}

func NewOrchestrator() *ConsistencyOrchestrator {
	return &ConsistencyOrchestrator{
		CurrentLevel: Strong, // default
	}
}

// Simulate monitoring: adjusts based on latency/error rate
func (co *ConsistencyOrchestrator) EvaluateNetwork(latency time.Duration, errorRate float64) {
	co.LastLatency = latency
	co.ErrorRate = errorRate

	if latency > 250*time.Millisecond || errorRate > 0.08 {
		co.CurrentLevel = Eventual
	} else if latency > 100*time.Millisecond || errorRate > 0.03 {
		co.CurrentLevel = Causal
	} else {
		co.CurrentLevel = Strong
	}
}

func (co *ConsistencyOrchestrator) PrintStatus() {
	fmt.Println("=== Consistency Orchestrator ===")
	fmt.Printf("Current Level: %s\n", co.CurrentLevel)
	fmt.Printf("Last Latency: %s\n", co.LastLatency)
	fmt.Printf("Last Error Rate: %.2f\n", co.ErrorRate)
}
