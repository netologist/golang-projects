# Read Repair

## Concept
On each read, check all replicas for the latest version. If any replica
returns stale data, asynchronously push the latest value to bring it up
to date. Combines eventual consistency with incremental repair driven
by client reads — stale data is fixed automatically as it's accessed.

## When to Use
- You use eventual consistency and want stale data to be repaired organically.
- Reads are frequent enough that stale replicas get repaired quickly.

## When NOT to Use
- Stale data must never be returned — use quorum reads instead.
- Reads are rare — unrepaired replicas may stay stale for a long time.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Self-healing on reads | First read of stale data returns correct value but stale replicas still persist briefly |
| No background repair process needed | Adds latency to reads (checking all replicas) |

## Go-Specific Notes
The store checks all replicas on read, picks the max-version value, and
fires a goroutine to repair stale nodes. A repair log tracks which nodes
were fixed with old → new version transitions.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Read repair is the simplest self-healing mechanism for replicated data.
- Repairs are asynchronous — the reader gets the correct value immediately.
- Works well with quorum reads (check N, repair stale, return latest).
