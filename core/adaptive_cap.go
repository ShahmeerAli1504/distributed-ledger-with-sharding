package core

import (
	"sync"
	"time"
)

// VectorClock represents a vector clock for tracking causal ordering
type VectorClock struct {
	clock map[string]uint64
	mu    sync.RWMutex
}

// NewVectorClock creates a new vector clock
func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make(map[string]uint64),
	}
}

// Update increments the clock value for the given node
func (vc *VectorClock) Update(nodeID string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.clock[nodeID]++
}

// Merge updates the vector clock with another vector clock
func (vc *VectorClock) Merge(other *VectorClock) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	for nodeID, timestamp := range other.clock {
		if vc.clock[nodeID] < timestamp {
			vc.clock[nodeID] = timestamp
		}
	}
}

// Clone creates a copy of the vector clock
func (vc *VectorClock) Clone() *VectorClock {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	clone := NewVectorClock()
	for nodeID, timestamp := range vc.clock {
		clone.clock[nodeID] = timestamp
	}
	return clone
}

// Get returns the timestamp for a specific node
func (vc *VectorClock) Get(nodeID string) uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.clock[nodeID]
}

// AdaptiveCapacityPolicy defines how nodes adjust capacity based on network conditions
type AdaptiveCapacityPolicy interface {
	AdjustCapacity(metrics NetworkMetrics) float64
}

// NetworkMetrics represents collected metrics about network performance
type NetworkMetrics struct {
	Latency     time.Duration
	Throughput  float64
	ErrorRate   float64
	NodeID      string
	Timestamp   time.Time
	VectorClock *VectorClock
}

// DefaultAdaptivePolicy implements a simple adaptive capacity policy
type DefaultAdaptivePolicy struct {
	baseCapacity  float64
	maxCapacity   float64
	latencyFactor float64
	errorFactor   float64
}

// NewDefaultAdaptivePolicy creates a default policy with reasonable parameters
func NewDefaultAdaptivePolicy() *DefaultAdaptivePolicy {
	return &DefaultAdaptivePolicy{
		baseCapacity:  100.0,
		maxCapacity:   1000.0,
		latencyFactor: 0.5,
		errorFactor:   2.0,
	}
}

// AdjustCapacity implements the AdaptiveCapacityPolicy interface
func (p *DefaultAdaptivePolicy) AdjustCapacity(metrics NetworkMetrics) float64 {
	capacity := p.baseCapacity

	// Reduce capacity when latency increases
	latencyMs := float64(metrics.Latency.Milliseconds())
	latencyAdjustment := p.latencyFactor * latencyMs / 100.0
	capacity -= latencyAdjustment

	// Reduce capacity more aggressively when errors increase
	errorAdjustment := p.errorFactor * metrics.ErrorRate * p.baseCapacity
	capacity -= errorAdjustment

	// Ensure capacity stays within bounds
	if capacity < 0 {
		capacity = 0
	}
	if capacity > p.maxCapacity {
		capacity = p.maxCapacity
	}

	return capacity
}

// AdaptiveCapacityManager manages adaptive capacity for blockchain nodes
type AdaptiveCapacityManager struct {
	nodeCapacities map[string]float64
	nodeLastUpdate map[string]time.Time
	metricHistory  map[string][]NetworkMetrics
	historyLimit   int
	policy         AdaptiveCapacityPolicy
	vectorClock    *VectorClock
	mu             sync.RWMutex
}

// NewAdaptiveCapacityManager creates a new adaptive capacity manager
func NewAdaptiveCapacityManager(nodeID string) *AdaptiveCapacityManager {
	return &AdaptiveCapacityManager{
		nodeCapacities: make(map[string]float64),
		nodeLastUpdate: make(map[string]time.Time),
		metricHistory:  make(map[string][]NetworkMetrics),
		historyLimit:   100,
		policy:         NewDefaultAdaptivePolicy(),
		vectorClock:    NewVectorClock(),
	}
}

// RecordMetrics records new network metrics for a node
func (acm *AdaptiveCapacityManager) RecordMetrics(metrics NetworkMetrics) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	// Update vector clock
	acm.vectorClock.Update(metrics.NodeID)

	// Merge the incoming vector clock
	if metrics.VectorClock != nil {
		acm.vectorClock.Merge(metrics.VectorClock)
	}

	// Add metrics to history
	if _, exists := acm.metricHistory[metrics.NodeID]; !exists {
		acm.metricHistory[metrics.NodeID] = make([]NetworkMetrics, 0)
	}

	acm.metricHistory[metrics.NodeID] = append(acm.metricHistory[metrics.NodeID], metrics)

	// Trim history if needed
	if len(acm.metricHistory[metrics.NodeID]) > acm.historyLimit {
		acm.metricHistory[metrics.NodeID] = acm.metricHistory[metrics.NodeID][1:]
	}

	// Update node capacity
	acm.nodeCapacities[metrics.NodeID] = acm.policy.AdjustCapacity(metrics)
	acm.nodeLastUpdate[metrics.NodeID] = time.Now()
}

// GetNodeCapacity returns the current capacity for a given node
func (acm *AdaptiveCapacityManager) GetNodeCapacity(nodeID string) float64 {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if capacity, exists := acm.nodeCapacities[nodeID]; exists {
		return capacity
	}

	// Return default capacity if node not known
	if policy, ok := acm.policy.(*DefaultAdaptivePolicy); ok {
		return policy.baseCapacity
	}
	return 100.0 // Fallback default
}

// SyncWithPeer syncs capacity information with another node
func (acm *AdaptiveCapacityManager) SyncWithPeer(peerMetrics map[string]NetworkMetrics, peerVC *VectorClock) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	// Merge vector clocks
	if peerVC != nil {
		acm.vectorClock.Merge(peerVC)
	}

	// Process metrics from peer
	for nodeID, metrics := range peerMetrics {
		// Only process metrics that are newer than what we have
		if lastUpdate, exists := acm.nodeLastUpdate[nodeID]; !exists || metrics.Timestamp.After(lastUpdate) {
			// Add to history
			if _, exists := acm.metricHistory[nodeID]; !exists {
				acm.metricHistory[nodeID] = make([]NetworkMetrics, 0)
			}

			acm.metricHistory[nodeID] = append(acm.metricHistory[nodeID], metrics)

			// Trim history if needed
			if len(acm.metricHistory[nodeID]) > acm.historyLimit {
				acm.metricHistory[nodeID] = acm.metricHistory[nodeID][1:]
			}

			// Update capacity
			acm.nodeCapacities[nodeID] = acm.policy.AdjustCapacity(metrics)
			acm.nodeLastUpdate[nodeID] = metrics.Timestamp
		}
	}
}

// GetGlobalView provides a snapshot of the current network state
func (acm *AdaptiveCapacityManager) GetGlobalView() map[string]float64 {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	view := make(map[string]float64)
	for nodeID, capacity := range acm.nodeCapacities {
		view[nodeID] = capacity
	}
	return view
}

// GetVectorClock returns the current vector clock
func (acm *AdaptiveCapacityManager) GetVectorClock() *VectorClock {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	return acm.vectorClock.Clone()
}

// SetPolicy changes the adaptive capacity policy
func (acm *AdaptiveCapacityManager) SetPolicy(policy AdaptiveCapacityPolicy) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.policy = policy
}
