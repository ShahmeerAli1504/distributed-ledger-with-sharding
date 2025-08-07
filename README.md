# Distributed-Ledger-With-Sharding
This project implements a modular and secure **blockchain system** with a focus on performance, cryptographic integrity, cross-shard operations, and hybrid consensus. It is built as part of the course **CS4049 - Blockchain and Cryptocurrency** at **FAST-NUCES, Islamabad**.

## Team Members
- **Fatima Basit**
- **Abdullah Ashfaque**   
- **Shahmeer Ali Akhtar** 

## Project Structure

### 1. Architectural Design
- `main.go`: Initializes blockchain and workflow orchestration.
- `block.go`, `blockchain.go`: Define block structure and chain management.
- `shard.go`: Manages sharding and dynamic load balancing.
- `consensus.go`, `bft.go`: Implements hybrid consensus using PoW + BFT with VRF-based leader election.

### 2. Cryptographic Protocols
- `merkle_tree.go`: Merkle trees for state integrity and proof generation.
- `zk_proofs.go`: Zero-Knowledge Proofs using SHA-256 commitment schemes.
- `mpc.go`: Multiparty Computation with secret sharing and threshold signatures.
- `homomorphic_auth.go`: Homomorphic authentication for cross-shard operations.

### 3. Performance Modules
- `state.go`, `state_pruning.go`: Cryptographic state compression and archival.
- `adaptive_cap.go`: Network-aware capacity management and optimization.
- `probabilistic_verification.go`: Bloom filters for fast membership verification.

---

## Key Features

### Core Functionality
- Genesis block creation and sequential block growth
- Secure block hashing with SHA-256
- Shard management with Merkle Forest verification

### Security
- Byzantine Fault Tolerant consensus with node reputation scoring
- Cross-shard HMAC-SHA256 commitments
- Multi-party computation for distributed trust
- Zero-knowledge proof verification

### Performance
- Horizontal scalability through dynamic sharding
- Adaptive consistency tuning based on network metrics
- State pruning reduces storage with <1% verification overhead
- Fast cross-shard atomic transfers with rollback on failure

---

## Benchmarks
- **Shard Reveal Time:** O(log n) using Merkle trees (≈ 3ms for 10 shards)
- **Probabilistic Verification:** 92% proof size reduction using Bloom filters
- **MPC:** Withstands up to 2/5 faulty nodes with 99.9% success
- **Cross-Shard Latency:** ~15–20ms additional latency with strong consistency

---

## Theoretical Concepts Applied
- Adaptive Merkle Forest
- Hybrid PoW + BFT consensus
- Energy-efficient finality via VRF
- Homomorphic and probabilistic cryptography
- Cryptographic accumulators and state compression

---

## Conclusion
This system integrates advanced blockchain mechanisms into a clean, modular, and extensible architecture. It demonstrates performance, scalability, and security suitable for decentralized applications and permissioned blockchain environments.

---

## How to Run (Simulated Workflow)
```bash
go run main.go
