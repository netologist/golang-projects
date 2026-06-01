# Paxos

## Concept
Classic distributed consensus algorithm. Proposers send proposals to acceptors;
a value is chosen when a majority of acceptors agree. Paxos guarantees safety
(no two nodes decide on different values) even with failures.

## When to Use
- You need strong consistency (safety) in an asynchronous distributed system.
- You can tolerate some liveness delay in exchange for correctness.

## When NOT to Use
- You need high throughput — Multi-Paxos or Raft are better optimised.
- You can tolerate eventual consistency — use simpler replication instead.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Proven safety under asynchrony | Complex to understand and implement correctly |
| Tolerates up to (N-1)/2 failures | Two phases per decision; higher latency |

## Go-Specific Notes
The implementation models proposers and acceptors as Go structs. Proposers
drive the protocol; acceptors track the highest-numbered proposal seen.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Paxos guarantees safety (consensus on a single value) even with crash failures.
- The two-phase protocol (Prepare/Accept) is the core mechanism.
- Real-world systems often use Multi-Paxos for log replication.
