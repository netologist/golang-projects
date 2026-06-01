# Raft

## Concept
Leader-based consensus algorithm designed for understandability. A leader
manages a replicated log; followers replicate entries. Leader election happens
via randomised timeouts. Once a leader is elected, it serves all client requests
and replicates log entries to followers.

## When to Use
- You need a production-grade consensus protocol that's easier to reason about than Paxos.
- You're building a replicated state machine or distributed key-value store.

## When NOT to Use
- A single node suffices — consensus overhead is unnecessary.
- The system is read-heavy and can tolerate eventual consistency.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Simple mental model with strong leadership | Leader is a bottleneck for writes |
| Log matching ensures safety | Network partitions can delay leader election |

## Go-Specific Notes
Nodes communicate via Go channels. Each node has a goroutine running
a ticker for election timeouts and heartbeat intervals.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Leader election uses randomised timeouts to prevent split votes.
- Raft guarantees linearisable semantics through its replicated log.
- The leader handles all writes; followers serve reads in the simple version.
