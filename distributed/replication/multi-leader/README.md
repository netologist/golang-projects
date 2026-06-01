# Multi-Leader Replication

## Concept
Multiple leaders accept writes concurrently. Each leader propagates its
writes to other leaders asynchronously. Conflicts are inevitable when two
leaders write to the same key concurrently — a conflict resolver (e.g.,
Last-Write-Wins) picks the winner on sync.

## When to Use
- You need multi-region writes with low latency (write locally, sync globally).
- You can tolerate occasional conflicts resolved automatically.

## When NOT to Use
- Conflicts are unacceptable — use single-primary with strong consistency.
- Your conflict resolution logic is too complex to automate (e.g., merging text documents).

## Trade-offs
| Benefit | Cost |
|---------|------|
| Low-latency writes in every region | Write conflicts that need resolution |
| Survives single-region outages | More complex than primary-replica |

## Go-Specific Notes
Leaders maintain independent stores. `Sync()` merges all keys from the
remote leader. For conflicting keys, the resolver (LWW by timestamp) picks
the winner. Conflict details (local, remote, winner) are returned.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Multi-leader = write-anywhere with async cross-region sync.
- Conflicts are expected — the resolver picks a deterministic winner.
- LWW is the simplest resolver; CRDTs or application logic handle complex conflicts.
