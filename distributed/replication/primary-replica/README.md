# Primary-Replica Replication

## Concept
A single primary accepts all writes and propagates them to replicas
asynchronously. Replicas serve reads. If a replica is behind, it catches
up during the propagation interval. Simple and widely used for read-heavy
workloads.

## When to Use
- Your workload is read-heavy and reads can tolerate slight staleness.
- You want simple operational model with a single write authority.

## When NOT to Use
- You need strong read consistency — use quorum reads or read-your-writes.
- You need multi-region writes — use multi-leader or leaderless replication.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Simple mental model (one writer) | Replication lag causes stale reads |
| Linear scaling of read throughput | Primary is a single point of failure for writes |

## Go-Specific Notes
The primary writes locally then propagates to replica nodes via goroutines.
Each replica simulates a configurable delay. Version tracking ensures
propagation can be verified.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Writes flow primary → replicas; reads can come from any replica.
- Asynchronous propagation means replicas may lag.
- Version numbers help detect staleness and monitor replication lag.
