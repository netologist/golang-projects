# CRDT (Conflict-Free Replicated Data Types)

## Concept
Data structures that achieve eventual consistency without coordination.
Operations commute, so any replica can accept writes independently.
Merging replicas always converges to the same value regardless of order.
Implements G-Counter, PN-Counter, LWW-Register, and 2P-Set.

## When to Use
- Replicas must accept writes independently and merge later.
- You want strong eventual consistency without conflict resolution logic.

## When NOT to Use
- You need linearisable consistency or transactions spanning multiple keys.
- Your data model doesn't map to available CRDT types.

## Trade-offs
| Benefit | Cost |
|---------|------|
| No coordination needed for writes | Limited set of data types |
| Automatic conflict resolution | Storage overhead for metadata (e.g., per-replica counters) |

## Go-Specific Notes
Each CRDT implements a `Merge()` method. The G-Counter tracks per-replica
counts; merging takes the max per replica. The LWW-Register uses timestamps
and replica IDs to pick the winner on conflict.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- CRDTs guarantee convergence: merge is commutative, associative, and idempotent.
- G-Counter = increment-only. PN-Counter = increments + decrements via two G-Counters.
- LWW-Register = last-write-wins using timestamps. 2P-Set = add-then-remove semantics.
